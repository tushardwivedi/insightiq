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
	// Use only 3 sample records to keep prompt short and fast
	sampleSize := min(3, len(data))
	sampleData := data[:sampleSize]

	dataJSON, _ := json.MarshalIndent(sampleData, "", "  ")

	prompt := fmt.Sprintf(`Analyze bike sales data (%d total records). Sample:
%s

Question: %s

Provide 2-3 key insights in 50 words or less.`, len(data), string(dataJSON), question)

	return oc.GenerateResponse(ctx, prompt)
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
