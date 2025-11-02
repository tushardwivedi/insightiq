'use client'

import { Sparkles, Clock, TrendingUp, Activity } from 'lucide-react'
import { useAuth } from '@/contexts/AuthContext'

export default function WelcomeHero() {
  const { user } = useAuth()

  const getGreeting = () => {
    const hour = new Date().getHours()
    if (hour < 12) return 'Good morning'
    if (hour < 18) return 'Good afternoon'
    return 'Good evening'
  }

  const stats = [
    { label: 'Queries Today', value: '12', icon: Activity, color: '#667eea' },
    { label: 'Avg Response', value: '340ms', icon: Clock, color: '#10b981' },
    { label: 'Success Rate', value: '98%', icon: TrendingUp, color: '#f59e0b' },
  ]

  return (
    <div className="mb-6">
      <div className="flex items-center justify-between flex-wrap gap-4">
        {/* Welcome Message */}
        <div>
          <h2 className="text-2xl font-bold flex items-center gap-2" style={{ color: 'var(--text-primary)' }}>
            <Sparkles className="w-6 h-6" style={{ color: 'var(--accent-color)' }} />
            {getGreeting()}, {user?.name?.split(' ')[0] || user?.email?.split('@')[0]}!
          </h2>
          <p className="mt-1 text-sm" style={{ color: 'var(--text-secondary)' }}>
            What would you like to analyze today?
          </p>
        </div>

        {/* Quick Stats */}
        <div className="flex gap-4">
          {stats.map((stat) => (
            <div
              key={stat.label}
              className="flex items-center gap-2 px-4 py-2 rounded-lg"
              style={{ background: 'var(--surface-color)', border: '1px solid var(--border-color)' }}
            >
              <div
                className="w-8 h-8 rounded-lg flex items-center justify-center"
                style={{ background: `${stat.color}15` }}
              >
                <stat.icon className="w-4 h-4" style={{ color: stat.color }} />
              </div>
              <div>
                <p className="text-xs" style={{ color: 'var(--text-secondary)' }}>
                  {stat.label}
                </p>
                <p className="text-sm font-semibold" style={{ color: 'var(--text-primary)' }}>
                  {stat.value}
                </p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
