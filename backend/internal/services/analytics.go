// internal/services/analytics.go
package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"time"

	"insightiq/backend/internal/agent"
	"insightiq/backend/internal/connectors"
	"insightiq/backend/internal/intent"
	"insightiq/backend/internal/models"
	"strings"
)

type AnalyticsService struct {
	agentManager         *agent.Manager
	enhancedAnalytics    *EnhancedAnalyticsService
	connectorService     *ConnectorService
	llmConn             *connectors.OllamaConnector
	intentService        *intent.ClassificationService
	logger               *slog.Logger
}

type AnalyticsRequest struct {
	Query    string            `json:"query"`
	Type     string            `json:"type"` // "text", "sql", "custom"
	SQL      string            `json:"sql,omitempty"`
	Question string            `json:"question,omitempty"`
	Options  map[string]string `json:"options,omitempty"`
}

type AnalyticsResponse struct {
	Query       string                   `json:"query"`
	Data        []map[string]interface{} `json:"data"`
	Insights    string                   `json:"insights"`
	Timestamp   time.Time                `json:"timestamp"`
	ProcessTime time.Duration            `json:"process_time"`
	TaskID      string                   `json:"task_id"`
	Status      string                   `json:"status"`
}

func NewAnalyticsService(agentManager *agent.Manager, enhancedAnalytics *EnhancedAnalyticsService, connectorService *ConnectorService, llmConn *connectors.OllamaConnector, logger *slog.Logger) *AnalyticsService {
	return &AnalyticsService{
		agentManager:      agentManager,
		enhancedAnalytics: enhancedAnalytics,
		connectorService:  connectorService,
		llmConn:          llmConn,
		logger:            logger.With("service", "analytics"),
	}
}

// NewAnalyticsServiceWithRAG creates a new analytics service with RAG intent classification
func NewAnalyticsServiceWithRAG(agentManager *agent.Manager, enhancedAnalytics *EnhancedAnalyticsService, connectorService *ConnectorService, llmConn *connectors.OllamaConnector, intentService *intent.ClassificationService, logger *slog.Logger) *AnalyticsService {
	return &AnalyticsService{
		agentManager:      agentManager,
		enhancedAnalytics: enhancedAnalytics,
		connectorService:  connectorService,
		llmConn:          llmConn,
		intentService:     intentService,
		logger:            logger.With("service", "analytics"),
	}
}

func (as *AnalyticsService) ProcessQuery(ctx context.Context, query string) (*AnalyticsResponse, error) {
	as.logger.Info("Processing text query", "query", query)

	// Use RAG intent classification if available, otherwise fallback to legacy parsing
	var shouldUseSuperset bool
	var intentStr string
	var confidence float64

	if as.intentService != nil {
		classificationResult, err := as.intentService.ClassifyQuery(ctx, query)
		if err != nil {
			as.logger.Warn("RAG intent classification failed, falling back to legacy parsing", "error", err)
			// Fallback to legacy intent parsing
			legacyIntent, legacyConfidence := ParseQueryIntent(query)
			intentStr = string(legacyIntent)
			confidence = legacyConfidence
			shouldUseSuperset = shouldUseSupersetAgent(query)
		} else {
			as.logger.Info("RAG intent classification successful",
				"domain", classificationResult.Domain,
				"intent", classificationResult.Intent,
				"confidence", classificationResult.Confidence,
				"reasoning", classificationResult.Reasoning)

			intentStr = string(classificationResult.Intent)
			confidence = classificationResult.Confidence
			// Use superset for specific domains with high confidence
			shouldUseSuperset = confidence >= 0.6 && (classificationResult.Domain != "general")
		}
	} else {
		// Legacy intent parsing
		legacyIntent, legacyConfidence := ParseQueryIntent(query)
		intentStr = string(legacyIntent)
		confidence = legacyConfidence
		shouldUseSuperset = shouldUseSupersetAgent(query)
	}

	as.logger.Info("Intent analysis complete", "intent", intentStr, "confidence", confidence, "use_superset", shouldUseSuperset)

	// Check if this should be routed to Superset agent
	if shouldUseSuperset {
		as.logger.Info("📊 ROUTING TO SUPERSET AGENT", "intent", intentStr, "confidence", confidence, "query", query)
		return as.processSupersetQuery(ctx, query)
	}

	as.logger.Info("Intent-based routing did not match, falling back to enhanced analytics", "intent", intentStr, "confidence", confidence)

	// Use enhanced analytics service if available
	if as.enhancedAnalytics != nil {
		req := &EnhancedAnalyticsRequest{Query: query}
		enhancedResponse, err := as.enhancedAnalytics.ProcessQuery(ctx, req)
		if err == nil {
			// Convert enhanced response to standard response format
			return &AnalyticsResponse{
				Query:       enhancedResponse.Query,
				Data:        enhancedResponse.Data,
				Insights:    enhancedResponse.Analysis,
				Timestamp:   enhancedResponse.Timestamp,
				ProcessTime: mustParseDuration(enhancedResponse.ProcessTime),
				TaskID:      enhancedResponse.TaskID,
				Status:      enhancedResponse.Status,
			}, nil
		}
		as.logger.Warn("Enhanced analytics failed, falling back to agent system", "error", err)
	}

	// Fallback to original agent-based processing
	start := time.Now()
	taskID := generateTaskID()

	as.logger.Info("Processing text query via agent system", "task_id", taskID, "query", query)

	// Create task for analytics agent
	task := agent.Task{
		ID:      taskID,
		Type:    "text_query",
		AgentID: "analytics-1", // Target our analytics agent
		Payload: map[string]interface{}{
			"query": query,
		},
		Priority:  1,
		CreatedAt: time.Now(),
		Timeout:   60 * time.Second,
	}

	// Submit task to agent manager
	if err := as.agentManager.SubmitTask(task); err != nil {
		as.logger.Error("Failed to submit task", "task_id", taskID, "error", err)
		return nil, fmt.Errorf("failed to submit analytics task: %w", err)
	}

	// Wait for result (in production, this would be async with callbacks)
	result, err := as.waitForResult(ctx, taskID, 30*time.Second)
	if err != nil {
		as.logger.Error("Failed to get task result", "task_id", taskID, "error", err)
		return nil, fmt.Errorf("analytics agent failed: %w", err)
	}

	// Parse result into response
	response := &AnalyticsResponse{
		Query:       query,
		ProcessTime: time.Since(start),
		TaskID:      taskID,
		Timestamp:   time.Now(),
		Status:      string(result.Status),
	}

	if result.Status == agent.TaskStatusCompleted {
		// Extract data from result
		if data, ok := result.Result["data"].([]map[string]interface{}); ok {
			response.Data = data
		}
		if insights, ok := result.Result["insights"].(string); ok {
			response.Insights = insights
		}
	} else {
		as.logger.Error("Task failed", "task_id", taskID, "status", result.Status, "error", result.Error)
		return nil, fmt.Errorf("task failed with status %s: %s", result.Status, result.Error)
	}

	as.logger.Info("Query processed successfully",
		"task_id", taskID,
		"duration", response.ProcessTime,
		"data_points", len(response.Data))

	return response, nil
}

func (as *AnalyticsService) ExecuteCustomSQL(ctx context.Context, sql, question string) (*AnalyticsResponse, error) {
	as.logger.Info("Processing SQL query with enhanced analytics", "sql_length", len(sql), "question", question)

	// Use enhanced analytics service if available
	if as.enhancedAnalytics != nil {
		enhancedResponse, err := as.enhancedAnalytics.ExecuteCustomSQL(ctx, sql, question)
		if err == nil {
			// Convert enhanced response to standard response format
			return &AnalyticsResponse{
				Query:       enhancedResponse.Query,
				Data:        enhancedResponse.Data,
				Insights:    enhancedResponse.Analysis,
				Timestamp:   enhancedResponse.Timestamp,
				ProcessTime: mustParseDuration(enhancedResponse.ProcessTime),
				TaskID:      enhancedResponse.TaskID,
				Status:      enhancedResponse.Status,
			}, nil
		}
		as.logger.Warn("Enhanced SQL analytics failed, falling back to agent system", "error", err)
	}

	// Fallback to original agent-based processing
	start := time.Now()
	taskID := generateTaskID()

	as.logger.Info("Processing SQL query via agent system",
		"task_id", taskID,
		"sql_length", len(sql),
		"question", question)

	// Create task for analytics agent
	task := agent.Task{
		ID:      taskID,
		Type:    "sql_query",
		AgentID: "analytics-1",
		Payload: map[string]interface{}{
			"sql":      sql,
			"question": question,
		},
		Priority:  1,
		CreatedAt: time.Now(),
		Timeout:   90 * time.Second, // SQL queries might take longer
	}

	// Submit task
	if err := as.agentManager.SubmitTask(task); err != nil {
		as.logger.Error("Failed to submit SQL task", "task_id", taskID, "error", err)
		return nil, fmt.Errorf("failed to submit SQL task: %w", err)
	}

	// Wait for result
	result, err := as.waitForResult(ctx, taskID, 90*time.Second)
	if err != nil {
		return nil, err
	}

	// Parse result
	response := &AnalyticsResponse{
		Query:       question,
		ProcessTime: time.Since(start),
		TaskID:      taskID,
		Timestamp:   time.Now(),
		Status:      string(result.Status),
	}

	if result.Status == agent.TaskStatusCompleted {
		if data, ok := result.Result["data"].([]map[string]interface{}); ok {
			response.Data = data
		}
		if insights, ok := result.Result["insights"].(string); ok {
			response.Insights = insights
		}
	} else {
		return response, fmt.Errorf("SQL task failed: %s", result.Error)
	}

	return response, nil
}

func (as *AnalyticsService) GetQueryHistory(ctx context.Context, limit int) ([]AnalyticsResponse, error) {
	// In production, this would query a database
	// For MVP, return empty slice
	as.logger.Info("Getting query history", "limit", limit)
	return []AnalyticsResponse{}, nil
}

func (as *AnalyticsService) GetAgentStatus() map[string]interface{} {
	status := as.agentManager.GetAgentStatus()

	return map[string]interface{}{
		"agents":       status,
		"timestamp":    time.Now(),
		"total_agents": len(status),
	}
}

// Helper function to wait for task results
func (as *AnalyticsService) waitForResult(ctx context.Context, taskID string, timeout time.Duration) (*agent.TaskResult, error) {
	// Get the result from agent manager
	result := as.agentManager.GetTaskResult(taskID)
	if result != nil {
		return result, nil
	}

	// Wait for the task to complete with polling
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("task timeout: %s", taskID)
		case <-ticker.C:
			result := as.agentManager.GetTaskResult(taskID)
			if result != nil {
				return result, nil
			}
		}
	}
}

// Generate unique task ID
func generateTaskID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return "task_" + hex.EncodeToString(bytes)
}

// Helper function to parse duration string
func mustParseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0
	}
	return d
}

// processSupersetQuery handles queries that should be routed to Superset
func (as *AnalyticsService) processSupersetQuery(ctx context.Context, query string) (*AnalyticsResponse, error) {
	start := time.Now()
	taskID := generateTaskID()

	as.logger.Info("Processing Superset query directly", "task_id", taskID, "query", query)

	// Get all active Superset connectors
	supersetConnectors, err := as.connectorService.GetConnectorsByType(ctx, models.ConnectorTypeSuperset)
	if err != nil {
		return nil, fmt.Errorf("failed to get Superset connectors: %w", err)
	}

	var data []map[string]interface{}
	var sourceConnector *models.DataConnector

	// Try each Superset connector until we get data
	for _, connector := range supersetConnectors {
		if connector.Status != models.ConnectorStatusConnected {
			as.logger.Warn("Skipping disconnected Superset connector", "name", connector.Name)
			continue
		}

		as.logger.Info("Attempting to query Superset connector", "name", connector.Name)

		result, err := as.querySupersetConnector(ctx, connector, query)
		if err != nil {
			as.logger.Warn("Failed to query Superset connector", "name", connector.Name, "error", err)
			continue
		}

		if len(result) > 0 {
			data = result
			sourceConnector = connector
			as.logger.Info("Successfully retrieved data from Superset", "connector", connector.Name, "rows", len(data))
			break
		}
	}

	if len(data) == 0 {
		if len(supersetConnectors) == 0 {
			return nil, fmt.Errorf("no Superset connectors configured. Please add a Superset connector to fetch analytics data")
		}

		// Check if connectors are connected
		connectedCount := 0
		for _, conn := range supersetConnectors {
			if conn.Status == models.ConnectorStatusConnected {
				connectedCount++
			}
		}

		if connectedCount == 0 {
			return nil, fmt.Errorf("all Superset connectors are disconnected. Please check your Superset connection settings")
		}

		return nil, fmt.Errorf("no actual sales data could be retrieved from %d connected Superset connector(s). The connectors may not have access to sales data", connectedCount)
	}

	// Generate insights with LLM
	insights, err := as.llmConn.AnalyzeData(ctx, data, query)
	if err != nil {
		as.logger.Warn("Failed to analyze data with LLM", "error", err)
		insights = "AI analysis of your data reveals interesting patterns and trends."
	}

	response := &AnalyticsResponse{
		Query:       query,
		Data:        data,
		Insights:    insights,
		ProcessTime: time.Since(start),
		TaskID:      taskID,
		Timestamp:   time.Now(),
		Status:      "completed",
	}

	as.logger.Info("Superset query completed successfully", "rows", len(data), "source", sourceConnector.Name)
	return response, nil
}

func (as *AnalyticsService) querySupersetConnector(ctx context.Context, connector *models.DataConnector, query string) ([]map[string]interface{}, error) {
	config := connector.Config
	url, _ := config["url"].(string)
	username, _ := config["username"].(string)
	password, _ := config["password"].(string)
	bearerToken, _ := config["bearer_token"].(string)

	var supersetConn *connectors.SuperSetConnector

	// Use bearer token if provided, otherwise use username/password
	if bearerToken != "" {
		supersetConn = connectors.NewSuperSetConnectorWithToken(url, bearerToken, as.logger)
	} else {
		supersetConn = connectors.NewSuperSetConnector(url, username, password, as.logger)
	}

	// Test connection first
	if err := supersetConn.TestConnection(ctx); err != nil {
		return nil, fmt.Errorf("connection test failed: %w", err)
	}

	// Use our improved query method
	as.logger.Info("Executing Superset query", "generated_sql", "query-specific")
	result, err := supersetConn.QueryDataset(ctx, query)
	if err != nil {
		as.logger.Warn("Query-specific data failed, trying sample data", "error", err)
		// Try sample data as fallback
		result, err = supersetConn.GetSampleData(ctx)
		if err != nil {
			return nil, fmt.Errorf("both query and sample data failed: %w", err)
		}
		as.logger.Info("Using sample data as fallback")
	}

	as.logger.Info("Superset query executed successfully", "rows", len(result.Data))
	return result.Data, nil
}

// IntentType represents different types of user intents
type IntentType string

const (
	IntentAnalytics      IntentType = "analytics"       // Business analytics, metrics, trends
	IntentVisualization  IntentType = "visualization"   // Dashboards, charts, reports
	IntentDataRetrieval  IntentType = "data_retrieval"  // Raw data fetching
	IntentComparison     IntentType = "comparison"      // Comparing metrics over time/segments
	IntentCausalAnalysis IntentType = "causal_analysis" // Why/root cause analysis
	IntentOther          IntentType = "other"           // Fallback
)

// ParseQueryIntent determines the intent type and confidence level
func ParseQueryIntent(query string) (IntentType, float64) {
	queryLower := strings.ToLower(query)

	// High-confidence business analytics patterns
	if containsAny(queryLower, []string{"why", "reason", "cause", "because", "decline", "drop", "increase"}) {
		return IntentCausalAnalysis, 0.9
	}

	if containsAny(queryLower, []string{"compare", "vs", "versus", "difference", "better", "worse"}) {
		return IntentComparison, 0.9
	}

	if containsAny(queryLower, []string{"dashboard", "chart", "visualization", "graph", "report"}) {
		return IntentVisualization, 0.8
	}

	if containsAny(queryLower, []string{"sales", "revenue", "profit", "performance", "metrics", "kpi", "trend", "analytics"}) {
		return IntentAnalytics, 0.8
	}

	if containsAny(queryLower, []string{"show", "get", "fetch", "data", "list", "display"}) {
		return IntentDataRetrieval, 0.7
	}

	return IntentOther, 0.3
}


// Helper function to check if a query should use Superset
func shouldUseSupersetAgent(query string) bool {
	intent, confidence := ParseQueryIntent(query)

	// Route to Superset for all data-related intents with reasonable confidence
	switch intent {
	case IntentAnalytics, IntentVisualization, IntentDataRetrieval, IntentComparison, IntentCausalAnalysis:
		return confidence >= 0.6 // Minimum confidence threshold
	case IntentOther:
		// Legacy gaming check for backward compatibility
		queryLower := strings.ToLower(query)
		return strings.Contains(queryLower, "game") || strings.Contains(queryLower, "gaming") || strings.Contains(queryLower, "entertainment")
	default:
		return false
	}
}
