package http

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"insightiq/backend/internal/repository"
	"insightiq/backend/internal/services"
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

func setupConnectorTest(t *testing.T, db *sqlx.DB) (*ConnectorHandlers, *repository.ConnectorRepository) {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))

	repo := repository.NewConnectorRepository(db)
	if err := repo.CreateTables(context.Background()); err != nil {
		t.Fatalf("CreateTables: %v", err)
	}

	svc := services.NewConnectorService(repo, logger)
	server := &Server{logger: logger}
	handlers := NewConnectorHandlers(svc, server)
	return handlers, repo
}

func TestDeleteConnector_ReturnsJSON(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	handlers, repo := setupConnectorTest(t, db)

	// Insert a test connector directly
	ctx := context.Background()
	_, err := db.ExecContext(ctx, `
		INSERT INTO data_connectors (id, name, type, status, config)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO NOTHING`,
		"test-delete-001", "Test Delete Connector", "superset", "disconnected",
		`{"url":"http://test.local"}`,
	)
	if err != nil {
		t.Fatalf("insert test connector: %v", err)
	}

	// Send DELETE request
	req := httptest.NewRequest(http.MethodDelete, "/api/connectors/test-delete-001", nil)
	rec := httptest.NewRecorder()
	handlers.handleDeleteConnector(rec, req)

	// Verify HTTP 200 with JSON body (not 204 No Content)
	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", ct)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("Response must be valid JSON (was causing frontend crash): %v", err)
	}

	if resp["success"] != true {
		t.Errorf("Expected success=true, got %v", resp["success"])
	}

	// Verify connector is actually deleted from DB
	connector, err := repo.GetByID(ctx, "test-delete-001")
	if err != nil {
		t.Fatalf("GetByID error: %v", err)
	}
	if connector != nil {
		t.Error("Connector should be deleted from database")
	}

	t.Log("Delete connector returns valid JSON response")
}

func TestDeleteConnector_NotFound(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	handlers, _ := setupConnectorTest(t, db)

	req := httptest.NewRequest(http.MethodDelete, "/api/connectors/nonexistent-id-999", nil)
	rec := httptest.NewRecorder()
	handlers.handleDeleteConnector(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 for non-existent connector, got %d", rec.Code)
	}

	t.Log("Delete non-existent connector returns 500")
}

func TestDeleteConnector_InvalidID(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	handlers, _ := setupConnectorTest(t, db)

	// Path with no ID segment
	req := httptest.NewRequest(http.MethodDelete, "/api/connectors/", nil)
	rec := httptest.NewRecorder()
	handlers.handleDeleteConnector(rec, req)

	// extractConnectorID returns "" for paths without a proper ID
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for empty connector ID, got %d", rec.Code)
	}

	t.Log("Delete with empty ID returns 400")
}

func TestDeleteConnector_WrongMethod(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	handlers, _ := setupConnectorTest(t, db)

	req := httptest.NewRequest(http.MethodGet, "/api/connectors/test-id", nil)
	rec := httptest.NewRecorder()
	handlers.handleDeleteConnector(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405 for GET, got %d", rec.Code)
	}

	t.Log("Delete with wrong method returns 405")
}
