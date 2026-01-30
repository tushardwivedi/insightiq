"use client";

import { useRef } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import {
  Send,
  Mic,
  MicOff,
  FileAudio,
  Database,
  Trash2,
} from "lucide-react";

interface ChatInputProps {
  input: string;
  onInputChange: (value: string) => void;
  onSend: () => void;
  onOpenSqlEditor: () => void;
  isSubmitting: boolean;
  audioProcessing: boolean;
  isRecording: boolean;
  recordedBlob: Blob | null;
  onStartRecording: () => void;
  onStopRecording: () => void;
  onProcessAudio: () => void;
  onDiscardAudio: () => void;
  onClearChat: () => void;
  hasMessages: boolean;
}

export function ChatInput({
  input,
  onInputChange,
  onSend,
  onOpenSqlEditor,
  isSubmitting,
  audioProcessing,
  isRecording,
  recordedBlob,
  onStartRecording,
  onStopRecording,
  onProcessAudio,
  onDiscardAudio,
  onClearChat,
  hasMessages,
}: ChatInputProps) {
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      onSend();
    }
  };

  return (
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
                <Button size="sm" variant="ghost" onClick={onDiscardAudio} className="h-8 rounded-lg">
                  Discard
                </Button>
                <Button size="sm" onClick={onProcessAudio} disabled={audioProcessing} className="h-8 rounded-lg bg-emerald-600 hover:bg-emerald-700 text-white">
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
              onClick={onOpenSqlEditor}
              className="flex-shrink-0 h-9 w-9 rounded-lg text-muted-foreground hover:text-foreground hover:bg-accent/50 transition-colors"
              title="Open SQL Editor"
            >
              <Database className="w-5 h-5" />
            </Button>

            <Textarea
              ref={textareaRef}
              value={input}
              onChange={(e) => onInputChange(e.target.value)}
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
                onClick={isRecording ? onStopRecording : onStartRecording}
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
                onClick={onSend}
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
          {hasMessages ? (
            <Button
              variant="ghost"
              size="sm"
              onClick={onClearChat}
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
  );
}
