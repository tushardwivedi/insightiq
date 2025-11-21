'use client'

import { useState, useEffect, useRef } from 'react'
import { File, Trash2, Clock, Database, Copy, CheckCircle, Loader2, Sparkles } from 'lucide-react'
import { apiClient } from '@/lib/api'
import { UploadedFile } from '@/types'

interface Props {
  refreshTrigger?: number
  onCountChange?: (count: number) => void
}

export default function UploadedFilesList({ refreshTrigger, onCountChange }: Props) {
  const [files, setFiles] = useState<UploadedFile[]>([])
  const [loading, setLoading] = useState(true)
  const [copiedTable, setCopiedTable] = useState<string | null>(null)
  const pollIntervalRef = useRef<NodeJS.Timeout | null>(null)

  useEffect(() => {
    loadFiles()

    // Start polling for ingestion status updates
    startPolling()

    return () => {
      // Cleanup polling on unmount
      if (pollIntervalRef.current) {
        clearInterval(pollIntervalRef.current)
      }
    }
  }, [refreshTrigger])

  const startPolling = () => {
    // Clear any existing interval
    if (pollIntervalRef.current) {
      clearInterval(pollIntervalRef.current)
    }

    // Poll every 3 seconds for files that are in progress
    pollIntervalRef.current = setInterval(() => {
      const inProgressFiles = files.filter(f =>
        f.ingestion_status === 'in_progress' ||
        (f.status === 'processing' && f.ingestion_status === 'pending')
      )

      if (inProgressFiles.length > 0) {
        loadFiles() // Refresh files if any are in progress
      }
    }, 3000)
  }

  const loadFiles = async () => {
    try {
      setLoading(true)
      const data = await apiClient.getUploadedFiles()
      setFiles(data)
      if (onCountChange) {
        onCountChange(data.length)
      }
    } catch (error) {
      console.error('Failed to load files:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async (fileId: string, fileName: string, tableName: string) => {
    if (!confirm(`Are you sure you want to delete "${fileName}"?\n\nThis will also delete the table "${tableName}" and all its data.`)) {
      return
    }

    try {
      await apiClient.deleteUploadedFile(fileId)
      const updatedFiles = files.filter(f => f.id !== fileId)
      setFiles(updatedFiles)
      if (onCountChange) {
        onCountChange(updatedFiles.length)
      }
    } catch (error) {
      console.error('Failed to delete file:', error)
      alert('Failed to delete file. Please try again.')
    }
  }

  const handleCopyTable = (tableName: string) => {
    navigator.clipboard.writeText(tableName)
    setCopiedTable(tableName)
    setTimeout(() => setCopiedTable(null), 2000)
  }

  const formatFileSize = (bytes: number) => {
    if (bytes < 1024) return bytes + ' B'
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed': return { bg: '#10b98120', color: '#10b981' }
      case 'processing': return { bg: '#f59e0b20', color: '#f59e0b' }
      case 'failed': return { bg: '#ef444420', color: '#ef4444' }
      default: return { bg: '#6b728020', color: '#6b7280' }
    }
  }

  const getIngestionStatusColor = (status: string) => {
    switch (status) {
      case 'completed': return { bg: '#10b98120', color: '#10b981' }
      case 'in_progress': return { bg: '#3b82f620', color: '#3b82f6' }
      case 'failed': return { bg: '#ef444420', color: '#ef4444' }
      default: return { bg: '#6b728020', color: '#6b7280' }
    }
  }

  const isRAGReady = (file: UploadedFile) => {
    return file.status === 'completed' && file.ingestion_status === 'completed'
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center p-12">
        <div
          className="animate-spin rounded-full h-8 w-8 border-2 border-t-transparent"
          style={{ borderColor: 'var(--accent-color)', borderTopColor: 'transparent' }}
        />
      </div>
    )
  }

  if (files.length === 0) {
    return (
      <div className="text-center p-12">
        <div
          className="w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-4"
          style={{ background: 'var(--hover-surface)' }}
        >
          <File className="w-8 h-8" style={{ color: 'var(--text-secondary)' }} />
        </div>
        <p className="text-lg font-medium mb-2" style={{ color: 'var(--text-primary)' }}>
          No files uploaded yet
        </p>
        <p className="text-sm" style={{ color: 'var(--text-secondary)' }}>
          Upload a CSV or Excel file to get started
        </p>
      </div>
    )
  }

  return (
    <div className="space-y-3">
      {files.map(file => {
        const statusStyle = getStatusColor(file.status)
        const ingestionStyle = getIngestionStatusColor(file.ingestion_status)
        const isCopied = copiedTable === file.table_name
        const ragReady = isRAGReady(file)

        return (
          <div
            key={file.id}
            className="p-4 rounded-lg border transition-all hover:shadow-md"
            style={{ background: 'var(--surface-color)', borderColor: 'var(--border-color)' }}
          >
            <div className="flex items-start justify-between gap-2 min-w-0">
              {/* Left: File Info */}
              <div className="flex items-start gap-2 flex-1 min-w-0 overflow-hidden">
                <div
                  className="w-10 h-10 rounded-lg flex items-center justify-center flex-shrink-0 mt-1"
                  style={{ background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' }}
                >
                  <File className="w-5 h-5 text-white" />
                </div>
                <div className="flex-1 min-w-0 overflow-hidden">
                  <div className="flex items-center gap-1.5 mb-1 min-w-0">
                    <p className="font-medium truncate text-sm flex-1 min-w-0" style={{ color: 'var(--text-primary)' }} title={file.original_filename}>
                      {file.original_filename}
                    </p>
                    {ragReady && (
                      <div
                        className="flex items-center gap-1 px-1.5 py-0.5 rounded-full text-xs font-medium flex-shrink-0 whitespace-nowrap"
                        style={{ background: '#10b98120', color: '#10b981' }}
                        title="RAG queries available"
                      >
                        <Sparkles className="w-3 h-3" />
                        <span className="hidden sm:inline">RAG</span>
                      </div>
                    )}
                  </div>
                  <div className="flex items-center gap-3 text-xs flex-wrap" style={{ color: 'var(--text-secondary)' }}>
                    <span className="flex items-center gap-1">
                      <Database className="w-3 h-3" />
                      {formatFileSize(file.file_size)}
                    </span>
                    <span>•</span>
                    <span>{file.file_type.toUpperCase()}</span>
                    {file.row_count > 0 && (
                      <>
                        <span>•</span>
                        <span>{file.row_count.toLocaleString()} rows</span>
                      </>
                    )}
                    {file.vector_count > 0 && (
                      <>
                        <span>•</span>
                        <span className="flex items-center gap-1">
                          <Sparkles className="w-3 h-3" />
                          {file.vector_count} vectors
                        </span>
                      </>
                    )}
                    <span>•</span>
                    <span className="flex items-center gap-1">
                      <Clock className="w-3 h-3" />
                      {new Date(file.created_at).toLocaleDateString()}
                    </span>
                  </div>

                  {/* Ingestion Progress Bar */}
                  {file.ingestion_status === 'in_progress' && (
                    <div className="mt-3">
                      <div className="flex items-center justify-between text-xs mb-1 gap-2 min-w-0" style={{ color: 'var(--text-secondary)' }}>
                        <span className="flex items-center gap-1 truncate flex-1 min-w-0">
                          <Loader2 className="w-3 h-3 animate-spin flex-shrink-0" />
                          <span className="truncate">Vectorizing...</span>
                        </span>
                        <span className="flex-shrink-0 font-medium">{file.ingestion_progress}%</span>
                      </div>
                      <div className="w-full h-2 rounded-full overflow-hidden" style={{ background: 'var(--hover-surface)' }}>
                        <div
                          className="h-full transition-all duration-300"
                          style={{
                            width: `${file.ingestion_progress}%`,
                            background: 'linear-gradient(90deg, #3b82f6, #8b5cf6)'
                          }}
                        />
                      </div>
                    </div>
                  )}

                  {/* Table Name - For Query Use */}
                  <div className="mt-2">
                    <p className="text-xs mb-1 truncate" style={{ color: 'var(--text-secondary)' }}>
                      Table name:
                    </p>
                    <button
                      onClick={() => handleCopyTable(file.table_name)}
                      className="flex items-center gap-1.5 px-2 py-1.5 rounded-md text-xs font-mono transition-all w-full min-w-0"
                      style={{
                        background: 'var(--primary-background)',
                        color: 'var(--text-primary)',
                        border: '1px solid var(--border-color)',
                        cursor: 'pointer'
                      }}
                      onMouseEnter={(e) => {
                        e.currentTarget.style.borderColor = 'var(--accent-color)'
                        e.currentTarget.style.background = 'var(--hover-surface)'
                      }}
                      onMouseLeave={(e) => {
                        e.currentTarget.style.borderColor = 'var(--border-color)'
                        e.currentTarget.style.background = 'var(--primary-background)'
                      }}
                      title={`Click to copy: ${file.table_name}`}
                    >
                      {isCopied ? (
                        <>
                          <CheckCircle className="w-3 h-3 flex-shrink-0" style={{ color: '#10b981' }} />
                          <span className="truncate flex-1 text-left text-xs" style={{ color: '#10b981' }}>Copied!</span>
                        </>
                      ) : (
                        <>
                          <Database className="w-3 h-3 flex-shrink-0" />
                          <span className="truncate flex-1 text-left min-w-0">{file.table_name}</span>
                          <Copy className="w-3 h-3 flex-shrink-0" />
                        </>
                      )}
                    </button>
                  </div>

                  {/* Error Message */}
                  {file.error_message && (
                    <div className="mt-2 p-2 rounded text-xs" style={{ background: '#ef444420', color: '#ef4444' }}>
                      {file.error_message}
                    </div>
                  )}
                </div>
              </div>

              {/* Right: Status & Actions */}
              <div className="flex items-start gap-1.5 flex-shrink-0">
                <div className="flex flex-col gap-1 items-end">
                  <span
                    className="text-xs px-2 py-0.5 rounded-full font-medium whitespace-nowrap"
                    style={{ background: statusStyle.bg, color: statusStyle.color }}
                  >
                    {file.status}
                  </span>
                  {file.status === 'completed' && file.ingestion_status !== 'pending' && (
                    <span
                      className="text-xs px-2 py-0.5 rounded-full font-medium whitespace-nowrap"
                      style={{ background: ingestionStyle.bg, color: ingestionStyle.color }}
                    >
                      {file.ingestion_status === 'in_progress' ? 'vector' : file.ingestion_status}
                    </span>
                  )}
                </div>
                <button
                  onClick={() => handleDelete(file.id, file.original_filename, file.table_name)}
                  className="p-1.5 rounded-lg transition-colors"
                  style={{ color: '#ef4444' }}
                  onMouseEnter={(e) => e.currentTarget.style.background = 'rgba(239, 68, 68, 0.1)'}
                  onMouseLeave={(e) => e.currentTarget.style.background = 'transparent'}
                  title="Delete file"
                >
                  <Trash2 className="w-4 h-4" />
                </button>
              </div>
            </div>
          </div>
        )
      })}
    </div>
  )
}
