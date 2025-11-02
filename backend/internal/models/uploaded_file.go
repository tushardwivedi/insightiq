package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// UploadedFile represents a file uploaded by a user
type UploadedFile struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	UserID           uuid.UUID  `json:"user_id" db:"user_id"`
	Filename         string     `json:"filename" db:"filename"`
	OriginalFilename string     `json:"original_filename" db:"original_filename"`
	FileSize         int64      `json:"file_size" db:"file_size"`
	FileType         string     `json:"file_type" db:"file_type"` // csv, excel
	MimeType         string     `json:"mime_type" db:"mime_type"`
	StoragePath      string     `json:"storage_path" db:"storage_path"`
	TableName        string     `json:"table_name" db:"table_name"`
	SchemaJSON       FileSchema `json:"schema_json" db:"schema_json"`
	RowCount         int        `json:"row_count" db:"row_count"`
	Status           string     `json:"status" db:"status"` // pending, processing, completed, failed
	ErrorMessage     *string    `json:"error_message,omitempty" db:"error_message"`
	Source           string     `json:"source" db:"source"` // local, google_drive
	DriveFileID      *string    `json:"drive_file_id,omitempty" db:"drive_file_id"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
	ProcessedAt      *time.Time `json:"processed_at,omitempty" db:"processed_at"`
}

// FileSchema represents the schema of an uploaded file
type FileSchema struct {
	Columns []ColumnInfo `json:"columns"`
}

// ColumnInfo represents information about a column in the file
type ColumnInfo struct {
	Name     string `json:"name"`
	Type     string `json:"type"` // text, integer, numeric, boolean, date, timestamp
	Nullable bool   `json:"nullable"`
}

// Scan implements the sql.Scanner interface for FileSchema
func (fs *FileSchema) Scan(value interface{}) error {
	if value == nil {
		fs.Columns = []ColumnInfo{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, fs)
}

// Value implements the driver.Valuer interface for FileSchema
func (fs FileSchema) Value() (driver.Value, error) {
	if len(fs.Columns) == 0 {
		return json.Marshal(map[string]interface{}{"columns": []ColumnInfo{}})
	}
	return json.Marshal(fs)
}
