"use client";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Moon, Sun, Monitor } from "lucide-react";

type Theme = "dark" | "light" | "system";
type AccentColor = "cyan" | "indigo" | "purple" | "pink" | "orange";

interface AppearanceSettingsProps {
  theme: Theme;
  accentColor: AccentColor;
  onThemeChange: (theme: Theme) => void;
  onAccentColorChange: (color: AccentColor) => void;
}

const accentColors = [
  { name: "cyan", color: "#4fd1c5" },
  { name: "indigo", color: "#6366f1" },
  { name: "purple", color: "#8b5cf6" },
  { name: "pink", color: "#ec4899" },
  { name: "orange", color: "#f59e0b" },
];

export function AppearanceSettings({ theme, accentColor, onThemeChange, onAccentColorChange }: AppearanceSettingsProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Appearance Settings</CardTitle>
        <CardDescription>Customize how InsightIQ looks</CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="space-y-2">
          <Label>Theme</Label>
          <div className="grid grid-cols-3 gap-4">
            <div
              className={`p-4 rounded-lg border-2 cursor-pointer transition-all ${
                theme === "dark" ? "border-primary bg-primary/5" : "border-border hover:border-primary/50"
              }`}
              onClick={() => onThemeChange("dark")}
            >
              <div className="h-16 rounded bg-slate-900 mb-2 flex items-center justify-center">
                <Moon className="w-6 h-6 text-slate-400" />
              </div>
              <p className="font-medium text-sm text-center">Dark</p>
            </div>
            <div
              className={`p-4 rounded-lg border-2 cursor-pointer transition-all ${
                theme === "light" ? "border-primary bg-primary/5" : "border-border hover:border-primary/50"
              }`}
              onClick={() => onThemeChange("light")}
            >
              <div className="h-16 rounded bg-white border mb-2 flex items-center justify-center">
                <Sun className="w-6 h-6 text-amber-500" />
              </div>
              <p className="font-medium text-sm text-center">Light</p>
            </div>
            <div
              className={`p-4 rounded-lg border-2 cursor-pointer transition-all ${
                theme === "system" ? "border-primary bg-primary/5" : "border-border hover:border-primary/50"
              }`}
              onClick={() => onThemeChange("system")}
            >
              <div className="h-16 rounded bg-gradient-to-br from-slate-900 to-white mb-2 flex items-center justify-center">
                <Monitor className="w-6 h-6 text-slate-600" />
              </div>
              <p className="font-medium text-sm text-center">System</p>
            </div>
          </div>
        </div>

        <div className="space-y-2">
          <Label>Accent Color</Label>
          <div className="flex gap-3">
            {accentColors.map((item) => (
              <button
                key={item.name}
                onClick={() => onAccentColorChange(item.name as AccentColor)}
                className={`w-12 h-12 rounded-full border-2 transition-all ${
                  accentColor === item.name
                    ? "border-white scale-110 shadow-lg"
                    : "border-transparent hover:scale-105"
                }`}
                style={{ backgroundColor: item.color }}
                title={item.name.charAt(0).toUpperCase() + item.name.slice(1)}
              />
            ))}
          </div>
          <p className="text-sm text-muted-foreground">
            Current: <span className="capitalize">{accentColor}</span>
          </p>
        </div>
      </CardContent>
    </Card>
  );
}
