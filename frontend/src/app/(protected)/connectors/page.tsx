"use client";

import { useState, useEffect, useCallback } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { apiClient } from "@/lib/api";
import { DataConnector } from "@/types";
import { motion } from "framer-motion";
import {
  Database,
  Plus,
  Trash2,
  RefreshCw,
  CheckCircle,
  XCircle,
  AlertCircle,
  Settings,
  ChevronRight,
  ExternalLink,
  TestTube
} from "lucide-react";

export default function ConnectorsPage() {
  const [connectors, setConnectors] = useState<DataConnector[]>([]);
  const [loading, setLoading] = useState(true);
  const [testing, setTesting] = useState<string | null>(null);

  const loadConnectors = useCallback(async () => {
    try {
      const response = await apiClient.getConnectors();
      setConnectors(response.data || []);
    } catch (error) {
      console.error("Failed to load connectors:", error);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadConnectors();
  }, [loadConnectors]);

  const handleTest = async (connector: DataConnector) => {
    setTesting(connector.id);
    try {
      const result = await apiClient.testConnector(connector.id);
      setConnectors(prev =>
        prev.map(c =>
          c.id === connector.id
            ? { ...c, status: result.success ? "connected" : "error", last_tested: new Date().toISOString() }
            : c
        )
      );
    } catch (error) {
      setConnectors(prev =>
        prev.map(c =>
          c.id === connector.id ? { ...c, status: "error" } : c
        )
      );
    } finally {
      setTesting(null);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm("Are you sure you want to delete this connector?")) return;
    try {
      await apiClient.deleteConnector(id);
      setConnectors(prev => prev.filter(c => c.id !== id));
    } catch (error) {
      console.error("Failed to delete:", error);
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case "connected": return "text-green-500";
      case "testing": return "text-yellow-500";
      case "error": return "text-red-500";
      default: return "text-gray-500";
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case "connected": return <CheckCircle className="w-4 h-4 text-green-500" />;
      case "testing": return <RefreshCw className="w-4 h-4 text-yellow-500 animate-spin" />;
      case "error": return <XCircle className="w-4 h-4 text-red-500" />;
      default: return <AlertCircle className="w-4 h-4 text-gray-500" />;
    }
  };

  const getConnectorIcon = (type: string) => {
    const icons: Record<string, string> = {
      postgres: "🐘",
      mysql: "🐬",
      superset: "📊",
      mongodb: "🍃",
      api: "🌐",
    };
    return icons[type] || "📁";
  };

  const formatDate = (date?: string) => {
    if (!date) return "Never";
    return new Date(date).toLocaleDateString(undefined, {
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  return (
    <div className="flex-1 overflow-auto">
      <div className="container mx-auto max-w-6xl space-y-6 p-6">
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          className="flex items-center justify-between"
        >
          <div>
            <h1 className="text-3xl font-bold flex items-center gap-3">
              <Database className="w-8 h-8 text-primary" />
              Data Sources
            </h1>
            <p className="text-muted-foreground mt-1">
              Manage your connected databases and data sources
            </p>
          </div>
          <Button className="bg-gradient-to-r from-primary to-cyan-500">
            <Plus className="w-4 h-4 mr-2" />
            Add Connector
          </Button>
        </motion.div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center gap-3">
                <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-blue-500/20 to-cyan-500/20 flex items-center justify-center">
                  <Database className="w-6 h-6 text-primary" />
                </div>
                <div>
                  <p className="text-2xl font-bold">{connectors.length}</p>
                  <p className="text-sm text-muted-foreground">Total Connectors</p>
                </div>
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center gap-3">
                <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-green-500/20 to-emerald-500/20 flex items-center justify-center">
                  <CheckCircle className="w-6 h-6 text-green-500" />
                </div>
                <div>
                  <p className="text-2xl font-bold">
                    {connectors.filter(c => c.status === "connected").length}
                  </p>
                  <p className="text-sm text-muted-foreground">Connected</p>
                </div>
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center gap-3">
                <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-red-500/20 to-orange-500/20 flex items-center justify-center">
                  <AlertCircle className="w-6 h-6 text-red-500" />
                </div>
                <div>
                  <p className="text-2xl font-bold">
                    {connectors.filter(c => c.status === "error").length}
                  </p>
                  <p className="text-sm text-muted-foreground">Errors</p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>Connected Sources</CardTitle>
            <CardDescription>Your active data source connections</CardDescription>
          </CardHeader>
          <CardContent>
            {loading ? (
              <div className="space-y-4">
                {[1, 2, 3].map((i) => (
                  <div key={i} className="h-24 bg-muted/30 rounded-lg animate-pulse" />
                ))}
              </div>
            ) : connectors.length === 0 ? (
              <div className="text-center py-12">
                <Database className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                <p className="text-muted-foreground mb-4">No data sources connected</p>
                <Button>
                  <Plus className="w-4 h-4 mr-2" />
                  Add Your First Connector
                </Button>
              </div>
            ) : (
              <div className="space-y-3">
                {connectors.map((connector, index) => (
                  <motion.div
                    key={connector.id}
                    initial={{ opacity: 0, y: 10 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: index * 0.05 }}
                    className="p-4 rounded-lg border hover:bg-muted/30 transition-colors"
                  >
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-4">
                        <div className="w-12 h-12 rounded-xl bg-muted flex items-center justify-center text-2xl">
                          {getConnectorIcon(connector.type)}
                        </div>
                        <div>
                          <div className="flex items-center gap-2">
                            <h3 className="font-semibold">{connector.name}</h3>
                            {getStatusIcon(connector.status)}
                          </div>
                          <p className="text-sm text-muted-foreground capitalize">
                            {connector.type}
                          </p>
                          {connector.last_tested && (
                            <p className="text-xs text-muted-foreground mt-1">
                              Last tested: {formatDate(connector.last_tested)}
                            </p>
                          )}
                        </div>
                      </div>
                      <div className="flex items-center gap-2">
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => handleTest(connector)}
                          disabled={testing === connector.id}
                        >
                          <TestTube className={`w-4 h-4 mr-1 ${testing === connector.id ? "animate-spin" : ""}`} />
                          Test
                        </Button>
                        <Button variant="outline" size="sm">
                          <Settings className="w-4 h-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => handleDelete(connector.id)}
                        >
                          <Trash2 className="w-4 h-4 text-destructive" />
                        </Button>
                      </div>
                    </div>
                  </motion.div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
