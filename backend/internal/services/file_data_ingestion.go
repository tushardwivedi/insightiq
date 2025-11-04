package services

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"insightiq/backend/internal/embedding"
	"insightiq/backend/internal/models"
	"insightiq/backend/internal/vectorstore"
)

const (
	// Vector ingestion configuration
	MaxRowsPerBatch       = 100  // Batch size for vector upsert
	MaxRowsToEmbed        = 1000 // Maximum rows to embed per file
	SampleRowsForSchema   = 10   // Number of sample rows for schema understanding
	FileDataCollectionName = "file_data_embeddings"
)

// FileDataIngestionService handles ingestion of uploaded file data into vector store
type FileDataIngestionService struct {
	db               *sqlx.DB
	vectorStore      vectorstore.VectorStore
	embeddingService embedding.EmbeddingService
	fileRepo         FileRepository
	logger           *slog.Logger
}

// FileRepository defines the interface for file repository operations
type FileRepository interface {
	StartIngestion(ctx context.Context, id uuid.UUID) error
	UpdateIngestionStatus(ctx context.Context, id uuid.UUID, status string, progress int, vectorCount int) error
	CompleteIngestion(ctx context.Context, id uuid.UUID, vectorCount int) error
	FailIngestion(ctx context.Context, id uuid.UUID, errorMsg string) error
}

// NewFileDataIngestionService creates a new file data ingestion service
func NewFileDataIngestionService(
	db *sqlx.DB,
	vectorStore vectorstore.VectorStore,
	embeddingService embedding.EmbeddingService,
	fileRepo FileRepository,
	logger *slog.Logger,
) *FileDataIngestionService {
	return &FileDataIngestionService{
		db:               db,
		vectorStore:      vectorStore,
		embeddingService: embeddingService,
		fileRepo:         fileRepo,
		logger:           logger.With("service", "file_data_ingestion"),
	}
}

// IngestFileData ingests uploaded file data into vector store for RAG queries
func (s *FileDataIngestionService) IngestFileData(ctx context.Context, file *models.UploadedFile) error {
	s.logger.Info("Starting vector ingestion for uploaded file",
		"file_id", file.ID,
		"filename", file.OriginalFilename,
		"table_name", file.TableName,
		"row_count", file.RowCount)

	startTime := time.Now()

	// Mark ingestion as started
	if s.fileRepo != nil {
		if err := s.fileRepo.StartIngestion(ctx, file.ID); err != nil {
			s.logger.Warn("Failed to update ingestion start status", "error", err)
		}
	}

	totalVectorCount := 0

	// Ensure collection exists
	dimension := s.embeddingService.GetDimension()
	if err := s.vectorStore.CreateCollection(ctx, FileDataCollectionName, dimension); err != nil {
		s.logger.Warn("Failed to create vector collection (may already exist)", "error", err)
	}

	// Step 1: Ingest schema metadata (for schema-aware queries) - 10% progress
	if err := s.ingestSchemaMetadata(ctx, file); err != nil {
		s.logger.Error("Failed to ingest schema metadata", "error", err)
		if s.fileRepo != nil {
			s.fileRepo.FailIngestion(ctx, file.ID, fmt.Sprintf("schema ingestion failed: %v", err))
		}
		return fmt.Errorf("schema ingestion failed: %w", err)
	}
	totalVectorCount++ // 1 vector for schema
	if s.fileRepo != nil {
		s.fileRepo.UpdateIngestionStatus(ctx, file.ID, "in_progress", 10, totalVectorCount)
	}

	// Step 2: Ingest file summary (for high-level queries) - 20% progress
	if err := s.ingestFileSummary(ctx, file); err != nil {
		s.logger.Error("Failed to ingest file summary", "error", err)
		if s.fileRepo != nil {
			s.fileRepo.FailIngestion(ctx, file.ID, fmt.Sprintf("file summary ingestion failed: %v", err))
		}
		return fmt.Errorf("file summary ingestion failed: %w", err)
	}
	totalVectorCount++ // 1 vector for summary
	if s.fileRepo != nil {
		s.fileRepo.UpdateIngestionStatus(ctx, file.ID, "in_progress", 20, totalVectorCount)
	}

	// Step 3: Ingest row-level data (for detailed queries) - 20-100% progress
	rowVectorCount, err := s.ingestRowData(ctx, file)
	if err != nil {
		s.logger.Error("Failed to ingest row data", "error", err)
		if s.fileRepo != nil {
			s.fileRepo.FailIngestion(ctx, file.ID, fmt.Sprintf("row data ingestion failed: %v", err))
		}
		return fmt.Errorf("row data ingestion failed: %w", err)
	}
	totalVectorCount += rowVectorCount

	// Mark ingestion as completed
	if s.fileRepo != nil {
		if err := s.fileRepo.CompleteIngestion(ctx, file.ID, totalVectorCount); err != nil {
			s.logger.Warn("Failed to update ingestion completion status", "error", err)
		}
	}

	duration := time.Since(startTime)
	s.logger.Info("Vector ingestion completed successfully",
		"file_id", file.ID,
		"filename", file.OriginalFilename,
		"total_vectors", totalVectorCount,
		"duration", duration.String())

	return nil
}

// ingestSchemaMetadata creates embeddings for column schema information
func (s *FileDataIngestionService) ingestSchemaMetadata(ctx context.Context, file *models.UploadedFile) error {
	s.logger.Info("Ingesting schema metadata", "file_id", file.ID)

	// Build schema description text
	var schemaBuilder strings.Builder
	schemaBuilder.WriteString(fmt.Sprintf("File: %s\n", file.OriginalFilename))
	schemaBuilder.WriteString(fmt.Sprintf("Table: %s\n", file.TableName))
	schemaBuilder.WriteString(fmt.Sprintf("Total Rows: %d\n", file.RowCount))
	schemaBuilder.WriteString("Schema:\n")

	for _, col := range file.SchemaJSON.Columns {
		schemaBuilder.WriteString(fmt.Sprintf("- %s (%s)\n", col.Name, col.Type))
	}

	schemaText := schemaBuilder.String()

	// Generate embedding
	embeddingResp, err := s.embeddingService.GenerateEmbedding(ctx, schemaText)
	if err != nil {
		return fmt.Errorf("failed to generate schema embedding: %w", err)
	}

	// Create vector with rich metadata
	// Use UUID format for vector ID (Qdrant requires UUID or uint64)
	vectorID := uuid.New().String()
	vector := vectorstore.Vector{
		ID:     vectorID,
		Values: embeddingResp.Embedding,
		Metadata: map[string]interface{}{
			"type":              "schema",
			"file_id":           file.ID.String(),
			"filename":          file.OriginalFilename,
			"table_name":        file.TableName,
			"row_count":         file.RowCount,
			"column_count":      len(file.SchemaJSON.Columns),
			"file_type":         file.FileType,
			"created_at":        file.CreatedAt.Unix(),
			"searchable_text":   schemaText,
			"semantic_id":       fmt.Sprintf("schema_%s", file.ID), // Store semantic ID in metadata
		},
	}

	// Upsert to vector store
	if err := s.vectorStore.UpsertVectors(ctx, FileDataCollectionName, []vectorstore.Vector{vector}); err != nil {
		return fmt.Errorf("failed to upsert schema vector: %w", err)
	}

	s.logger.Info("Schema metadata ingested", "file_id", file.ID, "vector_id", vector.ID)
	return nil
}

// ingestFileSummary creates embeddings for file-level summary statistics
func (s *FileDataIngestionService) ingestFileSummary(ctx context.Context, file *models.UploadedFile) error {
	s.logger.Info("Ingesting file summary", "file_id", file.ID)

	// Query sample data for summary
	query := fmt.Sprintf("SELECT * FROM %s LIMIT %d", file.TableName, SampleRowsForSchema)
	rows, err := s.db.QueryxContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to query sample data: %w", err)
	}
	defer rows.Close()

	// Collect sample rows
	var sampleRows []map[string]interface{}
	for rows.Next() {
		row := make(map[string]interface{})
		if err := rows.MapScan(row); err != nil {
			s.logger.Warn("Failed to scan row", "error", err)
			continue
		}
		sampleRows = append(sampleRows, row)
	}

	if len(sampleRows) == 0 {
		s.logger.Warn("No sample data available for summary")
		return nil
	}

	// Build summary text with sample data
	var summaryBuilder strings.Builder
	summaryBuilder.WriteString(fmt.Sprintf("File Summary: %s\n", file.OriginalFilename))
	summaryBuilder.WriteString(fmt.Sprintf("Contains %d records with %d columns\n", file.RowCount, len(file.SchemaJSON.Columns)))
	summaryBuilder.WriteString("\nColumns: ")

	columnNames := make([]string, 0, len(file.SchemaJSON.Columns))
	for _, col := range file.SchemaJSON.Columns {
		columnNames = append(columnNames, col.Name)
	}
	summaryBuilder.WriteString(strings.Join(columnNames, ", "))
	summaryBuilder.WriteString("\n\nSample Data:\n")

	// Add sample rows as context
	for i, row := range sampleRows {
		if i >= 3 { // Only include first 3 rows in summary
			break
		}
		summaryBuilder.WriteString(fmt.Sprintf("Row %d: ", i+1))
		for key, val := range row {
			summaryBuilder.WriteString(fmt.Sprintf("%s=%v, ", key, val))
		}
		summaryBuilder.WriteString("\n")
	}

	summaryText := summaryBuilder.String()

	// Generate embedding
	embeddingResp, err := s.embeddingService.GenerateEmbedding(ctx, summaryText)
	if err != nil {
		return fmt.Errorf("failed to generate summary embedding: %w", err)
	}

	// Create vector
	// Use UUID format for vector ID (Qdrant requires UUID or uint64)
	vectorID := uuid.New().String()
	vector := vectorstore.Vector{
		ID:     vectorID,
		Values: embeddingResp.Embedding,
		Metadata: map[string]interface{}{
			"type":              "summary",
			"file_id":           file.ID.String(),
			"filename":          file.OriginalFilename,
			"table_name":        file.TableName,
			"row_count":         file.RowCount,
			"column_count":      len(file.SchemaJSON.Columns),
			"sample_row_count":  len(sampleRows),
			"searchable_text":   summaryText,
			"semantic_id":       fmt.Sprintf("summary_%s", file.ID), // Store semantic ID in metadata
		},
	}

	// Upsert to vector store
	if err := s.vectorStore.UpsertVectors(ctx, FileDataCollectionName, []vectorstore.Vector{vector}); err != nil {
		return fmt.Errorf("failed to upsert summary vector: %w", err)
	}

	s.logger.Info("File summary ingested", "file_id", file.ID, "sample_rows", len(sampleRows))
	return nil
}

// ingestRowData creates embeddings for individual rows (chunked strategy)
func (s *FileDataIngestionService) ingestRowData(ctx context.Context, file *models.UploadedFile) (int, error) {
	s.logger.Info("Ingesting row-level data", "file_id", file.ID, "max_rows", MaxRowsToEmbed)

	// Determine how many rows to embed (cap at MaxRowsToEmbed)
	rowsToEmbed := file.RowCount
	if rowsToEmbed > MaxRowsToEmbed {
		s.logger.Info("Limiting row embeddings", "total_rows", file.RowCount, "embedding_rows", MaxRowsToEmbed)
		rowsToEmbed = MaxRowsToEmbed
	}

	// Query data in batches
	offset := 0
	batchNum := 0
	totalVectorsCreated := 0

	for offset < rowsToEmbed {
		limit := MaxRowsPerBatch
		if offset+limit > rowsToEmbed {
			limit = rowsToEmbed - offset
		}

		vectorsInBatch, err := s.ingestRowBatch(ctx, file, offset, limit, batchNum)
		if err != nil {
			s.logger.Error("Failed to ingest row batch", "offset", offset, "limit", limit, "error", err)
			return totalVectorsCreated, fmt.Errorf("row batch ingestion failed at offset %d: %w", offset, err)
		}

		totalVectorsCreated += vectorsInBatch
		offset += limit
		batchNum++

		// Update progress (20% base + 80% based on rows processed)
		// Progress goes from 20% to 100% during row ingestion
		progressPercent := 20 + int(float64(offset)/float64(rowsToEmbed)*80)
		if s.fileRepo != nil {
			s.fileRepo.UpdateIngestionStatus(ctx, file.ID, "in_progress", progressPercent, 2+totalVectorsCreated) // 2 = schema + summary
		}
	}

	s.logger.Info("Row-level data ingestion completed",
		"file_id", file.ID,
		"rows_embedded", rowsToEmbed,
		"batches", batchNum,
		"vectors_created", totalVectorsCreated)

	return totalVectorsCreated, nil
}

// ingestRowBatch ingests a batch of rows into vector store
func (s *FileDataIngestionService) ingestRowBatch(ctx context.Context, file *models.UploadedFile, offset, limit, batchNum int) (int, error) {
	s.logger.Info("Ingesting row batch", "batch", batchNum, "offset", offset, "limit", limit)

	// Query batch of rows
	query := fmt.Sprintf("SELECT * FROM %s OFFSET %d LIMIT %d", file.TableName, offset, limit)
	rows, err := s.db.QueryxContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to query row batch: %w", err)
	}
	defer rows.Close()

	// Collect vectors for batch upsert
	vectors := make([]vectorstore.Vector, 0, limit)
	rowNum := offset

	for rows.Next() {
		row := make(map[string]interface{})
		if err := rows.MapScan(row); err != nil {
			s.logger.Warn("Failed to scan row", "row_num", rowNum, "error", err)
			continue
		}

		// Create searchable text representation of the row
		rowText := s.createRowSearchText(file, row, rowNum)

		// Generate embedding
		embeddingResp, err := s.embeddingService.GenerateEmbedding(ctx, rowText)
		if err != nil {
			s.logger.Warn("Failed to generate embedding for row", "row_num", rowNum, "error", err)
			rowNum++
			continue
		}

		// Create vector with metadata
		// Use UUID format for vector ID (Qdrant requires UUID or uint64)
		vectorID := uuid.New().String()
		vector := vectorstore.Vector{
			ID:     vectorID,
			Values: embeddingResp.Embedding,
			Metadata: map[string]interface{}{
				"type":              "row",
				"file_id":           file.ID.String(),
				"filename":          file.OriginalFilename,
				"table_name":        file.TableName,
				"row_number":        rowNum,
				"batch_number":      batchNum,
				"searchable_text":   rowText,
				"row_data":          row, // Store actual row data for retrieval
				"semantic_id":       fmt.Sprintf("row_%s_%d", file.ID, rowNum), // Store semantic ID in metadata
			},
		}

		vectors = append(vectors, vector)
		rowNum++
	}

	if len(vectors) == 0 {
		s.logger.Warn("No vectors generated for batch", "batch", batchNum)
		return 0, nil
	}

	// Batch upsert to vector store
	if err := s.vectorStore.UpsertVectors(ctx, FileDataCollectionName, vectors); err != nil {
		return 0, fmt.Errorf("failed to upsert row vectors: %w", err)
	}

	s.logger.Info("Row batch ingested successfully",
		"batch", batchNum,
		"vectors_created", len(vectors))

	return len(vectors), nil
}

// createRowSearchText creates a searchable text representation of a row
func (s *FileDataIngestionService) createRowSearchText(file *models.UploadedFile, row map[string]interface{}, rowNum int) string {
	var textBuilder strings.Builder

	textBuilder.WriteString(fmt.Sprintf("Record %d from %s:\n", rowNum+1, file.OriginalFilename))

	// Add column-value pairs
	for _, col := range file.SchemaJSON.Columns {
		if val, exists := row[col.Name]; exists && val != nil {
			textBuilder.WriteString(fmt.Sprintf("%s: %v\n", col.Name, val))
		}
	}

	return textBuilder.String()
}

// DeleteFileVectors removes all vectors associated with a file
// Note: Since we now use UUID-based IDs, deletion requires querying by metadata filter
// This is a TODO for implementation - requires Qdrant scroll/filter API
func (s *FileDataIngestionService) DeleteFileVectors(ctx context.Context, fileID string) error {
	s.logger.Info("Deleting vectors for file (not implemented - requires metadata-based deletion)", "file_id", fileID)

	// TODO: Implement metadata-based deletion using Qdrant filter API
	// Example: DELETE WHERE payload.file_id = fileID
	// For now, vectors will remain in Qdrant but won't be matched in queries
	// since the file_id in metadata won't match any existing files

	s.logger.Warn("Vector deletion not fully implemented - vectors will be orphaned", "file_id", fileID)
	return nil
}
