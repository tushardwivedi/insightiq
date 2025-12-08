package services

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"insightiq/backend/internal/connectors"
	"insightiq/backend/internal/models"
)

// EnhancedAnalyticsService provides intelligent data source routing and RAG capabilities
type EnhancedAnalyticsService struct {
	connectorService   *ConnectorService
	plannerService     *PlannerService
	ragService         *RAGQueryService
	orchestrator       *DataSourceOrchestrator
	validationService  *DataValidationService
	llmConn            *connectors.OllamaConnector
	fallbackPostgres   *connectors.PostgresConnector
	fallbackSuperset   *connectors.SuperSetConnector
	db                 *sqlx.DB
	logger             *slog.Logger
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
	Query           string                   `json:"query"`
	Data            []map[string]interface{} `json:"data"`
	Sources         map[string]interface{}   `json:"sources"`
	Analysis        string                   `json:"analysis"`
	DataSources     []string                 `json:"data_sources"`
	Timestamp       time.Time                `json:"timestamp"`
	ProcessTime     string                   `json:"process_time"`
	TaskID          string                   `json:"task_id"`
	Status          string                   `json:"status"`
	Intent          *models.Intent           `json:"intent,omitempty"`
	TaskGraph       *models.TaskGraph        `json:"task_graph,omitempty"`
	PlanningTime    string                   `json:"planning_time,omitempty"`
	DataValidation  *ValidationResult        `json:"data_validation,omitempty"`
}

func NewEnhancedAnalyticsService(
	connectorService *ConnectorService,
	llm *connectors.OllamaConnector,
	fallbackPostgres *connectors.PostgresConnector,
	fallbackSuperset *connectors.SuperSetConnector,
	ragService *RAGQueryService,
	logger *slog.Logger,
) *EnhancedAnalyticsService {
	eas := &EnhancedAnalyticsService{
		connectorService: connectorService,
		llmConn:          llm,
		fallbackPostgres: fallbackPostgres,
		fallbackSuperset: fallbackSuperset,
		ragService:       ragService,
		db:               nil, // Will be set if needed for file upload queries
		logger:           logger.With("service", "enhanced_analytics"),
	}

	// Initialize planner service
	eas.plannerService = NewPlannerService(llm, connectorService, logger)

	// Initialize data source orchestrator
	eas.orchestrator = NewDataSourceOrchestrator(connectorService, ragService, logger)

	// Initialize data validation service
	eas.validationService = NewDataValidationService(logger)

	return eas
}

// SetDatabase sets the database connection for querying uploaded files
func (eas *EnhancedAnalyticsService) SetDatabase(db *sqlx.DB) {
	eas.db = db
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

	// 2.5. Perform RAG vector search for relevant file context
	if eas.ragService != nil {
		eas.logger.Info("Performing RAG vector search", "query", req.Query)
		ragStart := time.Now()
		ragResults, err := eas.ragService.GetRelevantContext(ctx, req.Query, 10)
		ragDuration := time.Since(ragStart)

		if err != nil {
			eas.logger.Warn("RAG search failed", "error", err, "duration", ragDuration)
		} else if len(ragResults) > 0 {
			eas.logger.Info("RAG search completed successfully",
				"results_count", len(ragResults),
				"duration", ragDuration)

			// Convert RAG results to data format and merge
			ragData := eas.convertRAGResultsToData(ragResults)
			allData["rag_vector_search"] = ragData
			combinedData = append(combinedData, ragData...)
			eas.logger.Info("RAG results merged into combined data",
				"rag_rows", len(ragData),
				"total_rows", len(combinedData))
		} else {
			eas.logger.Info("RAG search returned no results", "duration", ragDuration)
		}
	} else {
		eas.logger.Debug("RAG service not initialized, skipping vector search")
	}

	// 3. Check if any data was retrieved from connectors
	if len(combinedData) == 0 {
		return nil, fmt.Errorf("no data available from configured connectors. Please check your connector configuration and ensure they contain the requested data")
	}

	// 3.5. Validate data quality and accuracy
	var validationResult *ValidationResult
	if eas.validationService != nil {
		eas.logger.Info("Performing data validation", "record_count", len(combinedData))
		validationResult, err = eas.validationService.ValidateData(ctx, combinedData)
		if err != nil {
			eas.logger.Warn("Data validation failed", "error", err)
		} else {
			eas.logger.Info("Data validation completed",
				"quality_score", validationResult.QualityScore,
				"warnings", len(validationResult.Warnings),
				"errors", len(validationResult.Errors))

			// Log quality warnings
			if validationResult.QualityScore < 70 {
				eas.logger.Warn("Data quality score below threshold",
					"score", validationResult.QualityScore,
					"recommendations", validationResult.Recommendations)
			}
		}
	}

	// 4. Generate comprehensive analysis with enhanced RAG context using intent
	analysis, err := eas.generateAnalysisWithIntentRAG(ctx, combinedData, allData, req.Query, plannerResponse.Intent)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze data: %w", err)
	}

	return &EnhancedAnalyticsResponse{
		Query:          req.Query,
		Data:           combinedData,
		Sources:        allData,
		Analysis:       analysis,
		DataSources:    eas.getSourceNames(dataSources),
		Timestamp:      time.Now(),
		ProcessTime:    time.Since(start).String(),
		TaskID:         fmt.Sprintf("task_%d", start.Unix()),
		Status:         "completed",
		Intent:         &plannerResponse.Intent,
		TaskGraph:      &plannerResponse.TaskGraph,
		PlanningTime:   planningTime.String(),
		DataValidation: validationResult,
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
	case models.ConnectorTypeFileUpload:
		return eas.fetchFromUploadedFile(ctx, connector, query)
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

	eas.logger.Info("🔍 Fetching data from Superset based on user query", "query", query)

	// Try to find relevant dataset based on query
	dashboard, err := supersetConn.FindRelevantDataset(ctx, query)
	if err != nil {
		eas.logger.Warn("❌ No matching dashboard found, trying QueryDataset fallback", "error", err)
		// Fallback to QueryDataset which has its own heuristics
		result, err := supersetConn.QueryDataset(ctx, query)
		if err != nil {
			eas.logger.Error("Query failed", "error", err)
			return nil, fmt.Errorf("no matching dataset found and query failed: %w. Available dashboards may not match query keywords.", err)
		}
		eas.logger.Info("✅ Query executed using fallback method", "rows", len(result.Data))
		return result.Data, nil
	}

	// Extract dashboard ID and title
	dashboardID := 0
	if id, ok := dashboard["id"].(float64); ok {
		dashboardID = int(id)
	}

	dashTitle := ""
	if title, ok := dashboard["dashboard_title"].(string); ok {
		dashTitle = title
	}

	eas.logger.Info("✅ Found matching dashboard", "dashboard", dashTitle, "id", dashboardID)

	// Query data from the discovered dashboard
	result, err := supersetConn.QueryDashboardData(ctx, dashboardID, query)
	if err != nil {
		eas.logger.Error("Failed to query dashboard", "dashboard", dashTitle, "error", err)
		return nil, fmt.Errorf("failed to query dashboard '%s' (id: %d): %w", dashTitle, dashboardID, err)
	}

	eas.logger.Info("✅ Successfully retrieved data from Superset", "dashboard", dashTitle, "rows", len(result.Data))
	return result.Data, nil
}

// fetchFromPostgres retrieves data from a PostgreSQL connector with security controls
func (eas *EnhancedAnalyticsService) fetchFromPostgres(ctx context.Context, connector *models.DataConnector, query string) ([]map[string]interface{}, error) {
	config := connector.Config

	// SECURITY: Extract connection parameters
	host, _ := config["host"].(string)
	port, _ := config["port"].(string)
	database, _ := config["database"].(string)
	username, _ := config["username"].(string)
	password, _ := config["password"].(string)

	// SECURITY: Validate required parameters
	if host == "" || database == "" || username == "" {
		return nil, fmt.Errorf("invalid PostgreSQL connector configuration: missing required fields")
	}

	// SECURITY: Prevent connection to localhost/internal databases
	// This prevents access to the application's own database
	if eas.isInternalHost(host) {
		eas.logger.Warn("Blocked connection attempt to internal database",
			"host", host,
			"connector", connector.Name)
		return nil, fmt.Errorf("security policy violation: cannot connect to internal/localhost databases. Use external PostgreSQL instances only")
	}

	// Set default port if not specified
	if port == "" {
		port = "5432"
	}

	// Build connection string for EXTERNAL PostgreSQL only
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		host, port, username, password, database)

	eas.logger.Info("Connecting to external PostgreSQL database",
		"host", host,
		"port", port,
		"database", database,
		"connector", connector.Name)

	// Connect to external PostgreSQL
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		eas.logger.Error("Failed to connect to external PostgreSQL",
			"host", host,
			"error", err)
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}
	defer db.Close()

	// SECURITY: Set connection timeout
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(time.Minute * 5)

	// SECURITY: Use LLM to generate safe SQL query
	// This prevents direct SQL injection and ensures only SELECT queries
	sqlQuery, err := eas.generateSafeSQLQuery(ctx, query, database)
	if err != nil {
		eas.logger.Error("Failed to generate safe SQL query", "error", err)
		return nil, fmt.Errorf("failed to generate query: %w", err)
	}

	// SECURITY: Validate that query is read-only (SELECT only)
	if !eas.isReadOnlyQuery(sqlQuery) {
		eas.logger.Warn("Blocked non-SELECT query attempt",
			"query", sqlQuery,
			"connector", connector.Name)
		return nil, fmt.Errorf("security policy violation: only SELECT queries are allowed")
	}

	// SECURITY: Set query timeout
	queryCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Execute the query
	rows, err := db.QueryxContext(queryCtx, sqlQuery)
	if err != nil {
		eas.logger.Error("Failed to execute PostgreSQL query",
			"host", host,
			"error", err)
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Convert rows to map slice
	var results []map[string]interface{}
	for rows.Next() {
		row := make(map[string]interface{})
		if err := rows.MapScan(row); err != nil {
			eas.logger.Error("Failed to scan row", "error", err)
			continue
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating query results: %w", err)
	}

	eas.logger.Info("Successfully retrieved data from external PostgreSQL",
		"host", host,
		"connector", connector.Name,
		"rows", len(results))

	return results, nil
}

// isInternalHost checks if a host is internal/localhost
// SECURITY: This prevents connecting to the application's own database
func (eas *EnhancedAnalyticsService) isInternalHost(host string) bool {
	internalHosts := []string{
		"localhost",
		"127.0.0.1",
		"::1",
		"postgres", // Docker internal service name
		"0.0.0.0",
		"insightiq-postgres", // Our internal database
	}

	hostLower := strings.ToLower(strings.TrimSpace(host))

	for _, internal := range internalHosts {
		if hostLower == internal {
			return true
		}
	}

	// Check for 127.* addresses
	if strings.HasPrefix(hostLower, "127.") {
		return true
	}

	// Check for 10.*, 172.16-31.*, 192.168.* (private networks)
	// These are often used for internal services
	if strings.HasPrefix(hostLower, "10.") ||
	   strings.HasPrefix(hostLower, "192.168.") {
		return true
	}

	// Check for 172.16.* through 172.31.*
	if strings.HasPrefix(hostLower, "172.") {
		parts := strings.Split(hostLower, ".")
		if len(parts) >= 2 {
			if secondOctet, err := strconv.Atoi(parts[1]); err == nil {
				if secondOctet >= 16 && secondOctet <= 31 {
					return true
				}
			}
		}
	}

	return false
}

// isReadOnlyQuery validates that a SQL query is read-only (SELECT only)
// SECURITY: Prevents data modification through SQL injection
func (eas *EnhancedAnalyticsService) isReadOnlyQuery(sqlQuery string) bool {
	queryUpper := strings.ToUpper(strings.TrimSpace(sqlQuery))

	// Must start with SELECT
	if !strings.HasPrefix(queryUpper, "SELECT") {
		return false
	}

	// Check for dangerous keywords
	dangerousKeywords := []string{
		"INSERT", "UPDATE", "DELETE", "DROP", "CREATE", "ALTER",
		"TRUNCATE", "GRANT", "REVOKE", "EXEC", "EXECUTE",
		"CALL", "PROCEDURE", "FUNCTION",
	}

	for _, keyword := range dangerousKeywords {
		if strings.Contains(queryUpper, keyword) {
			return false
		}
	}

	return true
}

// generateSafeSQLQuery generates a safe SQL query using LLM
// SECURITY: Uses LLM to convert natural language to safe SELECT queries
func (eas *EnhancedAnalyticsService) generateSafeSQLQuery(ctx context.Context, query string, database string) (string, error) {
	// For now, if query looks like SQL, validate it
	// Otherwise, generate SELECT query
	queryTrimmed := strings.TrimSpace(query)
	queryUpper := strings.ToUpper(queryTrimmed)

	if strings.HasPrefix(queryUpper, "SELECT") {
		// Query is already SQL, validate it
		if !eas.isReadOnlyQuery(queryTrimmed) {
			return "", fmt.Errorf("query contains non-SELECT statements")
		}
		return queryTrimmed, nil
	}

	// For natural language queries, generate a safe query
	// In production, use LLM to convert natural language to SQL
	// For now, return a simple default query
	eas.logger.Warn("Natural language to SQL conversion not fully implemented",
		"query", query)

	return "", fmt.Errorf("please provide a SELECT query. Natural language to SQL conversion requires LLM integration")
}

// fetchFromUploadedFile retrieves data from an uploaded file connector
func (eas *EnhancedAnalyticsService) fetchFromUploadedFile(ctx context.Context, connector *models.DataConnector, query string) ([]map[string]interface{}, error) {
	if eas.db == nil {
		return nil, fmt.Errorf("database connection not available for file upload queries")
	}

	// Extract table name from connector config
	tableName, ok := connector.Config["table_name"].(string)
	if !ok || tableName == "" {
		return nil, fmt.Errorf("table_name not found in file upload connector config")
	}

	eas.logger.Info("Fetching data from uploaded file",
		"connector_name", connector.Name,
		"table_name", tableName,
		"query", query)

	// Determine if query contains SQL or if we need to query all data
	sqlQuery := query
	if !strings.Contains(strings.ToUpper(query), "SELECT") {
		// If no SELECT statement, query all data from the table
		sqlQuery = fmt.Sprintf("SELECT * FROM %s LIMIT 1000", tableName)
		eas.logger.Info("No SQL detected, querying all data from table", "sql", sqlQuery)
	} else {
		// Validate that the query references the correct table
		if !strings.Contains(strings.ToLower(query), strings.ToLower(tableName)) {
			eas.logger.Warn("Query doesn't reference table, adjusting query", "expected_table", tableName)
			// This is a basic validation - in production you'd want more sophisticated parsing
		}
	}

	// Execute the query
	rows, err := eas.db.QueryxContext(ctx, sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query on uploaded file table: %w", err)
	}
	defer rows.Close()

	// Convert rows to map slice
	var results []map[string]interface{}
	for rows.Next() {
		row := make(map[string]interface{})
		if err := rows.MapScan(row); err != nil {
			eas.logger.Error("Failed to scan row", "error", err)
			continue
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating query results: %w", err)
	}

	eas.logger.Info("Successfully retrieved data from uploaded file",
		"connector_name", connector.Name,
		"table_name", tableName,
		"rows", len(results))

	return results, nil
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
		// Prefer uploaded files and PostgreSQL for direct SQL queries
		for _, connector := range activeConnectors {
			if connector.Type == models.ConnectorTypeFileUpload || connector.Type == models.ConnectorTypePostgres {
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
			// Use uploaded files and PostgreSQL for detailed analytics
			for _, connector := range activeConnectors {
				if connector.Type == models.ConnectorTypeFileUpload || connector.Type == models.ConnectorTypePostgres {
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

	// Add source information with RAG context awareness
	contextBuilder.WriteString("Data Sources:\n")
	var hasRAGResults bool
	var ragResultCount int

	for sourceName, data := range sourceData {
		if dataSlice, ok := data.([]map[string]interface{}); ok {
			contextBuilder.WriteString(fmt.Sprintf("- %s: %d records\n", sourceName, len(dataSlice)))

			// Track RAG results separately
			if sourceName == "rag_vector_search" {
				hasRAGResults = true
				ragResultCount = len(dataSlice)
			}
		}
	}

	// Add RAG-specific context if present
	if hasRAGResults {
		contextBuilder.WriteString(fmt.Sprintf("\nVector Search Context:\n"))
		contextBuilder.WriteString(fmt.Sprintf("- Found %d relevant documents from uploaded files\n", ragResultCount))
		contextBuilder.WriteString("- RAG results include file schemas, summaries, and specific data rows\n")
		contextBuilder.WriteString("- Each result has a similarity score indicating relevance to the query\n")
		contextBuilder.WriteString("- Consider using these vectorized contexts to enrich your analysis\n")
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

// convertRAGResultsToData converts RAG search results to standard data format
func (eas *EnhancedAnalyticsService) convertRAGResultsToData(ragResults []RAGResult) []map[string]interface{} {
	data := make([]map[string]interface{}, 0, len(ragResults))

	for _, result := range ragResults {
		row := map[string]interface{}{
			"source":        "rag_vector_search",
			"file_id":       result.FileID,
			"file_name":     result.FileName,
			"table_name":    result.TableName,
			"content_type":  result.Type,
			"content":       result.Content,
			"score":         result.Score,
			"vector_id":     result.VectorID,
		}

		// Add data quality metrics if available
		if result.QualityScore != nil {
			row["data_quality_score"] = *result.QualityScore
		}
		if result.MissingPercent != nil {
			row["missing_data_percent"] = *result.MissingPercent
		}
		if result.DuplicateRows != nil {
			row["duplicate_rows"] = *result.DuplicateRows
		}

		// Add column categorization if available
		if len(result.NumericColumns) > 0 {
			row["numeric_columns"] = strings.Join(result.NumericColumns, ", ")
		}
		if len(result.CategoricalColumns) > 0 {
			row["categorical_columns"] = strings.Join(result.CategoricalColumns, ", ")
		}
		if len(result.DateColumns) > 0 {
			row["date_columns"] = strings.Join(result.DateColumns, ", ")
		}

		data = append(data, row)
	}

	eas.logger.Debug("Converted RAG results to data format", "count", len(data))
	return data
}