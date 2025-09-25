package agent

import (
	"context"
	"fmt"
	"log/slog"

	"insightiq/backend/internal/connectors"
)

type VoiceAgent struct {
	*BaseAgent
	whisperConn    *connectors.WhisperConnector
	analyticsAgent *AnalyticsAgent
}

func NewVoiceAgent(id string, whisper *connectors.WhisperConnector, analyticsAgent *AnalyticsAgent, logger *slog.Logger) *VoiceAgent {
	return &VoiceAgent{
		BaseAgent:      NewBaseAgent(id, AgentTypeVoice, logger),
		whisperConn:    whisper,
		analyticsAgent: analyticsAgent,
	}
}

func (va *VoiceAgent) ProcessTask(ctx context.Context, task Task) (*TaskResult, error) {
	switch task.Type {
	case "voice_query":
		return va.processVoiceQuery(ctx, task)
	default:
		return nil, fmt.Errorf("unsupported task type: %s", task.Type)
	}
}

func (va *VoiceAgent) processVoiceQuery(ctx context.Context, task Task) (*TaskResult, error) {
	// Extract audio data
	audioData, ok := task.Payload["audio_data"].([]byte)
	if !ok {
		return nil, fmt.Errorf("missing or invalid audio_data")
	}

	format, ok := task.Payload["format"].(string)
	if !ok {
		format = "wav"
	}

	va.logger.Info("Processing voice query",
		"task_id", task.ID,
		"audio_size", len(audioData),
		"format", format)

	// Step 1: Transcribe audio to text
	transcript, err := va.whisperConn.TranscribeAudio(ctx, audioData, format)
	if err != nil {
		return nil, fmt.Errorf("failed to transcribe audio: %w", err)
	}

	va.logger.Info("Audio transcribed", "task_id", task.ID, "transcript", transcript)

	// Step 2: Process text query through analytics agent
	analyticsTask := Task{
		ID:      task.ID + "_analytics",
		Type:    "text_query",
		AgentID: va.analyticsAgent.ID(),
		Payload: map[string]interface{}{
			"query": transcript,
		},
		Priority:  task.Priority,
		CreatedAt: task.CreatedAt,
		Timeout:   task.Timeout,
	}

	analyticsResult, err := va.analyticsAgent.ProcessTask(ctx, analyticsTask)
	if err != nil {
		return nil, fmt.Errorf("failed to process analytics: %w", err)
	}

	// Step 3: Combine results
	result := &TaskResult{
		TaskID:  task.ID,
		AgentID: va.ID(),
		Status:  TaskStatusCompleted,
		Result: map[string]interface{}{
			"transcript": transcript,
			"analytics":  analyticsResult.Result,
		},
		ProcessedAt: analyticsResult.ProcessedAt,
		Duration:    analyticsResult.Duration,
	}

	va.logger.Info("Voice query processed successfully", "task_id", task.ID)

	return result, nil
}
