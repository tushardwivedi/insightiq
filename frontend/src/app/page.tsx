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

export default function Dashboard() {
  const [health, setHealth] = useState<HealthCheck | null>(null)
  const [results, setResults] = useState<any>(null)
  const [loading, setLoading] = useState(false)

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
    <div className="container mx-auto px-4 py-8">
      <Header />
      <StatusIndicator health={health} />
      
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 mt-8">
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
  )
}