package connectors

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
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

func (pc *PostgresConnector) ExecuteQuery(ctx context.Context, query string) (*QueryResult, error) {
	pc.logger.Info("Executing query", "sql", query)

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