package schema

import (
	"context"
)

// ConnectorInfo represents basic connector information
type ConnectorInfo struct {
	ID     string
	Name   string
	Type   string
	Status string
	Config map[string]interface{}
}

// ConnectorService interface to avoid circular imports
type ConnectorService interface {
	ListConnectors(ctx context.Context) ([]ConnectorInfo, error)
	GetConnector(ctx context.Context, id string) (*ConnectorInfo, error)
}