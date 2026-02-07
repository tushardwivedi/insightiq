package security

import (
	"log/slog"
	"regexp"
	"strings"
)

// PIIRedactor handles the redaction of Personally Identifiable Information (PII)
// from query text and SQL statements to comply with GDPR, HIPAA, and other privacy regulations
type PIIRedactor interface {
	RedactQuery(query string) string
	RedactSQL(sql string) string
	RedactText(text string) string
}

// DefaultPIIRedactor is the default implementation of PIIRedactor
type DefaultPIIRedactor struct {
	logger   *slog.Logger
	patterns []RedactionPattern
}

// RedactionPattern defines a pattern to match and redact
type RedactionPattern struct {
	Name        string
	Pattern     *regexp.Regexp
	Replacement string
}

// NewDefaultPIIRedactor creates a new PII redactor with default patterns
func NewDefaultPIIRedactor(logger *slog.Logger) *DefaultPIIRedactor {
	redactor := &DefaultPIIRedactor{
		logger:   logger.With("component", "pii_redactor"),
		patterns: make([]RedactionPattern, 0),
	}

	// Initialize default redaction patterns
	redactor.initializePatterns()

	return redactor
}

// initializePatterns sets up all the PII detection patterns
// Order matters: more specific patterns should come first to avoid false matches
func (r *DefaultPIIRedactor) initializePatterns() {
	// Credit Card Numbers (Visa, MasterCard, Amex, Discover) - FIRST to avoid phone false positives
	// 16 digits with optional spaces or dashes
	r.patterns = append(r.patterns, RedactionPattern{
		Name:        "credit_card",
		Pattern:     regexp.MustCompile(`\b(?:\d{4}[-\s]?){3}\d{4}\b`),
		Replacement: "[CARD_REDACTED]",
	})

	// Email addresses
	r.patterns = append(r.patterns, RedactionPattern{
		Name:        "email",
		Pattern:     regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`),
		Replacement: "[EMAIL_REDACTED]",
	})

	// Phone numbers (various formats)
	// US: (123) 456-7890, 123-456-7890, 1234567890, +1-123-456-7890
	r.patterns = append(r.patterns, RedactionPattern{
		Name:        "phone",
		Pattern:     regexp.MustCompile(`(\+?\d{1,3}[-.\s]?)?\(?\d{3}\)?[-.\s]?\d{3}[-.\s]?\d{4}\b`),
		Replacement: "[PHONE_REDACTED]",
	})

	// US Social Security Numbers (SSN)
	// Formats: 123-45-6789, 123 45 6789, 123456789
	r.patterns = append(r.patterns, RedactionPattern{
		Name:        "ssn",
		Pattern:     regexp.MustCompile(`\b\d{3}[-\s]?\d{2}[-\s]?\d{4}\b`),
		Replacement: "[SSN_REDACTED]",
	})

	// IPv4 Addresses
	r.patterns = append(r.patterns, RedactionPattern{
		Name:        "ipv4",
		Pattern:     regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`),
		Replacement: "[IP_REDACTED]",
	})

	// Date of Birth patterns (MM/DD/YYYY, MM-DD-YYYY, etc.)
	r.patterns = append(r.patterns, RedactionPattern{
		Name:        "date_of_birth",
		Pattern:     regexp.MustCompile(`\b(?:0?[1-9]|1[0-2])[-/](?:0?[1-9]|[12]\d|3[01])[-/](?:19|20)\d{2}\b`),
		Replacement: "[DOB_REDACTED]",
	})

	// Passport Numbers (basic pattern - alphanumeric 6-9 characters)
	r.patterns = append(r.patterns, RedactionPattern{
		Name:        "passport",
		Pattern:     regexp.MustCompile(`\b[A-Z]{1,2}\d{6,9}\b`),
		Replacement: "[PASSPORT_REDACTED]",
	})

	// Bank Account Numbers (8-17 digits)
	r.patterns = append(r.patterns, RedactionPattern{
		Name:        "bank_account",
		Pattern:     regexp.MustCompile(`\b\d{8,17}\b`),
		Replacement: "[ACCOUNT_REDACTED]",
	})
}

// RedactQuery redacts PII from user query text
func (r *DefaultPIIRedactor) RedactQuery(query string) string {
	if query == "" {
		return query
	}

	redactedQuery := query
	redactedCount := 0

	// Apply all redaction patterns
	for _, pattern := range r.patterns {
		if pattern.Pattern.MatchString(redactedQuery) {
			matches := pattern.Pattern.FindAllString(redactedQuery, -1)
			redactedQuery = pattern.Pattern.ReplaceAllString(redactedQuery, pattern.Replacement)
			redactedCount += len(matches)

			r.logger.Debug("PII redacted from query",
				"pattern", pattern.Name,
				"matches", len(matches))
		}
	}

	if redactedCount > 0 {
		r.logger.Info("PII redaction applied to query",
			"total_redactions", redactedCount,
			"original_length", len(query),
			"redacted_length", len(redactedQuery))
	}

	return redactedQuery
}

// RedactSQL redacts PII from generated SQL statements
func (r *DefaultPIIRedactor) RedactSQL(sql string) string {
	if sql == "" {
		return sql
	}

	redactedSQL := sql
	redactedCount := 0

	// Apply all redaction patterns
	for _, pattern := range r.patterns {
		if pattern.Pattern.MatchString(redactedSQL) {
			matches := pattern.Pattern.FindAllString(redactedSQL, -1)
			redactedSQL = pattern.Pattern.ReplaceAllString(redactedSQL, pattern.Replacement)
			redactedCount += len(matches)

			r.logger.Debug("PII redacted from SQL",
				"pattern", pattern.Name,
				"matches", len(matches))
		}
	}

	if redactedCount > 0 {
		r.logger.Info("PII redaction applied to SQL",
			"total_redactions", redactedCount,
			"original_length", len(sql),
			"redacted_length", len(redactedSQL))
	}

	return redactedSQL
}

// RedactText is a general-purpose method for redacting PII from any text
func (r *DefaultPIIRedactor) RedactText(text string) string {
	return r.RedactQuery(text) // Same logic as query redaction
}

// AddCustomPattern allows adding custom redaction patterns at runtime
func (r *DefaultPIIRedactor) AddCustomPattern(name string, pattern string, replacement string) error {
	compiledPattern, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	r.patterns = append(r.patterns, RedactionPattern{
		Name:        name,
		Pattern:     compiledPattern,
		Replacement: replacement,
	})

	r.logger.Info("Custom redaction pattern added",
		"name", name,
		"pattern", pattern)

	return nil
}

// GetRedactionStats returns statistics about redaction patterns
func (r *DefaultPIIRedactor) GetRedactionStats() map[string]int {
	stats := make(map[string]int)
	for _, pattern := range r.patterns {
		stats[pattern.Name] = 0 // Initialize counter
	}
	return stats
}

// ValidateRedaction checks if redaction was successful by scanning for common PII patterns
func (r *DefaultPIIRedactor) ValidateRedaction(original, redacted string) bool {
	// If strings are identical, no redaction occurred
	if original == redacted {
		return true // No PII found, which is valid
	}

	// Check if any obvious PII patterns remain in redacted text
	for _, pattern := range r.patterns {
		if pattern.Pattern.MatchString(redacted) {
			r.logger.Warn("Potential PII leak detected after redaction",
				"pattern", pattern.Name)
			return false
		}
	}

	return true
}

// SensitivityLevel represents the level of PII sensitivity
type SensitivityLevel int

const (
	SensitivityLow SensitivityLevel = iota
	SensitivityMedium
	SensitivityHigh
	SensitivityCritical
)

// DetectSensitivityLevel analyzes text and determines PII sensitivity level
func (r *DefaultPIIRedactor) DetectSensitivityLevel(text string) SensitivityLevel {
	if text == "" {
		return SensitivityLow
	}

	highSensitivityPatterns := []string{"ssn", "credit_card", "passport", "bank_account"}
	mediumSensitivityPatterns := []string{"email", "phone", "date_of_birth"}

	matchedPatterns := make(map[string]bool)

	for _, pattern := range r.patterns {
		if pattern.Pattern.MatchString(text) {
			matchedPatterns[pattern.Name] = true
		}
	}

	// Check for critical patterns
	for _, name := range highSensitivityPatterns {
		if matchedPatterns[name] {
			return SensitivityCritical
		}
	}

	// Check for high sensitivity patterns (multiple medium patterns)
	mediumCount := 0
	for _, name := range mediumSensitivityPatterns {
		if matchedPatterns[name] {
			mediumCount++
		}
	}

	if mediumCount >= 2 {
		return SensitivityHigh
	} else if mediumCount == 1 {
		return SensitivityMedium
	}

	return SensitivityLow
}

// RedactWithContext provides contextual information about what was redacted
type RedactionResult struct {
	RedactedText     string
	OriginalLength   int
	RedactedLength   int
	RedactionCount   int
	PatternsMatched  []string
	SensitivityLevel SensitivityLevel
}

// RedactWithDetails performs redaction and returns detailed results
func (r *DefaultPIIRedactor) RedactWithDetails(text string) RedactionResult {
	result := RedactionResult{
		RedactedText:    text,
		OriginalLength:  len(text),
		PatternsMatched: make([]string, 0),
	}

	if text == "" {
		return result
	}

	redactedText := text
	totalRedactions := 0

	// Apply all redaction patterns
	for _, pattern := range r.patterns {
		if pattern.Pattern.MatchString(redactedText) {
			matches := pattern.Pattern.FindAllString(redactedText, -1)
			redactedText = pattern.Pattern.ReplaceAllString(redactedText, pattern.Replacement)

			if len(matches) > 0 {
				result.PatternsMatched = append(result.PatternsMatched, pattern.Name)
				totalRedactions += len(matches)
			}
		}
	}

	result.RedactedText = redactedText
	result.RedactedLength = len(redactedText)
	result.RedactionCount = totalRedactions
	result.SensitivityLevel = r.DetectSensitivityLevel(text)

	return result
}

// FormatRedactionResult creates a human-readable summary of redaction
func FormatRedactionResult(result RedactionResult) string {
	if result.RedactionCount == 0 {
		return "No PII detected"
	}

	parts := []string{
		"PII Redaction Summary:",
	}

	parts = append(parts, "- Redactions: "+string(rune(result.RedactionCount)))
	parts = append(parts, "- Patterns: "+strings.Join(result.PatternsMatched, ", "))

	var sensitivity string
	switch result.SensitivityLevel {
	case SensitivityCritical:
		sensitivity = "CRITICAL"
	case SensitivityHigh:
		sensitivity = "HIGH"
	case SensitivityMedium:
		sensitivity = "MEDIUM"
	default:
		sensitivity = "LOW"
	}
	parts = append(parts, "- Sensitivity: "+sensitivity)

	return strings.Join(parts, "\n")
}
