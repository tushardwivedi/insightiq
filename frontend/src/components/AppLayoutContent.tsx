'use client'

import { useState, useEffect, useCallback } from 'react'
import { useRouter } from 'next/navigation'
import Header from '@/components/Header'
import CompactSystemStatus from '@/components/CompactSystemStatus'
import ConnectorSidebar from '@/components/ConnectorSidebar'
import { apiClient } from '@/lib/api'
import { HealthCheck } from '@/types'
import { useAuth } from '@/contexts/AuthContext'

interface AppLayoutContentProps {
  children: React.ReactNode
}

export default function AppLayoutContent({ children }: AppLayoutContentProps) {
  const router = useRouter()
  const { isAuthenticated, isLoading } = useAuth()
  const [health, setHealth] = useState<HealthCheck | null>(null)
  const [sidebarOpen, setSidebarOpen] = useState(false)

  const handleCloseSidebar = useCallback(() => {
    setSidebarOpen(false)
  }, [])

  const handleOpenSidebar = useCallback(() => {
    setSidebarOpen(true)
  }, [])

  const checkHealth = useCallback(async () => {
    try {
      const healthData = await apiClient.healthCheck()
      setHealth(healthData)
    } catch (error) {
      console.error('Health check failed:', error)
    }
  }, [])

  useEffect(() => {
    // Check authentication
    if (!isLoading && !isAuthenticated) {
      router.push('/login')
      return
    }

    if (isAuthenticated) {
      checkHealth()
      const interval = setInterval(checkHealth, 30000)
      return () => clearInterval(interval)
    }
  }, [isAuthenticated, isLoading, router, checkHealth])

  // Show loading state while checking authentication
  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center" style={{ background: 'var(--primary-background)' }}>
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 mx-auto mb-4" style={{ borderColor: 'var(--accent-color)' }}></div>
          <p style={{ color: 'var(--text-secondary)' }}>Loading...</p>
        </div>
      </div>
    )
  }

  // Don't render protected content if not authenticated
  if (!isAuthenticated) {
    return null
  }

  return (
    <div className="flex h-screen" style={{ background: 'var(--primary-background)' }}>
      {/* Connector Sidebar */}
      <ConnectorSidebar isOpen={sidebarOpen} onClose={handleCloseSidebar} />

      {/* Main Content - Adjust for collapsible sidebar */}
      <div className="flex-1 flex flex-col ml-0 lg:ml-16 transition-all duration-300 ease-in-out">
        {/* Sticky Header with sidebar toggle */}
        <div className="sticky top-0 z-40 shadow-sm" style={{ background: 'var(--surface-color)' }}>
          <div className="container mx-auto px-4">
            <div className="flex items-center gap-4">
              <button
                onClick={handleOpenSidebar}
                className="lg:hidden p-2 rounded-md transition-colors"
                style={{ color: 'var(--text-primary)' }}
                onMouseEnter={(e) => e.currentTarget.style.backgroundColor = 'var(--hover-surface)'}
                onMouseLeave={(e) => e.currentTarget.style.backgroundColor = 'transparent'}
                title="Open Data Sources"
              >
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
                </svg>
              </button>
              <div className="hidden lg:block">
                <div className="flex items-center gap-2 px-2 py-1 text-sm" style={{ color: 'var(--text-secondary)' }}>
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 7v10c0 2.21 1.79 4 4 4h8c2.21 0 4-1.79 4-4V7c0-2.21-1.79-4-4-4H8c-2.21 0-4 1.79-4 4z" />
                  </svg>
                  <span className="text-xs">Hover left sidebar to expand</span>
                </div>
              </div>
              <div className="flex-1 flex items-center justify-between">
                <Header />
                <div className="ml-3">
                  <CompactSystemStatus health={health} />
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Main Content */}
        {children}
      </div>
    </div>
  )
}
