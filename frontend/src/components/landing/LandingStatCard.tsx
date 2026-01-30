'use client'

import { motion } from 'framer-motion'

interface LandingStatCardProps {
  value: string
  label: string
  icon: React.ElementType
  delay?: number
}

export function LandingStatCard({ value, label, icon: Icon, delay = 0 }: LandingStatCardProps) {
  return (
    <motion.div
      initial={{ opacity: 0, scale: 0.9 }}
      whileInView={{ opacity: 1, scale: 1 }}
      viewport={{ once: true }}
      transition={{ delay, duration: 0.5 }}
      className="text-center p-6"
    >
      <div className="inline-flex items-center justify-center w-14 h-14 rounded-2xl bg-[#4FD1C5]/10 mb-4">
        <Icon className="w-7 h-7 text-[#4FD1C5]" />
      </div>
      <div className="text-4xl md:text-5xl font-bold bg-gradient-to-r from-[#4FD1C5] to-[#38B2AC] bg-clip-text text-transparent mb-2">
        {value}
      </div>
      <div className="text-[#A0AEC0] text-sm">{label}</div>
    </motion.div>
  )
}
