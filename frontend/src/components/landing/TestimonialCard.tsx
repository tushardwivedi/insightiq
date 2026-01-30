'use client'

import { motion } from 'framer-motion'
import { Star } from 'lucide-react'

interface TestimonialCardProps {
  name: string
  role: string
  content: string
  image: string
  rating: number
  delay?: number
}

export function TestimonialCard({ name, role, content, rating, delay = 0 }: TestimonialCardProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true }}
      transition={{ delay, duration: 0.5 }}
      className="p-6 rounded-2xl bg-[#2D3748]/50 border border-[#4FD1C5]/20 backdrop-blur-sm hover:border-[#4FD1C5]/40 transition-all duration-300"
    >
      <div className="flex items-center gap-1 mb-4">
        {[...Array(5)].map((_, i) => (
          <Star key={i} className={`w-4 h-4 ${i < rating ? 'fill-[#4FD1C5] text-[#4FD1C5]' : 'text-[#4A5568]'}`} />
        ))}
      </div>
      <p className="text-[#A0AEC0] mb-6 leading-relaxed italic">&quot;{content}&quot;</p>
      <div className="flex items-center gap-3">
        <div className="w-12 h-12 rounded-full bg-gradient-to-br from-[#4FD1C5] to-[#38B2AC] flex items-center justify-center text-white font-bold">
          {name.charAt(0)}
        </div>
        <div>
          <div className="font-semibold text-white">{name}</div>
          <div className="text-sm text-[#A0AEC0]">{role}</div>
        </div>
      </div>
    </motion.div>
  )
}
