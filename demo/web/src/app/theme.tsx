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

    // モバイルブラウザの UI 色 (theme-color) を解決済みテーマへ追従させる。
    // 手動切替 (.light/.dark) は index.html の media 付き meta では追従できないため、
    // 地球儀 canvas と同じく getComputedStyle でトークン値を読む (DESIGN.md: 色はトークン駆動)。
    const media = window.matchMedia("(prefers-color-scheme: dark)");
    const syncThemeColor = () => {
      const color = getComputedStyle(document.body).backgroundColor;
      for (const meta of document.querySelectorAll<HTMLMetaElement>('meta[name="theme-color"]')) {
        meta.content = color;
      }
    };
    syncThemeColor();
    media.addEventListener("change", syncThemeColor);
    return () => media.removeEventListener("change", syncThemeColor);
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
