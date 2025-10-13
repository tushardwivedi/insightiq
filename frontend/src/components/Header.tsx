'use client'

import { Brain, BarChart3 } from 'lucide-react'

export default function Header() {
  return (
    <header className="text-center">
      <div className="flex items-center justify-center gap-3 mb-2">
        <div className="p-2 bg-gradient-to-r from-blue-500 to-purple-600 rounded-full">
          <Brain className="w-5 h-5 text-white" />
        </div>
        <h1 className="text-2xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
          InsightIQ Analytics
        </h1>
        <div className="p-2 bg-gradient-to-r from-purple-500 to-pink-600 rounded-full">
          <BarChart3 className="w-5 h-5 text-white" />
        </div>
      </div>
      <p className="text-sm text-gray-600 max-w-2xl mx-auto">
        AI-powered insights from your analytics data instantly.
      </p>
    </header>
  )
}