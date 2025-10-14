'use client'

import Link from 'next/link'
import { Brain, Database, Mic, Code, Shield, Zap, ArrowRight, CheckCircle } from 'lucide-react'
import { motion } from 'framer-motion'

export default function LandingPage() {
  const features = [
    {
      icon: Brain,
      title: 'AI-Powered Insights',
      description: 'Get intelligent answers from your data using natural language queries powered by local LLMs'
    },
    {
      icon: Database,
      title: 'Multiple Data Sources',
      description: 'Connect to PostgreSQL, MySQL, and other databases with built-in connector support'
    },
    {
      icon: Mic,
      title: 'Voice Input',
      description: 'Ask questions using your voice with integrated Whisper speech recognition'
    },
    {
      icon: Code,
      title: 'SQL Syntax Highlighting',
      description: 'Professional code editor with real-time syntax highlighting for SQL queries'
    },
    {
      icon: Shield,
      title: 'Self-Hosted & Secure',
      description: 'Run entirely on your infrastructure. Your data never leaves your servers'
    },
    {
      icon: Zap,
      title: 'Fast & Efficient',
      description: 'Built with Go and React for lightning-fast performance and real-time analytics'
    }
  ]

  return (
    <div className="landing-page">
      {/* Hero Section */}
      <section className="hero-section py-20 md:py-32">
        <div className="container mx-auto px-4">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8 }}
            className="text-center max-w-4xl mx-auto"
          >
            <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full mb-6" style={{ background: 'var(--surface-color)', border: '1px solid var(--border-color)' }}>
              <span className="w-2 h-2 rounded-full animate-pulse" style={{ background: 'var(--accent-color)' }}></span>
              <span className="text-sm" style={{ color: 'var(--text-secondary)' }}>Open Source • Self-Hosted • Privacy-First</span>
            </div>

            <h1 className="text-5xl md:text-7xl font-bold mb-6" style={{ color: 'var(--text-primary)' }}>
              Self-Hosted AI Analytics
              <br />
              <span style={{ color: 'var(--accent-color)' }}>for Your Data</span>
            </h1>

            <p className="text-xl md:text-2xl mb-8" style={{ color: 'var(--text-secondary)' }}>
              Chat with your databases using natural language. Generate insights, visualize data,
              and make decisions faster with AI-powered analytics.
            </p>

            <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
              <Link href="/login">
                <button className="btn-primary-large px-8 py-4 text-lg rounded-lg font-semibold flex items-center gap-2 transition-transform hover:scale-105">
                  Get Started
                  <ArrowRight className="w-5 h-5" />
                </button>
              </Link>
              <a href="#installation" className="px-8 py-4 text-lg rounded-lg font-semibold border transition-all" style={{ borderColor: 'var(--border-color)', color: 'var(--text-primary)' }}>
                View Installation
              </a>
            </div>

            {/* Hero Image Placeholder */}
            <motion.div
              initial={{ opacity: 0, y: 40 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.8, delay: 0.3 }}
              className="mt-16 rounded-2xl overflow-hidden shadow-2xl border"
              style={{ borderColor: 'var(--border-color)' }}
            >
              <div className="aspect-video flex items-center justify-center" style={{ background: 'var(--surface-color)' }}>
                <div className="text-center">
                  <Brain className="w-16 h-16 mx-auto mb-4" style={{ color: 'var(--accent-color)' }} />
                  <p style={{ color: 'var(--text-secondary)' }}>Dashboard Screenshot Coming Soon</p>
                </div>
              </div>
            </motion.div>
          </motion.div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-20" style={{ background: 'var(--surface-color)' }}>
        <div className="container mx-auto px-4">
          <div className="text-center max-w-3xl mx-auto mb-16">
            <h2 className="text-4xl md:text-5xl font-bold mb-4" style={{ color: 'var(--text-primary)' }}>
              Everything You Need
            </h2>
            <p className="text-xl" style={{ color: 'var(--text-secondary)' }}>
              Powerful features for modern data analytics, all running on your infrastructure
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            {features.map((feature, index) => (
              <motion.div
                key={index}
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5, delay: index * 0.1 }}
                viewport={{ once: true }}
                className="feature-card p-6 rounded-xl border transition-all hover:scale-105"
                style={{ background: 'var(--primary-background)', borderColor: 'var(--border-color)' }}
              >
                <div className="w-12 h-12 rounded-lg flex items-center justify-center mb-4" style={{ background: 'var(--accent-color)' }}>
                  <feature.icon className="w-6 h-6" style={{ color: 'var(--primary-background)' }} />
                </div>
                <h3 className="text-xl font-semibold mb-2" style={{ color: 'var(--text-primary)' }}>
                  {feature.title}
                </h3>
                <p style={{ color: 'var(--text-secondary)' }}>
                  {feature.description}
                </p>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* Installation Section */}
      <section id="installation" className="py-20">
        <div className="container mx-auto px-4">
          <div className="max-w-4xl mx-auto">
            <h2 className="text-4xl md:text-5xl font-bold mb-4 text-center" style={{ color: 'var(--text-primary)' }}>
              Get Started in Minutes
            </h2>
            <p className="text-xl text-center mb-12" style={{ color: 'var(--text-secondary)' }}>
              Deploy InsightIQ with a single Docker Compose command
            </p>

            <div className="space-y-8">
              {/* Step 1 */}
              <div className="flex gap-4">
                <div className="flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center font-bold" style={{ background: 'var(--accent-color)', color: 'var(--primary-background)' }}>
                  1
                </div>
                <div className="flex-1">
                  <h3 className="text-xl font-semibold mb-2" style={{ color: 'var(--text-primary)' }}>
                    Clone the repository
                  </h3>
                  <div className="p-4 rounded-lg font-mono text-sm" style={{ background: 'var(--surface-color)', color: 'var(--text-primary)' }}>
                    git clone https://github.com/yourusername/insightiq.git
                    <br />
                    cd insightiq
                  </div>
                </div>
              </div>

              {/* Step 2 */}
              <div className="flex gap-4">
                <div className="flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center font-bold" style={{ background: 'var(--accent-color)', color: 'var(--primary-background)' }}>
                  2
                </div>
                <div className="flex-1">
                  <h3 className="text-xl font-semibold mb-2" style={{ color: 'var(--text-primary)' }}>
                    Configure environment
                  </h3>
                  <div className="p-4 rounded-lg font-mono text-sm" style={{ background: 'var(--surface-color)', color: 'var(--text-primary)' }}>
                    cp .env.example .env
                    <br />
                    # Edit .env with your configuration
                  </div>
                </div>
              </div>

              {/* Step 3 */}
              <div className="flex gap-4">
                <div className="flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center font-bold" style={{ background: 'var(--accent-color)', color: 'var(--primary-background)' }}>
                  3
                </div>
                <div className="flex-1">
                  <h3 className="text-xl font-semibold mb-2" style={{ color: 'var(--text-primary)' }}>
                    Launch with Docker Compose
                  </h3>
                  <div className="p-4 rounded-lg font-mono text-sm" style={{ background: 'var(--surface-color)', color: 'var(--text-primary)' }}>
                    docker compose up -d
                  </div>
                </div>
              </div>

              {/* Step 4 */}
              <div className="flex gap-4">
                <div className="flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center font-bold" style={{ background: 'var(--accent-color)', color: 'var(--primary-background)' }}>
                  4
                </div>
                <div className="flex-1">
                  <h3 className="text-xl font-semibold mb-2" style={{ color: 'var(--text-primary)' }}>
                    Access your dashboard
                  </h3>
                  <div className="p-4 rounded-lg font-mono text-sm" style={{ background: 'var(--surface-color)', color: 'var(--text-primary)' }}>
                    Open http://localhost:3000
                  </div>
                </div>
              </div>
            </div>

            <div className="mt-12 p-6 rounded-xl border" style={{ background: 'var(--surface-color)', borderColor: 'var(--border-color)' }}>
              <h4 className="font-semibold mb-3 flex items-center gap-2" style={{ color: 'var(--text-primary)' }}>
                <CheckCircle className="w-5 h-5" style={{ color: 'var(--accent-color)' }} />
                What's Included
              </h4>
              <ul className="space-y-2" style={{ color: 'var(--text-secondary)' }}>
                <li>• PostgreSQL database with sample data</li>
                <li>• Ollama LLM server for AI-powered insights</li>
                <li>• Whisper ASR for voice input</li>
                <li>• Qdrant vector database for semantic search</li>
                <li>• Go backend API</li>
                <li>• Next.js frontend application</li>
              </ul>
            </div>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-20" style={{ background: 'var(--surface-color)' }}>
        <div className="container mx-auto px-4 text-center">
          <h2 className="text-4xl md:text-5xl font-bold mb-6" style={{ color: 'var(--text-primary)' }}>
            Ready to Transform Your Data?
          </h2>
          <p className="text-xl mb-8 max-w-2xl mx-auto" style={{ color: 'var(--text-secondary)' }}>
            Start analyzing your data with AI today. No cloud required, no data leaves your infrastructure.
          </p>
          <Link href="/login">
            <button className="btn-primary-large px-8 py-4 text-lg rounded-lg font-semibold transition-transform hover:scale-105">
              Get Started Now
            </button>
          </Link>
        </div>
      </section>
    </div>
  )
}
