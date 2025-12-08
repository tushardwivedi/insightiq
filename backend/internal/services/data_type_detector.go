package services

import (
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// DataTypeDetector provides enhanced data type detection with confidence scoring
type DataTypeDetector struct {
	logger *slog.Logger
}

// DataType represents detected data types
type DataType string

const (
	TypeInteger    DataType = "INTEGER"
	TypeBigInt     DataType = "BIGINT"
	TypeNumeric    DataType = "NUMERIC"
	TypeFloat      DataType = "FLOAT"
	TypeBoolean    DataType = "BOOLEAN"
	TypeTimestamp  DataType = "TIMESTAMP"
	TypeDate       DataType = "DATE"
	TypeTime       DataType = "TIME"
	TypeEmail      DataType = "EMAIL"
	TypeURL        DataType = "URL"
	TypeJSON       DataType = "JSON"
	TypeUUID       DataType = "UUID"
	TypeText       DataType = "TEXT"
	TypeVarchar    DataType = "VARCHAR"
	TypeUnknown    DataType = "UNKNOWN"
)

// TypeDetectionResult represents the result of type detection
type TypeDetectionResult struct {
	DetectedType   DataType          `json:"detected_type"`
	Confidence     float64           `json:"confidence"` // 0.0 to 1.0
	SampleValues   []string          `json:"sample_values,omitempty"`
	TypeCounts     map[DataType]int  `json:"type_counts"`
	IsMixedType    bool              `json:"is_mixed_type"`
	Recommendations []string          `json:"recommendations,omitempty"`
}

// NewDataTypeDetector creates a new data type detector
func NewDataTypeDetector(logger *slog.Logger) *DataTypeDetector {
	return &DataTypeDetector{
		logger: logger.With("service", "data_type_detector"),
	}
}

// DetectColumnType detects the data type of a column from sample values
func (dtd *DataTypeDetector) DetectColumnType(columnName string, values []interface{}) *TypeDetectionResult {
	result := &TypeDetectionResult{
		TypeCounts:   make(map[DataType]int),
		SampleValues: make([]string, 0),
	}

	if len(values) == 0 {
		result.DetectedType = TypeUnknown
		result.Confidence = 0.0
		return result
	}

	// Collect sample values and detect type for each value
	totalValues := 0
	nullCount := 0

	for _, value := range values {
		if value == nil {
			nullCount++
			continue
		}

		totalValues++

		// Convert to string for analysis
		strValue := fmt.Sprintf("%v", value)

		// Store first 5 samples
		if len(result.SampleValues) < 5 {
			result.SampleValues = append(result.SampleValues, strValue)
		}

		// Detect type of this value
		detectedType := dtd.detectValueType(strValue)
		result.TypeCounts[detectedType]++
	}

	// If all values are null
	if totalValues == 0 {
		result.DetectedType = TypeUnknown
		result.Confidence = 0.0
		result.Recommendations = append(result.Recommendations,
			"All values are null. Unable to determine type.")
		return result
	}

	// Determine final type based on majority
	finalType, confidence := dtd.determineFinalType(result.TypeCounts, totalValues)
	result.DetectedType = finalType
	result.Confidence = confidence

	// Check for mixed types
	if len(result.TypeCounts) > 1 {
		result.IsMixedType = true
		result.Recommendations = append(result.Recommendations,
			fmt.Sprintf("Column '%s' contains mixed types. Consider data cleaning.", columnName))
	}

	// Add recommendations based on detected type
	result.Recommendations = append(result.Recommendations,
		dtd.getTypeRecommendations(finalType, confidence)...)

	dtd.logger.Debug("Type detection completed",
		"column", columnName,
		"detected_type", finalType,
		"confidence", confidence,
		"is_mixed", result.IsMixedType)

	return result
}

// detectValueType detects the type of a single value
func (dtd *DataTypeDetector) detectValueType(value string) DataType {
	if value == "" {
		return TypeText
	}

	// Check for UUID
	if dtd.isUUID(value) {
		return TypeUUID
	}

	// Check for Boolean
	if dtd.isBoolean(value) {
		return TypeBoolean
	}

	// Check for Email
	if dtd.isEmail(value) {
		return TypeEmail
	}

	// Check for URL
	if dtd.isURL(value) {
		return TypeURL
	}

	// Check for JSON
	if dtd.isJSON(value) {
		return TypeJSON
	}

	// Check for Timestamp
	if dtd.isTimestamp(value) {
		return TypeTimestamp
	}

	// Check for Date
	if dtd.isDate(value) {
		return TypeDate
	}

	// Check for Time
	if dtd.isTime(value) {
		return TypeTime
	}

	// Check for Integer
	if dtd.isInteger(value) {
		// Check if it's too large for regular integer
		if len(value) > 10 {
			return TypeBigInt
		}
		return TypeInteger
	}

	// Check for Float/Numeric
	if dtd.isFloat(value) {
		return TypeNumeric
	}

	// Default to Text/Varchar
	if len(value) < 255 {
		return TypeVarchar
	}
	return TypeText
}

// Type checking functions

func (dtd *DataTypeDetector) isUUID(value string) bool {
	uuidPattern := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	return uuidPattern.MatchString(value)
}

func (dtd *DataTypeDetector) isBoolean(value string) bool {
	lower := strings.ToLower(strings.TrimSpace(value))
	return lower == "true" || lower == "false" ||
		   lower == "yes" || lower == "no" ||
		   lower == "1" || lower == "0" ||
		   lower == "t" || lower == "f"
}

func (dtd *DataTypeDetector) isEmail(value string) bool {
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailPattern.MatchString(value)
}

func (dtd *DataTypeDetector) isURL(value string) bool {
	urlPattern := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	return urlPattern.MatchString(value)
}

func (dtd *DataTypeDetector) isJSON(value string) bool {
	trimmed := strings.TrimSpace(value)
	return (strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}")) ||
		   (strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]"))
}

func (dtd *DataTypeDetector) isInteger(value string) bool {
	value = strings.TrimSpace(value)
	_, err := strconv.ParseInt(value, 10, 64)
	return err == nil
}

func (dtd *DataTypeDetector) isFloat(value string) bool {
	value = strings.TrimSpace(value)
	_, err := strconv.ParseFloat(value, 64)
	return err == nil
}

func (dtd *DataTypeDetector) isTimestamp(value string) bool {
	timestampFormats := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05.000",
		"01/02/2006 15:04:05",
		"02-Jan-2006 15:04:05",
	}

	for _, format := range timestampFormats {
		if _, err := time.Parse(format, value); err == nil {
			return true
		}
	}
	return false
}

func (dtd *DataTypeDetector) isDate(value string) bool {
	dateFormats := []string{
		"2006-01-02",
		"01/02/2006",
		"02-Jan-2006",
		"2006/01/02",
		"Jan 02, 2006",
	}

	for _, format := range dateFormats {
		if _, err := time.Parse(format, value); err == nil {
			return true
		}
	}
	return false
}

func (dtd *DataTypeDetector) isTime(value string) bool {
	timeFormats := []string{
		"15:04:05",
		"15:04",
		"3:04 PM",
		"3:04:05 PM",
	}

	for _, format := range timeFormats {
		if _, err := time.Parse(format, value); err == nil {
			return true
		}
	}
	return false
}

// determineFinalType determines the final type based on type counts
func (dtd *DataTypeDetector) determineFinalType(typeCounts map[DataType]int, totalValues int) (DataType, float64) {
	if len(typeCounts) == 0 {
		return TypeUnknown, 0.0
	}

	// Find the most common type
	var maxCount int
	var mostCommonType DataType

	for dataType, count := range typeCounts {
		if count > maxCount {
			maxCount = count
			mostCommonType = dataType
		}
	}

	// Calculate confidence based on percentage of values matching the detected type
	confidence := float64(maxCount) / float64(totalValues)

	// Special handling for numeric types
	// If we see INTEGER and NUMERIC, prefer NUMERIC for flexibility
	if typeCounts[TypeInteger] > 0 && typeCounts[TypeNumeric] > 0 {
		totalNumeric := typeCounts[TypeInteger] + typeCounts[TypeNumeric]
		confidence = float64(totalNumeric) / float64(totalValues)
		return TypeNumeric, confidence
	}

	// If we see INTEGER and BIGINT, prefer BIGINT for safety
	if typeCounts[TypeInteger] > 0 && typeCounts[TypeBigInt] > 0 {
		totalInt := typeCounts[TypeInteger] + typeCounts[TypeBigInt]
		confidence = float64(totalInt) / float64(totalValues)
		return TypeBigInt, confidence
	}

	// If we see specialized types (EMAIL, URL, UUID) mixed with TEXT, prefer specialized
	specializedTypes := []DataType{TypeEmail, TypeURL, TypeUUID, TypeJSON}
	for _, specialType := range specializedTypes {
		if typeCounts[specialType] > 0 {
			if float64(typeCounts[specialType])/float64(totalValues) > 0.8 {
				return specialType, float64(typeCounts[specialType]) / float64(totalValues)
			}
		}
	}

	return mostCommonType, confidence
}

// getTypeRecommendations provides recommendations based on detected type and confidence
func (dtd *DataTypeDetector) getTypeRecommendations(dataType DataType, confidence float64) []string {
	var recommendations []string

	// Low confidence warnings
	if confidence < 0.7 {
		recommendations = append(recommendations,
			fmt.Sprintf("Low confidence (%.1f%%) in type detection. Manual review recommended.", confidence*100))
	}

	// Type-specific recommendations
	switch dataType {
	case TypeInteger, TypeBigInt:
		if confidence < 0.95 {
			recommendations = append(recommendations,
				"Some non-numeric values detected. Consider using NUMERIC for flexibility.")
		}

	case TypeNumeric, TypeFloat:
		recommendations = append(recommendations,
			"Use NUMERIC type to preserve decimal precision.")

	case TypeTimestamp, TypeDate, TypeTime:
		recommendations = append(recommendations,
			"Ensure date/time format is consistent across all values.")

	case TypeEmail:
		recommendations = append(recommendations,
			"Consider adding email validation constraint.")

	case TypeURL:
		recommendations = append(recommendations,
			"Consider adding URL validation constraint.")

	case TypeBoolean:
		recommendations = append(recommendations,
			"Standardize boolean values to true/false or 1/0.")

	case TypeJSON:
		recommendations = append(recommendations,
			"Use JSON/JSONB type for efficient querying.")

	case TypeVarchar:
		recommendations = append(recommendations,
			"VARCHAR(255) recommended for short text fields.")

	case TypeText:
		recommendations = append(recommendations,
			"TEXT type recommended for long text content.")

	case TypeUnknown:
		recommendations = append(recommendations,
			"Unable to determine type. Default to TEXT type.")
	}

	return recommendations
}

// DetectBatchTypes detects types for multiple columns efficiently
func (dtd *DataTypeDetector) DetectBatchTypes(data map[string][]interface{}) map[string]*TypeDetectionResult {
	results := make(map[string]*TypeDetectionResult)

	for columnName, values := range data {
		results[columnName] = dtd.DetectColumnType(columnName, values)
	}

	dtd.logger.Info("Batch type detection completed",
		"columns_analyzed", len(results))

	return results
}

// GetSQLType converts detected type to SQL type string
func (dtd *DataTypeDetector) GetSQLType(detectedType DataType) string {
	switch detectedType {
	case TypeInteger:
		return "INTEGER"
	case TypeBigInt:
		return "BIGINT"
	case TypeNumeric, TypeFloat:
		return "NUMERIC"
	case TypeBoolean:
		return "BOOLEAN"
	case TypeTimestamp:
		return "TIMESTAMP"
	case TypeDate:
		return "DATE"
	case TypeTime:
		return "TIME"
	case TypeEmail, TypeURL, TypeUUID:
		return "VARCHAR(255)"
	case TypeJSON:
		return "JSONB"
	case TypeVarchar:
		return "VARCHAR(255)"
	case TypeText:
		return "TEXT"
	default:
		return "TEXT"
	}
}
