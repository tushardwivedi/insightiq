package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"insightiq/backend/internal/connectors"
	"insightiq/backend/internal/models"
	"insightiq/backend/internal/repository"
)

type ConnectorService struct {
	repo   *repository.ConnectorRepository
	logger *slog.Logger
}

func NewConnectorService(repo *repository.ConnectorRepository, logger *slog.Logger) *ConnectorService {
	return &ConnectorService{
		repo:   repo,
		logger: logger.With("service", "connector"),
	}
}

// GetConnectors retrieves all connectors
func (s *ConnectorService) GetConnectors(ctx context.Context) ([]*models.DataConnector, error) {
	return s.repo.GetAll(ctx)
}

// GetConnector retrieves a specific connector by ID
func (s *ConnectorService) GetConnector(ctx context.Context, id string) (*models.DataConnector, error) {
	return s.repo.GetByID(ctx, id)
}

// CreateConnector creates a new data connector
func (s *ConnectorService) CreateConnector(ctx context.Context, req *models.CreateConnectorRequest) (*models.DataConnector, error) {
	s.logger.Info("Creating new connector", "name", req.Name, "type", req.Type)

	// Validate the configuration based on connector type
	if err := s.validateConnectorConfig(req.Type, req.Config); err != nil {
		s.logger.Error("Invalid connector configuration", "error", err)
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	connector := &models.DataConnector{
		Name:   req.Name,
		Type:   req.Type,
		Config: req.Config,
	}

	if err := s.repo.Create(ctx, connector); err != nil {
		s.logger.Error("Failed to create connector", "error", err)
		return nil, fmt.Errorf("failed to create connector: %w", err)
	}

	s.logger.Info("Connector created successfully", "id", connector.ID, "name", connector.Name)
	return connector, nil
}

// UpdateConnector updates an existing connector
func (s *ConnectorService) UpdateConnector(ctx context.Context, id string, req *models.UpdateConnectorRequest) (*models.DataConnector, error) {
	s.logger.Info("Updating connector", "id", id)

	// Check if connector exists
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get connector: %w", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("connector not found")
	}

	// Validate configuration if provided
	if req.Config != nil {
		if err := s.validateConnectorConfig(existing.Type, req.Config); err != nil {
			s.logger.Error("Invalid connector configuration", "error", err)
			return nil, fmt.Errorf("invalid configuration: %w", err)
		}
	}

	updated, err := s.repo.Update(ctx, id, req)
	if err != nil {
		s.logger.Error("Failed to update connector", "error", err)
		return nil, fmt.Errorf("failed to update connector: %w", err)
	}

	s.logger.Info("Connector updated successfully", "id", id)
	return updated, nil
}

// DeleteConnector deletes a connector
func (s *ConnectorService) DeleteConnector(ctx context.Context, id string) error {
	s.logger.Info("Deleting connector", "id", id)

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete connector", "error", err)
		return fmt.Errorf("failed to delete connector: %w", err)
	}

	s.logger.Info("Connector deleted successfully", "id", id)
	return nil
}

// TestConnector tests the connection for an existing connector
func (s *ConnectorService) TestConnector(ctx context.Context, id string) (*models.ConnectorTestResult, error) {
	s.logger.Info("Testing connector", "id", id)

	connector, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get connector: %w", err)
	}
	if connector == nil {
		return nil, fmt.Errorf("connector not found")
	}

	// Update status to testing
	s.repo.UpdateStatus(ctx, id, models.ConnectorStatusTesting)

	// Test the connection
	result := s.testConnectorConfig(ctx, connector.Type, connector.Config)

	// Update status based on test result
	newStatus := models.ConnectorStatusError
	if result.Success {
		newStatus = models.ConnectorStatusConnected
	}
	s.repo.UpdateStatus(ctx, id, newStatus)

	s.logger.Info("Connector test completed", "id", id, "success", result.Success)
	return result, nil
}

// TestConnectorConfig tests a connector configuration without saving it
func (s *ConnectorService) TestConnectorConfig(ctx context.Context, req *models.TestConnectorConfigRequest) (*models.ConnectorTestResult, error) {
	s.logger.Info("Testing connector configuration", "type", req.Type)

	// Validate configuration
	if err := s.validateConnectorConfig(req.Type, req.Config); err != nil {
		return &models.ConnectorTestResult{
			Success: false,
			Message: "Invalid configuration",
			Error:   err.Error(),
		}, nil
	}

	result := s.testConnectorConfig(ctx, req.Type, req.Config)
	s.logger.Info("Configuration test completed", "type", req.Type, "success", result.Success)
	return result, nil
}

// GetActiveConnectors retrieves all active connectors
func (s *ConnectorService) GetActiveConnectors(ctx context.Context) ([]*models.DataConnector, error) {
	return s.repo.GetActive(ctx)
}

// GetConnectorsByType retrieves connectors by type
func (s *ConnectorService) GetConnectorsByType(ctx context.Context, connectorType models.ConnectorType) ([]*models.DataConnector, error) {
	return s.repo.GetByType(ctx, connectorType)
}

// validateConnectorConfig validates the configuration based on connector type
func (s *ConnectorService) validateConnectorConfig(connectorType models.ConnectorType, config models.ConnectorConfig) error {
	switch connectorType {
	case models.ConnectorTypeSuperset:
		return s.validateSupersetConfig(config)
	case models.ConnectorTypePostgres:
		return s.validatePostgresConfig(config)
	default:
		return fmt.Errorf("unsupported connector type: %s", connectorType)
	}
}

// validateSupersetConfig validates Superset connector configuration
func (s *ConnectorService) validateSupersetConfig(config models.ConnectorConfig) error {
	url, ok := config["url"].(string)
	if !ok || url == "" {
		return fmt.Errorf("url is required for Superset connector")
	}

	// Check if either username/password or bearer_token is provided
	username, hasUsername := config["username"].(string)
	password, hasPassword := config["password"].(string)
	token, hasToken := config["bearer_token"].(string)

	if (!hasUsername || !hasPassword || username == "" || password == "") && (!hasToken || token == "") {
		return fmt.Errorf("either username/password or bearer_token is required for Superset connector")
	}

	return nil
}

// validatePostgresConfig validates PostgreSQL connector configuration
func (s *ConnectorService) validatePostgresConfig(config models.ConnectorConfig) error {
	url, ok := config["url"].(string)
	if !ok || url == "" {
		return fmt.Errorf("url is required for PostgreSQL connector")
	}

	return nil
}

// testConnectorConfig tests the actual connection to the data source
func (s *ConnectorService) testConnectorConfig(ctx context.Context, connectorType models.ConnectorType, config models.ConnectorConfig) *models.ConnectorTestResult {
	start := time.Now()

	switch connectorType {
	case models.ConnectorTypeSuperset:
		return s.testSupersetConnection(ctx, config, start)
	case models.ConnectorTypePostgres:
		return s.testPostgresConnection(ctx, config, start)
	default:
		return &models.ConnectorTestResult{
			Success: false,
			Message: "Unsupported connector type",
			Error:   fmt.Sprintf("connector type %s is not supported", connectorType),
		}
	}
}

// testSupersetConnection tests connection to Superset
func (s *ConnectorService) testSupersetConnection(ctx context.Context, config models.ConnectorConfig, start time.Time) *models.ConnectorTestResult {
	url, _ := config["url"].(string)
	username, _ := config["username"].(string)
	password, _ := config["password"].(string)

	supersetConn := connectors.NewSuperSetConnector(url, username, password, s.logger)

	// Test connection with a simple health check or login attempt
	err := supersetConn.TestConnection(ctx)
	responseTime := time.Since(start).Milliseconds()

	if err != nil {
		return &models.ConnectorTestResult{
			Success:      false,
			Message:      "Failed to connect to Superset",
			ResponseTime: &responseTime,
			Error:        err.Error(),
		}
	}

	return &models.ConnectorTestResult{
		Success:      true,
		Message:      "Successfully connected to Superset",
		ResponseTime: &responseTime,
	}
}

// testPostgresConnection tests connection to PostgreSQL
func (s *ConnectorService) testPostgresConnection(ctx context.Context, config models.ConnectorConfig, start time.Time) *models.ConnectorTestResult {
	url, _ := config["url"].(string)

	postgresConn, err := connectors.NewPostgresConnector(url, s.logger)
	if err != nil {
		responseTime := time.Since(start).Milliseconds()
		return &models.ConnectorTestResult{
			Success:      false,
			Message:      "Failed to create PostgreSQL connection",
			ResponseTime: &responseTime,
			Error:        err.Error(),
		}
	}
	defer postgresConn.Close()

	// Test with a simple query
	_, err = postgresConn.TestConnection(ctx)
	responseTime := time.Since(start).Milliseconds()

	if err != nil {
		return &models.ConnectorTestResult{
			Success:      false,
			Message:      "Failed to connect to PostgreSQL",
			ResponseTime: &responseTime,
			Error:        err.Error(),
		}
	}

	return &models.ConnectorTestResult{
		Success:      true,
		Message:      "Successfully connected to PostgreSQL",
		ResponseTime: &responseTime,
	}
}