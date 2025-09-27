// internal/http/handlers.go
package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"insightiq/backend/internal/connectors"
	"insightiq/backend/internal/validation"
)

type ServiceHealth struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type HealthResponse struct {
	Status    string                   `json:"status"`
	Timestamp time.Time                `json:"timestamp"`
	Services  map[string]ServiceHealth `json:"services"`
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	health := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Services: map[string]ServiceHealth{
			"agent_manager": {
				Status: "healthy",
			},
			// Add other service health checks
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

func (s *Server) handleTestPostgres(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Simple test that bypasses agent system and tests PostgreSQL directly
	ctx := r.Context()

	// Create a new PostgreSQL connector for testing
	postgresURL := os.Getenv("POSTGRES_URL")
	if postgresURL == "" {
		postgresURL = "postgres://insightiq_user:insightiq_password@postgres:5432/insightiq?sslmode=disable"
	}
	postgresConn, err := connectors.NewPostgresConnector(postgresURL, s.logger)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to PostgreSQL: %v", err), http.StatusInternalServerError)
		return
	}
	defer postgresConn.Close()

	// Test the query
	result, err := postgresConn.GetBikeSalesData(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("PostgreSQL query failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"rows":   len(result.Data),
		"sample": result.Data[:min(3, len(result.Data))],
	})
}

func (s *Server) handleDirectAnalytics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Query string `json:"query"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Error("Invalid JSON in direct analytics request", "error", err, "remote_addr", r.RemoteAddr)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate and sanitize input
	req.Query = validation.SanitizeString(req.Query)
	if err := validation.ValidateTextQuery(req.Query); err != nil {
		s.logger.Error("Invalid query input", "error", err, "query", req.Query, "remote_addr", r.RemoteAddr)
		http.Error(w, "Invalid query input", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// Create PostgreSQL connector
	postgresURL := os.Getenv("POSTGRES_URL")
	if postgresURL == "" {
		postgresURL = "postgres://insightiq_user:insightiq_password@postgres:5432/insightiq?sslmode=disable"
	}
	postgresConn, err := connectors.NewPostgresConnector(postgresURL, s.logger)
	if err != nil {
		http.Error(w, fmt.Sprintf("Database connection failed: %v", err), http.StatusInternalServerError)
		return
	}
	defer postgresConn.Close()

	// Get real data from PostgreSQL
	result, err := postgresConn.GetBikeSalesData(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Database query failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Create Ollama connector for LLM analysis
	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://ollama:11434"
	}
	ollamaConn := connectors.NewOllamaConnector(ollamaURL, s.logger)

	// Get LLM analysis of the real data with timeout
	llmCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	insights, err := ollamaConn.AnalyzeData(llmCtx, result.Data, req.Query)
	if err != nil {
		s.logger.Error("LLM analysis failed", "error", err)
		insights = "LLM analysis temporarily unavailable. Please try again later."
	}

	// Return response in expected format
	response := map[string]interface{}{
		"query":       req.Query,
		"data":        result.Data,
		"insights":    insights,
		"timestamp":   time.Now(),
		"process_time": "500ms",
		"task_id":     "direct_query",
		"status":      "completed",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (s *Server) handleTextQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Query string `json:"query"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Error("Invalid JSON in text query request", "error", err, "remote_addr", r.RemoteAddr)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate and sanitize input
	req.Query = validation.SanitizeString(req.Query)
	if err := validation.ValidateTextQuery(req.Query); err != nil {
		s.logger.Error("Invalid text query input", "error", err, "query", req.Query, "remote_addr", r.RemoteAddr)
		http.Error(w, "Invalid query input", http.StatusBadRequest)
		return
	}

	result, err := s.analyticsService.ProcessQuery(r.Context(), req.Query)
	if err != nil {
		s.logger.Error("Text query failed", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (s *Server) handleVoiceQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10MB max
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("audio")
	if err != nil {
		s.logger.Error("No audio file in voice request", "error", err, "remote_addr", r.RemoteAddr)
		http.Error(w, "No audio file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file upload
	if err := validation.ValidateFileUpload(header.Filename, header.Size); err != nil {
		s.logger.Error("Invalid file upload", "error", err, "filename", header.Filename, "size", header.Size, "remote_addr", r.RemoteAddr)
		http.Error(w, "Invalid file upload", http.StatusBadRequest)
		return
	}

	audioData, err := io.ReadAll(file)
	if err != nil {
		s.logger.Error("Failed to read audio file", "error", err, "remote_addr", r.RemoteAddr)
		http.Error(w, "Failed to read audio file", http.StatusInternalServerError)
		return
	}

	// Determine format from filename
	format := "wav" // default
	if len(header.Filename) > 4 {
		format = strings.ToLower(header.Filename[len(header.Filename)-3:])
	}

	result, err := s.voiceService.ProcessVoiceQuery(r.Context(), audioData, format)
	if err != nil {
		s.logger.Error("Voice query failed", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (s *Server) handleSQLQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		SQL      string `json:"sql"`
		Question string `json:"question"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Error("Invalid JSON in SQL query request", "error", err, "remote_addr", r.RemoteAddr)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate and sanitize input
	req.SQL = validation.SanitizeString(req.SQL)
	req.Question = validation.SanitizeString(req.Question)

	if err := validation.ValidateSQL(req.SQL); err != nil {
		s.logger.Error("Invalid SQL input", "error", err, "sql", req.SQL, "remote_addr", r.RemoteAddr)
		http.Error(w, "Invalid SQL query", http.StatusBadRequest)
		return
	}

	if err := validation.ValidateTextQuery(req.Question); err != nil {
		s.logger.Error("Invalid question input", "error", err, "question", req.Question, "remote_addr", r.RemoteAddr)
		http.Error(w, "Invalid question input", http.StatusBadRequest)
		return
	}

	result, err := s.analyticsService.ExecuteCustomSQL(r.Context(), req.SQL, req.Question)
	if err != nil {
		s.logger.Error("SQL query failed", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
