"use client";

import { motion } from "framer-motion";
import { Bot, Database } from "lucide-react";

interface EditorEmptyStateProps {
  suggestedQueries: string[];
  onSelectQuery: (query: string) => void;
  onOpenSqlEditor: () => void;
}

export function EditorEmptyState({ suggestedQueries, onSelectQuery, onOpenSqlEditor }: EditorEmptyStateProps) {
  return (
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
          onClick={onOpenSqlEditor}
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
            onClick={() => onSelectQuery(query)}
            className="p-3 rounded-xl border border-border bg-card hover:bg-accent/50 transition-colors text-left text-sm"
          >
            {query}
          </button>
        ))}
      </div>
    </motion.div>
  );
}
