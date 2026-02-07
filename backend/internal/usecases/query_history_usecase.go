package usecases

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"insightiq/backend/internal/models"
	"insightiq/backend/internal/repository"
	"insightiq/backend/internal/security"
)

var (
	// ErrInvalidInput is returned when validation fails
	ErrInvalidInput = errors.New("invalid input")

	// ErrUnauthorized is returned when user is not authorized
	ErrUnauthorized = errors.New("unauthorized access")

	// ErrNotFound is returned when query history is not found
	ErrNotFound = errors.New("query history not found")
)

// QueryHistoryUseCase encapsulates business logic for query history management
// It sits between the event handler and the repository, providing:
// - Input validation
// - Business rule enforcement
// - PII redaction orchestration
// - Data access control
type QueryHistoryUseCase struct {
	repo     *repository.QueryHistoryRepository
	redactor security.PIIRedactor
	logger   *slog.Logger
}

// NewQueryHistoryUseCase creates a new use case instance
func NewQueryHistoryUseCase(
	repo *repository.QueryHistoryRepository,
	redactor security.PIIRedactor,
	logger *slog.Logger,
) *QueryHistoryUseCase {
	return &QueryHistoryUseCase{
		repo:     repo,
		redactor: redactor,
		logger:   logger.With("component", "query_history_usecase"),
	}
}

// RecordQueryRequest contains the data needed to record a query execution
type RecordQueryRequest struct {
	UserID        string
	ConnectorID   string
	ConnectorName string
	QueryType     string
	QueryText     string
	GeneratedSQL  string
	RowCount      int
	ExecutionTime int64
	Status        string
	ErrorMessage  string
	IPAddress     string
	UserAgent     string
	ResultPreview map[string]interface{}
}

// Validate checks if the request has all required fields
func (r *RecordQueryRequest) Validate() error {
	if r.UserID == "" {
		return fmt.Errorf("%w: user_id is required", ErrInvalidInput)
	}

	if r.QueryType == "" {
		return fmt.Errorf("%w: query_type is required", ErrInvalidInput)
	}

	if r.QueryText == "" {
		return fmt.Errorf("%w: query_text is required", ErrInvalidInput)
	}

	if r.Status == "" {
		return fmt.Errorf("%w: status is required", ErrInvalidInput)
	}

	// Validate status is one of the allowed values
	allowedStatuses := map[string]bool{
		"success": true,
		"error":   true,
		"timeout": true,
	}
	if !allowedStatuses[r.Status] {
		return fmt.Errorf("%w: status must be one of: success, error, timeout", ErrInvalidInput)
	}

	// If status is error, error message should be present
	if r.Status == "error" && r.ErrorMessage == "" {
		return fmt.Errorf("%w: error_message is required when status is error", ErrInvalidInput)
	}

	// Validate execution time is non-negative
	if r.ExecutionTime < 0 {
		return fmt.Errorf("%w: execution_time cannot be negative", ErrInvalidInput)
	}

	// Validate row count is non-negative
	if r.RowCount < 0 {
		return fmt.Errorf("%w: row_count cannot be negative", ErrInvalidInput)
	}

	return nil
}

// RecordQuery records a query execution with PII redaction
func (uc *QueryHistoryUseCase) RecordQuery(ctx context.Context, req RecordQueryRequest) (*models.QueryHistory, error) {
	// Validate input
	if err := req.Validate(); err != nil {
		uc.logger.Error("Validation failed for query recording",
			"user_id", req.UserID,
			"error", err)
		return nil, err
	}

	uc.logger.Info("Recording query execution",
		"user_id", req.UserID,
		"query_type", req.QueryType,
		"status", req.Status)

	// Perform PII redaction
	var dataRedacted bool
	redactedQueryText := req.QueryText
	redactedSQL := req.GeneratedSQL

	if uc.redactor != nil {
		// Redact query text
		originalQueryText := req.QueryText
		redactedQueryText = uc.redactor.RedactQuery(originalQueryText)

		// Redact generated SQL
		originalSQL := req.GeneratedSQL
		redactedSQL = uc.redactor.RedactSQL(originalSQL)

		// Set flag if any redaction occurred
		dataRedacted = (originalQueryText != redactedQueryText) || (originalSQL != redactedSQL)

		if dataRedacted {
			uc.logger.Info("PII redaction applied",
				"user_id", req.UserID,
				"query_text_changed", originalQueryText != redactedQueryText,
				"sql_changed", originalSQL != redactedSQL)
		}
	} else {
		uc.logger.Warn("PII redactor not configured, storing unredacted data",
			"user_id", req.UserID)
	}

	// Create query history model
	queryHistory := &models.QueryHistory{
		UserID:        req.UserID,
		ConnectorID:   req.ConnectorID,
		ConnectorName: req.ConnectorName,
		QueryType:     req.QueryType,
		QueryText:     redactedQueryText,
		GeneratedSQL:  redactedSQL,
		RowCount:      req.RowCount,
		ExecutionTime: req.ExecutionTime,
		Status:        req.Status,
		ErrorMessage:  req.ErrorMessage,
		IPAddress:     req.IPAddress,
		UserAgent:     req.UserAgent,
		DataRedacted:  dataRedacted,
		ResultPreview: req.ResultPreview,
	}

	// Save to repository
	if err := uc.repo.Create(ctx, queryHistory); err != nil {
		uc.logger.Error("Failed to save query history",
			"user_id", req.UserID,
			"error", err)
		return nil, fmt.Errorf("failed to save query history: %w", err)
	}

	uc.logger.Info("Query history saved successfully",
		"user_id", req.UserID,
		"query_history_id", queryHistory.ID,
		"data_redacted", dataRedacted)

	return queryHistory, nil
}

// GetUserHistory retrieves paginated query history for a user
func (uc *QueryHistoryUseCase) GetUserHistory(ctx context.Context, userID string, limit, offset int) ([]models.QueryHistoryListItem, error) {
	// Validate input
	if userID == "" {
		return nil, fmt.Errorf("%w: user_id is required", ErrInvalidInput)
	}

	// Validate pagination parameters
	if limit < 1 {
		limit = 10 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit
	}
	if offset < 0 {
		offset = 0
	}

	uc.logger.Debug("Retrieving query history",
		"user_id", userID,
		"limit", limit,
		"offset", offset)

	items, err := uc.repo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		uc.logger.Error("Failed to retrieve query history",
			"user_id", userID,
			"error", err)
		return nil, fmt.Errorf("failed to retrieve query history: %w", err)
	}

	uc.logger.Debug("Query history retrieved",
		"user_id", userID,
		"count", len(items))

	return items, nil
}

// GetQueryDetails retrieves full details of a specific query
func (uc *QueryHistoryUseCase) GetQueryDetails(ctx context.Context, queryID, userID string) (*models.QueryHistory, error) {
	// Validate input
	if queryID == "" {
		return nil, fmt.Errorf("%w: query_id is required", ErrInvalidInput)
	}
	if userID == "" {
		return nil, fmt.Errorf("%w: user_id is required", ErrInvalidInput)
	}

	uc.logger.Debug("Retrieving query details",
		"query_id", queryID,
		"user_id", userID)

	query, err := uc.repo.GetByID(ctx, queryID, userID)
	if err != nil {
		uc.logger.Error("Failed to retrieve query details",
			"query_id", queryID,
			"user_id", userID,
			"error", err)
		return nil, fmt.Errorf("%w: query not found", ErrNotFound)
	}

	// Verify the query belongs to the requesting user (defense in depth)
	if query.UserID != userID {
		uc.logger.Warn("Unauthorized access attempt",
			"query_id", queryID,
			"requesting_user", userID,
			"query_owner", query.UserID)
		return nil, ErrUnauthorized
	}

	uc.logger.Debug("Query details retrieved",
		"query_id", queryID,
		"user_id", userID)

	return query, nil
}

// GetUserStats retrieves aggregated statistics for a user's query history
func (uc *QueryHistoryUseCase) GetUserStats(ctx context.Context, userID string) (map[string]interface{}, error) {
	// Validate input
	if userID == "" {
		return nil, fmt.Errorf("%w: user_id is required", ErrInvalidInput)
	}

	uc.logger.Debug("Retrieving query statistics",
		"user_id", userID)

	stats, err := uc.repo.GetStats(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to retrieve query statistics",
			"user_id", userID,
			"error", err)
		return nil, fmt.Errorf("failed to retrieve statistics: %w", err)
	}

	uc.logger.Debug("Query statistics retrieved",
		"user_id", userID,
		"total_queries", stats["total_queries"])

	return stats, nil
}

// DeleteQuery removes a specific query from history
func (uc *QueryHistoryUseCase) DeleteQuery(ctx context.Context, queryID, userID string) error {
	// Validate input
	if queryID == "" {
		return fmt.Errorf("%w: query_id is required", ErrInvalidInput)
	}
	if userID == "" {
		return fmt.Errorf("%w: user_id is required", ErrInvalidInput)
	}

	uc.logger.Info("Deleting query history",
		"query_id", queryID,
		"user_id", userID)

	// Verify ownership before deletion (using GetByID which checks user_id)
	_, err := uc.repo.GetByID(ctx, queryID, userID)
	if err != nil {
		return fmt.Errorf("%w: query not found or unauthorized", ErrUnauthorized)
	}

	// Delete the query
	if err := uc.repo.Delete(ctx, queryID, userID); err != nil {
		uc.logger.Error("Failed to delete query history",
			"query_id", queryID,
			"user_id", userID,
			"error", err)
		return fmt.Errorf("failed to delete query: %w", err)
	}

	uc.logger.Info("Query history deleted successfully",
		"query_id", queryID,
		"user_id", userID)

	return nil
}

// DeleteOldQueries removes query history entries older than the specified duration
// This is used by the data retention policy
func (uc *QueryHistoryUseCase) DeleteOldQueries(ctx context.Context, retentionDuration time.Duration) (int64, error) {
	// Validate input
	if retentionDuration < 24*time.Hour {
		return 0, fmt.Errorf("%w: retention duration must be at least 24 hours", ErrInvalidInput)
	}

	uc.logger.Info("Deleting old query history",
		"retention_duration", retentionDuration.String())

	rowsDeleted, err := uc.repo.DeleteOlderThan(ctx, retentionDuration)
	if err != nil {
		uc.logger.Error("Failed to delete old queries",
			"retention_duration", retentionDuration.String(),
			"error", err)
		return 0, fmt.Errorf("failed to delete old queries: %w", err)
	}

	uc.logger.Info("Old query history deleted",
		"retention_duration", retentionDuration.String(),
		"rows_deleted", rowsDeleted)

	return rowsDeleted, nil
}

// DeleteUserHistory removes all query history for a specific user
// This is used for GDPR "right to erasure" compliance
func (uc *QueryHistoryUseCase) DeleteUserHistory(ctx context.Context, userID string) (int64, error) {
	// Validate input
	if userID == "" {
		return 0, fmt.Errorf("%w: user_id is required", ErrInvalidInput)
	}

	uc.logger.Info("Deleting all query history for user (GDPR erasure)",
		"user_id", userID)

	// Get all user's queries first for counting
	items, err := uc.repo.GetByUserID(ctx, userID, 10000, 0) // Large limit to get all
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve user queries: %w", err)
	}

	// Delete each query
	var deletedCount int64
	for _, item := range items {
		if err := uc.repo.Delete(ctx, item.ID, userID); err != nil {
			uc.logger.Error("Failed to delete query during user history deletion",
				"query_id", item.ID,
				"user_id", userID,
				"error", err)
			continue // Continue deleting other queries
		}
		deletedCount++
	}

	uc.logger.Info("User query history deleted (GDPR compliance)",
		"user_id", userID,
		"queries_deleted", deletedCount)

	return deletedCount, nil
}
