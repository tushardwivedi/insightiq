package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"insightiq/backend/internal/models"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrDuplicateEmail    = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateTables creates the users table if it doesn't exist
func (r *UserRepository) CreateTables(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(255) PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			name VARCHAR(255) NOT NULL,
			role VARCHAR(50) NOT NULL DEFAULT 'user',
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			last_login TIMESTAMP
		);

		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
		CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
	`

	_, err := r.db.ExecContext(ctx, query)
	return err
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if user.Role == "" {
		user.Role = "user"
	}

	query := `
		INSERT INTO users (id, email, password, name, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.Password, user.Name, user.Role,
		user.CreatedAt, user.UpdatedAt,
	)

	if err != nil {
		// Check for duplicate email
		if err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` {
			return ErrDuplicateEmail
		}
		return err
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	query := `SELECT id, email, password, name, role, created_at, updated_at, last_login FROM users WHERE id = $1`

	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, email, password, name, role, created_at, updated_at, last_login FROM users WHERE email = $1`

	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()

	query := `
		UPDATE users
		SET email = $1, name = $2, role = $3, updated_at = $4
		WHERE id = $5
	`

	result, err := r.db.ExecContext(ctx, query,
		user.Email, user.Name, user.Role, user.UpdatedAt, user.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrUserNotFound
	}

	return nil
}

// UpdatePassword updates a user's password
func (r *UserRepository) UpdatePassword(ctx context.Context, userID, hashedPassword string) error {
	query := `
		UPDATE users
		SET password = $1, updated_at = $2
		WHERE id = $3
	`

	result, err := r.db.ExecContext(ctx, query, hashedPassword, time.Now(), userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrUserNotFound
	}

	return nil
}

// UpdateLastLogin updates the user's last login timestamp
func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	now := time.Now()
	query := `
		UPDATE users
		SET last_login = $1, updated_at = $2
		WHERE id = $3
	`

	result, err := r.db.ExecContext(ctx, query, now, now, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrUserNotFound
	}

	return nil
}

// Delete deletes a user
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrUserNotFound
	}

	return nil
}

// List retrieves all users
func (r *UserRepository) List(ctx context.Context) ([]models.User, error) {
	var users []models.User
	query := `SELECT id, email, password, name, role, created_at, updated_at, last_login FROM users ORDER BY created_at DESC`

	err := r.db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, err
	}

	return users, nil
}
