package services

import (
	"context"
	"errors"
	"log/slog"

	"insightiq/backend/internal/auth"
	"insightiq/backend/internal/models"
	"insightiq/backend/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrWeakPassword       = errors.New("password does not meet strength requirements")
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
)

// AuthService handles authentication operations
type AuthService struct {
	userRepo   *repository.UserRepository
	jwtManager *auth.JWTManager
	logger     *slog.Logger
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo *repository.UserRepository, jwtManager *auth.JWTManager, logger *slog.Logger) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
		logger:     logger,
	}
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, req models.CreateUserRequest) (*models.User, error) {
	// Validate password strength
	if !auth.ValidatePasswordStrength(req.Password) {
		return nil, ErrWeakPassword
	}

	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil && err != repository.ErrUserNotFound {
		s.logger.Error("Failed to check existing user", "error", err)
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		s.logger.Error("Failed to hash password", "error", err)
		return nil, err
	}

	// Set default role if not provided
	role := req.Role
	if role == "" {
		role = "user"
	}

	// Create user
	user := &models.User{
		Email:    req.Email,
		Password: hashedPassword,
		Name:     req.Name,
		Role:     role,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Error("Failed to create user", "error", err)
		return nil, err
	}

	s.logger.Info("User registered successfully", "user_id", user.ID, "email", user.Email)
	return user, nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return nil, ErrInvalidCredentials
		}
		s.logger.Error("Failed to get user by email", "error", err)
		return nil, err
	}

	// Compare passwords
	if err := auth.ComparePassword(user.Password, req.Password); err != nil {
		s.logger.Warn("Invalid login attempt", "email", req.Email)
		return nil, ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := s.jwtManager.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		s.logger.Error("Failed to generate token", "error", err)
		return nil, err
	}

	// Update last login
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		s.logger.Warn("Failed to update last login", "error", err)
		// Don't fail the login if this fails
	}

	s.logger.Info("User logged in successfully", "user_id", user.ID, "email", user.Email)

	return &models.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *AuthService) ValidateToken(tokenString string) (*auth.Claims, error) {
	return s.jwtManager.ValidateToken(tokenString)
}

// RefreshToken generates a new token from an existing valid token
func (s *AuthService) RefreshToken(tokenString string) (string, error) {
	return s.jwtManager.RefreshToken(tokenString)
}

// GetUser retrieves a user by ID
func (s *AuthService) GetUser(ctx context.Context, userID string) (*models.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

// ChangePassword changes a user's password
func (s *AuthService) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify old password
	if err := auth.ComparePassword(user.Password, oldPassword); err != nil {
		return ErrInvalidCredentials
	}

	// Validate new password strength
	if !auth.ValidatePasswordStrength(newPassword) {
		return ErrWeakPassword
	}

	// Hash new password
	hashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		s.logger.Error("Failed to hash new password", "error", err)
		return err
	}

	// Update password
	if err := s.userRepo.UpdatePassword(ctx, userID, hashedPassword); err != nil {
		s.logger.Error("Failed to update password", "error", err)
		return err
	}

	s.logger.Info("Password changed successfully", "user_id", userID)
	return nil
}
