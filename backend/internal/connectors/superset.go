package connectors

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"insightiq/backend/internal/cache"
)

type SuperSetConnector struct {
	baseURL   string
	username  string
	password  string
	token     string
	csrfToken string // CSRF token extracted from JWT for POST requests
	client    *http.Client
	logger    *slog.Logger
	cache     *cache.RedisCache // Optional cache
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

// SetCache sets the Redis cache for this connector
func (sc *SuperSetConnector) SetCache(cache *cache.RedisCache) {
	sc.cache = cache
	sc.logger.Info("✅ Redis cache enabled for Superset connector")
}

// generateCacheKey creates a cache key from query components
func generateCacheKey(prefix string, components ...string) string {
	h := sha256.New()
	h.Write([]byte(prefix))
	for _, comp := range components {
		h.Write([]byte(comp))
	}
	return prefix + ":" + hex.EncodeToString(h.Sum(nil))[:16]
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

	// Extract CSRF token from JWT payload
	sc.csrfToken = sc.extractCSRFFromJWT(authResp.AccessToken)
	if sc.csrfToken != "" {
		sc.logger.Debug("Extracted CSRF token from JWT", "csrf_token", sc.csrfToken[:16]+"...")
	}

	sc.logger.Info("SuperSet authenticated successfully")
	return nil
}

// extractCSRFFromJWT extracts the CSRF token from the JWT payload
func (sc *SuperSetConnector) extractCSRFFromJWT(token string) string {
	// JWT format: header.payload.signature
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return ""
	}

	// Decode the payload (second part)
	payload, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		sc.logger.Warn("Failed to decode JWT payload", "error", err)
		return ""
	}

	// Parse JSON to extract CSRF
	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		sc.logger.Warn("Failed to parse JWT payload", "error", err)
		return ""
	}

	// Extract CSRF token
	if csrf, ok := claims["csrf"].(string); ok {
		return csrf
	}

	return ""
}

func (sc *SuperSetConnector) ExecuteSQL(ctx context.Context, sql string) (*SuperSetResponse, error) {
	// Try cache first if available (TTL: 5 minutes for query results)
	cacheKey := generateCacheKey("superset:query", sc.baseURL, sql)
	if sc.cache != nil {
		var cachedResult SuperSetResponse
		err := sc.cache.GetJSON(ctx, cacheKey, &cachedResult)
		if err == nil && len(cachedResult.Data) > 0 {
			sc.logger.Debug("✅ Cache hit: query result", "rows", len(cachedResult.Data))
			return &cachedResult, nil
		}
		sc.logger.Debug("Cache miss: query")
	}

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
	// Add CSRF token for POST requests
	if sc.csrfToken != "" {
		req.Header.Set("X-CSRFToken", sc.csrfToken)
	}

	resp, err := sc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result SuperSetResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Store in cache if available (TTL: 5 minutes)
	if sc.cache != nil && len(result.Data) > 0 {
		err := sc.cache.Set(ctx, cacheKey, result, 5*time.Minute)
		if err != nil {
			sc.logger.Warn("Failed to cache query result", "error", err)
		} else {
			sc.logger.Debug("✅ Cached query result", "rows", len(result.Data), "ttl", "5m")
		}
	}

	return &result, nil
}

// extractKeywords extracts meaningful keywords from a user query
func extractKeywords(query string) []string {
	// Remove common stop words
	stopWords := map[string]bool{
		"give": true, "me": true, "the": true, "a": true, "an": true,
		"from": true, "on": true, "in": true, "of": true, "to": true,
		"show": true, "get": true, "find": true, "what": true, "how": true,
		"is": true, "are": true, "was": true, "were": true, "be": true,
		"have": true, "has": true, "had": true, "do": true, "does": true,
		"did": true, "will": true, "would": true, "should": true, "could": true,
		"data": true, "s": true, // 's' from possessives like "Bank's"
	}

	// Clean and split query
	words := strings.Fields(strings.ToLower(query))
	var keywords []string

	for _, word := range words {
		// Remove punctuation
		word = strings.Trim(word, ".,!?;:'\"")

		// Skip stop words and short words
		if !stopWords[word] && len(word) > 2 {
			keywords = append(keywords, word)
		}
	}

	return keywords
}

// findBestMatchingDashboard scores dashboards by keyword matches and returns the best match
func findBestMatchingDashboard(dashboards []map[string]interface{}, keywords []string) (map[string]interface{}, int) {
	var bestMatch map[string]interface{}
	maxScore := 0

	for _, dash := range dashboards {
		score := 0

		// Extract dashboard name and title
		dashName := ""
		if name, ok := dash["dashboard_title"].(string); ok {
			dashName = strings.ToLower(name)
		} else if name, ok := dash["title"].(string); ok {
			dashName = strings.ToLower(name)
		}

		// Score based on keyword matches
		for _, keyword := range keywords {
			if strings.Contains(dashName, keyword) {
				score += 2 // Higher weight for exact matches
			}
			// Check for partial matches (e.g., "world" matches "worldwide")
			for _, dashWord := range strings.Fields(dashName) {
				if strings.Contains(dashWord, keyword) || strings.Contains(keyword, dashWord) {
					score++
				}
			}
		}

		if score > maxScore {
			maxScore = score
			bestMatch = dash
		}
	}

	return bestMatch, maxScore
}

// FindRelevantDataset searches for a dataset/dashboard based on query keywords
func (sc *SuperSetConnector) FindRelevantDataset(ctx context.Context, query string) (map[string]interface{}, error) {
	// Extract keywords from query
	keywords := extractKeywords(query)
	sc.logger.Info("Extracted keywords from query", "keywords", keywords)

	if len(keywords) == 0 {
		return nil, fmt.Errorf("no meaningful keywords found in query")
	}

	// Get all available dashboards
	dashboards, err := sc.GetDashboards(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboards: %w", err)
	}

	sc.logger.Info("Retrieved dashboards", "count", len(dashboards))

	// Find best matching dashboard
	bestMatch, score := findBestMatchingDashboard(dashboards, keywords)

	if bestMatch == nil || score == 0 {
		return nil, fmt.Errorf("no matching dashboard found for keywords: %v", keywords)
	}

	sc.logger.Info("Found matching dashboard",
		"dashboard", bestMatch["dashboard_title"],
		"score", score,
		"id", bestMatch["id"])

	return bestMatch, nil
}

// QueryDashboardData retrieves data from a dashboard by querying its datasets
func (sc *SuperSetConnector) QueryDashboardData(ctx context.Context, dashboardID int, query string) (*SuperSetResponse, error) {
	sc.logger.Info("Querying dashboard data", "dashboard_id", dashboardID)

	// Get dashboard details to find associated datasets
	dashboard, err := sc.GetDashboards(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboards: %w", err)
	}

	// Find the specific dashboard
	var targetDashboard map[string]interface{}
	for _, dash := range dashboard {
		if id, ok := dash["id"].(float64); ok && int(id) == dashboardID {
			targetDashboard = dash
			break
		}
	}

	if targetDashboard == nil {
		return nil, fmt.Errorf("dashboard %d not found", dashboardID)
	}

	// Extract dashboard title for proper error messages
	dashboardTitle := "Unknown Dashboard"
	if title, ok := targetDashboard["dashboard_title"].(string); ok {
		dashboardTitle = title
	} else if title, ok := targetDashboard["title"].(string); ok {
		dashboardTitle = title
	}

	sc.logger.Info("Processing dashboard", "title", dashboardTitle, "id", dashboardID)

	// Try to get data from dashboard charts
	data, err := sc.GetDashboardData(ctx, dashboardID)
	if err != nil {
		sc.logger.Warn("Failed to get dashboard chart data", "dashboard", dashboardTitle, "dashboard_id", dashboardID, "error", err)
	}

	if len(data) > 0 {
		sc.logger.Info("✅ Retrieved data from dashboard charts", "dashboard", dashboardTitle, "rows", len(data))
		return &SuperSetResponse{
			Data:   data,
			Status: "success",
		}, nil
	}

	// Fallback: Use generic SQL to explore the database
	sc.logger.Warn("No data from dashboard charts, using generic discovery query", "dashboard", dashboardTitle)

	// Try to discover what tables are available
	discoverSQL := `
	SELECT
		table_name,
		column_name,
		data_type
	FROM information_schema.columns
	WHERE table_schema = 'public'
		AND table_name NOT LIKE 'pg_%'
		AND table_name NOT LIKE 'sql_%'
	ORDER BY table_name, ordinal_position
	LIMIT 100
	`

	result, err := sc.ExecuteSQL(ctx, discoverSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to discover database schema: %w", err)
	}

	if len(result.Data) > 0 {
		// Return table information so user knows what's available
		sc.logger.Info("Returning database schema information", "dashboard", dashboardTitle, "tables_found", len(result.Data))
		return &SuperSetResponse{
			Data: []map[string]interface{}{
				{
					"message": fmt.Sprintf("Found '%s' dashboard but no data available from charts. Discovered %d tables/columns in the database. Please check the dashboard configuration in Superset.", dashboardTitle, len(result.Data)),
					"schema_info":    result.Data,
					"dashboard_title": dashboardTitle,
				},
			},
			Status: "partial",
		}, nil
	}

	return &SuperSetResponse{
		Data: []map[string]interface{}{
			{
				"error":           fmt.Sprintf("'%s' dashboard found but contains no accessible data. The dashboard may need to be configured with datasets in Superset.", dashboardTitle),
				"dashboard_id":    dashboardID,
				"dashboard_title": dashboardTitle,
			},
		},
		Status: "error",
	}, nil
}

// generateDashboardSQL generates appropriate SQL based on dashboard context
func (sc *SuperSetConnector) generateDashboardSQL(dashboardTitle, query string) string {
	// For World Bank / Health dashboards
	if strings.Contains(dashboardTitle, "world") || strings.Contains(dashboardTitle, "health") {
		return `
		SELECT
			country_name,
			country_code,
			year,
			sp_pop_totl as population,
			sh_dyn_mort as mortality_rate,
			sp_dyn_le00_in as life_expectancy,
			sh_xpd_chex_pc_cd as health_expenditure_per_capita
		FROM wb_health_population
		WHERE year >= 2015
		ORDER BY year DESC, country_name
		LIMIT 100
		`
	}

	// For bike/vehicle sales dashboards
	if strings.Contains(dashboardTitle, "bike") || strings.Contains(dashboardTitle, "vehicle") ||
		strings.Contains(dashboardTitle, "sales") {
		return `
		SELECT
			quarter,
			bike_category,
			total_revenue,
			total_bikes_sold
		FROM bike_sales
		ORDER BY quarter, bike_category
		LIMIT 100
		`
	}

	// Generic query - try to discover tables
	return `
	SELECT table_name, column_name, data_type
	FROM information_schema.columns
	WHERE table_schema = 'public'
	AND table_name NOT LIKE 'pg_%'
	ORDER BY table_name, ordinal_position
	LIMIT 50
	`
}

// DEPRECATED: Use FindRelevantDataset and query actual data instead
// This method always returns bike_sales data and should not be used
func (sc *SuperSetConnector) GetSampleData(ctx context.Context) (*SuperSetResponse, error) {
	sc.logger.Warn("⚠️  GetSampleData() is DEPRECATED - returns hardcoded bike_sales data")
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
	// Try cache first if available
	cacheKey := generateCacheKey("superset:dashboards", sc.baseURL)
	if sc.cache != nil {
		var cachedDashboards []map[string]interface{}
		err := sc.cache.GetJSON(ctx, cacheKey, &cachedDashboards)
		if err == nil && len(cachedDashboards) > 0 {
			sc.logger.Debug("✅ Cache hit: dashboards", "count", len(cachedDashboards))
			return cachedDashboards, nil
		}
		sc.logger.Debug("Cache miss: dashboards")
	}

	if sc.token == "" {
		if err := sc.Authenticate(ctx); err != nil {
			return nil, err
		}
	}

	// Add pagination parameters to avoid 422 errors
	url := fmt.Sprintf("%s/api/v1/dashboard/?q=(page:0,page_size:100)", sc.baseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
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
		// Read response body for better error diagnostics
		bodyBytes := make([]byte, 500) // Read first 500 bytes for error diagnostics
		n, _ := resp.Body.Read(bodyBytes)
		errorPreview := string(bodyBytes[:n])

		sc.logger.Error("Dashboard API returned non-OK status",
			"status", resp.StatusCode,
			"url", url,
			"response_preview", errorPreview)
		return nil, fmt.Errorf("dashboards API returned status %d", resp.StatusCode)
	}

	var result struct {
		Result []map[string]interface{} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse dashboards response: %w", err)
	}

	// Store in cache if available (TTL: 10 minutes)
	if sc.cache != nil && len(result.Result) > 0 {
		err := sc.cache.Set(ctx, cacheKey, result.Result, 10*time.Minute)
		if err != nil {
			sc.logger.Warn("Failed to cache dashboards", "error", err)
		} else {
			sc.logger.Debug("✅ Cached dashboards", "count", len(result.Result), "ttl", "10m")
		}
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
// Uses multiple strategies to get data, falling back when needed
func (sc *SuperSetConnector) GetChartData(ctx context.Context, chartID int) ([]map[string]interface{}, error) {
	if sc.token == "" {
		if err := sc.Authenticate(ctx); err != nil {
			return nil, err
		}
	}

	sc.logger.Debug("Fetching data for chart", "chart_id", chartID)

	// Strategy 1: Try GET /api/v1/chart/{id}/data/ (works if query_context is set)
	data, err := sc.getChartDataDirect(ctx, chartID)
	if err == nil && len(data) > 0 {
		sc.logger.Debug("✅ Got chart data via direct API", "chart_id", chartID, "rows", len(data))
		return data, nil
	}
	sc.logger.Debug("Direct chart data API failed or returned no data", "chart_id", chartID, "error", err)

	// Strategy 2: Get chart metadata and extract dataset info, then query via SQL Lab
	data, err = sc.getChartDataViaSQLLab(ctx, chartID)
	if err == nil && len(data) > 0 {
		sc.logger.Debug("✅ Got chart data via SQL Lab", "chart_id", chartID, "rows", len(data))
		return data, nil
	}
	sc.logger.Debug("SQL Lab query failed", "chart_id", chartID, "error", err)

	sc.logger.Warn("All methods to fetch chart data failed", "chart_id", chartID)
	return nil, fmt.Errorf("unable to fetch data for chart %d: all methods failed", chartID)
}

// getChartDataDirect tries to get chart data using GET /api/v1/chart/{id}/data/
func (sc *SuperSetConnector) getChartDataDirect(ctx context.Context, chartID int) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/chart/%d/data/", sc.baseURL, chartID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+sc.token)
	req.Header.Set("Accept", "application/json")

	resp, err := sc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	var result struct {
		Result []map[string]interface{} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Extract data from nested structure
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

// getChartDataViaSQLLab gets chart data by querying the underlying dataset via SQL Lab
func (sc *SuperSetConnector) getChartDataViaSQLLab(ctx context.Context, chartID int) ([]map[string]interface{}, error) {
	// First get chart metadata to find the datasource
	chartURL := fmt.Sprintf("%s/api/v1/chart/%d", sc.baseURL, chartID)

	req, err := http.NewRequestWithContext(ctx, "GET", chartURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+sc.token)
	req.Header.Set("Accept", "application/json")

	resp, err := sc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get chart metadata: status %d", resp.StatusCode)
	}

	var chartMeta struct {
		Result struct {
			DatasourceID   interface{} `json:"datasource_id"`
			DatasourceType string      `json:"datasource_type"`
			Params         string      `json:"params"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&chartMeta); err != nil {
		return nil, err
	}

	// Parse params to get form_data
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(chartMeta.Result.Params), &params); err != nil {
		sc.logger.Warn("Failed to parse chart params", "chart_id", chartID, "error", err)
		return nil, err
	}

	// Extract datasource info
	datasourceStr, _ := params["datasource"].(string)
	if datasourceStr == "" {
		return nil, fmt.Errorf("no datasource found in chart params")
	}

	// Parse datasource string (format: "id__type")
	parts := strings.Split(datasourceStr, "__")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid datasource format: %s", datasourceStr)
	}

	// Get the dataset metadata to find table name
	datasetID := parts[0]
	datasetURL := fmt.Sprintf("%s/api/v1/dataset/%s", sc.baseURL, datasetID)

	req2, err := http.NewRequestWithContext(ctx, "GET", datasetURL, nil)
	if err != nil {
		return nil, err
	}

	req2.Header.Set("Authorization", "Bearer "+sc.token)
	req2.Header.Set("Accept", "application/json")

	resp2, err := sc.client.Do(req2)
	if err != nil {
		return nil, err
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get dataset: status %d", resp2.StatusCode)
	}

	var dataset struct {
		Result struct {
			TableName  string `json:"table_name"`
			DatabaseID int    `json:"database_id"`
			SQL        string `json:"sql"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp2.Body).Decode(&dataset); err != nil {
		return nil, err
	}

	tableName := dataset.Result.TableName
	if tableName == "" && dataset.Result.SQL != "" {
		// Virtual dataset - use the SQL
		tableName = fmt.Sprintf("(%s) as virtual_table", dataset.Result.SQL)
	}

	if tableName == "" {
		return nil, fmt.Errorf("no table name or SQL found for dataset")
	}

	// Build SQL query based on chart params
	groupby, _ := params["groupby"].([]interface{})
	metric, _ := params["metric"].(string)
	rowLimit := 100 // default

	if limit, ok := params["row_limit"].(float64); ok && limit > 0 {
		rowLimit = int(limit)
	}

	// Construct SELECT clause
	var selectCols []string
	for _, col := range groupby {
		if colStr, ok := col.(string); ok {
			selectCols = append(selectCols, colStr)
		}
	}

	// Add metric
	if metric == "count" || metric == "" {
		selectCols = append(selectCols, "COUNT(*) as count")
	} else {
		selectCols = append(selectCols, fmt.Sprintf("SUM(%s) as %s", metric, metric))
	}

	// Build GROUP BY clause
	var groupByCols []string
	for _, col := range groupby {
		if colStr, ok := col.(string); ok {
			groupByCols = append(groupByCols, colStr)
		}
	}

	// Construct SQL
	var sql string
	if len(groupByCols) > 0 {
		sql = fmt.Sprintf("SELECT %s FROM %s GROUP BY %s LIMIT %d",
			strings.Join(selectCols, ", "),
			tableName,
			strings.Join(groupByCols, ", "),
			rowLimit)
	} else {
		sql = fmt.Sprintf("SELECT %s FROM %s LIMIT %d",
			strings.Join(selectCols, ", "),
			tableName,
			rowLimit)
	}

	sc.logger.Debug("Generated SQL for chart", "chart_id", chartID, "sql", sql)

	// Execute SQL via SQL Lab
	result, err := sc.ExecuteSQL(ctx, sql)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
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
	// Superset stores chart IDs in position_json, not in the charts array
	// The charts array contains chart names (strings), not chart objects
	var allData []map[string]interface{}
	var chartIDs []int

	// Try to extract chart IDs from position_json first (this is the correct way)
	if positionJSON, ok := dashboardResult.Result["position_json"].(string); ok {
		var positions map[string]interface{}
		if err := json.Unmarshal([]byte(positionJSON), &positions); err == nil {
			sc.logger.Debug("Parsing position_json for chart IDs")
			for key, value := range positions {
				if strings.HasPrefix(key, "CHART-") {
					if valueMap, ok := value.(map[string]interface{}); ok {
						if meta, ok := valueMap["meta"].(map[string]interface{}); ok {
							if chartID, ok := meta["chartId"].(float64); ok {
								chartIDs = append(chartIDs, int(chartID))
								sc.logger.Debug("Found chart ID in position_json", "chart_id", int(chartID))
							}
						}
					}
				}
			}
		}
	}

	// Fallback: Try old method (charts as array of objects) - for backward compatibility
	if len(chartIDs) == 0 {
		sc.logger.Debug("No charts found in position_json, trying legacy charts array")
		if charts, ok := dashboardResult.Result["charts"].([]interface{}); ok {
			for _, chart := range charts {
				if chartMap, ok := chart.(map[string]interface{}); ok {
					if chartID, ok := chartMap["id"].(float64); ok {
						chartIDs = append(chartIDs, int(chartID))
					}
				}
			}
		}
	}

	sc.logger.Info("Extracted chart IDs from dashboard", "dashboard_id", dashboardID, "chart_count", len(chartIDs), "chart_ids", chartIDs)

	// Fetch data from all charts
	for _, chartID := range chartIDs {
		chartData, err := sc.GetChartData(ctx, chartID)
		if err != nil {
			sc.logger.Warn("Failed to get chart data", "chart_id", chartID, "error", err)
			continue
		}
		if len(chartData) > 0 {
			sc.logger.Debug("Retrieved data from chart", "chart_id", chartID, "rows", len(chartData))
			allData = append(allData, chartData...)
		}
	}

	sc.logger.Info("Retrieved total data from dashboard", "dashboard_id", dashboardID, "total_rows", len(allData))
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
	sc.logger.Info("Starting comprehensive data retrieval", "query", userQuery)

	// Strategy 1: Try to discover available tables dynamically first (most reliable)
	// This was previously Strategy 2, but we prioritize it now
	sc.logger.Info("Attempting to discover available tables")
	metadataSQL := `
		SELECT table_name, column_name, data_type
		FROM information_schema.columns
		WHERE table_schema = 'public'
		AND table_name NOT LIKE 'pg_%'
		AND table_name NOT LIKE 'sql_%'
		ORDER BY table_name, ordinal_position
		LIMIT 50
	`

	sqlResult, err := sc.ExecuteSQL(ctx, metadataSQL)
	if err == nil && sqlResult != nil && len(sqlResult.Data) > 0 {
		sc.logger.Info("Successfully retrieved table metadata", "rows", len(sqlResult.Data))

		// Extract unique table names from metadata
		tableNames := make(map[string]bool)
		for _, row := range sqlResult.Data {
			if tableName, ok := row["table_name"].(string); ok {
				tableNames[tableName] = true
			}
		}

		sc.logger.Info("Found tables", "tables", tableNames)

		// Try to query the first few tables we found
		for tableName := range tableNames {
			querySQL := fmt.Sprintf("SELECT * FROM %s LIMIT 10", tableName)
			sc.logger.Info("Trying to query discovered table", "table", tableName)

			tableResult, err := sc.ExecuteSQL(ctx, querySQL)
			if err == nil && tableResult != nil && len(tableResult.Data) > 0 {
				sc.logger.Info("Successfully retrieved data from discovered table", "table", tableName, "rows", len(tableResult.Data))
				return tableResult.Data, nil
			}
		}

		// If we couldn't get actual data, return the metadata itself
		return sqlResult.Data, nil
	}

	// Strategy 2 (fallback): Return error instead of fake sample data
	// This ensures users know when real data isn't available rather than seeing misleading sample data
	sc.logger.Error("❌ Unable to retrieve data from Superset",
		"query", userQuery,
		"reason", "No accessible tables found and table discovery failed")

	return nil, fmt.Errorf("unable to retrieve data from Superset for query '%s'. Please verify: (1) Dashboard exists in Superset, (2) Database connection is configured, (3) User has necessary permissions", userQuery)
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
