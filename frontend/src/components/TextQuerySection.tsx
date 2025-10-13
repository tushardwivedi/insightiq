'use client'

import { useState } from 'react'
import { Send, MessageSquare } from 'lucide-react'
import { apiClient } from '@/lib/api'
import { AnalyticsResponse } from '@/types'

interface Props {
  onResult: (result: AnalyticsResponse) => void
  onLoading: (loading: boolean) => void
}

export default function TextQuerySection({ onResult, onLoading }: Props) {
  const [query, setQuery] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!query.trim() || isSubmitting) return

    setIsSubmitting(true)
    onLoading(true)

    try {
      const result = await apiClient.textQuery({ query })
      onResult(result)
    } catch (error) {
      console.error('Text query failed:', error)
      // Show error to user
      const errorMessage = error instanceof Error ? error.message : 'Failed to process query'
      onResult({
        query,
        data: [],
        insights: `Error: ${errorMessage}. Please check if the database is accessible and try again.`,
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
          <MessageSquare className="w-5 h-5" />
        </div>
        <h2 className="text-lg font-semibold" style={{ color: 'var(--text-primary)' }}>Text Query</h2>
      </div>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium mb-2" style={{ color: 'var(--text-secondary)' }}>
            Define Your Metric
          </label>
          <textarea
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Ask about your data... e.g., 'Show me sales trends for the last 6 months'"
            className="w-full p-3 rounded-lg resize-none"
            style={{
              background: 'var(--surface-color)',
              color: 'var(--text-primary)',
              border: '1px solid var(--border-color)',
              minHeight: '80px'
            }}
            rows={3}
            disabled={isSubmitting}
          />
        </div>

        <button
          type="submit"
          disabled={!query.trim() || isSubmitting}
          className="btn-primary w-full flex items-center justify-center gap-2"
        >
          {isSubmitting ? (
            <>
              <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
              Processing...
            </>
          ) : (
            <>
              <Send className="w-4 h-4" />
              Send Query
            </>
          )}
        </button>
      </form>

      <div className="mt-4 text-sm" style={{ color: 'var(--text-secondary)' }}>
        <p className="font-medium mb-1">Example queries:</p>
        <ul className="space-y-1">
          <li>• "What are the top performing products?"</li>
          <li>• "Show me revenue trends by month"</li>
          <li>• "Which customer segments are growing fastest?"</li>
        </ul>
      </div>
    </div>
  )
}