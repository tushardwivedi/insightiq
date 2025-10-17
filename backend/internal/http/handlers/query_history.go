package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"insightiq/backend/internal/repository"
)

type QueryHistoryHandler struct {
	repo   *repository.QueryHistoryRepository
	logger *slog.Logger
}

func NewQueryHistoryHandler(repo *repository.QueryHistoryRepository, logger *slog.Logger) *QueryHistoryHandler {
	return &QueryHistoryHandler{
		repo:   repo,
		logger: logger.With("handler", "query_history"),
	}
}

// ListQueryHistory returns paginated query history for the authenticated user
func (h *QueryHistoryHandler) ListQueryHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user ID from context (set by auth middleware)
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get pagination parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	offset := 0 // default
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Get query history
	history, err := h.repo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get query history", "error", err, "user_id", userID)
		http.Error(w, "Failed to retrieve query history", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":   history,
		"limit":  limit,
		"offset": offset,
		"count":  len(history),
	})
}

// GetQueryHistoryByID returns a single query history entry with full details
func (h *QueryHistoryHandler) GetQueryHistoryByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user ID from context
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get query ID from URL parameter
	queryID := r.URL.Query().Get("id")
	if queryID == "" {
		http.Error(w, "Missing query ID", http.StatusBadRequest)
		return
	}

	// Get query history entry
	entry, err := h.repo.GetByID(ctx, queryID, userID)
	if err != nil {
		h.logger.Error("Failed to get query history entry", "error", err, "query_id", queryID, "user_id", userID)
		http.Error(w, "Query not found", http.StatusNotFound)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": entry,
	})
}

// DeleteQueryHistory deletes a query history entry
func (h *QueryHistoryHandler) DeleteQueryHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user ID from context
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get query ID from URL parameter
	queryID := r.URL.Query().Get("id")
	if queryID == "" {
		http.Error(w, "Missing query ID", http.StatusBadRequest)
		return
	}

	// Delete query history entry
	err := h.repo.Delete(ctx, queryID, userID)
	if err != nil {
		h.logger.Error("Failed to delete query history entry", "error", err, "query_id", queryID, "user_id", userID)
		http.Error(w, "Failed to delete query", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Query history deleted successfully",
	})
}

// GetQueryHistoryStats returns statistics about the user's query history
func (h *QueryHistoryHandler) GetQueryHistoryStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user ID from context
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get stats
	stats, err := h.repo.GetStats(ctx, userID)
	if err != nil {
		h.logger.Error("Failed to get query history stats", "error", err, "user_id", userID)
		http.Error(w, "Failed to retrieve stats", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": stats,
	})
}
