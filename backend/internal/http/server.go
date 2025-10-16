// internal/http/server.go
package http

import (
	"log/slog"
	"net/http"
	"strings"

	"insightiq/backend/internal/auth"
	"insightiq/backend/internal/services"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type Server struct {
	analyticsService *services.AnalyticsService
	voiceService     *services.VoiceService
	connectorService *services.ConnectorService
	plannerService   *services.PlannerService
	authService      *services.AuthService
	logger           *slog.Logger
	mux              *http.ServeMux
}

func NewServer(analytics *services.AnalyticsService, voice *services.VoiceService, connector *services.ConnectorService, planner *services.PlannerService, auth *services.AuthService, logger *slog.Logger) *Server {
	s := &Server{
		analyticsService: analytics,
		voiceService:     voice,
		connectorService: connector,
		plannerService:   planner,
		authService:      auth,
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
	// Initialize SuperTokens
	if err := auth.InitSuperTokens(); err != nil {
		s.logger.Error("Failed to initialize SuperTokens", "error", err)
	}

	// SuperTokens routes - OAuth and email/password
	s.mux.Handle("/auth/", supertokens.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// SuperTokens handles all /auth/* routes including OAuth callbacks
		w.WriteHeader(http.StatusNotFound)
	})))

	// Public routes (no authentication required)
	s.mux.HandleFunc("/api/health", s.handleHealth)

	// Auth routes (public) - Keep existing JWT-based auth as fallback
	if s.authService != nil {
		authHandlers := NewAuthHandlers(s.authService, s)
		s.mux.HandleFunc("/api/auth/register", authHandlers.handleRegister)
		s.mux.HandleFunc("/api/auth/login", authHandlers.handleLogin)
		s.mux.HandleFunc("/api/auth/logout", authHandlers.handleLogout)
		s.mux.HandleFunc("/api/auth/refresh", authHandlers.handleRefreshToken)

		// Protected auth routes
		s.mux.HandleFunc("/api/auth/me", s.withAuth(authHandlers.handleGetCurrentUser))
		s.mux.HandleFunc("/api/auth/change-password", s.withAuth(authHandlers.handleChangePassword))
	}

	// Protected API routes (require authentication)
	s.mux.HandleFunc("/api/test-postgres", s.withAuth(s.handleTestPostgres))
	s.mux.HandleFunc("/api/direct-query", s.withAuth(s.handleDirectAnalytics))
	s.mux.HandleFunc("/api/query", s.withAuth(s.handleTextQuery))
	s.mux.HandleFunc("/api/voice", s.withAuth(s.handleVoiceQuery))
	s.mux.HandleFunc("/api/sql", s.withAuth(s.handleSQLQuery))

	// Protected connector routes
	if s.connectorService != nil {
		connectorHandlers := NewConnectorHandlers(s.connectorService, s)
		s.mux.HandleFunc("/api/connectors", s.withAuth(s.routeConnectors(connectorHandlers)))
		s.mux.HandleFunc("/api/connectors/", s.withAuth(s.routeConnectorsByID(connectorHandlers)))
	}

	// Protected planner routes
	if s.plannerService != nil {
		plannerHandlers := NewPlannerHandlers(s.plannerService, s)
		s.mux.HandleFunc("/api/planner/parse-intent", s.withAuth(plannerHandlers.handleParseIntent))
		s.mux.HandleFunc("/api/planner/analyze-query", s.withAuth(plannerHandlers.handleAnalyzeQuery))
	}
}

// withAuth wraps a handler with authentication middleware
func (s *Server) withAuth(handler http.HandlerFunc) http.HandlerFunc {
	if s.authService == nil {
		return handler
	}
	authMiddleware := s.authMiddleware(s.authService)
	return func(w http.ResponseWriter, r *http.Request) {
		authMiddleware(http.HandlerFunc(handler)).ServeHTTP(w, r)
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
