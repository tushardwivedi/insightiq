package connectors

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
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
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := sc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var authResp struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return err
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

	query := SuperSetQuery{
		SQL:        sql,
		DatabaseID: 1,
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

func (sc *SuperSetConnector) HealthCheck(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", sc.baseURL+"/health", nil)
	if err != nil {
		return err
	}

	resp, err := sc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("superset health check failed: status %d", resp.StatusCode)
	}

	return nil
}
