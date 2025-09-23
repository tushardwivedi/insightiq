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
)

type AnalyticsService struct {
	agentManager *agent.Manager
	logger       *slog.Logger
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

func NewAnalyticsService(agentManager *agent.Manager, logger *slog.Logger) *AnalyticsService {
	return &AnalyticsService{
		agentManager: agentManager,
		logger:       logger.With("service", "analytics"),
	}
}

func (as *AnalyticsService) ProcessQuery(ctx context.Context, query string) (*AnalyticsResponse, error) {
	start := time.Now()
	taskID := generateTaskID()

	as.logger.Info("Processing text query", "task_id", taskID, "query", query)

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
	result, err := as.waitForResult(ctx, taskID, 60*time.Second)
	if err != nil {
		return nil, err
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
		response.Status = "failed"
		return response, fmt.Errorf("task failed: %s", result.Error)
	}

	as.logger.Info("Query processed successfully",
		"task_id", taskID,
		"duration", response.ProcessTime,
		"data_points", len(response.Data))

	return response, nil
}

func (as *AnalyticsService) ExecuteCustomSQL(ctx context.Context, sql, question string) (*AnalyticsResponse, error) {
	start := time.Now()
	taskID := generateTaskID()

	as.logger.Info("Processing SQL query",
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
	// In production, this would be done with channels or callbacks
	// For MVP, we'll simulate waiting (this is a simplified approach)

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("task timeout: %s", taskID)
		case <-ticker.C:
			// In a real implementation, you'd check a result store or channel
			// For now, we'll simulate a completed task after a short delay

			// This is a mock - in production, you'd integrate with the agent manager's result system
			return &agent.TaskResult{
				TaskID:  taskID,
				AgentID: "analytics-1",
				Status:  agent.TaskStatusCompleted,
				Result: map[string]interface{}{
					"data": []map[string]interface{}{
						{"month": "2024-01", "revenue": 50000, "orders": 150},
						{"month": "2024-02", "revenue": 60000, "orders": 180},
						{"month": "2024-03", "revenue": 55000, "orders": 165},
					},
					"insights": "Revenue shows steady growth with Q1 averaging $55K per month. Order volume correlates positively with revenue, suggesting stable customer acquisition.",
				},
				ProcessedAt: time.Now(),
				Duration:    2 * time.Second,
			}, nil
		}
	}
}

// Generate unique task ID
func generateTaskID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return "task_" + hex.EncodeToString(bytes)
}
