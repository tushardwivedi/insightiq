// internal/agent/manager.go
package agent

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type Manager struct {
	agents      map[string]Agent
	taskQueue   chan Task
	resultQueue chan TaskResult
	results     map[string]*TaskResult
	logger      *slog.Logger
	mu          sync.RWMutex

	// Metrics
	tasksProcessed int64
	tasksInFlight  int64
}

func NewManager(logger *slog.Logger) *Manager {
	return &Manager{
		agents:      make(map[string]Agent),
		taskQueue:   make(chan Task, 1000),
		resultQueue: make(chan TaskResult, 1000),
		results:     make(map[string]*TaskResult),
		logger:      logger.With("component", "agent_manager"),
	}
}

func (m *Manager) RegisterAgent(agent Agent) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.agents[agent.ID()]; exists {
		return fmt.Errorf("agent %s already registered", agent.ID())
	}

	m.agents[agent.ID()] = agent
	m.logger.Info("Agent registered", "agent_id", agent.ID(), "agent_type", agent.Type())
	return nil
}

func (m *Manager) Start(ctx context.Context) error {
	m.logger.Info("Starting agent manager")

	// Start all agents
	for id, agent := range m.agents {
		if err := agent.Start(ctx); err != nil {
			m.logger.Error("Failed to start agent", "agent_id", id, "error", err)
			return err
		}
	}

	// Start task dispatcher
	go m.taskDispatcher(ctx)
	go m.resultCollector(ctx)
	go m.healthMonitor(ctx)

	return nil
}

func (m *Manager) SubmitTask(task Task) error {
	// Set default timeout if not specified
	if task.Timeout == 0 {
		task.Timeout = 30 * time.Second
	}

	// Set creation time
	task.CreatedAt = time.Now()

	select {
	case m.taskQueue <- task:
		m.logger.Info("Task submitted", "task_id", task.ID, "agent_id", task.AgentID)
		return nil
	default:
		return fmt.Errorf("task queue full")
	}
}

func (m *Manager) taskDispatcher(ctx context.Context) {
	for {
		select {
		case task := <-m.taskQueue:
			m.dispatchTask(task)
		case <-ctx.Done():
			m.logger.Info("Task dispatcher stopping")
			return
		}
	}
}

func (m *Manager) dispatchTask(task Task) {
	m.mu.RLock()
	agent, exists := m.agents[task.AgentID]
	m.mu.RUnlock()

	if !exists {
		m.logger.Error("Agent not found", "agent_id", task.AgentID, "task_id", task.ID)
		return
	}

	// Create task context with sufficient timeout for LLM processing
	ctx, cancel := context.WithTimeout(context.Background(), task.Timeout)
	defer cancel()

	go func() {
		result, err := agent.ProcessTask(ctx, task)
		if err != nil {
			m.logger.Error("Task processing failed",
				"task_id", task.ID,
				"agent_id", task.AgentID,
				"error", err)

			// Create a failed task result
			failedResult := TaskResult{
				TaskID:      task.ID,
				AgentID:     task.AgentID,
				Status:      TaskStatusFailed,
				Error:       err.Error(),
				ProcessedAt: time.Now(),
			}

			select {
			case m.resultQueue <- failedResult:
			case <-ctx.Done():
				m.logger.Warn("Failed result queue timeout", "task_id", task.ID)
			}
		} else if result != nil {
			select {
			case m.resultQueue <- *result:
			case <-ctx.Done():
				m.logger.Warn("Result queue timeout", "task_id", task.ID)
			}
		}
	}()
}

func (m *Manager) resultCollector(ctx context.Context) {
	for {
		select {
		case result := <-m.resultQueue:
			m.logger.Info("Task result collected",
				"task_id", result.TaskID,
				"agent_id", result.AgentID,
				"status", result.Status,
				"duration", result.Duration)

			// Store the result for retrieval
			m.mu.Lock()
			m.results[result.TaskID] = &result
			m.mu.Unlock()

		case <-ctx.Done():
			m.logger.Info("Result collector stopping")
			return
		}
	}
}

func (m *Manager) healthMonitor(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.checkAgentHealth()
		case <-ctx.Done():
			m.logger.Info("Health monitor stopping")
			return
		}
	}
}

func (m *Manager) checkAgentHealth() {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for id, agent := range m.agents {
		health := agent.Health()
		if health != HealthStatusHealthy {
			m.logger.Warn("Agent unhealthy", "agent_id", id, "status", health)
		}
	}
}

func (m *Manager) GetAgentStatus() map[string]HealthStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	status := make(map[string]HealthStatus)
	for id, agent := range m.agents {
		status[id] = agent.Health()
	}

	return status
}

func (m *Manager) GetTaskResult(taskID string) *TaskResult {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if result, exists := m.results[taskID]; exists {
		return result
	}
	return nil
}
