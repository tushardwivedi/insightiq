'use client'

import { useState, useEffect, useCallback } from 'react'
import { useRouter } from 'next/navigation'
import Header from '@/components/Header'
import HeaderActions from '@/components/HeaderActions'
import { apiClient } from '@/lib/api'
import { HealthCheck } from '@/types'
import { useAuth } from '@/contexts/AuthContext'
import { SidebarProvider, SidebarInset, SidebarTrigger } from '@/components/ui/sidebar'
import { AppSidebar } from '@/components/AppSidebar'

interface AppLayoutContentProps {
  children: React.ReactNode
}

export default function AppLayoutContent({ children }: AppLayoutContentProps) {
  const router = useRouter()
  const { isAuthenticated, isLoading } = useAuth()
  const [health, setHealth] = useState<HealthCheck | null>(null)

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
      <div className="min-h-screen flex items-center justify-center bg-background">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto mb-4"></div>
          <p className="text-muted-foreground">Loading...</p>
        </div>
      </div>
    )
  }

  // Don't render protected content if not authenticated
  if (!isAuthenticated) {
    return null
  }

  return (
    <SidebarProvider className="h-svh !min-h-0">
      <AppSidebar />
      <SidebarInset className="h-full overflow-hidden">
        {/* Fixed Header */}
        <header className="flex h-14 shrink-0 items-center gap-4 border-b px-4 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
          <SidebarTrigger className="-ml-1" />
          <div className="flex-1 flex items-center justify-between">
            {/* Logo & Brand on the left */}
            <Header />
            {/* All other actions on the right */}
            <HeaderActions health={health} />
          </div>
        </header>

        {/* Main Content */}
        <main className="flex-1 overflow-hidden flex flex-col">
          {children}
        </main>
      </SidebarInset>
    </SidebarProvider>
  )
}
