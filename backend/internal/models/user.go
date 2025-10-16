package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID        string    `db:"id" json:"id"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"-"` // Never expose in JSON
	Name      string    `db:"name" json:"name"`
	Role      string    `db:"role" json:"role"` // admin, user
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	LastLogin *time.Time `db:"last_login" json:"last_login,omitempty"`
}

// CreateUserRequest represents the payload for creating a user
type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     string `json:"role,omitempty"` // Optional, defaults to "user"
}

// LoginRequest represents the payload for user login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents the response after successful login
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
