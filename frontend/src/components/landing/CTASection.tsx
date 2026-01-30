'use client'

import Link from 'next/link'
import { motion } from 'framer-motion'
import { Zap, ArrowRight, Github } from 'lucide-react'
import { AnimatedOrb } from './AnimatedOrb'
import { GlowEffect } from './GlowEffect'

export function CTASection() {
  return (
    <section className="relative py-24 overflow-hidden">
      <div className="absolute inset-0 bg-gradient-to-br from-[#0F172A] via-[#1A202C] to-[#0F172A]" />
      <AnimatedOrb className="top-10 left-1/4" delay={0} size={350} color="teal" />
      <AnimatedOrb className="bottom-10 right-1/4" delay={1} size={300} color="cyan" />
      <GlowEffect variant="center" />

      <div className="container mx-auto px-4 text-center relative z-10">
        <motion.div initial={{ opacity: 0, y: 20 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true }} className="max-w-3xl mx-auto">
          <h2 className="text-4xl md:text-5xl lg:text-6xl font-bold mb-6 text-white">
            Start Analyzing Your Data
            <span className="block mt-2 bg-gradient-to-r from-[#4FD1C5] to-[#38B2AC] bg-clip-text text-transparent">Today</span>
          </h2>
          <p className="text-xl text-[#A0AEC0] mb-10">
            No cloud setup. No vendor lock-in. No data leaving your infrastructure.
            <br />
            <strong className="text-white">Just intelligent analytics.</strong>
          </p>
          <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
            <Link href="/login">
              <button className="group px-10 py-5 rounded-xl font-bold text-xl bg-gradient-to-r from-[#4FD1C5] to-[#38B2AC] text-[#1A202C] shadow-xl shadow-[#4FD1C5]/30 hover:shadow-[#4FD1C5]/50 transition-all duration-300 hover:scale-105 flex items-center gap-2">
                <Zap className="w-6 h-6" />
                Get Started Now
                <ArrowRight className="w-6 h-6 transition-transform group-hover:translate-x-1" />
              </button>
            </Link>
            <a href="https://github.com/tushardwivedi/insightiq" className="text-lg text-[#4FD1C5] hover:underline flex items-center gap-2">
              <Github className="w-5 h-5" />
              View on GitHub
            </a>
          </div>
        </motion.div>
      </div>
    </section>
  )
}
