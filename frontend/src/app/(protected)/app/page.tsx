"use client";

import { useState } from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import CommandBar from "@/components/CommandBar";
import SQLQuerySection from "@/components/SQLQuerySection";
import ResultsSection from "@/components/ResultsSection";
import WelcomeHero from "@/components/WelcomeHero";

export default function AnalyticsPage() {
  const [results, setResults] = useState<any>(null);
  const [loading, setLoading] = useState(false);

  return (
    <div className="flex flex-col h-full">
      <div className="flex-1 overflow-auto">
        <div className="container mx-auto max-w-7xl space-y-6">
          {/* Welcome Hero Section */}
          <WelcomeHero />

          {/* Main Workspace */}
          <div className="space-y-6">
            {/* Query Editor */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <span>Query Editor</span>
                  <Badge variant="secondary">SQL</Badge>
                </CardTitle>
                <CardDescription>
                  Write and execute SQL queries to analyze your data
                </CardDescription>
              </CardHeader>
              <CardContent>
                <SQLQuerySection onResult={setResults} onLoading={setLoading} />
              </CardContent>
            </Card>

            {/* Results Panel */}
            {(results || loading) && (
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <span>Analysis Results</span>
                    {loading && <Badge variant="outline">Processing</Badge>}
                  </CardTitle>
                  <CardDescription>
                    AI-powered insights and visualizations from your query
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <ResultsSection results={results} loading={loading} />
                </CardContent>
              </Card>
            )}
          </div>
        </div>
      </div>

      {/* Command Bar - Fixed at bottom */}
      <div className="border-t bg-card/95 backdrop-blur supports-[backdrop-filter]:bg-card/60 p-4">
        <div className="container mx-auto max-w-7xl">
          <CommandBar onResult={setResults} onLoading={setLoading} />
        </div>
      </div>
    </div>
  );
}
