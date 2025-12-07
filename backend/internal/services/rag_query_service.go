package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"insightiq/backend/internal/embedding"
	"insightiq/backend/internal/vectorstore"
)

// RAGQueryService handles RAG (Retrieval Augmented Generation) queries
// by searching vectorized file data for relevant context
type RAGQueryService struct {
	vectorStore      vectorstore.VectorStore
	embeddingService embedding.EmbeddingService
	logger           *slog.Logger
}

// RAGResult represents a result from RAG search
type RAGResult struct {
	// Vector data
	VectorID string                 `json:"vector_id"`
	Score    float64                `json:"score"`
	Metadata map[string]interface{} `json:"metadata"`

	// Extracted metadata fields for easier access
	FileID       string   `json:"file_id"`
	FileName     string   `json:"file_name"`
	TableName    string   `json:"table_name"`
	Type         string   `json:"type"` // schema, summary, row
	Content      string   `json:"content"`

	// Data quality metrics (if available)
	QualityScore     *float64 `json:"quality_score,omitempty"`
	MissingPercent   *float64 `json:"missing_percent,omitempty"`
	DuplicateRows    *int     `json:"duplicate_rows,omitempty"`

	// Column categorization (for schema type)
	NumericColumns      []string `json:"numeric_columns,omitempty"`
	CategoricalColumns  []string `json:"categorical_columns,omitempty"`
	DateColumns         []string `json:"date_columns,omitempty"`
}

// RAGSearchRequest represents a request for RAG search
type RAGSearchRequest struct {
	Query      string `json:"query"`
	Limit      int    `json:"limit"`
	MinScore   float64 `json:"min_score"`
	Collection string `json:"collection"`
	// Filter by specific metadata if needed
	FileID    string `json:"file_id,omitempty"`
	Type      string `json:"type,omitempty"` // schema, summary, row
}

// RAGSearchResponse contains the search results and metadata
type RAGSearchResponse struct {
	Results      []RAGResult `json:"results"`
	TotalResults int         `json:"total_results"`
	Query        string      `json:"query"`
	SearchTime   string      `json:"search_time"`
	Collection   string      `json:"collection"`
}

// NewRAGQueryService creates a new RAG query service
func NewRAGQueryService(
	vectorStore vectorstore.VectorStore,
	embeddingService embedding.EmbeddingService,
	logger *slog.Logger,
) *RAGQueryService {
	return &RAGQueryService{
		vectorStore:      vectorStore,
		embeddingService: embeddingService,
		logger:           logger.With("service", "rag_query"),
	}
}

// SearchFileData searches vectorized file data for relevant context
func (r *RAGQueryService) SearchFileData(ctx context.Context, query string, limit int) ([]RAGResult, error) {
	req := &RAGSearchRequest{
		Query:      query,
		Limit:      limit,
		MinScore:   0.0, // No minimum score filter by default
		Collection: "file_data_embeddings",
	}

	response, err := r.Search(ctx, req)
	if err != nil {
		return nil, err
	}

	return response.Results, nil
}

// Search performs a RAG search with more options
func (r *RAGQueryService) Search(ctx context.Context, req *RAGSearchRequest) (*RAGSearchResponse, error) {
	start := time.Now()

	r.logger.Info("Starting RAG search",
		"query", req.Query,
		"limit", req.Limit,
		"collection", req.Collection,
	)

	// 1. Generate embedding for the query
	embeddingResp, err := r.embeddingService.GenerateEmbedding(ctx, req.Query)
	if err != nil {
		r.logger.Error("Failed to generate query embedding", "error", err)
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	r.logger.Debug("Query embedding generated",
		"dimension", embeddingResp.Dimension,
		"model", embeddingResp.Model,
	)

	// 2. Search vector store
	searchResults, err := r.vectorStore.SearchVectors(ctx, req.Collection, embeddingResp.Embedding, req.Limit)
	if err != nil {
		r.logger.Error("Vector search failed", "error", err, "collection", req.Collection)
		return nil, fmt.Errorf("vector search failed: %w", err)
	}

	r.logger.Info("Vector search completed",
		"results_found", len(searchResults),
		"search_time", time.Since(start),
	)

	// 3. Convert search results to RAG results
	ragResults := make([]RAGResult, 0, len(searchResults))
	for _, result := range searchResults {
		// Apply minimum score filter
		if result.Score < req.MinScore {
			continue
		}

		// Apply metadata filters if specified
		if req.FileID != "" {
			if fileID, ok := result.Vector.Metadata["file_id"].(string); !ok || fileID != req.FileID {
				continue
			}
		}

		if req.Type != "" {
			if typeVal, ok := result.Vector.Metadata["type"].(string); !ok || typeVal != req.Type {
				continue
			}
		}

		ragResult := r.convertToRAGResult(result)
		ragResults = append(ragResults, ragResult)
	}

	searchTime := time.Since(start)

	r.logger.Info("RAG search completed",
		"total_results", len(ragResults),
		"search_time", searchTime,
	)

	return &RAGSearchResponse{
		Results:      ragResults,
		TotalResults: len(ragResults),
		Query:        req.Query,
		SearchTime:   searchTime.String(),
		Collection:   req.Collection,
	}, nil
}

// convertToRAGResult converts a vector search result to a RAG result
func (r *RAGQueryService) convertToRAGResult(result vectorstore.SearchResult) RAGResult {
	ragResult := RAGResult{
		VectorID: result.Vector.ID,
		Score:    result.Score,
		Metadata: result.Vector.Metadata,
	}

	// Extract common metadata fields
	if fileID, ok := result.Vector.Metadata["file_id"].(string); ok {
		ragResult.FileID = fileID
	}

	if fileName, ok := result.Vector.Metadata["filename"].(string); ok {
		ragResult.FileName = fileName
	}

	if tableName, ok := result.Vector.Metadata["table_name"].(string); ok {
		ragResult.TableName = tableName
	}

	if typeVal, ok := result.Vector.Metadata["type"].(string); ok {
		ragResult.Type = typeVal
	}

	if searchableText, ok := result.Vector.Metadata["searchable_text"].(string); ok {
		ragResult.Content = searchableText
	}

	// Extract data quality metrics if available
	if qualityScore, ok := result.Vector.Metadata["data_quality_score"].(float64); ok {
		ragResult.QualityScore = &qualityScore
	}

	if missingPercent, ok := result.Vector.Metadata["missing_data_percent"].(float64); ok {
		ragResult.MissingPercent = &missingPercent
	}

	if duplicateRows, ok := result.Vector.Metadata["duplicate_rows"].(float64); ok {
		// Convert float64 to int
		duplicates := int(duplicateRows)
		ragResult.DuplicateRows = &duplicates
	}

	// Extract column categorization (for schema type results)
	if numericCols, ok := result.Vector.Metadata["numeric_columns"].([]interface{}); ok {
		ragResult.NumericColumns = r.convertToStringSlice(numericCols)
	}

	if categoricalCols, ok := result.Vector.Metadata["categorical_columns"].([]interface{}); ok {
		ragResult.CategoricalColumns = r.convertToStringSlice(categoricalCols)
	}

	if dateCols, ok := result.Vector.Metadata["date_columns"].([]interface{}); ok {
		ragResult.DateColumns = r.convertToStringSlice(dateCols)
	}

	return ragResult
}

// convertToStringSlice converts []interface{} to []string
func (r *RAGQueryService) convertToStringSlice(input []interface{}) []string {
	result := make([]string, 0, len(input))
	for _, item := range input {
		if str, ok := item.(string); ok {
			result = append(result, str)
		}
	}
	return result
}

// SearchByFileID searches for vectors related to a specific file
func (r *RAGQueryService) SearchByFileID(ctx context.Context, fileID string, limit int) ([]RAGResult, error) {
	// Since we're filtering by file ID, we need a different approach
	// This is a placeholder - in production, you'd want to use metadata filters in Qdrant
	r.logger.Warn("SearchByFileID is not fully implemented - requires metadata filtering in vector store",
		"file_id", fileID,
		"limit", limit,
	)

	return nil, fmt.Errorf("SearchByFileID not yet implemented - requires vector store metadata filtering")
}

// GetSchemaContext retrieves schema information for files matching the query
func (r *RAGQueryService) GetSchemaContext(ctx context.Context, query string, limit int) ([]RAGResult, error) {
	req := &RAGSearchRequest{
		Query:      query,
		Limit:      limit,
		MinScore:   0.5, // Higher threshold for schema relevance
		Collection: "file_data_embeddings",
		Type:       "schema",
	}

	response, err := r.Search(ctx, req)
	if err != nil {
		return nil, err
	}

	r.logger.Info("Schema context retrieved",
		"query", query,
		"results", len(response.Results),
	)

	return response.Results, nil
}

// GetSummaryContext retrieves file summaries matching the query
func (r *RAGQueryService) GetSummaryContext(ctx context.Context, query string, limit int) ([]RAGResult, error) {
	req := &RAGSearchRequest{
		Query:      query,
		Limit:      limit,
		MinScore:   0.4, // Moderate threshold for summaries
		Collection: "file_data_embeddings",
		Type:       "summary",
	}

	response, err := r.Search(ctx, req)
	if err != nil {
		return nil, err
	}

	r.logger.Info("Summary context retrieved",
		"query", query,
		"results", len(response.Results),
	)

	return response.Results, nil
}

// GetRowContext retrieves specific row data matching the query
func (r *RAGQueryService) GetRowContext(ctx context.Context, query string, limit int) ([]RAGResult, error) {
	req := &RAGSearchRequest{
		Query:      query,
		Limit:      limit,
		MinScore:   0.6, // Higher threshold for specific row matches
		Collection: "file_data_embeddings",
		Type:       "row",
	}

	response, err := r.Search(ctx, req)
	if err != nil {
		return nil, err
	}

	r.logger.Info("Row context retrieved",
		"query", query,
		"results", len(response.Results),
	)

	return response.Results, nil
}

// GetRelevantContext intelligently retrieves the most relevant context
// by searching across all types (schema, summary, rows) and combining results
func (r *RAGQueryService) GetRelevantContext(ctx context.Context, query string, limit int) ([]RAGResult, error) {
	r.logger.Info("Getting relevant context across all types", "query", query)

	// Search across all types
	req := &RAGSearchRequest{
		Query:      query,
		Limit:      limit * 3, // Get more results to filter
		MinScore:   0.3,       // Lower threshold to get variety
		Collection: "file_data_embeddings",
	}

	response, err := r.Search(ctx, req)
	if err != nil {
		return nil, err
	}

	// Prioritize results: schema > summary > rows
	// This ensures we get structural understanding before specific data
	var schemaResults, summaryResults, rowResults []RAGResult

	for _, result := range response.Results {
		switch result.Type {
		case "schema":
			schemaResults = append(schemaResults, result)
		case "summary":
			summaryResults = append(summaryResults, result)
		case "row":
			rowResults = append(rowResults, result)
		}
	}

	// Combine with priority: schema first, then summary, then rows
	combinedResults := make([]RAGResult, 0, limit)

	// Add schemas (max 30% of limit)
	schemaLimit := limit / 3
	for i := 0; i < len(schemaResults) && i < schemaLimit; i++ {
		combinedResults = append(combinedResults, schemaResults[i])
	}

	// Add summaries (max 30% of limit)
	summaryLimit := limit / 3
	for i := 0; i < len(summaryResults) && i < summaryLimit; i++ {
		combinedResults = append(combinedResults, summaryResults[i])
	}

	// Add rows (remaining slots)
	rowLimit := limit - len(combinedResults)
	for i := 0; i < len(rowResults) && i < rowLimit; i++ {
		combinedResults = append(combinedResults, rowResults[i])
	}

	r.logger.Info("Relevant context retrieved",
		"total", len(combinedResults),
		"schemas", len(schemaResults),
		"summaries", len(summaryResults),
		"rows", len(rowResults),
	)

	return combinedResults, nil
}
