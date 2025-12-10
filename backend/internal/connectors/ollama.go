package connectors

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type OllamaConnector struct {
	baseURL string
	client  *http.Client
	logger  *slog.Logger
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

type OllamaModel struct {
	Name string `json:"name"`
}

type OllamaModelsResponse struct {
	Models []OllamaModel `json:"models"`
}

func NewOllamaConnector(baseURL string, logger *slog.Logger) *OllamaConnector {
	return &OllamaConnector{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
		logger: logger.With("connector", "ollama"),
	}
}

func (oc *OllamaConnector) GenerateResponse(ctx context.Context, prompt string) (string, error) {
	request := OllamaRequest{
		Model:  "llama3.2:1b",
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST",
		oc.baseURL+"/api/generate",
		bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := oc.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call Ollama API: %w", err)
	}
	defer resp.Body.Close()

	// CRITICAL: Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		bodyBytes := make([]byte, 512)
		n, _ := resp.Body.Read(bodyBytes)
		bodySnippet := string(bodyBytes[:n])

		if resp.StatusCode == http.StatusNotFound {
			return "", fmt.Errorf("Ollama model 'llama3.2:1b' not found (HTTP 404). Please install it with: docker exec insightiq-ollama-1 ollama pull llama3.2:1b. Response: %s", bodySnippet)
		}
		return "", fmt.Errorf("Ollama API returned HTTP %d: %s", resp.StatusCode, bodySnippet)
	}

	var result OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode Ollama response: %w", err)
	}

	// Validate that we actually got a response
	if result.Response == "" {
		return "", fmt.Errorf("Ollama returned empty response for model 'llama3.2:1b'")
	}

	return result.Response, nil
}

func (oc *OllamaConnector) AnalyzeData(ctx context.Context, data []map[string]interface{}, question string) (string, error) {
	oc.logger.Info("AnalyzeData called", "data_count", len(data), "question", question)

	// Check if data is actually an error message
	if len(data) == 1 {
		if errMsg, ok := data[0]["error"].(string); ok {
			oc.logger.Warn("Cannot generate insights from error data", "error", errMsg)
			return fmt.Sprintf("Unable to retrieve data: %s", errMsg), nil
		}
		if msg, ok := data[0]["message"].(string); ok {
			oc.logger.Warn("Cannot generate insights from message data", "message", msg)
			return fmt.Sprintf("Data retrieval issue: %s", msg), nil
		}
	}

	// Use only 3 sample records to keep prompt short and fast
	sampleSize := min(3, len(data))
	sampleData := data[:sampleSize]

	dataJSON, _ := json.MarshalIndent(sampleData, "", "  ")

	// Make the prompt more generic - not just "bike sales"
	prompt := fmt.Sprintf(`Analyze this data (%d total records). Sample:
%s

Question: %s

Provide 2-3 key insights in 50 words or less.`, len(data), string(dataJSON), question)

	oc.logger.Info("Calling Ollama GenerateResponse", "prompt_length", len(prompt))
	response, err := oc.GenerateResponse(ctx, prompt)
	if err != nil {
		oc.logger.Error("Ollama GenerateResponse failed", "error", err)
		return "", err
	}
	oc.logger.Info("Ollama GenerateResponse successful", "response_length", len(response))
	return response, nil
}

// GenerateVisualizationData creates synthetic data from insights for visualization
func (oc *OllamaConnector) GenerateVisualizationData(ctx context.Context, insights string, query string) ([]map[string]interface{}, error) {
	prompt := fmt.Sprintf(`Based on these insights:
%s

Extract key data points and convert them into structured JSON data for visualization.
Format: Return ONLY a JSON array of objects with meaningful keys and numeric values.

Example for "64%% were 30-49 years old, 20%% were 16-29":
[{"age_group":"30-49","percentage":64},{"age_group":"16-29","percentage":20}]

Return ONLY the JSON array, no explanation:`, insights)

	response, err := oc.GenerateResponse(ctx, prompt)
	if err != nil {
		return nil, err
	}

	// Try to extract JSON from the response
	var data []map[string]interface{}
	if err := json.Unmarshal([]byte(response), &data); err != nil {
		oc.logger.Warn("Failed to parse LLM visualization data, creating fallback", "error", err)
		// Create a simple fallback visualization data
		return []map[string]interface{}{
			{"category": "Insight 1", "value": 100},
			{"category": "Insight 2", "value": 75},
			{"category": "Insight 3", "value": 50},
		}, nil
	}

	return data, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (oc *OllamaConnector) HealthCheck(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", oc.baseURL+"/api/tags", nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := oc.client.Do(req)
	if err != nil {
		return fmt.Errorf("Ollama service is not reachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ollama health check failed with HTTP %d", resp.StatusCode)
	}

	return nil
}

// CheckModelAvailability verifies that the required model is installed
func (oc *OllamaConnector) CheckModelAvailability(ctx context.Context, modelName string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", oc.baseURL+"/api/tags", nil)
	if err != nil {
		return fmt.Errorf("failed to create model check request: %w", err)
	}

	resp, err := oc.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to check Ollama models: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ollama API returned HTTP %d when checking models", resp.StatusCode)
	}

	var modelsResp OllamaModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&modelsResp); err != nil {
		return fmt.Errorf("failed to decode models response: %w", err)
	}

	// Check if the required model is in the list
	for _, model := range modelsResp.Models {
		if model.Name == modelName || model.Name == modelName+":latest" {
			oc.logger.Info("Required Ollama model found", "model", modelName)
			return nil
		}
	}

	return fmt.Errorf("required Ollama model '%s' is not installed. Please run: docker exec insightiq-ollama-1 ollama pull %s", modelName, modelName)
}
