package models

import (
	"encoding/json"
	"time"
)

// Intent represents the parsed intent from user input
type Intent struct {
	Type        IntentType             `json:"type"`
	Confidence  float64                `json:"confidence"`
	Entities    map[string]interface{} `json:"entities"`
	Parameters  map[string]string      `json:"parameters"`
	Subintents  []Intent               `json:"subintents,omitempty"`
	ParsedQuery ParsedQuery            `json:"parsed_query"`
}

// IntentType represents different types of user intents
type IntentType string

const (
	IntentTypeAnalytics     IntentType = "analytics"
	IntentTypeSQL           IntentType = "sql"
	IntentTypeVisualization IntentType = "visualization"
	IntentTypeComparison    IntentType = "comparison"
	IntentTypeTrend         IntentType = "trend"
	IntentTypeFilter        IntentType = "filter"
	IntentTypeAggregation   IntentType = "aggregation"
	IntentTypeJoin          IntentType = "join"
	IntentTypeUnknown       IntentType = "unknown"
)

// ParsedQuery contains structured information extracted from the query
type ParsedQuery struct {
	MainAction    string                 `json:"main_action"`
	DataSources   []string               `json:"data_sources"`
	Metrics       []string               `json:"metrics"`
	Dimensions    []string               `json:"dimensions"`
	TimeRange     *TimeRange             `json:"time_range,omitempty"`
	Filters       []Filter               `json:"filters"`
	Aggregations  []Aggregation          `json:"aggregations"`
	SortBy        []SortCriteria         `json:"sort_by"`
	Limit         *int                   `json:"limit,omitempty"`
	GroupBy       []string               `json:"group_by"`
	Joins         []Join                 `json:"joins"`
	OutputFormat  string                 `json:"output_format"`
	Metadata      map[string]interface{} `json:"metadata"`
}

// TimeRange represents time-based filters
type TimeRange struct {
	Start    *time.Time `json:"start,omitempty"`
	End      *time.Time `json:"end,omitempty"`
	Period   string     `json:"period"`    // "last_week", "last_month", "last_quarter", "last_year"
	Relative string     `json:"relative"`  // "today", "yesterday", "last_7_days"
}

// Filter represents a data filter condition
type Filter struct {
	Field     string      `json:"field"`
	Operator  string      `json:"operator"` // "=", "!=", ">", "<", ">=", "<=", "IN", "LIKE"
	Value     interface{} `json:"value"`
	Condition string      `json:"condition"` // "AND", "OR"
}

// Aggregation represents data aggregation operations
type Aggregation struct {
	Function string `json:"function"` // "SUM", "COUNT", "AVG", "MAX", "MIN"
	Field    string `json:"field"`
	Alias    string `json:"alias,omitempty"`
}

// SortCriteria represents sorting requirements
type SortCriteria struct {
	Field     string `json:"field"`
	Direction string `json:"direction"` // "ASC", "DESC"
}

// Join represents table join operations
type Join struct {
	Type        string `json:"type"`         // "INNER", "LEFT", "RIGHT", "FULL"
	Table       string `json:"table"`
	OnCondition string `json:"on_condition"`
}

// TaskGraph represents the planned execution steps
type TaskGraph struct {
	ID          string      `json:"id"`
	Query       string      `json:"query"`
	Intent      Intent      `json:"intent"`
	Steps       []TaskStep  `json:"steps"`
	Dependencies map[string][]string `json:"dependencies"`
	EstimatedTime time.Duration `json:"estimated_time"`
	CreatedAt   time.Time   `json:"created_at"`
	Status      TaskStatus  `json:"status"`
}

// TaskStep represents an individual step in the execution plan
type TaskStep struct {
	ID          string                 `json:"id"`
	Type        TaskStepType           `json:"type"`
	Description string                 `json:"description"`
	Action      string                 `json:"action"`
	Parameters  map[string]interface{} `json:"parameters"`
	DataSources []string               `json:"data_sources"`
	Dependencies []string              `json:"dependencies"`
	EstimatedTime time.Duration        `json:"estimated_time"`
	Priority    int                    `json:"priority"`
	Retry       *RetryConfig           `json:"retry,omitempty"`
}

// TaskStepType represents different types of execution steps
type TaskStepType string

const (
	TaskStepTypeDataRetrieval TaskStepType = "data_retrieval"
	TaskStepTypeTransformation TaskStepType = "transformation"
	TaskStepTypeAggregation   TaskStepType = "aggregation"
	TaskStepTypeJoin          TaskStepType = "join"
	TaskStepTypeFilter        TaskStepType = "filter"
	TaskStepTypeSort          TaskStepType = "sort"
	TaskStepTypeVisualization TaskStepType = "visualization"
	TaskStepTypeAnalysis      TaskStepType = "analysis"
	TaskStepTypeValidation    TaskStepType = "validation"
)

// TaskStatus represents the status of task execution
type TaskStatus string

const (
	TaskStatusPlanned    TaskStatus = "planned"
	TaskStatusExecuting  TaskStatus = "executing"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusFailed     TaskStatus = "failed"
	TaskStatusCancelled  TaskStatus = "cancelled"
)

// RetryConfig represents retry configuration for tasks
type RetryConfig struct {
	MaxRetries int           `json:"max_retries"`
	Delay      time.Duration `json:"delay"`
	Backoff    string        `json:"backoff"` // "linear", "exponential"
}

// PlannerRequest represents a request to the planner service
type PlannerRequest struct {
	Query       string            `json:"query"`
	Context     map[string]interface{} `json:"context,omitempty"`
	UserID      string            `json:"user_id,omitempty"`
	SessionID   string            `json:"session_id,omitempty"`
	Preferences map[string]string `json:"preferences,omitempty"`
}

// PlannerResponse represents the response from the planner service
type PlannerResponse struct {
	Intent       Intent                 `json:"intent"`
	TaskGraph    TaskGraph              `json:"task_graph"`
	Confidence   float64                `json:"confidence"`
	Alternatives []TaskGraph            `json:"alternatives,omitempty"`
	Metadata     map[string]interface{} `json:"metadata"`
	ProcessTime  time.Duration          `json:"process_time"`
	CreatedAt    time.Time              `json:"created_at"`
}

// String returns a string representation of IntentType
func (it IntentType) String() string {
	return string(it)
}

// MarshalJSON implements json.Marshaler for TaskGraph
func (tg *TaskGraph) MarshalJSON() ([]byte, error) {
	type Alias TaskGraph
	return json.Marshal(&struct {
		EstimatedTime string `json:"estimated_time"`
		*Alias
	}{
		EstimatedTime: tg.EstimatedTime.String(),
		Alias:         (*Alias)(tg),
	})
}

// UnmarshalJSON implements json.Unmarshaler for TaskGraph
func (tg *TaskGraph) UnmarshalJSON(data []byte) error {
	type Alias TaskGraph
	aux := &struct {
		EstimatedTime string `json:"estimated_time"`
		*Alias
	}{
		Alias: (*Alias)(tg),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.EstimatedTime != "" {
		duration, err := time.ParseDuration(aux.EstimatedTime)
		if err != nil {
			return err
		}
		tg.EstimatedTime = duration
	}

	return nil
}