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
      initial={{ opacity: 0, y: 8 }}
      animate={{ opacity: 1, y: 0 }}
      className="py-6 bg-muted/20"
    >
      <div className="max-w-3xl mx-auto px-4">
        <div className="flex gap-4">
          <Avatar className="w-8 h-8 flex-shrink-0 mt-0.5">
            <AvatarFallback className="bg-gradient-to-br from-emerald-500 to-teal-600 text-white">
              <Bot className="w-4 h-4" />
            </AvatarFallback>
          </Avatar>

          <div className="flex-1 space-y-3 min-w-0">
            <div className="flex items-center gap-2">
              <span className="text-sm font-semibold text-foreground">InsightIQ</span>
            </div>
            <div className="flex items-center gap-3">
              <div className="flex items-center gap-1.5">
                {[0, 1, 2].map((i) => (
                  <motion.div
                    key={i}
                    animate={{
                      y: [0, -6, 0],
                      opacity: [0.4, 1, 0.4],
                    }}
                    transition={{
                      duration: 1.2,
                      repeat: Infinity,
                      delay: i * 0.15,
                      ease: "easeInOut",
                    }}
                    className="w-2 h-2 rounded-full bg-emerald-500"
                  />
                ))}
              </div>
              <span className="text-sm text-muted-foreground">
                {isAudioProcessing ? "Processing audio..." : "Thinking..."}
              </span>
            </div>
          </div>
        </div>
      </div>
    </motion.div>
  );
}
