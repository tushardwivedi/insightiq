package intent

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"insightiq/backend/internal/embedding"
	"insightiq/backend/internal/schema"
	"insightiq/backend/internal/vectorstore"
)

// ClassificationResult represents the result of intent classification
type ClassificationResult struct {
	Domain     schema.Domain `json:"domain"`
	Intent     schema.Intent `json:"intent"`
	Confidence float64       `json:"confidence"`
	Keywords   []string      `json:"keywords"`
	Reasoning  string        `json:"reasoning"`
	Context    []string      `json:"context,omitempty"`
}

// ClassificationService provides intent classification using RAG
type ClassificationService struct {
	vectorStore      vectorstore.VectorStore
	embeddingService embedding.EmbeddingService
	logger           *slog.Logger
}

// NewClassificationService creates a new intent classification service
func NewClassificationService(
	vectorStore vectorstore.VectorStore,
	embeddingService embedding.EmbeddingService,
	logger *slog.Logger,
) *ClassificationService {
	return &ClassificationService{
		vectorStore:      vectorStore,
		embeddingService: embeddingService,
		logger:           logger,
	}
}

// ClassifyQuery classifies a user query into domain and intent
func (s *ClassificationService) ClassifyQuery(ctx context.Context, query string) (*ClassificationResult, error) {
	if query == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}

	s.logger.Info("Classifying query", "query", query)

	// Step 1: Generate embedding for the query
	embeddingResp, err := s.embeddingService.GenerateEmbedding(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	// Step 2: Search for similar domain contexts
	searchResults, err := s.vectorStore.SearchVectors(ctx, "domain_contexts", embeddingResp.Embedding, 3)
	if err != nil {
		s.logger.Warn("Vector search failed, falling back to keyword matching", "error", err)
		return s.fallbackClassification(query), nil
	}

	// Step 3: Analyze search results and classify
	result := s.analyzeSearchResults(query, searchResults)

	s.logger.Info("Query classified",
		"query", query,
		"domain", result.Domain,
		"intent", result.Intent,
		"confidence", result.Confidence)

	return result, nil
}

// analyzeSearchResults analyzes vector search results to determine classification
func (s *ClassificationService) analyzeSearchResults(query string, results []vectorstore.SearchResult) *ClassificationResult {
	queryLower := strings.ToLower(query)

	if len(results) == 0 {
		return s.fallbackClassification(query)
	}

	// Get the best match
	bestMatch := results[0]
	domain := schema.Domain(bestMatch.Vector.Metadata["domain"].(string))
	confidence := bestMatch.Score

	// Determine intent based on query patterns
	intent := s.determineIntent(queryLower)

	// Extract keywords from the query
	keywords := s.extractKeywords(queryLower)

	// Create reasoning
	reasoning := fmt.Sprintf("Matched domain '%s' with confidence %.2f based on vector similarity", domain, confidence)

	// Add context from search results
	var context []string
	for i, result := range results {
		if i >= 2 { // Limit context to top 2 results
			break
		}
		if desc, ok := result.Vector.Metadata["description"].(string); ok {
			context = append(context, desc)
		}
	}

	return &ClassificationResult{
		Domain:     domain,
		Intent:     intent,
		Confidence: confidence,
		Keywords:   keywords,
		Reasoning:  reasoning,
		Context:    context,
	}
}

// determineIntent determines the intent type based on query patterns
func (s *ClassificationService) determineIntent(queryLower string) schema.Intent {
	// Analytics patterns
	if strings.Contains(queryLower, "insight") || strings.Contains(queryLower, "analysis") ||
		strings.Contains(queryLower, "analytics") || strings.Contains(queryLower, "show me") {
		return schema.IntentAnalytics
	}

	// Comparison patterns
	if strings.Contains(queryLower, "compare") || strings.Contains(queryLower, "vs") ||
		strings.Contains(queryLower, "difference") || strings.Contains(queryLower, "between") {
		return schema.IntentComparison
	}

	// Trend patterns
	if strings.Contains(queryLower, "trend") || strings.Contains(queryLower, "over time") ||
		strings.Contains(queryLower, "growth") || strings.Contains(queryLower, "change") {
		return schema.IntentTrend
	}

	// Visualization patterns
	if strings.Contains(queryLower, "chart") || strings.Contains(queryLower, "graph") ||
		strings.Contains(queryLower, "visualize") || strings.Contains(queryLower, "dashboard") {
		return schema.IntentVisualization
	}

	// Summary patterns
	if strings.Contains(queryLower, "summary") || strings.Contains(queryLower, "overview") ||
		strings.Contains(queryLower, "total") || strings.Contains(queryLower, "aggregate") {
		return schema.IntentSummary
	}

	// Default to analytics
	return schema.IntentAnalytics
}

// extractKeywords extracts relevant keywords from the query
func (s *ClassificationService) extractKeywords(queryLower string) []string {
	// Define keyword categories
	businessKeywords := []string{
		"sales", "revenue", "profit", "orders", "customers",
		"marketing", "campaign", "leads", "conversion",
		"finance", "cost", "budget", "expense", "accounting",
		"operations", "process", "efficiency", "productivity",
		"hr", "employee", "staff", "personnel",
		"product", "inventory", "catalog", "pricing",
		"game", "gaming", "video game", "entertainment",
		"slack", "channel", "messages", "communication",
		"covid", "vaccine", "vaccination", "pandemic",
		"birth", "names", "population", "demographics",
	}

	var foundKeywords []string
	for _, keyword := range businessKeywords {
		if strings.Contains(queryLower, keyword) {
			foundKeywords = append(foundKeywords, keyword)
		}
	}

	return foundKeywords
}

// fallbackClassification provides keyword-based classification when vector search fails
func (s *ClassificationService) fallbackClassification(query string) *ClassificationResult {
	queryLower := strings.ToLower(query)

	// Keyword-based domain classification
	var domain schema.Domain
	var confidence float64

	switch {
	case strings.Contains(queryLower, "finance") || strings.Contains(queryLower, "financial") ||
		strings.Contains(queryLower, "accounting") || strings.Contains(queryLower, "budget"):
		domain = schema.DomainFinance
		confidence = 0.8

	case strings.Contains(queryLower, "game") || strings.Contains(queryLower, "gaming") ||
		strings.Contains(queryLower, "video"):
		domain = schema.DomainGaming
		confidence = 0.8

	case strings.Contains(queryLower, "slack"):
		domain = schema.DomainSlack
		confidence = 0.9

	case strings.Contains(queryLower, "covid") || strings.Contains(queryLower, "vaccine"):
		domain = schema.DomainCOVID
		confidence = 0.8

	case strings.Contains(queryLower, "birth") || strings.Contains(queryLower, "name"):
		domain = schema.DomainBirthNames
		confidence = 0.8

	case strings.Contains(queryLower, "sales") || strings.Contains(queryLower, "revenue") ||
		strings.Contains(queryLower, "orders"):
		domain = schema.DomainSales
		confidence = 0.7

	default:
		domain = schema.DomainSales // Default fallback
		confidence = 0.5
	}

	intent := s.determineIntent(queryLower)
	keywords := s.extractKeywords(queryLower)

	return &ClassificationResult{
		Domain:     domain,
		Intent:     intent,
		Confidence: confidence,
		Keywords:   keywords,
		Reasoning:  fmt.Sprintf("Keyword-based classification (fallback) - matched '%s'", domain),
	}
}

// Health checks the health of all dependent services
func (s *ClassificationService) Health(ctx context.Context) error {
	// Check embedding service
	if err := s.embeddingService.Health(ctx); err != nil {
		return fmt.Errorf("embedding service unhealthy: %w", err)
	}

	// Check vector store (if it has a health method)
	// Note: The current interface doesn't include Health, but we can add it later

	s.logger.Debug("Intent classification service health check passed")
	return nil
}