'use client'

import { Brain } from 'lucide-react'

export default function Header() {
  return (
    <div className="flex items-center gap-2">
      <div className="w-8 h-8 rounded-lg flex items-center justify-center" style={{ background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' }}>
        <Brain className="w-5 h-5 text-white" />
      </div>
      <div>
        <h1 className="text-base font-semibold" style={{ color: 'var(--text-primary)' }}>
          InsightIQ
        </h1>
        <p className="text-xs" style={{ color: 'var(--text-secondary)' }}>Analytics Platform</p>
      </div>
    </div>
  )
}
