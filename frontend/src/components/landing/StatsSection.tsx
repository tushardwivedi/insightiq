'use client'

import { MessageSquare, Zap, Clock, Users } from 'lucide-react'
import { LandingStatCard } from './LandingStatCard'

const stats = [
  { value: '10M+', label: 'Queries Processed', icon: MessageSquare },
  { value: '99.9%', label: 'Uptime Guarantee', icon: Zap },
  { value: '< 100ms', label: 'Query Response Time', icon: Clock },
  { value: '24/7', label: 'Community Support', icon: Users }
]

export function StatsSection() {
  return (
    <section className="relative py-16 border-y border-[#4A5568]/30" style={{ background: '#1A202C' }}>
      <div className="container mx-auto px-4">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-8">
          {stats.map((stat, index) => (
            <LandingStatCard key={stat.label} {...stat} delay={index * 0.1} />
          ))}
        </div>
      </div>
    </section>
  )
}
