package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type ConnectorType string

const (
	ConnectorTypeSuperset ConnectorType = "superset"
	ConnectorTypePostgres ConnectorType = "postgres"
	ConnectorTypeMySQL    ConnectorType = "mysql"
	ConnectorTypeMongoDB  ConnectorType = "mongodb"
	ConnectorTypeAPI      ConnectorType = "api"
)

type ConnectorStatus string

const (
	ConnectorStatusConnected    ConnectorStatus = "connected"
	ConnectorStatusDisconnected ConnectorStatus = "disconnected"
	ConnectorStatusError        ConnectorStatus = "error"
	ConnectorStatusTesting      ConnectorStatus = "testing"
)

// ConnectorConfig represents the configuration for different connector types
type ConnectorConfig map[string]interface{}

// Value implements driver.Valuer for database storage
func (c ConnectorConfig) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// Scan implements sql.Scanner for database retrieval
func (c *ConnectorConfig) Scan(value interface{}) error {
	if value == nil {
		*c = make(ConnectorConfig)
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("cannot scan ConnectorConfig from type")
	}

	return json.Unmarshal(bytes, c)
}

// DataConnector represents a data source connector
type DataConnector struct {
	ID         string          `json:"id" db:"id"`
	Name       string          `json:"name" db:"name"`
	Type       ConnectorType   `json:"type" db:"type"`
	Status     ConnectorStatus `json:"status" db:"status"`
	Config     ConnectorConfig `json:"config" db:"config"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at" db:"updated_at"`
	LastTested *time.Time      `json:"last_tested,omitempty" db:"last_tested"`
}

// SupersetConfig represents Superset-specific configuration
type SupersetConfig struct {
	URL         string `json:"url"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	BearerToken string `json:"bearer_token,omitempty"`
}

// PostgresConfig represents PostgreSQL-specific configuration
type PostgresConfig struct {
	URL      string `json:"url"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Database string `json:"database,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
}

// ConnectorTestResult represents the result of testing a connector
type ConnectorTestResult struct {
	Success           bool     `json:"success"`
	Message           string   `json:"message"`
	ResponseTime      *int64   `json:"response_time,omitempty"` // in milliseconds
	AvailableDatasets []string `json:"available_datasets,omitempty"`
	Error             string   `json:"error,omitempty"`
}

// CreateConnectorRequest represents the request to create a new connector
type CreateConnectorRequest struct {
	Name   string          `json:"name" validate:"required,min=1,max=100"`
	Type   ConnectorType   `json:"type" validate:"required"`
	Config ConnectorConfig `json:"config" validate:"required"`
}

// UpdateConnectorRequest represents the request to update a connector
type UpdateConnectorRequest struct {
	Name   string          `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Config ConnectorConfig `json:"config,omitempty"`
	Status ConnectorStatus `json:"status,omitempty"`
}

// TestConnectorConfigRequest represents the request to test a connector configuration
type TestConnectorConfigRequest struct {
	Type   ConnectorType   `json:"type" validate:"required"`
	Config ConnectorConfig `json:"config" validate:"required"`
}