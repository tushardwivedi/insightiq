"use client";

import { useState, useRef, useEffect, useCallback } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Sheet, SheetContent, SheetHeader, SheetTitle } from "@/components/ui/sheet";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { motion, AnimatePresence } from "framer-motion";
import {
  Send,
  Mic,
  MicOff,
  FileAudio,
  Sparkles,
  Database,
  Zap,
  Copy,
  Check,
  Clock,
  Trash2,
  Bot,
  User,
  Wand2,
  Table,
  Code
} from "lucide-react";
import { AnalyticsResponse, VoiceResponse } from "@/types";
import { apiClient } from "@/lib/api";

interface Message {
  id: string;
  type: "user" | "assistant";
  content: string;
  sql?: string;
  results?: AnalyticsResponse | VoiceResponse;
  timestamp: Date;
}

export default function EditorPage() {
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isRecording, setIsRecording] = useState(false);
  const [recordedBlob, setRecordedBlob] = useState<Blob | null>(null);
  const [showSqlEditor, setShowSqlEditor] = useState(false);
  const [sqlQuery, setSqlQuery] = useState("");
  const [copied, setCopied] = useState<string | null>(null);
  const [mediaRecorder, setMediaRecorder] = useState<MediaRecorder | null>(null);
  const [audioProcessing, setAudioProcessing] = useState(false);

  const messagesEndRef = useRef<HTMLDivElement>(null);
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  const scrollToBottom = useCallback(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, []);

  useEffect(() => {
    scrollToBottom();
  }, [messages, scrollToBottom]);

  const startRecording = async () => {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
      const recorder = new MediaRecorder(stream);
      const newChunks: Blob[] = [];

      recorder.ondataavailable = (e) => {
        if (e.data.size > 0) newChunks.push(e.data);
      };

      recorder.onstop = () => {
        const blob = new Blob(newChunks, { type: "audio/webm" });
        setRecordedBlob(blob);
        stream.getTracks().forEach((t) => t.stop());
      };

      recorder.start(1000);
      setMediaRecorder(recorder);
      setIsRecording(true);
    } catch (error) {
      console.error("Error starting recording:", error);
    }
  };

  const stopRecording = () => {
    if (mediaRecorder && isRecording) {
      mediaRecorder.stop();
      setIsRecording(false);
    }
  };

  const processAudio = async (audioFile: File) => {
    setAudioProcessing(true);
    try {
      const result = await apiClient.voiceQuery(audioFile);
      const transcript = (result as VoiceResponse).transcript;
      const response = (result as VoiceResponse).response;

      const userMessage: Message = {
        id: Date.now().toString(),
        type: "user",
        content: transcript,
        timestamp: new Date(),
      };

      const assistantMessage: Message = {
        id: (Date.now() + 1).toString(),
        type: "assistant",
        content: response.insights || "Analysis complete",
        results: response,
        timestamp: new Date(),
      };

      setMessages((prev) => [...prev, userMessage, assistantMessage]);
    } catch (error) {
      const errorMessage: Message = {
        id: (Date.now() + 1).toString(),
        type: "assistant",
        content: `Error processing audio: ${error instanceof Error ? error.message : "Unknown error"}`,
        timestamp: new Date(),
      };
      setMessages((prev) => [...prev, errorMessage]);
    } finally {
      setAudioProcessing(false);
      setRecordedBlob(null);
    }
  };

  const handleSend = async () => {
    if (!input.trim() || isSubmitting) return;

    const userMessage: Message = {
      id: Date.now().toString(),
      type: "user",
      content: input,
      timestamp: new Date(),
    };

    setMessages((prev) => [...prev, userMessage]);
    setInput("");
    setIsSubmitting(true);

    try {
      const result = await apiClient.textQuery({ query: input });

      const assistantMessage: Message = {
        id: (Date.now() + 1).toString(),
        type: "assistant",
        content: (result as AnalyticsResponse).insights || "Analysis complete",
        results: result as AnalyticsResponse,
        timestamp: new Date(),
      };

      setMessages((prev) => [...prev, assistantMessage]);
    } catch (error) {
      const errorMessage: Message = {
        id: (Date.now() + 1).toString(),
        type: "assistant",
        content: `Error: ${error instanceof Error ? error.message : "Failed to process query"}`,
        timestamp: new Date(),
      };
      setMessages((prev) => [...prev, errorMessage]);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleSqlSubmit = async () => {
    if (!sqlQuery.trim() || isSubmitting) return;

    const userMessage: Message = {
      id: Date.now().toString(),
      type: "user",
      content: "SQL Query executed",
      sql: sqlQuery,
      timestamp: new Date(),
    };

    setMessages((prev) => [...prev, userMessage]);
    setShowSqlEditor(false);
    setSqlQuery("");
    setIsSubmitting(true);

    try {
      const result = await apiClient.sqlQuery({ sql: sqlQuery });

      const assistantMessage: Message = {
        id: (Date.now() + 1).toString(),
        type: "assistant",
        content: (result as AnalyticsResponse).insights || "Analysis complete",
        results: result as AnalyticsResponse,
        timestamp: new Date(),
      };

      setMessages((prev) => [...prev, assistantMessage]);
    } catch (error) {
      const errorMessage: Message = {
        id: (Date.now() + 1).toString(),
        type: "assistant",
        content: `Error: ${error instanceof Error ? error.message : "Failed to execute query"}`,
        timestamp: new Date(),
      };
      setMessages((prev) => [...prev, errorMessage]);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  const copyToClipboard = async (text: string, id: string) => {
    await navigator.clipboard.writeText(text);
    setCopied(id);
    setTimeout(() => setCopied(null), 2000);
  };

  const clearChat = () => {
    setMessages([]);
  };

  const processRecordedAudio = () => {
    if (recordedBlob) {
      const file = new File([recordedBlob], "recording.webm", { type: "audio/webm" });
      processAudio(file);
    }
  };

  const suggestedQueries = [
    "Show me sales trends for last quarter",
    "What are the top performing products?",
    "Analyze customer demographics",
    "Compare revenue across regions",
  ];

  const formatTime = (date: Date) => {
    return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
  };

  return (
    <div className="flex flex-col h-full bg-background">
      {/* Chat Messages Area - Scrollable */}
      <div className="flex-1 overflow-y-auto">
        <div className="max-w-3xl mx-auto space-y-0">
            <AnimatePresence mode="wait">
              {messages.length === 0 && (
                <motion.div
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: -20 }}
                  className="px-4 pt-12 pb-8"
                >
                  <div className="text-center mb-8">
                    <div className="w-12 h-12 mx-auto mb-4 rounded-xl bg-primary/10 flex items-center justify-center">
                      <Bot className="w-6 h-6 text-primary" />
                    </div>
                    <h1 className="text-2xl font-semibold mb-2">How can I help you today?</h1>
                    <p className="text-muted-foreground text-sm">Ask questions about your data or write SQL queries</p>
                  </div>

                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 max-w-2xl mx-auto">
                    <button
                      onClick={() => setShowSqlEditor(true)}
                      className="flex items-center gap-3 p-3 rounded-xl border border-border bg-card hover:bg-accent/50 transition-colors text-left"
                    >
                      <div className="w-8 h-8 rounded-lg bg-primary/10 flex items-center justify-center flex-shrink-0">
                        <Database className="w-4 h-4 text-primary" />
                      </div>
                      <div className="flex-1 min-w-0">
                        <p className="font-medium text-sm">Write SQL Query</p>
                        <p className="text-xs text-muted-foreground truncate">Open SQL editor</p>
                      </div>
                    </button>

                    {suggestedQueries.map((query, i) => (
                      <button
                        key={i}
                        onClick={() => setInput(query)}
                        className="p-3 rounded-xl border border-border bg-card hover:bg-accent/50 transition-colors text-left text-sm"
                      >
                        {query}
                      </button>
                    ))}
                  </div>
                </motion.div>
              )}
            </AnimatePresence>

            <div className="space-y-0">
              <AnimatePresence>
                {messages.map((message, index) => (
                  <motion.div
                    key={message.id}
                    initial={{ opacity: 0, y: 10 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, y: -10 }}
                    transition={{ delay: index * 0.03 }}
                    className={`${message.type === "assistant" ? "bg-muted/30" : "bg-transparent"}`}
                  >
                    <div className="max-w-3xl mx-auto px-4 py-5">
                      <div className="flex gap-4">
                        <Avatar className="w-8 h-8 flex-shrink-0">
                          <AvatarFallback className={message.type === "user" ? "bg-primary text-primary-foreground" : "bg-[#10a37f] text-white"}>
                            {message.type === "user" ? <User className="w-4 h-4" /> : <Bot className="w-4 h-4" />}
                          </AvatarFallback>
                        </Avatar>

                        <div className="flex-1 space-y-3 min-w-0">
                          <div className="flex items-center gap-2">
                            <span className="text-sm font-semibold">
                              {message.type === "user" ? "You" : "InsightIQ"}
                            </span>
                            <span className="text-xs text-muted-foreground">
                              {formatTime(message.timestamp)}
                            </span>
                          </div>

                          <div className="space-y-3">
                            {message.sql && (
                              <div className="space-y-2">
                                <div className="flex items-center justify-between">
                                  <Badge variant="secondary" className="flex items-center gap-1.5 text-xs">
                                    <Database className="w-3 h-3" />
                                    SQL Query
                                  </Badge>
                                  <Button
                                    variant="ghost"
                                    size="icon"
                                    className="h-6 w-6"
                                    onClick={() => copyToClipboard(message.sql!, message.id)}
                                  >
                                    {copied === message.id ? (
                                      <Check className="w-3 h-3 text-green-500" />
                                    ) : (
                                      <Copy className="w-3 h-3" />
                                    )}
                                  </Button>
                                </div>
                                <pre className="p-3 rounded-lg bg-black/5 dark:bg-white/5 text-sm overflow-x-auto border border-black/10 dark:border-white/10">
                                  <code className="text-foreground">{message.sql}</code>
                                </pre>
                              </div>
                            )}

                            <div className="prose dark:prose-invert max-w-none">
                              <p className="text-[15px] leading-7 whitespace-pre-wrap">{message.content}</p>
                            </div>

                            {message.results && message.type === "assistant" && (
                              <ResultsCard results={message.results} />
                            )}
                          </div>
                        </div>
                      </div>
                    </div>
                  </motion.div>
                ))}
              </AnimatePresence>

              {(isSubmitting || audioProcessing) && (
                <motion.div
                  initial={{ opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  className="bg-muted/30"
                >
                  <div className="max-w-3xl mx-auto px-4 py-5">
                    <div className="flex gap-4">
                      <Avatar className="w-8 h-8 flex-shrink-0">
                        <AvatarFallback className="bg-[#10a37f] text-white">
                          <Bot className="w-4 h-4" />
                        </AvatarFallback>
                      </Avatar>

                      <div className="flex-1 space-y-3 min-w-0">
                        <div className="flex items-center gap-2">
                          <span className="text-sm font-semibold">InsightIQ</span>
                        </div>
                        <div className="flex items-center gap-2">
                          <motion.div
                            animate={{ scale: [1, 1.2, 1] }}
                            transition={{ duration: 1, repeat: Infinity }}
                            className="w-2 h-2 rounded-full bg-[#10a37f]"
                          />
                          <motion.div
                            animate={{ scale: [1, 1.2, 1] }}
                            transition={{ duration: 1, repeat: Infinity, delay: 0.2 }}
                            className="w-2 h-2 rounded-full bg-[#10a37f]"
                          />
                          <motion.div
                            animate={{ scale: [1, 1.2, 1] }}
                            transition={{ duration: 1, repeat: Infinity, delay: 0.4 }}
                            className="w-2 h-2 rounded-full bg-[#10a37f]"
                          />
                          <span className="text-sm text-muted-foreground ml-1">
                            {audioProcessing ? "Processing audio..." : "Analyzing..."}
                          </span>
                        </div>
                      </div>
                    </div>
                  </div>
                </motion.div>
              )}

              <div ref={messagesEndRef} className="h-32" />
            </div>
          </div>
        </div>

      {/* Input Area - Fixed at Bottom */}
      <div className="flex-shrink-0 border-t border-border bg-background">
        <div className="max-w-3xl mx-auto px-4 py-4 space-y-3">
          <AnimatePresence>
            {recordedBlob && (
              <motion.div
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -10 }}
                className="flex items-center gap-3 p-3 rounded-xl border border-emerald-500/30 bg-emerald-500/10"
              >
                <FileAudio className="w-5 h-5 text-emerald-600" />
                <span className="flex-1 text-sm font-medium">Audio ready</span>
                <div className="flex gap-2">
                  <Button size="sm" variant="ghost" onClick={() => setRecordedBlob(null)} className="h-8 rounded-lg">
                    Discard
                  </Button>
                  <Button size="sm" onClick={processRecordedAudio} disabled={audioProcessing} className="h-8 rounded-lg bg-emerald-600 hover:bg-emerald-700 text-white">
                    {audioProcessing ? "Processing..." : "Process"}
                  </Button>
                </div>
              </motion.div>
            )}
          </AnimatePresence>

          <AnimatePresence>
            {isRecording && (
              <motion.div
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -10 }}
                className="flex items-center justify-center gap-3 p-3 rounded-xl border border-red-500/30 bg-red-500/10"
              >
                <motion.div
                  animate={{ scale: [1, 1.3, 1] }}
                  transition={{ duration: 1, repeat: Infinity }}
                  className="w-2.5 h-2.5 rounded-full bg-red-500"
                />
                <span className="text-sm font-medium text-red-600">Recording... Click mic to stop</span>
              </motion.div>
            )}
          </AnimatePresence>

          <div className="relative">
            <div className="flex items-end gap-3 rounded-2xl border border-border bg-card p-2 shadow-sm focus-within:border-primary/50 focus-within:ring-1 focus-within:ring-primary/20 transition-all">
              <Button
                variant="ghost"
                size="icon"
                onClick={() => setShowSqlEditor(true)}
                className="flex-shrink-0 h-9 w-9 rounded-lg text-muted-foreground hover:text-foreground hover:bg-accent/50 transition-colors"
                title="Open SQL Editor"
              >
                <Database className="w-5 h-5" />
              </Button>

              <Textarea
                ref={textareaRef}
                value={input}
                onChange={(e) => setInput(e.target.value)}
                onKeyDown={handleKeyDown}
                placeholder="Message InsightIQ..."
                rows={1}
                disabled={isSubmitting || audioProcessing}
                className="flex-1 min-h-[36px] max-h-[200px] resize-none border-0 shadow-none focus-visible:ring-0 py-2 px-1 text-[15px] leading-6 bg-transparent placeholder:text-muted-foreground/50 disabled:opacity-50"
              />

              <div className="flex items-center gap-1 flex-shrink-0">
                <Button
                  variant={isRecording ? "destructive" : "ghost"}
                  size="icon"
                  onClick={isRecording ? stopRecording : startRecording}
                  disabled={isSubmitting || audioProcessing}
                  className="h-9 w-9 rounded-lg hover:bg-accent/50 transition-colors"
                  title={isRecording ? "Stop recording" : "Start voice input"}
                >
                  {isRecording ? (
                    <MicOff className="w-5 h-5" />
                  ) : (
                    <Mic className="w-5 h-5" />
                  )}
                </Button>

                <Button
                  onClick={handleSend}
                  disabled={!input.trim() || isSubmitting || audioProcessing}
                  size="icon"
                  className="h-9 w-9 rounded-lg bg-primary hover:bg-primary/90 disabled:bg-muted disabled:text-muted-foreground transition-all"
                  title="Send message"
                >
                  {isSubmitting ? (
                    <motion.div
                      animate={{ rotate: 360 }}
                      transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
                      className="w-4 h-4 border-2 border-current border-t-transparent rounded-full"
                    />
                  ) : (
                    <Send className="w-5 h-5" />
                  )}
                </Button>
              </div>
            </div>
          </div>

          <div className="flex items-center justify-center text-xs text-muted-foreground/60 pt-2">
            {messages.length > 0 ? (
              <Button
                variant="ghost"
                size="sm"
                onClick={clearChat}
                className="h-7 text-xs text-muted-foreground hover:text-foreground rounded-lg transition-colors"
              >
                <Trash2 className="w-3 h-3 mr-1.5" />
                Clear chat
              </Button>
            ) : (
              <span>
                <kbd className="px-1.5 py-0.5 rounded bg-muted text-[10px] font-mono">Enter</kbd>
                {" "}to send · {" "}
                <kbd className="px-1.5 py-0.5 rounded bg-muted text-[10px] font-mono">Shift+Enter</kbd>
                {" "}for new line
              </span>
            )}
          </div>
        </div>
      </div>

      <Sheet open={showSqlEditor} onOpenChange={setShowSqlEditor}>
        <SheetContent className="flex flex-col sm:max-w-2xl bg-white dark:bg-[#1e1e1e] border-l border-border shadow-2xl w-full">
          <SheetHeader className="pb-4 border-b border-border">
            <SheetTitle className="flex items-center gap-2 text-lg font-semibold">
              <div className="p-2 rounded-lg bg-primary/10">
                <Database className="w-5 h-5 text-primary" />
              </div>
              SQL Query Editor
            </SheetTitle>
            <p className="text-sm text-muted-foreground mt-2">
              Write and execute custom SQL queries directly
            </p>
          </SheetHeader>

          <div className="flex-1 overflow-auto py-6 space-y-4">
            <div className="space-y-3">
              <Label htmlFor="sql-query" className="text-sm font-medium flex items-center gap-2">
                <Code className="w-4 h-4 text-primary" />
                SQL Query
              </Label>
              <Textarea
                id="sql-query"
                value={sqlQuery}
                onChange={(e) => setSqlQuery(e.target.value)}
                placeholder="SELECT * FROM your_table LIMIT 10;"
                className="min-h-[300px] font-mono text-sm bg-card border-2 border-border rounded-xl resize-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:border-primary p-4 shadow-sm"
              />
            </div>

            <Alert className="bg-blue-50 dark:bg-blue-950/30 border-blue-200 dark:border-blue-900">
              <Sparkles className="w-4 h-4 text-blue-600 dark:text-blue-400" />
              <AlertDescription className="text-sm text-blue-900 dark:text-blue-100">
                <strong>Pro tip:</strong> Only SELECT queries are allowed for security. Use joins, WHERE clauses, and GROUP BY to refine your results.
              </AlertDescription>
            </Alert>

            <Card className="bg-muted/50">
              <CardContent className="p-4">
                <h4 className="text-sm font-semibold mb-2 flex items-center gap-2">
                  <Zap className="w-4 h-4 text-amber-500" />
                  Example Queries
                </h4>
                <div className="space-y-2 text-xs font-mono">
                  <div className="p-2 bg-background rounded border">SELECT * FROM users WHERE created_at &gt; NOW() - INTERVAL '7 days';</div>
                  <div className="p-2 bg-background rounded border">SELECT COUNT(*), category FROM products GROUP BY category;</div>
                </div>
              </CardContent>
            </Card>
          </div>

          <div className="flex justify-between items-center gap-3 pt-4 border-t border-border bg-white dark:bg-[#1e1e1e]">
            <div className="text-xs text-muted-foreground flex items-center gap-1">
              <kbd className="px-2 py-1 rounded bg-muted text-[10px] font-mono">Ctrl+Enter</kbd>
              <span>to execute</span>
            </div>
            <div className="flex gap-2">
              <Button variant="outline" onClick={() => setShowSqlEditor(false)} className="min-w-[100px]">
                Cancel
              </Button>
              <Button
                onClick={handleSqlSubmit}
                disabled={!sqlQuery.trim() || isSubmitting}
                className="min-w-[140px] bg-[#10a37f] hover:bg-[#0d8c6f]"
              >
                {isSubmitting ? (
                  <>
                    <motion.div
                      animate={{ rotate: 360 }}
                      transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
                      className="w-4 h-4 border-2 border-current border-t-transparent rounded-full mr-2"
                    />
                    Executing...
                  </>
                ) : (
                  <>
                    <Zap className="w-4 h-4 mr-2" />
                    Execute Query
                  </>
                )}
              </Button>
            </div>
          </div>
        </SheetContent>
      </Sheet>
    </div>
  );
}

function ResultsCard({ results }: { results: AnalyticsResponse | VoiceResponse }) {
  const data = "response" in results ? (results as VoiceResponse).response.data : (results as AnalyticsResponse).data;
  const insights = "response" in results ? (results as VoiceResponse).response.insights : (results as AnalyticsResponse).insights;
  const processTime = "response" in results ? (results as VoiceResponse).response.process_time : (results as AnalyticsResponse).process_time;
  const recordCount = data?.length || 0;

  if (!data || data.length === 0) return null;

  return (
    <div className="space-y-4 pt-2">
      {insights && (
        <Alert className="bg-[#10a37f]/5 border-[#10a37f]/20 dark:bg-[#10a37f]/10">
          <Wand2 className="w-4 h-4 text-[#10a37f]" />
          <AlertDescription className="text-sm text-foreground">{insights}</AlertDescription>
        </Alert>
      )}

      <div className="flex items-center gap-4 text-xs text-muted-foreground">
        <Badge variant="secondary" className="flex items-center gap-1 font-normal">
          <Table className="w-3 h-3" />
          {recordCount} records
        </Badge>
        {processTime && (
          <Badge variant="secondary" className="flex items-center gap-1 font-normal">
            <Clock className="w-3 h-3" />
            {processTime}
          </Badge>
        )}
      </div>

      {data.length > 0 && (
        <div className="rounded-xl border-2 border-border overflow-hidden bg-card shadow-sm">
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead className="bg-muted/50 border-b border-border">
                <tr>
                  {Object.keys(data[0]).map((key) => (
                    <th
                      key={key}
                      className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wider text-foreground"
                    >
                      {key.replace(/_/g, " ")}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {data.slice(0, 5).map((row, i) => (
                  <tr key={i} className="hover:bg-muted/30 transition-colors">
                    {Object.values(row).map((value: any, j) => (
                      <td key={j} className="px-4 py-3 text-foreground">
                        {typeof value === "number" ? value.toLocaleString() : String(value)}
                      </td>
                    ))}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          {data.length > 5 && (
            <div className="px-4 py-2.5 bg-muted/50 text-xs text-center text-muted-foreground border-t border-border font-medium">
              Showing 5 of {data.length} records • <button className="text-primary hover:underline">View all</button>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
