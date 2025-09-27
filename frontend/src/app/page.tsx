'use client'

import { useState, useEffect } from 'react'
import { HealthCheck } from '@/types'
import { apiClient } from '@/lib/api'
import Header from '@/components/Header'
import StatusIndicator from '@/components/StatusIndicator'
import TextQuerySection from '@/components/TextQuerySection'
import VoiceQuerySection from '@/components/VoiceQuerySection'
import SQLQuerySection from '@/components/SQLQuerySection'
import ResultsSection from '@/components/ResultsSection'
import ConnectorSidebar from '@/components/ConnectorSidebar'

export default function Dashboard() {
  const [health, setHealth] = useState<HealthCheck | null>(null)
  const [results, setResults] = useState<any>(null)
  const [loading, setLoading] = useState(false)
  const [sidebarOpen, setSidebarOpen] = useState(false)

  useEffect(() => {
    checkHealth()
    const interval = setInterval(checkHealth, 30000)
    return () => clearInterval(interval)
  }, [])

  const checkHealth = async () => {
    try {
      const healthData = await apiClient.healthCheck()
      setHealth(healthData)
    } catch (error) {
      console.error('Health check failed:', error)
    }
  }

  return (
    <div className="flex h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
      {/* Connector Sidebar */}
      <ConnectorSidebar isOpen={sidebarOpen} onClose={() => setSidebarOpen(false)} />

      {/* Main Content */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Header with sidebar toggle */}
        <div className="bg-white shadow-sm border-b border-gray-200">
          <div className="container mx-auto px-4 py-4">
            <div className="flex items-center gap-4">
              <button
                onClick={() => setSidebarOpen(true)}
                className="lg:hidden p-2 hover:bg-gray-100 rounded-md transition-colors"
              >
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
                </svg>
              </button>
              <div className="hidden lg:block">
                <button
                  onClick={() => setSidebarOpen(true)}
                  className="flex items-center gap-2 px-3 py-2 text-sm bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
                >
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 7v10c0 2.21 1.79 4 4 4h8c2.21 0 4-1.79 4-4V7c0-2.21-1.79-4-4-4H8c-2.21 0-4 1.79-4 4z" />
                  </svg>
                  Data Sources
                </button>
              </div>
              <div className="flex-1">
                <Header />
              </div>
            </div>
            <StatusIndicator health={health} />
          </div>
        </div>

        {/* Scrollable Content Area */}
        <div className="flex-1 overflow-y-auto">
          <div className="container mx-auto px-4 py-8">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
              <div>
                <TextQuerySection
                  onResult={setResults}
                  onLoading={setLoading}
                />
              </div>

              <div>
                <VoiceQuerySection
                  onResult={setResults}
                  onLoading={setLoading}
                />
              </div>
            </div>

            <div className="mt-8">
              <SQLQuerySection
                onResult={setResults}
                onLoading={setLoading}
              />
            </div>

            <div className="mt-8">
              <ResultsSection results={results} loading={loading} />
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}