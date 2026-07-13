import { useCallback, useEffect, useState } from "react";

/** State persisted to localStorage (F-4-2 world-clock timezone list). */
export function useLocalStorage<T>(key: string, initial: T): [T, (value: T) => void] {
  const [value, setValue] = useState<T>(() => {
    try {
      const raw = localStorage.getItem(key);
      return raw ? (JSON.parse(raw) as T) : initial;
    } catch {
      return initial;
    }
  });

  useEffect(() => {
    try {
      localStorage.setItem(key, JSON.stringify(value));
    } catch {
      // Ignore quota / private-mode failures — persistence is best-effort.
    }
  }, [key, value]);

  const set = useCallback((next: T) => setValue(next), []);

  return [value, set];
}
