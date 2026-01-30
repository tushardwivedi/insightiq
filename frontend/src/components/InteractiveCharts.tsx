'use client'

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
import { Line, Bar, Doughnut } from 'react-chartjs-2'
import { motion } from 'framer-motion'
import { TrendingUp, BarChart3, PieChart } from 'lucide-react'
import { InsightsHeader } from '@/components/charts/InsightsHeader'
import { MetricsCards } from '@/components/charts/MetricsCards'
import { ChartCard } from '@/components/charts/ChartCard'

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
          color: getComputedStyle(document.documentElement).getPropertyValue('--text-primary') || '#374151',
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
          color: getComputedStyle(document.documentElement).getPropertyValue('--border-color') || 'rgba(0, 0, 0, 0.05)',
        },
        ticks: {
          color: getComputedStyle(document.documentElement).getPropertyValue('--text-secondary') || '#6B7280',
          font: {
            size: 11,
          },
        },
      },
      y: {
        grid: {
          color: getComputedStyle(document.documentElement).getPropertyValue('--border-color') || 'rgba(0, 0, 0, 0.05)',
        },
        ticks: {
          color: getComputedStyle(document.documentElement).getPropertyValue('--text-secondary') || '#6B7280',
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

  const doughnutOptions = {
    ...chartOptions,
    cutout: '60%',
    scales: undefined,
  }

  const growthOptions = {
    ...chartOptions,
    elements: {
      point: {
        radius: 8,
        hoverRadius: 10,
      },
    },
  }

  const processDataForCharts = () => {
    if (!data || data.length === 0) {
      return null
    }

    const keys = Object.keys(data[0])

    const hasQuarter = keys.some(k => k.toLowerCase().includes('quarter'))
    const hasRevenue = keys.some(k => k.toLowerCase().includes('revenue'))
    const hasCategory = keys.some(k => k.toLowerCase().includes('category'))
    const hasMonth = keys.some(k => k.toLowerCase().includes('month'))
    const hasOrders = keys.some(k => k.toLowerCase().includes('order'))
    const hasName = keys.some(k => k.toLowerCase().includes('name'))
    const hasBirths = keys.some(k => k.toLowerCase().includes('births'))
    const hasGender = keys.some(k => k.toLowerCase().includes('gender'))
    const hasGame = keys.some(k => k.toLowerCase().includes('game'))
    const hasSales = keys.some(k => k.toLowerCase().includes('sales'))
    const hasChannel = keys.some(k => k.toLowerCase().includes('channel'))
    const hasMessages = keys.some(k => k.toLowerCase().includes('messages'))
    const hasUsers = keys.some(k => k.toLowerCase().includes('users'))
    const hasState = keys.some(k => k.toLowerCase().includes('state'))
    const hasVaccinated = keys.some(k => k.toLowerCase().includes('vaccinated'))
    const hasPercentage = keys.some(k => k.toLowerCase().includes('percentage'))

    let charts: any = {}

    if (hasName && hasBirths && hasGender) {
      const genderData = data.reduce((acc: any, item) => {
        const gender = item.gender || item.Gender
        const births = item.births || item.Births
        if (!acc[gender]) acc[gender] = 0
        acc[gender] += Number(births) || 0
        return acc
      }, {})

      charts.genderBreakdown = {
        labels: Object.keys(genderData),
        datasets: [{
          label: 'Births by Gender',
          data: Object.values(genderData),
          backgroundColor: ['rgba(99, 102, 241, 0.8)', 'rgba(236, 72, 153, 0.8)'],
          borderColor: ['rgb(99, 102, 241)', 'rgb(236, 72, 153)'],
          borderWidth: 2,
          hoverOffset: 4,
        }],
      }

      const topNamesData = data.slice(0, 10)
      charts.topNames = {
        labels: topNamesData.map(item => item.name || item.Name),
        datasets: [{
          label: 'Number of Births',
          data: topNamesData.map(item => item.births || item.Births),
          backgroundColor: 'rgba(34, 197, 94, 0.8)',
          borderColor: 'rgb(34, 197, 94)',
          borderWidth: 2,
          borderRadius: 8,
        }],
      }
    }

    if (hasGame && hasSales) {
      const platformData = data.reduce((acc: any, item) => {
        const platform = item.platform || item.Platform
        const sales = item.sales || item.Sales
        if (!acc[platform]) acc[platform] = 0
        acc[platform] += Number(sales) || 0
        return acc
      }, {})

      const colors = [
        'rgba(239, 68, 68, 0.8)', 'rgba(245, 158, 11, 0.8)',
        'rgba(34, 197, 94, 0.8)', 'rgba(99, 102, 241, 0.8)',
        'rgba(236, 72, 153, 0.8)', 'rgba(14, 165, 233, 0.8)',
      ]

      charts.platformBreakdown = {
        labels: Object.keys(platformData),
        datasets: [{
          label: 'Sales by Platform (Millions)',
          data: Object.values(platformData),
          backgroundColor: colors,
          borderColor: colors.map(c => c.replace('0.8', '1')),
          borderWidth: 2,
          hoverOffset: 4,
        }],
      }

      const topGamesData = data.slice(0, 8)
      charts.topGames = {
        labels: topGamesData.map(item => item.game || item.Game),
        datasets: [{
          label: 'Sales (Millions)',
          data: topGamesData.map(item => item.sales || item.Sales),
          backgroundColor: 'rgba(99, 102, 241, 0.8)',
          borderColor: 'rgb(99, 102, 241)',
          borderWidth: 2,
          borderRadius: 8,
        }],
      }
    }

    if (hasChannel && hasMessages) {
      const channelData = data.reduce((acc: any, item) => {
        const channel = item.channel || item.Channel
        const messages = item.messages || item.Messages
        if (!acc[channel]) acc[channel] = 0
        acc[channel] += Number(messages) || 0
        return acc
      }, {})

      const colors = [
        'rgba(239, 68, 68, 0.8)', 'rgba(245, 158, 11, 0.8)',
        'rgba(34, 197, 94, 0.8)', 'rgba(99, 102, 241, 0.8)',
        'rgba(236, 72, 153, 0.8)', 'rgba(14, 165, 233, 0.8)',
      ]

      charts.channelActivity = {
        labels: Object.keys(channelData),
        datasets: [{
          label: 'Messages by Channel',
          data: Object.values(channelData),
          backgroundColor: colors,
          borderColor: colors.map(c => c.replace('0.8', '1')),
          borderWidth: 2,
          hoverOffset: 4,
        }],
      }

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
          datasets: [{
            label: 'Active Users',
            data: Object.values(userEngagement),
            backgroundColor: 'rgba(99, 102, 241, 0.8)',
            borderColor: 'rgb(99, 102, 241)',
            borderWidth: 2,
            borderRadius: 8,
          }],
        }
      }
    }

    if (hasState && hasVaccinated) {
      const topStates = data.slice(0, 10)
      charts.topStates = {
        labels: topStates.map(item => item.state || item.State),
        datasets: [{
          label: 'Vaccinated Population',
          data: topStates.map(item => (item.vaccinated || item.Vaccinated) / 1000000),
          backgroundColor: 'rgba(34, 197, 94, 0.8)',
          borderColor: 'rgb(34, 197, 94)',
          borderWidth: 2,
          borderRadius: 8,
        }],
      }

      if (hasPercentage) {
        charts.vaccinationPercentage = {
          labels: topStates.map(item => item.state || item.State),
          datasets: [{
            label: 'Vaccination Percentage',
            data: topStates.map(item => item.percentage || item.Percentage),
            backgroundColor: 'rgba(99, 102, 241, 0.8)',
            borderColor: 'rgb(99, 102, 241)',
            borderWidth: 2,
            borderRadius: 8,
          }],
        }
      }
    }

    if ((hasQuarter || hasMonth) && hasRevenue) {
      const timeData = data.reduce((acc: any, item) => {
        const timeKey = item.quarter || item.Quarter || item.month || item.Month || item.month_year
        const revenue = item.total_revenue || item.revenue || item.Revenue
        if (!acc[timeKey]) acc[timeKey] = 0
        acc[timeKey] += Number(revenue) || 0
        return acc
      }, {})

      charts.timeSeries = {
        labels: Object.keys(timeData),
        datasets: [{
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
        }],
      }
    }

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
        datasets: [{
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
        }],
      }
    }

    if (hasCategory && hasRevenue) {
      const categoryData = data.reduce((acc: any, item) => {
        const category = item.bike_category || item.category || item.Category
        const revenue = item.total_revenue || item.revenue || item.Revenue
        if (!acc[category]) acc[category] = 0
        acc[category] += Number(revenue) || 0
        return acc
      }, {})

      const colors = [
        'rgba(239, 68, 68, 0.8)', 'rgba(245, 158, 11, 0.8)',
        'rgba(34, 197, 94, 0.8)', 'rgba(99, 102, 241, 0.8)',
        'rgba(236, 72, 153, 0.8)', 'rgba(14, 165, 233, 0.8)',
      ]

      charts.categoryBreakdown = {
        labels: Object.keys(categoryData),
        datasets: [{
          label: 'Revenue by Category',
          data: Object.values(categoryData),
          backgroundColor: colors,
          borderColor: colors.map(c => c.replace('0.8', '1')),
          borderWidth: 2,
          hoverOffset: 4,
        }],
      }

      charts.categoryBar = {
        labels: Object.keys(categoryData),
        datasets: [{
          label: 'Revenue ($)',
          data: Object.values(categoryData),
          backgroundColor: colors,
          borderColor: colors.map(c => c.replace('0.8', '1')),
          borderWidth: 2,
          borderRadius: 8,
          borderSkipped: false,
        }],
      }
    }

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
        datasets: [{
          label: 'Revenue Trend',
          data: revenues,
          borderColor: 'rgb(34, 197, 94)',
          backgroundColor: 'rgba(34, 197, 94, 0.1)',
          borderWidth: 3,
          fill: true,
          tension: 0.4,
        }],
      }
    }

    // Generic fallback
    if (Object.keys(charts).length === 0 && data.length > 0) {
      const firstRow = data[0]
      const allKeys = Object.keys(firstRow)

      const numericCols: string[] = []
      const stringCols: string[] = []

      allKeys.forEach(key => {
        const sampleValue = firstRow[key]
        if (typeof sampleValue === 'number' || !isNaN(Number(sampleValue))) {
          numericCols.push(key)
        } else {
          stringCols.push(key)
        }
      })

      if (stringCols.length > 0 && numericCols.length > 0) {
        const categoricalKey = stringCols[0]
        const numericKey = numericCols[0]

        const aggregatedData = data.reduce((acc: any, item) => {
          const category = String(item[categoricalKey] || 'Unknown')
          const value = Number(item[numericKey]) || 0
          if (!acc[category]) acc[category] = 0
          acc[category] += value
          return acc
        }, {})

        const sortedEntries = Object.entries(aggregatedData)
          .sort((a: any, b: any) => b[1] - a[1])
          .slice(0, 10)

        const sortedData = Object.fromEntries(sortedEntries)

        const colors = [
          'rgba(99, 102, 241, 0.8)', 'rgba(34, 197, 94, 0.8)',
          'rgba(239, 68, 68, 0.8)', 'rgba(245, 158, 11, 0.8)',
          'rgba(236, 72, 153, 0.8)', 'rgba(14, 165, 233, 0.8)',
          'rgba(168, 85, 247, 0.8)', 'rgba(251, 146, 60, 0.8)',
          'rgba(20, 184, 166, 0.8)', 'rgba(244, 63, 94, 0.8)',
        ]

        charts.genericBar = {
          labels: Object.keys(sortedData),
          datasets: [{
            label: numericKey.replace(/_/g, ' ').replace(/\b\w/g, (l: string) => l.toUpperCase()),
            data: Object.values(sortedData),
            backgroundColor: colors,
            borderColor: colors.map(c => c.replace('0.8', '1')),
            borderWidth: 2,
            borderRadius: 8,
          }],
        }

        if (Object.keys(sortedData).length >= 2 && Object.keys(sortedData).length <= 10) {
          charts.genericPie = {
            labels: Object.keys(sortedData),
            datasets: [{
              label: numericKey.replace(/_/g, ' ').replace(/\b\w/g, (l: string) => l.toUpperCase()),
              data: Object.values(sortedData),
              backgroundColor: colors,
              borderColor: colors.map(c => c.replace('0.8', '1')),
              borderWidth: 2,
              hoverOffset: 4,
            }],
          }
        }

        if (numericCols.length > 1) {
          const topCategories = Object.keys(sortedData).slice(0, 5)
          const datasets = numericCols.slice(0, 3).map((col, idx) => {
            const categoryValues = topCategories.map(cat => {
              const items = data.filter(item => String(item[categoricalKey]) === cat)
              return items.reduce((sum, item) => sum + (Number(item[col]) || 0), 0)
            })

            return {
              label: col.replace(/_/g, ' ').replace(/\b\w/g, (l: string) => l.toUpperCase()),
              data: categoryValues,
              backgroundColor: colors[idx],
              borderColor: colors[idx].replace('0.8', '1'),
              borderWidth: 2,
              borderRadius: 8,
            }
          })

          charts.genericGrouped = {
            labels: topCategories,
            datasets: datasets,
          }
        }
      } else if (numericCols.length > 0) {
        const numericKey = numericCols[0]
        const topData = data.slice(0, 10)

        charts.genericBar = {
          labels: topData.map((_, idx) => `Record ${idx + 1}`),
          datasets: [{
            label: numericKey.replace(/_/g, ' ').replace(/\b\w/g, (l: string) => l.toUpperCase()),
            data: topData.map(item => Number(item[numericKey]) || 0),
            backgroundColor: 'rgba(99, 102, 241, 0.8)',
            borderColor: 'rgb(99, 102, 241)',
            borderWidth: 2,
            borderRadius: 8,
          }],
        }
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
      transition: { duration: 0.6, staggerChildren: 0.1 },
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
      <InsightsHeader insights={insights} />

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {charts.genderBreakdown && (
          <ChartCard icon={PieChart} iconClassName="text-purple-600" title="Births by Gender">
            <Doughnut data={charts.genderBreakdown} options={doughnutOptions} />
          </ChartCard>
        )}

        {charts.topNames && (
          <ChartCard icon={BarChart3} iconClassName="text-green-600" title="Top Baby Names">
            <Bar data={charts.topNames} options={chartOptions} />
          </ChartCard>
        )}

        {charts.platformBreakdown && (
          <ChartCard icon={PieChart} iconClassName="text-blue-600" title="Sales by Platform">
            <Doughnut data={charts.platformBreakdown} options={doughnutOptions} />
          </ChartCard>
        )}

        {charts.topGames && (
          <ChartCard icon={BarChart3} iconClassName="text-indigo-600" title="Top Video Games">
            <Bar data={charts.topGames} options={chartOptions} />
          </ChartCard>
        )}

        {charts.channelActivity && (
          <ChartCard icon={PieChart} iconClassName="text-orange-600" title="Channel Activity">
            <Doughnut data={charts.channelActivity} options={doughnutOptions} />
          </ChartCard>
        )}

        {charts.userEngagement && (
          <ChartCard icon={BarChart3} iconClassName="text-cyan-600" title="User Engagement">
            <Bar data={charts.userEngagement} options={chartOptions} />
          </ChartCard>
        )}

        {charts.topStates && (
          <ChartCard icon={BarChart3} iconClassName="text-emerald-600" title="Vaccination by State">
            <Bar data={charts.topStates} options={chartOptions} />
          </ChartCard>
        )}

        {charts.vaccinationPercentage && (
          <ChartCard icon={BarChart3} iconClassName="text-violet-600" title="Vaccination Percentage">
            <Bar data={charts.vaccinationPercentage} options={chartOptions} />
          </ChartCard>
        )}

        {charts.timeSeries && (
          <ChartCard icon={TrendingUp} iconClassName="text-blue-600" title="Revenue Trend">
            <Line data={charts.timeSeries} options={chartOptions} />
          </ChartCard>
        )}

        {charts.categoryBar && (
          <ChartCard icon={BarChart3} iconClassName="text-green-600" title="Category Performance">
            <Bar data={charts.categoryBar} options={chartOptions} />
          </ChartCard>
        )}

        {charts.categoryBreakdown && (
          <ChartCard icon={PieChart} iconClassName="text-purple-600" title="Market Share">
            <Doughnut data={charts.categoryBreakdown} options={doughnutOptions} />
          </ChartCard>
        )}

        {charts.ordersTrend && (
          <ChartCard icon={TrendingUp} iconClassName="text-green-600" title="Orders Trend">
            <Line data={charts.ordersTrend} options={chartOptions} />
          </ChartCard>
        )}

        {charts.growth && (
          <ChartCard icon={TrendingUp} iconClassName="text-indigo-600" title="Growth Analysis">
            <Line data={charts.growth} options={growthOptions} />
          </ChartCard>
        )}

        {charts.genericBar && (
          <ChartCard icon={BarChart3} iconClassName="text-indigo-600" title="Data Overview">
            <Bar data={charts.genericBar} options={chartOptions} />
          </ChartCard>
        )}

        {charts.genericPie && (
          <ChartCard icon={PieChart} iconClassName="text-purple-600" title="Distribution">
            <Doughnut data={charts.genericPie} options={doughnutOptions} />
          </ChartCard>
        )}

        {charts.genericGrouped && (
          <ChartCard icon={BarChart3} iconClassName="text-green-600" title="Comparative Analysis">
            <Bar data={charts.genericGrouped} options={chartOptions} />
          </ChartCard>
        )}
      </div>

      {data && data.length > 0 && (
        <MetricsCards data={data} />
      )}
    </motion.div>
  )
}
