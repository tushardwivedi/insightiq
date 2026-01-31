import { Badge } from "@/components/ui/badge";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { AnalyticsResponse, VoiceResponse } from "@/types";
import { Sparkles, Table, Clock, ChevronRight } from "lucide-react";

interface ResultsCardProps {
  results: AnalyticsResponse | VoiceResponse;
}

export function ResultsCard({ results }: ResultsCardProps) {
  const data = "response" in results ? (results as VoiceResponse).response.data : (results as AnalyticsResponse).data;
  const insights = "response" in results ? (results as VoiceResponse).response.insights : (results as AnalyticsResponse).insights;
  const processTime = "response" in results ? (results as VoiceResponse).response.process_time : (results as AnalyticsResponse).process_time;
  const recordCount = data?.length || 0;

  if (!data || data.length === 0) return null;

  return (
    <div className="space-y-3 pt-2">
      {/* Insights banner */}
      {insights && (
        <Alert className="bg-gradient-to-r from-emerald-500/5 to-teal-500/5 border-emerald-500/15 rounded-xl">
          <Sparkles className="w-4 h-4 text-emerald-500" />
          <AlertDescription className="text-sm text-foreground/90 leading-relaxed">
            {insights}
          </AlertDescription>
        </Alert>
      )}

      {/* Metadata badges */}
      <div className="flex items-center gap-2 text-xs">
        <Badge variant="secondary" className="flex items-center gap-1.5 font-normal rounded-full bg-muted/50 text-muted-foreground px-3 py-1">
          <Table className="w-3 h-3" />
          {recordCount} records
        </Badge>
        {processTime && (
          <Badge variant="secondary" className="flex items-center gap-1.5 font-normal rounded-full bg-muted/50 text-muted-foreground px-3 py-1">
            <Clock className="w-3 h-3" />
            {processTime}
          </Badge>
        )}
      </div>

      {/* Data table */}
      {data.length > 0 && (
        <div className="rounded-xl border border-border overflow-hidden bg-card shadow-sm">
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border bg-muted/30">
                  {Object.keys(data[0]).map((key) => (
                    <th
                      key={key}
                      className="px-4 py-2.5 text-left text-[11px] font-semibold uppercase tracking-wider text-muted-foreground"
                    >
                      {key.replace(/_/g, " ")}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {data.slice(0, 5).map((row, i) => (
                  <tr
                    key={i}
                    className="border-b border-border/50 last:border-0 hover:bg-muted/20 transition-colors"
                  >
                    {Object.values(row).map((value: any, j) => (
                      <td
                        key={j}
                        className="px-4 py-2.5 text-foreground/90 text-[13px] font-mono"
                      >
                        {typeof value === "number" ? value.toLocaleString() : String(value)}
                      </td>
                    ))}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          {data.length > 5 && (
            <div className="px-4 py-2 bg-muted/20 text-xs text-center border-t border-border/50">
              <button className="inline-flex items-center gap-1 text-primary hover:underline font-medium">
                View all {data.length} records
                <ChevronRight className="w-3 h-3" />
              </button>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
