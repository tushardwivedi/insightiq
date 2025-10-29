package connectors

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	_ "github.com/lib/pq"
)

var (
	// ErrDangerousQuery is returned when a query contains dangerous operations
	ErrDangerousQuery = errors.New("query contains dangerous or disallowed operations")
	// ErrInvalidQuery is returned when a query is malformed
	ErrInvalidQuery = errors.New("invalid or malformed query")
	// ErrQueryTooLong is returned when a query exceeds maximum length
	ErrQueryTooLong = errors.New("query exceeds maximum allowed length")
)

const (
	// Maximum query length to prevent DoS
	maxQueryLength = 10000
)

type PostgresConnector struct {
	db     *sql.DB
	logger *slog.Logger
}

type QueryResult struct {
	Data []map[string]interface{} `json:"data"`
}

func NewPostgresConnector(dbURL string, logger *slog.Logger) (*PostgresConnector, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresConnector{
		db:     db,
		logger: logger.With("connector", "postgres"),
	}, nil
}

func (pc *PostgresConnector) Close() error {
	return pc.db.Close()
}

// validateAndSanitizeQuery performs security checks on SQL queries
func (pc *PostgresConnector) validateAndSanitizeQuery(query string) error {
	// Check query length
	if len(query) > maxQueryLength {
		pc.logger.Warn("Query exceeds maximum length", "length", len(query))
		return ErrQueryTooLong
	}

	// Trim and normalize whitespace
	query = strings.TrimSpace(query)
	if query == "" {
		return ErrInvalidQuery
	}

	// Convert to uppercase for pattern matching
	queryUpper := strings.ToUpper(query)

	// Only allow SELECT statements (read-only)
	if !strings.HasPrefix(queryUpper, "SELECT") {
		pc.logger.Warn("Non-SELECT query attempted", "query_prefix", query[:min(50, len(query))])
		return ErrDangerousQuery
	}

	// Disallow dangerous keywords and patterns
	dangerousPatterns := []string{
		"DROP",
		"DELETE",
		"INSERT",
		"UPDATE",
		"ALTER",
		"CREATE",
		"TRUNCATE",
		"REPLACE",
		"EXEC",
		"EXECUTE",
		"GRANT",
		"REVOKE",
		"MERGE",
		"CALL",
		";DROP",        // SQL injection attempt
		"; DROP",       // SQL injection attempt
		"--",           // SQL comments can hide attacks
		"/*",           // Block comments
		"*/",
		"XP_",          // Extended stored procedures
		"SP_",          // System stored procedures
		"INFORMATION_SCHEMA", // Schema enumeration
		"PG_SLEEP",     // PostgreSQL sleep (DoS)
		"WAITFOR",      // SQL Server delay
		"BENCHMARK",    // MySQL DoS
		"INTO OUTFILE", // File system access
		"INTO DUMPFILE",
		"LOAD_FILE",
		"COPY ",        // PostgreSQL COPY command
		"\\COPY",       // PostgreSQL meta-command
	}

	for _, pattern := range dangerousPatterns {
		if strings.Contains(queryUpper, pattern) {
			pc.logger.Warn("Dangerous pattern detected in query",
				"pattern", pattern,
				"query_preview", query[:min(100, len(query))])
			return ErrDangerousQuery
		}
	}

	// Check for multiple statements (disallow semicolons except at the end)
	trimmedQuery := strings.TrimSuffix(query, ";")
	if strings.Contains(trimmedQuery, ";") {
		pc.logger.Warn("Multiple SQL statements detected")
		return ErrDangerousQuery
	}

	// Validate SQL syntax using regex - must look like a valid SELECT
	// This is a basic check, not a full SQL parser
	selectPattern := regexp.MustCompile(`(?i)^\s*SELECT\s+.+\s+FROM\s+[\w\.]+`)
	if !selectPattern.MatchString(query) {
		// Allow simple SELECT without FROM (e.g., SELECT 1)
		simpleSelectPattern := regexp.MustCompile(`(?i)^\s*SELECT\s+[\d\s\+\-\*\/]+(\s+as\s+\w+)?\s*$`)
		if !simpleSelectPattern.MatchString(query) {
			pc.logger.Warn("Query does not match expected SELECT pattern")
			return ErrInvalidQuery
		}
	}

	return nil
}

// ExecuteQuery executes a read-only SELECT query with security validation
// SECURITY: This function validates queries to prevent SQL injection
func (pc *PostgresConnector) ExecuteQuery(ctx context.Context, query string) (*QueryResult, error) {
	// Validate and sanitize the query before execution
	if err := pc.validateAndSanitizeQuery(query); err != nil {
		pc.logger.Error("Query validation failed",
			"error", err,
			"query_preview", query[:min(100, len(query))])
		return nil, fmt.Errorf("query validation failed: %w", err)
	}

	pc.logger.Info("Executing validated query", "query_length", len(query))

	rows, err := pc.db.QueryContext(ctx, query)
	if err != nil {
		pc.logger.Error("Query execution failed", "error", err)
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	var data []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		data = append(data, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	pc.logger.Info("Query executed successfully", "rows", len(data))
	return &QueryResult{Data: data}, nil
}

func (pc *PostgresConnector) GetBikeSalesData(ctx context.Context) (*QueryResult, error) {
	query := `
		SELECT
			quarter,
			bike_category,
			total_revenue,
			quantity as total_bikes_sold
		FROM bike_sales
		ORDER BY quarter, bike_category
	`
	return pc.ExecuteQuery(ctx, query)
}

func (pc *PostgresConnector) GetMonthlyActiveUsers(ctx context.Context) (*QueryResult, error) {
	query := `
		SELECT
			month_year,
			platform,
			active_users,
			region
		FROM monthly_active_users
		ORDER BY month_year, platform
	`
	return pc.ExecuteQuery(ctx, query)
}

func (pc *PostgresConnector) TestConnection(ctx context.Context) (*QueryResult, error) {
	// Simple test query to verify connection
	query := "SELECT 1 as test_connection"
	return pc.ExecuteQuery(ctx, query)
}

// ExecuteQueryWithParams executes a parameterized query (safer than raw SQL)
// SECURITY: Use this method when you need to include dynamic values in queries
// Example: ExecuteQueryWithParams(ctx, "SELECT * FROM users WHERE id = $1", userId)
func (pc *PostgresConnector) ExecuteQueryWithParams(ctx context.Context, query string, args ...interface{}) (*QueryResult, error) {
	// Still validate the query structure
	if err := pc.validateAndSanitizeQuery(query); err != nil {
		pc.logger.Error("Parameterized query validation failed",
			"error", err,
			"query_preview", query[:min(100, len(query))])
		return nil, fmt.Errorf("query validation failed: %w", err)
	}

	// Verify that the query uses placeholders if args are provided
	if len(args) > 0 {
		// Count placeholders ($1, $2, etc.)
		placeholderCount := 0
		for i := 1; i <= len(args); i++ {
			placeholder := fmt.Sprintf("$%d", i)
			if strings.Contains(query, placeholder) {
				placeholderCount++
			}
		}

		if placeholderCount != len(args) {
			pc.logger.Error("Placeholder count mismatch",
				"expected", len(args),
				"found", placeholderCount)
			return nil, fmt.Errorf("query has %d placeholders but %d arguments provided", placeholderCount, len(args))
		}
	}

	pc.logger.Info("Executing parameterized query",
		"query_length", len(query),
		"param_count", len(args))

	rows, err := pc.db.QueryContext(ctx, query, args...)
	if err != nil {
		pc.logger.Error("Parameterized query execution failed", "error", err)
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	var data []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		data = append(data, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	pc.logger.Info("Parameterized query executed successfully", "rows", len(data))
	return &QueryResult{Data: data}, nil
}