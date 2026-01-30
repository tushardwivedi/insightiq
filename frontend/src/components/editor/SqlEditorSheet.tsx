"use client";

import { motion } from "framer-motion";
import { Sheet, SheetContent, SheetHeader, SheetTitle } from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Card, CardContent } from "@/components/ui/card";
import { Database, Code, Sparkles, Zap } from "lucide-react";

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
      <SheetContent className="flex flex-col sm:max-w-2xl bg-white dark:bg-[#1e1e1e] border-l border-border shadow-2xl w-full">
        <SheetHeader className="pb-4 border-b border-border">
          <SheetTitle className="flex items-center gap-2 text-lg font-semibold">
            <div className="p-2 rounded-lg bg-primary/10">
              <Database className="w-5 h-5 text-primary" />
            </div>
            SQL Query Editor
          </SheetTitle>
          <p className="text-sm text-muted-foreground mt-2">
            Write and execute custom SQL queries directly
          </p>
        </SheetHeader>

        <div className="flex-1 overflow-auto py-6 space-y-4">
          <div className="space-y-3">
            <Label htmlFor="sql-query" className="text-sm font-medium flex items-center gap-2">
              <Code className="w-4 h-4 text-primary" />
              SQL Query
            </Label>
            <Textarea
              id="sql-query"
              value={sqlQuery}
              onChange={(e) => onSqlQueryChange(e.target.value)}
              placeholder="SELECT * FROM your_table LIMIT 10;"
              className="min-h-[300px] font-mono text-sm bg-card border-2 border-border rounded-xl resize-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:border-primary p-4 shadow-sm"
            />
          </div>

          <Alert className="bg-blue-50 dark:bg-blue-950/30 border-blue-200 dark:border-blue-900">
            <Sparkles className="w-4 h-4 text-blue-600 dark:text-blue-400" />
            <AlertDescription className="text-sm text-blue-900 dark:text-blue-100">
              <strong>Pro tip:</strong> Only SELECT queries are allowed for security. Use joins, WHERE clauses, and GROUP BY to refine your results.
            </AlertDescription>
          </Alert>

          <Card className="bg-muted/50">
            <CardContent className="p-4">
              <h4 className="text-sm font-semibold mb-2 flex items-center gap-2">
                <Zap className="w-4 h-4 text-amber-500" />
                Example Queries
              </h4>
              <div className="space-y-2 text-xs font-mono">
                <div className="p-2 bg-background rounded border">SELECT * FROM users WHERE created_at &gt; NOW() - INTERVAL '7 days';</div>
                <div className="p-2 bg-background rounded border">SELECT COUNT(*), category FROM products GROUP BY category;</div>
              </div>
            </CardContent>
          </Card>
        </div>

        <div className="flex justify-between items-center gap-3 pt-4 border-t border-border bg-white dark:bg-[#1e1e1e]">
          <div className="text-xs text-muted-foreground flex items-center gap-1">
            <kbd className="px-2 py-1 rounded bg-muted text-[10px] font-mono">Ctrl+Enter</kbd>
            <span>to execute</span>
          </div>
          <div className="flex gap-2">
            <Button variant="outline" onClick={() => onOpenChange(false)} className="min-w-[100px]">
              Cancel
            </Button>
            <Button
              onClick={onSubmit}
              disabled={!sqlQuery.trim() || isSubmitting}
              className="min-w-[140px] bg-[#10a37f] hover:bg-[#0d8c6f]"
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
                  <Zap className="w-4 h-4 mr-2" />
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
