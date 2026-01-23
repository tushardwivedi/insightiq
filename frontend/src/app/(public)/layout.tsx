'use client'

import Link from 'next/link'
import { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { Menu, X, Github, BookOpen, Zap } from 'lucide-react'

export default function PublicLayout({
  children,
}: {
  children: React.ReactNode
}) {
  const [isScrolled, setIsScrolled] = useState(false)
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)

  useEffect(() => {
    const handleScroll = () => {
      setIsScrolled(window.scrollY > 20)
    }
    window.addEventListener('scroll', handleScroll)
    return () => window.removeEventListener('scroll', handleScroll)
  }, [])

  const navLinks = [
    { href: '/docs', label: 'Documentation', icon: BookOpen },
    { href: 'https://github.com/tushardwivedi/insightiq', label: 'GitHub', icon: Github, external: true },
  ]

  return (
    <div className="min-h-screen flex flex-col bg-[#0F172A]">
      {/* Navigation Bar */}
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
            {/* Logo */}
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

            {/* Desktop Navigation */}
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

            {/* Mobile Menu Button */}
            <button
              onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
              className="md:hidden p-2 text-white hover:text-[#4FD1C5] transition-colors"
            >
              {isMobileMenuOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
            </button>
          </div>
        </div>

        {/* Mobile Menu */}
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

      {/* Main Content */}
      <main className="flex-1">
        {children}
      </main>

      {/* Footer */}
      <footer className="relative bg-[#0F172A] border-t border-[#4FD1C5]/10">
        {/* Gradient Line */}
        <div className="absolute top-0 left-0 right-0 h-px bg-gradient-to-r from-transparent via-[#4FD1C5]/50 to-transparent" />

        <div className="container mx-auto px-4 py-12">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
            {/* Brand */}
            <div className="md:col-span-1">
              <Link href="/" className="flex items-center gap-3 mb-4">
                <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-[#4FD1C5] to-[#38B2AC] flex items-center justify-center">
                  <span className="text-lg font-bold text-[#0F172A]">IQ</span>
                </div>
                <span className="text-xl font-bold text-white">InsightIQ</span>
              </Link>
              <p className="text-sm text-[#A0AEC0] leading-relaxed">
                Self-hosted AI analytics platform. Chat with your databases using natural language, voice, and SQL.
              </p>
            </div>

            {/* Product */}
            <div>
              <h4 className="font-semibold text-white mb-4">Product</h4>
              <ul className="space-y-3 text-sm">
                <li>
                  <Link href="/#features" className="text-[#A0AEC0] hover:text-[#4FD1C5] transition-colors">
                    Features
                  </Link>
                </li>
                <li>
                  <Link href="/#installation" className="text-[#A0AEC0] hover:text-[#4FD1C5] transition-colors">
                    Installation
                  </Link>
                </li>
                <li>
                  <Link href="/docs" className="text-[#A0AEC0] hover:text-[#4FD1C5] transition-colors">
                    Documentation
                  </Link>
                </li>
                <li>
                  <Link href="/docs/api" className="text-[#A0AEC0] hover:text-[#4FD1C5] transition-colors">
                    API Reference
                  </Link>
                </li>
              </ul>
            </div>

            {/* Resources */}
            <div>
              <h4 className="font-semibold text-white mb-4">Resources</h4>
              <ul className="space-y-3 text-sm">
                <li>
                  <a
                    href="https://github.com/tushardwivedi/insightiq"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-[#A0AEC0] hover:text-[#4FD1C5] transition-colors flex items-center gap-2"
                  >
                    <Github className="w-4 h-4" />
                    GitHub
                  </a>
                </li>
                <li>
                  <a
                    href="https://github.com/tushardwivedi/insightiq/issues"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-[#A0AEC0] hover:text-[#4FD1C5] transition-colors"
                  >
                    Report Issues
                  </a>
                </li>
                <li>
                  <a
                    href="https://github.com/tushardwivedi/insightiq/blob/main/LICENSE"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-[#A0AEC0] hover:text-[#4FD1C5] transition-colors"
                  >
                    MIT License
                  </a>
                </li>
              </ul>
            </div>

            {/* Tech Stack */}
            <div>
              <h4 className="font-semibold text-white mb-4">Powered By</h4>
              <div className="flex flex-wrap gap-2">
                {['Next.js', 'Go', 'PostgreSQL', 'Ollama', 'Whisper', 'Qdrant'].map((tech) => (
                  <span
                    key={tech}
                    className="px-3 py-1 text-xs rounded-full bg-[#2D3748]/50 text-[#A0AEC0] border border-[#4A5568]/30"
                  >
                    {tech}
                  </span>
                ))}
              </div>
            </div>
          </div>

          {/* Bottom Bar */}
          <div className="mt-12 pt-8 border-t border-[#4A5568]/30 flex flex-col md:flex-row items-center justify-between gap-4">
            <p className="text-sm text-[#A0AEC0]">
              &copy; {new Date().getFullYear()} InsightIQ. Open Source Software under MIT License.
            </p>
            <div className="flex items-center gap-6">
              <a
                href="https://github.com/tushardwivedi/insightiq"
                target="_blank"
                rel="noopener noreferrer"
                className="text-[#A0AEC0] hover:text-[#4FD1C5] transition-colors"
              >
                <Github className="w-5 h-5" />
              </a>
            </div>
          </div>
        </div>
      </footer>
    </div>
  )
}
