package http

import (
	"net/http"

	"insightiq/backend/internal/http/handlers"
	"insightiq/backend/internal/repository"
)

// handleQueryHistory routes query history requests based on HTTP method
func (s *Server) handleQueryHistory(w http.ResponseWriter, r *http.Request) {
	repo, ok := s.queryHistoryRepo.(*repository.QueryHistoryRepository)
	if !ok {
		http.Error(w, "Query history not available", http.StatusServiceUnavailable)
		return
	}

	handler := handlers.NewQueryHistoryHandler(repo, s.logger)

	switch r.Method {
	case http.MethodGet:
		if r.URL.Query().Get("id") != "" {
			handler.GetQueryHistoryByID(w, r)
		} else {
			handler.ListQueryHistory(w, r)
		}
	case http.MethodDelete:
		handler.DeleteQueryHistory(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleQueryHistoryStats returns query history statistics
func (s *Server) handleQueryHistoryStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	repo, ok := s.queryHistoryRepo.(*repository.QueryHistoryRepository)
	if !ok {
		http.Error(w, "Query history not available", http.StatusServiceUnavailable)
		return
	}

	handler := handlers.NewQueryHistoryHandler(repo, s.logger)
	handler.GetQueryHistoryStats(w, r)
}
