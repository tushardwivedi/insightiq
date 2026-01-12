'use client'

import { useState, useRef, useEffect } from 'react'
import { User, Settings, Key, HelpCircle, LogOut, ChevronDown, Bell, CheckCircle, AlertCircle, XCircle, Activity } from 'lucide-react'
import { useAuth } from '@/contexts/AuthContext'
import { HealthCheck } from '@/types'

interface HeaderActionsProps {
  health: HealthCheck | null
}

export default function HeaderActions({ health }: HeaderActionsProps) {
  const { user, logout } = useAuth()
  const [showUserMenu, setShowUserMenu] = useState(false)
  const [showStatusDetails, setShowStatusDetails] = useState(false)
  const userMenuRef = useRef<HTMLDivElement>(null)
  const statusMenuRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (userMenuRef.current && !userMenuRef.current.contains(event.target as Node)) {
        setShowUserMenu(false)
      }
      if (statusMenuRef.current && !statusMenuRef.current.contains(event.target as Node)) {
        setShowStatusDetails(false)
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

  const getStatusColor = () => {
    if (!health) return '#f59e0b' // warning
    switch (health.status) {
      case 'healthy':
        return '#10b981'
      default:
        return '#ef4444'
    }
  }

  const getStatusIcon = () => {
    if (!health) {
      return <AlertCircle className="w-4 h-4" style={{ color: getStatusColor() }} />
    }
    switch (health.status) {
      case 'healthy':
        return <CheckCircle className="w-4 h-4" style={{ color: getStatusColor() }} />
      default:
        return <XCircle className="w-4 h-4" style={{ color: getStatusColor() }} />
    }
  }

  return (
    <div className="flex items-center gap-2">
      {/* Status Details Dropdown */}
      <div className="relative" ref={statusMenuRef}>
        <button
          onClick={() => setShowStatusDetails(!showStatusDetails)}
          className="flex items-center gap-2 px-3 py-1.5 rounded-lg transition-colors"
          style={{
            background: showStatusDetails ? 'var(--hover-surface)' : 'var(--surface-color)',
            border: `1px solid ${showStatusDetails ? 'var(--accent-color)' : 'var(--border-color)'}`
          }}
          onMouseEnter={(e) => {
            if (!showStatusDetails) e.currentTarget.style.background = 'var(--hover-surface)'
          }}
          onMouseLeave={(e) => {
            if (!showStatusDetails) e.currentTarget.style.background = 'var(--surface-color)'
          }}
          title="System Status"
        >
          <div className="w-2 h-2 rounded-full" style={{ background: getStatusColor() }}></div>
          <span className="text-xs font-medium hidden md:inline" style={{ color: 'var(--text-primary)' }}>
            System {health?.status || 'checking'}
          </span>
          {getStatusIcon()}
        </button>

        {/* Status Details Dropdown */}
        {showStatusDetails && (
          <div
            className="absolute right-0 mt-2 w-80 rounded-lg shadow-lg border p-4 z-50"
            style={{
              background: 'var(--surface-color)',
              borderColor: 'var(--border-color)'
            }}
          >
            <div className="flex items-center justify-between mb-3">
              <h3 className="text-sm font-semibold flex items-center gap-2" style={{ color: 'var(--text-primary)' }}>
                <Activity className="w-4 h-4" />
                System Status
              </h3>
              <span className="text-xs px-2 py-0.5 rounded" style={{ background: `${getStatusColor()}20`, color: getStatusColor() }}>
                {(health?.status || 'checking').toUpperCase()}
              </span>
            </div>

            <p className="text-xs mb-3" style={{ color: 'var(--text-secondary)' }}>
              Last checked: {health?.timestamp ? new Date(health.timestamp).toLocaleTimeString() : 'N/A'}
            </p>

            <div className="space-y-2">
              {health?.services && Object.entries(health.services).map(([service, status]) => (
                <div
                  key={service}
                  className="flex items-center justify-between p-2 rounded"
                  style={{ background: 'var(--primary-background)' }}
                >
                  <div className="flex items-center gap-2">
                    {status.status === 'healthy' ? (
                      <CheckCircle className="w-3.5 h-3.5" style={{ color: '#10b981' }} />
                    ) : (
                      <XCircle className="w-3.5 h-3.5" style={{ color: '#ef4444' }} />
                    )}
                    <span className="text-xs font-medium capitalize" style={{ color: 'var(--text-primary)' }}>
                      {service}
                    </span>
                  </div>
                  <span
                    className="text-xs px-2 py-0.5 rounded"
                    style={{
                      background: status.status === 'healthy' ? '#10b98120' : '#ef444420',
                      color: status.status === 'healthy' ? '#10b981' : '#ef4444'
                    }}
                  >
                    {status.status}
                  </span>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>

      {/* Notifications */}
      <button
        className="relative p-2 rounded-lg transition-colors hover:bg-opacity-10 md:flex hidden"
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
      <div className="relative" ref={userMenuRef}>
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
  )
}