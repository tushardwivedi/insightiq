'use client'

import Link from 'next/link'
import { Code, ArrowRight, Lock, Database, Send, Trash2 } from 'lucide-react'

export default function APIReferencePage() {
  const endpoints = [
    {
      method: 'POST',
      path: '/api/auth/login',
      description: 'Authenticate user and get JWT token',
      auth: false,
      request: {
        email: 'admin@insightiq.local',
        password: 'admin123456'
      },
      response: {
        success: true,
        token: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...',
        user: {
          id: 'uuid',
          email: 'admin@insightiq.local',
          role: 'admin'
        }
      }
    },
    {
      method: 'POST',
      path: '/api/query',
      description: 'Execute natural language query',
      auth: true,
      request: {
        query: 'Show me top 10 sales by region',
        type: 'text'
      },
      response: {
        query: 'Show me top 10 sales by region',
        data: [
          { region: 'North', sales: 125000 },
          { region: 'South', sales: 98000 }
        ],
        insights: 'AI-generated insights about the data...',
        process_time: '1.2s'
      }
    },
    {
      method: 'GET',
      path: '/api/connectors',
      description: 'List all data connectors',
      auth: true,
      response: {
        data: [
          {
            id: 'uuid',
            name: 'Production DB',
            type: 'superset',
            status: 'connected',
            created_at: '2024-01-01T00:00:00Z'
          }
        ]
      }
    },
    {
      method: 'POST',
      path: '/api/connectors',
      description: 'Create new data connector',
      auth: true,
      request: {
        name: 'My Superset',
        type: 'superset',
        config: {
          url: 'http://localhost:8088',
          username: 'admin',
          password: 'admin'
        }
      },
      response: {
        success: true,
        connector: {
          id: 'uuid',
          name: 'My Superset',
          status: 'connected'
        }
      }
    },
    {
      method: 'GET',
      path: '/api/query-history',
      description: 'Get query history with pagination',
      auth: true,
      response: {
        count: 25,
        data: [
          {
            id: 'uuid',
            query: 'Show sales trends',
            status: 'success',
            created_at: '2024-01-01T00:00:00Z'
          }
        ],
        limit: 10,
        offset: 0
      }
    },
    {
      method: 'DELETE',
      path: '/api/query-history?id={id}',
      description: 'Delete query history entry',
      auth: true,
      response: {
        success: true,
        message: 'Query history deleted successfully'
      }
    }
  ]

  return (
    <div className="container mx-auto px-4 py-16 max-w-6xl">
      {/* Header */}
      <div className="mb-12">
        <Link href="/docs" className="text-sm mb-4 inline-flex items-center gap-2 hover:underline" style={{ color: 'var(--accent-color)' }}>
          ← Back to Documentation
        </Link>
        <h1 className="text-5xl font-bold mb-4" style={{ color: 'var(--text-primary)' }}>
          API Reference
        </h1>
        <p className="text-xl" style={{ color: 'var(--text-secondary)' }}>
          Complete reference for InsightIQ REST API
        </p>
      </div>

      {/* Base URL */}
      <div className="card p-6 mb-12">
        <h2 className="text-xl font-semibold mb-3" style={{ color: 'var(--text-primary)' }}>
          Base URL
        </h2>
        <code className="px-4 py-2 rounded inline-block text-lg" style={{ background: 'var(--hover-surface)', color: 'var(--accent-color)' }}>
          http://localhost:8080
        </code>
      </div>

      {/* Authentication */}
      <section className="mb-12">
        <h2 className="text-3xl font-bold mb-6 flex items-center gap-3" style={{ color: 'var(--text-primary)' }}>
          <Lock className="w-8 h-8" style={{ color: 'var(--accent-color)' }} />
          Authentication
        </h2>
        <div className="card p-6">
          <p className="mb-4" style={{ color: 'var(--text-secondary)' }}>
            Most endpoints require JWT authentication. Include the token in the Authorization header:
          </p>
          <div className="code-block">
            <pre className="p-4 rounded-lg overflow-x-auto" style={{ background: 'var(--hover-surface)' }}>
              <code style={{ color: 'var(--text-primary)' }}>{`Authorization: Bearer <your-jwt-token>`}</code>
            </pre>
          </div>
        </div>
      </section>

      {/* Endpoints */}
      <section>
        <h2 className="text-3xl font-bold mb-6 flex items-center gap-3" style={{ color: 'var(--text-primary)' }}>
          <Database className="w-8 h-8" style={{ color: 'var(--accent-color)' }} />
          Endpoints
        </h2>

        <div className="space-y-6">
          {endpoints.map((endpoint, idx) => (
            <div key={idx} className="card p-6">
              {/* Method and Path */}
              <div className="flex items-center gap-3 mb-4">
                <span
                  className={`px-3 py-1 rounded font-mono text-sm font-bold ${
                    endpoint.method === 'GET' ? 'bg-blue-500/20 text-blue-400' :
                    endpoint.method === 'POST' ? 'bg-green-500/20 text-green-400' :
                    'bg-red-500/20 text-red-400'
                  }`}
                >
                  {endpoint.method}
                </span>
                <code className="text-lg font-mono" style={{ color: 'var(--text-primary)' }}>
                  {endpoint.path}
                </code>
                {endpoint.auth && (
                  <span className="ml-auto flex items-center gap-1 text-xs px-2 py-1 rounded" style={{ background: 'rgba(79, 209, 197, 0.1)', color: 'var(--accent-color)' }}>
                    <Lock className="w-3 h-3" />
                    Auth Required
                  </span>
                )}
              </div>

              {/* Description */}
              <p className="mb-4" style={{ color: 'var(--text-secondary)' }}>
                {endpoint.description}
              </p>

              {/* Request */}
              {endpoint.request && (
                <div className="mb-4">
                  <h4 className="text-sm font-semibold mb-2 flex items-center gap-2" style={{ color: 'var(--text-primary)' }}>
                    <Send className="w-4 h-4" />
                    Request Body
                  </h4>
                  <div className="code-block">
                    <pre className="p-4 rounded-lg overflow-x-auto text-sm" style={{ background: 'var(--hover-surface)' }}>
                      <code style={{ color: 'var(--text-primary)' }}>
                        {JSON.stringify(endpoint.request, null, 2)}
                      </code>
                    </pre>
                  </div>
                </div>
              )}

              {/* Response */}
              <div>
                <h4 className="text-sm font-semibold mb-2 flex items-center gap-2" style={{ color: 'var(--text-primary)' }}>
                  <ArrowRight className="w-4 h-4" />
                  Response
                </h4>
                <div className="code-block">
                  <pre className="p-4 rounded-lg overflow-x-auto text-sm" style={{ background: 'var(--hover-surface)' }}>
                    <code style={{ color: 'var(--text-primary)' }}>
                      {JSON.stringify(endpoint.response, null, 2)}
                    </code>
                  </pre>
                </div>
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* Example Usage */}
      <section className="mt-12">
        <h2 className="text-3xl font-bold mb-6 flex items-center gap-3" style={{ color: 'var(--text-primary)' }}>
          <Code className="w-8 h-8" style={{ color: 'var(--accent-color)' }} />
          Example Usage
        </h2>

        <div className="card p-6">
          <h3 className="text-xl font-semibold mb-4" style={{ color: 'var(--text-primary)' }}>
            JavaScript / TypeScript
          </h3>
          <div className="code-block">
            <pre className="p-4 rounded-lg overflow-x-auto text-sm" style={{ background: 'var(--hover-surface)' }}>
              <code style={{ color: 'var(--text-primary)' }}>{`// Login and get token
const loginResponse = await fetch('http://localhost:8080/api/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    email: 'admin@insightiq.local',
    password: 'admin123456'
  })
})
const { token } = await loginResponse.json()

// Execute query
const queryResponse = await fetch('http://localhost:8080/api/query', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': \`Bearer \${token}\`
  },
  body: JSON.stringify({
    query: 'Show me sales by region',
    type: 'text'
  })
})
const results = await queryResponse.json()
console.log(results.data)`}</code>
            </pre>
          </div>
        </div>

        <div className="card p-6 mt-6">
          <h3 className="text-xl font-semibold mb-4" style={{ color: 'var(--text-primary)' }}>
            cURL
          </h3>
          <div className="code-block">
            <pre className="p-4 rounded-lg overflow-x-auto text-sm" style={{ background: 'var(--hover-surface)' }}>
              <code style={{ color: 'var(--text-primary)' }}>{`# Login
curl -X POST http://localhost:8080/api/auth/login \\
  -H "Content-Type: application/json" \\
  -d '{"email":"admin@insightiq.local","password":"admin123456"}'

# Execute query (replace TOKEN)
curl -X POST http://localhost:8080/api/query \\
  -H "Content-Type: application/json" \\
  -H "Authorization: Bearer TOKEN" \\
  -d '{"query":"Show me sales by region","type":"text"}'`}</code>
            </pre>
          </div>
        </div>
      </section>

      {/* Back to Docs */}
      <div className="mt-12 flex justify-center">
        <Link href="/docs" className="px-6 py-3 rounded-lg font-semibold flex items-center gap-2 transition-all hover:scale-105" style={{ background: 'var(--accent-color)', color: 'var(--primary-background)' }}>
          Back to Documentation
          <ArrowRight className="w-5 h-5" />
        </Link>
      </div>
    </div>
  )
}
