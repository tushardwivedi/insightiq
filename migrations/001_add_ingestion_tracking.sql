-- Migration: Add ingestion tracking fields to uploaded_files table
-- Created: 2025-01-04

-- Add new columns for ingestion tracking
ALTER TABLE uploaded_files
ADD COLUMN IF NOT EXISTS ingestion_status VARCHAR(50) NOT NULL DEFAULT 'pending',
ADD COLUMN IF NOT EXISTS ingestion_progress INTEGER NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS vector_count INTEGER NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS ingestion_started_at TIMESTAMP WITH TIME ZONE,
ADD COLUMN IF NOT EXISTS ingestion_completed_at TIMESTAMP WITH TIME ZONE;

-- Add index on ingestion_status for efficient queries
CREATE INDEX IF NOT EXISTS idx_uploaded_files_ingestion_status ON uploaded_files(ingestion_status);

-- Add comments to document the fields
COMMENT ON COLUMN uploaded_files.ingestion_status IS 'Status of vector ingestion: pending, in_progress, completed, failed';
COMMENT ON COLUMN uploaded_files.ingestion_progress IS 'Percentage of ingestion completion (0-100)';
COMMENT ON COLUMN uploaded_files.vector_count IS 'Number of vectors successfully ingested into Qdrant';
COMMENT ON COLUMN uploaded_files.ingestion_started_at IS 'Timestamp when vector ingestion started';
COMMENT ON COLUMN uploaded_files.ingestion_completed_at IS 'Timestamp when vector ingestion completed';
