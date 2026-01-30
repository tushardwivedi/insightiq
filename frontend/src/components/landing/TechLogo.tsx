'use client'

import { motion } from 'framer-motion'

interface TechLogoProps {
  name: string
  delay?: number
}

export function TechLogo({ name, delay = 0 }: TechLogoProps) {
  return (
    <motion.div initial={{ opacity: 0, y: 10 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true }} transition={{ delay, duration: 0.4 }} className="flex items-center justify-center">
      <div className="rounded-xl border border-[#4FD1C5]/10 bg-[#2D3748]/50 px-6 py-3 backdrop-blur-sm transition-all duration-300 hover:border-[#4FD1C5]/30 hover:bg-[#2D3748]/70">
        <span className="text-sm font-medium text-[#A0AEC0]">{name}</span>
      </div>
    </motion.div>
  )
}
