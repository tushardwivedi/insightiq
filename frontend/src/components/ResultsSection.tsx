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
        className="relative overflow-hidden bg-white rounded-xl shadow-lg p-8 border border-gray-200"
      >
        <div className="absolute inset-0 bg-gradient-to-r from-blue-50 via-purple-50 to-pink-50 opacity-50"></div>
        <div className="relative flex flex-col items-center justify-center space-y-6">
          {/* Modern pulsing dots animation */}
          <div className="flex space-x-2">
            {[0, 1, 2].map((index) => (
              <motion.div
                key={index}
                className="w-3 h-3 bg-gradient-to-r from-blue-500 to-purple-600 rounded-full"
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
              className="text-lg font-medium bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent"
              animate={{ opacity: [0.6, 1, 0.6] }}
              transition={{ duration: 2, repeat: Infinity, ease: "easeInOut" }}
            >
              Analyzing your data with AI
            </motion.span>
            <div className="flex items-center justify-center mt-2">
              <Brain className="w-5 h-5 text-blue-500 mr-2" />
              <motion.div className="flex space-x-1">
                {['●', '●', '●'].map((dot, index) => (
                  <motion.span
                    key={index}
                    className="text-blue-500"
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
        className="relative overflow-hidden bg-white rounded-xl shadow-lg p-8 border border-gray-200 text-center"
      >
        <div className="absolute inset-0 bg-gradient-to-br from-gray-50 via-blue-50 to-purple-50 opacity-50"></div>
        <div className="relative">
          <div className="w-16 h-16 mx-auto mb-4 bg-gradient-to-r from-blue-100 to-purple-100 rounded-full flex items-center justify-center">
            <Brain className="w-8 h-8 text-blue-600" />
          </div>
          <h3 className="text-xl font-semibold bg-gradient-to-r from-gray-800 to-gray-600 bg-clip-text text-transparent mb-2">
            Ready for AI Analysis
          </h3>
          <p className="text-gray-600">Submit a query above to see intelligent data insights and visualizations</p>
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
      <div className="relative overflow-hidden bg-white rounded-xl shadow-lg border border-gray-200">
        <div className="absolute inset-0 bg-gradient-to-r from-indigo-50 via-purple-50 to-pink-50 opacity-60"></div>
        <div className="relative p-6">
          <div className="flex items-center justify-between mb-4">
            <div className="flex items-center gap-3">
              <div className="relative">
                <motion.div
                  initial={{ scale: 0, opacity: 0 }}
                  animate={{ scale: 1, opacity: 1 }}
                  transition={{ type: "spring", duration: 0.8 }}
                  className="p-2 bg-gradient-to-r from-indigo-500 to-purple-600 rounded-lg shadow-lg relative z-10"
                >
                  <Brain className="w-5 h-5 text-white" />
                </motion.div>
                {/* Static glow effect - no animation */}
                <div className="absolute inset-0 bg-gradient-to-r from-indigo-400 to-purple-500 rounded-lg opacity-30 blur-sm scale-110"></div>
              </div>
              <motion.h2
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ delay: 0.3 }}
                className="text-2xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent"
              >
                AI Analysis Results
              </motion.h2>
            </div>

            {/* View Toggle */}
            <div className="flex items-center gap-2 bg-gray-100 rounded-lg p-1">
              <button
                onClick={() => setActiveView('charts')}
                className={`px-4 py-2 rounded-md text-sm font-medium transition-all duration-200 ${
                  activeView === 'charts'
                    ? 'bg-white text-indigo-600 shadow-sm'
                    : 'text-gray-600 hover:text-indigo-600'
                }`}
              >
                <BarChart3 className="w-4 h-4 inline mr-1" />
                Charts
              </button>
              <button
                onClick={() => setActiveView('data')}
                className={`px-4 py-2 rounded-md text-sm font-medium transition-all duration-200 ${
                  activeView === 'data'
                    ? 'bg-white text-indigo-600 shadow-sm'
                    : 'text-gray-600 hover:text-indigo-600'
                }`}
              >
                <Database className="w-4 h-4 inline mr-1" />
                Data
              </button>
            </div>
          </div>

          {/* Metadata */}
          <div className="flex items-center gap-4 text-sm text-gray-600">
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
      </div>

      {/* Voice Transcript */}
      <AnimatePresence>
        {transcript && (
          <motion.div
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: 'auto' }}
            exit={{ opacity: 0, height: 0 }}
            className="relative overflow-hidden bg-white rounded-xl shadow-lg border border-purple-200"
          >
            <div className="absolute inset-0 bg-gradient-to-r from-purple-50 to-pink-50 opacity-60"></div>
            <div className="relative p-6">
              <h3 className="font-semibold text-purple-800 mb-2 flex items-center gap-2">
                <div className="w-2 h-2 bg-purple-500 rounded-full animate-pulse"></div>
                Voice Transcript
              </h3>
              <p className="text-purple-700 italic text-lg">"{transcript}"</p>
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Query Display */}
      <motion.div
        initial={{ opacity: 0, x: -20 }}
        animate={{ opacity: 1, x: 0 }}
        transition={{ delay: 0.2 }}
        className="relative overflow-hidden bg-white rounded-xl shadow-lg border border-blue-200"
      >
        <div className="absolute inset-0 bg-gradient-to-r from-blue-50 to-cyan-50 opacity-60"></div>
        <div className="relative p-6">
          <h3 className="font-semibold text-blue-800 mb-2">Query</h3>
          <p className="text-blue-700 text-lg">{(analyticsData as AnalyticsResponse).query}</p>
        </div>
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
            className="bg-white rounded-xl shadow-lg border border-gray-200 overflow-hidden"
          >
            <div className="p-6">
              <h3 className="font-semibold text-gray-800 mb-4 flex items-center gap-2">
                <Database className="w-5 h-5 text-gray-600" />
                Raw Data ({(analyticsData as AnalyticsResponse).data?.length || 0} records)
              </h3>

              {(analyticsData as AnalyticsResponse).data && (analyticsData as AnalyticsResponse).data.length > 0 ? (
                <div className="max-h-96 overflow-auto border border-gray-200 rounded-lg">
                  <table className="min-w-full divide-y divide-gray-200">
                    <thead className="bg-gray-50 sticky top-0">
                      <tr>
                        {Object.keys((analyticsData as AnalyticsResponse).data[0]).map((key) => (
                          <th key={key} className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            {key.replace(/_/g, ' ')}
                          </th>
                        ))}
                      </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                      {(analyticsData as AnalyticsResponse).data.slice(0, 50).map((row, index) => (
                        <motion.tr
                          key={index}
                          initial={{ opacity: 0 }}
                          animate={{ opacity: 1 }}
                          transition={{ delay: index * 0.02 }}
                          className="hover:bg-gray-50 transition-colors"
                        >
                          {Object.values(row).map((value, cellIndex) => (
                            <td key={cellIndex} className="px-4 py-3 text-sm text-gray-900">
                              {typeof value === 'number' ? value.toLocaleString() : String(value)}
                            </td>
                          ))}
                        </motion.tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              ) : (
                <div className="text-center py-12 text-gray-500">
                  <Database className="w-12 h-12 mx-auto mb-4 text-gray-300" />
                  <p>No data returned from query</p>
                </div>
              )}
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </motion.div>
  )
}