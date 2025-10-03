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
    <div className="bg-white rounded-xl shadow-lg p-6 border border-gray-200">
      <div className="flex items-center gap-3 mb-4">
        <div className="p-2 bg-green-100 rounded-lg">
          <Database className="w-5 h-5 text-green-600" />
        </div>
        <h2 className="text-xl font-semibold text-gray-800">Custom SQL Query</h2>
      </div>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            SQL Query
          </label>
          <textarea
            value={sql}
            onChange={(e) => setSql(e.target.value)}
            className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent font-mono text-sm text-gray-900 bg-white placeholder-gray-500 dark:text-gray-900 dark:bg-white dark:placeholder-gray-500"
            style={{
              color: '#111827',
              backgroundColor: '#ffffff',
              border: '1px solid #d1d5db'
            }}
            rows={4}
            disabled={isSubmitting}
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Analysis Question
          </label>
          <input
            type="text"
            value={question}
            onChange={(e) => setQuestion(e.target.value)}
            className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent text-gray-900 bg-white placeholder-gray-500 dark:text-gray-900 dark:bg-white dark:placeholder-gray-500"
            style={{
              color: '#111827',
              backgroundColor: '#ffffff',
              border: '1px solid #d1d5db'
            }}
            disabled={isSubmitting}
          />
        </div>

        <button
          type="submit"
          disabled={!sql.trim() || !question.trim() || isSubmitting}
          className="w-full flex items-center justify-center gap-2 bg-gradient-to-r from-green-500 to-green-600 text-white py-3 px-4 rounded-lg hover:from-green-600 hover:to-green-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
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