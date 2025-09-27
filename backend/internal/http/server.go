// internal/http/server.go
package http

import (
	"log/slog"
	"net/http"
	"strings"

	"insightiq/backend/internal/services"
)

type Server struct {
	analyticsService *services.AnalyticsService
	voiceService     *services.VoiceService
	connectorService *services.ConnectorService
	plannerService   *services.PlannerService
	logger           *slog.Logger
	mux              *http.ServeMux
}

func NewServer(analytics *services.AnalyticsService, voice *services.VoiceService, connector *services.ConnectorService, planner *services.PlannerService, logger *slog.Logger) *Server {
	s := &Server{
		analyticsService: analytics,
		voiceService:     voice,
		connectorService: connector,
		plannerService:   planner,
		logger:           logger,
		mux:              http.NewServeMux(),
	}

	s.setupRoutes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Apply security middleware stack
	handler := s.corsMiddleware(
		s.securityMiddleware(
			s.rateLimitMiddleware(
				s.loggingMiddleware(s.mux))))
	handler.ServeHTTP(w, r)
}

func (s *Server) setupRoutes() {
	// API routes
	s.mux.HandleFunc("/api/health", s.handleHealth)
	s.mux.HandleFunc("/api/test-postgres", s.handleTestPostgres)
	s.mux.HandleFunc("/api/direct-query", s.handleDirectAnalytics)
	s.mux.HandleFunc("/api/query", s.handleTextQuery)
	s.mux.HandleFunc("/api/voice", s.handleVoiceQuery)
	s.mux.HandleFunc("/api/sql", s.handleSQLQuery)

	// Connector routes
	if s.connectorService != nil {
		connectorHandlers := NewConnectorHandlers(s.connectorService, s)
		s.mux.HandleFunc("/api/connectors", s.routeConnectors(connectorHandlers))
		s.mux.HandleFunc("/api/connectors/", s.routeConnectorsByID(connectorHandlers))
	}

	// Planner routes
	if s.plannerService != nil {
		plannerHandlers := NewPlannerHandlers(s.plannerService, s)
		s.mux.HandleFunc("/api/planner/parse-intent", plannerHandlers.handleParseIntent)
		s.mux.HandleFunc("/api/planner/analyze-query", plannerHandlers.handleAnalyzeQuery)
	}
}

// routeConnectors handles /api/connectors (collection endpoints)
func (s *Server) routeConnectors(handlers *ConnectorHandlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/connectors" {
			http.NotFound(w, r)
			return
		}

		switch r.Method {
		case http.MethodGet:
			handlers.handleGetConnectors(w, r)
		case http.MethodPost:
			handlers.handleCreateConnector(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// routeConnectorsByID handles /api/connectors/{id}/* endpoints
func (s *Server) routeConnectorsByID(handlers *ConnectorHandlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Handle test config endpoint: /api/connectors/test-config
		if path == "/api/connectors/test-config" {
			handlers.handleTestConnectorConfig(w, r)
			return
		}

		// Extract connector ID and action from path
		if len(path) < len("/api/connectors/") {
			http.NotFound(w, r)
			return
		}

		// Handle specific connector endpoints
		if strings.HasSuffix(path, "/test") {
			handlers.handleTestConnector(w, r)
		} else if strings.HasSuffix(path, "/data") {
			handlers.handleGetConnectorData(w, r)
		} else if strings.Count(path, "/") == 3 { // /api/connectors/{id}
			switch r.Method {
			case http.MethodGet:
				handlers.handleGetConnector(w, r)
			case http.MethodPut:
				handlers.handleUpdateConnector(w, r)
			case http.MethodDelete:
				handlers.handleDeleteConnector(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		} else {
			http.NotFound(w, r)
		}
	}
}
