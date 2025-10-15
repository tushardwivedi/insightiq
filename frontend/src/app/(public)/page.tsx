'use client'

import Link from 'next/link'
import { Brain, Database, Mic, Code, Shield, Zap, ArrowRight, CheckCircle, Github, Sparkles, Clock, Users, TrendingUp, BarChart3, MessageSquare, Lock } from 'lucide-react'
import { motion } from 'framer-motion'

export default function LandingPage() {
  const features = [
    {
      icon: MessageSquare,
      title: 'Natural Language Queries',
      description: 'Ask questions in plain English. No SQL knowledge required. Our AI translates your questions into optimized queries instantly.'
    },
    {
      icon: Brain,
      title: 'AI-Powered Insights',
      description: 'Get intelligent, contextual answers powered by local LLMs. Your AI assistant understands your data and business context.'
    },
    {
      icon: Database,
      title: 'Universal Connectors',
      description: 'Connect to PostgreSQL, MySQL, and more. One platform for all your data sources with unified querying.'
    },
    {
      icon: Mic,
      title: 'Voice Analytics',
      description: 'Speak your questions naturally with integrated Whisper speech recognition. Hands-free data exploration.'
    },
    {
      icon: Code,
      title: 'SQL Transparency',
      description: 'See and edit the generated SQL. Professional code editor with syntax highlighting for full control.'
    },
    {
      icon: Shield,
      title: 'Self-Hosted & Secure',
      description: 'Deploy on your infrastructure. Your data never leaves your servers. Complete privacy and compliance control.'
    },
    {
      icon: Zap,
      title: 'Real-Time Performance',
      description: 'Built with Go and React for lightning-fast analytics. Get insights in milliseconds, not minutes.'
    },
    {
      icon: BarChart3,
      title: 'Interactive Visualizations',
      description: 'Automatic chart generation from your data. Interactive dashboards that update in real-time.'
    }
  ]

  const useCases = [
    {
      title: 'Business Intelligence',
      description: 'Empower your entire team to explore data and make data-driven decisions without SQL knowledge.',
      icon: TrendingUp
    },
    {
      title: 'Data Engineering',
      description: 'Quickly prototype queries, validate data pipelines, and debug issues with AI assistance.',
      icon: Code
    },
    {
      title: 'Executive Dashboards',
      description: 'Ask strategic questions and get instant answers. Voice-enabled for on-the-go insights.',
      icon: Users
    }
  ]

  const comparisonPoints = [
    { traditional: 'Wait for data teams', insightiq: 'Self-service instantly' },
    { traditional: 'Learn SQL syntax', insightiq: 'Ask in plain English' },
    { traditional: 'Cloud dependency', insightiq: 'Run anywhere' },
    { traditional: 'Data leaves premises', insightiq: '100% on-premise' },
    { traditional: 'Expensive licenses', insightiq: 'Open source & free' },
    { traditional: 'Days to insights', insightiq: 'Seconds to insights' }
  ]

  return (
    <div className="landing-page">
      {/* Hero Section */}
      <section className="hero-section py-20 md:py-32 relative overflow-hidden">
        {/* Animated gradient background */}
        <div className="absolute inset-0 bg-gradient-radial"></div>

        <div className="container mx-auto px-4 relative z-10">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8 }}
            className="text-center max-w-5xl mx-auto"
          >
            {/* Badge */}
            <motion.div
              initial={{ opacity: 0, scale: 0.9 }}
              animate={{ opacity: 1, scale: 1 }}
              transition={{ duration: 0.5 }}
              className="inline-flex items-center gap-2 px-4 py-2 rounded-full mb-6 backdrop-blur-sm"
              style={{ background: 'rgba(79, 209, 197, 0.1)', border: '1px solid rgba(79, 209, 197, 0.3)' }}
            >
              <Sparkles className="w-4 h-4" style={{ color: 'var(--accent-color)' }} />
              <span className="text-sm font-medium" style={{ color: 'var(--accent-color)' }}>
                Open Source • Self-Hosted • Privacy-First
              </span>
            </motion.div>

            {/* Main Headline */}
            <h1 className="text-5xl md:text-7xl lg:text-8xl font-bold mb-6 leading-tight" style={{ color: 'var(--text-primary)' }}>
              Chat with Your Data.
              <br />
              <span className="bg-gradient-to-r from-[#4FD1C5] to-[#3BB5AB] bg-clip-text text-transparent">
                Get Instant Insights.
              </span>
            </h1>

            {/* Subheadline */}
            <p className="text-xl md:text-2xl mb-8 max-w-3xl mx-auto leading-relaxed" style={{ color: 'var(--text-secondary)' }}>
              The AI-powered analytics platform that runs entirely on your infrastructure.
              Ask questions in natural language, get answers in seconds. No cloud required.
            </p>

            {/* CTA Buttons */}
            <div className="flex flex-col sm:flex-row items-center justify-center gap-4 mb-12">
              <Link href="/login">
                <button className="btn-primary-large px-8 py-4 text-lg rounded-xl font-semibold flex items-center gap-2 transition-all hover:scale-105 shadow-lg">
                  <Zap className="w-5 h-5" />
                  Get Started Free
                  <ArrowRight className="w-5 h-5" />
                </button>
              </Link>
              <a href="#installation" className="px-8 py-4 text-lg rounded-xl font-semibold border-2 transition-all hover:scale-105 flex items-center gap-2" style={{ borderColor: 'var(--accent-color)', color: 'var(--text-primary)' }}>
                <Github className="w-5 h-5" />
                View on GitHub
              </a>
            </div>

            {/* Social Proof Metrics */}
            <div className="flex flex-wrap items-center justify-center gap-8 text-sm" style={{ color: 'var(--text-secondary)' }}>
              <div className="flex items-center gap-2">
                <Github className="w-4 h-4" />
                <span>Open Source</span>
              </div>
              <div className="flex items-center gap-2">
                <Shield className="w-4 h-4" />
                <span>100% On-Premise</span>
              </div>
              <div className="flex items-center gap-2">
                <Clock className="w-4 h-4" />
                <span>Deploy in 5 Minutes</span>
              </div>
            </div>

            {/* Hero Image/Demo */}
            <motion.div
              initial={{ opacity: 0, y: 40 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.8, delay: 0.3 }}
              className="mt-16 rounded-2xl overflow-hidden border-2 shadow-2xl"
              style={{ borderColor: 'var(--accent-color)' }}
            >
              <div className="aspect-video relative" style={{ background: 'linear-gradient(135deg, #1A202C 0%, #2D3748 100%)' }}>
                {/* Placeholder for dashboard screenshot */}
                <div className="absolute inset-0 flex items-center justify-center">
                  <div className="text-center p-8">
                    <Brain className="w-20 h-20 mx-auto mb-4 animate-pulse" style={{ color: 'var(--accent-color)' }} />
                    <p className="text-lg" style={{ color: 'var(--text-secondary)' }}>
                      Dashboard Preview
                    </p>
                    <p className="text-sm mt-2" style={{ color: 'var(--text-secondary)' }}>
                      Screenshot coming soon - Deploy now to see it in action!
                    </p>
                  </div>
                </div>
              </div>
            </motion.div>
          </motion.div>
        </div>
      </section>

      {/* Social Proof - Powered By */}
      <section className="py-12" style={{ background: 'var(--surface-color)', borderTop: '1px solid var(--border-color)', borderBottom: '1px solid var(--border-color)' }}>
        <div className="container mx-auto px-4">
          <p className="text-center text-sm mb-6" style={{ color: 'var(--text-secondary)' }}>
            POWERED BY BEST-IN-CLASS OPEN SOURCE TECHNOLOGIES
          </p>
          <div className="flex flex-wrap items-center justify-center gap-12 opacity-60">
            <div className="text-xl font-bold" style={{ color: 'var(--text-primary)' }}>PostgreSQL</div>
            <div className="text-xl font-bold" style={{ color: 'var(--text-primary)' }}>Ollama</div>
            <div className="text-xl font-bold" style={{ color: 'var(--text-primary)' }}>OpenAI Whisper</div>
            <div className="text-xl font-bold" style={{ color: 'var(--text-primary)' }}>Qdrant</div>
            <div className="text-xl font-bold" style={{ color: 'var(--text-primary)' }}>Next.js</div>
            <div className="text-xl font-bold" style={{ color: 'var(--text-primary)' }}>Go</div>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-24">
        <div className="container mx-auto px-4">
          <div className="text-center max-w-3xl mx-auto mb-16">
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ duration: 0.6 }}
            >
              <h2 className="text-4xl md:text-5xl lg:text-6xl font-bold mb-4" style={{ color: 'var(--text-primary)' }}>
                Everything You Need for
                <span className="block mt-2 bg-gradient-to-r from-[#4FD1C5] to-[#3BB5AB] bg-clip-text text-transparent">
                  Modern Analytics
                </span>
              </h2>
              <p className="text-xl mt-4" style={{ color: 'var(--text-secondary)' }}>
                A complete analytics platform that respects your privacy and runs on your infrastructure
              </p>
            </motion.div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {features.map((feature, index) => (
              <motion.div
                key={index}
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ duration: 0.5, delay: index * 0.05 }}
                className="feature-card-enhanced p-6 rounded-xl border transition-all group"
                style={{ background: 'var(--surface-color)', borderColor: 'var(--border-color)' }}
              >
                <div className="w-12 h-12 rounded-lg flex items-center justify-center mb-4 group-hover:scale-110 transition-transform" style={{ background: 'rgba(79, 209, 197, 0.1)' }}>
                  <feature.icon className="w-6 h-6" style={{ color: 'var(--accent-color)' }} />
                </div>
                <h3 className="text-lg font-semibold mb-2" style={{ color: 'var(--text-primary)' }}>
                  {feature.title}
                </h3>
                <p className="text-sm leading-relaxed" style={{ color: 'var(--text-secondary)' }}>
                  {feature.description}
                </p>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* Use Cases Section */}
      <section className="py-24" style={{ background: 'var(--surface-color)' }}>
        <div className="container mx-auto px-4">
          <div className="text-center max-w-3xl mx-auto mb-16">
            <h2 className="text-4xl md:text-5xl font-bold mb-4" style={{ color: 'var(--text-primary)' }}>
              Built for Every Team
            </h2>
            <p className="text-xl" style={{ color: 'var(--text-secondary)' }}>
              From analysts to executives, InsightIQ empowers everyone to work with data
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-8 max-w-6xl mx-auto">
            {useCases.map((useCase, index) => (
              <motion.div
                key={index}
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ duration: 0.6, delay: index * 0.1 }}
                className="use-case-card p-8 rounded-xl border-2"
                style={{ background: 'var(--primary-background)', borderColor: 'var(--border-color)' }}
              >
                <useCase.icon className="w-12 h-12 mb-4" style={{ color: 'var(--accent-color)' }} />
                <h3 className="text-2xl font-bold mb-3" style={{ color: 'var(--text-primary)' }}>
                  {useCase.title}
                </h3>
                <p className="text-lg" style={{ color: 'var(--text-secondary)' }}>
                  {useCase.description}
                </p>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* Comparison Section */}
      <section className="py-24">
        <div className="container mx-auto px-4">
          <div className="text-center max-w-3xl mx-auto mb-16">
            <h2 className="text-4xl md:text-5xl font-bold mb-4" style={{ color: 'var(--text-primary)' }}>
              Why Choose InsightIQ?
            </h2>
            <p className="text-xl" style={{ color: 'var(--text-secondary)' }}>
              Modern analytics without the complexity and vendor lock-in
            </p>
          </div>

          <div className="max-w-4xl mx-auto">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {/* Traditional BI Column */}
              <div className="p-6 rounded-xl" style={{ background: 'rgba(252, 129, 129, 0.05)', border: '1px solid rgba(252, 129, 129, 0.2)' }}>
                <h3 className="text-xl font-bold mb-4 flex items-center gap-2" style={{ color: 'var(--error-color)' }}>
                  <span>❌</span> Traditional BI Tools
                </h3>
                <ul className="space-y-3">
                  {comparisonPoints.map((point, index) => (
                    <li key={index} className="flex items-start gap-2" style={{ color: 'var(--text-secondary)' }}>
                      <span className="text-lg">•</span>
                      <span>{point.traditional}</span>
                    </li>
                  ))}
                </ul>
              </div>

              {/* InsightIQ Column */}
              <div className="p-6 rounded-xl" style={{ background: 'rgba(79, 209, 197, 0.05)', border: '2px solid var(--accent-color)' }}>
                <h3 className="text-xl font-bold mb-4 flex items-center gap-2" style={{ color: 'var(--accent-color)' }}>
                  <span>✅</span> InsightIQ
                </h3>
                <ul className="space-y-3">
                  {comparisonPoints.map((point, index) => (
                    <li key={index} className="flex items-start gap-2 font-medium" style={{ color: 'var(--text-primary)' }}>
                      <CheckCircle className="w-5 h-5 flex-shrink-0" style={{ color: 'var(--accent-color)' }} />
                      <span>{point.insightiq}</span>
                    </li>
                  ))}
                </ul>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Installation Section */}
      <section id="installation" className="py-24" style={{ background: 'var(--surface-color)' }}>
        <div className="container mx-auto px-4">
          <div className="max-w-4xl mx-auto">
            <div className="text-center mb-12">
              <h2 className="text-4xl md:text-5xl font-bold mb-4" style={{ color: 'var(--text-primary)' }}>
                Deploy in 5 Minutes
              </h2>
              <p className="text-xl" style={{ color: 'var(--text-secondary)' }}>
                Get started with a single Docker Compose command
              </p>
            </div>

            <div className="space-y-6">
              {[
                { step: 1, title: 'Clone the repository', code: 'git clone https://github.com/yourusername/insightiq.git\ncd insightiq' },
                { step: 2, title: 'Configure environment', code: 'cp .env.example .env\n# Edit .env with your configuration' },
                { step: 3, title: 'Launch with Docker Compose', code: 'docker compose up -d' },
                { step: 4, title: 'Access your dashboard', code: 'Open http://localhost:3000' }
              ].map((item, index) => (
                <motion.div
                  key={index}
                  initial={{ opacity: 0, x: -20 }}
                  whileInView={{ opacity: 1, x: 0 }}
                  viewport={{ once: true }}
                  transition={{ duration: 0.5, delay: index * 0.1 }}
                  className="flex gap-4"
                >
                  <div className="flex-shrink-0 w-10 h-10 rounded-full flex items-center justify-center font-bold text-lg" style={{ background: 'var(--accent-color)', color: 'var(--primary-background)' }}>
                    {item.step}
                  </div>
                  <div className="flex-1">
                    <h3 className="text-xl font-semibold mb-2" style={{ color: 'var(--text-primary)' }}>
                      {item.title}
                    </h3>
                    <div className="p-4 rounded-lg font-mono text-sm whitespace-pre-wrap" style={{ background: '#1e1e1e', color: '#d4d4d4' }}>
                      {item.code}
                    </div>
                  </div>
                </motion.div>
              ))}
            </div>

            <div className="mt-12 p-6 rounded-xl border" style={{ background: 'var(--primary-background)', borderColor: 'var(--border-color)' }}>
              <h4 className="font-semibold mb-3 flex items-center gap-2" style={{ color: 'var(--text-primary)' }}>
                <CheckCircle className="w-5 h-5" style={{ color: 'var(--accent-color)' }} />
                What's Included
              </h4>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-x-8 gap-y-2" style={{ color: 'var(--text-secondary)' }}>
                <div>• PostgreSQL with sample data</div>
                <div>• Next.js frontend</div>
                <div>• Ollama LLM server</div>
                <div>• Go backend API</div>
                <div>• Whisper speech recognition</div>
                <div>• Pre-configured networking</div>
                <div>• Qdrant vector database</div>
                <div>• Health monitoring</div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Final CTA Section */}
      <section className="py-24 relative overflow-hidden">
        <div className="absolute inset-0 bg-gradient-radial opacity-50"></div>
        <div className="container mx-auto px-4 text-center relative z-10">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ duration: 0.6 }}
            className="max-w-3xl mx-auto"
          >
            <h2 className="text-4xl md:text-5xl lg:text-6xl font-bold mb-6" style={{ color: 'var(--text-primary)' }}>
              Start Analyzing Your Data
              <span className="block mt-2 bg-gradient-to-r from-[#4FD1C5] to-[#3BB5AB] bg-clip-text text-transparent">
                in 5 Minutes
              </span>
            </h2>
            <p className="text-xl mb-8" style={{ color: 'var(--text-secondary)' }}>
              No cloud setup. No vendor lock-in. No data leaving your infrastructure.
              <br />
              <strong className="font-semibold" style={{ color: 'var(--text-primary)' }}>Just intelligent analytics.</strong>
            </p>
            <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
              <Link href="/login">
                <button className="btn-primary-large px-10 py-5 text-xl rounded-xl font-bold transition-transform hover:scale-105 shadow-2xl flex items-center gap-2">
                  <Zap className="w-6 h-6" />
                  Get Started Now
                  <ArrowRight className="w-6 h-6" />
                </button>
              </Link>
              <a href="#installation" className="text-lg hover:underline flex items-center gap-2" style={{ color: 'var(--accent-color)' }}>
                <Github className="w-5 h-5" />
                View Installation Guide
              </a>
            </div>
          </motion.div>
        </div>
      </section>
    </div>
  )
}
