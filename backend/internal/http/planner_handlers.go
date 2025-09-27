package http

import (
	"encoding/json"
	"net/http"

	"insightiq/backend/internal/models"
	"insightiq/backend/internal/services"
)

type PlannerHandlers struct {
	plannerService *services.PlannerService
	server         *Server
}

func NewPlannerHandlers(plannerService *services.PlannerService, server *Server) *PlannerHandlers {
	return &PlannerHandlers{
		plannerService: plannerService,
		server:         server,
	}
}

// handleParseIntent handles intent parsing requests
func (h *PlannerHandlers) handleParseIntent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.PlannerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.server.logger.Error("Failed to decode planner request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Query == "" {
		http.Error(w, "Query is required", http.StatusBadRequest)
		return
	}

	h.server.logger.Info("Processing intent parsing request", "query", req.Query)

	response, err := h.plannerService.ParseIntent(r.Context(), &req)
	if err != nil {
		h.server.logger.Error("Failed to parse intent", "error", err)
		http.Error(w, "Failed to parse intent", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.server.logger.Error("Failed to encode response", "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	h.server.logger.Info("Intent parsing completed successfully",
		"intent_type", response.Intent.Type,
		"confidence", response.Confidence,
		"steps", len(response.TaskGraph.Steps))
}

// handleAnalyzeQuery provides detailed query analysis
func (h *PlannerHandlers) handleAnalyzeQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Query string `json:"query"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.server.logger.Error("Failed to decode analyze request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Query == "" {
		http.Error(w, "Query is required", http.StatusBadRequest)
		return
	}

	h.server.logger.Info("Processing query analysis request", "query", req.Query)

	plannerReq := &models.PlannerRequest{
		Query: req.Query,
	}

	response, err := h.plannerService.ParseIntent(r.Context(), plannerReq)
	if err != nil {
		h.server.logger.Error("Failed to analyze query", "error", err)
		http.Error(w, "Failed to analyze query", http.StatusInternalServerError)
		return
	}

	// Return a simplified analysis response
	analysisResponse := map[string]interface{}{
		"intent": map[string]interface{}{
			"type":       response.Intent.Type,
			"confidence": response.Confidence,
			"entities":   response.Intent.Entities,
		},
		"parsed_query": response.Intent.ParsedQuery,
		"task_graph": map[string]interface{}{
			"steps":          len(response.TaskGraph.Steps),
			"estimated_time": response.TaskGraph.EstimatedTime.String(),
			"step_details":   response.TaskGraph.Steps,
		},
		"metadata": response.Metadata,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(analysisResponse); err != nil {
		h.server.logger.Error("Failed to encode analysis response", "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	h.server.logger.Info("Query analysis completed successfully",
		"intent_type", response.Intent.Type,
		"confidence", response.Confidence)
}