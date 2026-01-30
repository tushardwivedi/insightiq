'use client'

import { motion } from 'framer-motion'
import { Database, MessageSquare, TrendingUp, Rocket } from 'lucide-react'
import { TimelineStep } from './TimelineStep'
import { AnimatedOrb } from './AnimatedOrb'

const timelineSteps = [
  { number: 1, title: 'Install with Docker', description: 'Clone the repo and launch with a single docker compose command. All services configured and ready in minutes.', icon: Rocket },
  { number: 2, title: 'Connect Your Data', description: 'Add your database connections through the intuitive UI. Test connections before saving.', icon: Database },
  { number: 3, title: 'Ask Questions', description: 'Start querying your data using natural language, voice, or SQL. Get instant insights.', icon: MessageSquare },
  { number: 4, title: 'Get Insights', description: 'View results with auto-generated visualizations. Export data and share dashboards with your team.', icon: TrendingUp }
]

export function TimelineSection() {
  return (
    <section className="relative py-24" style={{ background: 'linear-gradient(to bottom, #0F172A, #1A202C)' }}>
      <AnimatedOrb className="top-20 left-10" delay={0} size={350} color="cyan" />

      <div className="container mx-auto px-4 relative z-10">
        <motion.div initial={{ opacity: 0, y: 20 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true }} className="text-center max-w-3xl mx-auto mb-16">
          <h2 className="text-4xl md:text-5xl font-bold mb-4 text-white">
            Get Started in <span className="bg-gradient-to-r from-[#4FD1C5] to-[#38B2AC] bg-clip-text text-transparent">4 Simple Steps</span>
          </h2>
          <p className="text-lg text-[#A0AEC0]">From installation to insights in minutes, not hours</p>
        </motion.div>

        <div className="max-w-3xl mx-auto">
          {timelineSteps.map((step, index) => (
            <TimelineStep key={step.number} {...step} delay={index * 0.15} />
          ))}
        </div>
      </div>
    </section>
  )
}
