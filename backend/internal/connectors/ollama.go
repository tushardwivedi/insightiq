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

	jsonData, _ := json.Marshal(request)

	req, err := http.NewRequestWithContext(ctx, "POST",
		oc.baseURL+"/api/generate",
		bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := oc.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Response, nil
}

func (oc *OllamaConnector) AnalyzeData(ctx context.Context, data []map[string]interface{}, question string) (string, error) {
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

	return oc.GenerateResponse(ctx, prompt)
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
		return err
	}

	resp, err := oc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
