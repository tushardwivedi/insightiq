package services

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"regexp"
	"strings"
	"time"
)

// DataValidationService provides data quality and accuracy validation
type DataValidationService struct {
	logger *slog.Logger
}

// ValidationResult represents the result of data validation
type ValidationResult struct {
	IsValid         bool                   `json:"is_valid"`
	QualityScore    float64                `json:"quality_score"` // 0-100
	Warnings        []ValidationWarning    `json:"warnings"`
	Errors          []ValidationError      `json:"errors"`
	Statistics      ValidationStatistics   `json:"statistics"`
	Recommendations []string               `json:"recommendations"`
}

// ValidationWarning represents a data quality warning
type ValidationWarning struct {
	Type        string `json:"type"`
	Severity    string `json:"severity"` // low, medium, high
	Message     string `json:"message"`
	FieldName   string `json:"field_name,omitempty"`
	RecordCount int    `json:"record_count,omitempty"`
}

// ValidationError represents a data quality error
type ValidationError struct {
	Type      string `json:"type"`
	Message   string `json:"message"`
	FieldName string `json:"field_name,omitempty"`
}

// ValidationStatistics provides statistical overview of data quality
type ValidationStatistics struct {
	TotalRecords       int                    `json:"total_records"`
	CompleteRecords    int                    `json:"complete_records"`
	IncompleteRecords  int                    `json:"incomplete_records"`
	MissingDataPercent float64                `json:"missing_data_percent"`
	FieldStatistics    map[string]FieldStats  `json:"field_statistics"`
}

// FieldStats represents statistics for a single field
type FieldStats struct {
	FieldName         string  `json:"field_name"`
	TotalValues       int     `json:"total_values"`
	NonNullValues     int     `json:"non_null_values"`
	NullValues        int     `json:"null_values"`
	NullPercent       float64 `json:"null_percent"`
	UniqueValues      int     `json:"unique_values"`
	DuplicatePercent  float64 `json:"duplicate_percent"`
	DataType          string  `json:"data_type"`
	HasInconsistency  bool    `json:"has_inconsistency"`
}

// NewDataValidationService creates a new data validation service
func NewDataValidationService(logger *slog.Logger) *DataValidationService {
	return &DataValidationService{
		logger: logger.With("service", "data_validation"),
	}
}

// ValidateData performs comprehensive data validation
func (dvs *DataValidationService) ValidateData(
	ctx context.Context,
	data []map[string]interface{},
) (*ValidationResult, error) {
	start := time.Now()
	dvs.logger.Info("Starting data validation", "record_count", len(data))

	if len(data) == 0 {
		return &ValidationResult{
			IsValid:      false,
			QualityScore: 0,
			Errors: []ValidationError{
				{Type: "empty_dataset", Message: "No data to validate"},
			},
		}, nil
	}

	result := &ValidationResult{
		IsValid:  true,
		Warnings: []ValidationWarning{},
		Errors:   []ValidationError{},
	}

	// 1. Detect missing data patterns
	dvs.detectMissingData(data, result)

	// 2. Check for data inconsistencies
	dvs.detectInconsistencies(data, result)

	// 3. Validate data types
	dvs.validateDataTypes(data, result)

	// 4. Detect potential anomalies
	dvs.detectAnomalies(data, result)

	// 5. Check for suspicious patterns
	dvs.detectSuspiciousPatterns(data, result)

	// 6. Calculate statistics
	result.Statistics = dvs.calculateStatistics(data)

	// 7. Calculate overall quality score
	result.QualityScore = dvs.calculateQualityScore(result)

	// 8. Generate recommendations
	result.Recommendations = dvs.generateRecommendations(result)

	// Set IsValid based on errors
	result.IsValid = len(result.Errors) == 0

	duration := time.Since(start)
	dvs.logger.Info("Data validation completed",
		"is_valid", result.IsValid,
		"quality_score", result.QualityScore,
		"warnings", len(result.Warnings),
		"errors", len(result.Errors),
		"duration", duration)

	return result, nil
}

// detectMissingData identifies missing data patterns
func (dvs *DataValidationService) detectMissingData(
	data []map[string]interface{},
	result *ValidationResult,
) {
	if len(data) == 0 {
		return
	}

	// Get all field names from first record
	fieldNames := make([]string, 0)
	for key := range data[0] {
		fieldNames = append(fieldNames, key)
	}

	// Check each field for missing data
	for _, fieldName := range fieldNames {
		nullCount := 0
		for _, record := range data {
			value, exists := record[fieldName]
			if !exists || value == nil || value == "" {
				nullCount++
			}
		}

		nullPercent := float64(nullCount) / float64(len(data)) * 100

		// Flag fields with high missing data percentage
		if nullPercent > 50 {
			result.Warnings = append(result.Warnings, ValidationWarning{
				Type:        "high_missing_data",
				Severity:    "high",
				Message:     fmt.Sprintf("Field '%s' has %.1f%% missing data", fieldName, nullPercent),
				FieldName:   fieldName,
				RecordCount: nullCount,
			})
		} else if nullPercent > 20 {
			result.Warnings = append(result.Warnings, ValidationWarning{
				Type:        "moderate_missing_data",
				Severity:    "medium",
				Message:     fmt.Sprintf("Field '%s' has %.1f%% missing data", fieldName, nullPercent),
				FieldName:   fieldName,
				RecordCount: nullCount,
			})
		}
	}

	// Check for completely empty records
	emptyRecords := 0
	for _, record := range data {
		isEmpty := true
		for _, value := range record {
			if value != nil && value != "" {
				isEmpty = false
				break
			}
		}
		if isEmpty {
			emptyRecords++
		}
	}

	if emptyRecords > 0 {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Type:        "empty_records",
			Severity:    "high",
			Message:     fmt.Sprintf("Found %d completely empty records", emptyRecords),
			RecordCount: emptyRecords,
		})
	}
}

// detectInconsistencies identifies data inconsistencies
func (dvs *DataValidationService) detectInconsistencies(
	data []map[string]interface{},
	result *ValidationResult,
) {
	if len(data) == 0 {
		return
	}

	// Get all field names
	fieldNames := make([]string, 0)
	for key := range data[0] {
		fieldNames = append(fieldNames, key)
	}

	// Check for type inconsistencies in each field
	for _, fieldName := range fieldNames {
		typeMap := make(map[string]int)

		for _, record := range data {
			value, exists := record[fieldName]
			if !exists || value == nil {
				continue
			}

			valueType := fmt.Sprintf("%T", value)
			typeMap[valueType]++
		}

		// If field has multiple types, it's inconsistent
		if len(typeMap) > 1 {
			result.Warnings = append(result.Warnings, ValidationWarning{
				Type:      "type_inconsistency",
				Severity:  "medium",
				Message:   fmt.Sprintf("Field '%s' has inconsistent data types", fieldName),
				FieldName: fieldName,
			})
		}
	}

	// Check for duplicate records (exact matches)
	seen := make(map[string]bool)
	duplicateCount := 0

	for _, record := range data {
		// Create a simple hash of the record
		hash := fmt.Sprintf("%v", record)
		if seen[hash] {
			duplicateCount++
		} else {
			seen[hash] = true
		}
	}

	if duplicateCount > 0 {
		duplicatePercent := float64(duplicateCount) / float64(len(data)) * 100
		severity := "low"
		if duplicatePercent > 10 {
			severity = "high"
		} else if duplicatePercent > 5 {
			severity = "medium"
		}

		result.Warnings = append(result.Warnings, ValidationWarning{
			Type:        "duplicate_records",
			Severity:    severity,
			Message:     fmt.Sprintf("Found %d duplicate records (%.1f%%)", duplicateCount, duplicatePercent),
			RecordCount: duplicateCount,
		})
	}
}

// validateDataTypes validates data types and formats
func (dvs *DataValidationService) validateDataTypes(
	data []map[string]interface{},
	result *ValidationResult,
) {
	// Email validation pattern
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	// URL validation pattern
	urlPattern := regexp.MustCompile(`^https?://`)

	for _, record := range data {
		for fieldName, value := range record {
			if value == nil {
				continue
			}

			strValue, ok := value.(string)
			if !ok {
				continue
			}

			// Check for malformed emails
			if strings.Contains(strings.ToLower(fieldName), "email") {
				if strValue != "" && !emailPattern.MatchString(strValue) {
					result.Warnings = append(result.Warnings, ValidationWarning{
						Type:      "invalid_email_format",
						Severity:  "low",
						Message:   fmt.Sprintf("Field '%s' contains invalid email format: %s", fieldName, strValue),
						FieldName: fieldName,
					})
				}
			}

			// Check for malformed URLs
			if strings.Contains(strings.ToLower(fieldName), "url") || strings.Contains(strings.ToLower(fieldName), "link") {
				if strValue != "" && !urlPattern.MatchString(strValue) {
					result.Warnings = append(result.Warnings, ValidationWarning{
						Type:      "invalid_url_format",
						Severity:  "low",
						Message:   fmt.Sprintf("Field '%s' contains invalid URL format: %s", fieldName, strValue),
						FieldName: fieldName,
					})
				}
			}
		}
	}
}

// detectAnomalies identifies statistical anomalies in numeric data
func (dvs *DataValidationService) detectAnomalies(
	data []map[string]interface{},
	result *ValidationResult,
) {
	if len(data) == 0 {
		return
	}

	// Get all field names
	fieldNames := make([]string, 0)
	for key := range data[0] {
		fieldNames = append(fieldNames, key)
	}

	// Check numeric fields for outliers
	for _, fieldName := range fieldNames {
		var values []float64

		for _, record := range data {
			value, exists := record[fieldName]
			if !exists || value == nil {
				continue
			}

			// Try to convert to float64
			var numValue float64
			switch v := value.(type) {
			case float64:
				numValue = v
			case int:
				numValue = float64(v)
			case int64:
				numValue = float64(v)
			default:
				continue
			}

			values = append(values, numValue)
		}

		// Need at least 10 values for meaningful statistics
		if len(values) < 10 {
			continue
		}

		// Calculate mean and standard deviation
		mean := dvs.calculateMean(values)
		stdDev := dvs.calculateStdDev(values, mean)

		// Check for outliers (values > 3 standard deviations from mean)
		outlierCount := 0
		for _, val := range values {
			if math.Abs(val-mean) > 3*stdDev {
				outlierCount++
			}
		}

		if outlierCount > 0 {
			outlierPercent := float64(outlierCount) / float64(len(values)) * 100
			if outlierPercent > 5 {
				result.Warnings = append(result.Warnings, ValidationWarning{
					Type:        "statistical_outliers",
					Severity:    "medium",
					Message:     fmt.Sprintf("Field '%s' has %d outliers (%.1f%% of numeric values)", fieldName, outlierCount, outlierPercent),
					FieldName:   fieldName,
					RecordCount: outlierCount,
				})
			}
		}
	}
}

// detectSuspiciousPatterns identifies suspicious data patterns
func (dvs *DataValidationService) detectSuspiciousPatterns(
	data []map[string]interface{},
	result *ValidationResult,
) {
	if len(data) == 0 {
		return
	}

	// Check for fields with all identical values (except nulls)
	fieldNames := make([]string, 0)
	for key := range data[0] {
		fieldNames = append(fieldNames, key)
	}

	for _, fieldName := range fieldNames {
		uniqueValues := make(map[string]bool)
		nonNullCount := 0

		for _, record := range data {
			value, exists := record[fieldName]
			if !exists || value == nil || value == "" {
				continue
			}

			nonNullCount++
			strValue := fmt.Sprintf("%v", value)
			uniqueValues[strValue] = true
		}

		// If all non-null values are identical, flag as suspicious
		if nonNullCount > 10 && len(uniqueValues) == 1 {
			result.Warnings = append(result.Warnings, ValidationWarning{
				Type:      "uniform_data",
				Severity:  "low",
				Message:   fmt.Sprintf("Field '%s' has identical values across all records", fieldName),
				FieldName: fieldName,
			})
		}
	}

	// Check for sequential patterns that might indicate synthetic data
	// This is a basic check - can be enhanced
	for i := 0; i < len(data)-5; i++ {
		// Simple check for identical consecutive records
		if fmt.Sprintf("%v", data[i]) == fmt.Sprintf("%v", data[i+1]) &&
			fmt.Sprintf("%v", data[i]) == fmt.Sprintf("%v", data[i+2]) {
			result.Warnings = append(result.Warnings, ValidationWarning{
				Type:     "suspicious_pattern",
				Severity: "medium",
				Message:  "Detected repeated identical consecutive records, possible data duplication",
			})
			break
		}
	}
}

// calculateStatistics computes validation statistics
func (dvs *DataValidationService) calculateStatistics(
	data []map[string]interface{},
) ValidationStatistics {
	stats := ValidationStatistics{
		TotalRecords:    len(data),
		FieldStatistics: make(map[string]FieldStats),
	}

	if len(data) == 0 {
		return stats
	}

	// Get all field names
	fieldNames := make([]string, 0)
	for key := range data[0] {
		fieldNames = append(fieldNames, key)
	}

	// Calculate statistics for each field
	completeRecords := 0
	for _, record := range data {
		isComplete := true
		for _, value := range record {
			if value == nil || value == "" {
				isComplete = false
				break
			}
		}
		if isComplete {
			completeRecords++
		}
	}

	stats.CompleteRecords = completeRecords
	stats.IncompleteRecords = stats.TotalRecords - completeRecords
	if stats.TotalRecords > 0 {
		stats.MissingDataPercent = float64(stats.IncompleteRecords) / float64(stats.TotalRecords) * 100
	}

	// Field-level statistics
	for _, fieldName := range fieldNames {
		fieldStat := FieldStats{
			FieldName:   fieldName,
			TotalValues: len(data),
		}

		uniqueValues := make(map[string]bool)
		for _, record := range data {
			value, exists := record[fieldName]
			if !exists || value == nil || value == "" {
				fieldStat.NullValues++
			} else {
				fieldStat.NonNullValues++
				uniqueValues[fmt.Sprintf("%v", value)] = true

				// Detect data type
				if fieldStat.DataType == "" {
					fieldStat.DataType = fmt.Sprintf("%T", value)
				}
			}
		}

		fieldStat.UniqueValues = len(uniqueValues)
		if fieldStat.TotalValues > 0 {
			fieldStat.NullPercent = float64(fieldStat.NullValues) / float64(fieldStat.TotalValues) * 100
			if fieldStat.NonNullValues > 0 {
				fieldStat.DuplicatePercent = (1 - float64(fieldStat.UniqueValues)/float64(fieldStat.NonNullValues)) * 100
			}
		}

		stats.FieldStatistics[fieldName] = fieldStat
	}

	return stats
}

// calculateQualityScore computes overall quality score (0-100)
func (dvs *DataValidationService) calculateQualityScore(result *ValidationResult) float64 {
	score := 100.0

	// Deduct points for errors
	score -= float64(len(result.Errors)) * 10

	// Deduct points for warnings based on severity
	for _, warning := range result.Warnings {
		switch warning.Severity {
		case "high":
			score -= 5
		case "medium":
			score -= 3
		case "low":
			score -= 1
		}
	}

	// Bonus for data completeness
	if result.Statistics.TotalRecords > 0 {
		completenessBonus := float64(result.Statistics.CompleteRecords) / float64(result.Statistics.TotalRecords) * 10
		score += completenessBonus
	}

	// Ensure score is in 0-100 range
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}

// generateRecommendations creates actionable recommendations
func (dvs *DataValidationService) generateRecommendations(result *ValidationResult) []string {
	recommendations := make([]string, 0)

	// Based on missing data
	if result.Statistics.MissingDataPercent > 20 {
		recommendations = append(recommendations,
			fmt.Sprintf("Address missing data: %.1f%% of records are incomplete. Consider data imputation or collection improvement.",
			result.Statistics.MissingDataPercent))
	}

	// Based on warnings
	highSeverityCount := 0
	for _, warning := range result.Warnings {
		if warning.Severity == "high" {
			highSeverityCount++
		}
	}

	if highSeverityCount > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("Resolve %d high-severity data quality issues before analysis", highSeverityCount))
	}

	// Based on quality score
	if result.QualityScore < 70 {
		recommendations = append(recommendations,
			"Data quality score is below acceptable threshold (70). Review and clean data before using for critical analysis.")
	}

	// Check for duplicate data
	for _, warning := range result.Warnings {
		if warning.Type == "duplicate_records" {
			recommendations = append(recommendations,
				"Remove duplicate records to improve data accuracy and reduce bias in analysis")
			break
		}
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Data quality is acceptable. No major issues detected.")
	}

	return recommendations
}

// Helper functions
func (dvs *DataValidationService) calculateMean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func (dvs *DataValidationService) calculateStdDev(values []float64, mean float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sumSquaredDiff := 0.0
	for _, v := range values {
		diff := v - mean
		sumSquaredDiff += diff * diff
	}
	variance := sumSquaredDiff / float64(len(values))
	return math.Sqrt(variance)
}
