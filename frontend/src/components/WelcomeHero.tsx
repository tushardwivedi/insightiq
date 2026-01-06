'use client'

import { Sparkles, Clock, TrendingUp, Activity } from 'lucide-react'
import { useAuth } from '@/contexts/AuthContext'
import { Card, CardContent } from '@/components/ui/card'

export default function WelcomeHero() {
  const { user } = useAuth()

  const getGreeting = () => {
    const hour = new Date().getHours()
    if (hour < 12) return 'Good morning'
    if (hour < 18) return 'Good afternoon'
    return 'Good evening'
  }

  const stats = [
    { label: 'Queries Today', value: '12', icon: Activity, color: 'text-blue-500' },
    { label: 'Avg Response', value: '340ms', icon: Clock, color: 'text-green-500' },
    { label: 'Success Rate', value: '98%', icon: TrendingUp, color: 'text-yellow-500' },
  ]

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between flex-wrap gap-4">
        {/* Welcome Message */}
        <div>
          <h1 className="text-3xl font-bold flex items-center gap-2">
            <Sparkles className="w-7 h-7 text-primary" />
            {getGreeting()}, {user?.name?.split(' ')[0] || user?.email?.split('@')[0]}!
          </h1>
          <p className="text-muted-foreground mt-1">
            What would you like to analyze today?
          </p>
        </div>

        {/* Quick Stats */}
        <div className="flex gap-3">
          {stats.map((stat) => (
            <Card key={stat.label} className="px-4 py-2">
              <CardContent className="flex items-center gap-3 p-0">
                <div className="w-8 h-8 rounded-lg bg-muted flex items-center justify-center">
                  <stat.icon className={`w-4 h-4 ${stat.color}`} />
                </div>
                <div>
                  <p className="text-xs text-muted-foreground">
                    {stat.label}
                  </p>
                  <p className="text-sm font-semibold">
                    {stat.value}
                  </p>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    </div>
  )
}
