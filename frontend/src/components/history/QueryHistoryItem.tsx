"use client";

import { motion } from "framer-motion";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Play, Copy, Check, Trash2, Calendar } from "lucide-react";

interface QueryHistoryItemProps {
  item: any;
  index: number;
  copied: string | null;
  onRerun: (query: string) => void;
  onCopy: (text: string, id: string) => void;
  onDelete: (id: string) => void;
}

function getQueryType(query: string) {
  if (query?.startsWith("SELECT")) return "SQL";
  if (query?.length > 0) return "Text";
  return "Unknown";
}

function formatDate(date: string) {
  return new Date(date).toLocaleDateString(undefined, {
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

export function QueryHistoryItem({ item, index, copied, onRerun, onCopy, onDelete }: QueryHistoryItemProps) {
  const query = item.query || item.sql;

  return (
    <motion.div
      key={item.id || index}
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: index * 0.03 }}
      className="p-4 rounded-lg border hover:bg-muted/30 transition-colors"
    >
      <div className="flex items-start justify-between gap-4">
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 mb-2">
            <Badge variant="secondary" className="flex items-center gap-1">
              {getQueryType(query)}
            </Badge>
            <span className="text-xs text-muted-foreground flex items-center gap-1">
              <Calendar className="w-3 h-3" />
              {formatDate(item.created_at || item.timestamp)}
            </span>
            {item.process_time && (
              <span className="text-xs text-muted-foreground">
                {item.process_time}
              </span>
            )}
          </div>
          <p className="text-sm font-mono bg-muted/50 p-2 rounded truncate">
            {query}
          </p>
          {item.insights && (
            <p className="text-sm text-muted-foreground mt-2 line-clamp-2">
              {item.insights}
            </p>
          )}
        </div>
        <div className="flex items-center gap-1">
          <Button
            variant="ghost"
            size="icon"
            onClick={() => onRerun(query)}
            title="Rerun query"
          >
            <Play className="w-4 h-4" />
          </Button>
          <Button
            variant="ghost"
            size="icon"
            onClick={() => onCopy(query, item.id)}
            title="Copy query"
          >
            {copied === item.id ? (
              <Check className="w-4 h-4 text-green-500" />
            ) : (
              <Copy className="w-4 h-4" />
            )}
          </Button>
          <Button
            variant="ghost"
            size="icon"
            onClick={() => onDelete(item.id)}
            title="Delete"
          >
            <Trash2 className="w-4 h-4 text-destructive" />
          </Button>
        </div>
      </div>
    </motion.div>
  );
}
