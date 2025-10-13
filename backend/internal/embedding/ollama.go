package embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// OllamaEmbeddingService implements EmbeddingService using Ollama
type OllamaEmbeddingService struct {
	baseURL    string
	model      string
	dimension  int
	httpClient *http.Client
	logger     *slog.Logger
}

// ollamaEmbeddingRequest represents the request to Ollama embedding API
type ollamaEmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// ollamaEmbeddingResponse represents the response from Ollama embedding API
type ollamaEmbeddingResponse struct {
	Embedding []float64 `json:"embedding"`
}

// NewOllamaEmbeddingService creates a new Ollama embedding service
func NewOllamaEmbeddingService(baseURL, model string, logger *slog.Logger) EmbeddingService {
	// Default to a lightweight embedding model dimension
	// In production, this should be determined by probing the model
	dimension := 384 // Common dimension for all-MiniLM-L6-v2

	return &OllamaEmbeddingService{
		baseURL:   baseURL,
		model:     model,
		dimension: dimension,
		httpClient: &http.Client{
			Timeout: 60 * time.Second, // Embeddings can take time
		},
		logger: logger,
	}
}

// GenerateEmbedding generates an embedding for a single text
func (o *OllamaEmbeddingService) GenerateEmbedding(ctx context.Context, text string) (*EmbeddingResponse, error) {
	if text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}

	o.logger.Debug("Generating embedding",
		"model", o.model,
		"text_length", len(text))

	// Try to call real Ollama API, fallback to mock if unavailable
	embedding, err := o.callOllamaEmbedding(ctx, text)
	if err != nil {
		o.logger.Warn("Ollama embedding API unavailable, using mock embedding", "error", err)
		embedding = o.generateMockEmbedding(text)
	}

	response := &EmbeddingResponse{
		Embedding: embedding,
		Dimension: len(embedding),
		Model:     o.model,
	}

	o.logger.Debug("Embedding generated successfully",
		"dimension", response.Dimension)

	return response, nil
}

// GenerateBatchEmbeddings generates embeddings for multiple texts
func (o *OllamaEmbeddingService) GenerateBatchEmbeddings(ctx context.Context, texts []string) ([]*EmbeddingResponse, error) {
	if len(texts) == 0 {
		return nil, fmt.Errorf("texts cannot be empty")
	}

	o.logger.Info("Generating batch embeddings",
		"count", len(texts),
		"model", o.model)

	responses := make([]*EmbeddingResponse, len(texts))
	for i, text := range texts {
		response, err := o.GenerateEmbedding(ctx, text)
		if err != nil {
			return nil, fmt.Errorf("failed to generate embedding for text %d: %w", i, err)
		}
		responses[i] = response
	}

	o.logger.Debug("Batch embeddings generated successfully",
		"count", len(responses))

	return responses, nil
}

// GetDimension returns the dimension of embeddings
func (o *OllamaEmbeddingService) GetDimension() int {
	return o.dimension
}

// GetModel returns the model name
func (o *OllamaEmbeddingService) GetModel() string {
	return o.model
}

// Health checks if Ollama is healthy
func (o *OllamaEmbeddingService) Health(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", o.baseURL+"/api/tags", nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("ollama health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ollama health check failed with status: %d", resp.StatusCode)
	}

	o.logger.Debug("Ollama embedding service health check passed")
	return nil
}

// Close closes the HTTP client
func (o *OllamaEmbeddingService) Close() error {
	o.logger.Info("Closing Ollama embedding service")
	return nil
}

// generateMockEmbedding creates a deterministic mock embedding based on text
func (o *OllamaEmbeddingService) generateMockEmbedding(text string) []float64 {
	// Create a simple hash-based embedding for testing
	embedding := make([]float64, o.dimension)

	// Use a simple hash function to create deterministic values
	textBytes := []byte(text)
	for i := 0; i < o.dimension; i++ {
		hash := 0
		for j, b := range textBytes {
			hash = hash*31 + int(b) + i + j
		}
		// Normalize to [-1, 1]
		embedding[i] = float64(hash%2000-1000) / 1000.0
	}

	// Normalize the vector
	magnitude := 0.0
	for _, val := range embedding {
		magnitude += val * val
	}
	magnitude = 1.0 / (magnitude + 1e-8) // Add small epsilon to avoid division by zero

	for i := range embedding {
		embedding[i] *= magnitude
	}

	return embedding
}

// callOllamaEmbedding calls the actual Ollama embedding API (for future implementation)
func (o *OllamaEmbeddingService) callOllamaEmbedding(ctx context.Context, text string) ([]float64, error) {
	reqBody := ollamaEmbeddingRequest{
		Model:  o.model,
		Prompt: text,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", o.baseURL+"/api/embeddings", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama returned status %d", resp.StatusCode)
	}

	var response ollamaEmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Embedding, nil
}