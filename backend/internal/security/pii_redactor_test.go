package security

import (
	"log/slog"
	"os"
	"strings"
	"testing"
)

func TestNewDefaultPIIRedactor(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	redactor := NewDefaultPIIRedactor(logger)

	if redactor == nil {
		t.Fatal("Expected redactor to be initialized")
	}

	if len(redactor.patterns) == 0 {
		t.Fatal("Expected patterns to be initialized")
	}

	t.Logf("Initialized with %d patterns", len(redactor.patterns))
}

func TestRedactEmail(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	redactor := NewDefaultPIIRedactor(logger)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple email",
			input:    "Contact john.doe@example.com for info",
			expected: "Contact [EMAIL_REDACTED] for info",
		},
		{
			name:     "Multiple emails",
			input:    "Email alice@test.com or bob@company.org",
			expected: "Email [EMAIL_REDACTED] or [EMAIL_REDACTED]",
		},
		{
			name:     "No email",
			input:    "No email here",
			expected: "No email here",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := redactor.RedactQuery(tt.input)
			if result != tt.expected {
				t.Errorf("Expected: %s, Got: %s", tt.expected, result)
			}
		})
	}
}

func TestRedactPhone(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	redactor := NewDefaultPIIRedactor(logger)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Formatted phone",
			input:    "Call me at (555) 123-4567",
			expected: "Call me at [PHONE_REDACTED]",
		},
		{
			name:     "Phone with dashes",
			input:    "My number is 555-123-4567",
			expected: "My number is [PHONE_REDACTED]",
		},
		{
			name:     "Phone no formatting",
			input:    "Contact 5551234567",
			expected: "Contact [PHONE_REDACTED]",
		},
		{
			name:     "International format",
			input:    "Call +1-555-123-4567",
			expected: "Call [PHONE_REDACTED]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := redactor.RedactQuery(tt.input)
			if result != tt.expected {
				t.Errorf("Expected: %s, Got: %s", tt.expected, result)
			}
		})
	}
}

func TestRedactSSN(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	redactor := NewDefaultPIIRedactor(logger)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "SSN with dashes",
			input:    "SSN: 123-45-6789",
			expected: "SSN: [SSN_REDACTED]",
		},
		{
			name:     "SSN with spaces",
			input:    "Social Security 123 45 6789",
			expected: "Social Security [SSN_REDACTED]",
		},
		{
			name:     "SSN no formatting",
			input:    "ID 123456789",
			expected: "ID [SSN_REDACTED]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := redactor.RedactQuery(tt.input)
			if result != tt.expected {
				t.Errorf("Expected: %s, Got: %s", tt.expected, result)
			}
		})
	}
}

func TestRedactCreditCard(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	redactor := NewDefaultPIIRedactor(logger)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Card with spaces",
			input:    "Card: 4111 1111 1111 1111",
			expected: "Card: [CARD_REDACTED]",
		},
		{
			name:     "Card with dashes",
			input:    "Pay with 4111-1111-1111-1111",
			expected: "Pay with [CARD_REDACTED]",
		},
		{
			name:     "Card no formatting",
			input:    "Number 4111111111111111",
			expected: "Number [CARD_REDACTED]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := redactor.RedactQuery(tt.input)
			if result != tt.expected {
				t.Errorf("Expected: %s, Got: %s", tt.expected, result)
			}
		})
	}
}

func TestRedactSQL(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	redactor := NewDefaultPIIRedactor(logger)

	tests := []struct {
		name     string
		input    string
		contains string
	}{
		{
			name:     "SQL with email",
			input:    "SELECT * FROM users WHERE email = 'john@example.com'",
			contains: "[EMAIL_REDACTED]",
		},
		{
			name:     "SQL with phone",
			input:    "INSERT INTO contacts (phone) VALUES ('555-123-4567')",
			contains: "[PHONE_REDACTED]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := redactor.RedactSQL(tt.input)
			if !strings.Contains(result, tt.contains) {
				t.Errorf("Expected result to contain: %s, Got: %s", tt.contains, result)
			}
		})
	}
}

func TestRedactMultiplePII(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	redactor := NewDefaultPIIRedactor(logger)

	input := "Contact John at john.doe@example.com or 555-123-4567. SSN: 123-45-6789, Card: 4111-1111-1111-1111"

	result := redactor.RedactQuery(input)

	// Check that all PII types were redacted
	expectedRedactions := []string{
		"[EMAIL_REDACTED]",
		"[PHONE_REDACTED]",
		"[SSN_REDACTED]",
		"[CARD_REDACTED]",
	}

	for _, expected := range expectedRedactions {
		if !strings.Contains(result, expected) {
			t.Errorf("Expected result to contain %s", expected)
		}
	}

	// Ensure original PII is not in result
	originalPII := []string{
		"john.doe@example.com",
		"555-123-4567",
		"123-45-6789",
		"4111-1111-1111-1111",
	}

	for _, pii := range originalPII {
		if strings.Contains(result, pii) {
			t.Errorf("Result should not contain original PII: %s", pii)
		}
	}
}

func TestDetectSensitivityLevel(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	redactor := NewDefaultPIIRedactor(logger)

	tests := []struct {
		name     string
		input    string
		expected SensitivityLevel
	}{
		{
			name:     "No PII",
			input:    "Show me the sales data",
			expected: SensitivityLow,
		},
		{
			name:     "Email only",
			input:    "Contact john@example.com",
			expected: SensitivityMedium,
		},
		{
			name:     "Email and phone",
			input:    "John: john@example.com, 555-123-4567",
			expected: SensitivityHigh,
		},
		{
			name:     "SSN present",
			input:    "SSN: 123-45-6789",
			expected: SensitivityCritical,
		},
		{
			name:     "Credit card present",
			input:    "Card: 4111-1111-1111-1111",
			expected: SensitivityCritical,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := redactor.DetectSensitivityLevel(tt.input)
			if result != tt.expected {
				t.Errorf("Expected sensitivity %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestRedactWithDetails(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	redactor := NewDefaultPIIRedactor(logger)

	input := "Contact alice@test.com or call 555-123-4567"
	result := redactor.RedactWithDetails(input)

	if result.OriginalLength == 0 {
		t.Error("Expected original length to be set")
	}

	if result.RedactionCount == 0 {
		t.Error("Expected redaction count > 0")
	}

	if len(result.PatternsMatched) == 0 {
		t.Error("Expected patterns matched to be populated")
	}

	if !strings.Contains(result.RedactedText, "[EMAIL_REDACTED]") {
		t.Error("Expected email to be redacted")
	}

	if !strings.Contains(result.RedactedText, "[PHONE_REDACTED]") {
		t.Error("Expected phone to be redacted")
	}

	t.Logf("Redaction result: %+v", result)
}

func TestAddCustomPattern(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	redactor := NewDefaultPIIRedactor(logger)

	// Add custom pattern for employee IDs
	err := redactor.AddCustomPattern("employee_id", `EMP\d{6}`, "[EMPID_REDACTED]")
	if err != nil {
		t.Fatalf("Failed to add custom pattern: %v", err)
	}

	input := "Employee EMP123456 accessed the system"
	result := redactor.RedactQuery(input)

	expected := "Employee [EMPID_REDACTED] accessed the system"
	if result != expected {
		t.Errorf("Expected: %s, Got: %s", expected, result)
	}
}

func TestEmptyInput(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	redactor := NewDefaultPIIRedactor(logger)

	result := redactor.RedactQuery("")
	if result != "" {
		t.Errorf("Expected empty string for empty input, got: %s", result)
	}

	result = redactor.RedactSQL("")
	if result != "" {
		t.Errorf("Expected empty string for empty SQL, got: %s", result)
	}
}

func BenchmarkRedactQuery(b *testing.B) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	redactor := NewDefaultPIIRedactor(logger)

	query := "Contact john.doe@example.com or call (555) 123-4567. SSN: 123-45-6789"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = redactor.RedactQuery(query)
	}
}

func BenchmarkRedactQueryComplex(b *testing.B) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	redactor := NewDefaultPIIRedactor(logger)

	query := `SELECT * FROM customers WHERE
		email = 'alice@example.com'
		AND phone = '555-123-4567'
		AND ssn = '123-45-6789'
		AND card_number = '4111-1111-1111-1111'
		AND ip_address = '192.168.1.1'`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = redactor.RedactSQL(query)
	}
}
