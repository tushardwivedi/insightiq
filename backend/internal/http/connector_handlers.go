package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"insightiq/backend/internal/models"
	"insightiq/backend/internal/services"
	"insightiq/backend/internal/validation"
)

type ConnectorHandlers struct {
	connectorService *services.ConnectorService
	server           *Server
}

func NewConnectorHandlers(connectorService *services.ConnectorService, server *Server) *ConnectorHandlers {
	return &ConnectorHandlers{
		connectorService: connectorService,
		server:           server,
	}
}

// handleGetConnectors retrieves all connectors
func (h *ConnectorHandlers) handleGetConnectors(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	connectors, err := h.connectorService.GetConnectors(r.Context())
	if err != nil {
		h.server.logger.Error("Failed to get connectors", "error", err)
		http.Error(w, "Failed to retrieve connectors", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"data":    connectors,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleCreateConnector creates a new connector
func (h *ConnectorHandlers) handleCreateConnector(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.CreateConnectorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.server.logger.Error("Invalid JSON in create connector request", "error", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate input
	req.Name = validation.SanitizeString(req.Name)
	if err := validation.ValidateConnectorName(req.Name); err != nil {
		h.server.logger.Error("Invalid connector name", "error", err, "name", req.Name)
		http.Error(w, "Invalid connector name", http.StatusBadRequest)
		return
	}

	connector, err := h.connectorService.CreateConnector(r.Context(), &req)
	if err != nil {
		h.server.logger.Error("Failed to create connector", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(connector)
}

// handleGetConnector retrieves a specific connector
func (h *ConnectorHandlers) handleGetConnector(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := h.extractConnectorID(r.URL.Path)
	if id == "" {
		http.Error(w, "Invalid connector ID", http.StatusBadRequest)
		return
	}

	connector, err := h.connectorService.GetConnector(r.Context(), id)
	if err != nil {
		h.server.logger.Error("Failed to get connector", "error", err, "id", id)
		http.Error(w, "Failed to retrieve connector", http.StatusInternalServerError)
		return
	}

	if connector == nil {
		http.Error(w, "Connector not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(connector)
}

// handleUpdateConnector updates an existing connector
func (h *ConnectorHandlers) handleUpdateConnector(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := h.extractConnectorID(r.URL.Path)
	if id == "" {
		http.Error(w, "Invalid connector ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateConnectorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.server.logger.Error("Invalid JSON in update connector request", "error", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Name != "" {
		req.Name = validation.SanitizeString(req.Name)
		if err := validation.ValidateConnectorName(req.Name); err != nil {
			h.server.logger.Error("Invalid connector name", "error", err, "name", req.Name)
			http.Error(w, "Invalid connector name", http.StatusBadRequest)
			return
		}
	}

	connector, err := h.connectorService.UpdateConnector(r.Context(), id, &req)
	if err != nil {
		h.server.logger.Error("Failed to update connector", "error", err, "id", id)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(connector)
}

// handleDeleteConnector deletes a connector
func (h *ConnectorHandlers) handleDeleteConnector(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := h.extractConnectorID(r.URL.Path)
	if id == "" {
		http.Error(w, "Invalid connector ID", http.StatusBadRequest)
		return
	}

	err := h.connectorService.DeleteConnector(r.Context(), id)
	if err != nil {
		h.server.logger.Error("Failed to delete connector", "error", err, "id", id)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handleTestConnector tests an existing connector
func (h *ConnectorHandlers) handleTestConnector(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := h.extractConnectorIDFromTestPath(r.URL.Path)
	if id == "" {
		http.Error(w, "Invalid connector ID", http.StatusBadRequest)
		return
	}

	result, err := h.connectorService.TestConnector(r.Context(), id)
	if err != nil {
		h.server.logger.Error("Failed to test connector", "error", err, "id", id)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// handleTestConnectorConfig tests a connector configuration without saving
func (h *ConnectorHandlers) handleTestConnectorConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.TestConnectorConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.server.logger.Error("Invalid JSON in test connector config request", "error", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	result, err := h.connectorService.TestConnectorConfig(r.Context(), &req)
	if err != nil {
		h.server.logger.Error("Failed to test connector config", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// handleGetConnectorData retrieves data from a specific connector
func (h *ConnectorHandlers) handleGetConnectorData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := h.extractConnectorIDFromDataPath(r.URL.Path)
	if id == "" {
		http.Error(w, "Invalid connector ID", http.StatusBadRequest)
		return
	}

	query := r.URL.Query().Get("query")
	if query != "" {
		query = validation.SanitizeString(query)
	}

	// TODO: Implement data retrieval through connector
	// This would involve getting the connector, establishing connection,
	// and executing the query through the appropriate connector type

	response := map[string]interface{}{
		"message": "Connector data retrieval not yet implemented",
		"id":      id,
		"query":   query,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Helper methods to extract IDs from URL paths

func (h *ConnectorHandlers) extractConnectorID(path string) string {
	// Extract ID from /api/connectors/{id}
	parts := strings.Split(path, "/")
	if len(parts) >= 4 && parts[2] == "connectors" {
		return parts[3]
	}
	return ""
}

func (h *ConnectorHandlers) extractConnectorIDFromTestPath(path string) string {
	// Extract ID from /api/connectors/{id}/test
	parts := strings.Split(path, "/")
	if len(parts) >= 5 && parts[2] == "connectors" && parts[4] == "test" {
		return parts[3]
	}
	return ""
}

func (h *ConnectorHandlers) extractConnectorIDFromDataPath(path string) string {
	// Extract ID from /api/connectors/{id}/data
	parts := strings.Split(path, "/")
	if len(parts) >= 5 && parts[2] == "connectors" && parts[4] == "data" {
		return parts[3]
	}
	return ""
}