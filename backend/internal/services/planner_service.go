package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"insightiq/backend/internal/connectors"
	"insightiq/backend/internal/models"
)

// PlannerService handles intent parsing and task graph generation
type PlannerService struct {
	llmConn          *connectors.OllamaConnector
	connectorService *ConnectorService
	logger           *slog.Logger
	intentPatterns   map[models.IntentType][]string
	entityExtractor  *EntityExtractor
}

// EntityExtractor handles entity extraction from queries
type EntityExtractor struct {
	timePatterns    []*regexp.Regexp
	metricPatterns  []*regexp.Regexp
	filterPatterns  []*regexp.Regexp
	aggregatePatterns []*regexp.Regexp
}

// NewPlannerService creates a new planner service instance
func NewPlannerService(
	llmConn *connectors.OllamaConnector,
	connectorService *ConnectorService,
	logger *slog.Logger,
) *PlannerService {
	ps := &PlannerService{
		llmConn:          llmConn,
		connectorService: connectorService,
		logger:           logger.With("service", "planner"),
		intentPatterns:   initializeIntentPatterns(),
		entityExtractor:  initializeEntityExtractor(),
	}

	return ps
}

// ParseIntent analyzes user input and determines intent with confidence
func (ps *PlannerService) ParseIntent(ctx context.Context, req *models.PlannerRequest) (*models.PlannerResponse, error) {
	start := time.Now()
	ps.logger.Info("Parsing intent for query", "query", req.Query)

	// Step 1: Fast pattern-based classification first
	primaryIntent := ps.classifyIntentWithPatterns(req.Query)

	// Only use LLM for uncertain or unknown intents
	if primaryIntent.Type == models.IntentTypeUnknown || primaryIntent.Confidence < 0.8 {
		ps.logger.Debug("Pattern matching uncertain, trying LLM classification")
		llmIntent, err := ps.classifyIntentWithLLM(ctx, req.Query)
		if err == nil && llmIntent.Confidence > primaryIntent.Confidence {
			primaryIntent = llmIntent
			ps.logger.Debug("LLM provided better classification", "llm_confidence", llmIntent.Confidence)
		} else if err != nil {
			ps.logger.Warn("LLM classification failed, using pattern result", "error", err)
		}
	}

	// Step 2: Extract entities and parameters
	entities := ps.entityExtractor.ExtractEntities(req.Query)
	parsedQuery := ps.parseQueryStructure(req.Query, entities)

	// Step 3: Generate task graph
	taskGraph, err := ps.generateTaskGraph(ctx, primaryIntent, parsedQuery, req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate task graph: %w", err)
	}

	response := &models.PlannerResponse{
		Intent:      primaryIntent,
		TaskGraph:   *taskGraph,
		Confidence:  primaryIntent.Confidence,
		ProcessTime: time.Since(start),
		CreatedAt:   time.Now(),
		Metadata: map[string]interface{}{
			"entities_found": len(entities),
			"steps_planned":  len(taskGraph.Steps),
		},
	}

	ps.logger.Info("Intent parsing completed",
		"intent_type", primaryIntent.Type,
		"confidence", primaryIntent.Confidence,
		"steps", len(taskGraph.Steps),
		"duration", response.ProcessTime)

	return response, nil
}

// classifyIntentWithLLM uses LLM for sophisticated intent classification
func (ps *PlannerService) classifyIntentWithLLM(ctx context.Context, query string) (models.Intent, error) {
	prompt := fmt.Sprintf(`Classify this query intent. Return ONLY JSON:
{"type": "analytics|sql|visualization|comparison|trend|filter|aggregation|join|unknown", "confidence": 0.0-1.0}

Query: "%s"

Rules:
- analytics: data analysis
- sql: SQL queries
- visualization: charts/dashboards
- comparison: comparing data
- trend: time-series analysis
- filter: filtering data
- aggregation: sum/count/avg
- join: combining sources
- unknown: unclear`, query)

	// Use a simple data structure for LLM response
	data := []map[string]interface{}{
		{"query": query, "task": "intent_classification"},
	}

	response, err := ps.llmConn.AnalyzeData(ctx, data, prompt)
	if err != nil {
		return models.Intent{}, err
	}

	// Try to parse JSON from response
	var intent models.Intent
	if err := ps.parseIntentFromLLMResponse(response, &intent); err != nil {
		ps.logger.Warn("Failed to parse LLM response as JSON", "response", response)
		// Fallback to pattern-based classification
		return ps.classifyIntentWithPatterns(query), nil
	}

	return intent, nil
}

// parseIntentFromLLMResponse extracts JSON from LLM response
func (ps *PlannerService) parseIntentFromLLMResponse(response string, intent *models.Intent) error {
	// Try to find JSON in the response
	startIdx := strings.Index(response, "{")
	endIdx := strings.LastIndex(response, "}")

	if startIdx == -1 || endIdx == -1 || startIdx >= endIdx {
		return fmt.Errorf("no valid JSON found in response")
	}

	jsonStr := response[startIdx : endIdx+1]

	// Parse the JSON
	var rawIntent map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &rawIntent); err != nil {
		return err
	}

	// Map to our Intent struct (simplified)
	if intentType, ok := rawIntent["type"].(string); ok {
		intent.Type = models.IntentType(intentType)
	}
	if confidence, ok := rawIntent["confidence"].(float64); ok {
		intent.Confidence = confidence
	}

	// Initialize basic structures
	intent.Entities = make(map[string]interface{})
	intent.Parameters = make(map[string]string)
	intent.ParsedQuery = models.ParsedQuery{
		MainAction:   "llm_classified",
		OutputFormat: "table",
		Metrics:      []string{},
		Dimensions:   []string{},
		DataSources:  []string{},
	}

	return nil
}

// classifyIntentWithPatterns provides fast pattern-based classification
func (ps *PlannerService) classifyIntentWithPatterns(query string) models.Intent {
	queryLower := strings.ToLower(query)

	// Score each intent type by pattern matches
	intentScores := make(map[models.IntentType]int)

	// Check patterns for each intent type and count matches
	for intentType, patterns := range ps.intentPatterns {
		score := 0
		for _, pattern := range patterns {
			if strings.Contains(queryLower, pattern) {
				score++
			}
		}
		if score > 0 {
			intentScores[intentType] = score
		}
	}

	// If no patterns matched, return unknown
	if len(intentScores) == 0 {
		return models.Intent{
			Type:        models.IntentTypeUnknown,
			Confidence:  0.3,
			Entities:    make(map[string]interface{}),
			Parameters:  make(map[string]string),
			ParsedQuery: models.ParsedQuery{MainAction: "unknown", OutputFormat: "table"},
		}
	}

	// Find the intent type with highest score
	var bestIntent models.IntentType
	var bestScore int
	for intentType, score := range intentScores {
		if score > bestScore {
			bestScore = score
			bestIntent = intentType
		}
	}

	// Calculate confidence based on score and query complexity
	confidence := float64(bestScore) * 0.3 // Base confidence from pattern matches
	if bestScore >= 3 {
		confidence = 0.9 // High confidence for multiple matches
	} else if bestScore == 2 {
		confidence = 0.8 // Good confidence for two matches
	} else {
		confidence = 0.7 // Moderate confidence for single match
	}

	// Enhance parsed query based on detected patterns
	parsedQuery := ps.enhanceParsedQueryFromPatterns(queryLower, bestIntent)

	return models.Intent{
		Type:        bestIntent,
		Confidence:  confidence,
		Entities:    ps.entityExtractor.ExtractEntities(query),
		Parameters:  make(map[string]string),
		ParsedQuery: parsedQuery,
	}
}

// enhanceParsedQueryFromPatterns creates better parsed query from pattern analysis
func (ps *PlannerService) enhanceParsedQueryFromPatterns(queryLower string, intentType models.IntentType) models.ParsedQuery {
	pq := models.ParsedQuery{
		OutputFormat: "table",
		Metrics:      []string{},
		Dimensions:   []string{},
		DataSources:  []string{},
	}

	switch intentType {
	case models.IntentTypeVisualization:
		pq.MainAction = "create_visualization"
		pq.OutputFormat = "chart"
		// Detect chart types
		if containsAny(queryLower, []string{"dashboard"}) {
			pq.OutputFormat = "dashboard"
		}
	case models.IntentTypeComparison:
		pq.MainAction = "compare_data"
		// Detect comparison targets
		if containsAny(queryLower, []string{"between", "vs", "versus"}) {
			pq.MainAction = "comparative_analysis"
		}
	case models.IntentTypeTrend:
		pq.MainAction = "analyze_trends"
		// Detect time periods
		if containsAny(queryLower, []string{"quarter", "monthly", "weekly"}) {
			pq.Dimensions = append(pq.Dimensions, "time")
		}
	case models.IntentTypeSQL:
		pq.MainAction = "execute_sql"
		pq.OutputFormat = "table"
	case models.IntentTypeAnalytics:
		pq.MainAction = "analyze_data"
		// Detect common metrics
		if containsAny(queryLower, []string{"sales", "revenue"}) {
			pq.Metrics = append(pq.Metrics, "revenue")
		}
		if containsAny(queryLower, []string{"performance", "kpi"}) {
			pq.MainAction = "performance_analysis"
		}
	default:
		pq.MainAction = "general_query"
	}

	return pq
}

// generateTaskGraph creates an execution plan based on intent and parsed query
func (ps *PlannerService) generateTaskGraph(
	ctx context.Context,
	intent models.Intent,
	parsedQuery models.ParsedQuery,
	req *models.PlannerRequest,
) (*models.TaskGraph, error) {

	taskGraph := &models.TaskGraph{
		ID:           fmt.Sprintf("plan_%d", time.Now().Unix()),
		Query:        req.Query,
		Intent:       intent,
		Steps:        []models.TaskStep{},
		Dependencies: make(map[string][]string),
		CreatedAt:    time.Now(),
		Status:       models.TaskStatusPlanned,
	}

	// Generate steps based on intent type
	switch intent.Type {
	case models.IntentTypeAnalytics:
		ps.addAnalyticsSteps(taskGraph, parsedQuery)
	case models.IntentTypeSQL:
		ps.addSQLSteps(taskGraph, parsedQuery)
	case models.IntentTypeVisualization:
		ps.addVisualizationSteps(taskGraph, parsedQuery)
	case models.IntentTypeComparison:
		ps.addComparisonSteps(taskGraph, parsedQuery)
	case models.IntentTypeTrend:
		ps.addTrendSteps(taskGraph, parsedQuery)
	default:
		ps.addDefaultSteps(taskGraph, parsedQuery)
	}

	// Calculate estimated execution time
	taskGraph.EstimatedTime = ps.calculateEstimatedTime(taskGraph.Steps)

	return taskGraph, nil
}

// Helper methods for different step types
func (ps *PlannerService) addAnalyticsSteps(taskGraph *models.TaskGraph, parsedQuery models.ParsedQuery) {
	steps := []models.TaskStep{
		{
			ID:          "data_discovery",
			Type:        models.TaskStepTypeDataRetrieval,
			Description: "Discover available data sources",
			Action:      "discover_sources",
			Priority:    1,
			EstimatedTime: 2 * time.Second,
		},
		{
			ID:          "data_retrieval",
			Type:        models.TaskStepTypeDataRetrieval,
			Description: "Retrieve data from identified sources",
			Action:      "fetch_data",
			Dependencies: []string{"data_discovery"},
			Priority:    2,
			EstimatedTime: 5 * time.Second,
		},
		{
			ID:          "data_analysis",
			Type:        models.TaskStepTypeAnalysis,
			Description: "Perform analytical processing on retrieved data",
			Action:      "analyze_data",
			Dependencies: []string{"data_retrieval"},
			Priority:    3,
			EstimatedTime: 8 * time.Second,
		},
	}

	taskGraph.Steps = append(taskGraph.Steps, steps...)
	ps.updateDependencies(taskGraph)
}

func (ps *PlannerService) addSQLSteps(taskGraph *models.TaskGraph, parsedQuery models.ParsedQuery) {
	steps := []models.TaskStep{
		{
			ID:          "sql_validation",
			Type:        models.TaskStepTypeValidation,
			Description: "Validate SQL query syntax and permissions",
			Action:      "validate_sql",
			Priority:    1,
			EstimatedTime: 1 * time.Second,
		},
		{
			ID:          "sql_execution",
			Type:        models.TaskStepTypeDataRetrieval,
			Description: "Execute SQL query on target database",
			Action:      "execute_sql",
			Dependencies: []string{"sql_validation"},
			Priority:    2,
			EstimatedTime: 10 * time.Second,
		},
	}

	taskGraph.Steps = append(taskGraph.Steps, steps...)
	ps.updateDependencies(taskGraph)
}

func (ps *PlannerService) addVisualizationSteps(taskGraph *models.TaskGraph, parsedQuery models.ParsedQuery) {
	steps := []models.TaskStep{
		{
			ID:          "data_preparation",
			Type:        models.TaskStepTypeTransformation,
			Description: "Prepare data for visualization",
			Action:      "prepare_viz_data",
			Priority:    1,
			EstimatedTime: 3 * time.Second,
		},
		{
			ID:          "chart_generation",
			Type:        models.TaskStepTypeVisualization,
			Description: "Generate chart or dashboard",
			Action:      "create_visualization",
			Dependencies: []string{"data_preparation"},
			Priority:    2,
			EstimatedTime: 5 * time.Second,
		},
	}

	taskGraph.Steps = append(taskGraph.Steps, steps...)
	ps.updateDependencies(taskGraph)
}

func (ps *PlannerService) addComparisonSteps(taskGraph *models.TaskGraph, parsedQuery models.ParsedQuery) {
	steps := []models.TaskStep{
		{
			ID:          "multi_source_retrieval",
			Type:        models.TaskStepTypeDataRetrieval,
			Description: "Retrieve data from multiple sources for comparison",
			Action:      "fetch_comparison_data",
			Priority:    1,
			EstimatedTime: 7 * time.Second,
		},
		{
			ID:          "data_alignment",
			Type:        models.TaskStepTypeTransformation,
			Description: "Align data schemas and formats",
			Action:      "align_data",
			Dependencies: []string{"multi_source_retrieval"},
			Priority:    2,
			EstimatedTime: 4 * time.Second,
		},
		{
			ID:          "comparison_analysis",
			Type:        models.TaskStepTypeAnalysis,
			Description: "Perform comparative analysis",
			Action:      "compare_data",
			Dependencies: []string{"data_alignment"},
			Priority:    3,
			EstimatedTime: 6 * time.Second,
		},
	}

	taskGraph.Steps = append(taskGraph.Steps, steps...)
	ps.updateDependencies(taskGraph)
}

func (ps *PlannerService) addTrendSteps(taskGraph *models.TaskGraph, parsedQuery models.ParsedQuery) {
	steps := []models.TaskStep{
		{
			ID:          "time_series_data",
			Type:        models.TaskStepTypeDataRetrieval,
			Description: "Retrieve time-series data",
			Action:      "fetch_time_series",
			Priority:    1,
			EstimatedTime: 5 * time.Second,
		},
		{
			ID:          "trend_analysis",
			Type:        models.TaskStepTypeAnalysis,
			Description: "Analyze trends and patterns over time",
			Action:      "analyze_trends",
			Dependencies: []string{"time_series_data"},
			Priority:    2,
			EstimatedTime: 8 * time.Second,
		},
	}

	taskGraph.Steps = append(taskGraph.Steps, steps...)
	ps.updateDependencies(taskGraph)
}

func (ps *PlannerService) addDefaultSteps(taskGraph *models.TaskGraph, parsedQuery models.ParsedQuery) {
	steps := []models.TaskStep{
		{
			ID:          "general_processing",
			Type:        models.TaskStepTypeAnalysis,
			Description: "General query processing",
			Action:      "process_query",
			Priority:    1,
			EstimatedTime: 5 * time.Second,
		},
	}

	taskGraph.Steps = append(taskGraph.Steps, steps...)
	ps.updateDependencies(taskGraph)
}

// Helper methods
func (ps *PlannerService) updateDependencies(taskGraph *models.TaskGraph) {
	taskGraph.Dependencies = make(map[string][]string)
	for _, step := range taskGraph.Steps {
		if len(step.Dependencies) > 0 {
			taskGraph.Dependencies[step.ID] = step.Dependencies
		}
	}
}

func (ps *PlannerService) calculateEstimatedTime(steps []models.TaskStep) time.Duration {
	var totalTime time.Duration
	for _, step := range steps {
		totalTime += step.EstimatedTime
	}
	return totalTime
}

func (ps *PlannerService) parseQueryStructure(query string, entities map[string]interface{}) models.ParsedQuery {
	return models.ParsedQuery{
		MainAction:   "analyze",
		DataSources:  []string{},
		Metrics:      []string{},
		Dimensions:   []string{},
		Filters:      []models.Filter{},
		Aggregations: []models.Aggregation{},
		SortBy:       []models.SortCriteria{},
		GroupBy:      []string{},
		Joins:        []models.Join{},
		OutputFormat: "table",
		Metadata:     entities,
	}
}

func (ps *PlannerService) mapToParsedQuery(data map[string]interface{}) models.ParsedQuery {
	pq := models.ParsedQuery{
		OutputFormat: "table",
	}

	if mainAction, ok := data["main_action"].(string); ok {
		pq.MainAction = mainAction
	}
	if outputFormat, ok := data["output_format"].(string); ok {
		pq.OutputFormat = outputFormat
	}

	// Handle arrays
	if dataSources, ok := data["data_sources"].([]interface{}); ok {
		for _, ds := range dataSources {
			if str, ok := ds.(string); ok {
				pq.DataSources = append(pq.DataSources, str)
			}
		}
	}

	if metrics, ok := data["metrics"].([]interface{}); ok {
		for _, m := range metrics {
			if str, ok := m.(string); ok {
				pq.Metrics = append(pq.Metrics, str)
			}
		}
	}

	if dimensions, ok := data["dimensions"].([]interface{}); ok {
		for _, d := range dimensions {
			if str, ok := d.(string); ok {
				pq.Dimensions = append(pq.Dimensions, str)
			}
		}
	}

	return pq
}

// ExtractEntities extracts entities from the query text
func (ee *EntityExtractor) ExtractEntities(query string) map[string]interface{} {
	entities := make(map[string]interface{})

	// Extract time-related entities
	for _, pattern := range ee.timePatterns {
		matches := pattern.FindAllString(query, -1)
		if len(matches) > 0 {
			entities["time_expressions"] = matches
		}
	}

	// Extract metric-related entities
	for _, pattern := range ee.metricPatterns {
		matches := pattern.FindAllString(query, -1)
		if len(matches) > 0 {
			entities["metrics"] = matches
		}
	}

	return entities
}

// Initialize functions
func initializeIntentPatterns() map[models.IntentType][]string {
	return map[models.IntentType][]string{
		models.IntentTypeAnalytics: {
			"analyze", "analysis", "insights", "understand", "explain", "what", "how", "why",
			"performance", "metrics", "kpi", "report", "summary",
		},
		models.IntentTypeSQL: {
			"select", "from", "where", "join", "group by", "order by", "having",
			"sql", "query", "execute", "run",
		},
		models.IntentTypeVisualization: {
			"chart", "graph", "plot", "dashboard", "visualize", "show", "display",
			"bar chart", "line chart", "pie chart", "scatter plot",
		},
		models.IntentTypeComparison: {
			"compare", "comparison", "versus", "vs", "difference", "between",
			"against", "relative to", "compared to",
		},
		models.IntentTypeTrend: {
			"trend", "over time", "time series", "growth", "decline", "pattern",
			"monthly", "weekly", "daily", "quarterly", "yearly",
		},
		models.IntentTypeFilter: {
			"filter", "where", "only", "exclude", "include", "containing",
			"matching", "equal to", "greater than", "less than",
		},
		models.IntentTypeAggregation: {
			"sum", "count", "average", "max", "min", "total", "aggregate",
			"group", "by", "per", "each",
		},
	}
}

func initializeEntityExtractor() *EntityExtractor {
	return &EntityExtractor{
		timePatterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)(last|past|previous)\s+(week|month|quarter|year|day|hour)`),
			regexp.MustCompile(`(?i)(this|current)\s+(week|month|quarter|year|day)`),
			regexp.MustCompile(`(?i)\d{4}-\d{2}-\d{2}`), // Date format
			regexp.MustCompile(`(?i)(yesterday|today|tomorrow)`),
		},
		metricPatterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)(sales|revenue|profit|cost|price|amount|value)`),
			regexp.MustCompile(`(?i)(count|number|quantity|volume|total)`),
			regexp.MustCompile(`(?i)(rate|percentage|ratio|proportion)`),
		},
		filterPatterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)(where|filter|only|exclude|include)`),
		},
		aggregatePatterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)(sum|count|avg|average|max|min|total)`),
		},
	}
}