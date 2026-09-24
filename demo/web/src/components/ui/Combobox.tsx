import { useEffect, useId, useMemo, useRef, useState } from "react";
import { cn } from "@/lib/utils";

interface ComboboxProps {
  value: string;
  options: string[];
  onChange: (value: string) => void;
  placeholder?: string;
  "aria-label"?: string;
  /** FormField の label (htmlFor) と関連付けるための input id。 */
  id?: string;
  className?: string;
}

/**
 * Searchable single-select (F-1-3 / F-1-4 pickers). Keyboard: ArrowUp/Down,
 * Enter to commit, Escape to close.
 */
function Combobox({
  value,
  options,
  onChange,
  placeholder,
  "aria-label": ariaLabel,
  id,
  className,
}: ComboboxProps) {
  const [open, setOpen] = useState(false);
  const [query, setQuery] = useState("");
  const [active, setActive] = useState(0);
  const rootRef = useRef<HTMLDivElement>(null);
  const listId = useId();

  const filtered = useMemo(() => {
    const q = query.trim().toLowerCase();
    const hits = q ? options.filter((o) => o.toLowerCase().includes(q)) : options;
    return hits.slice(0, 50);
  }, [options, query]);

  useEffect(() => {
    if (!open) {
      return;
    }
    const onPointerDown = (event: PointerEvent) => {
      if (rootRef.current && !rootRef.current.contains(event.target as Node)) {
        setOpen(false);
      }
    };
    document.addEventListener("pointerdown", onPointerDown);
    return () => document.removeEventListener("pointerdown", onPointerDown);
  }, [open]);

  const commit = (next: string) => {
    onChange(next);
    setQuery("");
    setOpen(false);
  };

  return (
    <div ref={rootRef} className={cn("relative", className)}>
      <input
        id={id}
        role="combobox"
        aria-expanded={open}
        aria-controls={listId}
        aria-autocomplete="list"
        aria-activedescendant={open && filtered[active] ? `${listId}-option-${active}` : undefined}
        aria-label={ariaLabel}
        placeholder={placeholder}
        value={open ? query : value}
        onFocus={() => {
          setOpen(true);
          setQuery("");
          setActive(0);
        }}
        onChange={(event) => {
          setQuery(event.target.value);
          setOpen(true);
          setActive(0);
        }}
        onKeyDown={(event) => {
          if (event.key === "ArrowDown") {
            event.preventDefault();
            setOpen(true);
            setActive((i) => Math.min(i + 1, filtered.length - 1));
          } else if (event.key === "ArrowUp") {
            event.preventDefault();
            setActive((i) => Math.max(i - 1, 0));
          } else if (event.key === "Enter" && open && filtered[active]) {
            event.preventDefault();
            commit(filtered[active]);
          } else if (event.key === "Escape") {
            setOpen(false);
          }
        }}
        className={cn(
          // モバイル (sm 未満) は 16px/44px 高。iOS Safari のフォーカス時強制ズームを防ぐ (Input と同基準)
          "h-11 w-full rounded-md border border-border bg-transparent px-3 font-mono text-base text-foreground tabular-nums placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background sm:h-9 sm:text-sm",
        )}
      />
      {open && filtered.length > 0 && (
        <div
          id={listId}
          role="listbox"
          className="absolute z-40 mt-1 max-h-64 w-full overflow-y-auto rounded-md border border-border bg-card py-1"
        >
          {filtered.map((option, index) => (
            <button
              key={option}
              id={`${listId}-option-${index}`}
              type="button"
              role="option"
              aria-selected={option === value}
              onMouseEnter={() => setActive(index)}
              onClick={() => commit(option)}
              className={cn(
                "block w-full px-3 py-2.5 text-left font-mono text-card-foreground text-sm sm:py-1.5",
                index === active && "bg-accent text-accent-foreground",
              )}
            >
              {option}
            </button>
          ))}
        </div>
      )}
    </div>
  );
}

export { Combobox };
