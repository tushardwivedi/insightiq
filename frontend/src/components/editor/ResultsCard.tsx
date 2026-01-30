import { Badge } from "@/components/ui/badge";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { AnalyticsResponse, VoiceResponse } from "@/types";
import { Wand2, Table, Clock } from "lucide-react";

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
    <div className="space-y-4 pt-2">
      {insights && (
        <Alert className="bg-[#10a37f]/5 border-[#10a37f]/20 dark:bg-[#10a37f]/10">
          <Wand2 className="w-4 h-4 text-[#10a37f]" />
          <AlertDescription className="text-sm text-foreground">{insights}</AlertDescription>
        </Alert>
      )}

      <div className="flex items-center gap-4 text-xs text-muted-foreground">
        <Badge variant="secondary" className="flex items-center gap-1 font-normal">
          <Table className="w-3 h-3" />
          {recordCount} records
        </Badge>
        {processTime && (
          <Badge variant="secondary" className="flex items-center gap-1 font-normal">
            <Clock className="w-3 h-3" />
            {processTime}
          </Badge>
        )}
      </div>

      {data.length > 0 && (
        <div className="rounded-xl border-2 border-border overflow-hidden bg-card shadow-sm">
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead className="bg-muted/50 border-b border-border">
                <tr>
                  {Object.keys(data[0]).map((key) => (
                    <th
                      key={key}
                      className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wider text-foreground"
                    >
                      {key.replace(/_/g, " ")}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {data.slice(0, 5).map((row, i) => (
                  <tr key={i} className="hover:bg-muted/30 transition-colors">
                    {Object.values(row).map((value: any, j) => (
                      <td key={j} className="px-4 py-3 text-foreground">
                        {typeof value === "number" ? value.toLocaleString() : String(value)}
                      </td>
                    ))}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          {data.length > 5 && (
            <div className="px-4 py-2.5 bg-muted/50 text-xs text-center text-muted-foreground border-t border-border font-medium">
              Showing 5 of {data.length} records • <button className="text-primary hover:underline">View all</button>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
