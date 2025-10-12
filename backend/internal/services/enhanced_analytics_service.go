package services

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"insightiq/backend/internal/connectors"
	"insightiq/backend/internal/models"
)

// EnhancedAnalyticsService provides intelligent data source routing and RAG capabilities
type EnhancedAnalyticsService struct {
	connectorService *ConnectorService
	plannerService   *PlannerService
	llmConn          *connectors.OllamaConnector
	fallbackPostgres *connectors.PostgresConnector
	fallbackSuperset *connectors.SuperSetConnector
	logger           *slog.Logger
}

// EnhancedAnalyticsRequest represents a request for analytics data
type EnhancedAnalyticsRequest struct {
	Query        string            `json:"query"`
	ConnectorIDs []string          `json:"connector_ids,omitempty"`
	SQL          string            `json:"sql,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// EnhancedAnalyticsResponse represents the response with data and insights
type EnhancedAnalyticsResponse struct {
	Query        string                   `json:"query"`
	Data         []map[string]interface{} `json:"data"`
	Sources      map[string]interface{}   `json:"sources"`
	Analysis     string                   `json:"analysis"`
	DataSources  []string                 `json:"data_sources"`
	Timestamp    time.Time                `json:"timestamp"`
	ProcessTime  string                   `json:"process_time"`
	TaskID       string                   `json:"task_id"`
	Status       string                   `json:"status"`
	Intent       *models.Intent           `json:"intent,omitempty"`
	TaskGraph    *models.TaskGraph        `json:"task_graph,omitempty"`
	PlanningTime string                   `json:"planning_time,omitempty"`
}

func NewEnhancedAnalyticsService(
	connectorService *ConnectorService,
	llm *connectors.OllamaConnector,
	fallbackPostgres *connectors.PostgresConnector,
	fallbackSuperset *connectors.SuperSetConnector,
	logger *slog.Logger,
) *EnhancedAnalyticsService {
	eas := &EnhancedAnalyticsService{
		connectorService: connectorService,
		llmConn:          llm,
		fallbackPostgres: fallbackPostgres,
		fallbackSuperset: fallbackSuperset,
		logger:           logger.With("service", "enhanced_analytics"),
	}

	// Initialize planner service
	eas.plannerService = NewPlannerService(llm, connectorService, logger)

	return eas
}

// ProcessQuery intelligently routes queries to appropriate data sources with RAG
func (eas *EnhancedAnalyticsService) ProcessQuery(ctx context.Context, req *EnhancedAnalyticsRequest) (*EnhancedAnalyticsResponse, error) {
	start := time.Now()
	eas.logger.Info("Processing enhanced analytics query with planner", "query", req.Query)

	// 1. Use Planner LLM for intent parsing and task graph generation
	plannerStart := time.Now()
	plannerReq := &models.PlannerRequest{
		Query:   req.Query,
		Context: map[string]interface{}{
			"connector_ids": req.ConnectorIDs,
			"metadata":      req.Metadata,
		},
	}

	plannerResponse, err := eas.plannerService.ParseIntent(ctx, plannerReq)
	planningTime := time.Since(plannerStart)

	if err != nil {
		eas.logger.Warn("Planner failed, falling back to basic routing", "error", err)
		// Fallback to original logic
		dataSources := eas.analyzeQueryForDataSources(ctx, req.Query, req.ConnectorIDs)
		return eas.processWithBasicRouting(ctx, req, dataSources, start, nil, "")
	}

	eas.logger.Info("Intent parsed successfully",
		"intent_type", plannerResponse.Intent.Type,
		"confidence", plannerResponse.Confidence,
		"steps", len(plannerResponse.TaskGraph.Steps),
		"planning_time", planningTime)

	// 2. Execute task graph based on parsed intent
	dataSources := eas.selectDataSourcesFromIntent(ctx, plannerResponse.Intent, req.ConnectorIDs)
	eas.logger.Info("Data sources selected for query", "count", len(dataSources), "query", req.Query)
	for i, ds := range dataSources {
		eas.logger.Info("Selected data source", "index", i, "name", ds.Name, "type", ds.Type, "status", ds.Status)
	}

	// 2. Retrieve data from multiple sources
	allData := make(map[string]interface{})
	var combinedData []map[string]interface{}

	for _, source := range dataSources {
		eas.logger.Info("Attempting to fetch data from source", "source", source.Name, "type", source.Type)
		data, err := eas.fetchDataFromSource(ctx, source, req.Query)
		if err != nil {
			eas.logger.Error("Failed to fetch from source", "source", source.Name, "type", source.Type, "error", err)
			continue
		}

		if len(data) > 0 {
			allData[source.Name] = data
			combinedData = append(combinedData, data...)
			eas.logger.Info("Retrieved data from source", "source", source.Name, "rows", len(data))
		} else {
			eas.logger.Warn("No data returned from source", "source", source.Name, "type", source.Type)
		}
	}

	// 3. Check if any data was retrieved from connectors
	if len(combinedData) == 0 {
		return nil, fmt.Errorf("no data available from configured connectors. Please check your connector configuration and ensure they contain the requested data")
	}

	// 4. Generate comprehensive analysis with enhanced RAG context using intent
	analysis, err := eas.generateAnalysisWithIntentRAG(ctx, combinedData, allData, req.Query, plannerResponse.Intent)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze data: %w", err)
	}

	return &EnhancedAnalyticsResponse{
		Query:        req.Query,
		Data:         combinedData,
		Sources:      allData,
		Analysis:     analysis,
		DataSources:  eas.getSourceNames(dataSources),
		Timestamp:    time.Now(),
		ProcessTime:  time.Since(start).String(),
		TaskID:       fmt.Sprintf("task_%d", start.Unix()),
		Status:       "completed",
		Intent:       &plannerResponse.Intent,
		TaskGraph:    &plannerResponse.TaskGraph,
		PlanningTime: planningTime.String(),
	}, nil
}

// ExecuteCustomSQL is disabled to prevent direct SQL execution on internal databases
func (eas *EnhancedAnalyticsService) ExecuteCustomSQL(ctx context.Context, sql, question string) (*EnhancedAnalyticsResponse, error) {
	eas.logger.Info("Custom SQL execution disabled - use configured external connectors only")

	// Custom SQL execution disabled to prevent direct database access
	return nil, fmt.Errorf("custom SQL execution is disabled. Please use the /api/query endpoint with natural language queries that will route to your configured external connectors")
}

// analyzeQueryForDataSources determines which data sources are most relevant for a query
func (eas *EnhancedAnalyticsService) analyzeQueryForDataSources(ctx context.Context, query string, connectorIDs []string) []*models.DataConnector {
	// If specific connectors are requested, use those
	if len(connectorIDs) > 0 {
		var requestedConnectors []*models.DataConnector
		for _, id := range connectorIDs {
			connector, err := eas.connectorService.GetConnector(ctx, id)
			if err == nil && connector != nil && connector.Status == models.ConnectorStatusConnected {
				requestedConnectors = append(requestedConnectors, connector)
			}
		}
		if len(requestedConnectors) > 0 {
			return requestedConnectors
		}
	}

	// Get all active connectors
	activeConnectors, err := eas.connectorService.GetActiveConnectors(ctx)
	if err != nil {
		eas.logger.Error("Failed to get active connectors", "error", err)
		return nil
	}

	// Simple keyword-based routing (can be enhanced with LLM-based routing)
	queryLower := strings.ToLower(query)
	var relevantSources []*models.DataConnector

	for _, connector := range activeConnectors {
		switch connector.Type {
		case models.ConnectorTypeSuperset:
			// Superset is good for dashboard queries, analytics, trends
			if containsAny(queryLower, []string{"dashboard", "chart", "trend", "analytics", "visualization"}) {
				relevantSources = append(relevantSources, connector)
			}
		case models.ConnectorTypePostgres:
			// PostgreSQL is good for detailed data queries, transactions
			if containsAny(queryLower, []string{"sales", "customers", "orders", "revenue", "data", "records"}) {
				relevantSources = append(relevantSources, connector)
			}
		default:
			// Include all other types for now
			relevantSources = append(relevantSources, connector)
		}
	}

	// If no specific matches, use all available sources
	if len(relevantSources) == 0 {
		relevantSources = activeConnectors
	}

	eas.logger.Info("Selected data sources for query", "count", len(relevantSources))
	return relevantSources
}

// fetchDataFromSource retrieves data from a specific connector
func (eas *EnhancedAnalyticsService) fetchDataFromSource(ctx context.Context, connector *models.DataConnector, query string) ([]map[string]interface{}, error) {
	switch connector.Type {
	case models.ConnectorTypeSuperset:
		return eas.fetchFromSuperset(ctx, connector, query)
	case models.ConnectorTypePostgres:
		return eas.fetchFromPostgres(ctx, connector, query)
	default:
		return nil, fmt.Errorf("unsupported connector type: %s", connector.Type)
	}
}

// fetchFromSuperset retrieves data from a Superset connector
func (eas *EnhancedAnalyticsService) fetchFromSuperset(ctx context.Context, connector *models.DataConnector, query string) ([]map[string]interface{}, error) {
	config := connector.Config
	url, _ := config["url"].(string)
	username, _ := config["username"].(string)
	password, _ := config["password"].(string)
	bearerToken, _ := config["bearer_token"].(string)

	var supersetConn *connectors.SuperSetConnector

	// Use bearer token if provided, otherwise use username/password
	if bearerToken != "" {
		supersetConn = connectors.NewSuperSetConnectorWithToken(url, bearerToken, eas.logger)
	} else {
		supersetConn = connectors.NewSuperSetConnector(url, username, password, eas.logger)
	}

	// Test connection first
	if err := supersetConn.TestConnection(ctx); err != nil {
		return nil, fmt.Errorf("superset connection failed: %w", err)
	}

	eas.logger.Info("Fetching data from Superset based on user query", "query", query)

	// Use query-specific data retrieval instead of hardcoded sample data
	result, err := supersetConn.QueryDataset(ctx, query)
	if err != nil {
		eas.logger.Error("Query-specific data failed", "error", err, "query", query)
		// Fallback to sample data if query-specific fails
		result, err = supersetConn.GetSampleData(ctx)
		if err != nil {
			eas.logger.Error("Sample data also failed", "error", err)
			return nil, fmt.Errorf("failed to get data from Superset: %w", err)
		}
		eas.logger.Info("Using sample data as fallback", "rows", len(result.Data))
	}

	eas.logger.Info("Successfully retrieved data from Superset", "rows", len(result.Data))
	return result.Data, nil
}

// fetchFromPostgres retrieves data from a PostgreSQL connector
func (eas *EnhancedAnalyticsService) fetchFromPostgres(ctx context.Context, connector *models.DataConnector, query string) ([]map[string]interface{}, error) {
	// PostgreSQL connectors are disabled to prevent fallback to internal database
	// All data should come from configured external connectors (Superset, etc.)
	return nil, fmt.Errorf("PostgreSQL connectors are disabled. Please use configured external connectors (Superset, etc.) for data retrieval")
}

// fetchFromFallbackSources is disabled to prevent any fallback to internal databases
func (eas *EnhancedAnalyticsService) fetchFromFallbackSources(ctx context.Context, query string) ([]map[string]interface{}, error) {
	eas.logger.Info("Fallback data sources are disabled - only configured external connectors allowed")

	// All fallback mechanisms disabled to ensure only external connectors are used
	return nil, fmt.Errorf("fallback data sources are disabled. Please ensure your external connectors (Superset, etc.) contain the required data")
}

// generateAnalysisWithRAG creates comprehensive analysis using RAG (Retrieval Augmented Generation)
func (eas *EnhancedAnalyticsService) generateAnalysisWithRAG(ctx context.Context, combinedData []map[string]interface{}, sourceData map[string]interface{}, query string) (string, error) {
	// Build context from multiple sources
	contextBuilder := strings.Builder{}
	contextBuilder.WriteString("Data Analysis Context:\n")
	contextBuilder.WriteString(fmt.Sprintf("User Query: %s\n", query))
	contextBuilder.WriteString(fmt.Sprintf("Total Records: %d\n", len(combinedData)))

	// Add source information
	contextBuilder.WriteString("Data Sources:\n")
	for sourceName, data := range sourceData {
		if dataSlice, ok := data.([]map[string]interface{}); ok {
			contextBuilder.WriteString(fmt.Sprintf("- %s: %d records\n", sourceName, len(dataSlice)))
		}
	}

	// Add sample data for context
	if len(combinedData) > 0 {
		contextBuilder.WriteString("\nSample Data Structure:\n")
		sampleRecord := combinedData[0]
		for key := range sampleRecord {
			contextBuilder.WriteString(fmt.Sprintf("- %s\n", key))
		}
	}

	// Generate analysis with enhanced context
	enhancedQuery := fmt.Sprintf(`%s

Please provide a comprehensive analysis that includes:
1. Key insights from the data
2. Trends and patterns identified
3. Recommendations based on findings
4. Data quality observations
5. Comparison between different data sources if applicable

Original Query: %s`, contextBuilder.String(), query)

	return eas.llmConn.AnalyzeData(ctx, combinedData, enhancedQuery)
}

// Helper functions

func containsAny(text string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}
	return false
}

func (eas *EnhancedAnalyticsService) getSourceNames(connectors []*models.DataConnector) []string {
	names := make([]string, len(connectors))
	for i, connector := range connectors {
		names[i] = connector.Name
	}
	return names
}

// selectDataSourcesFromIntent chooses data sources based on parsed intent
func (eas *EnhancedAnalyticsService) selectDataSourcesFromIntent(
	ctx context.Context,
	intent models.Intent,
	requestedConnectorIDs []string,
) []*models.DataConnector {
	eas.logger.Info("Selecting data sources based on intent",
		"intent_type", intent.Type,
		"confidence", intent.Confidence)

	// If specific connectors requested, use those first
	if len(requestedConnectorIDs) > 0 {
		var requestedConnectors []*models.DataConnector
		for _, id := range requestedConnectorIDs {
			connector, err := eas.connectorService.GetConnector(ctx, id)
			if err == nil && connector != nil && connector.Status == models.ConnectorStatusConnected {
				requestedConnectors = append(requestedConnectors, connector)
			}
		}
		if len(requestedConnectors) > 0 {
			return requestedConnectors
		}
	}

	// Get all active connectors
	activeConnectors, err := eas.connectorService.GetActiveConnectors(ctx)
	if err != nil {
		eas.logger.Error("Failed to get active connectors", "error", err)
		return nil
	}

	// Advanced intent-based routing
	var relevantSources []*models.DataConnector

	switch intent.Type {
	case models.IntentTypeVisualization:
		// Prefer Superset for visualization queries
		for _, connector := range activeConnectors {
			if connector.Type == models.ConnectorTypeSuperset {
				relevantSources = append(relevantSources, connector)
			}
		}
		// Add other sources as backup
		for _, connector := range activeConnectors {
			if connector.Type != models.ConnectorTypeSuperset {
				relevantSources = append(relevantSources, connector)
			}
		}

	case models.IntentTypeSQL:
		// Prefer PostgreSQL for direct SQL queries
		for _, connector := range activeConnectors {
			if connector.Type == models.ConnectorTypePostgres {
				relevantSources = append(relevantSources, connector)
			}
		}

	case models.IntentTypeTrend, models.IntentTypeComparison:
		// Use all sources for comprehensive trend analysis
		relevantSources = activeConnectors

	case models.IntentTypeAnalytics:
		// Check parsed query for specific data source hints
		queryLower := strings.ToLower(intent.ParsedQuery.MainAction)
		if containsAny(queryLower, []string{"dashboard", "chart"}) {
			// Prefer Superset
			for _, connector := range activeConnectors {
				if connector.Type == models.ConnectorTypeSuperset {
					relevantSources = append(relevantSources, connector)
				}
			}
		} else {
			// Use PostgreSQL for detailed analytics
			for _, connector := range activeConnectors {
				if connector.Type == models.ConnectorTypePostgres {
					relevantSources = append(relevantSources, connector)
				}
			}
		}

	default:
		// Default: use all available sources
		relevantSources = activeConnectors
	}

	// If no specific matches, use all available sources
	if len(relevantSources) == 0 {
		relevantSources = activeConnectors
	}

	eas.logger.Info("Selected data sources based on intent",
		"intent_type", intent.Type,
		"sources_count", len(relevantSources))

	return relevantSources
}

// processWithBasicRouting handles fallback to original routing logic
func (eas *EnhancedAnalyticsService) processWithBasicRouting(
	ctx context.Context,
	req *EnhancedAnalyticsRequest,
	dataSources []*models.DataConnector,
	start time.Time,
	intent *models.Intent,
	planningTime string,
) (*EnhancedAnalyticsResponse, error) {
	// Use original logic for data retrieval and analysis
	allData := make(map[string]interface{})
	var combinedData []map[string]interface{}

	for _, source := range dataSources {
		data, err := eas.fetchDataFromSource(ctx, source, req.Query)
		if err != nil {
			eas.logger.Warn("Failed to fetch from source", "source", source.Name, "error", err)
			continue
		}

		if len(data) > 0 {
			allData[source.Name] = data
			combinedData = append(combinedData, data...)
			eas.logger.Info("Retrieved data from source", "source", source.Name, "rows", len(data))
		}
	}

	// Check if any data was retrieved from connectors
	if len(combinedData) == 0 {
		return nil, fmt.Errorf("no data available from configured connectors. Please check your connector configuration and ensure they contain the requested data")
	}

	// Generate analysis using basic or intent-based RAG
	var analysis string
	var err error

	if intent != nil {
		analysis, err = eas.generateAnalysisWithIntentRAG(ctx, combinedData, allData, req.Query, *intent)
	} else {
		analysis, err = eas.generateAnalysisWithRAG(ctx, combinedData, allData, req.Query)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to analyze data: %w", err)
	}

	response := &EnhancedAnalyticsResponse{
		Query:       req.Query,
		Data:        combinedData,
		Sources:     allData,
		Analysis:    analysis,
		DataSources: eas.getSourceNames(dataSources),
		Timestamp:   time.Now(),
		ProcessTime: time.Since(start).String(),
		TaskID:      fmt.Sprintf("task_%d", start.Unix()),
		Status:      "completed",
	}

	if intent != nil {
		response.Intent = intent
		response.PlanningTime = planningTime
	}

	return response, nil
}

// generateAnalysisWithIntentRAG creates analysis using intent context
func (eas *EnhancedAnalyticsService) generateAnalysisWithIntentRAG(
	ctx context.Context,
	combinedData []map[string]interface{},
	sourceData map[string]interface{},
	query string,
	intent models.Intent,
) (string, error) {
	// Build enhanced context with intent information
	contextBuilder := strings.Builder{}
	contextBuilder.WriteString("Enhanced Data Analysis Context:\n")
	contextBuilder.WriteString(fmt.Sprintf("User Query: %s\n", query))
	contextBuilder.WriteString(fmt.Sprintf("Detected Intent: %s (Confidence: %.2f)\n", intent.Type, intent.Confidence))
	contextBuilder.WriteString(fmt.Sprintf("Main Action: %s\n", intent.ParsedQuery.MainAction))
	contextBuilder.WriteString(fmt.Sprintf("Total Records: %d\n", len(combinedData)))

	// Add intent-specific context
	if len(intent.ParsedQuery.Metrics) > 0 {
		contextBuilder.WriteString(fmt.Sprintf("Identified Metrics: %s\n", strings.Join(intent.ParsedQuery.Metrics, ", ")))
	}
	if len(intent.ParsedQuery.Dimensions) > 0 {
		contextBuilder.WriteString(fmt.Sprintf("Identified Dimensions: %s\n", strings.Join(intent.ParsedQuery.Dimensions, ", ")))
	}
	if intent.ParsedQuery.TimeRange != nil {
		contextBuilder.WriteString(fmt.Sprintf("Time Range: %s\n", intent.ParsedQuery.TimeRange.Period))
	}

	// Add source information
	contextBuilder.WriteString("Data Sources:\n")
	for sourceName, data := range sourceData {
		if dataSlice, ok := data.([]map[string]interface{}); ok {
			contextBuilder.WriteString(fmt.Sprintf("- %s: %d records\n", sourceName, len(dataSlice)))
		}
	}

	// Add sample data for context
	if len(combinedData) > 0 {
		contextBuilder.WriteString("\nSample Data Structure:\n")
		sampleRecord := combinedData[0]
		for key := range sampleRecord {
			contextBuilder.WriteString(fmt.Sprintf("- %s\n", key))
		}
	}

	// Generate intent-specific analysis prompt
	var analysisPrompt string
	switch intent.Type {
	case models.IntentTypeVisualization:
		analysisPrompt = `Please provide a comprehensive analysis focused on visualization:
1. Suggest appropriate chart types for this data
2. Identify key trends suitable for visual representation
3. Recommend dashboard layout and components
4. Data transformation suggestions for better visualization
5. Interactive features that would enhance user experience`

	case models.IntentTypeComparison:
		analysisPrompt = `Please provide a comprehensive comparative analysis:
1. Key differences and similarities in the data
2. Statistical comparisons and significance
3. Trend comparisons over time periods
4. Performance gaps and opportunities
5. Recommendations based on comparative insights`

	case models.IntentTypeTrend:
		analysisPrompt = `Please provide a comprehensive trend analysis:
1. Identify patterns and trends in time-series data
2. Growth rates and trend directions
3. Seasonal patterns and anomalies
4. Forecasting insights and predictions
5. Actionable recommendations based on trends`

	case models.IntentTypeSQL:
		analysisPrompt = `Please provide a comprehensive SQL analysis result:
1. Data quality assessment of query results
2. Statistical summary of returned data
3. Potential optimizations for the query
4. Insights from the specific data subset
5. Recommendations for further analysis`

	default:
		analysisPrompt = `Please provide a comprehensive analysis that includes:
1. Key insights from the data
2. Trends and patterns identified
3. Recommendations based on findings
4. Data quality observations
5. Comparison between different data sources if applicable`
	}

	// Generate analysis with enhanced context
	enhancedQuery := fmt.Sprintf(`%s

%s

Original Query: %s`, contextBuilder.String(), analysisPrompt, query)

	return eas.llmConn.AnalyzeData(ctx, combinedData, enhancedQuery)
}