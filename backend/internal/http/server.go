// internal/http/server.go
package http

import (
	"log/slog"
	"net/http"

	"insightiq/backend/internal/services"
)

type Server struct {
	analyticsService *services.AnalyticsService
	voiceService     *services.VoiceService
	logger           *slog.Logger
	mux              *http.ServeMux
}

func NewServer(analytics *services.AnalyticsService, voice *services.VoiceService, logger *slog.Logger) *Server {
	s := &Server{
		analyticsService: analytics,
		voiceService:     voice,
		logger:           logger,
		mux:              http.NewServeMux(),
	}

	s.setupRoutes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Apply middleware
	handler := s.corsMiddleware(s.loggingMiddleware(s.mux))
	handler.ServeHTTP(w, r)
}

func (s *Server) setupRoutes() {
	// API routes
	s.mux.HandleFunc("/api/health", s.handleHealth)
	s.mux.HandleFunc("/api/query", s.handleTextQuery)
	s.mux.HandleFunc("/api/voice", s.handleVoiceQuery)
	s.mux.HandleFunc("/api/sql", s.handleSQLQuery)
}
