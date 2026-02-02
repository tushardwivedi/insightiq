"use client";

import { useState, useEffect, useCallback } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { motion } from "framer-motion";
import { useRouter } from "next/navigation";
import {
  Database,
  Upload,
  FileText,
  BarChart3,
  TrendingUp,
  Clock,
  ChevronRight,
  Zap,
  History,
  Settings,
  Users,
  FileBarChart,
  ArrowRight
} from "lucide-react";
import { apiClient } from "@/lib/api";

export default function DashboardPage() {
  const router = useRouter()

  const [quickStats, setQuickStats] = useState([
    { label: "Total Queries", value: "-", change: "", icon: "⚡", color: "from-blue-500 to-cyan-500" },
    { label: "Data Sources", value: "-", change: "", icon: "📊", color: "from-purple-500 to-pink-500" },
    { label: "Successful Queries", value: "-", change: "", icon: "💡", color: "from-orange-500 to-red-500" },
  ])
  const [recentQueries, setRecentQueries] = useState<{ query: string; time: string; type: string }[]>([])
  const [dashboardLoading, setDashboardLoading] = useState(true)

  const formatRelativeTime = (dateString: string): string => {
    const date = new Date(dateString)
    const now = new Date()
    const diffMs = now.getTime() - date.getTime()
    const diffMins = Math.floor(diffMs / (1000 * 60))
    if (diffMins < 1) return "Just now"
    if (diffMins < 60) return `${diffMins} min ago`
    const diffHours = Math.floor(diffMins / 60)
    if (diffHours < 24) return `${diffHours} hour${diffHours > 1 ? "s" : ""} ago`
    const diffDays = Math.floor(diffHours / 24)
    return `${diffDays} day${diffDays > 1 ? "s" : ""} ago`
  }

  const loadDashboardData = useCallback(async () => {
    try {
      setDashboardLoading(true)
      const [statsResult, connectorsResult, historyResult] = await Promise.allSettled([
        apiClient.getQueryHistoryStats(),
        apiClient.getConnectors(),
        apiClient.getQueryHistory(3, 0),
      ])

      const stats = statsResult.status === "fulfilled" ? (statsResult.value as any)?.data || statsResult.value || {} : {}
      const connectors = connectorsResult.status === "fulfilled" ? (connectorsResult.value as any)?.data || [] : []
      const history = historyResult.status === "fulfilled" ? (historyResult.value as any)?.data || [] : []

      const totalQueries = stats.total_queries || 0
      const successfulQueries = stats.successful_queries || 0
      const connectorCount = Array.isArray(connectors) ? connectors.length : 0

      setQuickStats([
        {
          label: "Total Queries",
          value: String(totalQueries),
          change: totalQueries > 0 ? `${successfulQueries} successful` : "",
          icon: "⚡",
          color: "from-blue-500 to-cyan-500",
        },
        {
          label: "Data Sources",
          value: String(connectorCount),
          change: "",
          icon: "📊",
          color: "from-purple-500 to-pink-500",
        },
        {
          label: "Successful Queries",
          value: String(successfulQueries),
          change: stats.failed_queries ? `${stats.failed_queries} failed` : "",
          icon: "💡",
          color: "from-orange-500 to-red-500",
        },
      ])

      setRecentQueries(
        (Array.isArray(history) ? history : []).map((item: any) => ({
          query: item.query_text || "Untitled query",
          time: item.created_at ? formatRelativeTime(item.created_at) : "",
          type: (item.query_type || "text").charAt(0).toUpperCase() + (item.query_type || "text").slice(1),
        }))
      )
    } catch (error) {
      console.error("Failed to load dashboard data:", error)
    } finally {
      setDashboardLoading(false)
    }
  }, [])

  useEffect(() => {
    loadDashboardData()
  }, [loadDashboardData])

  const quickActions = [
    { icon: Database, label: "Connect Data", desc: "Add a new data source", color: "from-blue-500 to-cyan-500", href: "/connectors" },
    { icon: Upload, label: "Upload File", desc: "Import CSV or Excel", color: "from-purple-500 to-pink-500", action: "upload" },
    { icon: FileText, label: "Saved Queries", desc: "Your saved queries", color: "from-orange-500 to-red-500", href: "/queries" },
    { icon: History, label: "Query History", desc: "View past queries", color: "from-green-500 to-emerald-500", href: "/history" },
  ]

  const analyticsCards = [
    { icon: BarChart3, label: "Analytics Dashboard", desc: "View all your analytics", href: "/analytics", color: "text-blue-500" },
    { icon: FileBarChart, label: "Reports", desc: "Generate and view reports", href: "/reports", color: "text-purple-500" },
    { icon: Users, label: "Team", desc: "Collaborate with your team", href: "/team", color: "text-green-500" },
    { icon: Settings, label: "Settings", desc: "Manage your account", href: "/settings", color: "text-gray-500" },
  ]

  return (
    <div className="flex-1 overflow-auto">
      <div className="container mx-auto max-w-7xl space-y-8 p-6">
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          className="flex items-center justify-between"
        >
          <div>
            <h1 className="text-3xl font-bold">Dashboard</h1>
            <p className="text-muted-foreground mt-1">
              Welcome back! Here's an overview of your analytics activity.
            </p>
          </div>
          <Button
            onClick={() => router.push('/editor')}
            className="bg-gradient-to-r from-primary to-cyan-500 hover:from-primary/90 hover:to-cyan-600"
          >
            <Zap className="w-4 h-4 mr-2" />
            New Query
            <ArrowRight className="w-4 h-4 ml-2" />
          </Button>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1 }}
          className="grid grid-cols-1 md:grid-cols-3 gap-4"
        >
          {quickStats.map((stat, i) => (
            <Card key={stat.label} className="hover:border-primary/50 transition-colors cursor-pointer group">
              <CardContent className="p-4">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm text-muted-foreground">{stat.label}</p>
                    <p className="text-3xl font-bold mt-1">{stat.value}</p>
                    {stat.change && (
                      <p className="text-xs text-green-500 mt-1 flex items-center gap-1">
                        <TrendingUp className="w-3 h-3" />
                        {stat.change}
                      </p>
                    )}
                  </div>
                  <div className={`w-14 h-14 rounded-xl bg-gradient-to-br ${stat.color} flex items-center justify-center text-2xl group-hover:scale-110 transition-transform shadow-lg`}>
                    {stat.icon}
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </motion.div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2 }}
            className="lg:col-span-2"
          >
            <Card className="h-full">
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Zap className="w-5 h-5 text-primary" />
                  Quick Actions
                </CardTitle>
                <CardDescription>Get started with common tasks</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-2 gap-4">
                  {quickActions.map((action, i) => (
                    <motion.button
                      key={action.label}
                      initial={{ opacity: 0, scale: 0.95 }}
                      animate={{ opacity: 1, scale: 1 }}
                      transition={{ delay: 0.3 + i * 0.1 }}
                      onClick={() => action.href ? router.push(action.href) : null}
                      className="p-4 rounded-xl bg-muted/50 hover:bg-muted transition-all group text-left border border-transparent hover:border-primary/20"
                    >
                      <div className={`w-12 h-12 rounded-lg bg-gradient-to-br ${action.color} flex items-center justify-center mb-3 group-hover:scale-110 transition-transform`}>
                        <action.icon className="w-6 h-6 text-white" />
                      </div>
                      <p className="font-medium">{action.label}</p>
                      <p className="text-sm text-muted-foreground">{action.desc}</p>
                    </motion.button>
                  ))}
                </div>
              </CardContent>
            </Card>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.3 }}
          >
            <Card className="h-full">
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Clock className="w-5 h-5 text-primary" />
                  Recent Queries
                </CardTitle>
                <CardDescription>Your latest query activity</CardDescription>
              </CardHeader>
              <CardContent>
                {dashboardLoading ? (
                  <div className="space-y-3">
                    {[1, 2, 3].map((i) => (
                      <div key={i} className="h-12 bg-muted/30 rounded-lg animate-pulse" />
                    ))}
                  </div>
                ) : recentQueries.length === 0 ? (
                  <div className="text-center py-6 text-muted-foreground text-sm">
                    No recent queries yet. Start by asking a question!
                  </div>
                ) : (
                  <div className="space-y-3">
                    {recentQueries.map((q, i) => (
                      <div
                        key={i}
                        className="flex items-center justify-between p-3 rounded-lg bg-muted/30 hover:bg-muted/50 transition-colors cursor-pointer group"
                      >
                        <div className="flex items-center gap-3 min-w-0">
                          <span className="text-xs px-2 py-1 rounded-full bg-primary/10 text-primary font-medium">
                            {q.type}
                          </span>
                          <span className="text-sm truncate">{q.query}</span>
                        </div>
                        <span className="text-xs text-muted-foreground ml-2">{q.time}</span>
                      </div>
                    ))}
                  </div>
                )}
                <Button variant="ghost" className="w-full mt-4" onClick={() => router.push('/history')}>
                  View all queries
                  <ChevronRight className="w-4 h-4 ml-2" />
                </Button>
              </CardContent>
            </Card>
          </motion.div>
        </div>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.4 }}
        >
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <BarChart3 className="w-5 h-5 text-primary" />
                Explore More
              </CardTitle>
              <CardDescription>Navigate to other sections of the application</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                {analyticsCards.map((card, i) => (
                  <motion.button
                    key={card.label}
                    initial={{ opacity: 0, y: 10 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: 0.5 + i * 0.1 }}
                    onClick={() => router.push(card.href)}
                    className="p-4 rounded-xl border border-border/50 hover:border-primary/50 hover:bg-muted/30 transition-all group text-center"
                  >
                    <card.icon className={`w-8 h-8 mx-auto mb-3 ${card.color} group-hover:scale-110 transition-transform`} />
                    <p className="font-medium">{card.label}</p>
                    <p className="text-sm text-muted-foreground">{card.desc}</p>
                  </motion.button>
                ))}
              </div>
            </CardContent>
          </Card>
        </motion.div>
      </div>
    </div>
  );
}
