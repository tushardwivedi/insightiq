package services

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/xuri/excelize/v2"
	"insightiq/backend/internal/models"
	"insightiq/backend/internal/repository"
)

const (
	BatchSize         = 1000
	MaxRowsToSample   = 100
	MaxColumnNameLen  = 63    // PostgreSQL identifier length limit
	MaxPostgresParams = 32000 // Conservative PostgreSQL parameter limit with large buffer
	MinBatchSize      = 1     // Minimum batch size for very wide tables
)

type FileProcessor struct {
	db                    *sqlx.DB
	fileRepo              *repository.UploadedFileRepository
	connectorRepo         *repository.ConnectorRepository
	dataIngestionService  *FileDataIngestionService
}

func NewFileProcessor(db *sqlx.DB, fileRepo *repository.UploadedFileRepository, connectorRepo *repository.ConnectorRepository, dataIngestionService *FileDataIngestionService) *FileProcessor {
	return &FileProcessor{
		db:                   db,
		fileRepo:             fileRepo,
		connectorRepo:        connectorRepo,
		dataIngestionService: dataIngestionService,
	}
}

// ProcessFile is the main entry point for file processing
func (fp *FileProcessor) ProcessFile(ctx context.Context, file *models.UploadedFile) error {
	slog.Info("Starting file processing", "file_id", file.ID, "filename", file.OriginalFilename)

	// Update status to processing
	file.Status = "processing"
	if err := fp.fileRepo.UpdateStatus(ctx, file.ID, "processing", nil); err != nil {
		slog.Error("Failed to update file status to processing", "error", err)
		return err
	}

	var processErr error
	switch file.FileType {
	case "csv":
		processErr = fp.processCSV(ctx, file)
	case "excel":
		processErr = fp.processExcel(ctx, file)
	default:
		processErr = fmt.Errorf("unsupported file type: %s", file.FileType)
	}

	if processErr != nil {
		slog.Error("File processing failed", "file_id", file.ID, "error", processErr)
		file.Status = "failed"
		errMsg := processErr.Error()
		file.ErrorMessage = &errMsg
		if err := fp.fileRepo.UpdateStatus(ctx, file.ID, "failed", &errMsg); err != nil {
			slog.Error("Failed to update file status to failed", "error", err)
		}
		return processErr
	}

	// Update status to completed
	file.Status = "completed"
	now := time.Now()
	file.ProcessedAt = &now
	if err := fp.fileRepo.UpdateFileMetadata(ctx, file); err != nil {
		slog.Error("Failed to update file status to completed", "error", err)
		return err
	}

	// Create virtual connector for the uploaded file
	if err := fp.createVirtualConnector(ctx, file); err != nil {
		slog.Error("Failed to create virtual connector", "file_id", file.ID, "error", err)
		// Don't fail the processing if connector creation fails - file is still usable
	}

	// Ingest file data into vector store for RAG queries
	if fp.dataIngestionService != nil {
		slog.Info("Starting vector ingestion for file data", "file_id", file.ID)
		if err := fp.dataIngestionService.IngestFileData(ctx, file); err != nil {
			slog.Error("Failed to ingest file data into vector store", "file_id", file.ID, "error", err)
			// Don't fail the processing - SQL queries will still work
		} else {
			slog.Info("Vector ingestion completed successfully", "file_id", file.ID)
		}
	} else {
		slog.Warn("Data ingestion service not available, skipping vector ingestion", "file_id", file.ID)
	}

	slog.Info("File processing completed successfully", "file_id", file.ID, "table_name", file.TableName)
	return nil
}

// processCSV handles CSV file processing
func (fp *FileProcessor) processCSV(ctx context.Context, file *models.UploadedFile) error {
	f, err := os.Open(file.StoragePath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	// Read header
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Sanitize column names
	sanitizedHeaders := make([]string, len(headers))
	for i, h := range headers {
		sanitizedHeaders[i] = fp.sanitizeColumnName(h)
	}

	// Sample data for schema detection
	var sampleData [][]string
	for i := 0; i < MaxRowsToSample; i++ {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read CSV row: %w", err)
		}
		sampleData = append(sampleData, row)
	}

	// Detect schema
	schema := fp.detectSchema(sanitizedHeaders, sampleData)
	file.SchemaJSON = models.FileSchema{Columns: schema}

	// Create table
	if err := fp.createTable(ctx, file.TableName, schema); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	// Reset file reader to beginning (skip header)
	if _, err := f.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to reset file reader: %w", err)
	}
	reader = csv.NewReader(f)
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true
	reader.Read() // skip header

	// Insert data in batches
	rowCount, err := fp.insertCSVData(ctx, file.TableName, sanitizedHeaders, reader)
	if err != nil {
		return fmt.Errorf("failed to insert data: %w", err)
	}

	file.RowCount = rowCount

	// Perform data quality analysis
	slog.Info("Starting data quality analysis", "table", file.TableName)
	qualityService := NewDataQualityService(fp.db, slog.Default())
	qualityReport, err := qualityService.AnalyzeDataQuality(ctx, file.TableName, schema)
	if err != nil {
		slog.Warn("Data quality analysis failed", "error", err)
		// Don't fail the entire processing - continue without quality metrics
	} else {
		// Store quality metrics in file model
		file.DataQualityScore = &qualityReport.QualityScore
		file.MissingDataPercent = &qualityReport.MissingDataPercent
		file.DuplicateRows = &qualityReport.DuplicateRows
		file.RowsWithMissing = &qualityReport.RowsWithMissing

		// Log the quality metrics
		slog.Info("Data quality analysis completed",
			"quality_score", qualityReport.QualityScore,
			"missing_percent", qualityReport.MissingDataPercent,
			"duplicates", qualityReport.DuplicateRows,
			"complete_rows", qualityReport.CompleteRows,
		)
	}

	return nil
}

// processExcel handles Excel file processing
func (fp *FileProcessor) processExcel(ctx context.Context, file *models.UploadedFile) error {
	f, err := excelize.OpenFile(file.StoragePath)
	if err != nil {
		return fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer f.Close()

	// Get first sheet (we'll process only the first sheet for now)
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return fmt.Errorf("no sheets found in Excel file")
	}

	sheetName := sheets[0]
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("failed to read Excel rows: %w", err)
	}

	if len(rows) == 0 {
		return fmt.Errorf("Excel file is empty")
	}

	// First row is header
	headers := rows[0]
	sanitizedHeaders := make([]string, len(headers))
	for i, h := range headers {
		sanitizedHeaders[i] = fp.sanitizeColumnName(h)
	}

	// Sample data for schema detection
	sampleData := rows[1:]
	if len(sampleData) > MaxRowsToSample {
		sampleData = sampleData[:MaxRowsToSample]
	}

	// Detect schema
	schema := fp.detectSchema(sanitizedHeaders, sampleData)
	file.SchemaJSON = models.FileSchema{Columns: schema}

	// Create table
	if err := fp.createTable(ctx, file.TableName, schema); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	// Insert data in batches
	rowCount, err := fp.insertExcelData(ctx, file.TableName, sanitizedHeaders, rows[1:])
	if err != nil {
		return fmt.Errorf("failed to insert data: %w", err)
	}

	file.RowCount = rowCount

	// Perform data quality analysis
	slog.Info("Starting data quality analysis", "table", file.TableName)
	qualityService := NewDataQualityService(fp.db, slog.Default())
	qualityReport, err := qualityService.AnalyzeDataQuality(ctx, file.TableName, schema)
	if err != nil {
		slog.Warn("Data quality analysis failed", "error", err)
		// Don't fail the entire processing - continue without quality metrics
	} else {
		// Store quality metrics in file model
		file.DataQualityScore = &qualityReport.QualityScore
		file.MissingDataPercent = &qualityReport.MissingDataPercent
		file.DuplicateRows = &qualityReport.DuplicateRows
		file.RowsWithMissing = &qualityReport.RowsWithMissing

		// Log the quality metrics
		slog.Info("Data quality analysis completed",
			"quality_score", qualityReport.QualityScore,
			"missing_percent", qualityReport.MissingDataPercent,
			"duplicates", qualityReport.DuplicateRows,
			"complete_rows", qualityReport.CompleteRows,
		)
	}

	return nil
}

// sanitizeColumnName ensures column names are valid PostgreSQL identifiers
func (fp *FileProcessor) sanitizeColumnName(name string) string {
	// Remove leading/trailing whitespace
	name = strings.TrimSpace(name)

	// Replace spaces and special characters with underscores
	reg := regexp.MustCompile(`[^a-zA-Z0-9_]`)
	name = reg.ReplaceAllString(name, "_")

	// Remove consecutive underscores
	reg = regexp.MustCompile(`_+`)
	name = reg.ReplaceAllString(name, "_")

	// Ensure it starts with a letter or underscore
	if len(name) > 0 && !regexp.MustCompile(`^[a-zA-Z_]`).MatchString(name) {
		name = "col_" + name
	}

	// Convert to lowercase
	name = strings.ToLower(name)

	// Truncate if too long
	if len(name) > MaxColumnNameLen {
		name = name[:MaxColumnNameLen]
	}

	// Handle empty names
	if name == "" {
		name = "column"
	}

	return name
}

// detectSchema analyzes sample data to determine column types
func (fp *FileProcessor) detectSchema(headers []string, sampleData [][]string) []models.ColumnInfo {
	schema := make([]models.ColumnInfo, len(headers))

	for i, header := range headers {
		columnType := "TEXT" // Default type
		nullable := true

		// Analyze column values
		hasNonEmpty := false
		allIntegers := true
		allBigInts := true
		allFloats := true
		allBooleans := true
		allDates := true

		for _, row := range sampleData {
			if i >= len(row) {
				continue
			}

			value := strings.TrimSpace(row[i])
			if value == "" {
				continue
			}

			hasNonEmpty = true

			// Check if integer (with range validation)
			if allIntegers {
				if regexp.MustCompile(`^-?\d+$`).MatchString(value) {
					// Additional check: exclude very long numeric strings that are likely IDs
					// Ultra-strict validation: reject anything longer than 8 digits for safety
					if len(value) > 8 {
						allIntegers = false
						slog.Debug("Rejecting as integer due to length", "value", value, "length", len(value))
					} else if intVal, err := strconv.ParseInt(value, 10, 32); err != nil || intVal > 99999999 || intVal < -99999999 {
						allIntegers = false
						slog.Debug("Rejecting as integer due to range", "value", value)
					}
				} else {
					allIntegers = false
				}
			}

			// Check if big integer
			if allBigInts {
				if regexp.MustCompile(`^-?\d+$`).MatchString(value) {
					// Additional check: exclude very long numeric strings that are likely IDs
					// Ultra-strict validation: reject anything longer than 10 digits for safety
					if len(value) > 10 {
						allBigInts = false
						slog.Debug("Rejecting as bigint due to length", "value", value, "length", len(value))
					} else if bigIntVal, err := strconv.ParseInt(value, 10, 64); err != nil || bigIntVal > 9999999999 || bigIntVal < -9999999999 {
						allBigInts = false
						slog.Debug("Rejecting as bigint due to range", "value", value)
					}
				} else {
					allBigInts = false
				}
			}

			// Check if float
			if allFloats && !regexp.MustCompile(`^-?\d+\.?\d*$`).MatchString(value) {
				allFloats = false
			}

			// Check if boolean
			if allBooleans {
				lowerValue := strings.ToLower(value)
				if lowerValue != "true" && lowerValue != "false" && lowerValue != "1" && lowerValue != "0" && lowerValue != "yes" && lowerValue != "no" {
					allBooleans = false
				}
			}

			// Check if date
			if allDates {
				if !fp.isDate(value) {
					allDates = false
				}
			}
		}

		// Determine type based on analysis
		if hasNonEmpty {
			if allIntegers {
				columnType = "INTEGER"
			} else if allBigInts {
				columnType = "BIGINT"
			} else if allFloats {
				columnType = "NUMERIC"
			} else if allBooleans {
				columnType = "BOOLEAN"
			} else if allDates {
				columnType = "TIMESTAMP"
			}
		}

		schema[i] = models.ColumnInfo{
			Name:     header,
			Type:     columnType,
			Nullable: nullable,
		}

		// Log column type detection for debugging
		slog.Debug("Detected column type", "column", header, "type", columnType, "samples", len(sampleData))
	}

	return schema
}

// isDate checks if a string can be parsed as a date
func (fp *FileProcessor) isDate(value string) bool {
	dateFormats := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
		"01/02/2006",
		"02-01-2006",
		time.RFC3339,
	}

	for _, format := range dateFormats {
		if _, err := time.Parse(format, value); err == nil {
			return true
		}
	}
	return false
}

// createTable creates a new table with the detected schema
func (fp *FileProcessor) createTable(ctx context.Context, tableName string, schema []models.ColumnInfo) error {
	var columns []string
	for _, col := range schema {
		nullConstraint := "NULL"
		if !col.Nullable {
			nullConstraint = "NOT NULL"
		}
		columns = append(columns, fmt.Sprintf("%s %s %s", col.Name, col.Type, nullConstraint))
	}

	createSQL := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s (%s)",
		tableName,
		strings.Join(columns, ", "),
	)

	slog.Info("Creating table", "table_name", tableName, "sql", createSQL)

	if _, err := fp.db.ExecContext(ctx, createSQL); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

// insertCSVData inserts CSV data in batches with dynamic batch sizing
func (fp *FileProcessor) insertCSVData(ctx context.Context, tableName string, headers []string, reader *csv.Reader) (int, error) {
	totalRows := 0
	safeBatchSize := fp.calculateSafeBatchSize(len(headers))
	batch := make([][]string, 0, safeBatchSize)

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return totalRows, fmt.Errorf("failed to read CSV row: %w", err)
		}

		batch = append(batch, row)

		if len(batch) >= safeBatchSize {
			if err := fp.insertBatch(ctx, tableName, headers, batch); err != nil {
				return totalRows, err
			}
			totalRows += len(batch)
			batch = batch[:0]
		}
	}

	// Insert remaining rows
	if len(batch) > 0 {
		if err := fp.insertBatch(ctx, tableName, headers, batch); err != nil {
			return totalRows, err
		}
		totalRows += len(batch)
	}

	return totalRows, nil
}

// insertExcelData inserts Excel data in batches with dynamic batch sizing
func (fp *FileProcessor) insertExcelData(ctx context.Context, tableName string, headers []string, rows [][]string) (int, error) {
	totalRows := 0
	safeBatchSize := fp.calculateSafeBatchSize(len(headers))

	for i := 0; i < len(rows); i += safeBatchSize {
		end := i + safeBatchSize
		if end > len(rows) {
			end = len(rows)
		}

		batch := rows[i:end]
		if err := fp.insertBatch(ctx, tableName, headers, batch); err != nil {
			return totalRows, err
		}
		totalRows += len(batch)
	}

	return totalRows, nil
}

// calculateSafeBatchSize calculates batch size to stay under PostgreSQL parameter limit
func (fp *FileProcessor) calculateSafeBatchSize(columnCount int) int {
	// Calculate safe batch size based on column count
	safeBatchSize := MaxPostgresParams / columnCount

	// Ensure minimum batch size
	if safeBatchSize < MinBatchSize {
		safeBatchSize = MinBatchSize
	}

	// Cap at default batch size for normal tables
	if safeBatchSize > BatchSize {
		safeBatchSize = BatchSize
	}

	// Log when using reduced batch size
	if safeBatchSize < BatchSize {
		slog.Warn("Using reduced batch size due to large column count",
			"columns", columnCount,
			"safe_batch_size", safeBatchSize,
			"default_batch_size", BatchSize,
			"max_params_per_batch", safeBatchSize*columnCount,
			"postgres_limit", MaxPostgresParams)
	}

	return safeBatchSize
}

// insertBatch inserts a batch of rows into the table
func (fp *FileProcessor) insertBatch(ctx context.Context, tableName string, headers []string, batch [][]string) error {
	if len(batch) == 0 {
		return nil
	}

	// Validate parameter count before building query
	paramCount := len(batch) * len(headers)
	if paramCount > MaxPostgresParams {
		return fmt.Errorf("batch too large: %d parameters exceeds PostgreSQL limit of %d (batch: %d rows, columns: %d)",
			paramCount, MaxPostgresParams, len(batch), len(headers))
	}

	// Build INSERT statement with placeholders
	placeholders := make([]string, len(batch))
	values := make([]interface{}, 0, len(batch)*len(headers))

	for i, row := range batch {
		rowPlaceholders := make([]string, len(headers))
		for j := range headers {
			var value interface{}
			if j < len(row) && row[j] != "" {
				value = row[j]
			} else {
				value = nil
			}
			values = append(values, value)
			rowPlaceholders[j] = fmt.Sprintf("$%d", i*len(headers)+j+1)
		}
		placeholders[i] = fmt.Sprintf("(%s)", strings.Join(rowPlaceholders, ", "))
	}

	insertSQL := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES %s",
		tableName,
		strings.Join(headers, ", "),
		strings.Join(placeholders, ", "),
	)

	if _, err := fp.db.ExecContext(ctx, insertSQL, values...); err != nil {
		return fmt.Errorf("failed to insert batch: %w", err)
	}

	return nil
}

// ProcessFileAsync processes a file in the background
func (fp *FileProcessor) ProcessFileAsync(file *models.UploadedFile) {
	go func() {
		ctx := context.Background()
		if err := fp.ProcessFile(ctx, file); err != nil {
			slog.Error("Background file processing failed", "file_id", file.ID, "error", err)
		}
	}()
}

// createVirtualConnector creates a virtual file_upload connector for the uploaded file
// This allows the query system to treat uploaded files as queryable data sources
func (fp *FileProcessor) createVirtualConnector(ctx context.Context, file *models.UploadedFile) error {
	// Check if connector already exists for this file
	connectors, err := fp.connectorRepo.GetByType(ctx, models.ConnectorTypeFileUpload)
	if err != nil {
		return fmt.Errorf("failed to get existing file upload connectors: %w", err)
	}

	// Check if this file already has a connector
	for _, conn := range connectors {
		if fileID, ok := conn.Config["file_id"].(string); ok && fileID == file.ID.String() {
			slog.Info("Virtual connector already exists for file", "file_id", file.ID, "connector_id", conn.ID)
			return nil
		}
	}

	// Create connector config with file metadata
	config := models.ConnectorConfig{
		"table_name":        file.TableName,
		"file_id":           file.ID.String(),
		"original_filename": file.OriginalFilename,
		"row_count":         file.RowCount,
		"file_type":         file.FileType,
		"schema":            file.SchemaJSON,
	}

	// Create the virtual connector
	connector := &models.DataConnector{
		Name:   file.OriginalFilename,
		Type:   models.ConnectorTypeFileUpload,
		Status: models.ConnectorStatusConnected,
		Config: config,
	}

	if err := fp.connectorRepo.Create(ctx, connector); err != nil {
		return fmt.Errorf("failed to create virtual connector: %w", err)
	}

	slog.Info("Created virtual connector for uploaded file",
		"file_id", file.ID,
		"connector_id", connector.ID,
		"table_name", file.TableName,
		"filename", file.OriginalFilename)

	return nil
}
