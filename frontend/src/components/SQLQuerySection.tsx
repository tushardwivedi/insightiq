'use client'

import { useState } from 'react'
import { Database, Play } from 'lucide-react'
import { apiClient } from '@/lib/api'
import { AnalyticsResponse } from '@/types'

interface Props {
  onResult: (result: AnalyticsResponse) => void
  onLoading: (loading: boolean) => void
}

export default function SQLQuerySection({ onResult, onLoading }: Props) {
  const [sql, setSql] = useState("SELECT * FROM sample_sales LIMIT 10")
  const [question, setQuestion] = useState("What insights can you provide from this data?")
  const [isSubmitting, setIsSubmitting] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!sql.trim() || !question.trim() || isSubmitting) return

    const sqlLower = sql.toLowerCase();
    const dangerousPatterns = ['drop', 'delete', 'truncate', 'alter', 'create', 'insert', 'update'];
    const hasDangerousKeywords = dangerousPatterns.some(pattern => sqlLower.includes(pattern));

    if (hasDangerousKeywords) {
      alert('For security reasons, only SELECT queries are allowed');
      return;
    }

    setIsSubmitting(true)
    onLoading(true)

    try {
      const result = await apiClient.sqlQuery({ sql, question })
      onResult(result)
    } catch (error) {
      console.error('SQL query failed:', error)
      const errorMessage = error instanceof Error ? error.message : 'Failed to execute SQL query'
      onResult({
        query: sql,
        data: [],
        insights: `Error: ${errorMessage}. Please check your SQL syntax and try again.`,
        timestamp: new Date().toISOString(),
        process_time: '0ms',
        task_id: 'error',
        status: 'failed'
      })
    } finally {
      setIsSubmitting(false)
      onLoading(false)
    }
  }

  return (
    <div className="card">
      <div className="flex items-center gap-3 mb-4">
        <div className="p-2 rounded-lg" style={{ background: 'var(--accent-color)', color: 'var(--primary-background)' }}>
          <Database className="w-5 h-5" />
        </div>
        <h2 className="text-lg font-semibold" style={{ color: 'var(--text-primary)' }}>Custom SQL Query</h2>
      </div>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium mb-2" style={{ color: 'var(--text-secondary)' }}>
            SQL Query
          </label>
          <textarea
            value={sql}
            onChange={(e) => setSql(e.target.value)}
            className="w-full p-3 rounded-lg font-mono text-sm"
            style={{
              background: 'var(--surface-color)',
              color: 'var(--text-primary)',
              border: '1px solid var(--border-color)',
              minHeight: '100px'
            }}
            rows={4}
            disabled={isSubmitting}
          />
        </div>

        <div>
          <label className="block text-sm font-medium mb-2" style={{ color: 'var(--text-secondary)' }}>
            Analysis Question
          </label>
          <input
            type="text"
            value={question}
            onChange={(e) => setQuestion(e.target.value)}
            className="w-full p-3 rounded-lg"
            style={{
              background: 'var(--surface-color)',
              color: 'var(--text-primary)',
              border: '1px solid var(--border-color)'
            }}
            disabled={isSubmitting}
            placeholder="What insights can you provide from this data?"
          />
        </div>

        <button
          type="submit"
          disabled={!sql.trim() || !question.trim() || isSubmitting}
          className="btn-primary w-full flex items-center justify-center gap-2"
        >
          {isSubmitting ? (
            <>
              <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
              Executing...
            </>
          ) : (
            <>
              <Play className="w-4 h-4" />
              Execute & Analyze
            </>
          )}
        </button>
      </form>
    </div>
  )
}