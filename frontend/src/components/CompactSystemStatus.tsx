'use client'

import { useState, useRef, useEffect } from 'react'
import { CheckCircle, AlertCircle, XCircle, Activity } from 'lucide-react'
import { HealthCheck } from '@/types'

interface Props {
  health: HealthCheck | null
}

export default function CompactSystemStatus({ health }: Props) {
  const [showDetails, setShowDetails] = useState(false)
  const menuRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        setShowDetails(false)
      }
    }

    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  if (!health) {
    return (
      <button
        className="flex items-center gap-2 px-3 py-1.5 rounded-lg transition-colors"
        style={{ background: 'var(--surface-color)', border: '1px solid var(--border-color)' }}
      >
        <div className="w-2 h-2 rounded-full animate-pulse" style={{ background: 'var(--warning-color)' }}></div>
        <span className="text-xs font-medium" style={{ color: 'var(--text-secondary)' }}>Checking...</span>
      </button>
    )
  }

  const getStatusColor = () => {
    switch (health.status) {
      case 'healthy':
        return '#10b981'
      default:
        return '#ef4444'
    }
  }

  const getStatusIcon = () => {
    switch (health.status) {
      case 'healthy':
        return <CheckCircle className="w-4 h-4" style={{ color: getStatusColor() }} />
      default:
        return <XCircle className="w-4 h-4" style={{ color: getStatusColor() }} />
    }
  }

  return (
    <div className="relative" ref={menuRef}>
      <button
        onClick={() => setShowDetails(!showDetails)}
        className="flex items-center gap-2 px-3 py-1.5 rounded-lg transition-colors"
        style={{
          background: showDetails ? 'var(--hover-surface)' : 'var(--surface-color)',
          border: `1px solid ${showDetails ? 'var(--accent-color)' : 'var(--border-color)'}`
        }}
        onMouseEnter={(e) => {
          if (!showDetails) e.currentTarget.style.background = 'var(--hover-surface)'
        }}
        onMouseLeave={(e) => {
          if (!showDetails) e.currentTarget.style.background = 'var(--surface-color)'
        }}
      >
        <div className="w-2 h-2 rounded-full" style={{ background: getStatusColor() }}></div>
        <span className="text-xs font-medium hidden md:inline" style={{ color: 'var(--text-primary)' }}>
          System {health.status}
        </span>
        {getStatusIcon()}
      </button>

      {/* Status Details Dropdown */}
      {showDetails && (
        <div
          className="absolute right-0 mt-2 w-80 rounded-lg shadow-lg border p-4 z-50"
          style={{
            background: 'var(--surface-color)',
            borderColor: 'var(--border-color)'
          }}
        >
          <div className="flex items-center justify-between mb-3">
            <h3 className="text-sm font-semibold flex items-center gap-2" style={{ color: 'var(--text-primary)' }}>
              <Activity className="w-4 h-4" />
              System Status
            </h3>
            <span className="text-xs px-2 py-0.5 rounded" style={{ background: `${getStatusColor()}20`, color: getStatusColor() }}>
              {health.status.toUpperCase()}
            </span>
          </div>

          <p className="text-xs mb-3" style={{ color: 'var(--text-secondary)' }}>
            Last checked: {health.timestamp ? new Date(health.timestamp).toLocaleTimeString() : 'N/A'}
          </p>

          <div className="space-y-2">
            {Object.entries(health.services).map(([service, status]) => (
              <div
                key={service}
                className="flex items-center justify-between p-2 rounded"
                style={{ background: 'var(--primary-background)' }}
              >
                <div className="flex items-center gap-2">
                  {status.status === 'healthy' ? (
                    <CheckCircle className="w-3.5 h-3.5" style={{ color: '#10b981' }} />
                  ) : (
                    <XCircle className="w-3.5 h-3.5" style={{ color: '#ef4444' }} />
                  )}
                  <span className="text-xs font-medium capitalize" style={{ color: 'var(--text-primary)' }}>
                    {service}
                  </span>
                </div>
                <span
                  className="text-xs px-2 py-0.5 rounded"
                  style={{
                    background: status.status === 'healthy' ? '#10b98120' : '#ef444420',
                    color: status.status === 'healthy' ? '#10b981' : '#ef4444'
                  }}
                >
                  {status.status}
                </span>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}
