package connectors

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type SuperSetConnector struct {
	baseURL  string
	username string
	password string
	token    string
	client   *http.Client
	logger   *slog.Logger
}

type SuperSetQuery struct {
	SQL        string `json:"sql"`
	DatabaseID int    `json:"database_id"`
}

type SuperSetResponse struct {
	Data   []map[string]interface{} `json:"data"`
	Status string                   `json:"status"`
	Error  string                   `json:"error,omitempty"`
}

func NewSuperSetConnector(baseURL, username, password string, logger *slog.Logger) *SuperSetConnector {
	return &SuperSetConnector{
		baseURL:  baseURL,
		username: username,
		password: password,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
		logger: logger.With("connector", "superset"),
	}
}

func NewSuperSetConnectorWithToken(baseURL, token string, logger *slog.Logger) *SuperSetConnector {
	// Clean the token by removing any whitespace/newlines
	cleanToken := strings.TrimSpace(strings.ReplaceAll(token, "\n", ""))

	return &SuperSetConnector{
		baseURL: baseURL,
		token:   cleanToken,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
		logger: logger.With("connector", "superset"),
	}
}

func (sc *SuperSetConnector) Authenticate(ctx context.Context) error {
	authPayload := map[string]string{
		"username": sc.username,
		"password": sc.password,
		"provider": "db",
	}

	jsonData, _ := json.Marshal(authPayload)

	req, err := http.NewRequestWithContext(ctx, "POST",
		sc.baseURL+"/api/v1/security/login",
		bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create authentication request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := sc.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to Superset login endpoint: %w", err)
	}
	defer resp.Body.Close()

	sc.logger.Debug("Authentication response", "status", resp.StatusCode)

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("authentication failed: invalid username or password")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("authentication failed with status %d", resp.StatusCode)
	}

	var authResp struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return fmt.Errorf("failed to parse authentication response: %w", err)
	}

	if authResp.AccessToken == "" {
		return fmt.Errorf("no access token received from Superset")
	}

	sc.token = authResp.AccessToken
	sc.logger.Info("SuperSet authenticated successfully")
	return nil
}

func (sc *SuperSetConnector) ExecuteSQL(ctx context.Context, sql string) (*SuperSetResponse, error) {
	if sc.token == "" {
		if err := sc.Authenticate(ctx); err != nil {
			return nil, err
		}
	}

	// Try to discover the correct database ID or use default
	query := SuperSetQuery{
		SQL:        sql,
		DatabaseID: 1, // Most Superset instances start with database ID 1
	}

	jsonData, _ := json.Marshal(query)

	req, err := http.NewRequestWithContext(ctx, "POST",
		sc.baseURL+"/api/v1/sqllab/execute/",
		bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+sc.token)

	resp, err := sc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result SuperSetResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (sc *SuperSetConnector) GetSampleData(ctx context.Context) (*SuperSetResponse, error) {
	sql := `
    SELECT
        quarter,
        bike_category,
        total_revenue,
        total_bikes_sold
    FROM bike_sales
    ORDER BY quarter, bike_category
    `

	return sc.ExecuteSQL(ctx, sql)
}

// GetDatasets retrieves available datasets from Superset
func (sc *SuperSetConnector) GetDatasets(ctx context.Context) ([]map[string]interface{}, error) {
	if sc.token == "" {
		if err := sc.Authenticate(ctx); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, "GET", sc.baseURL+"/api/v1/dataset/", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create dataset request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+sc.token)
	req.Header.Set("Accept", "application/json")

	resp, err := sc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get datasets: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("datasets API returned status %d", resp.StatusCode)
	}

	var result struct {
		Result []map[string]interface{} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse datasets response: %w", err)
	}

	return result.Result, nil
}

// GetDashboards retrieves all dashboards from Superset
func (sc *SuperSetConnector) GetDashboards(ctx context.Context) ([]map[string]interface{}, error) {
	if sc.token == "" {
		if err := sc.Authenticate(ctx); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, "GET", sc.baseURL+"/api/v1/dashboard/", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create dashboard request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+sc.token)
	req.Header.Set("Accept", "application/json")

	resp, err := sc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboards: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dashboards API returned status %d", resp.StatusCode)
	}

	var result struct {
		Result []map[string]interface{} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse dashboards response: %w", err)
	}

	return result.Result, nil
}

// GetCharts retrieves all charts from Superset
func (sc *SuperSetConnector) GetCharts(ctx context.Context) ([]map[string]interface{}, error) {
	if sc.token == "" {
		if err := sc.Authenticate(ctx); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, "GET", sc.baseURL+"/api/v1/chart/", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create chart request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+sc.token)
	req.Header.Set("Accept", "application/json")

	resp, err := sc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get charts: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("charts API returned status %d", resp.StatusCode)
	}

	var result struct {
		Result []map[string]interface{} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse charts response: %w", err)
	}

	return result.Result, nil
}

// GetChartData retrieves data from a specific chart
func (sc *SuperSetConnector) GetChartData(ctx context.Context, chartID int) ([]map[string]interface{}, error) {
	if sc.token == "" {
		if err := sc.Authenticate(ctx); err != nil {
			return nil, err
		}
	}

	url := fmt.Sprintf("%s/api/v1/chart/data", sc.baseURL)

	payload := map[string]interface{}{
		"queries": []map[string]interface{}{
			{
				"datasource": map[string]interface{}{
					"id":   chartID,
					"type": "table",
				},
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chart data request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create chart data request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+sc.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := sc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get chart data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("chart data API returned status %d", resp.StatusCode)
	}

	var result struct {
		Result []map[string]interface{} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse chart data response: %w", err)
	}

	// Extract data from the nested result structure
	var data []map[string]interface{}
	if len(result.Result) > 0 {
		if queryResult, ok := result.Result[0]["data"].([]interface{}); ok {
			for _, item := range queryResult {
				if mapItem, ok := item.(map[string]interface{}); ok {
					data = append(data, mapItem)
				}
			}
		}
	}

	return data, nil
}

// GetDashboardData retrieves data from all charts in a dashboard
func (sc *SuperSetConnector) GetDashboardData(ctx context.Context, dashboardID int) ([]map[string]interface{}, error) {
	if sc.token == "" {
		if err := sc.Authenticate(ctx); err != nil {
			return nil, err
		}
	}

	// First get dashboard details to find charts
	dashboardURL := fmt.Sprintf("%s/api/v1/dashboard/%d", sc.baseURL, dashboardID)
	req, err := http.NewRequestWithContext(ctx, "GET", dashboardURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create dashboard detail request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+sc.token)
	req.Header.Set("Accept", "application/json")

	resp, err := sc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard details: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dashboard detail API returned status %d", resp.StatusCode)
	}

	var dashboardResult struct {
		Result map[string]interface{} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&dashboardResult); err != nil {
		return nil, fmt.Errorf("failed to parse dashboard response: %w", err)
	}

	// Extract chart IDs from dashboard
	var allData []map[string]interface{}
	if charts, ok := dashboardResult.Result["charts"].([]interface{}); ok {
		for _, chart := range charts {
			if chartMap, ok := chart.(map[string]interface{}); ok {
				if chartID, ok := chartMap["id"].(float64); ok {
					chartData, err := sc.GetChartData(ctx, int(chartID))
					if err != nil {
						sc.logger.Warn("Failed to get chart data", "chart_id", int(chartID), "error", err)
						continue
					}
					allData = append(allData, chartData...)
				}
			}
		}
	}

	return allData, nil
}

// GetDatabaseTables retrieves tables from a specific database
func (sc *SuperSetConnector) GetDatabaseTables(ctx context.Context, databaseID int) ([]map[string]interface{}, error) {
	if sc.token == "" {
		if err := sc.Authenticate(ctx); err != nil {
			return nil, err
		}
	}

	url := fmt.Sprintf("%s/api/v1/database/%d/tables/", sc.baseURL, databaseID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create tables request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+sc.token)
	req.Header.Set("Accept", "application/json")

	resp, err := sc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get tables: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tables API returned status %d", resp.StatusCode)
	}

	var result struct {
		Result []map[string]interface{} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse tables response: %w", err)
	}

	return result.Result, nil
}

// QueryDataset executes a query based on user intent using multiple data retrieval methods
func (sc *SuperSetConnector) QueryDataset(ctx context.Context, userQuery string) (*SuperSetResponse, error) {
	sc.logger.Info("🚀 ENTERING QueryDataset with comprehensive data fetching", "query", userQuery)

	// Try multiple data retrieval methods based on query type
	data, err := sc.getRelevantData(ctx, userQuery)
	if err != nil {
		sc.logger.Warn("🔄 Comprehensive data fetching failed, falling back to SQL", "error", err)
		// Fallback to SQL generation
		sql := sc.generateSQLFromQuery(userQuery)
		sc.logger.Info("Generated SQL from user query", "query", userQuery, "sql", sql)
		return sc.ExecuteSQL(ctx, sql)
	}

	sc.logger.Info("✅ Comprehensive data fetching succeeded", "rows", len(data))
	return &SuperSetResponse{
		Data:   data,
		Status: "success",
	}, nil
}

// getRelevantData tries multiple approaches to get relevant data based on user query
func (sc *SuperSetConnector) getRelevantData(ctx context.Context, userQuery string) ([]map[string]interface{}, error) {
	queryLower := strings.ToLower(userQuery)
	sc.logger.Info("Starting comprehensive data retrieval", "query", userQuery)

	// Strategy 1: Try to get data from dashboards (most comprehensive)
	if strings.Contains(queryLower, "dashboard") || strings.Contains(queryLower, "game") ||
	   strings.Contains(queryLower, "gaming") || strings.Contains(queryLower, "top") ||
	   strings.Contains(queryLower, "analysis") {

		sc.logger.Info("Attempting to fetch dashboard data")
		dashboards, err := sc.GetDashboards(ctx)
		sc.logger.Info("Dashboard API response", "error", err, "count", len(dashboards))
		if err == nil && len(dashboards) > 0 {
			// Try to get data from the first available dashboard
			for _, dashboard := range dashboards {
				if dashboardID, ok := dashboard["id"].(float64); ok {
					sc.logger.Info("Trying to get data from dashboard", "dashboard_id", int(dashboardID))
					dashboardData, err := sc.GetDashboardData(ctx, int(dashboardID))
					sc.logger.Info("Dashboard data retrieval result", "dashboard_id", int(dashboardID), "error", err, "rows", len(dashboardData))
					if err == nil && len(dashboardData) > 0 {
						sc.logger.Info("Successfully retrieved dashboard data", "dashboard_id", int(dashboardID), "rows", len(dashboardData))
						return dashboardData, nil
					}
				}
			}
		}
		sc.logger.Warn("Dashboard data retrieval failed", "error", err)
	}

	// Strategy 2: Try to get data from individual charts
	sc.logger.Info("Attempting to fetch chart data")
	charts, err := sc.GetCharts(ctx)
	sc.logger.Info("Chart API response", "error", err, "count", len(charts))
	if err == nil && len(charts) > 0 {
		var allChartData []map[string]interface{}
		for i, chart := range charts {
			if i >= 5 { // Limit to first 5 charts to avoid too much data
				break
			}
			if chartID, ok := chart["id"].(float64); ok {
				sc.logger.Info("Trying to get data from chart", "chart_id", int(chartID))
				chartData, err := sc.GetChartData(ctx, int(chartID))
				sc.logger.Info("Chart data retrieval result", "chart_id", int(chartID), "error", err, "rows", len(chartData))
				if err == nil && len(chartData) > 0 {
					allChartData = append(allChartData, chartData...)
				}
			}
		}
		if len(allChartData) > 0 {
			sc.logger.Info("Successfully retrieved chart data", "rows", len(allChartData))
			return allChartData, nil
		}
	}
	sc.logger.Warn("Chart data retrieval failed", "error", err)

	// Strategy 3: Try to get datasets and explore their data
	sc.logger.Info("Attempting to fetch dataset information")
	datasets, err := sc.GetDatasets(ctx)
	sc.logger.Info("Dataset API response", "error", err, "count", len(datasets))
	if err == nil && len(datasets) > 0 {
		// Return dataset metadata as structured data
		var datasetInfo []map[string]interface{}
		for _, dataset := range datasets {
			datasetInfo = append(datasetInfo, map[string]interface{}{
				"dataset_name":    dataset["table_name"],
				"database_id":     dataset["database_id"],
				"schema":          dataset["schema"],
				"table_name":      dataset["table_name"],
				"owners":          dataset["owners"],
				"created_on":      dataset["created_on_delta_humanized"],
				"changed_on":      dataset["changed_on_delta_humanized"],
			})
		}
		if len(datasetInfo) > 0 {
			sc.logger.Info("Successfully retrieved dataset information", "datasets", len(datasetInfo))
			return datasetInfo, nil
		}
	}
	sc.logger.Warn("Dataset retrieval failed", "error", err)

	// Strategy 4: Try to explore database tables
	sc.logger.Info("Attempting to fetch database tables")
	// Try database ID 1 (most common default)
	tables, err := sc.GetDatabaseTables(ctx, 1)
	sc.logger.Info("Database tables API response", "error", err, "count", len(tables))
	if err == nil && len(tables) > 0 {
		var tableInfo []map[string]interface{}
		for _, table := range tables {
			tableInfo = append(tableInfo, map[string]interface{}{
				"table_name": table["value"],
				"table_type": table["type"],
				"extra":      table["extra"],
			})
		}
		if len(tableInfo) > 0 {
			sc.logger.Info("Successfully retrieved table information", "tables", len(tableInfo))
			return tableInfo, nil
		}
	}
	sc.logger.Warn("Database tables retrieval failed", "error", err)

	sc.logger.Error("All comprehensive data retrieval strategies failed")
	return nil, fmt.Errorf("all data retrieval strategies failed")
}

// generateSQLFromQuery creates SQL queries based on user intent
func (sc *SuperSetConnector) generateSQLFromQuery(userQuery string) string {
	queryLower := strings.ToLower(userQuery)

	// Gaming/Entertainment related queries - try very simple approach first
	if strings.Contains(queryLower, "game") || strings.Contains(queryLower, "gaming") ||
	   strings.Contains(queryLower, "entertainment") || strings.Contains(queryLower, "top game") {

		// Try the simplest possible query first - just list tables
		return `SELECT 1 as test_value, 'gaming query test' as test_message`
	}

	// Sales performance queries
	if strings.Contains(queryLower, "sales") || strings.Contains(queryLower, "revenue") {
		return `
		SELECT *
		FROM (
			SELECT
				COALESCE(product_name, name, item_name) as name,
				COALESCE(category, type, class) as category,
				COALESCE(total_sales, sales, revenue, amount) as value,
				COALESCE(quarter, period, month) as time_period,
				COALESCE(region, location, area) as region
			FROM sales_data

			UNION ALL

			SELECT
				COALESCE(bike_model, product_name, name) as name,
				COALESCE(bike_category, category, type) as category,
				COALESCE(total_revenue, revenue, sales) as value,
				COALESCE(quarter, period, month_year) as time_period,
				COALESCE(region, store_id, location) as region
			FROM bike_sales
		) combined
		WHERE name IS NOT NULL AND value IS NOT NULL
		ORDER BY value DESC
		LIMIT 20`
	}

	// User activity queries
	if strings.Contains(queryLower, "user") || strings.Contains(queryLower, "active") {
		return `
		SELECT
			COALESCE(month_year, period, date) as time_period,
			COALESCE(active_users, users, count) as active_count,
			COALESCE(new_users, new_count, acquisitions) as new_count,
			COALESCE(platform, source, channel) as platform,
			COALESCE(region, location, area) as region
		FROM monthly_active_users
		WHERE month_year IS NOT NULL OR period IS NOT NULL
		ORDER BY time_period DESC
		LIMIT 12`
	}

	// Default: Explore available data with flexible column detection
	return `
	SELECT
		table_name,
		column_name,
		data_type
	FROM information_schema.columns
	WHERE table_schema = 'public'
	AND table_name NOT LIKE 'pg_%'
	AND table_name NOT LIKE 'sql_%'
	ORDER BY table_name, ordinal_position
	LIMIT 50`
}

func (sc *SuperSetConnector) HealthCheck(ctx context.Context) error {
	// Try multiple health check endpoints that Superset might use
	endpoints := []string{"/api/v1/security/csrf_token/", "/health", "/heartbeat"}

	var lastErr error
	for _, endpoint := range endpoints {
		req, err := http.NewRequestWithContext(ctx, "GET", sc.baseURL+endpoint, nil)
		if err != nil {
			lastErr = fmt.Errorf("failed to create request for %s: %w", endpoint, err)
			sc.logger.Debug("Failed to create request", "endpoint", endpoint, "error", err)
			continue
		}

		resp, err := sc.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("failed to connect to %s: %w", endpoint, err)
			sc.logger.Debug("Failed to connect", "endpoint", endpoint, "error", err)
			continue
		}
		resp.Body.Close()

		sc.logger.Debug("Health check response", "endpoint", endpoint, "status", resp.StatusCode)

		// If any endpoint responds with 200-299, consider it healthy
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			sc.logger.Info("Health check successful", "endpoint", endpoint, "status", resp.StatusCode)
			return nil
		}

		lastErr = fmt.Errorf("endpoint %s returned status %d", endpoint, resp.StatusCode)
	}

	if lastErr != nil {
		return fmt.Errorf("superset health check failed: %w", lastErr)
	}
	return fmt.Errorf("superset health check failed: no working endpoints found")
}

func (sc *SuperSetConnector) TestConnection(ctx context.Context) error {
	// First try health check
	if err := sc.HealthCheck(ctx); err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	// If we have a token, test with it directly, otherwise try authentication
	if sc.token != "" {
		return sc.testTokenAuth(ctx)
	}

	// Try authentication with username/password
	if err := sc.Authenticate(ctx); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	return nil
}

func (sc *SuperSetConnector) testTokenAuth(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", sc.baseURL+"/api/v1/chart/", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+sc.token)
	req.Header.Set("Accept", "application/json")

	resp, err := sc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("bearer token authentication failed: invalid or expired token")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	sc.logger.Info("SuperSet token authentication successful")
	return nil
}
