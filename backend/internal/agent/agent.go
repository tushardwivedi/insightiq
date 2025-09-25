// internal/agent/agent.go
package agent

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

type BaseAgent struct {
	id        string
	agentType AgentType
	status    HealthStatus
	logger    *slog.Logger

	// Task processing
	taskQueue   chan Task
	resultQueue chan TaskResult

	// Lifecycle
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	mu     sync.RWMutex
}

func NewBaseAgent(id string, agentType AgentType, logger *slog.Logger) *BaseAgent {
	return &BaseAgent{
		id:          id,
		agentType:   agentType,
		status:      HealthStatusHealthy,
		logger:      logger.With("agent_id", id, "agent_type", agentType),
		taskQueue:   make(chan Task, 100),
		resultQueue: make(chan TaskResult, 100),
	}
}

func (ba *BaseAgent) ID() string {
	return ba.id
}

func (ba *BaseAgent) Type() AgentType {
	return ba.agentType
}

func (ba *BaseAgent) Health() HealthStatus {
	ba.mu.RLock()
	defer ba.mu.RUnlock()
	return ba.status
}

func (ba *BaseAgent) Start(ctx context.Context) error {
	ba.ctx, ba.cancel = context.WithCancel(ctx)

	ba.wg.Add(1)
	go ba.taskProcessor()

	ba.logger.Info("Agent started")
	return nil
}

func (ba *BaseAgent) Stop(ctx context.Context) error {
	if ba.cancel != nil {
		ba.cancel()
	}

	// Wait for graceful shutdown with timeout
	done := make(chan struct{})
	go func() {
		ba.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		ba.logger.Info("Agent stopped gracefully")
	case <-ctx.Done():
		ba.logger.Warn("Agent stop timeout")
	}

	return nil
}

func (ba *BaseAgent) taskProcessor() {
	defer ba.wg.Done()

	for {
		select {
		case task := <-ba.taskQueue:
			ba.processTask(task)
		case <-ba.ctx.Done():
			ba.logger.Info("Task processor stopping")
			return
		}
	}
}

func (ba *BaseAgent) processTask(task Task) {
	start := time.Now()
	ba.logger.Info("Processing task", "task_id", task.ID, "task_type", task.Type)

	result := TaskResult{
		TaskID:      task.ID,
		AgentID:     ba.id,
		ProcessedAt: time.Now(),
	}

	// Create task-specific context with timeout (independent of agent context)
	taskCtx, cancel := context.WithTimeout(context.Background(), task.Timeout)
	defer cancel()

	// Process task (override in specific agents)
	taskResult, err := ba.ProcessTask(taskCtx, task)
	if err != nil {
		result.Status = TaskStatusFailed
		result.Error = err.Error()
		ba.logger.Error("Task failed", "task_id", task.ID, "error", err)
	} else {
		result.Status = TaskStatusCompleted
		if taskResult != nil {
			result.Result = taskResult.Result
		}
		ba.logger.Info("Task completed", "task_id", task.ID)
	}

	result.Duration = time.Since(start)

	// Send result
	select {
	case ba.resultQueue <- result:
	case <-taskCtx.Done():
		ba.logger.Warn("Result queue full, dropping result", "task_id", task.ID)
	}
}

// Default implementation - override in specific agents
func (ba *BaseAgent) ProcessTask(ctx context.Context, task Task) (*TaskResult, error) {
	return &TaskResult{
		TaskID:  task.ID,
		AgentID: ba.id,
		Status:  TaskStatusCompleted,
		Result: map[string]interface{}{
			"message": "Default task processing",
		},
	}, nil
}
