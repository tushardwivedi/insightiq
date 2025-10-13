package schema

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"insightiq/backend/internal/embedding"
	"insightiq/backend/internal/vectorstore"
)

// IngestionService handles schema and domain context ingestion into vector store
type IngestionService struct {
	vectorStore         vectorstore.VectorStore
	embeddingService    embedding.EmbeddingService
	domainGenerator     *DomainGeneratorService
	connectorService    ConnectorService
	logger              *slog.Logger
}

// NewIngestionService creates a new schema ingestion service
func NewIngestionService(
	vectorStore vectorstore.VectorStore,
	embeddingService embedding.EmbeddingService,
	logger *slog.Logger,
) *IngestionService {
	return &IngestionService{
		vectorStore:      vectorStore,
		embeddingService: embeddingService,
		logger:           logger,
	}
}

// NewEnhancedIngestionService creates an ingestion service with dynamic capabilities
func NewEnhancedIngestionService(
	vectorStore vectorstore.VectorStore,
	embeddingService embedding.EmbeddingService,
	domainGenerator *DomainGeneratorService,
	connectorService ConnectorService,
	logger *slog.Logger,
) *IngestionService {
	return &IngestionService{
		vectorStore:      vectorStore,
		embeddingService: embeddingService,
		domainGenerator:  domainGenerator,
		connectorService: connectorService,
		logger:           logger,
	}
}

// IngestDomainContexts ingests predefined domain contexts into the vector store
func (s *IngestionService) IngestDomainContexts(ctx context.Context) error {
	s.logger.Info("Starting domain context ingestion")

	// Get predefined domain contexts
	contexts := s.getPredefinedDomainContexts()

	// Create collection for domain contexts
	dimension := s.embeddingService.GetDimension()
	if err := s.vectorStore.CreateCollection(ctx, "domain_contexts", dimension); err != nil {
		s.logger.Warn("Failed to create domain_contexts collection (may already exist)", "error", err)
	}

	// Ingest each domain context
	for _, domainCtx := range contexts {
		if err := s.ingestDomainContext(ctx, domainCtx); err != nil {
			s.logger.Error("Failed to ingest domain context",
				"domain", domainCtx.Domain,
				"error", err)
			continue
		}
		s.logger.Info("Successfully ingested domain context", "domain", domainCtx.Domain)
	}

	s.logger.Info("Domain context ingestion completed", "domains", len(contexts))
	return nil
}

// ingestDomainContext ingests a single domain context
func (s *IngestionService) ingestDomainContext(ctx context.Context, domainCtx DomainContext) error {
	// Create searchable text representation of the domain
	searchText := s.createDomainSearchText(domainCtx)

	// Generate embedding for the domain context
	embeddingResp, err := s.embeddingService.GenerateEmbedding(ctx, searchText)
	if err != nil {
		return fmt.Errorf("failed to generate embedding for domain %s: %w", domainCtx.Domain, err)
	}

	// Create vector with metadata
	vector := vectorstore.Vector{
		ID:     fmt.Sprintf("domain_%s", domainCtx.Domain),
		Values: embeddingResp.Embedding,
		Metadata: map[string]interface{}{
			"type":        "domain_context",
			"domain":      string(domainCtx.Domain),
			"description": domainCtx.Description,
			"keywords":    domainCtx.Keywords,
			"metrics":     domainCtx.Metrics,
			"dimensions":  domainCtx.Dimensions,
			"table_count": len(domainCtx.Tables),
			"updated_at":  domainCtx.LastUpdated.Format(time.RFC3339),
		},
	}

	// Upsert the vector
	return s.vectorStore.UpsertVector(ctx, "domain_contexts", vector)
}

// createDomainSearchText creates a searchable text representation of a domain
func (s *IngestionService) createDomainSearchText(domainCtx DomainContext) string {
	var parts []string

	// Add domain name and description
	parts = append(parts, string(domainCtx.Domain))
	parts = append(parts, domainCtx.Description)

	// Add keywords
	parts = append(parts, strings.Join(domainCtx.Keywords, " "))

	// Add table names and descriptions
	for _, table := range domainCtx.Tables {
		parts = append(parts, table.TableName)
		parts = append(parts, table.Description)
		parts = append(parts, strings.Join(table.BusinessTags, " "))
	}

	// Add glossary terms
	for _, glossary := range domainCtx.Glossary {
		parts = append(parts, glossary.Term)
		parts = append(parts, glossary.Definition)
		parts = append(parts, strings.Join(glossary.Synonyms, " "))
	}

	// Add query patterns
	for _, pattern := range domainCtx.Patterns {
		parts = append(parts, pattern.Description)
		parts = append(parts, strings.Join(pattern.Keywords, " "))
		parts = append(parts, strings.Join(pattern.Examples, " "))
	}

	return strings.Join(parts, " ")
}

// IngestDynamicDomainContexts ingests dynamically generated domain contexts
func (s *IngestionService) IngestDynamicDomainContexts(ctx context.Context, connectorID string) error {
	if s.domainGenerator == nil {
		return fmt.Errorf("domain generator not available - use NewEnhancedIngestionService")
	}

	s.logger.Info("Starting dynamic domain context ingestion", "connector_id", connectorID)

	// Generate and ingest domain contexts for the connector
	return s.domainGenerator.GenerateAndIngestDomainContexts(ctx, connectorID)
}

// IngestAllConnectorContexts ingests contexts for all available connectors
func (s *IngestionService) IngestAllConnectorContexts(ctx context.Context) error {
	if s.connectorService == nil {
		return fmt.Errorf("connector service not available - use NewEnhancedIngestionService")
	}

	s.logger.Info("Starting ingestion for all connector contexts")

	// Get all connected connectors
	connectors, err := s.connectorService.ListConnectors(ctx)
	if err != nil {
		return fmt.Errorf("failed to list connectors: %w", err)
	}

	var errors []string
	successCount := 0

	for _, connector := range connectors {
		if connector.Status != "connected" {
			s.logger.Info("Skipping disconnected connector", "connector_id", connector.ID, "status", connector.Status)
			continue
		}

		if err := s.IngestDynamicDomainContexts(ctx, connector.ID); err != nil {
			s.logger.Error("Failed to ingest contexts for connector",
				"connector_id", connector.ID,
				"connector_name", connector.Name,
				"error", err)
			errors = append(errors, fmt.Sprintf("%s: %v", connector.Name, err))
		} else {
			successCount++
			s.logger.Info("Successfully ingested contexts for connector",
				"connector_id", connector.ID,
				"connector_name", connector.Name)
		}
	}

	s.logger.Info("Completed ingestion for all connectors",
		"total_connectors", len(connectors),
		"successful", successCount,
		"failed", len(errors))

	if len(errors) > 0 {
		return fmt.Errorf("failed to ingest contexts for %d connectors: %s", len(errors), strings.Join(errors, "; "))
	}

	return nil
}

// RefreshConnectorContext refreshes domain contexts for a specific connector
func (s *IngestionService) RefreshConnectorContext(ctx context.Context, connectorID string) error {
	if s.domainGenerator == nil {
		return fmt.Errorf("domain generator not available - use NewEnhancedIngestionService")
	}

	s.logger.Info("Refreshing connector context", "connector_id", connectorID)
	return s.domainGenerator.RefreshDomainContexts(ctx, connectorID)
}

// getPredefinedDomainContexts returns predefined domain contexts
func (s *IngestionService) getPredefinedDomainContexts() []DomainContext {
	now := time.Now()

	return []DomainContext{
		{
			Domain:      DomainSales,
			Description: "Sales and revenue analytics, customer orders, product performance",
			Keywords:    []string{"sales", "revenue", "orders", "customers", "products", "performance", "conversion"},
			Metrics:     []string{"revenue", "orders", "conversion_rate", "avg_order_value"},
			Dimensions:  []string{"product_category", "customer_segment", "time_period"},
			Tables: []TableContext{
				{
					TableName:    "orders",
					Domain:       DomainSales,
					Description:  "Customer order transactions",
					BusinessTags: []string{"sales", "transactions", "customers"},
				},
				{
					TableName:    "products",
					Domain:       DomainSales,
					Description:  "Product catalog and pricing",
					BusinessTags: []string{"products", "pricing", "catalog"},
				},
			},
			Glossary: []BusinessGlossary{
				{
					Term:       "Revenue",
					Definition: "Total income from sales",
					Domain:     DomainSales,
					Synonyms:   []string{"income", "sales", "earnings"},
				},
			},
			Patterns: []QueryPattern{
				{
					Pattern:     "sales performance analysis",
					Domain:      DomainSales,
					Intent:      IntentAnalytics,
					Keywords:    []string{"sales", "performance", "revenue"},
					Description: "Analysis of sales metrics and performance",
					Examples:    []string{"show me sales data", "sales performance", "revenue analysis"},
				},
			},
			LastUpdated: now,
		},
		{
			Domain:      DomainMarketing,
			Description: "Marketing campaign performance, lead generation, customer acquisition analytics",
			Keywords:    []string{"marketing", "campaign", "leads", "conversion", "acquisition", "cac", "roas", "attribution"},
			Metrics:     []string{"cac", "roas", "conversion_rate", "leads", "impressions", "clicks", "ctr"},
			Dimensions:  []string{"campaign", "channel", "audience", "time_period"},
			Tables: []TableContext{
				{
					TableName:    "campaigns",
					Domain:       DomainMarketing,
					Description:  "Marketing campaign data and performance",
					BusinessTags: []string{"campaigns", "marketing", "advertising"},
				},
				{
					TableName:    "leads",
					Domain:       DomainMarketing,
					Description:  "Lead generation and conversion data",
					BusinessTags: []string{"leads", "prospects", "conversion"},
				},
			},
			Glossary: []BusinessGlossary{
				{
					Term:       "CAC",
					Definition: "Customer Acquisition Cost - cost to acquire one customer",
					Domain:     DomainMarketing,
					Synonyms:   []string{"acquisition cost", "customer cost"},
				},
				{
					Term:       "ROAS",
					Definition: "Return on Ad Spend - revenue per dollar spent on advertising",
					Domain:     DomainMarketing,
					Synonyms:   []string{"ad return", "advertising roi"},
				},
			},
			Patterns: []QueryPattern{
				{
					Pattern:     "marketing performance analysis",
					Domain:      DomainMarketing,
					Intent:      IntentAnalytics,
					Keywords:    []string{"marketing", "campaign", "performance", "conversion"},
					Description: "Analysis of marketing campaign effectiveness and ROI",
					Examples:    []string{"campaign performance", "marketing roi", "lead generation"},
				},
			},
			LastUpdated: now,
		},
		{
			Domain:      DomainFinance,
			Description: "Financial reporting, accounting, revenue and expense tracking",
			Keywords:    []string{"finance", "accounting", "revenue", "expense", "profit", "cash", "budget"},
			Metrics:     []string{"revenue", "profit", "expenses", "cash_flow", "margin", "roi"},
			Dimensions:  []string{"account", "category", "period", "department"},
			Tables: []TableContext{
				{
					TableName:    "transactions",
					Domain:       DomainFinance,
					Description:  "Financial transaction records",
					BusinessTags: []string{"finance", "transactions", "accounting"},
				},
			},
			Glossary: []BusinessGlossary{
				{
					Term:       "ROI",
					Definition: "Return on Investment - profit relative to cost of investment",
					Domain:     DomainFinance,
					Synonyms:   []string{"return on investment", "roi"},
				},
			},
			Patterns: []QueryPattern{
				{
					Pattern:     "financial analysis",
					Domain:      DomainFinance,
					Intent:      IntentAnalytics,
					Keywords:    []string{"financial", "revenue", "profit", "expenses"},
					Description: "Analysis of financial performance and profitability",
					Examples:    []string{"financial dashboard", "revenue analysis", "profit margins"},
				},
			},
			LastUpdated: now,
		},
		{
			Domain:      DomainGaming,
			Description: "Video game sales, gaming industry analytics, platform performance",
			Keywords:    []string{"gaming", "video games", "game sales", "platform", "genre", "entertainment"},
			Metrics:     []string{"sales", "units_sold", "revenue", "market_share"},
			Dimensions:  []string{"platform", "genre", "year", "region"},
			Tables: []TableContext{
				{
					TableName:    "game_sales",
					Domain:       DomainGaming,
					Description:  "Video game sales data by platform and genre",
					BusinessTags: []string{"games", "sales", "entertainment"},
				},
			},
			Glossary: []BusinessGlossary{
				{
					Term:       "Platform",
					Definition: "Gaming system or console (e.g., PlayStation, Xbox, Nintendo)",
					Domain:     DomainGaming,
					Synonyms:   []string{"console", "system", "device"},
				},
			},
			Patterns: []QueryPattern{
				{
					Pattern:     "gaming analytics",
					Domain:      DomainGaming,
					Intent:      IntentAnalytics,
					Keywords:    []string{"game", "gaming", "video game", "sales"},
					Description: "Analysis of video game sales and performance",
					Examples:    []string{"video game sales", "gaming data", "game performance"},
				},
			},
			LastUpdated: now,
		},
		{
			Domain:      DomainSlack,
			Description: "Slack workspace analytics, channel activity, user engagement",
			Keywords:    []string{"slack", "channels", "messages", "users", "workspace", "communication"},
			Metrics:     []string{"messages", "active_users", "channel_activity", "response_time"},
			Dimensions:  []string{"channel", "user", "date", "team"},
			Tables: []TableContext{
				{
					TableName:    "slack_channels",
					Domain:       DomainSlack,
					Description:  "Slack channel activity and engagement metrics",
					BusinessTags: []string{"slack", "communication", "channels"},
				},
			},
			Patterns: []QueryPattern{
				{
					Pattern:     "slack analytics",
					Domain:      DomainSlack,
					Intent:      IntentAnalytics,
					Keywords:    []string{"slack", "channel", "messages", "communication"},
					Description: "Analysis of Slack workspace activity and engagement",
					Examples:    []string{"slack dashboard", "channel activity", "slack metrics"},
				},
			},
			LastUpdated: now,
		},
		{
			Domain:      DomainCOVID,
			Description: "COVID-19 vaccination data, pandemic statistics, public health metrics",
			Keywords:    []string{"covid", "vaccine", "vaccination", "pandemic", "health", "immunization"},
			Metrics:     []string{"vaccinated", "vaccination_rate", "doses_administered", "population_coverage"},
			Dimensions:  []string{"state", "country", "date", "vaccine_type"},
			Tables: []TableContext{
				{
					TableName:    "covid_vaccinations",
					Domain:       DomainCOVID,
					Description:  "COVID-19 vaccination data by location and time",
					BusinessTags: []string{"covid", "vaccination", "health"},
				},
			},
			Patterns: []QueryPattern{
				{
					Pattern:     "covid analytics",
					Domain:      DomainCOVID,
					Intent:      IntentAnalytics,
					Keywords:    []string{"covid", "vaccine", "vaccination", "pandemic"},
					Description: "Analysis of COVID-19 vaccination and health data",
					Examples:    []string{"covid vaccine dashboard", "vaccination data", "pandemic metrics"},
				},
			},
			LastUpdated: now,
		},
	}
}