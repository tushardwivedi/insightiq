// internal/services/voice.go
package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"insightiq/backend/internal/agent"
)

type VoiceService struct {
	agentManager     *agent.Manager
	analyticsService *AnalyticsService
	logger           *slog.Logger
}

type VoiceRequest struct {
	AudioData []byte `json:"audio_data"`
	Format    string `json:"format"` // "wav", "mp3", "m4a", etc.
	Language  string `json:"language,omitempty"`
}

type VoiceResponse struct {
	Transcript  string            `json:"transcript"`
	Response    AnalyticsResponse `json:"response"`
	AudioReply  []byte            `json:"audio_reply,omitempty"`
	TaskID      string            `json:"task_id"`
	ProcessTime time.Duration     `json:"process_time"`
	Status      string            `json:"status"`
}

func NewVoiceService(agentManager *agent.Manager, logger *slog.Logger) *VoiceService {
	return &VoiceService{
		agentManager: agentManager,
		logger:       logger.With("service", "voice"),
	}
}

func (vs *VoiceService) ProcessVoiceQuery(ctx context.Context, audioData []byte, format string) (*VoiceResponse, error) {
	start := time.Now()
	taskID := generateTaskID()

	vs.logger.Info("Processing voice query",
		"task_id", taskID,
		"audio_size", len(audioData),
		"format", format)

	// Validate input
	if len(audioData) == 0 {
		return nil, fmt.Errorf("empty audio data")
	}

	if format == "" {
		format = "wav" // default
	}

	// Create task for voice agent
	task := agent.Task{
		ID:      taskID,
		Type:    "voice_query",
		AgentID: "voice-1",
		Payload: map[string]interface{}{
			"audio_data": audioData,
			"format":     format,
			"language":   "en", // default to English
		},
		Priority:  1,
		CreatedAt: time.Now(),
		Timeout:   120 * time.Second, // Voice processing can take longer
	}

	// Submit task
	if err := vs.agentManager.SubmitTask(task); err != nil {
		vs.logger.Error("Failed to submit voice task", "task_id", taskID, "error", err)
		return nil, fmt.Errorf("failed to submit voice task: %w", err)
	}

	// Wait for result
	result, err := vs.waitForVoiceResult(ctx, taskID, 120*time.Second)
	if err != nil {
		return nil, err
	}

	// Parse result
	response := &VoiceResponse{
		TaskID:      taskID,
		ProcessTime: time.Since(start),
		Status:      string(result.Status),
	}

	if result.Status == agent.TaskStatusCompleted {
		// Extract transcript
		if transcript, ok := result.Result["transcript"].(string); ok {
			response.Transcript = transcript
		}

		// Extract analytics response
		if analyticsResult, ok := result.Result["analytics"].(map[string]interface{}); ok {
			response.Response = AnalyticsResponse{
				Query:       response.Transcript,
				Insights:    getString(analyticsResult, "insights"),
				ProcessTime: response.ProcessTime,
				TaskID:      taskID,
				Timestamp:   time.Now(),
				Status:      "completed",
			}

			if data, ok := analyticsResult["data"].([]map[string]interface{}); ok {
				response.Response.Data = data
			}
		}

		// Extract audio reply if available
		if audioReply, ok := result.Result["audio_reply"].([]byte); ok {
			response.AudioReply = audioReply
		}
	} else {
		response.Status = "failed"
		return response, fmt.Errorf("voice task failed: %s", result.Error)
	}

	vs.logger.Info("Voice query processed successfully",
		"task_id", taskID,
		"transcript", response.Transcript,
		"duration", response.ProcessTime)

	return response, nil
}

func (vs *VoiceService) ProcessAudioFile(ctx context.Context, audioData []byte, filename string) (*VoiceResponse, error) {
	// Extract format from filename
	format := "wav"
	if len(filename) > 4 {
		ext := filename[len(filename)-3:]
		switch ext {
		case "mp3", "wav", "m4a", "ogg":
			format = ext
		}
	}

	return vs.ProcessVoiceQuery(ctx, audioData, format)
}

func (vs *VoiceService) GetSupportedFormats() []string {
	return []string{"wav", "mp3", "m4a", "ogg", "flac"}
}

func (vs *VoiceService) GetVoiceHistory(ctx context.Context, limit int) ([]VoiceResponse, error) {
	// In production, this would query a database
	// For MVP, return empty slice
	vs.logger.Info("Getting voice history", "limit", limit)
	return []VoiceResponse{}, nil
}

func (vs *VoiceService) HealthCheck(ctx context.Context) map[string]interface{} {
	status := map[string]interface{}{
		"voice_agent":       "healthy",
		"supported_formats": vs.GetSupportedFormats(),
		"timestamp":         time.Now(),
	}

	// Check agent manager status
	agentStatus := vs.agentManager.GetAgentStatus()
	if voiceAgentStatus, ok := agentStatus["voice-1"]; ok {
		status["voice_agent"] = voiceAgentStatus
	}

	return status
}

// Helper function to wait for voice task results
func (vs *VoiceService) waitForVoiceResult(ctx context.Context, taskID string, timeout time.Duration) (*agent.TaskResult, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(1 * time.Second) // Check every second for voice tasks
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("voice task timeout: %s", taskID)
		case <-ticker.C:
			// Mock implementation - in production, integrate with real agent system
			return &agent.TaskResult{
				TaskID:  taskID,
				AgentID: "voice-1",
				Status:  agent.TaskStatusCompleted,
				Result: map[string]interface{}{
					"transcript": "Show me the sales data for the last quarter",
					"analytics": map[string]interface{}{
						"data": []map[string]interface{}{
							{"quarter": "Q1 2024", "revenue": 150000, "orders": 450},
							{"quarter": "Q2 2024", "revenue": 180000, "orders": 520},
							{"quarter": "Q3 2024", "revenue": 165000, "orders": 485},
						},
						"insights": "Q3 sales show a slight decline from Q2 peak but remain above Q1 levels. Revenue per order has remained stable around $340.",
					},
				},
				ProcessedAt: time.Now(),
				Duration:    5 * time.Second,
			}, nil
		}
	}
}

// Helper function to safely extract string from map
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}
