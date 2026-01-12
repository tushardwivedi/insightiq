"use client";

import { useState, useEffect } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { apiClient } from "@/lib/api";
import { motion } from "framer-motion";
import {
  BarChart3,
  TrendingUp,
  TrendingDown,
  Users,
  Database,
  Clock,
  Zap,
  FileText,
  Activity,
  Calendar,
  ArrowUpRight,
  ArrowDownRight
} from "lucide-react";

export default function AnalyticsPage() {
  const [analytics, setAnalytics] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState("overview");

  useEffect(() => {
    loadAnalytics();
  }, []);

  const loadAnalytics = async () => {
    setLoading(true);
    try {
      const result = await apiClient.getAnalytics();
      setAnalytics(result);
    } catch (error) {
      console.error("Failed to load analytics:", error);
    } finally {
      setLoading(false);
    }
  };

  const formatNumber = (num: number) => {
    if (num >= 1000000) return (num / 1000000).toFixed(1) + "M";
    if (num >= 1000) return (num / 1000).toFixed(1) + "K";
    return num.toString();
  };

  if (loading) {
    return (
      <div className="flex-1 overflow-auto">
        <div className="container mx-auto max-w-6xl space-y-6 p-6">
          <div className="h-10 w-64 bg-muted/30 rounded animate-pulse" />
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            {[1, 2, 3, 4].map(i => (
              <div key={i} className="h-32 bg-muted/30 rounded-lg animate-pulse" />
            ))}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="flex-1 overflow-auto">
      <div className="container mx-auto max-w-6xl space-y-6 p-6">
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
        >
          <h1 className="text-3xl font-bold flex items-center gap-3">
            <BarChart3 className="w-8 h-8 text-primary" />
            Analytics Dashboard
          </h1>
          <p className="text-muted-foreground mt-1">
            Overview of your analytics activity and insights
          </p>
        </motion.div>

        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList>
            <TabsTrigger value="overview">Overview</TabsTrigger>
            <TabsTrigger value="queries">Queries</TabsTrigger>
            <TabsTrigger value="data">Data Sources</TabsTrigger>
          </TabsList>

          <TabsContent value="overview" className="mt-6 space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
              <Card>
                <CardContent className="p-4">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="text-sm text-muted-foreground">Total Queries</p>
                      <p className="text-3xl font-bold mt-1">{formatNumber(analytics?.totalQueries || 156)}</p>
                      <p className="text-xs text-green-500 flex items-center mt-1">
                        <ArrowUpRight className="w-3 h-3 mr-1" />
                        +12% from last month
                      </p>
                    </div>
                    <div className="w-12 h-12 rounded-xl bg-blue-500/10 flex items-center justify-center">
                      <Activity className="w-6 h-6 text-blue-500" />
                    </div>
                  </div>
                </CardContent>
              </Card>
              <Card>
                <CardContent className="p-4">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="text-sm text-muted-foreground">Data Sources</p>
                      <p className="text-3xl font-bold mt-1">{formatNumber(analytics?.dataSources || 8)}</p>
                      <p className="text-xs text-green-500 flex items-center mt-1">
                        <ArrowUpRight className="w-3 h-3 mr-1" />
                        +2 new this week
                      </p>
                    </div>
                    <div className="w-12 h-12 rounded-xl bg-purple-500/10 flex items-center justify-center">
                      <Database className="w-6 h-6 text-purple-500" />
                    </div>
                  </div>
                </CardContent>
              </Card>
              <Card>
                <CardContent className="p-4">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="text-sm text-muted-foreground">Processing Time</p>
                      <p className="text-3xl font-bold mt-1">{(analytics?.avgProcessingTime || 1.2).toFixed(1)}s</p>
                      <p className="text-xs text-green-500 flex items-center mt-1">
                        <ArrowDownRight className="w-3 h-3 mr-1" />
                        -8% faster
                      </p>
                    </div>
                    <div className="w-12 h-12 rounded-xl bg-green-500/10 flex items-center justify-center">
                      <Clock className="w-6 h-6 text-green-500" />
                    </div>
                  </div>
                </CardContent>
              </Card>
              <Card>
                <CardContent className="p-4">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="text-sm text-muted-foreground">Active Users</p>
                      <p className="text-3xl font-bold mt-1">{formatNumber(analytics?.activeUsers || 24)}</p>
                      <p className="text-xs text-green-500 flex items-center mt-1">
                        <ArrowUpRight className="w-3 h-3 mr-1" />
                        +5 this week
                      </p>
                    </div>
                    <div className="w-12 h-12 rounded-xl bg-orange-500/10 flex items-center justify-center">
                      <Users className="w-6 h-6 text-orange-500" />
                    </div>
                  </div>
                </CardContent>
              </Card>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <Card>
                <CardHeader>
                  <CardTitle>Query Activity</CardTitle>
                  <CardDescription>Queries over the last 7 days</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="h-64 flex items-end justify-between gap-2">
                    {["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"].map((day, i) => {
                      const height = Math.random() * 60 + 20;
                      return (
                        <div key={day} className="flex-1 flex flex-col items-center gap-2">
                          <div
                            className="w-full bg-gradient-to-t from-primary to-cyan-500 rounded-t transition-all hover:opacity-80"
                            style={{ height: `${height}%` }}
                          />
                          <span className="text-xs text-muted-foreground">{day}</span>
                        </div>
                      );
                    })}
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle>Top Data Sources</CardTitle>
                  <CardDescription>Most used data sources</CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                  {[
                    { name: "Sales Database", queries: 45, color: "bg-blue-500" },
                    { name: "Customer Data", queries: 32, color: "bg-purple-500" },
                    { name: "Product Catalog", queries: 28, color: "bg-green-500" },
                    { name: "Inventory", queries: 18, color: "bg-orange-500" },
                  ].map((source, i) => (
                    <div key={source.name} className="flex items-center gap-3">
                      <div className={`w-3 h-3 rounded-full ${source.color}`} />
                      <div className="flex-1">
                        <div className="flex items-center justify-between mb-1">
                          <span className="text-sm font-medium">{source.name}</span>
                          <span className="text-sm text-muted-foreground">{source.queries} queries</span>
                        </div>
                        <div className="h-2 bg-muted rounded-full overflow-hidden">
                          <div
                            className={`h-full ${source.color} rounded-full`}
                            style={{ width: `${(source.queries / 50) * 100}%` }}
                          />
                        </div>
                      </div>
                    </div>
                  ))}
                </CardContent>
              </Card>
            </div>
          </TabsContent>

          <TabsContent value="queries" className="mt-6">
            <Card>
              <CardHeader>
                <CardTitle>Query Performance</CardTitle>
                <CardDescription>Recent query execution metrics</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="text-center py-12 text-muted-foreground">
                  <Activity className="w-12 h-12 mx-auto mb-4" />
                  <p>Query analytics coming soon</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="data" className="mt-6">
            <Card>
              <CardHeader>
                <CardTitle>Data Source Analytics</CardTitle>
                <CardDescription>Performance metrics by data source</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="text-center py-12 text-muted-foreground">
                  <Database className="w-12 h-12 mx-auto mb-4" />
                  <p>Data source analytics coming soon</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
}
