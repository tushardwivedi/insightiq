package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"insightiq/backend/internal/domain"
	"insightiq/backend/internal/events"
	"insightiq/backend/internal/repository"
	"insightiq/backend/internal/security"
	"insightiq/backend/internal/usecases"
)

func getTestDB(t *testing.T) *sqlx.DB {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://insightiq_user:insightiq_password@localhost:5432/insightiq?sslmode=disable"
	}
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		t.Skipf("Skipping: cannot connect to database: %v", err)
	}
	return db
}

// injectUserContext simulates the auth middleware by setting user_id, ip_address, user_agent in context.
func injectUserContext(next http.Handler, userID, ip, ua string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "user_id", userID)
		ctx = context.WithValue(ctx, "ip_address", ip)
		ctx = context.WithValue(ctx, "user_agent", ua)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// setupTestStack creates the full event pipeline and returns handler + publisher.
func setupTestStack(t *testing.T, db *sqlx.DB) (*QueryHistoryHandler, *events.SyncEventPublisher, *repository.QueryHistoryRepository) {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))

	repo := repository.NewQueryHistoryRepository(db)
	if err := repo.CreateTables(context.Background()); err != nil {
		t.Fatalf("CreateTables: %v", err)
	}
	// Clean test data
	db.Exec("DELETE FROM query_history WHERE user_id LIKE 'http-test-%'")

	redactor := security.NewDefaultPIIRedactor(logger)
	uc := usecases.NewQueryHistoryUseCase(repo, redactor, logger)
	publisher := events.NewSyncEventPublisher(logger)
	evtHandler := events.NewQueryHistoryEventHandler(uc, logger)
	publisher.Subscribe("query.executed", evtHandler)

	handler := NewQueryHistoryHandler(repo, logger)
	return handler, publisher, repo
}

// publishEvent is a helper to publish a query event through the event pipeline.
func publishEvent(t *testing.T, pub *events.SyncEventPublisher, userID, query, sql, ip, ua string) {
	t.Helper()
	evt := &domain.QueryExecutionEvent{
		EventID:       "http-evt-" + userID,
		UserID:        userID,
		Timestamp:     time.Now(),
		QueryType:     "text",
		QueryText:     query,
		GeneratedSQL:  sql,
		ConnectorName: "test",
		Status:        "success",
		RowCount:      1,
		ExecutionTime: 42,
		IPAddress:     ip,
		UserAgent:     ua,
		ResultPreview: map[string]interface{}{"rows": 1},
	}
	if err := pub.Publish(context.Background(), "query.executed", evt); err != nil {
		t.Fatalf("publish: %v", err)
	}
}

func TestHTTP_PIIRedactionEndToEnd(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	handler, publisher, _ := setupTestStack(t, db)

	userID := "http-test-pii-001"
	clientIP := "198.51.100.42"
	clientUA := "TestBrowser/1.0"

	// Publish a query containing PII through the event pipeline
	publishEvent(t, publisher, userID,
		"Show orders for john@example.com and phone 555-123-4567",
		"SELECT * FROM orders WHERE email = 'john@example.com'",
		clientIP, clientUA,
	)

	// Now call the HTTP handler to list query history
	req := httptest.NewRequest("GET", "/api/query-history?limit=10", nil)
	rec := httptest.NewRecorder()
	injectUserContext(http.HandlerFunc(handler.ListQueryHistory), userID, clientIP, clientUA).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Data []struct {
			ID        string `json:"id"`
			QueryText string `json:"query_text"`
		} `json:"data"`
		Count int `json:"count"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if resp.Count != 1 {
		t.Fatalf("Expected 1 result, got %d", resp.Count)
	}

	stored := resp.Data[0].QueryText

	// Verify PII is redacted in HTTP response
	if strings.Contains(stored, "john@example.com") {
		t.Errorf("Email not redacted in HTTP response: %s", stored)
	}
	if strings.Contains(stored, "555-123-4567") {
		t.Errorf("Phone not redacted in HTTP response: %s", stored)
	}
	if !strings.Contains(stored, "[EMAIL_REDACTED]") {
		t.Errorf("Expected [EMAIL_REDACTED], got: %s", stored)
	}
	if !strings.Contains(stored, "[PHONE_REDACTED]") {
		t.Errorf("Expected [PHONE_REDACTED], got: %s", stored)
	}

	// Verify full details via GetByID endpoint
	detailReq := httptest.NewRequest("GET", "/api/query-history/detail?id="+resp.Data[0].ID, nil)
	detailRec := httptest.NewRecorder()
	injectUserContext(http.HandlerFunc(handler.GetQueryHistoryByID), userID, clientIP, clientUA).ServeHTTP(detailRec, detailReq)

	if detailRec.Code != http.StatusOK {
		t.Fatalf("Detail expected 200, got %d", detailRec.Code)
	}

	var detailResp struct {
		Data struct {
			QueryText    string `json:"query_text"`
			GeneratedSQL string `json:"generated_sql"`
			IPAddress    string `json:"ip_address"`
			UserAgent    string `json:"user_agent"`
			DataRedacted bool   `json:"data_redacted"`
		} `json:"data"`
	}
	if err := json.NewDecoder(detailRec.Body).Decode(&detailResp); err != nil {
		t.Fatalf("decode detail: %v", err)
	}

	if !detailResp.Data.DataRedacted {
		t.Error("Expected data_redacted=true")
	}
	if strings.Contains(detailResp.Data.GeneratedSQL, "john@example.com") {
		t.Error("Email not redacted in generated_sql")
	}
	if detailResp.Data.IPAddress != clientIP {
		t.Errorf("Expected IP %s, got %s", clientIP, detailResp.Data.IPAddress)
	}
	if detailResp.Data.UserAgent != clientUA {
		t.Errorf("Expected UA %s, got %s", clientUA, detailResp.Data.UserAgent)
	}

	t.Logf("HTTP PII test passed - stored: %s", stored)
}

func TestHTTP_UserIsolation(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	handler, publisher, _ := setupTestStack(t, db)

	userA := "http-test-iso-A"
	userB := "http-test-iso-B"

	publishEvent(t, publisher, userA, "User A secret query", "SELECT 1", "10.0.0.1", "UA-A")
	publishEvent(t, publisher, userB, "User B secret query", "SELECT 2", "10.0.0.2", "UA-B")

	// User A lists history — should only see their own
	req := httptest.NewRequest("GET", "/api/query-history?limit=50", nil)
	rec := httptest.NewRecorder()
	injectUserContext(http.HandlerFunc(handler.ListQueryHistory), userA, "10.0.0.1", "UA-A").ServeHTTP(rec, req)

	var resp struct {
		Data []struct {
			ID        string `json:"id"`
			QueryText string `json:"query_text"`
		} `json:"data"`
		Count int `json:"count"`
	}
	json.NewDecoder(rec.Body).Decode(&resp)

	if resp.Count != 1 {
		t.Fatalf("User A should see 1 query, got %d", resp.Count)
	}
	if !strings.Contains(resp.Data[0].QueryText, "User A") {
		t.Errorf("User A should see own query, got: %s", resp.Data[0].QueryText)
	}

	// User A tries to get User B's query detail — should fail
	// First get User B's query ID
	reqB := httptest.NewRequest("GET", "/api/query-history?limit=50", nil)
	recB := httptest.NewRecorder()
	injectUserContext(http.HandlerFunc(handler.ListQueryHistory), userB, "10.0.0.2", "UA-B").ServeHTTP(recB, reqB)

	var respB struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	json.NewDecoder(recB.Body).Decode(&respB)

	if len(respB.Data) > 0 {
		crossReq := httptest.NewRequest("GET", "/api/query-history/detail?id="+respB.Data[0].ID, nil)
		crossRec := httptest.NewRecorder()
		injectUserContext(http.HandlerFunc(handler.GetQueryHistoryByID), userA, "10.0.0.1", "UA-A").ServeHTTP(crossRec, crossReq)

		if crossRec.Code == http.StatusOK {
			t.Error("User A should NOT access User B's query detail via HTTP")
		}
	}

	t.Log("HTTP user isolation test passed")
}

func TestHTTP_NoPII_QueryUnchanged(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	handler, publisher, _ := setupTestStack(t, db)

	userID := "http-test-nopii-001"
	originalQuery := "Show total revenue for Q4"

	publishEvent(t, publisher, userID, originalQuery, "SELECT SUM(revenue) FROM sales", "10.0.0.1", "TestUA")

	req := httptest.NewRequest("GET", "/api/query-history?limit=10", nil)
	rec := httptest.NewRecorder()
	injectUserContext(http.HandlerFunc(handler.ListQueryHistory), userID, "10.0.0.1", "TestUA").ServeHTTP(rec, req)

	var resp struct {
		Data []struct {
			QueryText string `json:"query_text"`
		} `json:"data"`
	}
	json.NewDecoder(rec.Body).Decode(&resp)

	if len(resp.Data) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(resp.Data))
	}
	if resp.Data[0].QueryText != originalQuery {
		t.Errorf("Non-PII query should be unchanged. Expected: %s, Got: %s", originalQuery, resp.Data[0].QueryText)
	}

	t.Log("HTTP no-PII test passed - query stored unchanged")
}
