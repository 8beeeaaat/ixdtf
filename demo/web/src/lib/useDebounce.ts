import { useEffect, useState } from "react";

/** Debounce a value by `delayMs` (F-2-1 real-time-but-throttled parsing). */
export function useDebounce<T>(value: T, delayMs: number): T {
  const [debounced, setDebounced] = useState(value);

  useEffect(() => {
    const id = setTimeout(() => setDebounced(value), delayMs);
    return () => clearTimeout(id);
  }, [value, delayMs]);

  return debounced;
}
