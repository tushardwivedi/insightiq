'use client'

import { motion } from 'framer-motion'
import { TestimonialCard } from './TestimonialCard'
import { AnimatedOrb } from './AnimatedOrb'

const testimonials = [
  { name: 'Sarah Chen', role: 'Data Analyst at TechCorp', content: 'InsightIQ transformed how we work with data. Natural language queries save us hours every day. The voice feature is a game-changer!', image: '', rating: 5 },
  { name: 'Michael Rodriguez', role: 'CTO at DataFlow', content: 'Being self-hosted was crucial for our compliance needs. InsightIQ delivers enterprise features without the enterprise price tag.', image: '', rating: 5 },
  { name: 'Emily Watson', role: 'Business Intelligence Lead', content: 'The AI-powered insights help us spot trends we would have missed. Our entire team loves how easy it is to use.', image: '', rating: 5 }
]

export function TestimonialsSection() {
  return (
    <section className="relative py-24" style={{ background: 'linear-gradient(to bottom, #1A202C, #0F172A)' }}>
      <AnimatedOrb className="bottom-20 right-10" delay={0.5} size={300} color="teal" />

      <div className="container mx-auto px-4 relative z-10">
        <motion.div initial={{ opacity: 0, y: 20 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true }} className="text-center max-w-3xl mx-auto mb-16">
          <h2 className="text-4xl md:text-5xl font-bold mb-4 text-white">
            Loved by <span className="bg-gradient-to-r from-[#4FD1C5] to-[#38B2AC] bg-clip-text text-transparent">Data Teams</span>
          </h2>
          <p className="text-lg text-[#A0AEC0]">See what our users are saying about InsightIQ</p>
        </motion.div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-6xl mx-auto">
          {testimonials.map((testimonial, index) => (
            <TestimonialCard key={testimonial.name} {...testimonial} delay={index * 0.1} />
          ))}
        </div>
      </div>
    </section>
  )
}
