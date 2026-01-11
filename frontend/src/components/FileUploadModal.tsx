'use client'

import { useState, useCallback, useRef } from 'react'
import { Upload, X, CheckCircle, AlertCircle, File, FileText, Table, Sparkles, ChevronRight } from 'lucide-react'
import { apiClient } from '@/lib/api'
import { UploadedFile } from '@/types'
import { motion, AnimatePresence } from 'framer-motion'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'

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
  const [selectedFile, setSelectedFile] = useState<File | null>(null)

  const inputRef = useRef<HTMLInputElement>(null)

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return '0 Bytes'
    const k = 1024
    const sizes = ['Bytes', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  }

  const getFileIcon = (filename: string) => {
    const ext = filename.split('.').pop()?.toLowerCase()
    if (ext === 'csv') return Table
    if (['xlsx', 'xls'].includes(ext || '')) return FileText
    return File
  }

  const handleUpload = useCallback(async (file: File) => {
    const validTypes = ['.csv', '.xlsx', '.xls']
    const ext = file.name.substring(file.name.lastIndexOf('.')).toLowerCase()
    if (!validTypes.includes(ext)) {
      setUploadStatus('error')
      setErrorMessage('Only CSV and Excel files are supported')
      return
    }

    const maxSize = 50 * 1024 * 1024
    if (file.size > maxSize) {
      setUploadStatus('error')
      setErrorMessage('File size must be less than 50MB')
      return
    }

    setSelectedFile(file)
    setUploading(true)
    setUploadStatus('uploading')
    setUploadProgress(0)
    setErrorMessage(null)

    try {
      const formData = new FormData()
      formData.append('file', file)

      const progressInterval = setInterval(() => {
        setUploadProgress(prev => Math.min(prev + Math.random() * 15, 90))
      }, 150)

      const result = await apiClient.uploadFile(formData)

      clearInterval(progressInterval)
      setUploadProgress(100)
      setUploadStatus('success')
      setUploadedFile(result)

      if (onUploadComplete) {
        onUploadComplete(result)
      }

      setTimeout(() => {
        handleClose()
      }, 2500)
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
    setSelectedFile(null)
    onClose()
  }

  if (!isOpen) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        exit={{ opacity: 0 }}
        className="absolute inset-0 bg-black/60 backdrop-blur-sm"
        onClick={handleClose}
      />

      <motion.div
        initial={{ opacity: 0, scale: 0.95, y: 20 }}
        animate={{ opacity: 1, scale: 1, y: 0 }}
        exit={{ opacity: 0, scale: 0.95, y: 20 }}
        transition={{ type: "spring", duration: 0.4 }}
        className="relative w-full max-w-lg mx-4 rounded-2xl shadow-2xl overflow-hidden"
        style={{ background: 'var(--surface-color)' }}
      >
        <div className="absolute top-0 left-0 right-0 h-1">
          <motion.div
            className="h-full bg-gradient-to-r from-cyan-400 via-blue-500 to-purple-600"
            initial={{ width: 0 }}
            animate={{ width: uploadStatus === 'uploading' ? `${uploadProgress}%` : '100%' }}
            transition={{ duration: 0.3 }}
          />
        </div>

        <div className="flex items-center justify-between p-5 border-b" style={{ borderColor: 'var(--border-color)' }}>
          <div className="flex items-center gap-3">
            <motion.div
              initial={{ scale: 0 }}
              animate={{ scale: 1 }}
              transition={{ type: "spring", delay: 0.1 }}
              className="w-10 h-10 rounded-xl flex items-center justify-center bg-gradient-to-br from-cyan-500 to-blue-600"
            >
              <Upload className="w-5 h-5 text-white" />
            </motion.div>
            <div>
              <h2 className="text-lg font-semibold" style={{ color: 'var(--text-primary)' }}>
                Upload Data File
              </h2>
              <p className="text-sm" style={{ color: 'var(--text-secondary)' }}>
                Supports CSV and Excel files up to 50MB
              </p>
            </div>
          </div>
          <motion.button
            whileHover={{ scale: 1.1 }}
            whileTap={{ scale: 0.9 }}
            onClick={handleClose}
            className="p-2 rounded-lg transition-colors hover:bg-white/10"
            style={{ color: 'var(--text-secondary)' }}
          >
            <X className="w-5 h-5" />
          </motion.button>
        </div>

        <div className="p-5">
          <AnimatePresence mode="wait">
            {uploadStatus === 'idle' && (
              <motion.div
                key="idle"
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -10 }}
              >
                <div
                  className={`relative border-2 border-dashed rounded-xl p-8 transition-all duration-300 cursor-pointer ${
                    dragActive ? 'border-primary scale-[1.02] bg-primary/5' : ''
                  }`}
                  style={{
                    borderColor: dragActive ? 'var(--primary)' : 'var(--border-color)',
                  }}
                  onDragEnter={handleDrag}
                  onDragLeave={handleDrag}
                  onDragOver={handleDrag}
                  onDrop={handleDrop}
                  onClick={() => inputRef.current?.click()}
                  onMouseEnter={(e) => {
                    if (!dragActive) {
                      e.currentTarget.style.borderColor = 'var(--primary)'
                      e.currentTarget.style.background = 'var(--hover-surface)'
                    }
                  }}
                  onMouseLeave={(e) => {
                    if (!dragActive) {
                      e.currentTarget.style.borderColor = 'var(--border-color)'
                      e.currentTarget.style.background = 'transparent'
                    }
                  }}
                >
                  <input
                    ref={inputRef}
                    type="file"
                    className="hidden"
                    accept=".csv,.xlsx,.xls"
                    onChange={handleFileSelect}
                    disabled={uploading}
                  />

                  <div className="flex flex-col items-center text-center">
                    <motion.div
                      whileHover={{ rotate: 180 }}
                      transition={{ duration: 0.3 }}
                      className="w-16 h-16 rounded-2xl bg-gradient-to-br from-primary/20 to-purple-500/20 flex items-center justify-center mb-4"
                    >
                      <File className="w-8 h-8 text-primary" />
                    </motion.div>

                    <p className="text-base font-medium mb-1" style={{ color: 'var(--text-primary)' }}>
                      Drop your file here or click to browse
                    </p>
                    <p className="text-sm mb-4" style={{ color: 'var(--text-secondary)' }}>
                      CSV, XLSX, XLS up to 50MB
                    </p>

                    <motion.button
                      whileHover={{ scale: 1.05 }}
                      whileTap={{ scale: 0.95 }}
                      className="px-5 py-2.5 rounded-lg font-medium flex items-center gap-2 bg-gradient-to-r from-primary to-cyan-500 text-white shadow-lg shadow-primary/25"
                    >
                      <File className="w-4 h-4" />
                      Browse Files
                    </motion.button>
                  </div>
                </div>

                <div className="mt-5 space-y-3">
                  <div className="flex items-center gap-3 p-3 rounded-lg bg-muted/50">
                    <div className="w-8 h-8 rounded-lg bg-green-500/20 flex items-center justify-center">
                      <CheckCircle className="w-4 h-4 text-green-500" />
                    </div>
                    <div className="flex-1">
                      <p className="text-sm font-medium" style={{ color: 'var(--text-primary)' }}>Auto Table Creation</p>
                      <p className="text-xs" style={{ color: 'var(--text-secondary)' }}>Your file is automatically converted to a queryable table</p>
                    </div>
                  </div>

                  <div className="flex items-center gap-3 p-3 rounded-lg bg-muted/50">
                    <div className="w-8 h-8 rounded-lg bg-cyan-500/20 flex items-center justify-center">
                      <Sparkles className="w-4 h-4 text-cyan-500" />
                    </div>
                    <div className="flex-1">
                      <p className="text-sm font-medium" style={{ color: 'var(--text-primary)' }}>AI-Powered Analysis</p>
                      <p className="text-xs" style={{ color: 'var(--text-secondary)' }}>Get instant insights once your data is uploaded</p>
                    </div>
                    <ChevronRight className="w-4 h-4 text-muted-foreground" />
                  </div>
                </div>
              </motion.div>
            )}

            {uploadStatus === 'uploading' && (
              <motion.div
                key="uploading"
                initial={{ opacity: 0, scale: 0.95 }}
                animate={{ opacity: 1, scale: 1 }}
                exit={{ opacity: 0, scale: 0.95 }}
                className="py-8"
              >
                {selectedFile && (
                  <div className="flex items-center gap-4 mb-6 p-4 rounded-xl bg-muted/30">
                    <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-primary/20 to-purple-500/20 flex items-center justify-center">
                      {(() => {
                        const Icon = getFileIcon(selectedFile.name)
                        return <Icon className="w-6 h-6 text-primary" />
                      })()}
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium truncate" style={{ color: 'var(--text-primary)' }}>
                        {selectedFile.name}
                      </p>
                      <p className="text-xs" style={{ color: 'var(--text-secondary)' }}>
                        {formatFileSize(selectedFile.size)}
                      </p>
                    </div>
                    <motion.div
                      animate={{ rotate: 360 }}
                      transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
                      className="w-5 h-5 border-2 border-primary border-t-transparent rounded-full"
                    />
                  </div>
                )}

                <div className="space-y-3">
                  <div className="flex justify-between items-center">
                    <span className="text-sm font-medium" style={{ color: 'var(--text-primary)' }}>
                      Uploading...
                    </span>
                    <span className="text-sm font-bold bg-gradient-to-r from-primary to-cyan-500 bg-clip-text text-transparent">
                      {Math.round(uploadProgress)}%
                    </span>
                  </div>
                  <div className="h-2 rounded-full bg-muted overflow-hidden">
                    <motion.div
                      className="h-full rounded-full bg-gradient-to-r from-primary to-cyan-500"
                      initial={{ width: 0 }}
                      animate={{ width: `${uploadProgress}%` }}
                      transition={{ duration: 0.3 }}
                    />
                  </div>
                  <p className="text-xs text-center" style={{ color: 'var(--text-secondary)' }}>
                    Processing and analyzing your data
                  </p>
                </div>
              </motion.div>
            )}

            {uploadStatus === 'success' && uploadedFile && (
              <motion.div
                key="success"
                initial={{ opacity: 0, scale: 0.95 }}
                animate={{ opacity: 1, scale: 1 }}
                exit={{ opacity: 0, scale: 0.95 }}
                className="py-4"
              >
                <div className="flex flex-col items-center text-center">
                  <motion.div
                    initial={{ scale: 0 }}
                    animate={{ scale: 1 }}
                    transition={{ type: "spring", delay: 0.1 }}
                    className="w-20 h-20 rounded-full bg-gradient-to-br from-green-500/20 to-emerald-500/20 flex items-center justify-center mb-4"
                  >
                    <CheckCircle className="w-10 h-10 text-green-500" />
                  </motion.div>

                  <h3 className="text-xl font-semibold mb-1" style={{ color: 'var(--text-primary)' }}>
                    Upload Complete!
                  </h3>
                  <p className="text-sm mb-4" style={{ color: 'var(--text-secondary)' }}>
                    {selectedFile?.name}
                  </p>

                  <div className="w-full p-4 rounded-xl bg-gradient-to-br from-primary/10 to-purple-500/10 border border-primary/20 mb-4">
                    <div className="flex items-center justify-center gap-2 mb-2">
                      <Table className="w-4 h-4 text-primary" />
                      <span className="text-sm font-medium" style={{ color: 'var(--text-primary)' }}>Table Created</span>
                    </div>
                    <code className="text-lg font-mono bg-black/20 px-3 py-1.5 rounded-lg block" style={{ color: 'var(--primary)' }}>
                      {uploadedFile.table_name}
                    </code>
                  </div>

                  {uploadedFile.data_quality_score !== undefined && (
                    <Card className="w-full p-4 border-0 bg-muted/30">
                      <h4 className="text-sm font-semibold mb-3 flex items-center gap-2" style={{ color: 'var(--text-primary)' }}>
                        <Sparkles className="w-4 h-4 text-cyan-500" />
                        Data Quality Analysis
                      </h4>
                      <div className="grid grid-cols-2 gap-3">
                        <div className="flex items-center justify-between p-2 rounded-lg bg-background/50">
                          <span className="text-xs" style={{ color: 'var(--text-secondary)' }}>Quality Score</span>
                          <span className="text-sm font-bold" style={{
                            color: uploadedFile.data_quality_score >= 0.8 ? '#10b981' : uploadedFile.data_quality_score >= 0.6 ? '#f59e0b' : '#ef4444'
                          }}>
                            {(uploadedFile.data_quality_score * 100).toFixed(0)}%
                          </span>
                        </div>
                        {uploadedFile.missing_data_percent !== undefined && (
                          <div className="flex items-center justify-between p-2 rounded-lg bg-background/50">
                            <span className="text-xs" style={{ color: 'var(--text-secondary)' }}>Missing Data</span>
                            <span className="text-sm font-semibold" style={{ color: 'var(--text-primary)' }}>
                              {uploadedFile.missing_data_percent.toFixed(1)}%
                            </span>
                          </div>
                        )}
                        {uploadedFile.row_count !== undefined && (
                          <div className="flex items-center justify-between p-2 rounded-lg bg-background/50">
                            <span className="text-xs" style={{ color: 'var(--text-secondary)' }}>Total Rows</span>
                            <span className="text-sm font-semibold" style={{ color: 'var(--text-primary)' }}>
                              {uploadedFile.row_count.toLocaleString()}
                            </span>
                          </div>
                        )}
                        {uploadedFile.schema_json?.columns !== undefined && (
                          <div className="flex items-center justify-between p-2 rounded-lg bg-background/50">
                            <span className="text-xs" style={{ color: 'var(--text-secondary)' }}>Columns</span>
                            <span className="text-sm font-semibold" style={{ color: 'var(--text-primary)' }}>
                              {uploadedFile.schema_json.columns.length}
                            </span>
                          </div>
                        )}
                      </div>
                    </Card>
                  )}
                </div>
              </motion.div>
            )}

            {uploadStatus === 'error' && (
              <motion.div
                key="error"
                initial={{ opacity: 0, scale: 0.95 }}
                animate={{ opacity: 1, scale: 1 }}
                exit={{ opacity: 0, scale: 0.95 }}
                className="py-4"
              >
                <div className="flex flex-col items-center text-center">
                  <motion.div
                    initial={{ scale: 0 }}
                    animate={{ scale: 1 }}
                    transition={{ type: "spring" }}
                    className="w-16 h-16 rounded-full bg-red-500/20 flex items-center justify-center mb-4"
                  >
                    <AlertCircle className="w-8 h-8 text-red-500" />
                  </motion.div>

                  <h3 className="text-lg font-semibold mb-2" style={{ color: 'var(--text-primary)' }}>
                    Upload Failed
                  </h3>
                  <p className="text-sm mb-4 px-4 py-2 rounded-lg bg-red-500/10 text-red-400 border border-red-500/20">
                    {errorMessage}
                  </p>

                  <div className="flex gap-3">
                    <Button
                      variant="outline"
                      onClick={() => setUploadStatus('idle')}
                      className="flex-1"
                    >
                      Try Again
                    </Button>
                    <Button
                      variant="ghost"
                      onClick={handleClose}
                      className="flex-1"
                    >
                      Cancel
                    </Button>
                  </div>
                </div>
              </motion.div>
            )}
          </AnimatePresence>
        </div>
      </motion.div>
    </div>
  )
}
