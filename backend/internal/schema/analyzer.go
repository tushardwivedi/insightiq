package schema

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"time"

	"insightiq/backend/internal/connectors"
)

// AnalyzerService handles business context analysis and domain generation
type AnalyzerService struct {
	scannerService *ScannerService
	llmConn       *connectors.OllamaConnector
	logger        *slog.Logger
}

// NewAnalyzerService creates a new business context analyzer
func NewAnalyzerService(scannerService *ScannerService, llmConn *connectors.OllamaConnector, logger *slog.Logger) *AnalyzerService {
	return &AnalyzerService{
		scannerService: scannerService,
		llmConn:       llmConn,
		logger:        logger,
	}
}

// AnalyzeBusinessContext performs comprehensive business context analysis
func (a *AnalyzerService) AnalyzeBusinessContext(ctx context.Context, connectorID string) (*SchemaContext, error) {
	a.logger.Info("Starting business context analysis", "connector_id", connectorID)

	// Step 1: Scan schema structure
	schemaContext, err := a.scannerService.ScanDataSource(ctx, connectorID)
	if err != nil {
		return nil, fmt.Errorf("failed to scan schema: %w", err)
	}

	// Step 2: Classify primary domain
	primaryDomain, confidence := a.classifyPrimaryDomain(schemaContext)
	schemaContext.PrimaryDomain = primaryDomain
	schemaContext.Confidence = confidence

	// Step 3: Generate domain contexts
	domainContexts, err := a.generateDomainContexts(ctx, schemaContext)
	if err != nil {
		return nil, fmt.Errorf("failed to generate domain contexts: %w", err)
	}
	schemaContext.DetectedDomains = domainContexts

	// Step 4: Enhance sample queries with LLM
	enhancedQueries, err := a.enhanceSampleQueries(ctx, schemaContext)
	if err != nil {
		a.logger.Warn("Failed to enhance queries with LLM, using generated queries", "error", err)
	} else {
		schemaContext.SampleQueries = enhancedQueries
	}

	a.logger.Info("Business context analysis completed",
		"connector_id", connectorID,
		"primary_domain", primaryDomain,
		"confidence", confidence,
		"domains", len(domainContexts))

	return schemaContext, nil
}

// classifyPrimaryDomain determines the primary business domain
func (a *AnalyzerService) classifyPrimaryDomain(schemaContext *SchemaContext) (Domain, float64) {
	domainScores := make(map[Domain]float64)

	// Analyze table patterns for domain classification
	for _, table := range schemaContext.Tables {
		tableDomain := a.classifyTableDomain(table)
		domainScores[tableDomain] += 1.0

		// Weight by number of business metrics in table
		metricCount := 0
		for _, col := range table.Columns {
			if col.IsMetric {
				metricCount++
			}
		}
		domainScores[tableDomain] += float64(metricCount) * 0.5
	}

	// Analyze business metrics for additional domain signals
	for _, metric := range schemaContext.BusinessMetrics {
		domainScores[metric.Domain] += 0.5
	}

	// Find domain with highest score
	var bestDomain Domain
	var maxScore float64
	var totalScore float64

	for domain, score := range domainScores {
		totalScore += score
		if score > maxScore {
			maxScore = score
			bestDomain = domain
		}
	}

	// Calculate confidence
	confidence := 0.5 // baseline confidence
	if totalScore > 0 {
		confidence = maxScore / totalScore
	}

	// Ensure minimum confidence for non-general domains
	if bestDomain == DomainGeneral && confidence < 0.7 {
		confidence = 0.5
	}

	return bestDomain, confidence
}

// classifyTableDomain classifies a single table's domain
func (a *AnalyzerService) classifyTableDomain(table TableContext) Domain {
	tableName := strings.ToLower(table.TableName)

	// Domain classification patterns
	domainPatterns := map[Domain][]string{
		DomainSales: {
			"order", "sale", "transaction", "purchase", "invoice",
			"receipt", "payment", "billing", "checkout",
		},
		DomainMarketing: {
			"campaign", "lead", "prospect", "marketing", "advertisement",
			"email", "newsletter", "promotion", "coupon", "affiliate",
		},
		DomainCustomer: {
			"customer", "user", "client", "account", "profile",
			"contact", "member", "subscriber",
		},
		DomainProduct: {
			"product", "item", "catalog", "inventory", "stock",
			"sku", "variant", "category", "brand",
		},
		DomainFinance: {
			"finance", "accounting", "revenue", "expense", "budget",
			"cost", "profit", "tax", "payroll", "investment",
		},
		DomainOperations: {
			"operation", "process", "workflow", "task", "job",
			"schedule", "resource", "asset", "facility",
		},
		DomainHR: {
			"employee", "staff", "personnel", "hr", "human",
			"department", "position", "salary", "benefit",
		},
	}

	// Check table name against patterns
	for domain, patterns := range domainPatterns {
		for _, pattern := range patterns {
			if strings.Contains(tableName, pattern) {
				return domain
			}
		}
	}

	// Check business tags
	for _, tag := range table.BusinessTags {
		tag = strings.ToLower(tag)
		for domain, patterns := range domainPatterns {
			for _, pattern := range patterns {
				if strings.Contains(tag, pattern) || pattern == tag {
					return domain
				}
			}
		}
	}

	// Analyze column patterns
	for _, col := range table.Columns {
		colName := strings.ToLower(col.Name)
		for domain, patterns := range domainPatterns {
			for _, pattern := range patterns {
				if strings.Contains(colName, pattern) {
					return domain
				}
			}
		}
	}

	return DomainGeneral
}

// generateDomainContexts creates comprehensive domain contexts from schema analysis
func (a *AnalyzerService) generateDomainContexts(ctx context.Context, schemaContext *SchemaContext) ([]DomainContext, error) {
	domainMap := make(map[Domain]*DomainContext)

	// Group tables by domain
	for _, table := range schemaContext.Tables {
		domain := a.classifyTableDomain(table)
		table.Domain = domain

		if _, exists := domainMap[domain]; !exists {
			domainMap[domain] = &DomainContext{
				Domain:        domain,
				AutoGenerated: true,
				LastUpdated:   time.Now(),
				Tables:        []TableContext{},
				Keywords:      []string{},
				Metrics:       []string{},
				Dimensions:    []string{},
			}
		}

		domainMap[domain].Tables = append(domainMap[domain].Tables, table)
	}

	// Generate context for each domain
	var domainContexts []DomainContext
	for _, context := range domainMap {
		a.enrichDomainContext(context)
		domainContexts = append(domainContexts, *context)
	}

	// Sort by relevance (number of tables)
	sort.Slice(domainContexts, func(i, j int) bool {
		return len(domainContexts[i].Tables) > len(domainContexts[j].Tables)
	})

	return domainContexts, nil
}

// enrichDomainContext adds business intelligence to a domain context
func (a *AnalyzerService) enrichDomainContext(context *DomainContext) {
	// Generate description based on domain type
	context.Description = a.generateDomainDescription(context.Domain, context.Tables)

	// Extract keywords from table names and columns
	keywordSet := make(map[string]bool)
	metricSet := make(map[string]bool)
	dimensionSet := make(map[string]bool)

	for _, table := range context.Tables {
		// Add table name as keyword
		keywordSet[strings.ToLower(table.TableName)] = true

		// Add business tags
		for _, tag := range table.BusinessTags {
			keywordSet[strings.ToLower(tag)] = true
		}

		// Extract metrics and dimensions
		for _, col := range table.Columns {
			if col.IsMetric {
				metricSet[col.Name] = true
				keywordSet[strings.ToLower(col.Name)] = true
			}
			if col.IsDimension {
				dimensionSet[col.Name] = true
				keywordSet[strings.ToLower(col.Name)] = true
			}
		}
	}

	// Convert sets to slices
	for keyword := range keywordSet {
		context.Keywords = append(context.Keywords, keyword)
	}
	for metric := range metricSet {
		context.Metrics = append(context.Metrics, metric)
	}
	for dimension := range dimensionSet {
		context.Dimensions = append(context.Dimensions, dimension)
	}

	// Generate business glossary
	context.Glossary = a.generateBusinessGlossary(context.Domain, context.Keywords)

	// Generate query patterns
	context.Patterns = a.generateQueryPatterns(context.Domain, context.Keywords, context.Metrics)

	// Calculate confidence based on data richness
	context.Confidence = a.calculateDomainConfidence(context)
}

// generateDomainDescription creates a business description for the domain
func (a *AnalyzerService) generateDomainDescription(domain Domain, tables []TableContext) string {
	tableCount := len(tables)

	baseDescriptions := map[Domain]string{
		DomainSales:      "Sales and revenue analytics, order management, customer transactions",
		DomainMarketing:  "Marketing campaign performance, lead generation, customer acquisition",
		DomainCustomer:   "Customer relationship management, user profiles, engagement metrics",
		DomainProduct:    "Product catalog, inventory management, pricing and categorization",
		DomainFinance:    "Financial reporting, accounting, revenue and expense tracking",
		DomainOperations: "Operational efficiency, process management, resource utilization",
		DomainHR:         "Human resources, employee management, organizational analytics",
		DomainGeneral:    "General business data and analytics",
	}

	description := baseDescriptions[domain]
	if description == "" {
		description = "Business data and analytics"
	}

	// Add context about data richness
	if tableCount > 1 {
		description += fmt.Sprintf(" with %d related data sources", tableCount)
	}

	return description
}

// generateBusinessGlossary creates business term definitions
func (a *AnalyzerService) generateBusinessGlossary(domain Domain, keywords []string) []BusinessGlossary {
	var glossary []BusinessGlossary

	// Predefined business terms by domain
	domainGlossary := map[Domain][]BusinessGlossary{
		DomainSales: {
			{
				Term:       "Revenue",
				Definition: "Total income generated from sales transactions",
				Domain:     domain,
				Synonyms:   []string{"income", "sales", "earnings"},
			},
			{
				Term:       "Conversion Rate",
				Definition: "Percentage of prospects that become customers",
				Domain:     domain,
				Synonyms:   []string{"conversion", "close rate"},
			},
		},
		DomainMarketing: {
			{
				Term:       "CAC",
				Definition: "Customer Acquisition Cost - cost to acquire one customer",
				Domain:     domain,
				Synonyms:   []string{"acquisition cost", "customer cost"},
			},
			{
				Term:       "ROAS",
				Definition: "Return on Ad Spend - revenue per dollar spent on advertising",
				Domain:     domain,
				Synonyms:   []string{"ad return", "advertising roi"},
			},
		},
		DomainCustomer: {
			{
				Term:       "LTV",
				Definition: "Customer Lifetime Value - total revenue from customer relationship",
				Domain:     domain,
				Synonyms:   []string{"lifetime value", "customer value"},
			},
			{
				Term:       "Churn Rate",
				Definition: "Percentage of customers who stop using service",
				Domain:     domain,
				Synonyms:   []string{"churn", "attrition rate"},
			},
		},
	}

	if terms, exists := domainGlossary[domain]; exists {
		glossary = append(glossary, terms...)
	}

	return glossary
}

// generateQueryPatterns creates common query patterns for the domain
func (a *AnalyzerService) generateQueryPatterns(domain Domain, keywords []string, metrics []string) []QueryPattern {
	var patterns []QueryPattern

	// Generate patterns based on domain
	switch domain {
	case DomainSales:
		patterns = append(patterns, QueryPattern{
			Pattern:     "sales performance analysis",
			Domain:      domain,
			Intent:      IntentAnalytics,
			Keywords:    []string{"sales", "performance", "revenue", "orders"},
			Description: "Analysis of sales metrics and performance trends",
			Examples:    []string{"show me sales data", "sales performance", "revenue analysis"},
		})

	case DomainMarketing:
		patterns = append(patterns, QueryPattern{
			Pattern:     "marketing campaign analysis",
			Domain:      domain,
			Intent:      IntentAnalytics,
			Keywords:    []string{"marketing", "campaign", "leads", "conversion"},
			Description: "Analysis of marketing campaign effectiveness",
			Examples:    []string{"campaign performance", "lead generation", "marketing roi"},
		})

	case DomainCustomer:
		patterns = append(patterns, QueryPattern{
			Pattern:     "customer behavior analysis",
			Domain:      domain,
			Intent:      IntentAnalytics,
			Keywords:    []string{"customer", "behavior", "engagement", "retention"},
			Description: "Analysis of customer behavior and engagement patterns",
			Examples:    []string{"customer insights", "user engagement", "customer segmentation"},
		})
	}

	// Add trend analysis pattern for all domains
	if len(metrics) > 0 {
		patterns = append(patterns, QueryPattern{
			Pattern:     fmt.Sprintf("%s trend analysis", strings.ToLower(string(domain))),
			Domain:      domain,
			Intent:      IntentTrend,
			Keywords:    append(keywords[:min(3, len(keywords))], "trend", "over time"),
			Description: fmt.Sprintf("Time-based trend analysis for %s metrics", strings.ToLower(string(domain))),
			Examples:    []string{fmt.Sprintf("%s trends", strings.ToLower(string(domain))), "trends over time"},
		})
	}

	return patterns
}

// calculateDomainConfidence calculates confidence score for domain classification
func (a *AnalyzerService) calculateDomainConfidence(context *DomainContext) float64 {
	score := 0.0

	// Base score for having tables
	score += float64(len(context.Tables)) * 0.2

	// Score for having metrics
	score += float64(len(context.Metrics)) * 0.1

	// Score for having dimensions
	score += float64(len(context.Dimensions)) * 0.05

	// Score for domain-specific keywords
	domainKeywords := map[Domain][]string{
		DomainSales:     {"sales", "order", "revenue", "customer"},
		DomainMarketing: {"campaign", "lead", "marketing", "conversion"},
		DomainCustomer:  {"customer", "user", "profile", "engagement"},
		DomainProduct:   {"product", "inventory", "catalog", "price"},
		DomainFinance:   {"finance", "payment", "invoice", "revenue"},
	}

	if expectedKeywords, exists := domainKeywords[context.Domain]; exists {
		matchCount := 0
		for _, expected := range expectedKeywords {
			for _, actual := range context.Keywords {
				if strings.Contains(actual, expected) {
					matchCount++
					break
				}
			}
		}
		score += float64(matchCount) / float64(len(expectedKeywords)) * 0.3
	}

	// Normalize to 0-1 range
	if score > 1.0 {
		score = 1.0
	}

	// Ensure minimum confidence
	if score < 0.3 {
		score = 0.3
	}

	return score
}

// enhanceSampleQueries uses LLM to improve and expand sample queries
func (a *AnalyzerService) enhanceSampleQueries(ctx context.Context, schemaContext *SchemaContext) ([]string, error) {
	if a.llmConn == nil {
		return schemaContext.SampleQueries, nil // Return original queries if no LLM
	}

	// Create context prompt for LLM
	prompt := a.buildQueryEnhancementPrompt(schemaContext)

	// Call LLM to enhance queries (this would be a real LLM call in production)
	response, err := a.llmConn.GenerateResponse(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to enhance queries with LLM: %w", err)
	}

	// Parse LLM response to extract enhanced queries
	enhancedQueries := a.parseLLMQueryResponse(response)

	// Fallback to original queries if parsing fails
	if len(enhancedQueries) == 0 {
		return schemaContext.SampleQueries, nil
	}

	return enhancedQueries, nil
}

// buildQueryEnhancementPrompt creates a prompt for query enhancement
func (a *AnalyzerService) buildQueryEnhancementPrompt(schemaContext *SchemaContext) string {
	var builder strings.Builder

	builder.WriteString("Given this business data schema, generate natural language queries that users might ask:\n\n")

	builder.WriteString(fmt.Sprintf("Primary Domain: %s\n", schemaContext.PrimaryDomain))
	builder.WriteString(fmt.Sprintf("Tables: %s\n", a.extractTableNames(schemaContext.Tables)))
	builder.WriteString(fmt.Sprintf("Key Metrics: %s\n", strings.Join(a.extractMetricNames(schemaContext.BusinessMetrics), ", ")))

	builder.WriteString("\nGenerate 5-8 sample queries that business users would naturally ask about this data.")

	return builder.String()
}

// extractTableNames extracts table names for prompt
func (a *AnalyzerService) extractTableNames(tables []TableContext) string {
	var names []string
	for _, table := range tables {
		names = append(names, table.TableName)
	}
	return strings.Join(names, ", ")
}

// extractMetricNames extracts metric names for prompt
func (a *AnalyzerService) extractMetricNames(metrics []BusinessMetric) []string {
	var names []string
	for _, metric := range metrics {
		names = append(names, metric.Name)
	}
	return names
}

// parseLLMQueryResponse parses the LLM response to extract queries
func (a *AnalyzerService) parseLLMQueryResponse(response string) []string {
	// Simple parsing - in production this would be more sophisticated
	lines := strings.Split(response, "\n")
	var queries []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 10 && (strings.Contains(line, "?") || strings.Contains(line, "show") || strings.Contains(line, "what")) {
			queries = append(queries, line)
		}
	}

	return queries
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}