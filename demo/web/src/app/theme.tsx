import { createContext, type ReactNode, useCallback, useContext, useEffect, useState } from "react";

export type ThemeMode = "auto" | "light" | "dark";

interface ThemeContextValue {
  mode: ThemeMode;
  setMode: (mode: ThemeMode) => void;
}

const ThemeContext = createContext<ThemeContextValue | null>(null);
const STORAGE_KEY = "theme";

function readMode(): ThemeMode {
  const stored = localStorage.getItem(STORAGE_KEY);
  return stored === "light" || stored === "dark" || stored === "auto" ? stored : "auto";
}

/**
 * Theme provider (F-0-3). `auto` follows `prefers-color-scheme`; `light`/`dark`
 * force it by stamping a class on <html>. Tokens switch via CSS (no `dark:`).
 */
export function ThemeProvider({ children }: { children: ReactNode }) {
  const [mode, setModeState] = useState<ThemeMode>(readMode);

  useEffect(() => {
    const root = document.documentElement;
    root.classList.remove("light", "dark");
    if (mode === "light") {
      root.classList.add("light");
    } else if (mode === "dark") {
      root.classList.add("dark");
    }
    localStorage.setItem(STORAGE_KEY, mode);
  }, [mode]);

  const setMode = useCallback((next: ThemeMode) => setModeState(next), []);

  return <ThemeContext.Provider value={{ mode, setMode }}>{children}</ThemeContext.Provider>;
}

export function useTheme(): ThemeContextValue {
  const context = useContext(ThemeContext);
  if (!context) {
    throw new Error("useTheme must be used within ThemeProvider");
  }
  return context;
}
