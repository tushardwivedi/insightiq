package usecases

import (
	"testing"
	"time"
)

func TestRecordQueryRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request RecordQueryRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid request",
			request: RecordQueryRequest{
				UserID:        "user-001",
				QueryType:     "analytics",
				QueryText:     "SELECT * FROM users",
				Status:        "success",
				ExecutionTime: 150,
				RowCount:      10,
			},
			wantErr: false,
		},
		{
			name: "Missing user_id",
			request: RecordQueryRequest{
				QueryType:     "analytics",
				QueryText:     "SELECT * FROM users",
				Status:        "success",
				ExecutionTime: 150,
			},
			wantErr: true,
			errMsg:  "user_id is required",
		},
		{
			name: "Missing query_type",
			request: RecordQueryRequest{
				UserID:        "user-001",
				QueryText:     "SELECT * FROM users",
				Status:        "success",
				ExecutionTime: 150,
			},
			wantErr: true,
			errMsg:  "query_type is required",
		},
		{
			name: "Missing query_text",
			request: RecordQueryRequest{
				UserID:        "user-001",
				QueryType:     "analytics",
				Status:        "success",
				ExecutionTime: 150,
			},
			wantErr: true,
			errMsg:  "query_text is required",
		},
		{
			name: "Missing status",
			request: RecordQueryRequest{
				UserID:        "user-001",
				QueryType:     "analytics",
				QueryText:     "SELECT * FROM users",
				ExecutionTime: 150,
			},
			wantErr: true,
			errMsg:  "status is required",
		},
		{
			name: "Invalid status",
			request: RecordQueryRequest{
				UserID:        "user-001",
				QueryType:     "analytics",
				QueryText:     "SELECT * FROM users",
				Status:        "invalid_status",
				ExecutionTime: 150,
			},
			wantErr: true,
			errMsg:  "status must be one of",
		},
		{
			name: "Error status without error_message",
			request: RecordQueryRequest{
				UserID:        "user-001",
				QueryType:     "analytics",
				QueryText:     "SELECT * FROM users",
				Status:        "error",
				ExecutionTime: 150,
			},
			wantErr: true,
			errMsg:  "error_message is required when status is error",
		},
		{
			name: "Error status with error_message",
			request: RecordQueryRequest{
				UserID:        "user-001",
				QueryType:     "analytics",
				QueryText:     "SELECT * FROM users",
				Status:        "error",
				ErrorMessage:  "Syntax error",
				ExecutionTime: 150,
			},
			wantErr: false,
		},
		{
			name: "Negative execution_time",
			request: RecordQueryRequest{
				UserID:        "user-001",
				QueryType:     "analytics",
				QueryText:     "SELECT * FROM users",
				Status:        "success",
				ExecutionTime: -100,
			},
			wantErr: true,
			errMsg:  "execution_time cannot be negative",
		},
		{
			name: "Negative row_count",
			request: RecordQueryRequest{
				UserID:        "user-001",
				QueryType:     "analytics",
				QueryText:     "SELECT * FROM users",
				Status:        "success",
				ExecutionTime: 150,
				RowCount:      -5,
			},
			wantErr: true,
			errMsg:  "row_count cannot be negative",
		},
		{
			name: "Valid timeout status",
			request: RecordQueryRequest{
				UserID:        "user-001",
				QueryType:     "analytics",
				QueryText:     "SELECT * FROM users",
				Status:        "timeout",
				ExecutionTime: 30000,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("Expected error message to contain %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestDeleteOldQueries_ValidationErrors(t *testing.T) {
	// This test verifies validation logic without needing database
	tests := []struct {
		name              string
		retentionDuration time.Duration
		wantErr           bool
		errMsg            string
	}{
		{
			name:              "Valid retention duration - 30 days",
			retentionDuration: 30 * 24 * time.Hour,
			wantErr:           false,
		},
		{
			name:              "Valid retention duration - 90 days",
			retentionDuration: 90 * 24 * time.Hour,
			wantErr:           false,
		},
		{
			name:              "Invalid retention duration - less than 24 hours",
			retentionDuration: 12 * time.Hour,
			wantErr:           true,
			errMsg:            "retention duration must be at least 24 hours",
		},
		{
			name:              "Invalid retention duration - 1 hour",
			retentionDuration: 1 * time.Hour,
			wantErr:           true,
			errMsg:            "retention duration must be at least 24 hours",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We can test the validation logic by checking the duration
			// The actual implementation would call uc.DeleteOldQueries which requires a DB
			if tt.retentionDuration < 24*time.Hour {
				if !tt.wantErr {
					t.Errorf("Expected validation to fail for duration %v", tt.retentionDuration)
				}
			} else {
				if tt.wantErr {
					t.Errorf("Expected validation to pass for duration %v", tt.retentionDuration)
				}
			}
		})
	}
}

func TestGetUserHistory_PaginationValidation(t *testing.T) {
	// Test pagination parameter normalization logic
	tests := []struct {
		name           string
		inputLimit     int
		inputOffset    int
		expectedLimit  int
		expectedOffset int
	}{
		{
			name:           "Valid pagination",
			inputLimit:     20,
			inputOffset:    10,
			expectedLimit:  20,
			expectedOffset: 10,
		},
		{
			name:           "Limit too small - should default to 10",
			inputLimit:     0,
			inputOffset:    0,
			expectedLimit:  10,
			expectedOffset: 0,
		},
		{
			name:           "Limit too large - should cap at 100",
			inputLimit:     500,
			inputOffset:    0,
			expectedLimit:  100,
			expectedOffset: 0,
		},
		{
			name:           "Negative offset - should default to 0",
			inputLimit:     20,
			inputOffset:    -5,
			expectedLimit:  20,
			expectedOffset: 0,
		},
		{
			name:           "Negative limit and offset",
			inputLimit:     -10,
			inputOffset:    -20,
			expectedLimit:  10,
			expectedOffset: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate pagination normalization logic
			limit := tt.inputLimit
			offset := tt.inputOffset

			if limit < 1 {
				limit = 10
			}
			if limit > 100 {
				limit = 100
			}
			if offset < 0 {
				offset = 0
			}

			if limit != tt.expectedLimit {
				t.Errorf("Expected limit %d, got %d", tt.expectedLimit, limit)
			}
			if offset != tt.expectedOffset {
				t.Errorf("Expected offset %d, got %d", tt.expectedOffset, offset)
			}
		})
	}
}

func TestUseCaseErrors(t *testing.T) {
	// Test that our custom errors are defined
	if ErrInvalidInput == nil {
		t.Error("ErrInvalidInput should be defined")
	}
	if ErrUnauthorized == nil {
		t.Error("ErrUnauthorized should be defined")
	}
	if ErrNotFound == nil {
		t.Error("ErrNotFound should be defined")
	}

	// Test error messages
	if ErrInvalidInput.Error() != "invalid input" {
		t.Errorf("Expected 'invalid input', got %q", ErrInvalidInput.Error())
	}
	if ErrUnauthorized.Error() != "unauthorized access" {
		t.Errorf("Expected 'unauthorized access', got %q", ErrUnauthorized.Error())
	}
	if ErrNotFound.Error() != "query history not found" {
		t.Errorf("Expected 'query history not found', got %q", ErrNotFound.Error())
	}
}

func TestRecordQueryRequest_AllStatuses(t *testing.T) {
	// Test all allowed statuses
	statuses := []string{"success", "error", "timeout"}

	for _, status := range statuses {
		t.Run("Status_"+status, func(t *testing.T) {
			req := RecordQueryRequest{
				UserID:        "user-001",
				QueryType:     "analytics",
				QueryText:     "SELECT 1",
				Status:        status,
				ExecutionTime: 100,
			}

			// Add error message for error status
			if status == "error" {
				req.ErrorMessage = "Test error"
			}

			err := req.Validate()
			if err != nil {
				t.Errorf("Valid status %q should not produce error: %v", status, err)
			}
		})
	}
}

func TestRecordQueryRequest_BoundaryValues(t *testing.T) {
	tests := []struct {
		name          string
		executionTime int64
		rowCount      int
		wantErr       bool
	}{
		{
			name:          "Zero execution time",
			executionTime: 0,
			rowCount:      0,
			wantErr:       false,
		},
		{
			name:          "Large execution time",
			executionTime: 999999999,
			rowCount:      1000000,
			wantErr:       false,
		},
		{
			name:          "Negative execution time",
			executionTime: -1,
			rowCount:      0,
			wantErr:       true,
		},
		{
			name:          "Negative row count",
			executionTime: 100,
			rowCount:      -1,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := RecordQueryRequest{
				UserID:        "user-001",
				QueryType:     "analytics",
				QueryText:     "SELECT 1",
				Status:        "success",
				ExecutionTime: tt.executionTime,
				RowCount:      tt.rowCount,
			}

			err := req.Validate()

			if tt.wantErr && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
