"use client";

import { useState, useRef, useEffect, useCallback } from "react";
import { Card, CardContent, CardHeader, CardTitle, CardDescription, CardFooter } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetFooter } from "@/components/ui/sheet";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { motion, AnimatePresence } from "framer-motion";
import {
  Send,
  Mic,
  MicOff,
  Upload,
  FileAudio,
  Sparkles,
  Database,
  Zap,
  Plus,
  X,
  Copy,
  Check,
  Clock,
  Trash2,
  Database as DatabaseIcon,
  Bot,
  User,
  Wand2,
  Table,
  BarChart3,
  ChevronDown,
  MoreHorizontal
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
      <div className="flex-1 overflow-hidden flex flex-col">
        <div className="flex-1 overflow-auto p-4">
          <div className="max-w-3xl mx-auto space-y-6">
            <AnimatePresence mode="wait">
              {messages.length === 0 && (
                <motion.div
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: -20 }}
                  className="space-y-8"
                >
                  <div className="text-center space-y-4 pt-8">
                    <div className="w-16 h-16 mx-auto rounded-2xl bg-primary/10 flex items-center justify-center">
                      <Bot className="w-8 h-8 text-primary" />
                    </div>
                    <div>
                      <h1 className="text-2xl font-semibold">How can I help you today?</h1>
                      <p className="text-muted-foreground mt-2">
                        Ask questions about your data or write SQL queries
                      </p>
                    </div>
                  </div>

                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                    <Card
                      className="cursor-pointer hover:bg-muted/50 transition-colors border-dashed"
                      onClick={() => setShowSqlEditor(true)}
                    >
                      <CardContent className="p-4 flex items-center gap-4">
                        <div className="w-10 h-10 rounded-lg bg-blue-500/10 flex items-center justify-center">
                          <Database className="w-5 h-5 text-blue-500" />
                        </div>
                        <div className="flex-1">
                          <p className="font-medium text-sm">Write SQL Query</p>
                          <p className="text-xs text-muted-foreground">Use the SQL editor</p>
                        </div>
                        <Plus className="w-4 h-4 text-muted-foreground" />
                      </CardContent>
                    </Card>

                    {suggestedQueries.map((query, i) => (
                      <Card
                        key={i}
                        className="cursor-pointer hover:bg-muted/50 transition-colors"
                        onClick={() => setInput(query)}
                      >
                        <CardContent className="p-4">
                          <p className="text-sm text-muted-foreground group-hover:text-foreground transition-colors">
                            {query}
                          </p>
                        </CardContent>
                      </Card>
                    ))}
                  </div>

                  <div className="flex items-center justify-center gap-4 pt-4">
                    <Badge variant="secondary" className="flex items-center gap-2 px-3 py-1">
                      <DatabaseIcon className="w-3 h-3" />
                      Connected to 5 data sources
                    </Badge>
                    <Badge variant="secondary" className="flex items-center gap-2 px-3 py-1">
                      <Sparkles className="w-3 h-3" />
                      AI Enabled
                    </Badge>
                  </div>
                </motion.div>
              )}
            </AnimatePresence>

            <div className="space-y-4">
              <AnimatePresence>
                {messages.map((message, index) => (
                  <motion.div
                    key={message.id}
                    initial={{ opacity: 0, y: 10 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, y: -10 }}
                    transition={{ delay: index * 0.03 }}
                  >
                    <div className={`flex gap-3 ${message.type === "assistant" ? "flex-row" : "flex-row-reverse"}`}>
                      <Avatar className="w-8 h-8 mt-1">
                        <AvatarFallback>
                          {message.type === "user" ? <User className="w-4 h-4" /> : <Bot className="w-4 h-4" />}
                        </AvatarFallback>
                      </Avatar>

                      <div className={`flex-1 max-w-[85%] ${message.type === "user" ? "text-right" : ""}`}>
                        <div className="flex items-center gap-2 mb-1.5">
                          <span className="text-sm font-medium">
                            {message.type === "user" ? "You" : "InsightIQ"}
                          </span>
                          <span className="text-xs text-muted-foreground">
                            {formatTime(message.timestamp)}
                          </span>
                        </div>

                        <Card className={`${message.type === "assistant" ? "" : "bg-primary text-primary-foreground"}`}>
                          <CardContent className="p-4 space-y-3">
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
                                <pre className="p-3 rounded-lg bg-muted text-sm overflow-x-auto">
                                  <code className="text-foreground">{message.sql}</code>
                                </pre>
                              </div>
                            )}

                            {message.type === "user" ? (
                              <p className="text-sm whitespace-pre-wrap">{message.content}</p>
                            ) : (
                              <div className="space-y-3">
                                <p className="text-sm whitespace-pre-wrap">{message.content}</p>
                                {message.results && (
                                  <ResultsCard results={message.results} />
                                )}
                              </div>
                            )}
                          </CardContent>
                        </Card>
                      </div>
                    </div>
                  </motion.div>
                ))}
              </AnimatePresence>

              {(isSubmitting || audioProcessing) && (
                <motion.div
                  initial={{ opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  className="flex gap-3"
                >
                  <Avatar className="w-8 h-8 mt-1">
                    <AvatarFallback>
                      <Bot className="w-4 h-4" />
                    </AvatarFallback>
                  </Avatar>

                  <div className="flex-1 max-w-[85%]">
                    <div className="flex items-center gap-2 mb-1.5">
                      <span className="text-sm font-medium">InsightIQ</span>
                    </div>
                    <Card>
                      <CardContent className="p-4">
                        <div className="flex items-center gap-2">
                          <motion.div
                            animate={{ scale: [1, 1.2, 1] }}
                            transition={{ duration: 1, repeat: Infinity }}
                            className="w-2 h-2 rounded-full bg-primary"
                          />
                          <motion.div
                            animate={{ scale: [1, 1.2, 1] }}
                            transition={{ duration: 1, repeat: Infinity, delay: 0.2 }}
                            className="w-2 h-2 rounded-full bg-primary"
                          />
                          <motion.div
                            animate={{ scale: [1, 1.2, 1] }}
                            transition={{ duration: 1, repeat: Infinity, delay: 0.4 }}
                            className="w-2 h-2 rounded-full bg-primary"
                          />
                          <span className="text-sm text-muted-foreground">
                            {audioProcessing ? "Processing audio..." : "Analyzing..."}
                          </span>
                        </div>
                      </CardContent>
                    </Card>
                  </div>
                </motion.div>
              )}

              <div ref={messagesEndRef} />
            </div>
          </div>
        </div>

        <div className="border-t bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
          <div className="max-w-3xl mx-auto p-4 space-y-3">
            <AnimatePresence>
              {recordedBlob && (
                <motion.div
                  initial={{ opacity: 0, height: 0 }}
                  animate={{ opacity: 1, height: "auto" }}
                  exit={{ opacity: 0, height: 0 }}
                  className="flex items-center gap-3 p-3 rounded-lg border bg-green-500/5"
                >
                  <FileAudio className="w-4 h-4 text-green-500" />
                  <span className="text-sm flex-1">Audio recorded</span>
                  <div className="flex gap-2">
                    <Button size="sm" variant="ghost" onClick={() => setRecordedBlob(null)}>
                      Discard
                    </Button>
                    <Button size="sm" onClick={processRecordedAudio} disabled={audioProcessing}>
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
                  className="flex items-center justify-center"
                >
                  <div className="flex items-center gap-2 px-4 py-2 rounded-full bg-red-500/10 text-red-500 text-sm">
                    <motion.div
                      animate={{ scale: [1, 1.5, 1] }}
                      transition={{ duration: 1, repeat: Infinity }}
                      className="w-2 h-2 rounded-full bg-red-500"
                    />
                    Recording... Click mic to stop
                  </div>
                </motion.div>
              )}
            </AnimatePresence>

            <Card className="border-2 focus-within:border-primary/50 transition-colors">
              <CardContent className="p-3">
                <div className="flex items-end gap-2">
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => setShowSqlEditor(true)}
                    className="flex-shrink-0 mb-1 text-muted-foreground hover:text-foreground"
                    title="Open SQL Editor"
                  >
                    <Database className="w-5 h-5" />
                  </Button>

                  <Textarea
                    ref={textareaRef}
                    value={input}
                    onChange={(e) => setInput(e.target.value)}
                    onKeyDown={handleKeyDown}
                    placeholder="Ask a question or describe the analysis you need..."
                    rows={1}
                    disabled={isSubmitting || audioProcessing}
                    className="min-h-[44px] max-h-[150px] resize-none border-0 shadow-none focus-visible:ring-0 p-0 text-sm bg-transparent"
                  />

                  <div className="flex items-center gap-1 flex-shrink-0">
                    <Button
                      variant={isRecording ? "destructive" : "ghost"}
                      size="icon"
                      onClick={isRecording ? stopRecording : startRecording}
                      disabled={isSubmitting || audioProcessing}
                      className="mb-1"
                    >
                      {isRecording ? (
                        <MicOff className="w-4 h-4" />
                      ) : (
                        <Mic className="w-4 h-4" />
                      )}
                    </Button>

                    <Button
                      onClick={handleSend}
                      disabled={!input.trim() || isSubmitting || audioProcessing}
                      size="icon"
                      className="mb-1"
                    >
                      {isSubmitting ? (
                        <motion.div
                          animate={{ rotate: 360 }}
                          transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
                          className="w-4 h-4 border-2 border-current border-t-transparent rounded-full"
                        />
                      ) : (
                        <Send className="w-4 h-4" />
                      )}
                    </Button>
                  </div>
                </div>
              </CardContent>
            </Card>

            {messages.length > 0 && (
              <div className="flex items-center justify-between text-xs text-muted-foreground px-1">
                <Button variant="ghost" size="sm" onClick={clearChat} className="h-7 text-muted-foreground hover:text-foreground">
                  <Trash2 className="w-3 h-3 mr-1.5" />
                  Clear chat
                </Button>
                <span className="flex items-center gap-1">
                  <kbd className="px-1.5 py-0.5 rounded bg-muted text-[10px]">↵</kbd>
                  to send
                  <kbd className="px-1.5 py-0.5 rounded bg-muted text-[10px] ml-2">⇧</kbd>
                  <kbd className="px-1.5 py-0.5 rounded bg-muted text-[10px]">↵</kbd>
                  for new line
                </span>
              </div>
            )}
          </div>
        </div>
      </div>

      <Sheet open={showSqlEditor} onOpenChange={setShowSqlEditor}>
        <SheetContent className="flex flex-col sm:max-w-xl bg-background border-l w-full">
          <SheetHeader className="pb-4 border-b">
            <SheetTitle className="flex items-center gap-2 text-lg">
              <Database className="w-5 h-5 text-primary" />
              SQL Query Editor
            </SheetTitle>
          </SheetHeader>

          <div className="flex-1 overflow-auto py-4 space-y-4">
            <div className="space-y-2">
              <Label htmlFor="sql-query" className="text-sm font-medium">SQL Query</Label>
              <Textarea
                id="sql-query"
                value={sqlQuery}
                onChange={(e) => setSqlQuery(e.target.value)}
                placeholder="SELECT * FROM your_table LIMIT 10"
                className="min-h-[200px] font-mono text-sm bg-background border resize-none focus-visible:ring-1 focus-visible:ring-primary"
              />
            </div>

            <Alert className="bg-muted/50 border-muted">
              <AlertDescription className="text-xs text-muted-foreground">
                Only SELECT queries are allowed for security reasons.
              </AlertDescription>
            </Alert>
          </div>

          <div className="flex justify-end gap-2 pt-4 border-t bg-background">
            <Button variant="outline" onClick={() => setShowSqlEditor(false)}>
              Cancel
            </Button>
            <Button
              onClick={handleSqlSubmit}
              disabled={!sqlQuery.trim() || isSubmitting}
              className="min-w-[120px]"
            >
              {isSubmitting ? (
                <>
                  <motion.div
                    animate={{ rotate: 360 }}
                    transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
                    className="w-4 h-4 border-2 border-current border-t-transparent rounded-full mr-2"
                  />
                  Processing...
                </>
              ) : (
                <>
                  <Zap className="w-4 h-4 mr-2" />
                  Execute Query
                </>
              )}
            </Button>
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
    <div className="space-y-3 pt-2">
      {insights && (
        <Alert className="bg-primary/5 border-primary/20">
          <Wand2 className="w-4 h-4 text-primary" />
          <AlertDescription className="text-sm">{insights}</AlertDescription>
        </Alert>
      )}

      <div className="flex items-center gap-4 text-xs text-muted-foreground">
        <span className="flex items-center gap-1">
          <Table className="w-3 h-3" />
          {recordCount} records
        </span>
        {processTime && (
          <span className="flex items-center gap-1">
            <Clock className="w-3 h-3" />
            {processTime}
          </span>
        )}
      </div>

      {data.length > 0 && (
        <div className="rounded-md border overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead className="bg-muted/50">
                <tr>
                  {Object.keys(data[0]).map((key) => (
                    <th
                      key={key}
                      className="px-3 py-2 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground"
                    >
                      {key.replace(/_/g, " ")}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody className="divide-y">
                {data.slice(0, 5).map((row, i) => (
                  <tr key={i} className="hover:bg-muted/30 transition-colors">
                    {Object.values(row).map((value: any, j) => (
                      <td key={j} className="px-3 py-2">
                        {typeof value === "number" ? value.toLocaleString() : String(value)}
                      </td>
                    ))}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          {data.length > 5 && (
            <div className="px-3 py-2 bg-muted/30 text-xs text-center text-muted-foreground">
              Showing 5 of {data.length} records
            </div>
          )}
        </div>
      )}
    </div>
  );
}
