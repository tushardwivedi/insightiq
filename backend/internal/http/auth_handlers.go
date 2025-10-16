package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"insightiq/backend/internal/models"
	"insightiq/backend/internal/services"
)

type AuthHandlers struct {
	authService *services.AuthService
	server      *Server
}

func NewAuthHandlers(authService *services.AuthService, server *Server) *AuthHandlers {
	return &AuthHandlers{
		authService: authService,
		server:      server,
	}
}

// handleRegister handles user registration
func (h *AuthHandlers) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.server.logger.Error("Failed to decode register request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" || req.Name == "" {
		http.Error(w, "Email, password, and name are required", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Register(r.Context(), req)
	if err != nil {
		if err == services.ErrWeakPassword {
			http.Error(w, "Password must be at least 8 characters long", http.StatusBadRequest)
			return
		}
		if err == services.ErrUserAlreadyExists {
			http.Error(w, "User with this email already exists", http.StatusConflict)
			return
		}
		h.server.logger.Error("Failed to register user", "error", err)
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User registered successfully",
		"user":    user,
	})
}

// handleLogin handles user login
func (h *AuthHandlers) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.server.logger.Error("Failed to decode login request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	response, err := h.authService.Login(r.Context(), req)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
		h.server.logger.Error("Failed to login user", "error", err)
		http.Error(w, "Failed to login", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleLogout handles user logout (client-side token deletion)
func (h *AuthHandlers) handleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// In a JWT-based system, logout is primarily handled client-side
	// by deleting the token. This endpoint is mainly for logging purposes.

	// Extract user info from context if available (set by auth middleware)
	userID := r.Context().Value("user_id")
	if userID != nil {
		h.server.logger.Info("User logged out", "user_id", userID)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logged out successfully",
	})
}

// handleRefreshToken handles token refresh
func (h *AuthHandlers) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
		return
	}

	token := parts[1]
	newToken, err := h.authService.RefreshToken(token)
	if err != nil {
		http.Error(w, "Failed to refresh token", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": newToken,
	})
}

// handleGetCurrentUser handles getting the current authenticated user
func (h *AuthHandlers) handleGetCurrentUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	user, err := h.authService.GetUser(r.Context(), userID)
	if err != nil {
		h.server.logger.Error("Failed to get user", "error", err)
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// handleChangePassword handles password change
func (h *AuthHandlers) handleChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.server.logger.Error("Failed to decode password change request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		http.Error(w, "Old password and new password are required", http.StatusBadRequest)
		return
	}

	if err := h.authService.ChangePassword(r.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		if err == services.ErrInvalidCredentials {
			http.Error(w, "Invalid old password", http.StatusUnauthorized)
			return
		}
		if err == services.ErrWeakPassword {
			http.Error(w, "New password must be at least 8 characters long", http.StatusBadRequest)
			return
		}
		h.server.logger.Error("Failed to change password", "error", err)
		http.Error(w, "Failed to change password", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password changed successfully",
	})
}
