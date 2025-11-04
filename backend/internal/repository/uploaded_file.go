package repository

import (
	"context"
	"database/sql"
	"fmt"

	"insightiq/backend/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// UploadedFileRepository handles database operations for uploaded files
type UploadedFileRepository struct {
	db *sqlx.DB
}

// NewUploadedFileRepository creates a new UploadedFileRepository
func NewUploadedFileRepository(db *sqlx.DB) *UploadedFileRepository {
	return &UploadedFileRepository{db: db}
}

// Create creates a new uploaded file record
func (r *UploadedFileRepository) Create(ctx context.Context, file *models.UploadedFile) error {
	query := `
		INSERT INTO uploaded_files (
			user_id, filename, original_filename, file_size, file_type,
			mime_type, storage_path, table_name, schema_json, status, source, drive_file_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		file.UserID, file.Filename, file.OriginalFilename, file.FileSize,
		file.FileType, file.MimeType, file.StoragePath, file.TableName,
		file.SchemaJSON, file.Status, file.Source, file.DriveFileID,
	).Scan(&file.ID, &file.CreatedAt, &file.UpdatedAt)
}

// GetByID retrieves an uploaded file by ID
func (r *UploadedFileRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.UploadedFile, error) {
	var file models.UploadedFile
	query := `SELECT * FROM uploaded_files WHERE id = $1`
	err := r.db.GetContext(ctx, &file, query, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("file not found")
	}
	return &file, err
}

// ListByUser retrieves all uploaded files for a user
func (r *UploadedFileRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]models.UploadedFile, error) {
	var files []models.UploadedFile
	query := `SELECT * FROM uploaded_files WHERE user_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &files, query, userID)
	if err != nil {
		return nil, err
	}
	// Return empty slice instead of nil if no results
	if files == nil {
		files = []models.UploadedFile{}
	}
	return files, nil
}

// UpdateStatus updates the status of an uploaded file
func (r *UploadedFileRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string, errorMsg *string) error {
	query := `UPDATE uploaded_files SET status = $1, error_message = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, status, errorMsg, id)
	return err
}

// UpdateFileMetadata updates the file metadata after processing
func (r *UploadedFileRepository) UpdateFileMetadata(ctx context.Context, file *models.UploadedFile) error {
	query := `
		UPDATE uploaded_files
		SET schema_json = $1, row_count = $2, status = $3, processed_at = $4, updated_at = NOW()
		WHERE id = $5
	`
	_, err := r.db.ExecContext(ctx, query, file.SchemaJSON, file.RowCount, file.Status, file.ProcessedAt, file.ID)
	return err
}

// Delete deletes an uploaded file record and its associated table
func (r *UploadedFileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// First get the file to know the table name
	file, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Drop the data table if it exists
	if file.TableName != "" {
		dropQuery := fmt.Sprintf("DROP TABLE IF EXISTS %s", file.TableName)
		_, err = tx.ExecContext(ctx, dropQuery)
		if err != nil {
			return fmt.Errorf("failed to drop table: %w", err)
		}
	}

	// Delete the file record
	deleteQuery := `DELETE FROM uploaded_files WHERE id = $1`
	_, err = tx.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete file record: %w", err)
	}

	return tx.Commit()
}

// GetByTableName retrieves an uploaded file by its table name
func (r *UploadedFileRepository) GetByTableName(ctx context.Context, tableName string) (*models.UploadedFile, error) {
	var file models.UploadedFile
	query := `SELECT * FROM uploaded_files WHERE table_name = $1`
	err := r.db.GetContext(ctx, &file, query, tableName)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("file not found")
	}
	return &file, err
}

// UpdateIngestionStatus updates the ingestion status and related fields
func (r *UploadedFileRepository) UpdateIngestionStatus(ctx context.Context, id uuid.UUID, status string, progress int, vectorCount int) error {
	query := `
		UPDATE uploaded_files
		SET ingestion_status = $1,
		    ingestion_progress = $2,
		    vector_count = $3,
		    updated_at = NOW()
		WHERE id = $4
	`
	_, err := r.db.ExecContext(ctx, query, status, progress, vectorCount, id)
	return err
}

// StartIngestion marks the start of ingestion process
func (r *UploadedFileRepository) StartIngestion(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE uploaded_files
		SET ingestion_status = 'in_progress',
		    ingestion_progress = 0,
		    ingestion_started_at = NOW(),
		    updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// CompleteIngestion marks the completion of ingestion process
func (r *UploadedFileRepository) CompleteIngestion(ctx context.Context, id uuid.UUID, vectorCount int) error {
	query := `
		UPDATE uploaded_files
		SET ingestion_status = 'completed',
		    ingestion_progress = 100,
		    vector_count = $1,
		    ingestion_completed_at = NOW(),
		    updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, vectorCount, id)
	return err
}

// FailIngestion marks the ingestion as failed with an error message
func (r *UploadedFileRepository) FailIngestion(ctx context.Context, id uuid.UUID, errorMsg string) error {
	query := `
		UPDATE uploaded_files
		SET ingestion_status = 'failed',
		    error_message = $1,
		    updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, errorMsg, id)
	return err
}
