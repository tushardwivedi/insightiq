'use client'

import { useState } from 'react'
import { DataConnector } from '@/types'
import { apiClient } from '@/lib/api'

interface ConnectorCardProps {
  connector: DataConnector
  onTest: () => void
  onEdit: (connector: DataConnector) => void
  onDelete: () => void
}

export default function ConnectorCard({ connector, onTest, onEdit, onDelete }: ConnectorCardProps) {
  const [showMenu, setShowMenu] = useState(false)
  const [deleting, setDeleting] = useState(false)

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'connected':
        return { background: '#10b98120', color: '#10b981', border: '#10b98140' }
      case 'disconnected':
        return { background: 'var(--hover-surface)', color: 'var(--text-secondary)', border: 'var(--border-color)' }
      case 'error':
        return { background: '#ef444420', color: '#ef4444', border: '#ef444440' }
      case 'testing':
        return { background: '#f59e0b20', color: '#f59e0b', border: '#f59e0b40' }
      default:
        return { background: 'var(--hover-surface)', color: 'var(--text-secondary)', border: 'var(--border-color)' }
    }
  }

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'connected':
        return '✓'
      case 'error':
        return '✗'
      case 'testing':
        return '⏳'
      default:
        return '◯'
    }
  }

  const getTypeIcon = (type: string) => {
    switch (type) {
      case 'superset':
        return '📊'
      case 'postgres':
        return '🐘'
      case 'file_upload':
        return '📁'
      default:
        return '🔗'
    }
  }

  const handleDelete = async () => {
    if (!confirm(`Are you sure you want to delete "${connector.name}"?`)) {
      return
    }

    try {
      setDeleting(true)
      await apiClient.deleteConnector(connector.id)
      onDelete()
    } catch (error) {
      console.error('Failed to delete connector:', error)
      alert('Failed to delete connector. Please try again.')
    } finally {
      setDeleting(false)
    }
  }

  const formatLastTested = (lastTested?: string) => {
    if (!lastTested) return 'Never tested'
    const date = new Date(lastTested)
    const now = new Date()
    const diffMs = now.getTime() - date.getTime()
    const diffMins = Math.floor(diffMs / (1000 * 60))

    if (diffMins < 1) return 'Just now'
    if (diffMins < 60) return `${diffMins}m ago`
    if (diffMins < 1440) return `${Math.floor(diffMins / 60)}h ago`
    return `${Math.floor(diffMins / 1440)}d ago`
  }

  return (
    <div
      className="relative border rounded-lg p-4 transition-all hover:shadow-md"
      style={{
        background: 'var(--surface-color)',
        borderColor: 'var(--border-color)'
      }}
      onMouseEnter={(e) => e.currentTarget.style.borderColor = 'var(--accent-color)'}
      onMouseLeave={(e) => e.currentTarget.style.borderColor = 'var(--border-color)'}
    >
      {/* Header */}
      <div className="flex items-start justify-between mb-3">
        <div className="flex items-center gap-3">
          <span className="text-2xl">{getTypeIcon(connector.type)}</span>
          <div className="flex-1">
            <h3 className="font-medium" style={{ color: 'var(--text-primary)' }}>{connector.name}</h3>
            <p className="text-sm capitalize" style={{ color: 'var(--text-secondary)' }}>{connector.type}</p>
          </div>
        </div>

        <div className="relative">
          <button
            onClick={() => setShowMenu(!showMenu)}
            className="p-1 rounded transition-colors"
            style={{ color: 'var(--text-secondary)' }}
            onMouseEnter={(e) => e.currentTarget.style.background = 'var(--hover-surface)'}
            onMouseLeave={(e) => e.currentTarget.style.background = 'transparent'}
          >
            <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
              <path d="M10 6a2 2 0 110-4 2 2 0 010 4zM10 12a2 2 0 110-4 2 2 0 010 4zM10 18a2 2 0 110-4 2 2 0 010 4z" />
            </svg>
          </button>

          {showMenu && (
            <div
              className="absolute right-0 top-8 w-32 border rounded-md shadow-lg z-10"
              style={{
                background: 'var(--surface-color)',
                borderColor: 'var(--border-color)'
              }}
            >
              <button
                onClick={() => {
                  onTest()
                  setShowMenu(false)
                }}
                className="w-full px-3 py-2 text-left text-sm transition-colors"
                style={{ color: 'var(--text-primary)' }}
                onMouseEnter={(e) => e.currentTarget.style.background = 'var(--hover-surface)'}
                onMouseLeave={(e) => e.currentTarget.style.background = 'transparent'}
                disabled={connector.status === 'testing'}
              >
                Test Connection
              </button>
              <button
                onClick={() => {
                  // TODO: Implement edit functionality
                  setShowMenu(false)
                }}
                className="w-full px-3 py-2 text-left text-sm transition-colors"
                style={{ color: 'var(--text-primary)' }}
                onMouseEnter={(e) => e.currentTarget.style.background = 'var(--hover-surface)'}
                onMouseLeave={(e) => e.currentTarget.style.background = 'transparent'}
              >
                Edit
              </button>
              <button
                onClick={() => {
                  handleDelete()
                  setShowMenu(false)
                }}
                className="w-full px-3 py-2 text-left text-sm transition-colors"
                style={{ color: '#ef4444' }}
                onMouseEnter={(e) => e.currentTarget.style.background = '#ef444420'}
                onMouseLeave={(e) => e.currentTarget.style.background = 'transparent'}
                disabled={deleting}
              >
                {deleting ? 'Deleting...' : 'Delete'}
              </button>
            </div>
          )}
        </div>
      </div>

      {/* Status */}
      <div className="flex items-center justify-between mb-2">
        <span
          className="inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs font-medium border"
          style={getStatusColor(connector.status)}
        >
          <span>{getStatusIcon(connector.status)}</span>
          {connector.status.charAt(0).toUpperCase() + connector.status.slice(1)}
        </span>
      </div>

      {/* Connection Details */}
      <div className="text-xs space-y-1" style={{ color: 'var(--text-secondary)' }}>
        <div className="flex justify-between">
          <span>URL:</span>
          <span className="truncate ml-2" title={connector.config.url}>
            {connector.config.url}
          </span>
        </div>
        <div className="flex justify-between">
          <span>Last tested:</span>
          <span>{formatLastTested(connector.last_tested)}</span>
        </div>
      </div>

      {/* Test Button */}
      <button
        onClick={onTest}
        disabled={connector.status === 'testing'}
        className="w-full mt-3 px-3 py-1.5 text-sm bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
      >
        {connector.status === 'testing' ? 'Testing...' : 'Test Connection'}
      </button>

      {/* Backdrop to close menu */}
      {showMenu && (
        <div
          className="fixed inset-0 z-0"
          onClick={() => setShowMenu(false)}
        />
      )}
    </div>
  )
}