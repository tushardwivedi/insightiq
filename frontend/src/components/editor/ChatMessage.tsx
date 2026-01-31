"use client";

import { motion } from "framer-motion";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Database, Copy, Check, Bot, User } from "lucide-react";
import { AnalyticsResponse, VoiceResponse } from "@/types";
import { ResultsCard } from "./ResultsCard";

interface Message {
  id: string;
  type: "user" | "assistant";
  content: string;
  sql?: string;
  results?: AnalyticsResponse | VoiceResponse;
  timestamp: Date;
}

interface ChatMessageProps {
  message: Message;
  index: number;
  copied: string | null;
  onCopy: (text: string, id: string) => void;
}

export function ChatMessage({ message, index, copied, onCopy }: ChatMessageProps) {
  const formatTime = (date: Date) => {
    return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
  };

  const isUser = message.type === "user";

  return (
    <motion.div
      key={message.id}
      initial={{ opacity: 0, y: 8 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, y: -8 }}
      transition={{ delay: index * 0.03, duration: 0.3 }}
      className={`py-6 ${!isUser ? "bg-muted/20" : ""}`}
    >
      <div className="max-w-3xl mx-auto px-4">
        <div className="flex gap-4">
          {/* Avatar */}
          <Avatar className="w-8 h-8 flex-shrink-0 mt-0.5">
            <AvatarFallback
              className={
                isUser
                  ? "bg-gradient-to-br from-primary to-primary/80 text-primary-foreground"
                  : "bg-gradient-to-br from-emerald-500 to-teal-600 text-white"
              }
            >
              {isUser ? <User className="w-4 h-4" /> : <Bot className="w-4 h-4" />}
            </AvatarFallback>
          </Avatar>

          {/* Content */}
          <div className="flex-1 space-y-3 min-w-0">
            {/* Name & time */}
            <div className="flex items-center gap-2">
              <span className="text-sm font-semibold text-foreground">
                {isUser ? "You" : "InsightIQ"}
              </span>
              <span className="text-[11px] text-muted-foreground/60">
                {formatTime(message.timestamp)}
              </span>
            </div>

            <div className="space-y-3">
              {/* SQL code block */}
              {message.sql && (
                <div className="space-y-2">
                  <div className="flex items-center justify-between">
                    <Badge
                      variant="secondary"
                      className="flex items-center gap-1.5 text-xs font-medium bg-primary/5 text-primary border-primary/10"
                    >
                      <Database className="w-3 h-3" />
                      SQL Query
                    </Badge>
                    <Button
                      variant="ghost"
                      size="icon"
                      className="h-7 w-7 rounded-lg hover:bg-primary/5"
                      onClick={() => onCopy(message.sql!, message.id)}
                    >
                      {copied === message.id ? (
                        <Check className="w-3.5 h-3.5 text-emerald-500" />
                      ) : (
                        <Copy className="w-3.5 h-3.5 text-muted-foreground" />
                      )}
                    </Button>
                  </div>
                  <div className="relative rounded-xl overflow-hidden border border-border/50">
                    <div className="absolute top-0 left-0 right-0 h-8 bg-muted/50 flex items-center px-3 border-b border-border/50">
                      <div className="flex gap-1.5">
                        <div className="w-2.5 h-2.5 rounded-full bg-red-400/50" />
                        <div className="w-2.5 h-2.5 rounded-full bg-yellow-400/50" />
                        <div className="w-2.5 h-2.5 rounded-full bg-green-400/50" />
                      </div>
                    </div>
                    <pre className="p-4 pt-12 bg-muted/20 text-sm overflow-x-auto">
                      <code className="text-foreground font-mono text-[13px] leading-6">{message.sql}</code>
                    </pre>
                  </div>
                </div>
              )}

              {/* Message content */}
              <div className="prose dark:prose-invert max-w-none">
                <p className="text-[15px] leading-7 whitespace-pre-wrap text-foreground/90">{message.content}</p>
              </div>

              {/* Results */}
              {message.results && !isUser && (
                <ResultsCard results={message.results} />
              )}
            </div>
          </div>
        </div>
      </div>
    </motion.div>
  );
}
