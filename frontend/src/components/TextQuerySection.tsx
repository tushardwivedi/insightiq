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
    <div className="bg-white rounded-xl shadow-lg p-6 border border-gray-200">
      <div className="flex items-center gap-3 mb-4">
        <div className="p-2 bg-blue-100 rounded-lg">
          <MessageSquare className="w-5 h-5 text-blue-600" />
        </div>
        <h2 className="text-xl font-semibold text-gray-800">Text Query</h2>
      </div>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <textarea
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Ask about your data... e.g., 'Show me sales trends for the last 6 months'"
            className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none text-gray-900 bg-white placeholder-gray-500 dark:text-gray-900 dark:bg-white dark:placeholder-gray-500"
            style={{
              color: '#111827',
              backgroundColor: '#ffffff',
              border: '1px solid #d1d5db'
            }}
            rows={3}
            disabled={isSubmitting}
          />
        </div>

        <button
          type="submit"
          disabled={!query.trim() || isSubmitting}
          className="w-full flex items-center justify-center gap-2 bg-gradient-to-r from-blue-500 to-blue-600 text-white py-3 px-4 rounded-lg hover:from-blue-600 hover:to-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
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

      <div className="mt-4 text-sm text-gray-500">
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