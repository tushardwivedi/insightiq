package http

import (
	"net/http"

	"insightiq/backend/internal/http/handlers"
	"insightiq/backend/internal/usecases"
)

// handleGDPRUserData routes GDPR user data requests
// DELETE /api/gdpr/user-data - Right to Erasure (GDPR Article 17)
func (s *Server) handleGDPRUserData(w http.ResponseWriter, r *http.Request) {
	useCase, ok := s.queryHistoryUseCase.(*usecases.QueryHistoryUseCase)
	if !ok {
		http.Error(w, "GDPR services not available", http.StatusServiceUnavailable)
		return
	}

	handler := handlers.NewGDPRHandler(useCase, s.logger)

	switch r.Method {
	case http.MethodDelete:
		handler.HandleDeleteUserData(w, r)
	default:
		http.Error(w, "Method not allowed. Use DELETE to erase user data", http.StatusMethodNotAllowed)
	}
}

// handleGDPRExport handles GDPR data export requests
// GET /api/gdpr/export?format=json|csv - Right to Access (GDPR Article 15)
func (s *Server) handleGDPRExport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	useCase, ok := s.queryHistoryUseCase.(*usecases.QueryHistoryUseCase)
	if !ok {
		http.Error(w, "GDPR services not available", http.StatusServiceUnavailable)
		return
	}

	handler := handlers.NewGDPRHandler(useCase, s.logger)
	handler.HandleExportUserData(w, r)
}

// handleGDPRPortableData handles GDPR data portability requests
// GET /api/gdpr/portable-data - Right to Data Portability (GDPR Article 20)
func (s *Server) handleGDPRPortableData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	useCase, ok := s.queryHistoryUseCase.(*usecases.QueryHistoryUseCase)
	if !ok {
		http.Error(w, "GDPR services not available", http.StatusServiceUnavailable)
		return
	}

	handler := handlers.NewGDPRHandler(useCase, s.logger)
	handler.HandleDataPortabilityRequest(w, r)
}

// handleGDPRComplianceReport generates a GDPR compliance report
// GET /api/gdpr/compliance-report
func (s *Server) handleGDPRComplianceReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	useCase, ok := s.queryHistoryUseCase.(*usecases.QueryHistoryUseCase)
	if !ok {
		http.Error(w, "GDPR services not available", http.StatusServiceUnavailable)
		return
	}

	handler := handlers.NewGDPRHandler(useCase, s.logger)
	handler.HandleComplianceReport(w, r)
}
