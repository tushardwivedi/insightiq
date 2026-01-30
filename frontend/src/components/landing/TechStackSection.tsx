'use client'

import { motion } from 'framer-motion'
import { TechLogo } from './TechLogo'

const techStack = ['PostgreSQL', 'Ollama', 'Whisper', 'Qdrant', 'Next.js', 'Go']

export function TechStackSection() {
  return (
    <section className="relative py-16 border-y border-[#4A5568]/30" style={{ background: 'linear-gradient(to right, #1A202C, #2D3748, #1A202C)' }}>
      <div className="container mx-auto px-4">
        <motion.div initial={{ opacity: 0, y: 20 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true }} className="text-center">
          <p className="text-xs uppercase tracking-[0.3em] text-[#A0AEC0] mb-8">Powered by Best-in-Class Open Source Technologies</p>
          <div className="flex flex-wrap items-center justify-center gap-4">
            {techStack.map((tech, index) => (
              <TechLogo key={tech} name={tech} delay={index * 0.1} />
            ))}
          </div>
        </motion.div>
      </div>
    </section>
  )
}
