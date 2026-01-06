'use client'

import { useState } from 'react'
import { Database, Play, History, Bookmark } from 'lucide-react'
import { apiClient } from '@/lib/api'
import { AnalyticsResponse } from '@/types'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import SQLEditor from './SQLEditor'

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
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="p-2 rounded-lg bg-gradient-to-br from-blue-500 to-purple-600 text-white">
            <Database className="w-5 h-5" />
          </div>
          <div>
            <h3 className="text-lg font-semibold">Query Workspace</h3>
            <p className="text-sm text-muted-foreground">Write and execute SQL queries</p>
          </div>
        </div>

        {/* Quick Actions */}
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm">
            <History className="w-4 h-4 mr-2" />
            History
          </Button>
          <Button variant="outline" size="sm">
            <Bookmark className="w-4 h-4 mr-2" />
            Saved
          </Button>
        </div>
      </div>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="space-y-2">
          <Label htmlFor="sql-query">SQL Query</Label>
          <SQLEditor
            value={sql}
            onChange={setSql}
            disabled={isSubmitting}
          />
        </div>

        <div className="space-y-2">
          <Label htmlFor="analysis-question">Analysis Question</Label>
          <Input
            id="analysis-question"
            type="text"
            value={question}
            onChange={(e) => setQuestion(e.target.value)}
            disabled={isSubmitting}
            placeholder="What insights can you provide from this data?"
          />
        </div>

        <div className="flex gap-3">
          <Button
            type="submit"
            disabled={!sql.trim() || !question.trim() || isSubmitting}
            className="flex-1"
            size="lg"
          >
            {isSubmitting ? (
              <>
                <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                Executing Query...
              </>
            ) : (
              <>
                <Play className="w-4 h-4 mr-2" />
                Execute & Analyze
              </>
            )}
          </Button>
          <Button
            type="button"
            variant="outline"
            onClick={() => {
              setSql("SELECT * FROM sample_sales LIMIT 10")
              setQuestion("What insights can you provide from this data?")
            }}
            disabled={isSubmitting}
            size="lg"
          >
            Clear
          </Button>
        </div>
      </form>
    </div>
  )
}