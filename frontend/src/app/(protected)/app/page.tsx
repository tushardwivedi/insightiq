'use client'

import { useState } from 'react'
import CommandBar from '@/components/CommandBar'
import SQLQuerySection from '@/components/SQLQuerySection'
import ResultsSection from '@/components/ResultsSection'
import WelcomeHero from '@/components/WelcomeHero'

export default function AnalyticsPage() {
  const [results, setResults] = useState<any>(null)
  const [loading, setLoading] = useState(false)

  return (
    <div className="flex-1 overflow-y-auto relative">
      <div className="container mx-auto px-4 py-6 pb-24 max-w-7xl">
        {/* Welcome Hero Section */}
        <WelcomeHero />

        {/* Main Workspace */}
        <div className="grid grid-cols-1 gap-6">
          {/* Query Editor */}
          <div className="col-span-1">
            <SQLQuerySection
              onResult={setResults}
              onLoading={setLoading}
            />
          </div>

          {/* Results Panel */}
          {(results || loading) && (
            <div className="col-span-1">
              <ResultsSection results={results} loading={loading} />
            </div>
          )}
        </div>
      </div>

      {/* Command Bar - Fixed at bottom of content area */}
      <div className="sticky bottom-0 left-0 right-0 border-t" style={{ background: 'var(--surface-color)', borderColor: 'var(--border-color)' }}>
        <div className="container mx-auto px-4 max-w-7xl">
          <CommandBar
            onResult={setResults}
            onLoading={setLoading}
          />
        </div>
      </div>
    </div>
  )
}
