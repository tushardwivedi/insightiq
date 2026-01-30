"use client";

import { motion } from "framer-motion";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Bot } from "lucide-react";

interface TypingIndicatorProps {
  isAudioProcessing?: boolean;
}

export function TypingIndicator({ isAudioProcessing = false }: TypingIndicatorProps) {
  return (
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
                {isAudioProcessing ? "Processing audio..." : "Analyzing..."}
              </span>
            </div>
          </div>
        </div>
      </div>
    </motion.div>
  );
}
