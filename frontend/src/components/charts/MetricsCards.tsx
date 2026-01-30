'use client'

import { motion } from 'framer-motion'

interface ChartData {
  [key: string]: any
}

interface MetricsCardsProps {
  data: ChartData[];
}

const itemVariants = {
  hidden: { opacity: 0, scale: 0.95 },
  visible: {
    opacity: 1,
    scale: 1,
    transition: { duration: 0.5 },
  },
};

export function MetricsCards({ data }: MetricsCardsProps) {
  if (!data || data.length === 0) return null;

  const keys = Object.keys(data[0])
  const hasRevenue = keys.some(k => k.toLowerCase().includes('revenue'))
  const hasSales = keys.some(k => k.toLowerCase().includes('sales'))
  const hasBirths = keys.some(k => k.toLowerCase().includes('births'))
  const hasMessages = keys.some(k => k.toLowerCase().includes('messages'))
  const hasVaccinated = keys.some(k => k.toLowerCase().includes('vaccinated'))

  let valueField = '', valueLabel = '', totalLabel = '', avgLabel = '', maxLabel = ''

  if (hasRevenue) {
    valueField = 'revenue'
    valueLabel = '$'
    totalLabel = 'Total Revenue'
    avgLabel = 'Average Revenue'
    maxLabel = 'Peak Performance'
  } else if (hasSales) {
    valueField = 'sales'
    valueLabel = 'M'
    totalLabel = 'Total Sales'
    avgLabel = 'Average Sales'
    maxLabel = 'Top Game Sales'
  } else if (hasBirths) {
    valueField = 'births'
    valueLabel = ''
    totalLabel = 'Total Births'
    avgLabel = 'Average Births'
    maxLabel = 'Most Popular Name'
  } else if (hasMessages) {
    valueField = 'messages'
    valueLabel = ''
    totalLabel = 'Total Messages'
    avgLabel = 'Average Messages'
    maxLabel = 'Most Active Channel'
  } else if (hasVaccinated) {
    valueField = 'vaccinated'
    valueLabel = ''
    totalLabel = 'Total Vaccinated'
    avgLabel = 'Average per State'
    maxLabel = 'Highest State'
  } else {
    valueField = 'revenue'
    valueLabel = '$'
    totalLabel = 'Total Revenue'
    avgLabel = 'Average Revenue'
    maxLabel = 'Peak Performance'
  }

  const totalValue = data.reduce((sum, item) =>
    sum + (Number(item[valueField] || item.total_revenue || item.revenue || item.Revenue) || 0), 0
  )
  const avgValue = totalValue / data.length
  const maxValue = Math.max(...data.map(item =>
    Number(item[valueField] || item.total_revenue || item.revenue || item.Revenue) || 0
  ))

  const metrics = [
    {
      title: totalLabel,
      value: `${valueLabel}${Math.round(totalValue).toLocaleString()}`,
      color: 'from-green-500 to-emerald-600',
      icon: '\u{1F4B0}',
    },
    {
      title: avgLabel,
      value: `${valueLabel}${Math.round(avgValue).toLocaleString()}`,
      color: 'from-blue-500 to-cyan-600',
      icon: '\u{1F4CA}',
    },
    {
      title: maxLabel,
      value: `${valueLabel}${Math.round(maxValue).toLocaleString()}`,
      color: 'from-purple-500 to-pink-600',
      icon: '\u{1F680}',
    },
  ]

  return (
    <motion.div variants={itemVariants} className="grid grid-cols-1 md:grid-cols-3 gap-6">
      {metrics.map((metric, index) => (
        <div
          key={index}
          className="relative overflow-hidden card p-6  "
        >
          <div className={`absolute inset-0 bg-gradient-to-br ${metric.color} opacity-5`}></div>
          <div className="relative">
            <div className="flex items-center justify-between mb-2">
              <span className="text-2xl">{metric.icon}</span>
              <div className={`w-2 h-2 bg-gradient-to-r ${metric.color} rounded-full`}></div>
            </div>
            <h5 className="text-sm font-medium mb-1" style={{ color: 'var(--text-secondary)' }}>{metric.title}</h5>
            <p className="text-2xl font-bold" style={{ color: 'var(--text-primary)' }}>{metric.value}</p>
          </div>
        </div>
      ))}
    </motion.div>
  );
}
