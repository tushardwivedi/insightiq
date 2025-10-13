package vectorstore

import (
	"context"
)

// Vector represents a vector with metadata
type Vector struct {
	ID       string                 `json:"id"`
	Values   []float64              `json:"values"`
	Metadata map[string]interface{} `json:"metadata"`
}

// SearchResult represents a search result with similarity score
type SearchResult struct {
	Vector Vector  `json:"vector"`
	Score  float64 `json:"score"`
}

// VectorStore defines the interface for vector database operations
type VectorStore interface {
	// UpsertVector inserts or updates a vector
	UpsertVector(ctx context.Context, collectionName string, vector Vector) error

	// UpsertVectors inserts or updates multiple vectors
	UpsertVectors(ctx context.Context, collectionName string, vectors []Vector) error

	// SearchVectors performs similarity search
	SearchVectors(ctx context.Context, collectionName string, queryVector []float64, limit int) ([]SearchResult, error)

	// DeleteVector deletes a vector by ID
	DeleteVector(ctx context.Context, collectionName string, vectorID string) error

	// CreateCollection creates a new collection with specified dimension
	CreateCollection(ctx context.Context, collectionName string, dimension int) error

	// DeleteCollection deletes a collection
	DeleteCollection(ctx context.Context, collectionName string) error

	// ListCollections returns all collection names
	ListCollections(ctx context.Context) ([]string, error)

	// Close closes the connection
	Close() error
}