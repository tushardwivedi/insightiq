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
        <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-500"></div>
        <span className="ml-2 text-gray-600">Checking system status...</span>
      </div>
    )
  }

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'healthy':
        return <CheckCircle className="w-5 h-5 text-green-500" />
      case 'degraded':
        return <AlertCircle className="w-5 h-5 text-yellow-500" />
      default:
        return <XCircle className="w-5 h-5 text-red-500" />
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'healthy':
        return 'bg-green-50 border-green-200'
      case 'degraded':
        return 'bg-yellow-50 border-yellow-200'
      default:
        return 'bg-red-50 border-red-200'
    }
  }

  return (
    <div className={`rounded-lg p-4 border ${getStatusColor(health.status)}`}>
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          {getStatusIcon(health.status)}
          <span className="font-medium capitalize">System {health.status}</span>
        </div>
        <div className="text-sm text-gray-500">
          Last checked: {health.timestamp ? new Date(health.timestamp).toLocaleTimeString() : 'N/A'}
        </div>
      </div>
      
      <div className="mt-3 grid grid-cols-1 md:grid-cols-3 gap-2">
        {Object.entries(health.services).map(([service, status]) => (
          <div key={service} className="flex items-center gap-2 text-sm">
            {status.status === 'healthy' ? (
              <CheckCircle className="w-4 h-4 text-green-500" />
            ) : (
              <XCircle className="w-4 h-4 text-red-500" />
            )}
            <span className="capitalize font-medium">{service}</span>
            <span className="text-gray-500">
              {status.status === 'healthy' ? '✓' : '✗'}
            </span>
          </div>
        ))}
      </div>
    </div>
  )
}