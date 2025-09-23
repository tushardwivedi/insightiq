// internal/agent/types.go
package agent

import (
	"context"
	"time"
)

type Agent interface {
	ID() string
	Type() AgentType
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	ProcessTask(ctx context.Context, task Task) (*TaskResult, error)
	Health() HealthStatus
}

type AgentType string

const (
	AgentTypeAnalytics AgentType = "analytics"
	AgentTypeVoice     AgentType = "voice"
	AgentTypeData      AgentType = "data"
)

type Task struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	AgentID   string                 `json:"agent_id"`
	Payload   map[string]interface{} `json:"payload"`
	Priority  int                    `json:"priority"`
	CreatedAt time.Time              `json:"created_at"`
	Timeout   time.Duration          `json:"timeout"`
}

type TaskResult struct {
	TaskID      string                 `json:"task_id"`
	AgentID     string                 `json:"agent_id"`
	Status      TaskStatus             `json:"status"`
	Result      map[string]interface{} `json:"result,omitempty"`
	Error       string                 `json:"error,omitempty"`
	ProcessedAt time.Time              `json:"processed_at"`
	Duration    time.Duration          `json:"duration"`
}

type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
)

type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusDegraded  HealthStatus = "degraded"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
)
