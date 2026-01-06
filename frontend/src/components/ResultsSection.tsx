'use client'

import { Clock, Database, Brain, BarChart3, Sparkles } from 'lucide-react'
import { AnalyticsResponse, VoiceResponse } from '@/types'
import { motion, AnimatePresence } from 'framer-motion'
import { useState } from 'react'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Card, CardContent } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import InteractiveCharts from './InteractiveCharts'

interface Props {
  results: AnalyticsResponse | VoiceResponse | null
  loading: boolean
}

export default function ResultsSection({ results, loading }: Props) {
  const [activeView, setActiveView] = useState<'charts' | 'data'>('charts')

  if (loading) {
    return (
      <Card>
        <CardContent className="flex flex-col items-center justify-center space-y-6 py-12">
          {/* Modern pulsing dots animation */}
          <div className="flex space-x-2">
            {[0, 1, 2].map((index) => (
              <motion.div
                key={index}
                className="w-3 h-3 rounded-full bg-primary"
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
            className="text-center space-y-2"
          >
            <motion.span
              className="text-lg font-medium"
              animate={{ opacity: [0.6, 1, 0.6] }}
              transition={{ duration: 2, repeat: Infinity, ease: "easeInOut" }}
            >
              Analyzing your data with AI
            </motion.span>
            <div className="flex items-center justify-center">
              <Brain className="w-5 h-5 mr-2 text-primary" />
              <motion.div className="flex space-x-1">
                {['●', '●', '●'].map((dot, index) => (
                  <motion.span
                    key={index}
                    className="text-primary"
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
        </CardContent>
      </Card>
    )
  }

  if (!results) {
    return (
      <Card>
        <CardContent className="text-center py-12">
          <div className="w-16 h-16 mx-auto mb-4 rounded-full bg-primary flex items-center justify-center">
            <Brain className="w-8 h-8 text-primary-foreground" />
          </div>
          <h3 className="text-xl font-semibold mb-2">
            Ready for AI Analysis
          </h3>
          <p className="text-muted-foreground">
            Submit a query above to see intelligent data insights and visualizations
          </p>
        </CardContent>
      </Card>
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
      <Card>
        <CardContent className="pt-6">
          <div className="flex items-center justify-between mb-4">
            <div className="flex items-center gap-3">
              <motion.div
                initial={{ scale: 0, opacity: 0 }}
                animate={{ scale: 1, opacity: 1 }}
                transition={{ type: "spring", duration: 0.8 }}
                className="p-2 rounded-lg bg-primary text-primary-foreground"
              >
                <Brain className="w-5 h-5" />
              </motion.div>
              <motion.h2
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ delay: 0.3 }}
                className="text-2xl font-bold"
              >
                AI Analysis Results
              </motion.h2>
            </div>

            {/* View Toggle */}
            <Tabs value={activeView} onValueChange={(value) => setActiveView(value as 'charts' | 'data')} className="w-auto">
              <TabsList>
                <TabsTrigger value="charts" className="flex items-center gap-1">
                  <BarChart3 className="w-4 h-4" />
                  Charts
                </TabsTrigger>
                <TabsTrigger value="data" className="flex items-center gap-1">
                  <Database className="w-4 h-4" />
                  Data
                </TabsTrigger>
              </TabsList>
            </Tabs>
          </div>

          {/* Metadata */}
          <div className="flex items-center gap-4 text-sm text-muted-foreground">
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
        </CardContent>
      </Card>

      {/* Voice Transcript */}
      <AnimatePresence>
        {transcript && (
          <motion.div
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: 'auto' }}
            exit={{ opacity: 0, height: 0 }}
          >
            <Card>
              <CardContent className="pt-6">
                <h3 className="font-semibold mb-2 flex items-center gap-2">
                  <div className="w-2 h-2 rounded-full bg-primary animate-pulse"></div>
                  Voice Transcript
                </h3>
                <p className="italic text-lg text-muted-foreground">"{transcript}"</p>
              </CardContent>
            </Card>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Query Display */}
      <motion.div
        initial={{ opacity: 0, x: -20 }}
        animate={{ opacity: 1, x: 0 }}
        transition={{ delay: 0.2 }}
      >
        <Card>
          <CardContent className="pt-6">
            <h3 className="font-semibold mb-2">Query</h3>
            <p className="text-lg text-muted-foreground">{(analyticsData as AnalyticsResponse).query}</p>
          </CardContent>
        </Card>
      </motion.div>

      {/* Main Content Area */}
      <Tabs value={activeView} onValueChange={(value) => setActiveView(value as 'charts' | 'data')}>
        <TabsContent value="charts" className="mt-0">
          <motion.div
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
        </TabsContent>

        <TabsContent value="data" className="mt-0">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0, scale: 0.95 }}
            transition={{ duration: 0.3 }}
          >
            <Card>
              <CardContent className="pt-6">
                <h3 className="font-semibold mb-4 flex items-center gap-2">
                  <Database className="w-5 h-5 text-muted-foreground" />
                  Raw Data ({(analyticsData as AnalyticsResponse).data?.length || 0} records)
                </h3>

                {(analyticsData as AnalyticsResponse).data && (analyticsData as AnalyticsResponse).data.length > 0 ? (
                  <div className="max-h-96 overflow-auto border rounded-lg">
                    <table className="min-w-full divide-y">
                      <thead className="sticky top-0 bg-muted/50">
                        <tr>
                          {Object.keys((analyticsData as AnalyticsResponse).data[0]).map((key) => (
                            <th key={key} className="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground">
                              {key.replace(/_/g, ' ')}
                            </th>
                          ))}
                        </tr>
                      </thead>
                      <tbody className="divide-y">
                        {(analyticsData as AnalyticsResponse).data.slice(0, 50).map((row, index) => (
                          <motion.tr
                            key={index}
                            initial={{ opacity: 0 }}
                            animate={{ opacity: 1 }}
                            transition={{ delay: index * 0.02 }}
                            className="hover:bg-muted/50 transition-colors"
                          >
                            {Object.values(row).map((value, cellIndex) => (
                              <td key={cellIndex} className="px-4 py-3 text-sm">
                                {typeof value === 'number' ? value.toLocaleString() : String(value)}
                              </td>
                            ))}
                          </motion.tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                ) : (
                  <div className="text-center py-12 text-muted-foreground">
                    <Database className="w-12 h-12 mx-auto mb-4 text-border" />
                    <p>No data returned from query</p>
                  </div>
                )}
              </CardContent>
            </Card>
          </motion.div>
        </TabsContent>
      </Tabs>
    </motion.div>
  )
}