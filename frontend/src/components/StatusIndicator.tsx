'use client'

import { CheckCircle, AlertCircle, XCircle } from 'lucide-react'
import { HealthCheck } from '@/types'

interface Props {
  health: HealthCheck | null
}

export default function StatusIndicator({ health }: Props) {
  if (!health) {
    return (
      <div className="flex items-center justify-center">
        <div className="animate-spin rounded-full h-6 w-6 border-b-2" style={{ borderColor: 'var(--accent-color)' }}></div>
        <span className="ml-2" style={{ color: 'var(--text-secondary)' }}>Checking system status...</span>
      </div>
    )
  }

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'healthy':
        return <CheckCircle className="w-5 h-5" style={{ color: 'var(--success-color)' }} />
      case 'degraded':
        return <AlertCircle className="w-5 h-5" style={{ color: 'var(--warning-color)' }} />
      default:
        return <XCircle className="w-5 h-5" style={{ color: 'var(--error-color)' }} />
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'healthy':
        return { background: 'rgba(104, 211, 145, 0.1)', borderColor: 'var(--success-color)' }
      case 'degraded':
        return { background: 'rgba(246, 224, 94, 0.1)', borderColor: 'var(--warning-color)' }
      default:
        return { background: 'rgba(252, 129, 129, 0.1)', borderColor: 'var(--error-color)' }
    }
  }

  return (
    <div className="rounded-lg p-4 border" style={getStatusColor(health.status)}>
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          {getStatusIcon(health.status)}
          <span className="font-medium capitalize" style={{ color: 'var(--text-primary)' }}>System {health.status}</span>
        </div>
        <div className="text-sm" style={{ color: 'var(--text-secondary)' }}>
          Last checked: {health.timestamp ? new Date(health.timestamp).toLocaleTimeString() : 'N/A'}
        </div>
      </div>
      
      <div className="mt-3 grid grid-cols-1 md:grid-cols-3 gap-2">
        {Object.entries(health.services).map(([service, status]) => (
          <div key={service} className="flex items-center gap-2 text-sm">
            {status.status === 'healthy' ? (
              <CheckCircle className="w-4 h-4" style={{ color: 'var(--success-color)' }} />
            ) : (
              <XCircle className="w-4 h-4" style={{ color: 'var(--error-color)' }} />
            )}
            <span className="capitalize font-medium" style={{ color: 'var(--text-primary)' }}>{service}</span>
            <span style={{ color: 'var(--text-secondary)' }}>
              {status.status === 'healthy' ? '✓' : '✗'}
            </span>
          </div>
        ))}
      </div>
    </div>
  )
}