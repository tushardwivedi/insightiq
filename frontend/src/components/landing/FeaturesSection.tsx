'use client'

import { motion } from 'framer-motion'
import { MessageSquare, Brain, Mic, Database, Code, Shield } from 'lucide-react'
import { BentoCard } from './BentoCard'
import { GlowEffect } from './GlowEffect'

const features = [
  { icon: MessageSquare, title: 'Natural Language Queries', description: 'Ask questions in plain English. Our AI translates your questions into optimized SQL queries instantly.', gradient: 'from-[#4FD1C5]/10 to-transparent', className: 'md:col-span-2 md:row-span-2' },
  { icon: Brain, title: 'AI-Powered Insights', description: 'Local LLM delivers intelligent, contextual answers that understand your business.', gradient: 'from-cyan-500/10 to-transparent', className: 'md:col-span-1' },
  { icon: Mic, title: 'Voice Analytics', description: 'Speak your queries naturally with integrated Whisper speech recognition.', gradient: 'from-blue-500/10 to-transparent', className: 'md:col-span-1' },
  { icon: Database, title: 'Universal Connectors', description: 'PostgreSQL, MySQL, and more. One platform for all your data sources.', gradient: 'from-indigo-500/10 to-transparent', className: 'md:col-span-1' },
  { icon: Code, title: 'SQL Transparency', description: 'See and edit generated SQL with professional code editor and syntax highlighting.', gradient: 'from-purple-500/10 to-transparent', className: 'md:col-span-1' },
  { icon: Shield, title: 'Self-Hosted & Secure', description: 'Deploy on your infrastructure. Your data never leaves your servers. Complete privacy control.', gradient: 'from-[#4FD1C5]/10 to-transparent', className: 'md:col-span-2' }
]

export function FeaturesSection() {
  return (
    <section className="relative py-24">
      <div className="absolute inset-0 bg-gradient-to-b from-[#1A202C] via-[#0F172A] to-[#1A202C]" />
      <GlowEffect variant="top" />

      <div className="container mx-auto px-4 relative z-10">
        <motion.div initial={{ opacity: 0, y: 20 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true }} className="text-center max-w-3xl mx-auto mb-16">
          <h2 className="text-4xl md:text-5xl lg:text-6xl font-bold mb-6 text-white">
            Everything You Need for <span className="bg-gradient-to-r from-[#4FD1C5] to-[#38B2AC] bg-clip-text text-transparent">Modern Analytics</span>
          </h2>
          <p className="text-lg text-[#A0AEC0]">A complete analytics platform that respects your privacy and runs entirely on your infrastructure</p>
        </motion.div>

        <div className="grid grid-cols-1 md:grid-cols-3 auto-rows-[180px] gap-4 max-w-6xl mx-auto">
          {features.map((feature, index) => (
            <BentoCard key={feature.title} {...feature} delay={index * 0.1} />
          ))}
        </div>
      </div>
    </section>
  )
}
