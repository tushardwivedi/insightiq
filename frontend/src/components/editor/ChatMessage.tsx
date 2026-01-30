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

  return (
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
                      onClick={() => onCopy(message.sql!, message.id)}
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
  );
}
