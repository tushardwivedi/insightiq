'use client'

import { motion } from 'framer-motion'
import { Target } from 'lucide-react'

interface InsightsHeaderProps {
  insights: string;
}

const itemVariants = {
  hidden: { opacity: 0, scale: 0.95 },
  visible: {
    opacity: 1,
    scale: 1,
    transition: { duration: 0.5 },
  },
};

export function InsightsHeader({ insights }: InsightsHeaderProps) {
  return (
    <motion.div variants={itemVariants} className="relative overflow-hidden">
      <div className="absolute inset-0 bg-gradient-to-r from-blue-600/10 via-purple-600/10 to-pink-600/10 rounded-2xl"></div>
      <div className="relative card backdrop-blur-sm rounded-2xl p-6  ">
        <div className="flex items-center gap-3 mb-4">
          <div className="p-2 bg-gradient-to-r from-blue-500 to-purple-600 rounded-lg">
            <Target className="w-5 h-5 text-white" />
          </div>
          <h3 className="text-xl font-semibold" style={{ color: 'var(--text-primary)' }}>
            AI Insights & Visualizations
          </h3>
        </div>
        <div className="leading-relaxed text-base" style={{ color: 'var(--text-primary)' }}>
          {insights || 'AI analysis of your data reveals interesting patterns and trends.'}
        </div>
      </div>
    </motion.div>
  );
}
