package events

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"insightiq/backend/internal/domain"
	"insightiq/backend/internal/repository"
	"insightiq/backend/internal/security"
	"insightiq/backend/internal/usecases"
)

// Integration tests for end-to-end query history event flow with PII redaction
// These tests require a running PostgreSQL instance

func getTestDB(t *testing.T) *sqlx.DB {
	// Get database URL from environment or use default
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://insightiq_user:insightiq_password@localhost:5432/insightiq?sslmode=disable"
	}

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		t.Skipf("Skipping integration test: cannot connect to database: %v", err)
	}

	return db
}

func setupTestDB(t *testing.T, db *sqlx.DB) *repository.QueryHistoryRepository {
	repo := repository.NewQueryHistoryRepository(db)

	// Create tables
	ctx := context.Background()
	if err := repo.CreateTables(ctx); err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}

	// Clean up existing test data
	_, err := db.Exec("DELETE FROM query_history WHERE user_id LIKE 'test-user-%'")
	if err != nil {
		t.Fatalf("Failed to clean up test data: %v", err)
	}

	return repo
}

func TestIntegration_EndToEndQueryHistoryWithPII(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	repo := setupTestDB(t, db)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	// Initialize PII redactor
	redactor := security.NewDefaultPIIRedactor(logger)

	// Initialize use case
	useCase := usecases.NewQueryHistoryUseCase(repo, redactor, logger)

	// Initialize event publisher and handler
	publisher := NewSyncEventPublisher(logger)
	handler := NewQueryHistoryEventHandler(useCase, logger)
	publisher.Subscribe("query.executed", handler)

	ctx := context.Background()

	// Test 1: Query with email PII
	t.Run("EmailRedaction", func(t *testing.T) {
		userID := "test-user-email-001"
		queryWithEmail := "Show me orders for john.doe@example.com"

		event := &domain.QueryExecutionEvent{
			EventID:       "evt-email-001",
			UserID:        userID,
			Timestamp:     time.Now(),
			QueryType:     "analytics",
			QueryText:     queryWithEmail,
			GeneratedSQL:  "SELECT * FROM orders WHERE email = 'john.doe@example.com'",
			ConnectorID:   "conn-001",
			ConnectorName: "Test Connector",
			RowCount:      5,
			ExecutionTime: 150,
			Status:        "success",
			IPAddress:     "192.168.1.100",
			UserAgent:     "Mozilla/5.0 Test Browser",
			ResultPreview: map[string]interface{}{"sample": "data"},
		}

		// Publish event
		err := publisher.Publish(ctx, "query.executed", event)
		if err != nil {
			t.Fatalf("Failed to publish event: %v", err)
		}

		// Wait for async processing
		time.Sleep(100 * time.Millisecond)

		// Retrieve from database
		items, err := repo.GetByUserID(ctx, userID, 10, 0)
		if err != nil {
			t.Fatalf("Failed to retrieve query history: %v", err)
		}

		if len(items) != 1 {
			t.Fatalf("Expected 1 query history item, got %d", len(items))
		}

		// Verify email was redacted
		if strings.Contains(items[0].QueryText, "john.doe@example.com") {
			t.Errorf("Email should have been redacted in query_text")
		}

		if !strings.Contains(items[0].QueryText, "[EMAIL_REDACTED]") {
			t.Errorf("Expected [EMAIL_REDACTED] in query_text, got: %s", items[0].QueryText)
		}

		// Verify full details
		details, err := repo.GetByID(ctx, items[0].ID, userID)
		if err != nil {
			t.Fatalf("Failed to get query details: %v", err)
		}

		if strings.Contains(details.GeneratedSQL, "john.doe@example.com") {
			t.Errorf("Email should have been redacted in generated_sql")
		}

		if !details.DataRedacted {
			t.Errorf("Expected data_redacted flag to be true")
		}

		if details.IPAddress != "192.168.1.100" {
			t.Errorf("Expected IP address 192.168.1.100, got %s", details.IPAddress)
		}

		if details.UserAgent != "Mozilla/5.0 Test Browser" {
			t.Errorf("Expected User-Agent to be stored")
		}

		t.Logf("✅ Email redaction test passed - Stored: %s", items[0].QueryText)
	})

	// Test 2: Query with phone number PII
	t.Run("PhoneRedaction", func(t *testing.T) {
		userID := "test-user-phone-001"
		queryWithPhone := "Find customer with phone (555) 123-4567"

		event := &domain.QueryExecutionEvent{
			EventID:       "evt-phone-001",
			UserID:        userID,
			Timestamp:     time.Now(),
			QueryType:     "analytics",
			QueryText:     queryWithPhone,
			GeneratedSQL:  "SELECT * FROM customers WHERE phone = '555-123-4567'",
			ConnectorID:   "conn-001",
			ConnectorName: "Test Connector",
			RowCount:      1,
			ExecutionTime: 50,
			Status:        "success",
			IPAddress:     "10.0.0.5",
			UserAgent:     "Test Client 1.0",
			ResultPreview: map[string]interface{}{"data": "test"},
		}

		err := publisher.Publish(ctx, "query.executed", event)
		if err != nil {
			t.Fatalf("Failed to publish event: %v", err)
		}

		time.Sleep(100 * time.Millisecond)

		items, err := repo.GetByUserID(ctx, userID, 10, 0)
		if err != nil {
			t.Fatalf("Failed to retrieve query history: %v", err)
		}

		if len(items) != 1 {
			t.Fatalf("Expected 1 query history item, got %d", len(items))
		}

		if strings.Contains(items[0].QueryText, "555") || strings.Contains(items[0].QueryText, "123-4567") {
			t.Errorf("Phone number should have been redacted")
		}

		if !strings.Contains(items[0].QueryText, "[PHONE_REDACTED]") {
			t.Errorf("Expected [PHONE_REDACTED] in query_text")
		}

		details, err := repo.GetByID(ctx, items[0].ID, userID)
		if err != nil {
			t.Fatalf("Failed to get query details: %v", err)
		}

		if !details.DataRedacted {
			t.Errorf("Expected data_redacted flag to be true")
		}

		t.Logf("✅ Phone redaction test passed - Stored: %s", items[0].QueryText)
	})

	// Test 3: Query with multiple PII types
	t.Run("MultiplePIIRedaction", func(t *testing.T) {
		userID := "test-user-multi-001"
		queryWithMultiplePII := "Contact John at john@test.com or 555-123-4567. SSN: 123-45-6789"

		event := &domain.QueryExecutionEvent{
			EventID:       "evt-multi-001",
			UserID:        userID,
			Timestamp:     time.Now(),
			QueryType:     "analytics",
			QueryText:     queryWithMultiplePII,
			GeneratedSQL:  "SELECT * FROM contacts WHERE email = 'john@test.com'",
			ConnectorID:   "conn-001",
			ConnectorName: "Test Connector",
			RowCount:      1,
			ExecutionTime: 75,
			Status:        "success",
			IPAddress:     "172.16.0.1",
			UserAgent:     "Integration Test",
			ResultPreview: map[string]interface{}{"data": "test"},
		}

		err := publisher.Publish(ctx, "query.executed", event)
		if err != nil {
			t.Fatalf("Failed to publish event: %v", err)
		}

		time.Sleep(100 * time.Millisecond)

		items, err := repo.GetByUserID(ctx, userID, 10, 0)
		if err != nil {
			t.Fatalf("Failed to retrieve query history: %v", err)
		}

		if len(items) != 1 {
			t.Fatalf("Expected 1 query history item, got %d", len(items))
		}

		// Verify all PII types were redacted
		expectedRedactions := []string{
			"[EMAIL_REDACTED]",
			"[PHONE_REDACTED]",
			"[SSN_REDACTED]",
		}

		for _, expected := range expectedRedactions {
			if !strings.Contains(items[0].QueryText, expected) {
				t.Errorf("Expected %s in query_text, got: %s", expected, items[0].QueryText)
			}
		}

		// Ensure original PII is not present
		forbiddenStrings := []string{
			"john@test.com",
			"555-123-4567",
			"123-45-6789",
		}

		for _, forbidden := range forbiddenStrings {
			if strings.Contains(items[0].QueryText, forbidden) {
				t.Errorf("Original PII should not be present: %s", forbidden)
			}
		}

		details, err := repo.GetByID(ctx, items[0].ID, userID)
		if err != nil {
			t.Fatalf("Failed to get query details: %v", err)
		}

		if !details.DataRedacted {
			t.Errorf("Expected data_redacted flag to be true")
		}

		t.Logf("✅ Multiple PII redaction test passed - Stored: %s", items[0].QueryText)
	})

	// Test 4: Query without PII (should not be flagged as redacted)
	t.Run("NoPIIQuery", func(t *testing.T) {
		userID := "test-user-no-pii-001"
		queryWithoutPII := "Show me total sales for last month"

		event := &domain.QueryExecutionEvent{
			EventID:       "evt-no-pii-001",
			UserID:        userID,
			Timestamp:     time.Now(),
			QueryType:     "analytics",
			QueryText:     queryWithoutPII,
			GeneratedSQL:  "SELECT SUM(amount) FROM sales WHERE created_at > NOW() - INTERVAL '1 month'",
			ConnectorID:   "conn-001",
			ConnectorName: "Test Connector",
			RowCount:      1,
			ExecutionTime: 25,
			Status:        "success",
			IPAddress:     "192.168.1.50",
			UserAgent:     "Test Browser",
			ResultPreview: map[string]interface{}{"data": "test"},
		}

		err := publisher.Publish(ctx, "query.executed", event)
		if err != nil {
			t.Fatalf("Failed to publish event: %v", err)
		}

		time.Sleep(100 * time.Millisecond)

		items, err := repo.GetByUserID(ctx, userID, 10, 0)
		if err != nil {
			t.Fatalf("Failed to retrieve query history: %v", err)
		}

		if len(items) != 1 {
			t.Fatalf("Expected 1 query history item, got %d", len(items))
		}

		// Query text should be unchanged
		if items[0].QueryText != queryWithoutPII {
			t.Errorf("Query text should not have been modified. Expected: %s, Got: %s",
				queryWithoutPII, items[0].QueryText)
		}

		details, err := repo.GetByID(ctx, items[0].ID, userID)
		if err != nil {
			t.Fatalf("Failed to get query details: %v", err)
		}

		if details.DataRedacted {
			t.Errorf("Expected data_redacted flag to be false for query without PII")
		}

		t.Logf("✅ No PII query test passed - Stored unchanged: %s", items[0].QueryText)
	})

	// Test 5: User isolation - User A should not see User B's queries
	t.Run("UserIsolation", func(t *testing.T) {
		userA := "test-user-isolation-A"
		userB := "test-user-isolation-B"

		// Create event for User A
		eventA := &domain.QueryExecutionEvent{
			EventID:       "evt-isolation-A",
			UserID:        userA,
			Timestamp:     time.Now(),
			QueryType:     "analytics",
			QueryText:     "User A query with alice@example.com",
			GeneratedSQL:  "SELECT * FROM orders WHERE user = 'A'",
			ConnectorID:   "conn-001",
			ConnectorName: "Test Connector",
			Status:        "success",
			IPAddress:     "192.168.1.10",
			UserAgent:     "User A Browser",
			ResultPreview: map[string]interface{}{"data": "test"},
		}

		// Create event for User B
		eventB := &domain.QueryExecutionEvent{
			EventID:       "evt-isolation-B",
			UserID:        userB,
			Timestamp:     time.Now(),
			QueryType:     "analytics",
			QueryText:     "User B query with bob@example.com",
			GeneratedSQL:  "SELECT * FROM orders WHERE user = 'B'",
			ConnectorID:   "conn-001",
			ConnectorName: "Test Connector",
			Status:        "success",
			IPAddress:     "192.168.1.20",
			UserAgent:     "User B Browser",
			ResultPreview: map[string]interface{}{"data": "test"},
		}

		// Publish both events
		if err := publisher.Publish(ctx, "query.executed", eventA); err != nil {
			t.Fatalf("Failed to publish event A: %v", err)
		}
		if err := publisher.Publish(ctx, "query.executed", eventB); err != nil {
			t.Fatalf("Failed to publish event B: %v", err)
		}

		time.Sleep(100 * time.Millisecond)

		// User A should only see their own queries
		itemsA, err := repo.GetByUserID(ctx, userA, 10, 0)
		if err != nil {
			t.Fatalf("Failed to retrieve User A queries: %v", err)
		}

		if len(itemsA) != 1 {
			t.Errorf("User A should see exactly 1 query, got %d", len(itemsA))
		}

		if len(itemsA) > 0 {
			if !strings.Contains(itemsA[0].QueryText, "User A query") {
				t.Errorf("User A should only see their own queries")
			}
			if strings.Contains(itemsA[0].QueryText, "User B query") {
				t.Errorf("User A should NOT see User B's queries")
			}
		}

		// User B should only see their own queries
		itemsB, err := repo.GetByUserID(ctx, userB, 10, 0)
		if err != nil {
			t.Fatalf("Failed to retrieve User B queries: %v", err)
		}

		if len(itemsB) != 1 {
			t.Errorf("User B should see exactly 1 query, got %d", len(itemsB))
		}

		if len(itemsB) > 0 {
			if !strings.Contains(itemsB[0].QueryText, "User B query") {
				t.Errorf("User B should only see their own queries")
			}
			if strings.Contains(itemsB[0].QueryText, "User A query") {
				t.Errorf("User B should NOT see User A's queries")
			}
		}

		// User A should not be able to access User B's query by ID (repo-level)
		if len(itemsB) > 0 {
			queryBID := itemsB[0].ID
			_, err := repo.GetByID(ctx, queryBID, userA)
			if err == nil {
				t.Errorf("User A should NOT be able to access User B's query via repo")
			}
		}

		// Use-case level isolation: GetQueryDetails enforces owner check
		if len(itemsB) > 0 {
			queryBID := itemsB[0].ID
			_, err := useCase.GetQueryDetails(ctx, queryBID, userA)
			if err == nil {
				t.Errorf("Use case should reject User A accessing User B's query")
			}
		}

		// Delete isolation: User A cannot delete User B's query
		if len(itemsB) > 0 {
			queryBID := itemsB[0].ID
			err := useCase.DeleteQuery(ctx, queryBID, userA)
			if err == nil {
				t.Errorf("User A should NOT be able to delete User B's query")
			}
		}

		// Verify User B's query still exists after User A's delete attempt
		itemsBAfter, err := repo.GetByUserID(ctx, userB, 10, 0)
		if err != nil {
			t.Fatalf("Failed to retrieve User B queries after delete attempt: %v", err)
		}
		if len(itemsBAfter) != 1 {
			t.Errorf("User B's query should survive User A's delete attempt, got %d", len(itemsBAfter))
		}

		t.Logf("✅ User isolation test passed - repo, use-case, and delete levels verified")
	})

	// Test 6: IP Address and User-Agent capture
	t.Run("AuditTrailCapture", func(t *testing.T) {
		userID := "test-user-audit-001"

		testCases := []struct {
			name      string
			ipAddress string
			userAgent string
		}{
			{
				name:      "IPv4 Address",
				ipAddress: "203.0.113.42",
				userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			},
			{
				name:      "Private IP",
				ipAddress: "10.0.0.100",
				userAgent: "curl/7.68.0",
			},
			{
				name:      "Localhost",
				ipAddress: "127.0.0.1",
				userAgent: "PostmanRuntime/7.28.4",
			},
		}

		for i, tc := range testCases {
			event := &domain.QueryExecutionEvent{
				EventID:       fmt.Sprintf("evt-audit-%d", i),
				UserID:        userID,
				Timestamp:     time.Now(),
				QueryType:     "analytics",
				QueryText:     "Test query " + tc.name,
				GeneratedSQL:  "SELECT 1",
				ConnectorID:   "conn-001",
				ConnectorName: "Test Connector",
				Status:        "success",
				IPAddress:     tc.ipAddress,
				UserAgent:     tc.userAgent,
				ResultPreview: map[string]interface{}{"data": "test"},
			}

			if err := publisher.Publish(ctx, "query.executed", event); err != nil {
				t.Fatalf("Failed to publish event: %v", err)
			}
		}

		time.Sleep(200 * time.Millisecond)

		items, err := repo.GetByUserID(ctx, userID, 10, 0)
		if err != nil {
			t.Fatalf("Failed to retrieve query history: %v", err)
		}

		if len(items) != len(testCases) {
			t.Fatalf("Expected %d query history items, got %d", len(testCases), len(items))
		}

		// Verify each query has correct audit trail
		for i, item := range items {
			details, err := repo.GetByID(ctx, item.ID, userID)
			if err != nil {
				t.Fatalf("Failed to get query details: %v", err)
			}

			// Note: items are in DESC order, so reverse index
			tcIndex := len(testCases) - 1 - i
			tc := testCases[tcIndex]

			if details.IPAddress != tc.ipAddress {
				t.Errorf("Test %s: Expected IP %s, got %s", tc.name, tc.ipAddress, details.IPAddress)
			}

			if details.UserAgent != tc.userAgent {
				t.Errorf("Test %s: Expected User-Agent %s, got %s", tc.name, tc.userAgent, details.UserAgent)
			}

			t.Logf("✅ Audit trail captured for %s - IP: %s, UA: %s", tc.name, details.IPAddress, details.UserAgent)
		}
	})

	// Test 7: Credit card redaction
	t.Run("CreditCardRedaction", func(t *testing.T) {
		userID := "test-user-cc-001"
		queryWithCC := "Find transactions for card 4111-1111-1111-1111"

		event := &domain.QueryExecutionEvent{
			EventID:       "evt-cc-001",
			UserID:        userID,
			Timestamp:     time.Now(),
			QueryType:     "analytics",
			QueryText:     queryWithCC,
			GeneratedSQL:  "SELECT * FROM transactions WHERE card = '4111111111111111'",
			ConnectorID:   "conn-001",
			ConnectorName: "Test Connector",
			Status:        "success",
			IPAddress:     "192.168.1.99",
			UserAgent:     "Test Browser",
			ResultPreview: map[string]interface{}{"data": "test"},
		}

		err := publisher.Publish(ctx, "query.executed", event)
		if err != nil {
			t.Fatalf("Failed to publish event: %v", err)
		}

		time.Sleep(100 * time.Millisecond)

		items, err := repo.GetByUserID(ctx, userID, 10, 0)
		if err != nil {
			t.Fatalf("Failed to retrieve query history: %v", err)
		}

		if len(items) != 1 {
			t.Fatalf("Expected 1 query history item, got %d", len(items))
		}

		if strings.Contains(items[0].QueryText, "4111") {
			t.Errorf("Credit card number should have been redacted")
		}

		if !strings.Contains(items[0].QueryText, "[CARD_REDACTED]") {
			t.Errorf("Expected [CARD_REDACTED] in query_text, got: %s", items[0].QueryText)
		}

		details, err := repo.GetByID(ctx, items[0].ID, userID)
		if err != nil {
			t.Fatalf("Failed to get query details: %v", err)
		}

		if !details.DataRedacted {
			t.Errorf("Expected data_redacted flag to be true")
		}

		if strings.Contains(details.GeneratedSQL, "4111111111111111") {
			t.Errorf("Credit card number should have been redacted in SQL")
		}

		t.Logf("✅ Credit card redaction test passed - Stored: %s", items[0].QueryText)
	})
}

func TestIntegration_QueryHistoryStats(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	repo := setupTestDB(t, db)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	redactor := security.NewDefaultPIIRedactor(logger)
	useCase := usecases.NewQueryHistoryUseCase(repo, redactor, logger)
	publisher := NewSyncEventPublisher(logger)
	handler := NewQueryHistoryEventHandler(useCase, logger)
	publisher.Subscribe("query.executed", handler)

	ctx := context.Background()
	userID := "test-user-stats-001"

	// Create multiple query events
	events := []*domain.QueryExecutionEvent{
		{
			EventID:       "evt-stats-1",
			UserID:        userID,
			Timestamp:     time.Now(),
			QueryType:     "analytics",
			QueryText:     "Query 1",
			GeneratedSQL:  "SELECT 1",
			Status:        "success",
			RowCount:      10,
			ExecutionTime: 50,
			ResultPreview: map[string]interface{}{"data": "test"},
		},
		{
			EventID:       "evt-stats-2",
			UserID:        userID,
			Timestamp:     time.Now(),
			QueryType:     "analytics",
			QueryText:     "Query 2",
			GeneratedSQL:  "SELECT 2",
			Status:        "success",
			RowCount:      20,
			ExecutionTime: 75,
			ResultPreview: map[string]interface{}{"data": "test"},
		},
		{
			EventID:       "evt-stats-3",
			UserID:        userID,
			Timestamp:     time.Now(),
			QueryType:     "analytics",
			QueryText:     "Query 3",
			GeneratedSQL:  "SELECT 3",
			Status:        "error",
			ErrorMessage:  "Test error",
			ExecutionTime: 25,
			ResultPreview: map[string]interface{}{"data": "test"},
		},
	}

	for _, event := range events {
		if err := publisher.Publish(ctx, "query.executed", event); err != nil {
			t.Fatalf("Failed to publish event: %v", err)
		}
	}

	time.Sleep(200 * time.Millisecond)

	// Get stats
	stats, err := repo.GetStats(ctx, userID)
	if err != nil {
		t.Fatalf("Failed to get stats: %v", err)
	}

	if stats["total_queries"].(int) != 3 {
		t.Errorf("Expected 3 total queries, got %v", stats["total_queries"])
	}

	if stats["successful_queries"].(int) != 2 {
		t.Errorf("Expected 2 successful queries, got %v", stats["successful_queries"])
	}

	if stats["failed_queries"].(int) != 1 {
		t.Errorf("Expected 1 failed query, got %v", stats["failed_queries"])
	}

	if stats["total_rows_returned"].(int64) != 30 {
		t.Errorf("Expected 30 total rows, got %v", stats["total_rows_returned"])
	}

	t.Logf("✅ Stats test passed - Total: %v, Success: %v, Failed: %v",
		stats["total_queries"], stats["successful_queries"], stats["failed_queries"])
}
