'use client'

import { motion } from 'framer-motion'

interface TimelineStepProps {
  number: number
  title: string
  description: string
  icon: React.ElementType
  delay?: number
}

export function TimelineStep({ number, title, description, icon: Icon, delay = 0 }: TimelineStepProps) {
  return (
    <motion.div
      initial={{ opacity: 0, x: -20 }}
      whileInView={{ opacity: 1, x: 0 }}
      viewport={{ once: true }}
      transition={{ delay, duration: 0.5 }}
      className="flex gap-6 relative"
    >
      <div className="flex flex-col items-center">
        <div className="flex h-12 w-12 items-center justify-center rounded-full bg-gradient-to-br from-[#4FD1C5] to-[#38B2AC] shadow-lg shadow-[#4FD1C5]/30">
          <Icon className="w-6 h-6 text-[#0F172A]" />
        </div>
        <div className="w-px h-full bg-gradient-to-b from-[#4FD1C5] to-transparent mt-4" />
      </div>
      <div className="flex-1 pb-12">
        <div className="inline-block px-3 py-1 rounded-full bg-[#4FD1C5]/10 border border-[#4FD1C5]/30 text-xs text-[#4FD1C5] mb-3">
          Step {number}
        </div>
        <h3 className="text-xl font-semibold text-white mb-2">{title}</h3>
        <p className="text-[#A0AEC0] leading-relaxed">{description}</p>
      </div>
    </motion.div>
  )
}
