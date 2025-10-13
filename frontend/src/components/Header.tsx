'use client'

import { Brain, BarChart3 } from 'lucide-react'

export default function Header() {
  return (
    <header className="text-center">
      <div className="flex items-center justify-center gap-3 mb-1">
        <div className="p-1.5 rounded-lg" style={{ background: 'var(--accent-color)' }}>
          <Brain className="w-4 h-4" style={{ color: 'var(--primary-background)' }} />
        </div>
        <h1 className="text-xl font-bold" style={{ color: 'var(--text-primary)' }}>
          InsightIQ Analytics
        </h1>
        <div className="p-1.5 rounded-lg" style={{ background: 'var(--accent-color)' }}>
          <BarChart3 className="w-4 h-4" style={{ color: 'var(--primary-background)' }} />
        </div>
      </div>
      <p className="text-xs max-w-2xl mx-auto" style={{ color: 'var(--text-secondary)' }}>
        AI-powered insights from your analytics data instantly.
      </p>
    </header>
  )
}