package models

import (
	"time"
)

// QueryHistory represents a stored query execution
type QueryHistory struct {
	ID            string                 `json:"id" db:"id"`
	UserID        string                 `json:"user_id" db:"user_id"`
	ConnectorID   string                 `json:"connector_id" db:"connector_id"`
	ConnectorName string                 `json:"connector_name" db:"connector_name"`
	QueryType     string                 `json:"query_type" db:"query_type"` // "text", "sql", "voice"
	QueryText     string                 `json:"query_text" db:"query_text"`
	GeneratedSQL  string                 `json:"generated_sql,omitempty" db:"generated_sql"`
	ResultPreview map[string]interface{} `json:"result_preview,omitempty" db:"result_preview"`
	RowCount      int                    `json:"row_count" db:"row_count"`
	ExecutionTime int64                  `json:"execution_time_ms" db:"execution_time_ms"` // milliseconds
	Status        string                 `json:"status" db:"status"`                       // "success", "error", "partial"
	ErrorMessage  string                 `json:"error_message,omitempty" db:"error_message"`
	CreatedAt     time.Time              `json:"created_at" db:"created_at"`
}

// QueryHistoryListItem is a simplified version for list views
type QueryHistoryListItem struct {
	ID            string    `json:"id" db:"id"`
	QueryType     string    `json:"query_type" db:"query_type"`
	QueryText     string    `json:"query_text" db:"query_text"`
	ConnectorName string    `json:"connector_name" db:"connector_name"`
	RowCount      int       `json:"row_count" db:"row_count"`
	Status        string    `json:"status" db:"status"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}
