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
	postgresConn *connectors.PostgresConnector
	llmConn      *connectors.OllamaConnector
}

func NewAnalyticsAgent(id string, superset *connectors.SuperSetConnector, postgres *connectors.PostgresConnector, llm *connectors.OllamaConnector, logger *slog.Logger) *AnalyticsAgent {
	return &AnalyticsAgent{
		BaseAgent:    NewBaseAgent(id, AgentTypeAnalytics, logger),
		supersetConn: superset,
		postgresConn: postgres,
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

	// Try to get data from PostgreSQL first, fallback to SuperSet
	var data []map[string]interface{}
	var err error

	if aa.postgresConn != nil {
		result, pgErr := aa.postgresConn.GetBikeSalesData(ctx)
		if pgErr != nil {
			aa.logger.Warn("PostgreSQL query failed, trying SuperSet", "error", pgErr)
		} else {
			data = result.Data
			aa.logger.Info("Got data from PostgreSQL", "rows", len(data))
		}
	}

	// Fallback to SuperSet if PostgreSQL failed
	if len(data) == 0 && aa.supersetConn != nil {
		supersetData, err := aa.supersetConn.GetSampleData(ctx)
		if err != nil {
			aa.logger.Error("Both PostgreSQL and SuperSet failed", "error", err)
			return nil, fmt.Errorf("failed to get data from any source: %w", err)
		}
		data = supersetData.Data
		aa.logger.Info("Got data from SuperSet", "rows", len(data))
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("no data available from any source")
	}

	// Generate insights with LLM
	insights, err := aa.llmConn.AnalyzeData(ctx, data, query)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze data: %w", err)
	}

	return &TaskResult{
		TaskID:  task.ID,
		AgentID: aa.ID(),
		Status:  TaskStatusCompleted,
		Result: map[string]interface{}{
			"query":     query,
			"data":      data,
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
