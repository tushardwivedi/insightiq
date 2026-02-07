package events

import (
	"context"
	"fmt"
	"log/slog"

	"insightiq/backend/internal/domain"
	"insightiq/backend/internal/usecases"
)

// QueryHistoryEventHandler handles query execution events
// and delegates business logic to the use case layer
type QueryHistoryEventHandler struct {
	useCase *usecases.QueryHistoryUseCase
	logger  *slog.Logger
}

// NewQueryHistoryEventHandler creates a new query history event handler
func NewQueryHistoryEventHandler(useCase *usecases.QueryHistoryUseCase, logger *slog.Logger) *QueryHistoryEventHandler {
	return &QueryHistoryEventHandler{
		useCase: useCase,
		logger:  logger,
	}
}

// Handle processes a QueryExecutionEvent and saves it to the database
func (h *QueryHistoryEventHandler) Handle(ctx context.Context, event interface{}) error {
	// Type assertion to get the actual event
	queryEvent, ok := event.(*domain.QueryExecutionEvent)
	if !ok {
		return fmt.Errorf("invalid event type: expected *domain.QueryExecutionEvent, got %T", event)
	}

	h.logger.Info("Processing query execution event",
		"event_id", queryEvent.EventID,
		"user_id", queryEvent.UserID,
		"query_type", queryEvent.QueryType,
		"status", queryEvent.Status)

	// Convert domain event to use case request
	var resultPreview map[string]interface{}
	if queryEvent.ResultPreview != nil {
		if previewMap, ok := queryEvent.ResultPreview.(map[string]interface{}); ok {
			resultPreview = previewMap
		}
	}

	request := usecases.RecordQueryRequest{
		UserID:        queryEvent.UserID,
		ConnectorID:   queryEvent.ConnectorID,
		ConnectorName: queryEvent.ConnectorName,
		QueryType:     queryEvent.QueryType,
		QueryText:     queryEvent.QueryText,
		GeneratedSQL:  queryEvent.GeneratedSQL,
		RowCount:      queryEvent.RowCount,
		ExecutionTime: queryEvent.ExecutionTime,
		Status:        queryEvent.Status,
		ErrorMessage:  queryEvent.ErrorMessage,
		IPAddress:     queryEvent.IPAddress,
		UserAgent:     queryEvent.UserAgent,
		ResultPreview: resultPreview,
	}

	// Delegate to use case layer (handles validation, PII redaction, and persistence)
	queryHistory, err := h.useCase.RecordQuery(ctx, request)
	if err != nil {
		h.logger.Error("Failed to record query via use case",
			"event_id", queryEvent.EventID,
			"user_id", queryEvent.UserID,
			"error", err)
		return fmt.Errorf("failed to record query: %w", err)
	}

	h.logger.Info("Query history saved successfully",
		"event_id", queryEvent.EventID,
		"user_id", queryEvent.UserID,
		"query_history_id", queryHistory.ID)

	return nil
}
