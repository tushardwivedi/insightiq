package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
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

// setupGDPRTestStack creates the full event pipeline and returns GDPR handler + publisher.
func setupGDPRTestStack(t *testing.T, db *sqlx.DB) (*GDPRHandler, *events.SyncEventPublisher, *repository.QueryHistoryRepository) {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))

	repo := repository.NewQueryHistoryRepository(db)
	if err := repo.CreateTables(context.Background()); err != nil {
		t.Fatalf("CreateTables: %v", err)
	}
	// Clean test data
	db.Exec("DELETE FROM query_history WHERE user_id LIKE 'gdpr-test-%'")

	redactor := security.NewDefaultPIIRedactor(logger)
	uc := usecases.NewQueryHistoryUseCase(repo, redactor, logger)
	publisher := events.NewSyncEventPublisher(logger)
	evtHandler := events.NewQueryHistoryEventHandler(uc, logger)
	publisher.Subscribe("query.executed", evtHandler)

	handler := NewGDPRHandler(uc, logger)
	return handler, publisher, repo
}

// publishGDPREvent publishes a query event for GDPR testing.
func publishGDPREvent(t *testing.T, pub *events.SyncEventPublisher, userID, query, sql, ip, ua string) {
	t.Helper()
	evt := &domain.QueryExecutionEvent{
		EventID:       "gdpr-evt-" + userID + "-" + time.Now().Format("150405.000"),
		UserID:        userID,
		Timestamp:     time.Now(),
		QueryType:     "text",
		QueryText:     query,
		GeneratedSQL:  sql,
		ConnectorName: "test-connector",
		Status:        "success",
		RowCount:      5,
		ExecutionTime: 100,
		IPAddress:     ip,
		UserAgent:     ua,
		ResultPreview: map[string]interface{}{"rows": 5},
	}
	if err := pub.Publish(context.Background(), "query.executed", evt); err != nil {
		t.Fatalf("publish: %v", err)
	}
}

func TestGDPR_DeleteUserData(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	handler, publisher, _ := setupGDPRTestStack(t, db)

	userID := "gdpr-test-delete-001"

	// Publish some events for this user
	publishGDPREvent(t, publisher, userID, "Show total revenue", "SELECT SUM(revenue) FROM sales", "10.0.0.1", "TestBrowser/1.0")
	publishGDPREvent(t, publisher, userID, "List all customers", "SELECT * FROM customers", "10.0.0.1", "TestBrowser/1.0")
	publishGDPREvent(t, publisher, userID, "Monthly sales report", "SELECT month, SUM(amount) FROM orders GROUP BY month", "10.0.0.1", "TestBrowser/1.0")

	// Call DELETE /api/gdpr/user-data
	req := httptest.NewRequest("DELETE", "/api/gdpr/user-data", nil)
	rec := httptest.NewRecorder()
	injectUserContext(http.HandlerFunc(handler.HandleDeleteUserData), userID, "10.0.0.1", "TestBrowser/1.0").ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}

	// Verify response structure
	if resp["success"] != true {
		t.Errorf("Expected success=true, got %v", resp["success"])
	}

	deletedCount, ok := resp["records_deleted"].(float64)
	if !ok || deletedCount != 3 {
		t.Errorf("Expected records_deleted=3, got %v", resp["records_deleted"])
	}

	if resp["message"] == nil || resp["message"] == "" {
		t.Error("Expected non-empty message")
	}

	if resp["timestamp"] == nil {
		t.Error("Expected timestamp in response")
	}

	// Verify compliance metadata
	compliance, ok := resp["compliance"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected compliance object in response")
	}
	if compliance["regulation"] != "GDPR Article 17" {
		t.Errorf("Expected regulation 'GDPR Article 17', got %v", compliance["regulation"])
	}
	if compliance["right"] != "Right to Erasure" {
		t.Errorf("Expected right 'Right to Erasure', got %v", compliance["right"])
	}
	if compliance["status"] != "completed" {
		t.Errorf("Expected status 'completed', got %v", compliance["status"])
	}

	t.Logf("GDPR delete test passed - deleted %v records", deletedCount)
}

func TestGDPR_ExportUserData_JSON(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	handler, publisher, _ := setupGDPRTestStack(t, db)

	userID := "gdpr-test-export-json-001"

	// Publish events
	publishGDPREvent(t, publisher, userID, "Show revenue", "SELECT SUM(revenue) FROM sales", "10.0.0.1", "TestBrowser/1.0")
	publishGDPREvent(t, publisher, userID, "List products", "SELECT * FROM products", "10.0.0.1", "TestBrowser/1.0")

	// Call GET /api/gdpr/export?format=json
	req := httptest.NewRequest("GET", "/api/gdpr/export?format=json", nil)
	rec := httptest.NewRecorder()
	injectUserContext(http.HandlerFunc(handler.HandleExportUserData), userID, "10.0.0.1", "TestBrowser/1.0").ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	// Verify Content-Type
	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	// Verify Content-Disposition (download attachment)
	contentDisposition := rec.Header().Get("Content-Disposition")
	if contentDisposition == "" {
		t.Error("Expected Content-Disposition header for file download")
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}

	// Verify export_info
	exportInfo, ok := resp["export_info"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected export_info in response")
	}
	if exportInfo["user_id"] != userID {
		t.Errorf("Expected user_id %s, got %v", userID, exportInfo["user_id"])
	}
	if exportInfo["format"] != "json" {
		t.Errorf("Expected format 'json', got %v", exportInfo["format"])
	}
	if exportInfo["regulation"] != "GDPR Article 15" {
		t.Errorf("Expected regulation 'GDPR Article 15', got %v", exportInfo["regulation"])
	}
	if exportInfo["right"] != "Right to Access" {
		t.Errorf("Expected right 'Right to Access', got %v", exportInfo["right"])
	}

	// Verify query_history data
	history, ok := resp["query_history"].([]interface{})
	if !ok {
		t.Fatal("Expected query_history array in response")
	}
	if len(history) != 2 {
		t.Errorf("Expected 2 query history items, got %d", len(history))
	}

	// Verify metadata
	metadata, ok := resp["metadata"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected metadata in response")
	}
	totalQueries, ok := metadata["total_queries"].(float64)
	if !ok || totalQueries != 2 {
		t.Errorf("Expected total_queries=2, got %v", metadata["total_queries"])
	}

	t.Logf("GDPR JSON export test passed - exported %d records", len(history))
}

func TestGDPR_ExportUserData_CSV(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	handler, publisher, _ := setupGDPRTestStack(t, db)

	userID := "gdpr-test-export-csv-001"

	// Publish events
	publishGDPREvent(t, publisher, userID, "Show revenue", "SELECT SUM(revenue) FROM sales", "10.0.0.1", "TestBrowser/1.0")

	// Call GET /api/gdpr/export?format=csv
	req := httptest.NewRequest("GET", "/api/gdpr/export?format=csv", nil)
	rec := httptest.NewRecorder()
	injectUserContext(http.HandlerFunc(handler.HandleExportUserData), userID, "10.0.0.1", "TestBrowser/1.0").ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	// Verify Content-Type
	contentType := rec.Header().Get("Content-Type")
	if contentType != "text/csv" {
		t.Errorf("Expected Content-Type text/csv, got %s", contentType)
	}

	// Verify Content-Disposition
	contentDisposition := rec.Header().Get("Content-Disposition")
	if contentDisposition == "" {
		t.Error("Expected Content-Disposition header for CSV download")
	}

	// Verify CSV header row exists
	body := rec.Body.String()
	if body == "" {
		t.Error("Expected non-empty CSV body")
	}

	t.Logf("GDPR CSV export test passed - got %d bytes", len(body))
}

func TestGDPR_ExportUserData_UnsupportedFormat(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	handler, publisher, _ := setupGDPRTestStack(t, db)

	userID := "gdpr-test-export-bad-001"

	publishGDPREvent(t, publisher, userID, "Show revenue", "SELECT 1", "10.0.0.1", "TestBrowser/1.0")

	// Call GET /api/gdpr/export?format=xml (unsupported)
	req := httptest.NewRequest("GET", "/api/gdpr/export?format=xml", nil)
	rec := httptest.NewRecorder()
	injectUserContext(http.HandlerFunc(handler.HandleExportUserData), userID, "10.0.0.1", "TestBrowser/1.0").ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for unsupported format, got %d", rec.Code)
	}

	t.Log("GDPR unsupported format test passed")
}

func TestGDPR_DataPortability(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	handler, publisher, _ := setupGDPRTestStack(t, db)

	userID := "gdpr-test-portable-001"

	// Publish events
	publishGDPREvent(t, publisher, userID, "Show orders", "SELECT * FROM orders", "10.0.0.1", "TestBrowser/1.0")
	publishGDPREvent(t, publisher, userID, "List users", "SELECT * FROM users", "10.0.0.1", "TestBrowser/1.0")

	// Call GET /api/gdpr/portable-data
	req := httptest.NewRequest("GET", "/api/gdpr/portable-data", nil)
	rec := httptest.NewRecorder()
	injectUserContext(http.HandlerFunc(handler.HandleDataPortabilityRequest), userID, "10.0.0.1", "TestBrowser/1.0").ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	// Verify Content-Type
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", ct)
	}

	// Verify Content-Disposition (download)
	if cd := rec.Header().Get("Content-Disposition"); cd == "" {
		t.Error("Expected Content-Disposition header for portable data download")
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}

	// Verify metadata
	metadata, ok := resp["metadata"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected metadata in portable data response")
	}
	if metadata["standard"] != "GDPR Article 20 - Data Portability" {
		t.Errorf("Expected standard 'GDPR Article 20 - Data Portability', got %v", metadata["standard"])
	}
	if metadata["user_id"] != userID {
		t.Errorf("Expected user_id %s, got %v", userID, metadata["user_id"])
	}

	// Verify data section
	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected data section in portable data response")
	}
	history, ok := data["query_history"].([]interface{})
	if !ok {
		t.Fatal("Expected query_history in data section")
	}
	if len(history) != 2 {
		t.Errorf("Expected 2 query history items, got %d", len(history))
	}

	// Verify data_types in metadata
	dataTypes, ok := metadata["data_types"].([]interface{})
	if !ok {
		t.Fatal("Expected data_types array in metadata")
	}
	if len(dataTypes) < 1 {
		t.Error("Expected at least 1 data type")
	}

	t.Logf("GDPR data portability test passed - %d records", len(history))
}

func TestGDPR_ComplianceReport(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	handler, publisher, _ := setupGDPRTestStack(t, db)

	userID := "gdpr-test-compliance-001"

	// Publish an event so stats have data
	publishGDPREvent(t, publisher, userID, "Show revenue", "SELECT SUM(revenue) FROM sales", "10.0.0.1", "TestBrowser/1.0")

	// Call GET /api/gdpr/compliance-report
	req := httptest.NewRequest("GET", "/api/gdpr/compliance-report", nil)
	rec := httptest.NewRecorder()
	injectUserContext(http.HandlerFunc(handler.HandleComplianceReport), userID, "10.0.0.1", "TestBrowser/1.0").ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}

	// Verify user_id
	if resp["user_id"] != userID {
		t.Errorf("Expected user_id %s, got %v", userID, resp["user_id"])
	}

	// Verify generated_at timestamp
	if resp["generated_at"] == nil {
		t.Error("Expected generated_at timestamp")
	}

	// Verify GDPR rights
	gdprRights, ok := resp["gdpr_rights"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected gdpr_rights in compliance report")
	}

	// Right to Access
	rightToAccess, ok := gdprRights["right_to_access"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected right_to_access in gdpr_rights")
	}
	if rightToAccess["available"] != true {
		t.Error("Expected right_to_access.available=true")
	}

	// Right to Erasure
	rightToErasure, ok := gdprRights["right_to_erasure"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected right_to_erasure in gdpr_rights")
	}
	if rightToErasure["available"] != true {
		t.Error("Expected right_to_erasure.available=true")
	}

	// Right to Portability
	rightToPortability, ok := gdprRights["right_to_portability"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected right_to_portability in gdpr_rights")
	}
	if rightToPortability["available"] != true {
		t.Error("Expected right_to_portability.available=true")
	}

	// Verify data_retention
	retention, ok := resp["data_retention"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected data_retention in compliance report")
	}
	if retention["policy"] != "90 days" {
		t.Errorf("Expected retention policy '90 days', got %v", retention["policy"])
	}
	if retention["automatic_deletion"] != true {
		t.Error("Expected automatic_deletion=true")
	}

	// Verify security_measures
	secMeasures, ok := resp["security_measures"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected security_measures in compliance report")
	}
	if secMeasures["pii_redaction"] != "Automatic" {
		t.Errorf("Expected pii_redaction 'Automatic', got %v", secMeasures["pii_redaction"])
	}

	// Verify data_summary (stats)
	if resp["data_summary"] == nil {
		t.Error("Expected data_summary in compliance report")
	}

	t.Log("GDPR compliance report test passed")
}

func TestGDPR_Unauthorized(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	handler, _, _ := setupGDPRTestStack(t, db)

	tests := []struct {
		name       string
		method     string
		path       string
		handlerFn  http.HandlerFunc
	}{
		{"Delete without auth", "DELETE", "/api/gdpr/user-data", handler.HandleDeleteUserData},
		{"Export without auth", "GET", "/api/gdpr/export?format=json", handler.HandleExportUserData},
		{"Portability without auth", "GET", "/api/gdpr/portable-data", handler.HandleDataPortabilityRequest},
		{"Compliance report without auth", "GET", "/api/gdpr/compliance-report", handler.HandleComplianceReport},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Request without user context - no injectUserContext wrapper
			req := httptest.NewRequest(tc.method, tc.path, nil)
			rec := httptest.NewRecorder()
			tc.handlerFn(rec, req)

			if rec.Code != http.StatusUnauthorized {
				t.Errorf("Expected 401, got %d: %s", rec.Code, rec.Body.String())
			}
		})
	}

	t.Log("GDPR unauthorized test passed - all endpoints require auth")
}

func TestGDPR_UserIsolation(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	handler, publisher, _ := setupGDPRTestStack(t, db)

	userA := "gdpr-test-iso-A"
	userB := "gdpr-test-iso-B"

	// Publish events for both users
	publishGDPREvent(t, publisher, userA, "User A query 1", "SELECT 1", "10.0.0.1", "UA-A")
	publishGDPREvent(t, publisher, userA, "User A query 2", "SELECT 2", "10.0.0.1", "UA-A")
	publishGDPREvent(t, publisher, userB, "User B query 1", "SELECT 3", "10.0.0.2", "UA-B")
	publishGDPREvent(t, publisher, userB, "User B query 2", "SELECT 4", "10.0.0.2", "UA-B")
	publishGDPREvent(t, publisher, userB, "User B query 3", "SELECT 5", "10.0.0.2", "UA-B")

	// User A exports their data - should only see 2 records
	req := httptest.NewRequest("GET", "/api/gdpr/export?format=json", nil)
	rec := httptest.NewRecorder()
	injectUserContext(http.HandlerFunc(handler.HandleExportUserData), userA, "10.0.0.1", "UA-A").ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("User A export expected 200, got %d", rec.Code)
	}

	var respA map[string]interface{}
	json.NewDecoder(rec.Body).Decode(&respA)

	historyA, ok := respA["query_history"].([]interface{})
	if !ok {
		t.Fatal("Expected query_history for User A")
	}
	if len(historyA) != 2 {
		t.Errorf("User A should see 2 records, got %d", len(historyA))
	}

	// User A deletes their data
	delReq := httptest.NewRequest("DELETE", "/api/gdpr/user-data", nil)
	delRec := httptest.NewRecorder()
	injectUserContext(http.HandlerFunc(handler.HandleDeleteUserData), userA, "10.0.0.1", "UA-A").ServeHTTP(delRec, delReq)

	if delRec.Code != http.StatusOK {
		t.Fatalf("User A delete expected 200, got %d", delRec.Code)
	}

	var delResp map[string]interface{}
	json.NewDecoder(delRec.Body).Decode(&delResp)

	deletedCount, _ := delResp["records_deleted"].(float64)
	if deletedCount != 2 {
		t.Errorf("User A should have deleted 2 records, deleted %v", deletedCount)
	}

	// User B's data should be unaffected - export should still show 3 records
	reqB := httptest.NewRequest("GET", "/api/gdpr/export?format=json", nil)
	recB := httptest.NewRecorder()
	injectUserContext(http.HandlerFunc(handler.HandleExportUserData), userB, "10.0.0.2", "UA-B").ServeHTTP(recB, reqB)

	if recB.Code != http.StatusOK {
		t.Fatalf("User B export expected 200, got %d", recB.Code)
	}

	var respB map[string]interface{}
	json.NewDecoder(recB.Body).Decode(&respB)

	historyB, ok := respB["query_history"].([]interface{})
	if !ok {
		t.Fatal("Expected query_history for User B")
	}
	if len(historyB) != 3 {
		t.Errorf("User B should still see 3 records after User A's deletion, got %d", len(historyB))
	}

	// Verify User A has no data left
	reqAAfter := httptest.NewRequest("GET", "/api/gdpr/export?format=json", nil)
	recAAfter := httptest.NewRecorder()
	injectUserContext(http.HandlerFunc(handler.HandleExportUserData), userA, "10.0.0.1", "UA-A").ServeHTTP(recAAfter, reqAAfter)

	var respAAfter map[string]interface{}
	json.NewDecoder(recAAfter.Body).Decode(&respAAfter)

	historyAAfter, ok := respAAfter["query_history"].([]interface{})
	if !ok {
		// query_history could be null/nil when no records exist
		historyAAfter = []interface{}{}
	}
	if len(historyAAfter) != 0 {
		t.Errorf("User A should see 0 records after deletion, got %d", len(historyAAfter))
	}

	t.Log("GDPR user isolation test passed - User A's deletion did not affect User B")
}

func TestGDPR_DeleteThenExport(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	handler, publisher, _ := setupGDPRTestStack(t, db)

	userID := "gdpr-test-del-export-001"

	// Publish events
	publishGDPREvent(t, publisher, userID, "Show data", "SELECT 1", "10.0.0.1", "TestBrowser/1.0")

	// Delete user data
	delReq := httptest.NewRequest("DELETE", "/api/gdpr/user-data", nil)
	delRec := httptest.NewRecorder()
	injectUserContext(http.HandlerFunc(handler.HandleDeleteUserData), userID, "10.0.0.1", "TestBrowser/1.0").ServeHTTP(delRec, delReq)

	if delRec.Code != http.StatusOK {
		t.Fatalf("Delete expected 200, got %d", delRec.Code)
	}

	// Export after deletion - should return empty data
	req := httptest.NewRequest("GET", "/api/gdpr/export?format=json", nil)
	rec := httptest.NewRecorder()
	injectUserContext(http.HandlerFunc(handler.HandleExportUserData), userID, "10.0.0.1", "TestBrowser/1.0").ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Export expected 200, got %d", rec.Code)
	}

	var resp map[string]interface{}
	json.NewDecoder(rec.Body).Decode(&resp)

	history, ok := resp["query_history"].([]interface{})
	if !ok {
		// query_history could be null/nil for empty results
		history = []interface{}{}
	}
	if len(history) != 0 {
		t.Errorf("Expected 0 records after deletion, got %d", len(history))
	}

	t.Log("GDPR delete-then-export test passed - no data remains after erasure")
}

func TestGDPR_ExportDefaultFormat(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	handler, publisher, _ := setupGDPRTestStack(t, db)

	userID := "gdpr-test-default-fmt-001"

	publishGDPREvent(t, publisher, userID, "Show data", "SELECT 1", "10.0.0.1", "TestBrowser/1.0")

	// Call export without format parameter - should default to JSON
	req := httptest.NewRequest("GET", "/api/gdpr/export", nil)
	rec := httptest.NewRecorder()
	injectUserContext(http.HandlerFunc(handler.HandleExportUserData), userID, "10.0.0.1", "TestBrowser/1.0").ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	// Should be JSON by default
	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected default format to be JSON, got Content-Type: %s", contentType)
	}

	t.Log("GDPR default format test passed - defaults to JSON")
}
