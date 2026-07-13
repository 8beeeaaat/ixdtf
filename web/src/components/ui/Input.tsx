import type { InputHTMLAttributes } from "react";
import { cn } from "@/lib/utils";

type InputProps = InputHTMLAttributes<HTMLInputElement>;

/** Transparent-background input with focus ring (DESIGN.md). */
function Input({ className, ...props }: InputProps) {
  return (
    <input
      className={cn(
        "h-9 w-full rounded-md border border-border bg-transparent px-3 font-mono text-foreground text-sm tabular-nums placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background disabled:opacity-50",
        className,
      )}
      {...props}
    />
  );
}

export { Input };
