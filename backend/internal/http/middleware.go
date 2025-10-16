// internal/http/middleware.go
package http

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"insightiq/backend/internal/services"
)

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		s.logger.Info("HTTP Request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start),
			"remote_addr", r.RemoteAddr,
		)
	})
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		allowedOrigins := []string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
			"http://frontend:3000",
		}

		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		if allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// SecurityMiddleware adds security headers and limits
func (s *Server) securityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'")

		// Request size limit (10MB)
		r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

		next.ServeHTTP(w, r)
	})
}

// Simple rate limiter
type rateLimiter struct {
	visitors map[string]*visitor
	mu       sync.Mutex
}

type visitor struct {
	lastSeen time.Time
	count    int
}

var limiter = &rateLimiter{
	visitors: make(map[string]*visitor),
}

// RateLimitMiddleware implements basic rate limiting
func (s *Server) rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			ip = forwarded
		}

		limiter.mu.Lock()
		v, exists := limiter.visitors[ip]
		if !exists {
			limiter.visitors[ip] = &visitor{lastSeen: time.Now(), count: 1}
			limiter.mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		// Reset count if last request was more than 1 minute ago
		if time.Since(v.lastSeen) > time.Minute {
			v.count = 1
			v.lastSeen = time.Now()
			limiter.mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		// Limit to 60 requests per minute
		if v.count >= 60 {
			limiter.mu.Unlock()
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		v.count++
		v.lastSeen = time.Now()
		limiter.mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware validates JWT tokens and adds user context
func (s *Server) authMiddleware(authService *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

			// Validate token
			claims, err := authService.ValidateToken(token)
			if err != nil {
				s.logger.Warn("Invalid token", "error", err)
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Add user information to context
			ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
			ctx = context.WithValue(ctx, "user_email", claims.Email)
			ctx = context.WithValue(ctx, "user_role", claims.Role)

			// Call next handler with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// OptionalAuthMiddleware adds user context if token is present but doesn't require it
func (s *Server) optionalAuthMiddleware(authService *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				// No token, continue without auth context
				next.ServeHTTP(w, r)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				// Invalid format, continue without auth context
				next.ServeHTTP(w, r)
				return
			}

			token := parts[1]

			// Validate token
			claims, err := authService.ValidateToken(token)
			if err != nil {
				// Invalid token, continue without auth context
				next.ServeHTTP(w, r)
				return
			}

			// Add user information to context
			ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
			ctx = context.WithValue(ctx, "user_email", claims.Email)
			ctx = context.WithValue(ctx, "user_role", claims.Role)

			// Call next handler with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RoleMiddleware checks if user has required role
func (s *Server) roleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value("user_role").(string)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if user role is in allowed roles
			allowed := false
			for _, role := range allowedRoles {
				if userRole == role {
					allowed = true
					break
				}
			}

			if !allowed {
				http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
