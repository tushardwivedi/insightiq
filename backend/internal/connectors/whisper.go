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
			Timeout: 120 * time.Second, // Increase timeout for longer audio files
		},
		logger: logger.With("connector", "whisper"),
	}
}

func (wc *WhisperConnector) TranscribeAudio(ctx context.Context, audioData []byte, format string) (string, error) {
	wc.logger.Info("Starting audio transcription",
		"audio_size", len(audioData),
		"format", format,
		"whisper_url", wc.baseURL+"/asr")

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	fileWriter, err := writer.CreateFormFile("audio_file", "audio."+format)
	if err != nil {
		wc.logger.Error("Failed to create form file", "error", err)
		return "", err
	}

	_, err = fileWriter.Write(audioData)
	if err != nil {
		wc.logger.Error("Failed to write audio data", "error", err)
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
		wc.logger.Error("Failed to send request to Whisper", "error", err)
		return "", err
	}
	defer resp.Body.Close()

	wc.logger.Info("Received response from Whisper", "status_code", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		wc.logger.Error("Failed to read response body", "error", err)
		return "", err
	}

	if resp.StatusCode != 200 {
		wc.logger.Error("Whisper returned error", "status_code", resp.StatusCode, "body", string(body))
		return "", fmt.Errorf("whisper returned status %d: %s", resp.StatusCode, string(body))
	}

	wc.logger.Info("Whisper response", "body_length", len(body), "body", string(body))

	var result WhisperResponse
	if err := json.Unmarshal(body, &result); err != nil {
		wc.logger.Error("Failed to parse Whisper response", "error", err, "body", string(body))
		return "", fmt.Errorf("failed to parse response: %w, body: %s", err, string(body))
	}

	wc.logger.Info("Audio transcription completed", "transcript_length", len(result.Text), "transcript", result.Text)
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
