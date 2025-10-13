'use client'

import { useEffect, useRef } from 'react'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  Title,
  Tooltip,
  Legend,
  ArcElement,
  RadialLinearScale,
} from 'chart.js'
import { Line, Bar, Doughnut, PolarArea } from 'react-chartjs-2'
import { motion } from 'framer-motion'
import { TrendingUp, BarChart3, PieChart, Target } from 'lucide-react'

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  Title,
  Tooltip,
  Legend,
  ArcElement,
  RadialLinearScale
)

interface ChartData {
  [key: string]: any
}

interface Props {
  data: ChartData[]
  insights: string
}

export default function InteractiveCharts({ data, insights }: Props) {
  const chartOptions = {
    responsive: true,
    interaction: {
      mode: 'index' as const,
      intersect: false,
    },
    plugins: {
      legend: {
        position: 'top' as const,
        labels: {
          font: {
            size: 12,
            family: 'system-ui, sans-serif',
          },
          color: '#374151',
        },
      },
      title: {
        display: false,
      },
      tooltip: {
        backgroundColor: 'rgba(0, 0, 0, 0.8)',
        titleColor: '#fff',
        bodyColor: '#fff',
        cornerRadius: 8,
      },
    },
    scales: {
      x: {
        grid: {
          color: 'rgba(0, 0, 0, 0.05)',
        },
        ticks: {
          color: '#6B7280',
          font: {
            size: 11,
          },
        },
      },
      y: {
        grid: {
          color: 'rgba(0, 0, 0, 0.05)',
        },
        ticks: {
          color: '#6B7280',
          font: {
            size: 11,
          },
        },
      },
    },
    animation: {
      duration: 1500,
      easing: 'easeInOutQuart' as const,
    },
  }

  // Process data for different chart types
  const processDataForCharts = () => {
    if (!data || data.length === 0) {
      console.log('No data available for charts')
      return null
    }

    console.log('Chart data received:', data)

    // Detect data type and create appropriate visualizations
    const keys = Object.keys(data[0])
    console.log('Data keys:', keys)

    const hasQuarter = keys.some(k => k.toLowerCase().includes('quarter'))
    const hasRevenue = keys.some(k => k.toLowerCase().includes('revenue'))
    const hasCategory = keys.some(k => k.toLowerCase().includes('category'))
    const hasMonth = keys.some(k => k.toLowerCase().includes('month'))
    const hasOrders = keys.some(k => k.toLowerCase().includes('order'))
    const hasYear = keys.some(k => k.toLowerCase().includes('year'))
    const hasName = keys.some(k => k.toLowerCase().includes('name'))
    const hasBirths = keys.some(k => k.toLowerCase().includes('births'))
    const hasGender = keys.some(k => k.toLowerCase().includes('gender'))
    const hasGame = keys.some(k => k.toLowerCase().includes('game'))
    const hasSales = keys.some(k => k.toLowerCase().includes('sales'))
    const hasPlatform = keys.some(k => k.toLowerCase().includes('platform'))
    const hasChannel = keys.some(k => k.toLowerCase().includes('channel'))
    const hasMessages = keys.some(k => k.toLowerCase().includes('messages'))
    const hasUsers = keys.some(k => k.toLowerCase().includes('users'))
    const hasState = keys.some(k => k.toLowerCase().includes('state'))
    const hasVaccinated = keys.some(k => k.toLowerCase().includes('vaccinated'))
    const hasPercentage = keys.some(k => k.toLowerCase().includes('percentage'))

    console.log('Data flags:', { hasQuarter, hasRevenue, hasCategory, hasMonth, hasOrders })

    let charts: any = {}

    // USA Birth Names data visualization
    if (hasName && hasBirths && hasGender) {
      // Gender breakdown
      const genderData = data.reduce((acc: any, item) => {
        const gender = item.gender || item.Gender
        const births = item.births || item.Births
        if (!acc[gender]) acc[gender] = 0
        acc[gender] += Number(births) || 0
        return acc
      }, {})

      charts.genderBreakdown = {
        labels: Object.keys(genderData),
        datasets: [
          {
            label: 'Births by Gender',
            data: Object.values(genderData),
            backgroundColor: ['rgba(99, 102, 241, 0.8)', 'rgba(236, 72, 153, 0.8)'],
            borderColor: ['rgb(99, 102, 241)', 'rgb(236, 72, 153)'],
            borderWidth: 2,
            hoverOffset: 4,
          },
        ],
      }

      // Top names chart
      const topNamesData = data.slice(0, 10)
      charts.topNames = {
        labels: topNamesData.map(item => item.name || item.Name),
        datasets: [
          {
            label: 'Number of Births',
            data: topNamesData.map(item => item.births || item.Births),
            backgroundColor: 'rgba(34, 197, 94, 0.8)',
            borderColor: 'rgb(34, 197, 94)',
            borderWidth: 2,
            borderRadius: 8,
          },
        ],
      }
    }

    // Video Game Sales data visualization
    if (hasGame && hasSales) {
      // Platform breakdown
      const platformData = data.reduce((acc: any, item) => {
        const platform = item.platform || item.Platform
        const sales = item.sales || item.Sales
        if (!acc[platform]) acc[platform] = 0
        acc[platform] += Number(sales) || 0
        return acc
      }, {})

      const colors = [
        'rgba(239, 68, 68, 0.8)',
        'rgba(245, 158, 11, 0.8)',
        'rgba(34, 197, 94, 0.8)',
        'rgba(99, 102, 241, 0.8)',
        'rgba(236, 72, 153, 0.8)',
        'rgba(14, 165, 233, 0.8)',
      ]

      charts.platformBreakdown = {
        labels: Object.keys(platformData),
        datasets: [
          {
            label: 'Sales by Platform (Millions)',
            data: Object.values(platformData),
            backgroundColor: colors,
            borderColor: colors.map(c => c.replace('0.8', '1')),
            borderWidth: 2,
            hoverOffset: 4,
          },
        ],
      }

      // Top games chart
      const topGamesData = data.slice(0, 8)
      charts.topGames = {
        labels: topGamesData.map(item => item.game || item.Game),
        datasets: [
          {
            label: 'Sales (Millions)',
            data: topGamesData.map(item => item.sales || item.Sales),
            backgroundColor: 'rgba(99, 102, 241, 0.8)',
            borderColor: 'rgb(99, 102, 241)',
            borderWidth: 2,
            borderRadius: 8,
          },
        ],
      }
    }

    // Slack Dashboard data visualization
    if (hasChannel && hasMessages) {
      // Channel activity breakdown
      const channelData = data.reduce((acc: any, item) => {
        const channel = item.channel || item.Channel
        const messages = item.messages || item.Messages
        if (!acc[channel]) acc[channel] = 0
        acc[channel] += Number(messages) || 0
        return acc
      }, {})

      const colors = [
        'rgba(239, 68, 68, 0.8)',
        'rgba(245, 158, 11, 0.8)',
        'rgba(34, 197, 94, 0.8)',
        'rgba(99, 102, 241, 0.8)',
        'rgba(236, 72, 153, 0.8)',
        'rgba(14, 165, 233, 0.8)',
      ]

      charts.channelActivity = {
        labels: Object.keys(channelData),
        datasets: [
          {
            label: 'Messages by Channel',
            data: Object.values(channelData),
            backgroundColor: colors,
            borderColor: colors.map(c => c.replace('0.8', '1')),
            borderWidth: 2,
            hoverOffset: 4,
          },
        ],
      }

      // User engagement chart
      if (hasUsers) {
        const userEngagement = data.reduce((acc: any, item) => {
          const channel = item.channel || item.Channel
          const users = item.users || item.Users
          if (!acc[channel]) acc[channel] = 0
          acc[channel] += Number(users) || 0
          return acc
        }, {})

        charts.userEngagement = {
          labels: Object.keys(userEngagement),
          datasets: [
            {
              label: 'Active Users',
              data: Object.values(userEngagement),
              backgroundColor: 'rgba(99, 102, 241, 0.8)',
              borderColor: 'rgb(99, 102, 241)',
              borderWidth: 2,
              borderRadius: 8,
            },
          ],
        }
      }
    }

    // COVID Vaccine Dashboard data visualization
    if (hasState && hasVaccinated) {
      // Top states by vaccination
      const topStates = data.slice(0, 10)
      charts.topStates = {
        labels: topStates.map(item => item.state || item.State),
        datasets: [
          {
            label: 'Vaccinated Population',
            data: topStates.map(item => (item.vaccinated || item.Vaccinated) / 1000000), // Convert to millions
            backgroundColor: 'rgba(34, 197, 94, 0.8)',
            borderColor: 'rgb(34, 197, 94)',
            borderWidth: 2,
            borderRadius: 8,
          },
        ],
      }

      // Vaccination percentage chart
      if (hasPercentage) {
        charts.vaccinationPercentage = {
          labels: topStates.map(item => item.state || item.State),
          datasets: [
            {
              label: 'Vaccination Percentage',
              data: topStates.map(item => item.percentage || item.Percentage),
              backgroundColor: 'rgba(99, 102, 241, 0.8)',
              borderColor: 'rgb(99, 102, 241)',
              borderWidth: 2,
              borderRadius: 8,
            },
          ],
        }
      }
    }

    // Time series data (quarters, months, or any time-based data)
    if ((hasQuarter || hasMonth) && hasRevenue) {
      const timeData = data.reduce((acc: any, item) => {
        // Handle different time field names
        const timeKey = item.quarter || item.Quarter || item.month || item.Month || item.month_year
        const revenue = item.total_revenue || item.revenue || item.Revenue
        if (!acc[timeKey]) acc[timeKey] = 0
        acc[timeKey] += Number(revenue) || 0
        return acc
      }, {})

      charts.timeSeries = {
        labels: Object.keys(timeData),
        datasets: [
          {
            label: 'Revenue ($)',
            data: Object.values(timeData),
            borderColor: 'rgb(99, 102, 241)',
            backgroundColor: 'rgba(99, 102, 241, 0.1)',
            borderWidth: 3,
            fill: true,
            tension: 0.4,
            pointBackgroundColor: 'rgb(99, 102, 241)',
            pointBorderColor: 'white',
            pointBorderWidth: 2,
            pointRadius: 6,
          },
        ],
      }
    }

    // Add orders trend if available
    if ((hasQuarter || hasMonth) && hasOrders) {
      const timeData = data.reduce((acc: any, item) => {
        const timeKey = item.quarter || item.Quarter || item.month || item.Month || item.month_year
        const orders = item.orders || item.Orders || item.quantity || item.total_bikes_sold
        if (!acc[timeKey]) acc[timeKey] = 0
        acc[timeKey] += Number(orders) || 0
        return acc
      }, {})

      charts.ordersTrend = {
        labels: Object.keys(timeData),
        datasets: [
          {
            label: 'Orders',
            data: Object.values(timeData),
            borderColor: 'rgb(34, 197, 94)',
            backgroundColor: 'rgba(34, 197, 94, 0.1)',
            borderWidth: 3,
            fill: true,
            tension: 0.4,
            pointBackgroundColor: 'rgb(34, 197, 94)',
            pointBorderColor: 'white',
            pointBorderWidth: 2,
            pointRadius: 6,
          },
        ],
      }
    }

    // Category breakdown
    if (hasCategory && hasRevenue) {
      const categoryData = data.reduce((acc: any, item) => {
        const category = item.bike_category || item.category || item.Category
        const revenue = item.total_revenue || item.revenue || item.Revenue
        if (!acc[category]) acc[category] = 0
        acc[category] += Number(revenue) || 0
        return acc
      }, {})

      const colors = [
        'rgba(239, 68, 68, 0.8)',
        'rgba(245, 158, 11, 0.8)',
        'rgba(34, 197, 94, 0.8)',
        'rgba(99, 102, 241, 0.8)',
        'rgba(236, 72, 153, 0.8)',
        'rgba(14, 165, 233, 0.8)',
      ]

      charts.categoryBreakdown = {
        labels: Object.keys(categoryData),
        datasets: [
          {
            label: 'Revenue by Category',
            data: Object.values(categoryData),
            backgroundColor: colors,
            borderColor: colors.map(c => c.replace('0.8', '1')),
            borderWidth: 2,
            hoverOffset: 4,
          },
        ],
      }

      // Bar chart version
      charts.categoryBar = {
        labels: Object.keys(categoryData),
        datasets: [
          {
            label: 'Revenue ($)',
            data: Object.values(categoryData),
            backgroundColor: colors,
            borderColor: colors.map(c => c.replace('0.8', '1')),
            borderWidth: 2,
            borderRadius: 8,
            borderSkipped: false,
          },
        ],
      }
    }

    // Growth analysis
    if (hasQuarter || hasMonth) {
      const timeKey = hasQuarter ? 'quarter' : 'month'
      const sortedData = [...data].sort((a, b) => {
        const aTime = a[timeKey] || a[timeKey.charAt(0).toUpperCase() + timeKey.slice(1)]
        const bTime = b[timeKey] || b[timeKey.charAt(0).toUpperCase() + timeKey.slice(1)]
        return aTime?.localeCompare(bTime) || 0
      })

      const revenues = sortedData.map(item =>
        Number(item.total_revenue || item.revenue || item.Revenue) || 0
      )
      const labels = sortedData.map(item =>
        item[timeKey] || item[timeKey.charAt(0).toUpperCase() + timeKey.slice(1)]
      )

      charts.growth = {
        labels,
        datasets: [
          {
            label: 'Revenue Trend',
            data: revenues,
            borderColor: 'rgb(34, 197, 94)',
            backgroundColor: 'rgba(34, 197, 94, 0.1)',
            borderWidth: 3,
            fill: true,
            tension: 0.4,
          },
        ],
      }
    }

    return charts
  }

  const charts = processDataForCharts()

  const containerVariants = {
    hidden: { opacity: 0, y: 50 },
    visible: {
      opacity: 1,
      y: 0,
      transition: {
        duration: 0.6,
        staggerChildren: 0.1,
      },
    },
  }

  const itemVariants = {
    hidden: { opacity: 0, scale: 0.95 },
    visible: {
      opacity: 1,
      scale: 1,
      transition: {
        duration: 0.5,
      },
    },
  }

  if (!charts) {
    return (
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        className="text-center py-8 text-gray-500"
      >
        No chart data available
      </motion.div>
    )
  }

  return (
    <motion.div
      variants={containerVariants}
      initial="hidden"
      animate="visible"
      className="space-y-8"
    >
      {/* Header with AI Insights */}
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
          <div className="leading-relaxed" style={{ color: 'var(--text-secondary)' }}>
            {insights || 'AI analysis of your data reveals interesting patterns and trends.'}
          </div>
        </div>
      </motion.div>

      {/* Charts Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {charts.genderBreakdown && (
          <motion.div variants={itemVariants} className="card p-6  ">
            <div className="flex items-center gap-2 mb-4">
              <PieChart className="w-5 h-5 text-purple-600" />
              <h4 className="font-semibold" style={{ color: 'var(--text-primary)' }}>Births by Gender</h4>
            </div>
            <div className="h-64">
              <Doughnut
                data={charts.genderBreakdown}
                options={{
                  ...chartOptions,
                  cutout: '60%',
                  scales: undefined,
                }}
              />
            </div>
          </motion.div>
        )}

        {charts.topNames && (
          <motion.div variants={itemVariants} className="card p-6  ">
            <div className="flex items-center gap-2 mb-4">
              <BarChart3 className="w-5 h-5 text-green-600" />
              <h4 className="font-semibold" style={{ color: 'var(--text-primary)' }}>Top Baby Names</h4>
            </div>
            <div className="h-64">
              <Bar data={charts.topNames} options={chartOptions} />
            </div>
          </motion.div>
        )}

        {charts.platformBreakdown && (
          <motion.div variants={itemVariants} className="card p-6  ">
            <div className="flex items-center gap-2 mb-4">
              <PieChart className="w-5 h-5 text-blue-600" />
              <h4 className="font-semibold" style={{ color: 'var(--text-primary)' }}>Sales by Platform</h4>
            </div>
            <div className="h-64">
              <Doughnut
                data={charts.platformBreakdown}
                options={{
                  ...chartOptions,
                  cutout: '60%',
                  scales: undefined,
                }}
              />
            </div>
          </motion.div>
        )}

        {charts.topGames && (
          <motion.div variants={itemVariants} className="card p-6  ">
            <div className="flex items-center gap-2 mb-4">
              <BarChart3 className="w-5 h-5 text-indigo-600" />
              <h4 className="font-semibold" style={{ color: 'var(--text-primary)' }}>Top Video Games</h4>
            </div>
            <div className="h-64">
              <Bar data={charts.topGames} options={chartOptions} />
            </div>
          </motion.div>
        )}

        {charts.channelActivity && (
          <motion.div variants={itemVariants} className="card p-6  ">
            <div className="flex items-center gap-2 mb-4">
              <PieChart className="w-5 h-5 text-orange-600" />
              <h4 className="font-semibold" style={{ color: 'var(--text-primary)' }}>Channel Activity</h4>
            </div>
            <div className="h-64">
              <Doughnut
                data={charts.channelActivity}
                options={{
                  ...chartOptions,
                  cutout: '60%',
                  scales: undefined,
                }}
              />
            </div>
          </motion.div>
        )}

        {charts.userEngagement && (
          <motion.div variants={itemVariants} className="card p-6  ">
            <div className="flex items-center gap-2 mb-4">
              <BarChart3 className="w-5 h-5 text-cyan-600" />
              <h4 className="font-semibold" style={{ color: 'var(--text-primary)' }}>User Engagement</h4>
            </div>
            <div className="h-64">
              <Bar data={charts.userEngagement} options={chartOptions} />
            </div>
          </motion.div>
        )}

        {charts.topStates && (
          <motion.div variants={itemVariants} className="card p-6  ">
            <div className="flex items-center gap-2 mb-4">
              <BarChart3 className="w-5 h-5 text-emerald-600" />
              <h4 className="font-semibold" style={{ color: 'var(--text-primary)' }}>Vaccination by State</h4>
            </div>
            <div className="h-64">
              <Bar data={charts.topStates} options={chartOptions} />
            </div>
          </motion.div>
        )}

        {charts.vaccinationPercentage && (
          <motion.div variants={itemVariants} className="card p-6  ">
            <div className="flex items-center gap-2 mb-4">
              <BarChart3 className="w-5 h-5 text-violet-600" />
              <h4 className="font-semibold" style={{ color: 'var(--text-primary)' }}>Vaccination Percentage</h4>
            </div>
            <div className="h-64">
              <Bar data={charts.vaccinationPercentage} options={chartOptions} />
            </div>
          </motion.div>
        )}

        {charts.timeSeries && (
          <motion.div variants={itemVariants} className="card p-6  ">
            <div className="flex items-center gap-2 mb-4">
              <TrendingUp className="w-5 h-5 text-blue-600" />
              <h4 className="font-semibold" style={{ color: 'var(--text-primary)' }}>Revenue Trend</h4>
            </div>
            <div className="h-64">
              <Line data={charts.timeSeries} options={chartOptions} />
            </div>
          </motion.div>
        )}

        {charts.categoryBar && (
          <motion.div variants={itemVariants} className="card p-6  ">
            <div className="flex items-center gap-2 mb-4">
              <BarChart3 className="w-5 h-5 text-green-600" />
              <h4 className="font-semibold" style={{ color: 'var(--text-primary)' }}>Category Performance</h4>
            </div>
            <div className="h-64">
              <Bar data={charts.categoryBar} options={chartOptions} />
            </div>
          </motion.div>
        )}

        {charts.categoryBreakdown && (
          <motion.div variants={itemVariants} className="card p-6  ">
            <div className="flex items-center gap-2 mb-4">
              <PieChart className="w-5 h-5 text-purple-600" />
              <h4 className="font-semibold" style={{ color: 'var(--text-primary)' }}>Market Share</h4>
            </div>
            <div className="h-64">
              <Doughnut
                data={charts.categoryBreakdown}
                options={{
                  ...chartOptions,
                  cutout: '60%',
                  scales: undefined,
                }}
              />
            </div>
          </motion.div>
        )}

        {charts.ordersTrend && (
          <motion.div variants={itemVariants} className="card p-6  ">
            <div className="flex items-center gap-2 mb-4">
              <TrendingUp className="w-5 h-5 text-green-600" />
              <h4 className="font-semibold" style={{ color: 'var(--text-primary)' }}>Orders Trend</h4>
            </div>
            <div className="h-64">
              <Line data={charts.ordersTrend} options={chartOptions} />
            </div>
          </motion.div>
        )}

        {charts.growth && (
          <motion.div variants={itemVariants} className="card p-6  ">
            <div className="flex items-center gap-2 mb-4">
              <TrendingUp className="w-5 h-5 text-indigo-600" />
              <h4 className="font-semibold" style={{ color: 'var(--text-primary)' }}>Growth Analysis</h4>
            </div>
            <div className="h-64">
              <Line
                data={charts.growth}
                options={{
                  ...chartOptions,
                  elements: {
                    point: {
                      radius: 8,
                      hoverRadius: 10,
                    },
                  },
                }}
              />
            </div>
          </motion.div>
        )}
      </div>

      {/* Key Metrics Cards */}
      {data && data.length > 0 && (
        <motion.div variants={itemVariants} className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {(() => {
            // Detect data type for appropriate metric calculations
            const keys = Object.keys(data[0])
            const hasRevenue = keys.some(k => k.toLowerCase().includes('revenue'))
            const hasSales = keys.some(k => k.toLowerCase().includes('sales'))
            const hasBirths = keys.some(k => k.toLowerCase().includes('births'))
            const hasMessages = keys.some(k => k.toLowerCase().includes('messages'))
            const hasVaccinated = keys.some(k => k.toLowerCase().includes('vaccinated'))

            // Determine the main value field and its label
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
              // Fallback
              valueField = 'revenue'
              valueLabel = '$'
              totalLabel = 'Total Revenue'
              avgLabel = 'Average Revenue'
              maxLabel = 'Peak Performance'
            }

            const totalValue = data.reduce((sum, item) =>
              sum + (Number(item[valueField] || item.total_revenue || item.revenue || item.Revenue) || 0), 0
            )
            const totalOrders = data.reduce((sum, item) =>
              sum + (Number(item.orders || item.Orders || item.quantity || item.total_bikes_sold || item.users || item.messages) || 0), 0
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
                icon: '💰',
              },
              {
                title: avgLabel,
                value: `${valueLabel}${Math.round(avgValue).toLocaleString()}`,
                color: 'from-blue-500 to-cyan-600',
                icon: '📊',
              },
              {
                title: maxLabel,
                value: `${valueLabel}${Math.round(maxValue).toLocaleString()}`,
                color: 'from-purple-500 to-pink-600',
                icon: '🚀',
              },
            ]

            // Add total orders if data has orders
            if (totalOrders > 0) {
              metrics.push({
                title: 'Total Orders',
                value: totalOrders.toLocaleString(),
                color: 'from-orange-500 to-red-600',
                icon: '📦',
              })
            }

            return metrics.slice(0, 3).map((metric, index) => (
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
                  <h5 className="text-sm font-medium style={{ color: 'var(--text-secondary)' }} mb-1">{metric.title}</h5>
                  <p className="text-2xl font-bold style={{ color: 'var(--text-primary)' }}">{metric.value}</p>
                </div>
              </div>
            ))
          })()}
        </motion.div>
      )}
    </motion.div>
  )
}