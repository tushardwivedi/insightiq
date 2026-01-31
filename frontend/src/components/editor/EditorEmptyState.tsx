"use client";

import { motion } from "framer-motion";
import { Bot, Database, BarChart3, TrendingUp, Users, Sparkles } from "lucide-react";

interface EditorEmptyStateProps {
  suggestedQueries: string[];
  onSelectQuery: (query: string) => void;
  onOpenSqlEditor: () => void;
}

const queryIcons = [TrendingUp, BarChart3, Users, Sparkles];

export function EditorEmptyState({ suggestedQueries, onSelectQuery, onOpenSqlEditor }: EditorEmptyStateProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, y: -20 }}
      className="flex flex-col items-center justify-center px-4 pt-16 pb-8"
    >
      {/* Animated hero icon */}
      <motion.div
        initial={{ scale: 0.8, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        transition={{ delay: 0.1, type: "spring", stiffness: 200, damping: 20 }}
        className="relative mb-6"
      >
        <div className="absolute inset-0 rounded-2xl bg-gradient-to-br from-primary/20 via-cyan-500/20 to-purple-500/20 blur-xl" />
        <div className="relative w-16 h-16 rounded-2xl bg-gradient-to-br from-primary via-cyan-500 to-primary flex items-center justify-center shadow-lg shadow-primary/25">
          <Bot className="w-8 h-8 text-white" />
        </div>
        <motion.div
          animate={{ scale: [1, 1.2, 1], opacity: [0.5, 0.8, 0.5] }}
          transition={{ duration: 3, repeat: Infinity, ease: "easeInOut" }}
          className="absolute -inset-2 rounded-3xl bg-gradient-to-br from-primary/10 to-cyan-500/10 -z-10"
        />
      </motion.div>

      {/* Heading */}
      <motion.div
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.2 }}
        className="text-center mb-8"
      >
        <h1 className="text-3xl font-bold mb-2 bg-gradient-to-r from-foreground via-foreground to-muted-foreground bg-clip-text">
          Chat with your data
        </h1>
        <p className="text-muted-foreground text-base max-w-md">
          Ask questions in natural language, run SQL queries, or use voice input to explore your datasets
        </p>
      </motion.div>

      {/* Action cards grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 max-w-2xl mx-auto w-full">
        {/* SQL Editor card */}
        <motion.button
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.3 }}
          whileHover={{ scale: 1.02, y: -2 }}
          whileTap={{ scale: 0.98 }}
          onClick={onOpenSqlEditor}
          className="group relative flex items-center gap-4 p-4 rounded-2xl border border-border bg-card hover:border-primary/40 transition-all text-left overflow-hidden"
        >
          <div className="absolute inset-0 bg-gradient-to-br from-primary/5 to-cyan-500/5 opacity-0 group-hover:opacity-100 transition-opacity" />
          <div className="relative w-10 h-10 rounded-xl bg-gradient-to-br from-primary/10 to-cyan-500/10 flex items-center justify-center flex-shrink-0 group-hover:from-primary/20 group-hover:to-cyan-500/20 transition-colors">
            <Database className="w-5 h-5 text-primary" />
          </div>
          <div className="relative flex-1 min-w-0">
            <p className="font-semibold text-sm text-foreground">Write SQL Query</p>
            <p className="text-xs text-muted-foreground mt-0.5">Open the SQL editor for custom queries</p>
          </div>
        </motion.button>

        {/* Suggested query cards */}
        {suggestedQueries.map((query, i) => {
          const Icon = queryIcons[i % queryIcons.length];
          return (
            <motion.button
              key={i}
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.35 + i * 0.05 }}
              whileHover={{ scale: 1.02, y: -2 }}
              whileTap={{ scale: 0.98 }}
              onClick={() => onSelectQuery(query)}
              className="group relative flex items-center gap-4 p-4 rounded-2xl border border-border bg-card hover:border-primary/40 transition-all text-left overflow-hidden"
            >
              <div className="absolute inset-0 bg-gradient-to-br from-primary/5 to-cyan-500/5 opacity-0 group-hover:opacity-100 transition-opacity" />
              <div className="relative w-10 h-10 rounded-xl bg-muted/50 flex items-center justify-center flex-shrink-0 group-hover:bg-primary/10 transition-colors">
                <Icon className="w-5 h-5 text-muted-foreground group-hover:text-primary transition-colors" />
              </div>
              <div className="relative flex-1 min-w-0">
                <p className="text-sm text-foreground leading-snug">{query}</p>
              </div>
            </motion.button>
          );
        })}
      </div>
    </motion.div>
  );
}
