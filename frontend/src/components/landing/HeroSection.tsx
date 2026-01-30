'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { motion } from 'framer-motion'
import { Brain, Zap, ArrowRight, Sparkles, Play, Lock, Globe, Terminal, Github } from 'lucide-react'
import { AnimatedTextCycle } from './AnimatedTextCycle'
import { AnimatedOrb } from './AnimatedOrb'
import { GlowEffect } from './GlowEffect'
import { GridBackground } from './GridBackground'

export function HeroSection() {
  const [mousePosition, setMousePosition] = useState({ x: 0, y: 0 })

  useEffect(() => {
    const handleMouseMove = (e: MouseEvent) => {
      setMousePosition({ x: e.clientX, y: e.clientY })
    }
    window.addEventListener('mousemove', handleMouseMove)
    return () => window.removeEventListener('mousemove', handleMouseMove)
  }, [])

  return (
    <section className="relative min-h-screen flex items-center justify-center py-20 md:py-32">
      <div className="absolute inset-0 bg-gradient-to-br from-[#0F172A] via-[#1A202C] to-[#0F172A]" />
      <GridBackground />
      <AnimatedOrb className="top-20 left-10" delay={0} size={400} color="teal" />
      <AnimatedOrb className="bottom-40 right-20" delay={1.5} size={350} color="cyan" />
      <AnimatedOrb className="top-1/3 right-1/4" delay={0.8} size={300} color="indigo" />
      <GlowEffect variant="center" />

      <div className="pointer-events-none fixed inset-0 z-10 transition-opacity duration-500" style={{ background: `radial-gradient(800px circle at ${mousePosition.x}px ${mousePosition.y}px, rgba(79, 209, 197, 0.08), transparent 40%)` }} />

      <div className="container mx-auto px-4 relative z-20">
        <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.8 }} className="text-center max-w-5xl mx-auto">
          <motion.div initial={{ opacity: 0, scale: 0.9 }} animate={{ opacity: 1, scale: 1 }} transition={{ delay: 0.2, duration: 0.5 }} className="inline-flex items-center gap-2 px-5 py-2.5 rounded-full mb-8 border border-[#4FD1C5]/30 bg-[#4FD1C5]/10 backdrop-blur-sm">
            <Sparkles className="w-4 h-4 text-[#4FD1C5]" />
            <span className="text-sm font-medium text-[#4FD1C5]">Open Source &bull; Self-Hosted &bull; Privacy-First</span>
          </motion.div>

          <motion.h1 initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.3, duration: 0.8 }} className="text-5xl md:text-6xl lg:text-7xl xl:text-8xl font-bold mb-6 leading-[1.1] tracking-tight">
            <span className="text-white">Transform </span>
            <AnimatedTextCycle words={['Data', 'Analytics', 'Insights', 'Business']} className="bg-gradient-to-r from-[#4FD1C5] via-[#38B2AC] to-[#4FD1C5] bg-clip-text text-transparent" />
            <br />
            <span className="bg-gradient-to-r from-[#4FD1C5] via-[#38B2AC] to-[#4FD1C5] bg-clip-text text-transparent">Into Instant Insights</span>
          </motion.h1>

          <motion.p initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.4, duration: 0.8 }} className="text-lg md:text-xl lg:text-2xl text-[#A0AEC0] mb-10 max-w-3xl mx-auto leading-relaxed">
            The AI-powered analytics platform that runs entirely on your infrastructure. Ask questions in natural language, get answers in seconds.
          </motion.p>

          <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.5, duration: 0.8 }} className="flex flex-col sm:flex-row items-center justify-center gap-4 mb-12">
            <Link href="/login">
              <button className="group px-8 py-4 rounded-xl font-semibold text-lg bg-gradient-to-r from-[#4FD1C5] to-[#38B2AC] text-[#1A202C] shadow-lg shadow-[#4FD1C5]/30 hover:shadow-[#4FD1C5]/50 transition-all duration-300 hover:scale-105 flex items-center gap-2">
                <Zap className="w-5 h-5" />
                Get Started Free
                <ArrowRight className="w-5 h-5 transition-transform group-hover:translate-x-1" />
              </button>
            </Link>
            <a href="https://github.com/tushardwivedi/insightiq" target="_blank" rel="noopener noreferrer" className="px-8 py-4 rounded-xl font-semibold text-lg border-2 border-[#4FD1C5]/30 text-white bg-white/5 backdrop-blur-sm hover:bg-white/10 hover:border-[#4FD1C5]/50 transition-all duration-300 flex items-center gap-2">
              <Github className="w-5 h-5" />
              View on GitHub
            </a>
          </motion.div>

          <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} transition={{ delay: 0.6, duration: 0.8 }} className="flex flex-wrap items-center justify-center gap-8 text-sm text-[#A0AEC0]">
            <div className="flex items-center gap-2"><Lock className="w-4 h-4 text-[#4FD1C5]" /><span>100% On-Premise</span></div>
            <div className="flex items-center gap-2"><Terminal className="w-4 h-4 text-[#4FD1C5]" /><span>Deploy in 5 Minutes</span></div>
            <div className="flex items-center gap-2"><Globe className="w-4 h-4 text-[#4FD1C5]" /><span>No Cloud Required</span></div>
          </motion.div>

          <motion.div initial={{ opacity: 0, y: 40 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.7, duration: 0.8 }} className="mt-16 relative">
            <div className="relative rounded-2xl overflow-hidden border-2 border-[#4FD1C5]/30 shadow-2xl shadow-[#4FD1C5]/20">
              <div className="aspect-video bg-gradient-to-br from-[#1A202C] to-[#2D3748] flex items-center justify-center">
                <div className="text-center p-8">
                  <div className="relative">
                    <Brain className="w-20 h-20 mx-auto mb-4 text-[#4FD1C5]" />
                    <motion.div animate={{ scale: [1, 1.2, 1], opacity: [0.5, 0.8, 0.5] }} transition={{ duration: 2, repeat: Infinity }} className="absolute inset-0 rounded-full bg-[#4FD1C5]/20 blur-xl" />
                  </div>
                  <p className="text-lg text-[#A0AEC0] mb-2">Interactive Dashboard Preview</p>
                  <Link href="/login" className="inline-flex items-center gap-2 text-[#4FD1C5] hover:underline">
                    <Play className="w-4 h-4" />
                    Try it now
                  </Link>
                </div>
              </div>
            </div>
            <div className="absolute -bottom-4 left-1/2 -translate-x-1/2 w-3/4 h-4 bg-gradient-to-r from-transparent via-[#4FD1C5]/50 to-transparent blur-xl" />
          </motion.div>
        </motion.div>
      </div>
    </section>
  )
}
