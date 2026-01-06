'use client'

import { useState, useCallback } from 'react'
import { Upload, X, CheckCircle, AlertCircle, File } from 'lucide-react'
import { apiClient } from '@/lib/api'
import { UploadedFile } from '@/types'

interface Props {
  isOpen: boolean
  onClose: () => void
  onUploadComplete?: (file: UploadedFile) => void
}

export default function FileUploadModal({ isOpen, onClose, onUploadComplete }: Props) {
  const [dragActive, setDragActive] = useState(false)
  const [uploading, setUploading] = useState(false)
  const [uploadProgress, setUploadProgress] = useState(0)
  const [uploadStatus, setUploadStatus] = useState<'idle' | 'uploading' | 'success' | 'error'>('idle')
  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [uploadedFile, setUploadedFile] = useState<UploadedFile | null>(null)

  const handleUpload = useCallback(async (file: File) => {
    // Validate file type
    const validTypes = ['.csv', '.xlsx', '.xls']
    const ext = file.name.substring(file.name.lastIndexOf('.')).toLowerCase()
    if (!validTypes.includes(ext)) {
      setUploadStatus('error')
      setErrorMessage('Only CSV and Excel files are supported')
      return
    }

    // Validate file size (50MB max)
    const maxSize = 50 * 1024 * 1024
    if (file.size > maxSize) {
      setUploadStatus('error')
      setErrorMessage('File size must be less than 50MB')
      return
    }

    setUploading(true)
    setUploadStatus('uploading')
    setUploadProgress(0)
    setErrorMessage(null)

    try {
      const formData = new FormData()
      formData.append('file', file)

      // Simulate progress
      const progressInterval = setInterval(() => {
        setUploadProgress(prev => Math.min(prev + 10, 90))
      }, 200)

      const result = await apiClient.uploadFile(formData)

      clearInterval(progressInterval)
      setUploadProgress(100)
      setUploadStatus('success')
      setUploadedFile(result)

      if (onUploadComplete) {
        onUploadComplete(result)
      }

      // Auto-close after 2 seconds
      setTimeout(() => {
        handleClose()
      }, 2000)
    } catch (error) {
      setUploadStatus('error')
      setErrorMessage(error instanceof Error ? error.message : 'Upload failed')
    } finally {
      setUploading(false)
    }
  }, [onUploadComplete])

  const handleDrag = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    if (e.type === 'dragenter' || e.type === 'dragover') {
      setDragActive(true)
    } else if (e.type === 'dragleave') {
      setDragActive(false)
    }
  }, [])

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setDragActive(false)

    const files = e.dataTransfer.files
    if (files && files.length > 0) {
      handleUpload(files[0])
    }
  }, [handleUpload])

  const handleFileSelect = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files
    if (files && files.length > 0) {
      handleUpload(files[0])
    }
  }, [handleUpload])

  const handleClose = () => {
    setUploadStatus('idle')
    setUploadProgress(0)
    setErrorMessage(null)
    setUploadedFile(null)
    onClose()
  }

  if (!isOpen) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      {/* Backdrop */}
      <div
        className="absolute inset-0 bg-black/50 backdrop-blur-sm"
        onClick={handleClose}
      />

      {/* Modal */}
      <div
        className="relative w-full max-w-2xl mx-4 rounded-xl shadow-2xl"
        style={{ background: 'var(--surface-color)', border: '1px solid var(--border-color)' }}
      >
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b" style={{ borderColor: 'var(--border-color)' }}>
          <div className="flex items-center gap-3">
            <div
              className="w-10 h-10 rounded-lg flex items-center justify-center"
              style={{ background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' }}
            >
              <Upload className="w-5 h-5 text-white" />
            </div>
            <div>
              <h2 className="text-xl font-semibold" style={{ color: 'var(--text-primary)' }}>
                Upload Data File
              </h2>
              <p className="text-sm" style={{ color: 'var(--text-secondary)' }}>
                Upload CSV or Excel files to query with SQL
              </p>
            </div>
          </div>
          <button
            onClick={handleClose}
            className="p-2 rounded-lg transition-colors hover:bg-opacity-10"
            style={{ color: 'var(--text-secondary)' }}
            onMouseEnter={(e) => e.currentTarget.style.background = 'var(--hover-surface)'}
            onMouseLeave={(e) => e.currentTarget.style.background = 'transparent'}
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Content */}
        <div className="p-6">
          <div
            className={`relative border-2 border-dashed rounded-lg p-12 transition-all ${
              dragActive ? 'border-accent scale-[1.02]' : 'border'
            }`}
            style={{
              borderColor: dragActive ? 'var(--accent-color)' : 'var(--border-color)',
              background: dragActive ? 'var(--hover-surface)' : 'transparent'
            }}
            onDragEnter={handleDrag}
            onDragLeave={handleDrag}
            onDragOver={handleDrag}
            onDrop={handleDrop}
          >
            {uploadStatus === 'idle' && (
              <div className="flex flex-col items-center gap-4 text-center">
                <div
                  className="w-20 h-20 rounded-full flex items-center justify-center"
                  style={{ background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' }}
                >
                  <Upload className="w-10 h-10 text-white" />
                </div>
                <div>
                  <p className="text-lg font-semibold mb-2" style={{ color: 'var(--text-primary)' }}>
                    Drop your file here or click to browse
                  </p>
                  <p className="text-sm mb-4" style={{ color: 'var(--text-secondary)' }}>
                    Supports CSV (.csv) and Excel (.xlsx, .xls) files up to 50MB
                  </p>
                </div>
                <label
                  className="px-6 py-3 rounded-lg font-medium cursor-pointer inline-flex items-center gap-2 transition-all hover:shadow-lg"
                  style={{ background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)', color: 'white' }}
                >
                  <File className="w-5 h-5" />
                  Browse Files
                  <input
                    type="file"
                    className="hidden"
                    accept=".csv,.xlsx,.xls"
                    onChange={handleFileSelect}
                    disabled={uploading}
                  />
                </label>
              </div>
            )}

            {uploadStatus === 'uploading' && (
              <div className="flex flex-col items-center gap-4">
                <div
                  className="animate-spin rounded-full h-16 w-16 border-4 border-t-transparent"
                  style={{ borderColor: 'var(--accent-color)', borderTopColor: 'transparent' }}
                />
                <div className="w-full max-w-md">
                  <div className="flex justify-between mb-2">
                    <span className="text-sm font-medium" style={{ color: 'var(--text-primary)' }}>
                      Uploading file...
                    </span>
                    <span className="text-sm font-medium" style={{ color: 'var(--text-secondary)' }}>
                      {uploadProgress}%
                    </span>
                  </div>
                  <div className="w-full rounded-full h-3" style={{ background: 'var(--border-color)' }}>
                    <div
                      className="h-3 rounded-full transition-all duration-300"
                      style={{
                        width: `${uploadProgress}%`,
                        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)'
                      }}
                    />
                  </div>
                </div>
              </div>
            )}

            {uploadStatus === 'success' && uploadedFile && (
              <div className="flex flex-col items-center gap-4 text-center">
                <div
                  className="w-20 h-20 rounded-full flex items-center justify-center"
                  style={{ background: '#10b98120' }}
                >
                  <CheckCircle className="w-12 h-12" style={{ color: '#10b981' }} />
                </div>
                <div className="w-full">
                  <p className="text-xl font-semibold mb-2" style={{ color: 'var(--text-primary)' }}>
                    Upload Complete!
                  </p>
                  <p className="text-sm mb-4" style={{ color: 'var(--text-secondary)' }}>
                    {uploadedFile.original_filename} uploaded successfully
                  </p>
                  <div
                    className="inline-block px-4 py-2 rounded-lg text-sm font-mono mb-4"
                    style={{ background: 'var(--primary-background)', color: 'var(--text-primary)' }}
                  >
                    Table: <span style={{ color: 'var(--accent-color)' }}>{uploadedFile.table_name}</span>
                  </div>

                  {/* Data Quality Metrics */}
                  {uploadedFile.data_quality_score !== undefined && (
                    <div
                      className="mt-4 p-4 rounded-lg border w-full max-w-md mx-auto"
                      style={{ background: 'var(--surface-color)', borderColor: 'var(--border-color)' }}
                    >
                      <h4 className="font-semibold mb-3" style={{ color: 'var(--text-primary)' }}>
                        Data Quality Analysis
                      </h4>
                      <div className="grid grid-cols-2 gap-3 text-sm">
                        <div className="text-left">
                          <span style={{ color: 'var(--text-secondary)' }}>Quality Score:</span>
                          <span
                            className="ml-2 font-bold"
                            style={{
                              color: uploadedFile.data_quality_score >= 0.8 ? '#10b981' : uploadedFile.data_quality_score >= 0.6 ? '#f59e0b' : '#ef4444'
                            }}
                          >
                            {(uploadedFile.data_quality_score * 100).toFixed(0)}%
                          </span>
                        </div>
                        {uploadedFile.missing_data_percent !== undefined && (
                          <div className="text-left">
                            <span style={{ color: 'var(--text-secondary)' }}>Missing Data:</span>
                            <span className="ml-2 font-semibold" style={{ color: 'var(--text-primary)' }}>
                              {uploadedFile.missing_data_percent.toFixed(1)}%
                            </span>
                          </div>
                        )}
                        {uploadedFile.duplicate_rows !== undefined && (
                          <div className="text-left">
                            <span style={{ color: 'var(--text-secondary)' }}>Duplicates:</span>
                            <span className="ml-2 font-semibold" style={{ color: 'var(--text-primary)' }}>
                              {uploadedFile.duplicate_rows}
                            </span>
                          </div>
                        )}
                        {uploadedFile.rows_with_missing !== undefined && (
                          <div className="text-left">
                            <span style={{ color: 'var(--text-secondary)' }}>Complete Rows:</span>
                            <span className="ml-2 font-semibold" style={{ color: 'var(--text-primary)' }}>
                              {uploadedFile.row_count - uploadedFile.rows_with_missing} / {uploadedFile.row_count}
                            </span>
                          </div>
                        )}
                      </div>
                    </div>
                  )}
                </div>
              </div>
            )}

            {uploadStatus === 'error' && (
              <div className="flex flex-col items-center gap-4 text-center">
                <div
                  className="w-20 h-20 rounded-full flex items-center justify-center"
                  style={{ background: '#ef444420' }}
                >
                  <AlertCircle className="w-12 h-12" style={{ color: '#ef4444' }} />
                </div>
                <div>
                  <p className="text-xl font-semibold mb-2" style={{ color: 'var(--text-primary)' }}>
                    Upload Failed
                  </p>
                  <p className="text-sm mb-4" style={{ color: '#ef4444' }}>
                    {errorMessage}
                  </p>
                  <button
                    onClick={() => setUploadStatus('idle')}
                    className="px-6 py-3 rounded-lg font-medium transition-colors"
                    style={{ background: 'var(--hover-surface)', color: 'var(--text-primary)' }}
                    onMouseEnter={(e) => e.currentTarget.style.background = 'var(--border-color)'}
                    onMouseLeave={(e) => e.currentTarget.style.background = 'var(--hover-surface)'}
                  >
                    Try Again
                  </button>
                </div>
              </div>
            )}
          </div>

          {/* Info */}
          {uploadStatus === 'idle' && (
            <div className="mt-6 p-4 rounded-lg" style={{ background: 'var(--primary-background)' }}>
              <p className="text-xs font-medium mb-2" style={{ color: 'var(--text-secondary)' }}>
                📋 What happens after upload:
              </p>
              <ul className="text-xs space-y-1" style={{ color: 'var(--text-secondary)' }}>
                <li>• Your file will be stored securely in the database</li>
                <li>• A unique table will be created for your data</li>
                <li>• You can query it using SQL or natural language</li>
                <li>• Files are private and only visible to you</li>
              </ul>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
