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
  ArrowUp,
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
    <div className="flex-shrink-0 bg-gradient-to-t from-background via-background to-background/80 pt-2 pb-4">
      <div className="max-w-3xl mx-auto px-4 space-y-2">
        {/* Audio banners */}
        <AnimatePresence>
          {recordedBlob && (
            <motion.div
              initial={{ opacity: 0, y: 10, scale: 0.95 }}
              animate={{ opacity: 1, y: 0, scale: 1 }}
              exit={{ opacity: 0, y: -10, scale: 0.95 }}
              className="flex items-center gap-3 p-3 rounded-2xl border border-emerald-500/20 bg-emerald-500/5 backdrop-blur-sm"
            >
              <div className="w-8 h-8 rounded-full bg-emerald-500/10 flex items-center justify-center">
                <FileAudio className="w-4 h-4 text-emerald-500" />
              </div>
              <span className="flex-1 text-sm font-medium text-emerald-700 dark:text-emerald-400">
                Audio ready to process
              </span>
              <div className="flex gap-2">
                <Button
                  size="sm"
                  variant="ghost"
                  onClick={onDiscardAudio}
                  className="h-8 rounded-full text-xs hover:bg-red-500/10 hover:text-red-600"
                >
                  Discard
                </Button>
                <Button
                  size="sm"
                  onClick={onProcessAudio}
                  disabled={audioProcessing}
                  className="h-8 rounded-full bg-emerald-600 hover:bg-emerald-700 text-white text-xs shadow-sm"
                >
                  {audioProcessing ? "Processing..." : "Process Audio"}
                </Button>
              </div>
            </motion.div>
          )}
        </AnimatePresence>

        <AnimatePresence>
          {isRecording && (
            <motion.div
              initial={{ opacity: 0, y: 10, scale: 0.95 }}
              animate={{ opacity: 1, y: 0, scale: 1 }}
              exit={{ opacity: 0, y: -10, scale: 0.95 }}
              className="flex items-center justify-center gap-3 p-3 rounded-2xl border border-red-500/20 bg-red-500/5 backdrop-blur-sm"
            >
              <div className="flex items-center gap-2">
                <motion.div
                  animate={{ scale: [1, 1.4, 1], opacity: [1, 0.5, 1] }}
                  transition={{ duration: 1.5, repeat: Infinity }}
                  className="w-2 h-2 rounded-full bg-red-500"
                />
                <motion.div
                  animate={{ scale: [1, 1.4, 1], opacity: [1, 0.5, 1] }}
                  transition={{ duration: 1.5, repeat: Infinity, delay: 0.2 }}
                  className="w-1.5 h-1.5 rounded-full bg-red-400"
                />
                <motion.div
                  animate={{ scale: [1, 1.4, 1], opacity: [1, 0.5, 1] }}
                  transition={{ duration: 1.5, repeat: Infinity, delay: 0.4 }}
                  className="w-1 h-1 rounded-full bg-red-300"
                />
              </div>
              <span className="text-sm font-medium text-red-600 dark:text-red-400">
                Recording... Click mic to stop
              </span>
            </motion.div>
          )}
        </AnimatePresence>

        {/* Main input container */}
        <div className="relative">
          <div className="rounded-2xl border border-border bg-card shadow-lg shadow-black/5 dark:shadow-black/20 focus-within:border-primary/30 focus-within:shadow-primary/5 transition-all duration-200">
            {/* Textarea row */}
            <div className="flex items-end gap-2 p-3">
              <Textarea
                ref={textareaRef}
                value={input}
                onChange={(e) => onInputChange(e.target.value)}
                onKeyDown={handleKeyDown}
                placeholder="Ask anything about your data..."
                rows={1}
                disabled={isSubmitting || audioProcessing}
                className="flex-1 min-h-[44px] max-h-[200px] resize-none border-0 shadow-none focus-visible:ring-0 py-2.5 px-2 text-[15px] leading-6 bg-transparent placeholder:text-muted-foreground/50 disabled:opacity-50"
              />
              <motion.button
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
                onClick={onSend}
                disabled={!input.trim() || isSubmitting || audioProcessing}
                className="flex-shrink-0 h-9 w-9 rounded-full bg-primary text-primary-foreground flex items-center justify-center disabled:bg-muted disabled:text-muted-foreground transition-colors shadow-sm"
                title="Send message"
              >
                {isSubmitting ? (
                  <motion.div
                    animate={{ rotate: 360 }}
                    transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
                    className="w-4 h-4 border-2 border-current border-t-transparent rounded-full"
                  />
                ) : (
                  <ArrowUp className="w-4 h-4" />
                )}
              </motion.button>
            </div>

            {/* Bottom actions bar */}
            <div className="flex items-center justify-between px-3 pb-3">
              <div className="flex items-center gap-1">
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={onOpenSqlEditor}
                  className="h-8 rounded-full text-xs text-muted-foreground hover:text-foreground gap-1.5 px-3"
                  title="Open SQL Editor"
                >
                  <Database className="w-3.5 h-3.5" />
                  SQL
                </Button>

                <Button
                  variant={isRecording ? "destructive" : "ghost"}
                  size="sm"
                  onClick={isRecording ? onStopRecording : onStartRecording}
                  disabled={isSubmitting || audioProcessing}
                  className="h-8 rounded-full text-xs gap-1.5 px-3"
                  title={isRecording ? "Stop recording" : "Voice input"}
                >
                  {isRecording ? (
                    <MicOff className="w-3.5 h-3.5" />
                  ) : (
                    <Mic className="w-3.5 h-3.5" />
                  )}
                  {isRecording ? "Stop" : "Voice"}
                </Button>
              </div>

              <div className="flex items-center gap-2">
                {hasMessages ? (
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={onClearChat}
                    className="h-8 rounded-full text-xs text-muted-foreground hover:text-red-600 hover:bg-red-500/10 gap-1.5 px-3 transition-colors"
                  >
                    <Trash2 className="w-3 h-3" />
                    Clear
                  </Button>
                ) : (
                  <span className="text-[11px] text-muted-foreground/50 pr-1">
                    <kbd className="px-1.5 py-0.5 rounded bg-muted/50 font-mono text-[10px]">Enter</kbd>
                    {" "}send{" "}
                    <kbd className="px-1.5 py-0.5 rounded bg-muted/50 font-mono text-[10px]">Shift+Enter</kbd>
                    {" "}new line
                  </span>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
