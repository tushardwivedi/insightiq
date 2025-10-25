'use client'

import Link from 'next/link'
import { BookOpen, Code, Database, Zap, Shield, ArrowRight, Terminal, CheckCircle } from 'lucide-react'

export default function DocsPage() {
  return (
    <div className="container mx-auto px-4 py-16 max-w-6xl">
      {/* Header */}
      <div className="mb-12">
        <Link href="/" className="text-sm mb-4 inline-flex items-center gap-2 hover:underline" style={{ color: 'var(--accent-color)' }}>
          ← Back to Home
        </Link>
        <h1 className="text-5xl font-bold mb-4" style={{ color: 'var(--text-primary)' }}>
          Documentation
        </h1>
        <p className="text-xl" style={{ color: 'var(--text-secondary)' }}>
          Everything you need to get started with InsightIQ
        </p>
      </div>

      {/* Quick Links Grid */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-16">
        <Link href="/docs#installation" className="card p-6 hover:shadow-xl transition-all group">
          <Terminal className="w-10 h-10 mb-4" style={{ color: 'var(--accent-color)' }} />
          <h3 className="text-xl font-bold mb-2 group-hover:underline" style={{ color: 'var(--text-primary)' }}>
            Installation
          </h3>
          <p className="text-sm" style={{ color: 'var(--text-secondary)' }}>
            Get up and running in 5 minutes with Docker
          </p>
        </Link>

        <Link href="/docs/api" className="card p-6 hover:shadow-xl transition-all group">
          <Code className="w-10 h-10 mb-4" style={{ color: 'var(--accent-color)' }} />
          <h3 className="text-xl font-bold mb-2 group-hover:underline" style={{ color: 'var(--text-primary)' }}>
            API Reference
          </h3>
          <p className="text-sm" style={{ color: 'var(--text-secondary)' }}>
            Complete API documentation and examples
          </p>
        </Link>

        <Link href="/docs#architecture" className="card p-6 hover:shadow-xl transition-all group">
          <Database className="w-10 h-10 mb-4" style={{ color: 'var(--accent-color)' }} />
          <h3 className="text-xl font-bold mb-2 group-hover:underline" style={{ color: 'var(--text-primary)' }}>
            Architecture
          </h3>
          <p className="text-sm" style={{ color: 'var(--text-secondary)' }}>
            Understand how InsightIQ works under the hood
          </p>
        </Link>
      </div>

      {/* Installation Section */}
      <section id="installation" className="mb-16 scroll-mt-20">
        <h2 className="text-3xl font-bold mb-6 flex items-center gap-3" style={{ color: 'var(--text-primary)' }}>
          <Terminal className="w-8 h-8" style={{ color: 'var(--accent-color)' }} />
          Installation
        </h2>

        <div className="card p-8 mb-6">
          <h3 className="text-xl font-semibold mb-4" style={{ color: 'var(--text-primary)' }}>
            Quick Start with Docker
          </h3>
          <p className="mb-4" style={{ color: 'var(--text-secondary)' }}>
            The fastest way to get InsightIQ running is with Docker Compose:
          </p>

          <div className="code-block mb-6">
            <pre className="p-4 rounded-lg overflow-x-auto" style={{ background: 'var(--hover-surface)' }}>
              <code style={{ color: 'var(--text-primary)' }}>{`# Clone the repository
git clone https://github.com/yourusername/insightiq.git
cd insightiq

# Start all services
docker compose up -d

# Access InsightIQ
open http://localhost:3000`}</code>
            </pre>
          </div>

          <div className="flex items-start gap-3 p-4 rounded-lg mb-4" style={{ background: 'rgba(79, 209, 197, 0.1)', border: '1px solid rgba(79, 209, 197, 0.3)' }}>
            <CheckCircle className="w-5 h-5 mt-0.5" style={{ color: 'var(--accent-color)' }} />
            <div>
              <p className="font-semibold mb-1" style={{ color: 'var(--text-primary)' }}>Default Credentials</p>
              <p className="text-sm" style={{ color: 'var(--text-secondary)' }}>
                Email: <code className="px-2 py-1 rounded" style={{ background: 'var(--hover-surface)' }}>admin@insightiq.local</code>
                <br />
                Password: <code className="px-2 py-1 rounded" style={{ background: 'var(--hover-surface)' }}>admin123456</code>
              </p>
            </div>
          </div>
        </div>

        <div className="card p-8">
          <h3 className="text-xl font-semibold mb-4" style={{ color: 'var(--text-primary)' }}>
            System Requirements
          </h3>
          <ul className="space-y-2" style={{ color: 'var(--text-secondary)' }}>
            <li className="flex items-center gap-2">
              <CheckCircle className="w-5 h-5" style={{ color: 'var(--accent-color)' }} />
              Docker 20.10+ and Docker Compose
            </li>
            <li className="flex items-center gap-2">
              <CheckCircle className="w-5 h-5" style={{ color: 'var(--accent-color)' }} />
              4GB RAM minimum (8GB recommended)
            </li>
            <li className="flex items-center gap-2">
              <CheckCircle className="w-5 h-5" style={{ color: 'var(--accent-color)' }} />
              10GB disk space for services and models
            </li>
          </ul>
        </div>
      </section>

      {/* Architecture Section */}
      <section id="architecture" className="mb-16 scroll-mt-20">
        <h2 className="text-3xl font-bold mb-6 flex items-center gap-3" style={{ color: 'var(--text-primary)' }}>
          <Database className="w-8 h-8" style={{ color: 'var(--accent-color)' }} />
          Architecture
        </h2>

        <div className="card p-8">
          <p className="mb-6" style={{ color: 'var(--text-secondary)' }}>
            InsightIQ is built with a modern, scalable architecture:
          </p>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="p-4 rounded-lg" style={{ background: 'var(--hover-surface)' }}>
              <h4 className="font-semibold mb-2 flex items-center gap-2" style={{ color: 'var(--text-primary)' }}>
                <Code className="w-5 h-5" style={{ color: 'var(--accent-color)' }} />
                Frontend
              </h4>
              <ul className="text-sm space-y-1" style={{ color: 'var(--text-secondary)' }}>
                <li>• Next.js 14 (App Router)</li>
                <li>• React 18 with TypeScript</li>
                <li>• Tailwind CSS</li>
                <li>• Chart.js for visualizations</li>
              </ul>
            </div>

            <div className="p-4 rounded-lg" style={{ background: 'var(--hover-surface)' }}>
              <h4 className="font-semibold mb-2 flex items-center gap-2" style={{ color: 'var(--text-primary)' }}>
                <Zap className="w-5 h-5" style={{ color: 'var(--accent-color)' }} />
                Backend
              </h4>
              <ul className="text-sm space-y-1" style={{ color: 'var(--text-secondary)' }}>
                <li>• Go with Chi Router</li>
                <li>• PostgreSQL for metadata</li>
                <li>• Redis for caching</li>
                <li>• JWT authentication</li>
              </ul>
            </div>

            <div className="p-4 rounded-lg" style={{ background: 'var(--hover-surface)' }}>
              <h4 className="font-semibold mb-2 flex items-center gap-2" style={{ color: 'var(--text-primary)' }}>
                <BookOpen className="w-5 h-5" style={{ color: 'var(--accent-color)' }} />
                AI Services
              </h4>
              <ul className="text-sm space-y-1" style={{ color: 'var(--text-secondary)' }}>
                <li>• Ollama (LLM inference)</li>
                <li>• Whisper (Speech-to-text)</li>
                <li>• Qdrant (Vector search)</li>
              </ul>
            </div>

            <div className="p-4 rounded-lg" style={{ background: 'var(--hover-surface)' }}>
              <h4 className="font-semibold mb-2 flex items-center gap-2" style={{ color: 'var(--text-primary)' }}>
                <Shield className="w-5 h-5" style={{ color: 'var(--accent-color)' }} />
                Security
              </h4>
              <ul className="text-sm space-y-1" style={{ color: 'var(--text-secondary)' }}>
                <li>• SuperTokens for OAuth</li>
                <li>• bcrypt password hashing</li>
                <li>• JWT token authentication</li>
                <li>• Role-based access control</li>
              </ul>
            </div>
          </div>
        </div>
      </section>

      {/* Configuration Section */}
      <section id="configuration" className="mb-16 scroll-mt-20">
        <h2 className="text-3xl font-bold mb-6" style={{ color: 'var(--text-primary)' }}>
          Configuration
        </h2>

        <div className="card p-8">
          <h3 className="text-xl font-semibold mb-4" style={{ color: 'var(--text-primary)' }}>
            Environment Variables
          </h3>
          <p className="mb-4" style={{ color: 'var(--text-secondary)' }}>
            Configure InsightIQ by editing the <code className="px-2 py-1 rounded" style={{ background: 'var(--hover-surface)' }}>.env</code> file:
          </p>

          <div className="code-block">
            <pre className="p-4 rounded-lg overflow-x-auto text-sm" style={{ background: 'var(--hover-surface)' }}>
              <code style={{ color: 'var(--text-primary)' }}>{`# Database
POSTGRES_URL=postgresql://user:password@postgres:5432/insightiq

# Authentication
JWT_SECRET=your-secret-key-change-this
SESSION_DURATION=24h

# Services
OLLAMA_URL=http://ollama:11434
WHISPER_URL=http://whisper:9000
REDIS_URL=redis://redis:6379`}</code>
            </pre>
          </div>
        </div>
      </section>

      {/* Next Steps */}
      <section className="card p-8">
        <h2 className="text-2xl font-bold mb-4" style={{ color: 'var(--text-primary)' }}>
          Next Steps
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <Link href="/docs/api" className="p-4 rounded-lg flex items-center justify-between group hover:shadow-lg transition-all" style={{ background: 'var(--hover-surface)' }}>
            <div className="flex items-center gap-3">
              <Code className="w-6 h-6" style={{ color: 'var(--accent-color)' }} />
              <span className="font-semibold group-hover:underline" style={{ color: 'var(--text-primary)' }}>
                Explore API Reference
              </span>
            </div>
            <ArrowRight className="w-5 h-5" style={{ color: 'var(--text-secondary)' }} />
          </Link>

          <Link href="/login" className="p-4 rounded-lg flex items-center justify-between group hover:shadow-lg transition-all" style={{ background: 'var(--accent-color)' }}>
            <div className="flex items-center gap-3">
              <Zap className="w-6 h-6" style={{ color: 'var(--primary-background)' }} />
              <span className="font-semibold" style={{ color: 'var(--primary-background)' }}>
                Start Using InsightIQ
              </span>
            </div>
            <ArrowRight className="w-5 h-5" style={{ color: 'var(--primary-background)' }} />
          </Link>
        </div>
      </section>
    </div>
  )
}
