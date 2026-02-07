package domain

import "time"

// QueryExecutionEvent represents a query execution event in the domain
// This event is published when a query is successfully executed
type QueryExecutionEvent struct {
	// EventID is a unique identifier for this event
	EventID string `json:"event_id"`

	// UserID is the ID of the user who executed the query
	UserID string `json:"user_id"`

	// Timestamp is when the query was executed
	Timestamp time.Time `json:"timestamp"`

	// QueryType indicates the type of query (text/sql/voice)
	QueryType string `json:"query_type"`

	// QueryText is the original query text from the user
	QueryText string `json:"query_text"`

	// GeneratedSQL is the SQL query generated from the user's query
	GeneratedSQL string `json:"generated_sql"`

	// ConnectorID is the ID of the connector used
	ConnectorID string `json:"connector_id"`

	// ConnectorName is the name of the connector used
	ConnectorName string `json:"connector_name"`

	// RowCount is the number of rows returned by the query
	RowCount int `json:"row_count"`

	// ExecutionTime is the query execution time in milliseconds
	ExecutionTime int64 `json:"execution_time"`

	// Status indicates the execution status (success/error/partial)
	Status string `json:"status"`

	// ErrorMessage contains error details if status is error
	ErrorMessage string `json:"error_message,omitempty"`

	// IPAddress is the IP address of the client (for audit trail)
	IPAddress string `json:"ip_address"`

	// UserAgent is the user agent string of the client (for audit trail)
	UserAgent string `json:"user_agent"`

	// ResultPreview contains a preview of the query results
	ResultPreview interface{} `json:"result_preview,omitempty"`
}

// Event types
const (
	EventTypeQueryExecuted = "query.executed"
)
