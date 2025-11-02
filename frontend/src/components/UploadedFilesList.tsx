'use client'

import { useState, useEffect } from 'react'
import { File, Trash2, Clock, Database, Copy, CheckCircle } from 'lucide-react'
import { apiClient } from '@/lib/api'
import { UploadedFile } from '@/types'

interface Props {
  refreshTrigger?: number
}

export default function UploadedFilesList({ refreshTrigger }: Props) {
  const [files, setFiles] = useState<UploadedFile[]>([])
  const [loading, setLoading] = useState(true)
  const [copiedTable, setCopiedTable] = useState<string | null>(null)

  useEffect(() => {
    loadFiles()
  }, [refreshTrigger])

  const loadFiles = async () => {
    try {
      setLoading(true)
      const data = await apiClient.getUploadedFiles()
      setFiles(data)
    } catch (error) {
      console.error('Failed to load files:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async (file: UploadedFile) => {
    if (!confirm(`Are you sure you want to delete "${file.original_filename}"?\n\nThis will also delete the table "${file.table_name}" and all its data.`)) {
      return
    }

    try {
      await apiClient.deleteUploadedFile(file.id)
      setFiles(files.filter(f => f.id !== file.id))
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
        const isCopied = copiedTable === file.table_name

        return (
          <div
            key={file.id}
            className="p-4 rounded-lg border transition-all hover:shadow-md"
            style={{ background: 'var(--surface-color)', borderColor: 'var(--border-color)' }}
          >
            <div className="flex items-start justify-between gap-4">
              {/* Left: File Info */}
              <div className="flex items-start gap-3 flex-1">
                <div
                  className="w-12 h-12 rounded-lg flex items-center justify-center flex-shrink-0 mt-1"
                  style={{ background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' }}
                >
                  <File className="w-6 h-6 text-white" />
                </div>
                <div className="flex-1 min-w-0">
                  <p className="font-medium truncate mb-1" style={{ color: 'var(--text-primary)' }}>
                    {file.original_filename}
                  </p>
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
                    <span>•</span>
                    <span className="flex items-center gap-1">
                      <Clock className="w-3 h-3" />
                      {new Date(file.created_at).toLocaleDateString()}
                    </span>
                  </div>

                  {/* Table Name */}
                  <div className="mt-2 flex items-center gap-2">
                    <button
                      onClick={() => handleCopyTable(file.table_name)}
                      className="flex items-center gap-2 px-3 py-1.5 rounded-md text-xs font-mono transition-all"
                      style={{
                        background: 'var(--primary-background)',
                        color: 'var(--text-primary)',
                        border: '1px solid var(--border-color)'
                      }}
                      onMouseEnter={(e) => e.currentTarget.style.borderColor = 'var(--accent-color)'}
                      onMouseLeave={(e) => e.currentTarget.style.borderColor = 'var(--border-color)'}
                      title="Click to copy table name"
                    >
                      {isCopied ? (
                        <>
                          <CheckCircle className="w-3 h-3" style={{ color: '#10b981' }} />
                          <span>Copied!</span>
                        </>
                      ) : (
                        <>
                          <Database className="w-3 h-3" />
                          <span>{file.table_name}</span>
                          <Copy className="w-3 h-3" />
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
              <div className="flex items-center gap-2 flex-shrink-0">
                <span
                  className="text-xs px-3 py-1 rounded-full font-medium"
                  style={{ background: statusStyle.bg, color: statusStyle.color }}
                >
                  {file.status}
                </span>
                <button
                  onClick={() => handleDelete(file)}
                  className="p-2 rounded-lg transition-colors"
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
