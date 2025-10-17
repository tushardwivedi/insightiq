package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"insightiq/backend/internal/models"
)

type QueryHistoryRepository struct {
	db *sqlx.DB
}

func NewQueryHistoryRepository(db *sqlx.DB) *QueryHistoryRepository {
	return &QueryHistoryRepository{db: db}
}

// CreateTables creates the query_history table if it doesn't exist
func (r *QueryHistoryRepository) CreateTables(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS query_history (
			id VARCHAR(255) PRIMARY KEY,
			user_id VARCHAR(255) NOT NULL,
			connector_id VARCHAR(255),
			connector_name VARCHAR(255),
			query_type VARCHAR(50) NOT NULL,
			query_text TEXT NOT NULL,
			generated_sql TEXT,
			result_preview JSONB,
			row_count INTEGER DEFAULT 0,
			execution_time_ms BIGINT DEFAULT 0,
			status VARCHAR(50) NOT NULL DEFAULT 'success',
			error_message TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS idx_query_history_user_id ON query_history(user_id);
		CREATE INDEX IF NOT EXISTS idx_query_history_created_at ON query_history(created_at DESC);
		CREATE INDEX IF NOT EXISTS idx_query_history_status ON query_history(status);
		CREATE INDEX IF NOT EXISTS idx_query_history_query_type ON query_history(query_type);
	`

	_, err := r.db.ExecContext(ctx, query)
	return err
}

// Create saves a new query history entry
func (r *QueryHistoryRepository) Create(ctx context.Context, qh *models.QueryHistory) error {
	qh.ID = uuid.New().String()
	qh.CreatedAt = time.Now()

	// Convert result_preview to JSON
	var resultPreviewJSON []byte
	var err error
	if qh.ResultPreview != nil {
		resultPreviewJSON, err = json.Marshal(qh.ResultPreview)
		if err != nil {
			return err
		}
	}

	query := `
		INSERT INTO query_history (
			id, user_id, connector_id, connector_name, query_type, query_text,
			generated_sql, result_preview, row_count, execution_time_ms, status, error_message, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err = r.db.ExecContext(ctx, query,
		qh.ID, qh.UserID, qh.ConnectorID, qh.ConnectorName, qh.QueryType, qh.QueryText,
		qh.GeneratedSQL, resultPreviewJSON, qh.RowCount, qh.ExecutionTime, qh.Status, qh.ErrorMessage, qh.CreatedAt,
	)

	return err
}

// GetByUserID retrieves query history for a specific user with pagination
func (r *QueryHistoryRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]models.QueryHistoryListItem, error) {
	query := `
		SELECT id, query_type, query_text, connector_name, row_count, status, created_at
		FROM query_history
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	var items []models.QueryHistoryListItem
	err := r.db.SelectContext(ctx, &items, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// GetByID retrieves a single query history entry with full details
func (r *QueryHistoryRepository) GetByID(ctx context.Context, id string, userID string) (*models.QueryHistory, error) {
	query := `
		SELECT id, user_id, connector_id, connector_name, query_type, query_text,
		       generated_sql, result_preview, row_count, execution_time_ms, status, error_message, created_at
		FROM query_history
		WHERE id = $1 AND user_id = $2
	`

	var qh models.QueryHistory
	var resultPreviewJSON []byte

	err := r.db.QueryRowContext(ctx, query, id, userID).Scan(
		&qh.ID, &qh.UserID, &qh.ConnectorID, &qh.ConnectorName, &qh.QueryType, &qh.QueryText,
		&qh.GeneratedSQL, &resultPreviewJSON, &qh.RowCount, &qh.ExecutionTime, &qh.Status, &qh.ErrorMessage, &qh.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Parse result_preview JSON
	if resultPreviewJSON != nil {
		err = json.Unmarshal(resultPreviewJSON, &qh.ResultPreview)
		if err != nil {
			return nil, err
		}
	}

	return &qh, nil
}

// Delete removes a query history entry
func (r *QueryHistoryRepository) Delete(ctx context.Context, id string, userID string) error {
	query := `DELETE FROM query_history WHERE id = $1 AND user_id = $2`
	_, err := r.db.ExecContext(ctx, query, id, userID)
	return err
}

// DeleteOlderThan removes query history entries older than a specified duration
func (r *QueryHistoryRepository) DeleteOlderThan(ctx context.Context, duration time.Duration) (int64, error) {
	query := `DELETE FROM query_history WHERE created_at < $1`
	cutoff := time.Now().Add(-duration)

	result, err := r.db.ExecContext(ctx, query, cutoff)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	return rowsAffected, err
}

// GetStats returns statistics about query history for a user
func (r *QueryHistoryRepository) GetStats(ctx context.Context, userID string) (map[string]interface{}, error) {
	query := `
		SELECT
			COUNT(*) as total_queries,
			COUNT(CASE WHEN status = 'success' THEN 1 END) as successful_queries,
			COUNT(CASE WHEN status = 'error' THEN 1 END) as failed_queries,
			AVG(execution_time_ms) as avg_execution_time_ms,
			SUM(row_count) as total_rows_returned
		FROM query_history
		WHERE user_id = $1
	`

	var stats struct {
		TotalQueries       int     `db:"total_queries"`
		SuccessfulQueries  int     `db:"successful_queries"`
		FailedQueries      int     `db:"failed_queries"`
		AvgExecutionTimeMS float64 `db:"avg_execution_time_ms"`
		TotalRowsReturned  int64   `db:"total_rows_returned"`
	}

	err := r.db.GetContext(ctx, &stats, query, userID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_queries":         stats.TotalQueries,
		"successful_queries":    stats.SuccessfulQueries,
		"failed_queries":        stats.FailedQueries,
		"avg_execution_time_ms": stats.AvgExecutionTimeMS,
		"total_rows_returned":   stats.TotalRowsReturned,
	}, nil
}
