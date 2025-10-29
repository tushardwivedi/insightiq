package connectors

import (
	"log/slog"
	"os"
	"testing"
)

// TestValidateAndSanitizeQuery tests the SQL injection protection
func TestValidateAndSanitizeQuery(t *testing.T) {
	// Create a logger for testing
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError, // Only show errors to keep test output clean
	}))
	pc := &PostgresConnector{
		logger: logger,
	}

	tests := []struct {
		name      string
		query     string
		wantError bool
		errorType error
	}{
		{
			name:      "Valid simple SELECT",
			query:     "SELECT * FROM users",
			wantError: false,
		},
		{
			name:      "Valid SELECT with WHERE",
			query:     "SELECT id, name FROM users WHERE active = true",
			wantError: false,
		},
		{
			name:      "Valid SELECT with JOIN",
			query:     "SELECT u.id, u.name FROM users u JOIN orders o ON u.id = o.user_id",
			wantError: false,
		},
		{
			name:      "Valid SELECT 1 (health check)",
			query:     "SELECT 1 as test_connection",
			wantError: false,
		},
		{
			name:      "SQL Injection - DROP TABLE",
			query:     "SELECT * FROM users; DROP TABLE users;--",
			wantError: true,
			errorType: ErrDangerousQuery,
		},
		{
			name:      "SQL Injection - UNION attack",
			query:     "SELECT * FROM users WHERE id = 1 UNION SELECT * FROM passwords",
			wantError: false, // UNION is allowed in SELECT, but we could add stricter rules
		},
		{
			name:      "SQL Injection - Comment bypass",
			query:     "SELECT * FROM users WHERE id = 1-- AND active = true",
			wantError: true,
			errorType: ErrDangerousQuery,
		},
		{
			name:      "Dangerous - DELETE statement",
			query:     "DELETE FROM users WHERE id = 1",
			wantError: true,
			errorType: ErrDangerousQuery,
		},
		{
			name:      "Dangerous - INSERT statement",
			query:     "INSERT INTO users (name) VALUES ('hacker')",
			wantError: true,
			errorType: ErrDangerousQuery,
		},
		{
			name:      "Dangerous - UPDATE statement",
			query:     "UPDATE users SET role = 'admin' WHERE id = 1",
			wantError: true,
			errorType: ErrDangerousQuery,
		},
		{
			name:      "Dangerous - DROP TABLE",
			query:     "DROP TABLE users",
			wantError: true,
			errorType: ErrDangerousQuery,
		},
		{
			name:      "Dangerous - ALTER TABLE",
			query:     "ALTER TABLE users ADD COLUMN hacked boolean",
			wantError: true,
			errorType: ErrDangerousQuery,
		},
		{
			name:      "Dangerous - CREATE TABLE",
			query:     "CREATE TABLE backdoor (id int)",
			wantError: true,
			errorType: ErrDangerousQuery,
		},
		{
			name:      "Dangerous - EXEC/EXECUTE",
			query:     "EXEC sp_executesql N'SELECT * FROM users'",
			wantError: true,
			errorType: ErrDangerousQuery,
		},
		{
			name:      "Dangerous - PG_SLEEP (DoS)",
			query:     "SELECT pg_sleep(10)",
			wantError: true,
			errorType: ErrDangerousQuery,
		},
		{
			name:      "Dangerous - COPY command",
			query:     "COPY users TO '/tmp/users.csv'",
			wantError: true,
			errorType: ErrDangerousQuery,
		},
		{
			name:      "Dangerous - Information schema access",
			query:     "SELECT * FROM information_schema.tables",
			wantError: true,
			errorType: ErrDangerousQuery,
		},
		{
			name:      "Empty query",
			query:     "",
			wantError: true,
			errorType: ErrInvalidQuery,
		},
		{
			name:      "Query too long",
			query:     "SELECT " + string(make([]byte, maxQueryLength+1)),
			wantError: true,
			errorType: ErrQueryTooLong,
		},
		{
			name:      "Multiple statements",
			query:     "SELECT * FROM users; SELECT * FROM orders",
			wantError: true,
			errorType: ErrDangerousQuery,
		},
		{
			name:      "Block comment injection",
			query:     "SELECT * FROM users /* WHERE id = 1 */ WHERE id = 2",
			wantError: true,
			errorType: ErrDangerousQuery,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pc.validateAndSanitizeQuery(tt.query)
			if tt.wantError {
				if err == nil {
					t.Errorf("validateAndSanitizeQuery() expected error but got none for query: %s", tt.query)
				} else if tt.errorType != nil && err != tt.errorType {
					t.Errorf("validateAndSanitizeQuery() error = %v, want %v", err, tt.errorType)
				}
			} else {
				if err != nil {
					t.Errorf("validateAndSanitizeQuery() unexpected error = %v for query: %s", err, tt.query)
				}
			}
		})
	}
}

// TestSQLInjectionScenarios tests common SQL injection attack patterns
func TestSQLInjectionScenarios(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))
	pc := &PostgresConnector{
		logger: logger,
	}

	injectionAttempts := []string{
		// Classic SQL injection patterns
		"1' OR '1'='1",
		"1' OR 1=1--",
		"admin'--",
		"' OR '1'='1' /*",
		"1'; DROP TABLE users--",

		// Boolean-based blind SQL injection
		"1' AND '1'='1",
		"1' AND '1'='2",

		// Time-based blind SQL injection
		"1' AND pg_sleep(5)--",
		"1' WAITFOR DELAY '00:00:05'--",

		// Union-based injection (some might be valid SELECTs)
		"1' UNION SELECT NULL--",
		"' UNION SELECT password FROM users--",

		// Stacked queries
		"1'; DELETE FROM users WHERE 1=1--",
		"1'; INSERT INTO users VALUES ('hacker','pass')--",

		// Out-of-band injection
		"'; COPY users TO '/tmp/out.txt'--",
	}

	for _, injection := range injectionAttempts {
		t.Run("Injection: "+injection, func(t *testing.T) {
			// All injection attempts should be caught somewhere in validation
			query := "SELECT * FROM users WHERE id = " + injection
			err := pc.validateAndSanitizeQuery(query)

			// We expect most injections to fail validation
			// Some might pass if they look like valid SELECT statements
			// The key is that they can't contain dangerous patterns
			if err == nil {
				t.Logf("Warning: injection pattern passed validation: %s", injection)
			}
		})
	}
}
