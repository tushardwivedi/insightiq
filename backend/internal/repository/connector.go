package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"insightiq/backend/internal/models"
)

type ConnectorRepository struct {
	db *sqlx.DB
}

func NewConnectorRepository(db *sqlx.DB) *ConnectorRepository {
	return &ConnectorRepository{db: db}
}

// CreateTables creates the connector tables if they don't exist
func (r *ConnectorRepository) CreateTables(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS data_connectors (
		id VARCHAR(36) PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		type VARCHAR(20) NOT NULL,
		status VARCHAR(20) NOT NULL DEFAULT 'disconnected',
		config JSONB NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		last_tested TIMESTAMP WITH TIME ZONE
	);

	CREATE INDEX IF NOT EXISTS idx_data_connectors_type ON data_connectors(type);
	CREATE INDEX IF NOT EXISTS idx_data_connectors_status ON data_connectors(status);
	CREATE INDEX IF NOT EXISTS idx_data_connectors_created_at ON data_connectors(created_at);
	`

	_, err := r.db.ExecContext(ctx, query)
	return err
}

// Create creates a new data connector
func (r *ConnectorRepository) Create(ctx context.Context, connector *models.DataConnector) error {
	connector.ID = uuid.New().String()
	connector.Status = models.ConnectorStatusDisconnected
	connector.CreatedAt = time.Now()
	connector.UpdatedAt = time.Now()

	query := `
		INSERT INTO data_connectors (id, name, type, status, config, created_at, updated_at)
		VALUES (:id, :name, :type, :status, :config, :created_at, :updated_at)
	`

	_, err := r.db.NamedExecContext(ctx, query, connector)
	return err
}

// GetByID retrieves a connector by ID
func (r *ConnectorRepository) GetByID(ctx context.Context, id string) (*models.DataConnector, error) {
	var connector models.DataConnector
	query := `
		SELECT id, name, type, status, config, created_at, updated_at, last_tested
		FROM data_connectors
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, &connector, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &connector, nil
}

// GetAll retrieves all connectors
func (r *ConnectorRepository) GetAll(ctx context.Context) ([]*models.DataConnector, error) {
	var connectors []*models.DataConnector
	query := `
		SELECT id, name, type, status, config, created_at, updated_at, last_tested
		FROM data_connectors
		ORDER BY created_at DESC
	`

	err := r.db.SelectContext(ctx, &connectors, query)
	if err != nil {
		return nil, err
	}

	return connectors, nil
}

// GetByType retrieves connectors by type
func (r *ConnectorRepository) GetByType(ctx context.Context, connectorType models.ConnectorType) ([]*models.DataConnector, error) {
	var connectors []*models.DataConnector
	query := `
		SELECT id, name, type, status, config, created_at, updated_at, last_tested
		FROM data_connectors
		WHERE type = $1
		ORDER BY created_at DESC
	`

	err := r.db.SelectContext(ctx, &connectors, query, connectorType)
	if err != nil {
		return nil, err
	}

	return connectors, nil
}

// GetActive retrieves all active (connected) connectors
func (r *ConnectorRepository) GetActive(ctx context.Context) ([]*models.DataConnector, error) {
	var connectors []*models.DataConnector
	query := `
		SELECT id, name, type, status, config, created_at, updated_at, last_tested
		FROM data_connectors
		WHERE status = $1
		ORDER BY created_at DESC
	`

	err := r.db.SelectContext(ctx, &connectors, query, models.ConnectorStatusConnected)
	if err != nil {
		return nil, err
	}

	return connectors, nil
}

// Update updates an existing connector
func (r *ConnectorRepository) Update(ctx context.Context, id string, updates *models.UpdateConnectorRequest) (*models.DataConnector, error) {
	// Start a transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Build dynamic update query
	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{}
	argIndex := 1

	if updates.Name != "" {
		setParts = append(setParts, "name = $"+getArgNumber(argIndex))
		args = append(args, updates.Name)
		argIndex++
	}

	if updates.Config != nil {
		setParts = append(setParts, "config = $"+getArgNumber(argIndex))
		args = append(args, updates.Config)
		argIndex++
	}

	if updates.Status != "" {
		setParts = append(setParts, "status = $"+getArgNumber(argIndex))
		args = append(args, updates.Status)
		argIndex++
	}

	// Add the ID for the WHERE clause
	args = append(args, id)
	whereArgIndex := argIndex

	query := `
		UPDATE data_connectors
		SET ` + joinStrings(setParts, ", ") + `
		WHERE id = $` + getArgNumber(whereArgIndex) + `
		RETURNING id, name, type, status, config, created_at, updated_at, last_tested
	`

	var connector models.DataConnector
	err = tx.GetContext(ctx, &connector, query, args...)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &connector, nil
}

// UpdateStatus updates the connector status and last tested time
func (r *ConnectorRepository) UpdateStatus(ctx context.Context, id string, status models.ConnectorStatus) error {
	query := `
		UPDATE data_connectors
		SET status = $1, last_tested = NOW(), updated_at = NOW()
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

// Delete deletes a connector
func (r *ConnectorRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM data_connectors WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Helper functions
func getArgNumber(index int) string {
	return string(rune('0' + index))
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}