package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"insightiq/backend/internal/models"
)

type DataQualityService struct {
	db     *sqlx.DB
	logger *slog.Logger
}

// DataQualityReport contains comprehensive quality metrics for a table
type DataQualityReport struct {
	TableName         string
	TotalRows         int
	CompleteRows      int // Rows with no missing values
	RowsWithMissing   int
	DuplicateRows     int
	QualityScore      float64
	MissingDataPercent float64
	ColumnReports     map[string]*ColumnQualityReport
	ColumnsWithNulls  []string
}

// ColumnQualityReport contains quality metrics for a single column
type ColumnQualityReport struct {
	ColumnName      string
	DataType        string
	MissingCount    int
	MissingPercent  float64
	UniqueCount     int
	Cardinality     float64 // Unique values / Total rows
	Statistics      interface{} // Can be NumericStatistics or CategoricalStatistics
}

// NumericStatistics contains statistical metrics for numeric columns
type NumericStatistics struct {
	Mean     float64
	Median   float64
	StdDev   float64
	Min      float64
	Max      float64
	Q1       float64 // 25th percentile
	Q3       float64 // 75th percentile
	Outliers []interface{}
}

// CategoricalStatistics contains distribution metrics for text columns
type CategoricalStatistics struct {
	MostFrequent      string
	MostFrequentCount int
	Distribution      map[string]int // Top 10 values and their counts
}

// DateStatistics contains metrics for timestamp columns
type DateStatistics struct {
	MinDate      string
	MaxDate      string
	MostFrequentDate string
	DateRange    string
}

// NewDataQualityService creates a new data quality service
func NewDataQualityService(db *sqlx.DB, logger *slog.Logger) *DataQualityService {
	return &DataQualityService{
		db:     db,
		logger: logger,
	}
}

// AnalyzeDataQuality performs comprehensive data quality analysis on a table
func (dq *DataQualityService) AnalyzeDataQuality(
	ctx context.Context,
	tableName string,
	schema []models.ColumnInfo,
) (*DataQualityReport, error) {

	dq.logger.Info("Starting data quality analysis", "table", tableName)

	report := &DataQualityReport{
		TableName:     tableName,
		ColumnReports: make(map[string]*ColumnQualityReport),
		ColumnsWithNulls: []string{},
	}

	// 1. Count total rows
	totalRows, err := dq.countRows(ctx, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to count rows: %w", err)
	}
	report.TotalRows = totalRows

	if totalRows == 0 {
		dq.logger.Warn("Table is empty", "table", tableName)
		report.QualityScore = 1.0 // Empty table is technically "clean"
		return report, nil
	}

	// 2. Analyze each column
	totalMissingValues := 0
	for _, col := range schema {
		colReport, err := dq.analyzeColumn(ctx, tableName, col, totalRows)
		if err != nil {
			dq.logger.Warn("Failed to analyze column", "column", col.Name, "error", err)
			continue
		}
		report.ColumnReports[col.Name] = colReport
		totalMissingValues += colReport.MissingCount

		if colReport.MissingCount > 0 {
			report.ColumnsWithNulls = append(report.ColumnsWithNulls, col.Name)
		}
	}

	// 3. Calculate rows with missing values
	rowsWithMissing, err := dq.countRowsWithMissing(ctx, tableName, schema)
	if err != nil {
		dq.logger.Warn("Failed to count rows with missing values", "error", err)
	}
	report.RowsWithMissing = rowsWithMissing
	report.CompleteRows = totalRows - rowsWithMissing

	// 4. Calculate missing data percentage
	totalCells := totalRows * len(schema)
	if totalCells > 0 {
		report.MissingDataPercent = float64(totalMissingValues) / float64(totalCells) * 100.0
	}

	// 5. Detect duplicates
	duplicates, err := dq.countDuplicates(ctx, tableName, schema)
	if err != nil {
		dq.logger.Warn("Failed to count duplicates", "error", err)
	}
	report.DuplicateRows = duplicates

	// 6. Calculate overall quality score
	report.QualityScore = dq.calculateQualityScore(report)

	dq.logger.Info("Data quality analysis complete",
		"table", tableName,
		"quality_score", report.QualityScore,
		"missing_percent", report.MissingDataPercent,
		"duplicates", report.DuplicateRows,
	)

	return report, nil
}

// countRows returns the total number of rows in the table
func (dq *DataQualityService) countRows(ctx context.Context, tableName string) (int, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	var count int
	err := dq.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// analyzeColumn performs quality analysis on a single column
func (dq *DataQualityService) analyzeColumn(
	ctx context.Context,
	tableName string,
	col models.ColumnInfo,
	totalRows int,
) (*ColumnQualityReport, error) {

	report := &ColumnQualityReport{
		ColumnName: col.Name,
		DataType:   col.Type,
	}

	// Count nulls and empty values
	query := fmt.Sprintf(
		"SELECT COUNT(*) FROM %s WHERE %s IS NULL OR TRIM(CAST(%s AS TEXT)) = ''",
		tableName, col.Name, col.Name,
	)
	err := dq.db.QueryRowContext(ctx, query).Scan(&report.MissingCount)
	if err != nil {
		return nil, fmt.Errorf("failed to count missing values: %w", err)
	}

	if totalRows > 0 {
		report.MissingPercent = float64(report.MissingCount) / float64(totalRows) * 100.0
	}

	// Count unique values (cardinality)
	query = fmt.Sprintf("SELECT COUNT(DISTINCT %s) FROM %s WHERE %s IS NOT NULL",
		col.Name, tableName, col.Name)
	err = dq.db.QueryRowContext(ctx, query).Scan(&report.UniqueCount)
	if err != nil {
		return nil, fmt.Errorf("failed to count unique values: %w", err)
	}

	if totalRows > 0 {
		report.Cardinality = float64(report.UniqueCount) / float64(totalRows)
	}

	// Type-specific analysis
	switch col.Type {
	case "INTEGER", "BIGINT", "NUMERIC":
		stats, err := dq.calculateNumericStats(ctx, tableName, col.Name)
		if err != nil {
			dq.logger.Warn("Failed to calculate numeric stats", "column", col.Name, "error", err)
		} else {
			report.Statistics = stats
		}
	case "TEXT":
		stats, err := dq.calculateCategoricalStats(ctx, tableName, col.Name)
		if err != nil {
			dq.logger.Warn("Failed to calculate categorical stats", "column", col.Name, "error", err)
		} else {
			report.Statistics = stats
		}
	case "TIMESTAMP":
		stats, err := dq.calculateDateStats(ctx, tableName, col.Name)
		if err != nil {
			dq.logger.Warn("Failed to calculate date stats", "column", col.Name, "error", err)
		} else {
			report.Statistics = stats
		}
	}

	return report, nil
}

// calculateNumericStats calculates statistical metrics for numeric columns
func (dq *DataQualityService) calculateNumericStats(
	ctx context.Context,
	tableName string,
	columnName string,
) (*NumericStatistics, error) {

	stats := &NumericStatistics{}

	query := fmt.Sprintf(`
		SELECT
			AVG(%s::NUMERIC) as mean,
			COALESCE(STDDEV(%s::NUMERIC), 0) as stddev,
			MIN(%s::NUMERIC) as min,
			MAX(%s::NUMERIC) as max,
			PERCENTILE_CONT(0.25) WITHIN GROUP (ORDER BY %s::NUMERIC) as q1,
			PERCENTILE_CONT(0.50) WITHIN GROUP (ORDER BY %s::NUMERIC) as median,
			PERCENTILE_CONT(0.75) WITHIN GROUP (ORDER BY %s::NUMERIC) as q3
		FROM %s
		WHERE %s IS NOT NULL
	`, columnName, columnName, columnName, columnName,
		columnName, columnName, columnName, tableName, columnName)

	err := dq.db.QueryRowContext(ctx, query).Scan(
		&stats.Mean, &stats.StdDev, &stats.Min, &stats.Max,
		&stats.Q1, &stats.Median, &stats.Q3,
	)
	if err != nil {
		return nil, err
	}

	// Detect outliers using IQR method
	iqr := stats.Q3 - stats.Q1
	lowerBound := stats.Q1 - 1.5*iqr
	upperBound := stats.Q3 + 1.5*iqr

	outlierQuery := fmt.Sprintf(`
		SELECT %s FROM %s
		WHERE %s IS NOT NULL AND (%s::NUMERIC < %f OR %s::NUMERIC > %f)
		LIMIT 10
	`, columnName, tableName, columnName, columnName, lowerBound, columnName, upperBound)

	rows, err := dq.db.QueryContext(ctx, outlierQuery)
	if err == nil {
		defer rows.Close()
		stats.Outliers = []interface{}{}
		for rows.Next() {
			var val interface{}
			if err := rows.Scan(&val); err == nil {
				stats.Outliers = append(stats.Outliers, val)
			}
		}
	}

	return stats, nil
}

// calculateCategoricalStats calculates distribution metrics for text columns
func (dq *DataQualityService) calculateCategoricalStats(
	ctx context.Context,
	tableName string,
	columnName string,
) (*CategoricalStatistics, error) {

	stats := &CategoricalStatistics{
		Distribution: make(map[string]int),
	}

	// Get top 10 most frequent values
	query := fmt.Sprintf(`
		SELECT %s, COUNT(*) as freq
		FROM %s
		WHERE %s IS NOT NULL AND TRIM(CAST(%s AS TEXT)) != ''
		GROUP BY %s
		ORDER BY freq DESC
		LIMIT 10
	`, columnName, tableName, columnName, columnName, columnName)

	rows, err := dq.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var value string
		var count int
		if err := rows.Scan(&value, &count); err != nil {
			continue
		}
		stats.Distribution[value] = count

		if stats.MostFrequent == "" {
			stats.MostFrequent = value
			stats.MostFrequentCount = count
		}
	}

	return stats, nil
}

// calculateDateStats calculates metrics for timestamp columns
func (dq *DataQualityService) calculateDateStats(
	ctx context.Context,
	tableName string,
	columnName string,
) (*DateStatistics, error) {

	stats := &DateStatistics{}

	query := fmt.Sprintf(`
		SELECT
			MIN(%s)::TEXT as min_date,
			MAX(%s)::TEXT as max_date
		FROM %s
		WHERE %s IS NOT NULL
	`, columnName, columnName, tableName, columnName)

	err := dq.db.QueryRowContext(ctx, query).Scan(&stats.MinDate, &stats.MaxDate)
	if err != nil {
		return nil, err
	}

	// Get most frequent date
	freqQuery := fmt.Sprintf(`
		SELECT %s::TEXT, COUNT(*) as freq
		FROM %s
		WHERE %s IS NOT NULL
		GROUP BY %s
		ORDER BY freq DESC
		LIMIT 1
	`, columnName, tableName, columnName, columnName)

	err = dq.db.QueryRowContext(ctx, freqQuery).Scan(&stats.MostFrequentDate, new(int))
	if err != nil {
		stats.MostFrequentDate = ""
	}

	stats.DateRange = fmt.Sprintf("%s to %s", stats.MinDate, stats.MaxDate)

	return stats, nil
}

// countRowsWithMissing counts rows that have at least one missing value
func (dq *DataQualityService) countRowsWithMissing(
	ctx context.Context,
	tableName string,
	schema []models.ColumnInfo,
) (int, error) {

	if len(schema) == 0 {
		return 0, nil
	}

	// Build WHERE clause for rows with any null values
	conditions := []string{}
	for _, col := range schema {
		condition := fmt.Sprintf("(%s IS NULL OR TRIM(CAST(%s AS TEXT)) = '')", col.Name, col.Name)
		conditions = append(conditions, condition)
	}

	query := fmt.Sprintf(
		"SELECT COUNT(*) FROM %s WHERE %s",
		tableName,
		joinConditions(conditions, " OR "),
	)

	var count int
	err := dq.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// countDuplicates counts duplicate rows in the table
func (dq *DataQualityService) countDuplicates(
	ctx context.Context,
	tableName string,
	schema []models.ColumnInfo,
) (int, error) {

	if len(schema) == 0 {
		return 0, nil
	}

	// Build column list for grouping
	columnList := []string{}
	for _, col := range schema {
		columnList = append(columnList, col.Name)
	}
	columns := joinConditions(columnList, ", ")

	query := fmt.Sprintf(`
		SELECT COUNT(*) - COUNT(DISTINCT (%s))
		FROM %s
	`, columns, tableName)

	var duplicates int
	err := dq.db.QueryRowContext(ctx, query).Scan(&duplicates)
	if err != nil {
		// If DISTINCT on multiple columns fails, try a different approach
		query = fmt.Sprintf(`
			SELECT SUM(cnt - 1) FROM (
				SELECT COUNT(*) as cnt
				FROM %s
				GROUP BY %s
				HAVING COUNT(*) > 1
			) duplicates
		`, tableName, columns)

		err = dq.db.QueryRowContext(ctx, query).Scan(&duplicates)
		if err != nil {
			return 0, err
		}
	}

	return duplicates, nil
}

// calculateQualityScore calculates an overall quality score (0-1)
func (dq *DataQualityService) calculateQualityScore(report *DataQualityReport) float64 {
	if report.TotalRows == 0 {
		return 1.0
	}

	// Scoring factors:
	// 1. Completeness: (100 - missing_percent) / 100
	completenessScore := (100.0 - report.MissingDataPercent) / 100.0

	// 2. Uniqueness: (total_rows - duplicates) / total_rows
	uniquenessScore := float64(report.TotalRows-report.DuplicateRows) / float64(report.TotalRows)

	// Weighted average: 70% completeness, 30% uniqueness
	qualityScore := (completenessScore * 0.7) + (uniquenessScore * 0.3)

	// Ensure score is between 0 and 1
	if qualityScore < 0 {
		qualityScore = 0
	}
	if qualityScore > 1 {
		qualityScore = 1
	}

	return qualityScore
}

// Helper function to join conditions
func joinConditions(conditions []string, separator string) string {
	result := ""
	for i, cond := range conditions {
		if i > 0 {
			result += separator
		}
		result += cond
	}
	return result
}
