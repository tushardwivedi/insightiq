package services

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"insightiq/backend/internal/models"
)

// DataSourceOrchestrator manages intelligent routing and fetching from multiple data sources
// with priority-based selection and fallback mechanisms
type DataSourceOrchestrator struct {
	connectorService *ConnectorService
	ragService       *RAGQueryService
	logger           *slog.Logger
}

// DataSourcePriority defines the priority order for data sources
type DataSourcePriority int

const (
	PriorityRAG DataSourcePriority = iota + 1
	PrioritySuperset
	PriorityFileUpload
	PriorityPostgres
)

// FetchResult represents a result from fetching data from a source
type FetchResult struct {
	SourceName   string
	SourceType   models.ConnectorType
	Data         []map[string]interface{}
	Priority     DataSourcePriority
	FetchTime    time.Duration
	Error        error
	RecordCount  int
	IsRAGResult  bool
}

// OrchestrationStrategy defines how to fetch data from sources
type OrchestrationStrategy struct {
	// ParallelFetch enables concurrent fetching from multiple sources
	ParallelFetch bool
	// MaxConcurrency limits concurrent fetch operations
	MaxConcurrency int
	// UseRAG enables RAG vector search
	UseRAG bool
	// RAGLimit sets maximum RAG results
	RAGLimit int
	// Timeout for each fetch operation
	FetchTimeout time.Duration
	// PreferredSources lists connector IDs to prioritize
	PreferredSources []string
}

// NewDataSourceOrchestrator creates a new data source orchestrator
func NewDataSourceOrchestrator(
	connectorService *ConnectorService,
	ragService *RAGQueryService,
	logger *slog.Logger,
) *DataSourceOrchestrator {
	return &DataSourceOrchestrator{
		connectorService: connectorService,
		ragService:       ragService,
		logger:           logger.With("service", "data_source_orchestrator"),
	}
}

// DefaultStrategy returns a default orchestration strategy
func DefaultStrategy() *OrchestrationStrategy {
	return &OrchestrationStrategy{
		ParallelFetch:  true,
		MaxConcurrency: 5,
		UseRAG:         true,
		RAGLimit:       10,
		FetchTimeout:   30 * time.Second,
	}
}

// FetchFromAllSources orchestrates data fetching from all available sources
func (dso *DataSourceOrchestrator) FetchFromAllSources(
	ctx context.Context,
	query string,
	intent models.Intent,
	strategy *OrchestrationStrategy,
	fetchFunc func(context.Context, *models.DataConnector, string) ([]map[string]interface{}, error),
) ([]FetchResult, error) {
	if strategy == nil {
		strategy = DefaultStrategy()
	}

	dso.logger.Info("Starting orchestrated data fetch",
		"query", query,
		"intent_type", intent.Type,
		"parallel", strategy.ParallelFetch,
		"use_rag", strategy.UseRAG)

	// Get available connectors
	connectors, err := dso.selectConnectorsByIntent(ctx, intent, strategy.PreferredSources)
	if err != nil {
		return nil, fmt.Errorf("failed to select connectors: %w", err)
	}

	dso.logger.Info("Selected connectors for orchestration",
		"connector_count", len(connectors))

	var results []FetchResult
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Channel to limit concurrency
	semaphore := make(chan struct{}, strategy.MaxConcurrency)

	// Fetch from RAG if enabled
	if strategy.UseRAG && dso.ragService != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			ragResult := dso.fetchFromRAG(ctx, query, strategy.RAGLimit, strategy.FetchTimeout)

			mu.Lock()
			results = append(results, ragResult)
			mu.Unlock()
		}()
	}

	// Fetch from connectors
	for _, connector := range connectors {
		wg.Add(1)
		conn := connector // Capture loop variable

		if strategy.ParallelFetch {
			// Parallel execution
			go func() {
				defer wg.Done()
				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				result := dso.fetchFromConnector(ctx, conn, query, fetchFunc, strategy.FetchTimeout)

				mu.Lock()
				results = append(results, result)
				mu.Unlock()
			}()
		} else {
			// Sequential execution
			result := dso.fetchFromConnector(ctx, conn, query, fetchFunc, strategy.FetchTimeout)
			results = append(results, result)
			wg.Done()
		}
	}

	wg.Wait()

	// Sort results by priority
	dso.sortResultsByPriority(results)

	// Log summary
	successCount := 0
	totalRecords := 0
	for _, result := range results {
		if result.Error == nil {
			successCount++
			totalRecords += result.RecordCount
		}
	}

	dso.logger.Info("Orchestrated fetch completed",
		"total_sources", len(results),
		"successful", successCount,
		"total_records", totalRecords)

	return results, nil
}

// fetchFromRAG performs RAG vector search
func (dso *DataSourceOrchestrator) fetchFromRAG(
	ctx context.Context,
	query string,
	limit int,
	timeout time.Duration,
) FetchResult {
	start := time.Now()
	result := FetchResult{
		SourceName:  "rag_vector_search",
		SourceType:  "rag",
		Priority:    PriorityRAG,
		IsRAGResult: true,
	}

	// Create context with timeout
	fetchCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	dso.logger.Info("Fetching from RAG", "query", query, "limit", limit)

	ragResults, err := dso.ragService.GetRelevantContext(fetchCtx, query, limit)
	if err != nil {
		result.Error = err
		result.FetchTime = time.Since(start)
		dso.logger.Warn("RAG fetch failed", "error", err, "duration", result.FetchTime)
		return result
	}

	// Convert RAG results to standard format
	data := make([]map[string]interface{}, 0, len(ragResults))
	for _, ragResult := range ragResults {
		row := map[string]interface{}{
			"source":        "rag_vector_search",
			"file_id":       ragResult.FileID,
			"file_name":     ragResult.FileName,
			"table_name":    ragResult.TableName,
			"content_type":  ragResult.Type,
			"content":       ragResult.Content,
			"score":         ragResult.Score,
			"vector_id":     ragResult.VectorID,
		}

		// Add quality metrics if available
		if ragResult.QualityScore != nil {
			row["data_quality_score"] = *ragResult.QualityScore
		}
		if ragResult.MissingPercent != nil {
			row["missing_data_percent"] = *ragResult.MissingPercent
		}

		data = append(data, row)
	}

	result.Data = data
	result.RecordCount = len(data)
	result.FetchTime = time.Since(start)

	dso.logger.Info("RAG fetch successful",
		"records", result.RecordCount,
		"duration", result.FetchTime)

	return result
}

// fetchFromConnector fetches data from a specific connector
func (dso *DataSourceOrchestrator) fetchFromConnector(
	ctx context.Context,
	connector *models.DataConnector,
	query string,
	fetchFunc func(context.Context, *models.DataConnector, string) ([]map[string]interface{}, error),
	timeout time.Duration,
) FetchResult {
	start := time.Now()
	result := FetchResult{
		SourceName:  connector.Name,
		SourceType:  connector.Type,
		Priority:    dso.getPriority(connector.Type),
		IsRAGResult: false,
	}

	// Create context with timeout
	fetchCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	dso.logger.Info("Fetching from connector",
		"connector", connector.Name,
		"type", connector.Type)

	// Use provided fetch function
	data, err := fetchFunc(fetchCtx, connector, query)
	if err != nil {
		result.Error = err
		result.FetchTime = time.Since(start)
		dso.logger.Warn("Connector fetch failed",
			"connector", connector.Name,
			"error", err,
			"duration", result.FetchTime)
		return result
	}

	result.Data = data
	result.RecordCount = len(data)
	result.FetchTime = time.Since(start)

	dso.logger.Info("Connector fetch successful",
		"connector", connector.Name,
		"records", result.RecordCount,
		"duration", result.FetchTime)

	return result
}

// selectConnectorsByIntent selects connectors based on intent type
func (dso *DataSourceOrchestrator) selectConnectorsByIntent(
	ctx context.Context,
	intent models.Intent,
	preferredSources []string,
) ([]*models.DataConnector, error) {
	// If preferred sources specified, use those first
	if len(preferredSources) > 0 {
		var connectors []*models.DataConnector
		for _, id := range preferredSources {
			conn, err := dso.connectorService.GetConnector(ctx, id)
			if err == nil && conn != nil && conn.Status == models.ConnectorStatusConnected {
				connectors = append(connectors, conn)
			}
		}
		if len(connectors) > 0 {
			return connectors, nil
		}
	}

	// Get all active connectors
	activeConnectors, err := dso.connectorService.GetActiveConnectors(ctx)
	if err != nil {
		return nil, err
	}

	// Filter by intent type
	var selected []*models.DataConnector

	switch intent.Type {
	case models.IntentTypeVisualization:
		// Prefer Superset for visualizations
		for _, conn := range activeConnectors {
			if conn.Type == models.ConnectorTypeSuperset {
				selected = append(selected, conn)
			}
		}
		// Add others as fallback
		for _, conn := range activeConnectors {
			if conn.Type != models.ConnectorTypeSuperset {
				selected = append(selected, conn)
			}
		}

	case models.IntentTypeSQL:
		// Prefer file uploads and PostgreSQL for SQL queries
		for _, conn := range activeConnectors {
			if conn.Type == models.ConnectorTypeFileUpload || conn.Type == models.ConnectorTypePostgres {
				selected = append(selected, conn)
			}
		}

	case models.IntentTypeTrend, models.IntentTypeComparison, models.IntentTypeAnalytics:
		// Use all sources for comprehensive analysis
		selected = activeConnectors

	default:
		// Default: use all available sources
		selected = activeConnectors
	}

	dso.logger.Info("Connectors selected by intent",
		"intent_type", intent.Type,
		"total", len(selected))

	return selected, nil
}

// getPriority returns the priority for a connector type
func (dso *DataSourceOrchestrator) getPriority(connectorType models.ConnectorType) DataSourcePriority {
	switch connectorType {
	case models.ConnectorTypeSuperset:
		return PrioritySuperset
	case models.ConnectorTypeFileUpload:
		return PriorityFileUpload
	case models.ConnectorTypePostgres:
		return PriorityPostgres
	default:
		return PriorityFileUpload
	}
}

// sortResultsByPriority sorts fetch results by priority and success
func (dso *DataSourceOrchestrator) sortResultsByPriority(results []FetchResult) {
	// Simple bubble sort - prioritize successful results first, then by priority
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			// Successful results come first
			if results[i].Error != nil && results[j].Error == nil {
				results[i], results[j] = results[j], results[i]
			} else if (results[i].Error == nil && results[j].Error == nil) ||
				(results[i].Error != nil && results[j].Error != nil) {
				// If both successful or both failed, sort by priority
				if results[i].Priority > results[j].Priority {
					results[i], results[j] = results[j], results[i]
				}
			}
		}
	}
}

// GetSuccessfulResults returns only successful fetch results
func (dso *DataSourceOrchestrator) GetSuccessfulResults(results []FetchResult) []FetchResult {
	var successful []FetchResult
	for _, result := range results {
		if result.Error == nil && result.RecordCount > 0 {
			successful = append(successful, result)
		}
	}
	return successful
}

// MergeResults combines data from multiple fetch results
func (dso *DataSourceOrchestrator) MergeResults(results []FetchResult) ([]map[string]interface{}, map[string]interface{}) {
	allData := make(map[string]interface{})
	var combinedData []map[string]interface{}

	for _, result := range results {
		if result.Error == nil && len(result.Data) > 0 {
			allData[result.SourceName] = result.Data
			combinedData = append(combinedData, result.Data...)
		}
	}

	dso.logger.Info("Results merged",
		"sources", len(allData),
		"total_records", len(combinedData))

	return combinedData, allData
}

// DeduplicateResults removes duplicate records across multiple data sources
// using a configurable set of key fields for matching
func (dso *DataSourceOrchestrator) DeduplicateResults(
	data []map[string]interface{},
	keyFields []string,
) []map[string]interface{} {
	if len(data) == 0 {
		return data
	}

	// If no key fields specified, use common identifier fields
	if len(keyFields) == 0 {
		keyFields = []string{"id", "name", "file_id", "vector_id"}
	}

	dso.logger.Info("Starting deduplication",
		"total_records", len(data),
		"key_fields", keyFields)

	seen := make(map[string]bool)
	var deduplicated []map[string]interface{}
	duplicateCount := 0

	for _, record := range data {
		// Generate composite key from specified fields
		key := dso.generateRecordKey(record, keyFields)

		if key == "" {
			// If no key can be generated, include the record (it's unique enough)
			deduplicated = append(deduplicated, record)
			continue
		}

		if !seen[key] {
			seen[key] = true
			deduplicated = append(deduplicated, record)
		} else {
			duplicateCount++
			dso.logger.Debug("Duplicate record found",
				"key", key,
				"record", fmt.Sprintf("%v", record))
		}
	}

	dso.logger.Info("Deduplication completed",
		"original_count", len(data),
		"deduplicated_count", len(deduplicated),
		"duplicates_removed", duplicateCount)

	return deduplicated
}

// generateRecordKey creates a composite key from specified fields
// Returns empty string if no key fields found in record
func (dso *DataSourceOrchestrator) generateRecordKey(
	record map[string]interface{},
	keyFields []string,
) string {
	var keyParts []string

	for _, field := range keyFields {
		if value, exists := record[field]; exists && value != nil {
			// Convert value to string for key generation
			keyParts = append(keyParts, fmt.Sprintf("%v", value))
		}
	}

	if len(keyParts) == 0 {
		return ""
	}

	// Join key parts with separator
	return fmt.Sprintf("%s", keyParts)
}

// DeduplicateAndMerge combines deduplication and merging in one operation
func (dso *DataSourceOrchestrator) DeduplicateAndMerge(
	results []FetchResult,
	keyFields []string,
) ([]map[string]interface{}, map[string]interface{}) {
	combinedData, allData := dso.MergeResults(results)

	// Apply deduplication
	deduplicatedData := dso.DeduplicateResults(combinedData, keyFields)

	dso.logger.Info("Deduplicate and merge completed",
		"sources", len(allData),
		"original_records", len(combinedData),
		"deduplicated_records", len(deduplicatedData))

	return deduplicatedData, allData
}

// RelevanceScore represents a relevance score for a data record
type RelevanceScore struct {
	Score       float64
	Source      string
	Reasons     []string
	IsRAGResult bool
}

// CalculateRelevanceScores assigns relevance scores to records based on multiple factors
func (dso *DataSourceOrchestrator) CalculateRelevanceScores(
	data []map[string]interface{},
	query string,
) []map[string]interface{} {
	if len(data) == 0 {
		return data
	}

	dso.logger.Info("Calculating relevance scores",
		"total_records", len(data),
		"query", query)

	queryLower := strings.ToLower(query)
	scoredData := make([]map[string]interface{}, len(data))
	copy(scoredData, data)

	for i, record := range scoredData {
		score := dso.calculateRecordRelevance(record, queryLower)
		record["relevance_score"] = score.Score
		record["relevance_reasons"] = score.Reasons
		scoredData[i] = record

		dso.logger.Debug("Relevance calculated",
			"record_index", i,
			"score", score.Score,
			"source", score.Source,
			"reasons", score.Reasons)
	}

	// Sort by relevance score (highest first)
	dso.sortByRelevance(scoredData)

	dso.logger.Info("Relevance scoring completed",
		"records_scored", len(scoredData))

	return scoredData
}

// calculateRecordRelevance calculates relevance score for a single record
func (dso *DataSourceOrchestrator) calculateRecordRelevance(
	record map[string]interface{},
	queryLower string,
) RelevanceScore {
	score := 0.0
	var reasons []string
	source := ""
	isRAGResult := false

	// Extract source information
	if src, ok := record["source"].(string); ok {
		source = src
		if src == "rag_vector_search" {
			isRAGResult = true
		}
	}

	// Factor 1: RAG similarity score (if available)
	if ragScore, ok := record["score"].(float64); ok && isRAGResult {
		score += ragScore * 50.0 // RAG scores are typically 0-1, scale to 0-50
		reasons = append(reasons, fmt.Sprintf("RAG similarity: %.2f", ragScore))
	}

	// Factor 2: Data quality score (if available)
	if qualityScore, ok := record["data_quality_score"].(float64); ok {
		score += qualityScore * 20.0 // Quality scores are 0-1, scale to 0-20
		reasons = append(reasons, fmt.Sprintf("Data quality: %.2f", qualityScore))
	}

	// Factor 3: Content matching
	contentMatchScore := dso.calculateContentMatch(record, queryLower)
	if contentMatchScore > 0 {
		score += contentMatchScore
		reasons = append(reasons, fmt.Sprintf("Content match: %.2f", contentMatchScore))
	}

	// Factor 4: Source priority boost
	sourcePriorityScore := dso.getSourcePriorityScore(source)
	score += sourcePriorityScore
	if sourcePriorityScore > 0 {
		reasons = append(reasons, fmt.Sprintf("Source priority: %.2f", sourcePriorityScore))
	}

	// Factor 5: Freshness (if timestamp available)
	if timestamp, ok := record["timestamp"].(string); ok {
		freshnessScore := dso.calculateFreshnessScore(timestamp)
		score += freshnessScore
		if freshnessScore > 0 {
			reasons = append(reasons, fmt.Sprintf("Freshness: %.2f", freshnessScore))
		}
	}

	// Normalize score to 0-100 range
	if score > 100 {
		score = 100
	}

	return RelevanceScore{
		Score:       score,
		Source:      source,
		Reasons:     reasons,
		IsRAGResult: isRAGResult,
	}
}

// calculateContentMatch scores how well record content matches the query
func (dso *DataSourceOrchestrator) calculateContentMatch(
	record map[string]interface{},
	queryLower string,
) float64 {
	score := 0.0
	queryTerms := strings.Fields(queryLower)

	// Check all string fields in the record
	for _, value := range record {
		if strValue, ok := value.(string); ok {
			strLower := strings.ToLower(strValue)

			// Exact match bonus
			if strings.Contains(strLower, queryLower) {
				score += 15.0
			}

			// Term matching
			for _, term := range queryTerms {
				if len(term) > 2 && strings.Contains(strLower, term) {
					score += 3.0
				}
			}
		}
	}

	// Cap content match score at 30
	if score > 30 {
		score = 30
	}

	return score
}

// getSourcePriorityScore returns a priority boost based on source type
func (dso *DataSourceOrchestrator) getSourcePriorityScore(source string) float64 {
	switch source {
	case "rag_vector_search":
		return 10.0 // RAG gets highest priority
	case "superset":
		return 7.0
	case "file_upload":
		return 5.0
	case "postgres":
		return 3.0
	default:
		return 0.0
	}
}

// calculateFreshnessScore scores records based on timestamp recency
func (dso *DataSourceOrchestrator) calculateFreshnessScore(timestamp string) float64 {
	// Parse timestamp and calculate age
	// For now, return a base score - can be enhanced with actual time parsing
	return 5.0
}

// sortByRelevance sorts records by relevance score in descending order
func (dso *DataSourceOrchestrator) sortByRelevance(data []map[string]interface{}) {
	// Simple bubble sort based on relevance_score
	for i := 0; i < len(data); i++ {
		for j := i + 1; j < len(data); j++ {
			scoreI, okI := data[i]["relevance_score"].(float64)
			scoreJ, okJ := data[j]["relevance_score"].(float64)

			if okI && okJ && scoreI < scoreJ {
				data[i], data[j] = data[j], data[i]
			}
		}
	}
}
