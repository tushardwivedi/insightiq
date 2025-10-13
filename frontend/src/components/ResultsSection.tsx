'use client'

import { Clock, Database, Brain, BarChart3, Sparkles } from 'lucide-react'
import { AnalyticsResponse, VoiceResponse } from '@/types'
import { motion, AnimatePresence } from 'framer-motion'
import { useState } from 'react'
import InteractiveCharts from './InteractiveCharts'

interface Props {
  results: AnalyticsResponse | VoiceResponse | null
  loading: boolean
}

export default function ResultsSection({ results, loading }: Props) {
  const [activeView, setActiveView] = useState<'charts' | 'data'>('charts')

  if (loading) {
    return (
      <motion.div
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        className="card"
      >
        <div className="flex flex-col items-center justify-center space-y-6">
          {/* Modern pulsing dots animation */}
          <div className="flex space-x-2">
            {[0, 1, 2].map((index) => (
              <motion.div
                key={index}
                className="w-3 h-3 rounded-full"
                style={{ background: 'var(--accent-color)' }}
                animate={{
                  scale: [1, 1.3, 1],
                  opacity: [0.7, 1, 0.7]
                }}
                transition={{
                  duration: 0.8,
                  repeat: Infinity,
                  delay: index * 0.2,
                  ease: "easeInOut"
                }}
              />
            ))}
          </div>

          {/* Animated text */}
          <motion.div
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.3 }}
            className="text-center"
          >
            <motion.span
              className="text-lg font-medium"
              style={{ color: 'var(--text-primary)' }}
              animate={{ opacity: [0.6, 1, 0.6] }}
              transition={{ duration: 2, repeat: Infinity, ease: "easeInOut" }}
            >
              Analyzing your data with AI
            </motion.span>
            <div className="flex items-center justify-center mt-2">
              <Brain className="w-5 h-5 mr-2" style={{ color: 'var(--accent-color)' }} />
              <motion.div className="flex space-x-1">
                {['●', '●', '●'].map((dot, index) => (
                  <motion.span
                    key={index}
                    style={{ color: 'var(--accent-color)' }}
                    animate={{ opacity: [0.3, 1, 0.3] }}
                    transition={{
                      duration: 1.2,
                      repeat: Infinity,
                      delay: index * 0.3,
                      ease: "easeInOut"
                    }}
                  >
                    {dot}
                  </motion.span>
                ))}
              </motion.div>
            </div>
          </motion.div>
        </div>
      </motion.div>
    )
  }

  if (!results) {
    return (
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="card text-center"
      >
        <div className="relative">
          <div className="w-16 h-16 mx-auto mb-4 rounded-full flex items-center justify-center" style={{ background: 'var(--accent-color)' }}>
            <Brain className="w-8 h-8" style={{ color: 'var(--primary-background)' }} />
          </div>
          <h3 className="text-xl font-semibold mb-2" style={{ color: 'var(--text-primary)' }}>
            Ready for AI Analysis
          </h3>
          <p style={{ color: 'var(--text-secondary)' }}>Submit a query above to see intelligent data insights and visualizations</p>
        </div>
      </motion.div>
    )
  }

  const analyticsData = 'response' in results ? results.response : results
  const transcript = 'transcript' in results ? results.transcript : null

  return (
    <motion.div
      initial={{ opacity: 0, y: 30 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.7, ease: "easeOut" }}
      className="space-y-6"
    >
      {/* Header Section */}
      <div className="card">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-3">
            <div className="relative">
              <motion.div
                initial={{ scale: 0, opacity: 0 }}
                animate={{ scale: 1, opacity: 1 }}
                transition={{ type: "spring", duration: 0.8 }}
                className="p-2 rounded-lg shadow-lg relative z-10"
                style={{ background: 'var(--accent-color)' }}
              >
                <Brain className="w-5 h-5" style={{ color: 'var(--primary-background)' }} />
              </motion.div>
            </div>
            <motion.h2
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: 0.3 }}
              className="text-2xl font-bold"
              style={{ color: 'var(--text-primary)' }}
            >
              AI Analysis Results
            </motion.h2>
          </div>

          {/* View Toggle */}
          <div className="flex items-center gap-2 rounded-lg p-1" style={{ background: 'var(--hover-surface)' }}>
            <button
              onClick={() => setActiveView('charts')}
              className="px-4 py-2 rounded-md text-sm font-medium transition-all duration-200"
              style={{
                background: activeView === 'charts' ? 'var(--accent-color)' : 'transparent',
                color: activeView === 'charts' ? 'var(--primary-background)' : 'var(--text-secondary)'
              }}
            >
              <BarChart3 className="w-4 h-4 inline mr-1" />
              Charts
            </button>
            <button
              onClick={() => setActiveView('data')}
              className="px-4 py-2 rounded-md text-sm font-medium transition-all duration-200"
              style={{
                background: activeView === 'data' ? 'var(--accent-color)' : 'transparent',
                color: activeView === 'data' ? 'var(--primary-background)' : 'var(--text-secondary)'
              }}
            >
              <Database className="w-4 h-4 inline mr-1" />
              Data
            </button>
          </div>
        </div>

        {/* Metadata */}
        <div className="flex items-center gap-4 text-sm" style={{ color: 'var(--text-secondary)' }}>
          <div className="flex items-center gap-1">
            <Clock className="w-4 h-4" />
            {(analyticsData as AnalyticsResponse).process_time || 'N/A'}
          </div>
          <div className="flex items-center gap-1">
            <Database className="w-4 h-4" />
            {(analyticsData as AnalyticsResponse).data?.length || 0} records
          </div>
          {(analyticsData as AnalyticsResponse).timestamp && (
            <div>
              Processed: {new Date((analyticsData as AnalyticsResponse).timestamp!).toLocaleTimeString()}
            </div>
          )}
        </div>
      </div>

      {/* Voice Transcript */}
      <AnimatePresence>
        {transcript && (
          <motion.div
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: 'auto' }}
            exit={{ opacity: 0, height: 0 }}
            className="card"
          >
            <h3 className="font-semibold mb-2 flex items-center gap-2" style={{ color: 'var(--text-primary)' }}>
              <div className="w-2 h-2 rounded-full animate-pulse" style={{ background: 'var(--accent-color)' }}></div>
              Voice Transcript
            </h3>
            <p className="italic text-lg" style={{ color: 'var(--text-secondary)' }}>"{transcript}"</p>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Query Display */}
      <motion.div
        initial={{ opacity: 0, x: -20 }}
        animate={{ opacity: 1, x: 0 }}
        transition={{ delay: 0.2 }}
        className="card"
      >
        <h3 className="font-semibold mb-2" style={{ color: 'var(--text-primary)' }}>Query</h3>
        <p className="text-lg" style={{ color: 'var(--text-secondary)' }}>{(analyticsData as AnalyticsResponse).query}</p>
      </motion.div>

      {/* Main Content Area */}
      <AnimatePresence mode="wait">
        {activeView === 'charts' ? (
          <motion.div
            key="charts"
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0, scale: 0.95 }}
            transition={{ duration: 0.3 }}
          >
            <InteractiveCharts
              data={(analyticsData as AnalyticsResponse).data || []}
              insights={(analyticsData as AnalyticsResponse).insights || ''}
            />
          </motion.div>
        ) : (
          <motion.div
            key="data"
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0, scale: 0.95 }}
            transition={{ duration: 0.3 }}
            className="card overflow-hidden"
          >
            <h3 className="font-semibold mb-4 flex items-center gap-2" style={{ color: 'var(--text-primary)' }}>
              <Database className="w-5 h-5" style={{ color: 'var(--text-secondary)' }} />
              Raw Data ({(analyticsData as AnalyticsResponse).data?.length || 0} records)
            </h3>

            {(analyticsData as AnalyticsResponse).data && (analyticsData as AnalyticsResponse).data.length > 0 ? (
              <div className="max-h-96 overflow-auto border rounded-lg" style={{ borderColor: 'var(--border-color)' }}>
                <table className="min-w-full divide-y" style={{ background: 'var(--surface-color)', borderColor: 'var(--border-color)' }}>
                  <thead className="sticky top-0" style={{ background: 'var(--hover-surface)' }}>
                    <tr>
                      {Object.keys((analyticsData as AnalyticsResponse).data[0]).map((key) => (
                        <th key={key} className="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider" style={{ color: 'var(--text-secondary)' }}>
                          {key.replace(/_/g, ' ')}
                        </th>
                      ))}
                    </tr>
                  </thead>
                  <tbody className="divide-y" style={{ background: 'var(--surface-color)', borderColor: 'var(--border-color)' }}>
                    {(analyticsData as AnalyticsResponse).data.slice(0, 50).map((row, index) => (
                      <motion.tr
                        key={index}
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        transition={{ delay: index * 0.02 }}
                        className="transition-colors"
                        style={{ background: 'var(--surface-color)' }}
                        onMouseEnter={(e) => e.currentTarget.style.background = 'var(--hover-surface)'}
                        onMouseLeave={(e) => e.currentTarget.style.background = 'var(--surface-color)'}
                      >
                        {Object.values(row).map((value, cellIndex) => (
                          <td key={cellIndex} className="px-4 py-3 text-sm" style={{ color: 'var(--text-primary)' }}>
                            {typeof value === 'number' ? value.toLocaleString() : String(value)}
                          </td>
                        ))}
                      </motion.tr>
                    ))}
                  </tbody>
                </table>
              </div>
            ) : (
              <div className="text-center py-12" style={{ color: 'var(--text-secondary)' }}>
                <Database className="w-12 h-12 mx-auto mb-4" style={{ color: 'var(--border-color)' }} />
                <p>No data returned from query</p>
              </div>
            )}
          </motion.div>
        )}
      </AnimatePresence>
    </motion.div>
  )
}