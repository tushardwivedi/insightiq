"use client";

import { motion } from "framer-motion";
import { Sheet, SheetContent, SheetHeader, SheetTitle } from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Database, Code, Sparkles, Zap, Play } from "lucide-react";

interface SqlEditorSheetProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  sqlQuery: string;
  onSqlQueryChange: (value: string) => void;
  onSubmit: () => void;
  isSubmitting: boolean;
}

export function SqlEditorSheet({
  open,
  onOpenChange,
  sqlQuery,
  onSqlQueryChange,
  onSubmit,
  isSubmitting,
}: SqlEditorSheetProps) {
  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent className="flex flex-col sm:max-w-2xl bg-background border-l border-border shadow-2xl w-full">
        <SheetHeader className="pb-4 border-b border-border">
          <SheetTitle className="flex items-center gap-3 text-lg font-semibold">
            <div className="p-2 rounded-xl bg-gradient-to-br from-primary/10 to-cyan-500/10">
              <Database className="w-5 h-5 text-primary" />
            </div>
            <div>
              <span>SQL Query Editor</span>
              <p className="text-sm text-muted-foreground font-normal mt-0.5">
                Write and execute custom SQL queries
              </p>
            </div>
          </SheetTitle>
        </SheetHeader>

        <div className="flex-1 overflow-auto py-6 space-y-4">
          <div className="space-y-3">
            <Label htmlFor="sql-query" className="text-sm font-medium flex items-center gap-2">
              <Code className="w-4 h-4 text-primary" />
              SQL Query
            </Label>
            {/* Editor with faux title bar */}
            <div className="rounded-xl overflow-hidden border border-border">
              <div className="h-8 bg-muted/50 flex items-center justify-between px-3 border-b border-border">
                <div className="flex gap-1.5">
                  <div className="w-2.5 h-2.5 rounded-full bg-red-400/60" />
                  <div className="w-2.5 h-2.5 rounded-full bg-yellow-400/60" />
                  <div className="w-2.5 h-2.5 rounded-full bg-green-400/60" />
                </div>
                <span className="text-[10px] text-muted-foreground font-mono">query.sql</span>
              </div>
              <Textarea
                id="sql-query"
                value={sqlQuery}
                onChange={(e) => onSqlQueryChange(e.target.value)}
                placeholder="SELECT * FROM your_table LIMIT 10;"
                className="min-h-[280px] font-mono text-sm bg-card border-0 rounded-none resize-none focus-visible:ring-0 p-4 shadow-none"
              />
            </div>
          </div>

          <Alert className="bg-primary/5 border-primary/10 rounded-xl">
            <Sparkles className="w-4 h-4 text-primary" />
            <AlertDescription className="text-sm">
              <strong>Pro tip:</strong> Only SELECT queries are allowed for security. Use joins, WHERE clauses, and GROUP BY to refine your results.
            </AlertDescription>
          </Alert>

          <div className="space-y-2">
            <h4 className="text-sm font-semibold flex items-center gap-2">
              <Zap className="w-4 h-4 text-amber-500" />
              Example Queries
            </h4>
            <div className="space-y-2">
              {[
                "SELECT * FROM users WHERE created_at > NOW() - INTERVAL '7 days';",
                "SELECT COUNT(*), category FROM products GROUP BY category;",
              ].map((example, i) => (
                <button
                  key={i}
                  onClick={() => onSqlQueryChange(example)}
                  className="w-full text-left p-3 bg-muted/30 hover:bg-muted/50 rounded-xl border border-border/50 text-xs font-mono text-muted-foreground hover:text-foreground transition-colors"
                >
                  {example}
                </button>
              ))}
            </div>
          </div>
        </div>

        <div className="flex justify-between items-center gap-3 pt-4 border-t border-border">
          <div className="text-xs text-muted-foreground flex items-center gap-1.5">
            <kbd className="px-2 py-1 rounded-md bg-muted text-[10px] font-mono border border-border/50">Ctrl+Enter</kbd>
            <span>to execute</span>
          </div>
          <div className="flex gap-2">
            <Button variant="outline" onClick={() => onOpenChange(false)} className="min-w-[100px] rounded-xl">
              Cancel
            </Button>
            <Button
              onClick={onSubmit}
              disabled={!sqlQuery.trim() || isSubmitting}
              className="min-w-[140px] rounded-xl bg-gradient-to-r from-primary to-cyan-600 hover:from-primary/90 hover:to-cyan-600/90 shadow-sm"
            >
              {isSubmitting ? (
                <>
                  <motion.div
                    animate={{ rotate: 360 }}
                    transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
                    className="w-4 h-4 border-2 border-current border-t-transparent rounded-full mr-2"
                  />
                  Executing...
                </>
              ) : (
                <>
                  <Play className="w-4 h-4 mr-2" />
                  Execute Query
                </>
              )}
            </Button>
          </div>
        </div>
      </SheetContent>
    </Sheet>
  );
}
