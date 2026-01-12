"use client";

import React, { createContext, useContext, useEffect, useState } from "react";

type Theme = "dark" | "light" | "system";

type AccentColor = 
  | "cyan" 
  | "indigo" 
  | "purple" 
  | "pink" 
  | "orange";

interface ThemeContextType {
  theme: Theme;
  accentColor: AccentColor;
  setTheme: (theme: Theme) => void;
  setAccentColor: (color: AccentColor) => void;
  resolvedTheme: "dark" | "light";
}

const ThemeContext = createContext<ThemeContextType | undefined>(undefined);

const accentColors: Record<AccentColor, { primary: string; from: string; to: string }> = {
  cyan: { primary: "#4fd1c5", from: "#4fd1c5", to: "#38b2ac" },
  indigo: { primary: "#6366f1", from: "#6366f1", to: "#4f46e5" },
  purple: { primary: "#8b5cf6", from: "#8b5cf6", to: "#7c3aed" },
  pink: { primary: "#ec4899", from: "#ec4899", to: "#db2777" },
  orange: { primary: "#f59e0b", from: "#f59e0b", to: "#d97706" },
};

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const [theme, setTheme] = useState<Theme>("dark");
  const [accentColor, setAccentColor] = useState<AccentColor>("cyan");
  const [resolvedTheme, setResolvedTheme] = useState<"dark" | "light">("dark");

  useEffect(() => {
    const storedTheme = localStorage.getItem("theme") as Theme | null;
    const storedAccentColor = localStorage.getItem("accentColor") as AccentColor | null;
    
    if (storedTheme) {
      setTheme(storedTheme);
    }
    if (storedAccentColor) {
      setAccentColor(storedAccentColor);
    }
  }, []);

  useEffect(() => {
    const root = document.documentElement;
    
    const systemTheme = window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light";
    const currentTheme = theme === "system" ? systemTheme : theme;
    
    setResolvedTheme(currentTheme);
    
    root.classList.remove("dark", "light");
    root.classList.add(currentTheme);
    
    localStorage.setItem("theme", theme);
  }, [theme]);

  useEffect(() => {
    const root = document.documentElement;
    const colors = accentColors[accentColor];
    
    root.style.setProperty("--accent-color", colors.primary);
    root.style.setProperty("--accent-from", colors.from);
    root.style.setProperty("--accent-to", colors.to);
    
    root.style.setProperty("--primary", colors.primary);
    root.style.setProperty("--primary-foreground", "#ffffff");
    
    localStorage.setItem("accentColor", accentColor);
  }, [accentColor]);

  useEffect(() => {
    const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
    const handleChange = () => {
      if (theme === "system") {
        const newTheme = mediaQuery.matches ? "dark" : "light";
        setResolvedTheme(newTheme);
        document.documentElement.classList.remove("dark", "light");
        document.documentElement.classList.add(newTheme);
      }
    };

    mediaQuery.addEventListener("change", handleChange);
    return () => mediaQuery.removeEventListener("change", handleChange);
  }, [theme]);

  return (
    <ThemeContext.Provider
      value={{
        theme,
        accentColor,
        setTheme,
        setAccentColor,
        resolvedTheme,
      }}
    >
      {children}
    </ThemeContext.Provider>
  );
}

export function useTheme() {
  const context = useContext(ThemeContext);
  if (context === undefined) {
    throw new Error("useTheme must be used within a ThemeProvider");
  }
  return context;
}

export { accentColors };
