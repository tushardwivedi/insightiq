'use client'

import Link from 'next/link'
import { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { Menu, X, Github, BookOpen, Zap } from 'lucide-react'

const navLinks = [
  { href: '/docs', label: 'Documentation', icon: BookOpen },
  { href: 'https://github.com/tushardwivedi/insightiq', label: 'GitHub', icon: Github, external: true },
]

export function Navbar() {
  const [isScrolled, setIsScrolled] = useState(false)
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)

  useEffect(() => {
    const handleScroll = () => {
      setIsScrolled(window.scrollY > 20)
    }
    window.addEventListener('scroll', handleScroll)
    return () => window.removeEventListener('scroll', handleScroll)
  }, [])

  return (
    <motion.nav
      initial={{ y: -100 }}
      animate={{ y: 0 }}
      transition={{ duration: 0.5 }}
      className={`fixed top-0 left-0 right-0 z-50 transition-all duration-300 ${
        isScrolled
          ? 'bg-[#0F172A]/80 backdrop-blur-xl border-b border-[#4FD1C5]/10 shadow-lg shadow-black/20'
          : 'bg-transparent'
      }`}
    >
      <div className="container mx-auto px-4">
        <div className="flex items-center justify-between h-16 md:h-20">
          <Link href="/" className="flex items-center gap-3 group">
            <div className="relative">
              <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-[#4FD1C5] to-[#38B2AC] flex items-center justify-center shadow-lg shadow-[#4FD1C5]/30 group-hover:shadow-[#4FD1C5]/50 transition-all duration-300">
                <span className="text-lg font-bold text-[#0F172A]">IQ</span>
              </div>
              <motion.div
                className="absolute inset-0 rounded-xl bg-[#4FD1C5]/20 blur-lg"
                animate={{ scale: [1, 1.2, 1], opacity: [0.5, 0.8, 0.5] }}
                transition={{ duration: 2, repeat: Infinity }}
              />
            </div>
            <span className="text-xl font-bold text-white">InsightIQ</span>
          </Link>

          <div className="hidden md:flex items-center gap-6">
            {navLinks.map((link) => (
              <Link
                key={link.href}
                href={link.href}
                target={link.external ? '_blank' : undefined}
                rel={link.external ? 'noopener noreferrer' : undefined}
                className="flex items-center gap-2 text-[#A0AEC0] hover:text-white transition-colors duration-300 text-sm font-medium"
              >
                <link.icon className="w-4 h-4" />
                {link.label}
              </Link>
            ))}
            <Link href="/login">
              <button className="flex items-center gap-2 px-5 py-2.5 rounded-xl font-semibold text-sm bg-gradient-to-r from-[#4FD1C5] to-[#38B2AC] text-[#0F172A] shadow-lg shadow-[#4FD1C5]/30 hover:shadow-[#4FD1C5]/50 transition-all duration-300 hover:scale-105">
                <Zap className="w-4 h-4" />
                Get Started
              </button>
            </Link>
          </div>

          <button
            onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
            className="md:hidden p-2 text-white hover:text-[#4FD1C5] transition-colors"
          >
            {isMobileMenuOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
          </button>
        </div>
      </div>

      <AnimatePresence>
        {isMobileMenuOpen && (
          <motion.div
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: 'auto' }}
            exit={{ opacity: 0, height: 0 }}
            transition={{ duration: 0.3 }}
            className="md:hidden bg-[#0F172A]/95 backdrop-blur-xl border-t border-[#4FD1C5]/10"
          >
            <div className="container mx-auto px-4 py-6 space-y-4">
              {navLinks.map((link) => (
                <Link
                  key={link.href}
                  href={link.href}
                  target={link.external ? '_blank' : undefined}
                  rel={link.external ? 'noopener noreferrer' : undefined}
                  onClick={() => setIsMobileMenuOpen(false)}
                  className="flex items-center gap-3 text-[#A0AEC0] hover:text-white transition-colors duration-300 py-2"
                >
                  <link.icon className="w-5 h-5" />
                  {link.label}
                </Link>
              ))}
              <Link href="/login" onClick={() => setIsMobileMenuOpen(false)}>
                <button className="w-full flex items-center justify-center gap-2 px-5 py-3 rounded-xl font-semibold bg-gradient-to-r from-[#4FD1C5] to-[#38B2AC] text-[#0F172A] shadow-lg shadow-[#4FD1C5]/30">
                  <Zap className="w-4 h-4" />
                  Get Started
                </button>
              </Link>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </motion.nav>
  )
}
