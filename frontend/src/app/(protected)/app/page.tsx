'use client'

import { useState } from 'react'
import CommandBar from '@/components/CommandBar'
import SQLQuerySection from '@/components/SQLQuerySection'
import ResultsSection from '@/components/ResultsSection'

export default function AnalyticsPage() {
  const [results, setResults] = useState<any>(null)
  const [loading, setLoading] = useState(false)

  return (
    <div className="flex-1 overflow-y-auto pt-4 relative">
      <div className="container mx-auto px-4 py-8 pb-24">
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

      {/* Command Bar - Fixed at bottom of content area */}
      <div className="sticky bottom-0 left-0 right-0" style={{ background: 'var(--primary-background)' }}>
        <CommandBar
          onResult={setResults}
          onLoading={setLoading}
        />
      </div>
    </div>
  )
}
