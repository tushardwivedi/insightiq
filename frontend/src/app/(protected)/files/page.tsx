"use client";

import { useState, useEffect, useCallback } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { apiClient } from "@/lib/api";
import { UploadedFile } from "@/types";
import { motion } from "framer-motion";
import FileUploadModal from "@/components/FileUploadModal";
import {
  FolderOpen,
  Upload,
  Trash2,
  File,
  FileText,
  Table,
  Download,
  Database,
  Clock,
  BarChart3,
  CheckCircle,
  AlertCircle,
  MoreVertical,
  Eye,
  Sparkles
} from "lucide-react";

export default function FilesPage() {
  const [files, setFiles] = useState<UploadedFile[]>([]);
  const [loading, setLoading] = useState(true);
  const [showUploadModal, setShowUploadModal] = useState(false);
  const [refreshTrigger, setRefreshTrigger] = useState(0);

  const loadFiles = useCallback(async () => {
    setLoading(true);
    try {
      const result = await apiClient.getUploadedFiles();
      setFiles(result);
    } catch (error) {
      console.error("Failed to load files:", error);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadFiles();
  }, [loadFiles, refreshTrigger]);

  const handleDelete = async (id: string) => {
    if (!confirm("Are you sure you want to delete this file?")) return;
    try {
      await apiClient.deleteUploadedFile(id);
      setFiles(prev => prev.filter(f => f.id !== id));
    } catch (error) {
      console.error("Failed to delete:", error);
    }
  };

  const handleUploadComplete = (file: UploadedFile) => {
    setRefreshTrigger(prev => prev + 1);
  };

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return "0 Bytes";
    const k = 1024;
    const sizes = ["Bytes", "KB", "MB", "GB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
  };

  const formatDate = (date?: string) => {
    if (!date) return "";
    return new Date(date).toLocaleDateString(undefined, {
      month: "short",
      day: "numeric",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const getFileIcon = (filename: string) => {
    const ext = filename.split(".").pop()?.toLowerCase();
    if (["xlsx", "xls"].includes(ext || "")) return Table;
    return FileText;
  };

  const getQualityColor = (score?: number) => {
    if (!score) return "text-gray-500";
    if (score >= 0.8) return "text-green-500";
    if (score >= 0.6) return "text-yellow-500";
    return "text-red-500";
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
              <FolderOpen className="w-8 h-8 text-primary" />
              Uploaded Files
            </h1>
            <p className="text-muted-foreground mt-1">
              Manage your uploaded data files
            </p>
          </div>
          <Button className="bg-gradient-to-r from-primary to-cyan-500" onClick={() => setShowUploadModal(true)}>
            <Upload className="w-4 h-4 mr-2" />
            Upload File
          </Button>
        </motion.div>

        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center gap-3">
                <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-blue-500/20 to-cyan-500/20 flex items-center justify-center">
                  <FolderOpen className="w-6 h-6 text-primary" />
                </div>
                <div>
                  <p className="text-2xl font-bold">{files.length}</p>
                  <p className="text-sm text-muted-foreground">Total Files</p>
                </div>
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center gap-3">
                <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-green-500/20 to-emerald-500/20 flex items-center justify-center">
                  <Database className="w-6 h-6 text-green-500" />
                </div>
                <div>
                  <p className="text-2xl font-bold">
                    {files.filter(f => f.status === "completed").length}
                  </p>
                  <p className="text-sm text-muted-foreground">Processed</p>
                </div>
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center gap-3">
                <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-purple-500/20 to-pink-500/20 flex items-center justify-center">
                  <BarChart3 className="w-6 h-6 text-purple-500" />
                </div>
                <div>
                  <p className="text-2xl font-bold">
                    {files.reduce((sum, f) => sum + (f.row_count || 0), 0).toLocaleString()}
                  </p>
                  <p className="text-sm text-muted-foreground">Total Rows</p>
                </div>
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center gap-3">
                <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-orange-500/20 to-red-500/20 flex items-center justify-center">
                  <Sparkles className="w-6 h-6 text-orange-500" />
                </div>
                <div>
                  <p className="text-2xl font-bold">
                    {files.filter(f => (f.data_quality_score || 0) >= 0.8).length}
                  </p>
                  <p className="text-sm text-muted-foreground">High Quality</p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>All Files</CardTitle>
            <CardDescription>Your uploaded and processed data files</CardDescription>
          </CardHeader>
          <CardContent>
            {loading ? (
              <div className="space-y-4">
                {[1, 2, 3].map((i) => (
                  <div key={i} className="h-24 bg-muted/30 rounded-lg animate-pulse" />
                ))}
              </div>
            ) : files.length === 0 ? (
              <div className="text-center py-12">
                <FolderOpen className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                <p className="text-muted-foreground mb-4">No files uploaded yet</p>
                <Button onClick={() => setShowUploadModal(true)}>
                  <Upload className="w-4 h-4 mr-2" />
                  Upload Your First File
                </Button>
              </div>
            ) : (
              <div className="space-y-3">
                {files.map((file, index) => {
                  const FileIcon = getFileIcon(file.original_filename);
                  return (
                    <motion.div
                      key={file.id}
                      initial={{ opacity: 0, y: 10 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ delay: index * 0.03 }}
                      className="p-4 rounded-lg border hover:bg-muted/30 transition-colors"
                    >
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-4">
                          <div className="w-12 h-12 rounded-xl bg-muted flex items-center justify-center">
                            <FileIcon className="w-6 h-6 text-primary" />
                          </div>
                          <div>
                            <div className="flex items-center gap-2">
                              <h3 className="font-semibold">{file.original_filename}</h3>
                              <Badge variant="secondary" className="text-xs">
                                {file.file_type?.toUpperCase()}
                              </Badge>
                              {file.status === "completed" ? (
                                <CheckCircle className="w-4 h-4 text-green-500" />
                              ) : file.status === "processing" ? (
                                <Clock className="w-4 h-4 text-yellow-500 animate-spin" />
                              ) : (
                                <AlertCircle className="w-4 h-4 text-red-500" />
                              )}
                            </div>
                            <div className="flex items-center gap-4 mt-1 text-sm text-muted-foreground">
                              <span className="flex items-center gap-1">
                                {formatFileSize(file.file_size)}
                              </span>
                              <span className="flex items-center gap-1">
                                <Database className="w-3 h-3" />
                                {file.table_name}
                              </span>
                              <span className="flex items-center gap-1">
                                <BarChart3 className="w-3 h-3" />
                                {(file.row_count || 0).toLocaleString()} rows
                              </span>
                              <span>{formatDate(file.created_at)}</span>
                            </div>
                          </div>
                        </div>
                        <div className="flex items-center gap-2">
                          {file.data_quality_score !== undefined && (
                            <Badge
                              variant="outline"
                              className={`${getQualityColor(file.data_quality_score)} border-current`}
                            >
                              Quality: {(file.data_quality_score * 100).toFixed(0)}%
                            </Badge>
                          )}
                          <Button variant="ghost" size="icon" title="View">
                            <Eye className="w-4 h-4" />
                          </Button>
                          <Button variant="ghost" size="icon" title="Query">
                            <Database className="w-4 h-4" />
                          </Button>
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => handleDelete(file.id)}
                            title="Delete"
                          >
                            <Trash2 className="w-4 h-4 text-destructive" />
                          </Button>
                        </div>
                      </div>
                    </motion.div>
                  );
                })}
              </div>
            )}
          </CardContent>
        </Card>

        <FileUploadModal
          isOpen={showUploadModal}
          onClose={() => setShowUploadModal(false)}
          onUploadComplete={handleUploadComplete}
        />
      </div>
    </div>
  );
}
