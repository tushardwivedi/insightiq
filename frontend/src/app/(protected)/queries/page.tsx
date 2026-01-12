"use client";

import { useState, useEffect } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { apiClient } from "@/lib/api";
import { motion } from "framer-motion";
import { useRouter } from "next/navigation";
import {
  Bookmark,
  Plus,
  Trash2,
  Search,
  Play,
  Copy,
  Check,
  Clock,
  FileText,
  MoreVertical,
  Star,
  Folder
} from "lucide-react";

interface SavedQuery {
  id: string;
  name: string;
  sql: string;
  description?: string;
  created_at: string;
  updated_at: string;
}

export default function QueriesPage() {
  const router = useRouter();
  const [queries, setQueries] = useState<SavedQuery[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");
  const [copied, setCopied] = useState<string | null>(null);

  useEffect(() => {
    loadQueries();
  }, []);

  const loadQueries = async () => {
    setLoading(true);
    try {
      const result = await apiClient.getSavedQueries();
      setQueries(result);
    } catch (error) {
      console.error("Failed to load queries:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm("Are you sure you want to delete this saved query?")) return;
    try {
      await apiClient.deleteSavedQuery(id);
      setQueries(prev => prev.filter(q => q.id !== id));
    } catch (error) {
      console.error("Failed to delete:", error);
    }
  };

  const handleCopy = async (sql: string, id: string) => {
    await navigator.clipboard.writeText(sql);
    setCopied(id);
    setTimeout(() => setCopied(null), 2000);
  };

  const handleRun = (sql: string) => {
    router.push(`/editor?query=${encodeURIComponent(sql)}`);
  };

  const filteredQueries = queries.filter(q =>
    q.name?.toLowerCase().includes(search.toLowerCase()) ||
    q.sql?.toLowerCase().includes(search.toLowerCase()) ||
    q.description?.toLowerCase().includes(search.toLowerCase())
  );

  const formatDate = (date?: string) => {
    if (!date) return "";
    return new Date(date).toLocaleDateString(undefined, {
      month: "short",
      day: "numeric",
      year: "numeric",
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
              <Bookmark className="w-8 h-8 text-primary" />
              Saved Queries
            </h1>
            <p className="text-muted-foreground mt-1">
              Your collection of saved SQL queries
            </p>
          </div>
          <Button className="bg-gradient-to-r from-primary to-cyan-500">
            <Plus className="w-4 h-4 mr-2" />
            Save Query
          </Button>
        </motion.div>

        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle>All Saved Queries</CardTitle>
                <CardDescription>Quick access to your frequently used queries</CardDescription>
              </div>
              <div className="relative">
                <Search className="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground" />
                <Input
                  placeholder="Search queries..."
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                  className="pl-9 w-[250px]"
                />
              </div>
            </div>
          </CardHeader>
          <CardContent>
            {loading ? (
              <div className="space-y-4">
                {[1, 2, 3].map((i) => (
                  <div key={i} className="h-24 bg-muted/30 rounded-lg animate-pulse" />
                ))}
              </div>
            ) : filteredQueries.length === 0 ? (
              <div className="text-center py-12">
                <Bookmark className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                <p className="text-muted-foreground mb-4">No saved queries yet</p>
                <p className="text-sm text-muted-foreground mb-4">
                  Save your frequently used queries for quick access
                </p>
                <Button>
                  <Plus className="w-4 h-4 mr-2" />
                  Save Your First Query
                </Button>
              </div>
            ) : (
              <div className="space-y-3">
                {filteredQueries.map((query, index) => (
                  <motion.div
                    key={query.id}
                    initial={{ opacity: 0, y: 10 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: index * 0.03 }}
                    className="p-4 rounded-lg border hover:bg-muted/30 transition-colors"
                  >
                    <div className="flex items-start justify-between gap-4">
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2 mb-2">
                          <Star className="w-4 h-4 text-yellow-500" />
                          <h3 className="font-semibold">{query.name}</h3>
                          <Badge variant="secondary">SQL</Badge>
                          <span className="text-xs text-muted-foreground">
                            {formatDate(query.created_at)}
                          </span>
                        </div>
                        {query.description && (
                          <p className="text-sm text-muted-foreground mb-2">
                            {query.description}
                          </p>
                        )}
                        <pre className="text-sm font-mono bg-muted/50 p-2 rounded truncate">
                          {query.sql}
                        </pre>
                      </div>
                      <div className="flex items-center gap-1">
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => handleRun(query.sql)}
                          title="Run query"
                        >
                          <Play className="w-4 h-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => handleCopy(query.sql, query.id)}
                          title="Copy query"
                        >
                          {copied === query.id ? (
                            <Check className="w-4 h-4 text-green-500" />
                          ) : (
                            <Copy className="w-4 h-4" />
                          )}
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => handleDelete(query.id)}
                          title="Delete"
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
