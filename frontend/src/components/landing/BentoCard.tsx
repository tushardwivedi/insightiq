'use client'

import { motion } from 'framer-motion'

interface BentoCardProps {
  icon: React.ElementType
  title: string
  description: string
  className?: string
  gradient: string
  delay?: number
}

export function BentoCard({ icon: Icon, title, description, className, gradient, delay = 0 }: BentoCardProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true }}
      transition={{ delay, duration: 0.5 }}
      className={`group relative overflow-hidden rounded-2xl border border-[#4FD1C5]/20 bg-[#1A202C]/80 backdrop-blur-sm transition-all duration-300 hover:border-[#4FD1C5]/40 hover:shadow-lg hover:shadow-[#4FD1C5]/10 ${className}`}
    >
      <div className={`absolute inset-0 bg-gradient-to-br ${gradient} opacity-0 transition-opacity duration-300 group-hover:opacity-100`} />
      <div className="relative z-10 flex h-full flex-col p-6">
        <div className="mb-4 inline-flex h-12 w-12 items-center justify-center rounded-xl bg-[#4FD1C5]/10 text-[#4FD1C5] transition-transform duration-300 group-hover:scale-110">
          <Icon className="h-6 w-6" />
        </div>
        <h3 className="mb-2 text-xl font-semibold text-white">{title}</h3>
        <p className="text-[#A0AEC0] text-sm leading-relaxed">{description}</p>
      </div>
      <div className="absolute bottom-0 left-0 h-1 w-0 bg-gradient-to-r from-[#4FD1C5] to-[#38B2AC] transition-all duration-300 group-hover:w-full" />
    </motion.div>
  )
}
