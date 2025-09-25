package validation

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrInvalidInput    = errors.New("invalid input")
	ErrInputTooLong    = errors.New("input too long")
	ErrInvalidFormat   = errors.New("invalid format")
	ErrSuspiciousInput = errors.New("suspicious input detected")
)

// ValidateTextQuery validates user text queries
func ValidateTextQuery(query string) error {
	if query == "" {
		return ErrInvalidInput
	}

	// Length check
	if len(query) > 1000 {
		return ErrInputTooLong
	}

	// Check for suspicious patterns
	suspiciousPatterns := []string{
		"<script",
		"javascript:",
		"vbscript:",
		"onload=",
		"onerror=",
		"eval(",
		"exec(",
		"system(",
		"DROP TABLE",
		"DELETE FROM",
		"INSERT INTO",
		"UPDATE SET",
		"UNION SELECT",
		"--",
		";--",
		"/*",
		"*/",
		"xp_",
		"sp_",
	}

	queryLower := strings.ToLower(query)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(queryLower, strings.ToLower(pattern)) {
			return ErrSuspiciousInput
		}
	}

	return nil
}

// ValidateSQL validates SQL queries (more restrictive)
func ValidateSQL(sql string) error {
	if sql == "" {
		return ErrInvalidInput
	}

	// Length check
	if len(sql) > 2000 {
		return ErrInputTooLong
	}

	// Only allow SELECT statements
	sqlTrimmed := strings.TrimSpace(strings.ToUpper(sql))
	if !strings.HasPrefix(sqlTrimmed, "SELECT") {
		return ErrInvalidFormat
	}

	// Disallow dangerous keywords
	dangerousKeywords := []string{
		"DROP", "DELETE", "INSERT", "UPDATE", "ALTER", "CREATE",
		"EXEC", "EXECUTE", "XP_", "SP_", "SHUTDOWN", "GRANT", "REVOKE",
	}

	sqlUpper := strings.ToUpper(sql)
	for _, keyword := range dangerousKeywords {
		if strings.Contains(sqlUpper, keyword) {
			return ErrSuspiciousInput
		}
	}

	return nil
}

// SanitizeString removes potentially dangerous characters
func SanitizeString(input string) string {
	// Remove null bytes and control characters
	re := regexp.MustCompile(`[\x00-\x1f\x7f-\x9f]`)
	sanitized := re.ReplaceAllString(input, "")

	// Trim whitespace
	sanitized = strings.TrimSpace(sanitized)

	return sanitized
}

// ValidateFileUpload validates file uploads
func ValidateFileUpload(filename string, size int64) error {
	// Size limit: 10MB
	if size > 10*1024*1024 {
		return ErrInputTooLong
	}

	// Allowed extensions for audio files
	allowedExtensions := []string{".wav", ".mp3", ".m4a", ".webm"}
	filename = strings.ToLower(filename)

	validExtension := false
	for _, ext := range allowedExtensions {
		if strings.HasSuffix(filename, ext) {
			validExtension = true
			break
		}
	}

	if !validExtension {
		return ErrInvalidFormat
	}

	return nil
}