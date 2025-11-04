package http

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"insightiq/backend/internal/models"
	"insightiq/backend/internal/repository"
	"insightiq/backend/internal/services"

	"github.com/google/uuid"
)

const (
	MaxUploadSize = 50 << 20 // 50 MB
	UploadDir     = "./uploads"
)

// FileUploadHandler handles file upload operations
type FileUploadHandler struct {
	fileRepo      *repository.UploadedFileRepository
	fileProcessor *services.FileProcessor
	logger        *slog.Logger
}

// NewFileUploadHandler creates a new FileUploadHandler
func NewFileUploadHandler(fileRepo *repository.UploadedFileRepository, fileProcessor *services.FileProcessor, logger *slog.Logger) *FileUploadHandler {
	return &FileUploadHandler{
		fileRepo:      fileRepo,
		fileProcessor: fileProcessor,
		logger:        logger,
	}
}

// HandleUpload handles file upload requests
func (h *FileUploadHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse multipart form with size limit
	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)
	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		h.logger.Error("Failed to parse multipart form", "error", err)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "File too large or invalid form data"})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		h.logger.Error("User ID not found in context")
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		h.logger.Error("Invalid user ID", "error", err)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		return
	}

	// Get file from form
	file, header, err := r.FormFile("file")
	if err != nil {
		h.logger.Error("Failed to read file", "error", err)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Failed to read file"})
		return
	}
	defer file.Close()

	// Validate file type
	fileType, err := detectFileType(header.Filename, header.Header.Get("Content-Type"))
	if err != nil {
		h.logger.Error("Invalid file type", "error", err)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s_%s", uuid.New().String(), sanitizeFilename(header.Filename))
	storagePath := filepath.Join(UploadDir, filename)

	// Ensure upload directory exists
	if err := os.MkdirAll(UploadDir, 0755); err != nil {
		h.logger.Error("Failed to create upload directory", "error", err)
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create upload directory"})
		return
	}

	// Save file to disk
	dst, err := os.Create(storagePath)
	if err != nil {
		h.logger.Error("Failed to create destination file", "error", err)
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to save file"})
		return
	}
	defer dst.Close()

	fileSize, err := io.Copy(dst, file)
	if err != nil {
		os.Remove(storagePath) // Cleanup
		h.logger.Error("Failed to write file", "error", err)
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to save file"})
		return
	}

	// Generate table name
	tableName := generateTableName(header.Filename)

	// Create file record
	uploadedFile := &models.UploadedFile{
		UserID:           userUUID,
		Filename:         filename,
		OriginalFilename: header.Filename,
		FileSize:         fileSize,
		FileType:         fileType,
		MimeType:         header.Header.Get("Content-Type"),
		StoragePath:      storagePath,
		TableName:        tableName,
		SchemaJSON:       models.FileSchema{Columns: []models.ColumnInfo{}},
		Status:           "pending",
		Source:           "local",
	}

	if err := h.fileRepo.Create(ctx, uploadedFile); err != nil {
		os.Remove(storagePath) // Cleanup
		h.logger.Error("Failed to create file record", "error", err)
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create file record"})
		return
	}

	h.logger.Info("File uploaded successfully", "file_id", uploadedFile.ID, "filename", uploadedFile.OriginalFilename)

	// Trigger async file processing
	h.fileProcessor.ProcessFileAsync(uploadedFile)

	respondJSON(w, http.StatusCreated, uploadedFile)
}

// HandleList lists all uploaded files for the current user
func (h *FileUploadHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user ID from context
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		h.logger.Error("User ID not found in context")
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		h.logger.Error("Invalid user ID", "error", err)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		return
	}

	files, err := h.fileRepo.ListByUser(ctx, userUUID)
	if err != nil {
		h.logger.Error("Failed to list files", "error", err)
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to list files"})
		return
	}

	respondJSON(w, http.StatusOK, files)
}

// HandleGetByID retrieves a specific uploaded file
func (h *FileUploadHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get file ID from URL path
	path := r.URL.Path
	parts := strings.Split(strings.TrimPrefix(path, "/api/files/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		h.logger.Error("File ID not provided in URL")
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "File ID required"})
		return
	}

	fileIDStr := parts[0]
	fileID, err := uuid.Parse(fileIDStr)
	if err != nil {
		h.logger.Error("Invalid file ID", "error", err)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid file ID"})
		return
	}

	file, err := h.fileRepo.GetByID(ctx, fileID)
	if err != nil {
		h.logger.Error("Failed to get file", "error", err, "file_id", fileID)
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "File not found"})
		return
	}

	// Get user ID from context and verify ownership
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		h.logger.Error("User ID not found in context")
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	if file.UserID.String() != userID {
		h.logger.Warn("User attempting to access file they don't own", "user_id", userID, "file_id", fileID)
		respondJSON(w, http.StatusForbidden, map[string]string{"error": "Access denied"})
		return
	}

	respondJSON(w, http.StatusOK, file)
}

// HandleGetIngestionStatus retrieves the ingestion status of an uploaded file
func (h *FileUploadHandler) HandleGetIngestionStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get file ID from URL path
	path := r.URL.Path
	parts := strings.Split(strings.TrimPrefix(path, "/api/files/"), "/")
	if len(parts) < 2 || parts[0] == "" {
		h.logger.Error("File ID not provided in URL")
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "File ID required"})
		return
	}

	fileIDStr := parts[0]
	fileID, err := uuid.Parse(fileIDStr)
	if err != nil {
		h.logger.Error("Invalid file ID", "error", err)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid file ID"})
		return
	}

	file, err := h.fileRepo.GetByID(ctx, fileID)
	if err != nil {
		h.logger.Error("Failed to get file", "error", err, "file_id", fileID)
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "File not found"})
		return
	}

	// Get user ID from context and verify ownership
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		h.logger.Error("User ID not found in context")
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	if file.UserID.String() != userID {
		h.logger.Warn("User attempting to access ingestion status for file they don't own", "user_id", userID, "file_id", fileID)
		respondJSON(w, http.StatusForbidden, map[string]string{"error": "Access denied"})
		return
	}

	// Return ingestion status information
	statusResponse := map[string]interface{}{
		"file_id":              file.ID,
		"filename":             file.OriginalFilename,
		"status":               file.Status,
		"ingestion_status":     file.IngestionStatus,
		"ingestion_progress":   file.IngestionProgress,
		"vector_count":         file.VectorCount,
		"ingestion_started_at": file.IngestionStartedAt,
		"ingestion_completed_at": file.IngestionCompletedAt,
		"row_count":            file.RowCount,
	}

	respondJSON(w, http.StatusOK, statusResponse)
}

// HandleDelete deletes an uploaded file
func (h *FileUploadHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get file ID from URL path
	path := r.URL.Path
	parts := strings.Split(strings.TrimPrefix(path, "/api/files/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		h.logger.Error("File ID not provided in URL")
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "File ID required"})
		return
	}

	fileIDStr := parts[0]
	fileID, err := uuid.Parse(fileIDStr)
	if err != nil {
		h.logger.Error("Invalid file ID", "error", err)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid file ID"})
		return
	}

	// Get file to verify ownership and get storage path
	file, err := h.fileRepo.GetByID(ctx, fileID)
	if err != nil {
		h.logger.Error("Failed to get file", "error", err, "file_id", fileID)
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "File not found"})
		return
	}

	// Verify ownership
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		h.logger.Error("User ID not found in context")
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	if file.UserID.String() != userID {
		h.logger.Warn("User attempting to delete file they don't own", "user_id", userID, "file_id", fileID)
		respondJSON(w, http.StatusForbidden, map[string]string{"error": "Access denied"})
		return
	}

	// Delete file from disk
	if err := os.Remove(file.StoragePath); err != nil && !os.IsNotExist(err) {
		h.logger.Warn("Failed to delete file from disk", "error", err, "path", file.StoragePath)
		// Continue with database deletion even if file deletion fails
	}

	// Delete from database (this will also drop the table)
	if err := h.fileRepo.Delete(ctx, fileID); err != nil {
		h.logger.Error("Failed to delete file record", "error", err, "file_id", fileID)
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete file"})
		return
	}

	h.logger.Info("File deleted successfully", "file_id", fileID, "filename", file.OriginalFilename)
	respondJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// Helper functions

func detectFileType(filename, mimeType string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".csv":
		return "csv", nil
	case ".xlsx", ".xls":
		return "excel", nil
	default:
		return "", fmt.Errorf("unsupported file type: %s. Only CSV and Excel files are allowed", ext)
	}
}

func sanitizeFilename(filename string) string {
	// Remove path separators and other dangerous characters
	filename = filepath.Base(filename)
	filename = strings.ReplaceAll(filename, " ", "_")
	// Remove any characters that aren't alphanumeric, underscore, dash, or dot
	var safe strings.Builder
	for _, r := range filename {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' || r == '.' {
			safe.WriteRune(r)
		}
	}
	result := safe.String()
	if result == "" || result == "." {
		result = "file"
	}
	return result
}

func generateTableName(filename string) string {
	// Generate safe table name from filename
	base := strings.TrimSuffix(filename, filepath.Ext(filename))
	base = strings.ToLower(base)
	base = strings.ReplaceAll(base, " ", "_")
	base = strings.ReplaceAll(base, "-", "_")

	// Remove non-alphanumeric characters
	var safe strings.Builder
	for _, r := range base {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			safe.WriteRune(r)
		}
	}

	result := safe.String()
	if result == "" || result == "_" {
		result = "file"
	}

	// Ensure it starts with a letter (SQL requirement)
	if len(result) > 0 && result[0] >= '0' && result[0] <= '9' {
		result = "t_" + result
	}

	// Add timestamp to ensure uniqueness
	return fmt.Sprintf("uploaded_%s_%d", result, time.Now().Unix())
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
