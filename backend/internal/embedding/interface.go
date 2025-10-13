package embedding

import (
	"context"
)

// EmbeddingRequest represents a request to generate embeddings
type EmbeddingRequest struct {
	Text  string `json:"text"`
	Model string `json:"model,omitempty"`
}

// EmbeddingResponse represents the response from embedding generation
type EmbeddingResponse struct {
	Embedding []float64 `json:"embedding"`
	Dimension int       `json:"dimension"`
	Model     string    `json:"model"`
}

// EmbeddingService defines the interface for text embedding generation
type EmbeddingService interface {
	// GenerateEmbedding generates an embedding for a single text
	GenerateEmbedding(ctx context.Context, text string) (*EmbeddingResponse, error)

	// GenerateBatchEmbeddings generates embeddings for multiple texts
	GenerateBatchEmbeddings(ctx context.Context, texts []string) ([]*EmbeddingResponse, error)

	// GetDimension returns the dimension of embeddings produced by this service
	GetDimension() int

	// GetModel returns the model name used by this service
	GetModel() string

	// Health checks if the embedding service is healthy
	Health(ctx context.Context) error

	// Close closes any open connections
	Close() error
}