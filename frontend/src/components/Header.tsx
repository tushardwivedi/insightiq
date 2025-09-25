'use client'

import { Brain, BarChart3 } from 'lucide-react'

export default function Header() {
  return (
    <header className="text-center mb-8">
      <div className="flex items-center justify-center gap-4 mb-4">
        <div className="p-3 bg-gradient-to-r from-blue-500 to-purple-600 rounded-full">
          <Brain className="w-8 h-8 text-white" />
        </div>
        <h1 className="text-4xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
          InsightIQ Analytics
        </h1>
        <div className="p-3 bg-gradient-to-r from-purple-500 to-pink-600 rounded-full">
          <BarChart3 className="w-8 h-8 text-white" />
        </div>
      </div>
      <p className="text-lg text-gray-600 max-w-2xl mx-auto">
        Interact with your data using natural language and voice commands. 
        Get AI-powered insights from your analytics data instantly.
      </p>
    </header>
  )
}