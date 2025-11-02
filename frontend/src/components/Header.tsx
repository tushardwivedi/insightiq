'use client'

import { useState, useRef, useEffect } from 'react'
import { Brain, User, Settings, Key, HelpCircle, LogOut, ChevronDown, Bell } from 'lucide-react'
import { useAuth } from '@/contexts/AuthContext'

export default function Header() {
  const { user, logout } = useAuth()
  const [showUserMenu, setShowUserMenu] = useState(false)
  const menuRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        setShowUserMenu(false)
      }
    }

    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  const getInitials = (name: string | undefined) => {
    if (!name) return 'U'
    return name
      .split(' ')
      .map(n => n[0])
      .join('')
      .toUpperCase()
      .slice(0, 2)
  }

  return (
    <header className="border-b" style={{ borderColor: 'var(--border-color)' }}>
      <div className="flex items-center justify-between h-14">
        {/* Logo & Brand */}
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

        {/* Right Section: Status + Notifications + User */}
        <div className="flex items-center gap-3">
          {/* Notifications */}
          <button
            className="relative p-2 rounded-lg transition-colors hover:bg-opacity-10"
            style={{ color: 'var(--text-secondary)' }}
            onMouseEnter={(e) => e.currentTarget.style.background = 'var(--hover-surface)'}
            onMouseLeave={(e) => e.currentTarget.style.background = 'transparent'}
            title="Notifications"
          >
            <Bell className="w-5 h-5" />
            {/* Notification badge (optional) */}
            {/* <span className="absolute top-1.5 right-1.5 w-2 h-2 bg-red-500 rounded-full"></span> */}
          </button>

          {/* User Menu */}
          <div className="relative" ref={menuRef}>
            <button
              onClick={() => setShowUserMenu(!showUserMenu)}
              className="flex items-center gap-2 px-2 py-1.5 rounded-lg transition-colors"
              style={{
                background: showUserMenu ? 'var(--hover-surface)' : 'transparent',
                border: `1px solid ${showUserMenu ? 'var(--accent-color)' : 'var(--border-color)'}`
              }}
              onMouseEnter={(e) => {
                if (!showUserMenu) e.currentTarget.style.background = 'var(--hover-surface)'
              }}
              onMouseLeave={(e) => {
                if (!showUserMenu) e.currentTarget.style.background = 'transparent'
              }}
            >
              <div
                className="w-7 h-7 rounded-full flex items-center justify-center text-xs font-semibold text-white"
                style={{ background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' }}
              >
                {getInitials(user?.name || user?.email)}
              </div>
              <div className="hidden md:flex flex-col items-start">
                <span className="text-xs font-medium leading-none mb-0.5" style={{ color: 'var(--text-primary)' }}>
                  {user?.name || user?.email?.split('@')[0]}
                </span>
                <span className="text-xs leading-none" style={{ color: 'var(--text-secondary)' }}>
                  {user?.role || 'User'}
                </span>
              </div>
              <ChevronDown
                className={`w-4 h-4 transition-transform ${showUserMenu ? 'rotate-180' : ''}`}
                style={{ color: 'var(--text-secondary)' }}
              />
            </button>

            {/* Dropdown Menu */}
            {showUserMenu && (
              <div
                className="absolute right-0 mt-2 w-56 rounded-lg shadow-lg border py-1 z-50"
                style={{
                  background: 'var(--surface-color)',
                  borderColor: 'var(--border-color)'
                }}
              >
                {/* User Info */}
                <div className="px-3 py-2 border-b" style={{ borderColor: 'var(--border-color)' }}>
                  <p className="text-sm font-medium" style={{ color: 'var(--text-primary)' }}>
                    {user?.name || 'User'}
                  </p>
                  <p className="text-xs" style={{ color: 'var(--text-secondary)' }}>
                    {user?.email}
                  </p>
                  {user?.role === 'admin' && (
                    <span className="inline-block mt-1 text-xs px-2 py-0.5 rounded" style={{ background: 'var(--accent-color)', color: 'white' }}>
                      Admin
                    </span>
                  )}
                </div>

                {/* Menu Items */}
                <div className="py-1">
                  <button
                    className="w-full px-3 py-2 text-left flex items-center gap-2 text-sm transition-colors"
                    style={{ color: 'var(--text-primary)' }}
                    onMouseEnter={(e) => e.currentTarget.style.background = 'var(--hover-surface)'}
                    onMouseLeave={(e) => e.currentTarget.style.background = 'transparent'}
                  >
                    <User className="w-4 h-4" />
                    Profile Settings
                  </button>
                  <button
                    className="w-full px-3 py-2 text-left flex items-center gap-2 text-sm transition-colors"
                    style={{ color: 'var(--text-primary)' }}
                    onMouseEnter={(e) => e.currentTarget.style.background = 'var(--hover-surface)'}
                    onMouseLeave={(e) => e.currentTarget.style.background = 'transparent'}
                  >
                    <Settings className="w-4 h-4" />
                    Account Preferences
                  </button>
                  <button
                    className="w-full px-3 py-2 text-left flex items-center gap-2 text-sm transition-colors"
                    style={{ color: 'var(--text-primary)' }}
                    onMouseEnter={(e) => e.currentTarget.style.background = 'var(--hover-surface)'}
                    onMouseLeave={(e) => e.currentTarget.style.background = 'transparent'}
                  >
                    <Key className="w-4 h-4" />
                    API Keys
                  </button>
                  <button
                    className="w-full px-3 py-2 text-left flex items-center gap-2 text-sm transition-colors"
                    style={{ color: 'var(--text-primary)' }}
                    onMouseEnter={(e) => e.currentTarget.style.background = 'var(--hover-surface)'}
                    onMouseLeave={(e) => e.currentTarget.style.background = 'transparent'}
                  >
                    <HelpCircle className="w-4 h-4" />
                    Help & Documentation
                  </button>
                </div>

                {/* Logout */}
                <div className="border-t" style={{ borderColor: 'var(--border-color)' }}>
                  <button
                    onClick={logout}
                    className="w-full px-3 py-2 text-left flex items-center gap-2 text-sm transition-colors"
                    style={{ color: '#ef4444' }}
                    onMouseEnter={(e) => {
                      e.currentTarget.style.background = 'rgba(239, 68, 68, 0.1)'
                    }}
                    onMouseLeave={(e) => e.currentTarget.style.background = 'transparent'}
                  >
                    <LogOut className="w-4 h-4" />
                    Sign Out
                  </button>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </header>
  )
}