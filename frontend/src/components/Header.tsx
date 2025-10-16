'use client'

import { Brain, BarChart3, LogOut, User } from 'lucide-react'
import { useAuth } from '@/contexts/AuthContext'

export default function Header() {
  const { user, logout } = useAuth()

  return (
    <header>
      <div className="flex items-center justify-between mb-1">
        <div className="flex items-center gap-3">
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

        {/* User Menu */}
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2 px-3 py-1.5 rounded-lg" style={{ background: 'var(--surface-color)', border: '1px solid var(--border-color)' }}>
            <User className="w-4 h-4" style={{ color: 'var(--text-secondary)' }} />
            <span className="text-sm font-medium" style={{ color: 'var(--text-primary)' }}>
              {user?.name || user?.email}
            </span>
            {user?.role === 'admin' && (
              <span className="text-xs px-2 py-0.5 rounded" style={{ background: 'var(--accent-color)', color: 'var(--primary-background)' }}>
                Admin
              </span>
            )}
          </div>
          <button
            onClick={logout}
            className="flex items-center gap-2 px-3 py-1.5 rounded-lg transition-colors"
            style={{ color: 'var(--text-secondary)', border: '1px solid var(--border-color)' }}
            onMouseEnter={(e) => {
              e.currentTarget.style.background = 'var(--hover-surface)'
              e.currentTarget.style.color = 'var(--text-primary)'
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.background = 'transparent'
              e.currentTarget.style.color = 'var(--text-secondary)'
            }}
            title="Logout"
          >
            <LogOut className="w-4 h-4" />
            <span className="text-sm font-medium">Logout</span>
          </button>
        </div>
      </div>
      <p className="text-xs max-w-2xl" style={{ color: 'var(--text-secondary)' }}>
        AI-powered insights from your analytics data instantly.
      </p>
    </header>
  )
}