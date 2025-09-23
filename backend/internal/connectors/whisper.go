package connectors

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"time"
)

type WhisperConnector struct {
	baseURL string
	client  *http.Client
	logger  *slog.Logger
}

type WhisperResponse struct {
	Text string `json:"text"`
}

func NewWhisperConnector(baseURL string, logger *slog.Logger) *WhisperConnector {
	return &WhisperConnector{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
		logger: logger.With("connector", "whisper"),
	}
}

func (wc *WhisperConnector) TranscribeAudio(ctx context.Context, audioData []byte, format string) (string, error) {
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	fileWriter, err := writer.CreateFormFile("audio_file", "audio."+format)
	if err != nil {
		return "", err
	}

	_, err = fileWriter.Write(audioData)
	if err != nil {
		return "", err
	}

	writer.WriteField("task", "transcribe")
	writer.WriteField("language", "en")
	writer.WriteField("output", "json")

	writer.Close()

	req, err := http.NewRequestWithContext(ctx, "POST",
		wc.baseURL+"/asr",
		&b)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := wc.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result WhisperResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w, body: %s", err, string(body))
	}

	return result.Text, nil
}

func (wc *WhisperConnector) HealthCheck(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", wc.baseURL+"/health", nil)
	if err != nil {
		return err
	}

	resp, err := wc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
