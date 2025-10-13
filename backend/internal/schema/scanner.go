package schema

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

// ScannerService handles database schema scanning and analysis
type ScannerService struct {
	connectorService ConnectorService
	logger          *slog.Logger
}

// NewScannerService creates a new schema scanner service
func NewScannerService(connectorService ConnectorService, logger *slog.Logger) *ScannerService {
	return &ScannerService{
		connectorService: connectorService,
		logger:          logger,
	}
}

// ScanDataSource performs comprehensive schema analysis for a data source
func (s *ScannerService) ScanDataSource(ctx context.Context, connectorID string) (*SchemaContext, error) {
	s.logger.Info("Starting schema scan", "connector_id", connectorID)

	connector, err := s.connectorService.GetConnector(ctx, connectorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get connector: %w", err)
	}

	// Scan tables based on connector type
	var tables []TableContext
	switch connector.Type {
	case "postgres", "mysql":
		tables, err = s.scanDatabaseTables(ctx, connector)
	case "superset":
		tables, err = s.scanSupersetTables(ctx, connector)
	case "api":
		tables, err = s.scanAPIEndpoints(ctx, connector)
	default:
		return nil, fmt.Errorf("unsupported connector type: %s", connector.Type)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to scan tables: %w", err)
	}

	// Analyze relationships between tables (not used in current implementation)
	// relationships := s.analyzeTableRelationships(tables)

	// Generate business metrics from schema
	businessMetrics := s.extractBusinessMetrics(tables)

	// Generate sample queries
	sampleQueries := s.generateSampleQueries(tables, businessMetrics)

	// Create schema context
	schemaContext := &SchemaContext{
		ConnectorID:     connectorID,
		ConnectorName:   connector.Name,
		ConnectorType:   string(connector.Type),
		Tables:          tables,
		BusinessMetrics: businessMetrics,
		SampleQueries:   sampleQueries,
		AnalyzedAt:      time.Now(),
	}

	s.logger.Info("Schema scan completed",
		"connector_id", connectorID,
		"tables", len(tables),
		"metrics", len(businessMetrics))

	return schemaContext, nil
}

// scanDatabaseTables scans PostgreSQL/MySQL database tables
func (s *ScannerService) scanDatabaseTables(ctx context.Context, connector *ConnectorInfo) ([]TableContext, error) {
	// For now, return mock data representing a typical business database
	// In production, this would connect to the actual database and scan schema
	tables := []TableContext{
		{
			TableName:   "orders",
			Schema:      "public",
			Description: "Customer order transactions",
			Columns: []ColumnInfo{
				{
					Name:         "order_id",
					Type:         "INTEGER",
					DataType:     "int",
					IsID:         true,
					Unique:       true,
					SampleValues: []string{"1001", "1002", "1003"},
				},
				{
					Name:         "customer_id",
					Type:         "INTEGER",
					DataType:     "int",
					IsID:         true,
					SampleValues: []string{"501", "502", "503"},
				},
				{
					Name:         "order_date",
					Type:         "TIMESTAMP",
					DataType:     "timestamp",
					IsDatetime:   true,
					SampleValues: []string{"2024-01-15", "2024-01-16", "2024-01-17"},
				},
				{
					Name:         "total_amount",
					Type:         "DECIMAL",
					DataType:     "decimal",
					IsMetric:     true,
					IsCurrency:   true,
					SampleValues: []string{"299.99", "149.50", "899.00"},
				},
				{
					Name:         "status",
					Type:         "VARCHAR",
					DataType:     "varchar",
					IsDimension:  true,
					SampleValues: []string{"completed", "pending", "cancelled"},
				},
			},
			BusinessTags: []string{"sales", "orders", "revenue", "transactions"},
		},
		{
			TableName:   "customers",
			Schema:      "public",
			Description: "Customer information and demographics",
			Columns: []ColumnInfo{
				{
					Name:         "customer_id",
					Type:         "INTEGER",
					DataType:     "int",
					IsID:         true,
					Unique:       true,
					SampleValues: []string{"501", "502", "503"},
				},
				{
					Name:         "email",
					Type:         "VARCHAR",
					DataType:     "varchar",
					SampleValues: []string{"john@example.com", "jane@company.com"},
				},
				{
					Name:         "registration_date",
					Type:         "TIMESTAMP",
					DataType:     "timestamp",
					IsDatetime:   true,
					SampleValues: []string{"2023-12-01", "2023-11-15", "2024-01-03"},
				},
				{
					Name:         "customer_segment",
					Type:         "VARCHAR",
					DataType:     "varchar",
					IsDimension:  true,
					SampleValues: []string{"premium", "standard", "basic"},
				},
				{
					Name:         "lifetime_value",
					Type:         "DECIMAL",
					DataType:     "decimal",
					IsMetric:     true,
					IsCurrency:   true,
					SampleValues: []string{"1299.99", "599.50", "299.00"},
				},
			},
			BusinessTags: []string{"customers", "crm", "segmentation", "ltv"},
		},
		{
			TableName:   "products",
			Schema:      "public",
			Description: "Product catalog and inventory",
			Columns: []ColumnInfo{
				{
					Name:         "product_id",
					Type:         "INTEGER",
					DataType:     "int",
					IsID:         true,
					Unique:       true,
					SampleValues: []string{"101", "102", "103"},
				},
				{
					Name:         "product_name",
					Type:         "VARCHAR",
					DataType:     "varchar",
					SampleValues: []string{"Laptop Pro", "Wireless Mouse", "Monitor 4K"},
				},
				{
					Name:         "category",
					Type:         "VARCHAR",
					DataType:     "varchar",
					IsDimension:  true,
					SampleValues: []string{"electronics", "accessories", "computers"},
				},
				{
					Name:         "price",
					Type:         "DECIMAL",
					DataType:     "decimal",
					IsMetric:     true,
					IsCurrency:   true,
					SampleValues: []string{"1299.99", "29.99", "599.00"},
				},
				{
					Name:         "stock_quantity",
					Type:         "INTEGER",
					DataType:     "int",
					IsMetric:     true,
					SampleValues: []string{"50", "200", "25"},
				},
			},
			BusinessTags: []string{"products", "inventory", "catalog", "pricing"},
		},
	}

	return tables, nil
}

// scanSupersetTables scans Superset dashboards and datasets
func (s *ScannerService) scanSupersetTables(ctx context.Context, connector *ConnectorInfo) ([]TableContext, error) {
	// This would connect to Superset API and extract available datasets
	// For now, return representative data from typical Superset setup
	return s.scanDatabaseTables(ctx, connector) // Fallback to database scanning
}

// scanAPIEndpoints scans API endpoints to understand data structure
func (s *ScannerService) scanAPIEndpoints(ctx context.Context, connector *ConnectorInfo) ([]TableContext, error) {
	// This would call API endpoints and analyze response structures
	// For now, return basic API table representation
	tables := []TableContext{
		{
			TableName:   "api_response",
			Description: "API endpoint response data",
			Columns: []ColumnInfo{
				{
					Name:         "id",
					Type:         "STRING",
					DataType:     "string",
					IsID:         true,
					SampleValues: []string{"api_001", "api_002"},
				},
				{
					Name:         "timestamp",
					Type:         "TIMESTAMP",
					DataType:     "timestamp",
					IsDatetime:   true,
					SampleValues: []string{"2024-01-15T10:30:00Z"},
				},
			},
			BusinessTags: []string{"api", "external_data"},
		},
	}

	return tables, nil
}

// analyzeTableRelationships identifies relationships between tables
func (s *ScannerService) analyzeTableRelationships(tables []TableContext) []TableRelationship {
	var relationships []TableRelationship

	// Simple relationship detection based on column names
	for i, table1 := range tables {
		for j, table2 := range tables {
			if i >= j {
				continue // Avoid self-comparison and duplicates
			}

			// Look for foreign key relationships
			for _, col1 := range table1.Columns {
				for _, col2 := range table2.Columns {
					if s.isPotentialRelationship(col1, col2, table1.TableName, table2.TableName) {
						relationships = append(relationships, TableRelationship{
							FromTable:    table1.TableName,
							ToTable:      table2.TableName,
							FromColumn:   col1.Name,
							ToColumn:     col2.Name,
							RelationType: "one_to_many",
							Confidence:   0.8,
						})
					}
				}
			}
		}
	}

	return relationships
}

// isPotentialRelationship checks if two columns represent a foreign key relationship
func (s *ScannerService) isPotentialRelationship(col1, col2 ColumnInfo, table1, table2 string) bool {
	// Check for matching ID columns
	if col1.IsID && col2.IsID {
		// customer_id in orders table -> customer_id in customers table
		if col1.Name == col2.Name {
			return true
		}
		// order_id matches id in orders table
		if strings.Contains(col1.Name, strings.TrimSuffix(table2, "s")) ||
		   strings.Contains(col2.Name, strings.TrimSuffix(table1, "s")) {
			return true
		}
	}
	return false
}

// extractBusinessMetrics identifies potential business metrics from table schemas
func (s *ScannerService) extractBusinessMetrics(tables []TableContext) []BusinessMetric {
	var metrics []BusinessMetric

	// Metric patterns for classification (currently using simpler approach)
	// metricPatterns := map[string][]string{
	//	"revenue":     {"total", "amount", "revenue", "sales", "price"},
	//	"count":       {"quantity", "count", "num", "total"},
	//	"rate":        {"rate", "percentage", "ratio"},
	//	"conversion":  {"conversion", "convert"},
	//	"engagement":  {"views", "clicks", "engagement", "interactions"},
	// }

	for _, table := range tables {
		domain := s.inferDomainFromTable(table)

		for _, col := range table.Columns {
			if col.IsMetric {
				metricType := s.inferMetricType(col.Name)

				// Find related dimension columns in same table
				var dimensions []string
				for _, dimCol := range table.Columns {
					if dimCol.IsDimension {
						dimensions = append(dimensions, dimCol.Name)
					}
				}

				metric := BusinessMetric{
					Name:        s.generateMetricName(col.Name, table.TableName),
					Description: s.generateMetricDescription(col.Name, table.TableName),
					Type:        metricType,
					Table:       table.TableName,
					Column:      col.Name,
					Dimensions:  dimensions,
					Domain:      domain,
					Keywords:    s.generateMetricKeywords(col.Name, table.TableName),
				}

				metrics = append(metrics, metric)
			}
		}
	}

	return metrics
}

// inferDomainFromTable infers business domain from table structure
func (s *ScannerService) inferDomainFromTable(table TableContext) Domain {
	tableName := strings.ToLower(table.TableName)

	// Check table name patterns
	switch {
	case strings.Contains(tableName, "order") || strings.Contains(tableName, "sale") || strings.Contains(tableName, "transaction"):
		return DomainSales
	case strings.Contains(tableName, "campaign") || strings.Contains(tableName, "lead") || strings.Contains(tableName, "marketing"):
		return DomainMarketing
	case strings.Contains(tableName, "customer") || strings.Contains(tableName, "user") || strings.Contains(tableName, "client"):
		return DomainCustomer
	case strings.Contains(tableName, "product") || strings.Contains(tableName, "inventory") || strings.Contains(tableName, "catalog"):
		return DomainProduct
	case strings.Contains(tableName, "employee") || strings.Contains(tableName, "hr") || strings.Contains(tableName, "staff"):
		return DomainHR
	case strings.Contains(tableName, "finance") || strings.Contains(tableName, "payment") || strings.Contains(tableName, "invoice"):
		return DomainFinance
	default:
		return DomainGeneral
	}
}

// inferMetricType determines the aggregation type for a metric
func (s *ScannerService) inferMetricType(columnName string) string {
	name := strings.ToLower(columnName)

	switch {
	case strings.Contains(name, "amount") || strings.Contains(name, "total") || strings.Contains(name, "revenue") || strings.Contains(name, "price"):
		return "sum"
	case strings.Contains(name, "count") || strings.Contains(name, "quantity") || strings.Contains(name, "num"):
		return "count"
	case strings.Contains(name, "rate") || strings.Contains(name, "percentage"):
		return "avg"
	case strings.Contains(name, "ratio"):
		return "ratio"
	default:
		return "sum"
	}
}

// generateMetricName creates a business-friendly metric name
func (s *ScannerService) generateMetricName(columnName, tableName string) string {
	// Convert snake_case to Title Case
	name := strings.ReplaceAll(columnName, "_", " ")
	name = strings.Title(name)

	// Add context from table name if helpful
	if !strings.Contains(strings.ToLower(name), strings.ToLower(tableName)) {
		context := strings.Title(strings.TrimSuffix(tableName, "s"))
		name = context + " " + name
	}

	return name
}

// generateMetricDescription creates a description for the metric
func (s *ScannerService) generateMetricDescription(columnName, tableName string) string {
	name := strings.ToLower(columnName)
	table := strings.ToLower(tableName)

	switch {
	case strings.Contains(name, "total") || strings.Contains(name, "amount"):
		return fmt.Sprintf("Total %s from %s", strings.ReplaceAll(columnName, "_", " "), table)
	case strings.Contains(name, "count"):
		return fmt.Sprintf("Number of records in %s", table)
	case strings.Contains(name, "revenue"):
		return fmt.Sprintf("Revenue generated from %s", table)
	case strings.Contains(name, "price"):
		return fmt.Sprintf("Price information from %s", table)
	default:
		return fmt.Sprintf("%s metric from %s table", strings.Title(strings.ReplaceAll(columnName, "_", " ")), table)
	}
}

// generateMetricKeywords creates search keywords for the metric
func (s *ScannerService) generateMetricKeywords(columnName, tableName string) []string {
	keywords := []string{
		strings.ToLower(columnName),
		strings.ToLower(strings.ReplaceAll(columnName, "_", " ")),
		strings.ToLower(tableName),
		strings.ToLower(strings.TrimSuffix(tableName, "s")),
	}

	// Add domain-specific keywords
	name := strings.ToLower(columnName)
	switch {
	case strings.Contains(name, "amount") || strings.Contains(name, "revenue"):
		keywords = append(keywords, "money", "sales", "income", "revenue")
	case strings.Contains(name, "count") || strings.Contains(name, "quantity"):
		keywords = append(keywords, "number", "total", "count")
	case strings.Contains(name, "rate"):
		keywords = append(keywords, "percentage", "ratio", "rate")
	}

	return keywords
}

// generateSampleQueries creates sample queries based on discovered schema
func (s *ScannerService) generateSampleQueries(tables []TableContext, metrics []BusinessMetric) []string {
	var queries []string

	// Generate queries based on detected metrics
	for _, metric := range metrics {
		switch metric.Domain {
		case DomainSales:
			queries = append(queries,
				"Show me total sales revenue",
				"What are our top selling products?",
				"How many orders did we have this month?",
			)
		case DomainMarketing:
			queries = append(queries,
				"Show me campaign performance",
				"What's our lead conversion rate?",
				"Which marketing channels are most effective?",
			)
		case DomainCustomer:
			queries = append(queries,
				"Show me customer segmentation",
				"What's our customer lifetime value?",
				"How many new customers this quarter?",
			)
		}
	}

	// Remove duplicates
	seen := make(map[string]bool)
	var uniqueQueries []string
	for _, query := range queries {
		if !seen[query] {
			seen[query] = true
			uniqueQueries = append(uniqueQueries, query)
		}
	}

	return uniqueQueries
}