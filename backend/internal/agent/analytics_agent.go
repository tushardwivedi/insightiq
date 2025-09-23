// internal/agent/analytics_agent.go
package agent

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"insightiq/backend/internal/connectors"
)

type AnalyticsAgent struct {
	*BaseAgent
	supersetConn *connectors.SuperSetConnector
	llmConn      *connectors.OllamaConnector
}

func NewAnalyticsAgent(id string, superset *connectors.SuperSetConnector, llm *connectors.OllamaConnector, logger *slog.Logger) *AnalyticsAgent {
	return &AnalyticsAgent{
		BaseAgent:    NewBaseAgent(id, AgentTypeAnalytics, logger),
		supersetConn: superset,
		llmConn:      llm,
	}
}

func (aa *AnalyticsAgent) ProcessTask(ctx context.Context, task Task) (*TaskResult, error) {
	switch task.Type {
	case "text_query":
		return aa.processTextQuery(ctx, task)
	case "sql_query":
		return aa.processSQLQuery(ctx, task)
	default:
		return nil, fmt.Errorf("unsupported task type: %s", task.Type)
	}
}

func (aa *AnalyticsAgent) processTextQuery(ctx context.Context, task Task) (*TaskResult, error) {
	query, ok := task.Payload["query"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid query parameter")
	}

	aa.logger.Info("Processing text query", "query", query)

	// Get data from SuperSet
	data, err := aa.supersetConn.GetSampleData(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get data: %w", err)
	}

	// Generate insights with LLM
	insights, err := aa.llmConn.AnalyzeData(ctx, data.Data, query)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze data: %w", err)
	}

	return &TaskResult{
		TaskID:  task.ID,
		AgentID: aa.ID(),
		Status:  TaskStatusCompleted,
		Result: map[string]interface{}{
			"query":     query,
			"data":      data.Data,
			"insights":  insights,
			"timestamp": time.Now(),
		},
	}, nil
}

func (aa *AnalyticsAgent) processSQLQuery(ctx context.Context, task Task) (*TaskResult, error) {
	sql, ok := task.Payload["sql"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid sql parameter")
	}

	question, ok := task.Payload["question"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid question parameter")
	}

	aa.logger.Info("Processing SQL query", "sql", sql)

	// Execute SQL
	data, err := aa.supersetConn.ExecuteSQL(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SQL: %w", err)
	}

	// Generate insights
	insights, err := aa.llmConn.AnalyzeData(ctx, data.Data, question)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze data: %w", err)
	}

	return &TaskResult{
		TaskID:  task.ID,
		AgentID: aa.ID(),
		Status:  TaskStatusCompleted,
		Result: map[string]interface{}{
			"query":     question,
			"data":      data.Data,
			"insights":  insights,
			"timestamp": time.Now(),
		},
	}, nil
}
