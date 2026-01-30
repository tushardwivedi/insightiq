'use client'

import { motion } from 'framer-motion'

interface AnimatedOrbProps {
  className?: string
  delay?: number
  size?: number
  color?: 'teal' | 'cyan' | 'indigo'
}

const colorMap = {
  teal: 'rgba(79, 209, 197, 0.4)',
  cyan: 'rgba(6, 182, 212, 0.3)',
  indigo: 'rgba(99, 102, 241, 0.25)'
}

export function AnimatedOrb({ className, delay = 0, size = 300, color = 'teal' }: AnimatedOrbProps) {
  return (
    <motion.div
      initial={{ opacity: 0, scale: 0.8 }}
      animate={{ opacity: [0.4, 0.8, 0.4], scale: [0.9, 1.1, 0.9] }}
      transition={{ duration: 8, delay, repeat: Infinity, repeatType: 'reverse', ease: 'easeInOut' }}
      className={`absolute rounded-full blur-3xl pointer-events-none ${className}`}
      style={{ width: size, height: size, background: `radial-gradient(circle, ${colorMap[color]}, transparent)` }}
    />
  )
}
