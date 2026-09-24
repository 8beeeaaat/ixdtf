import type { InputHTMLAttributes } from "react";
import { cn } from "@/lib/utils";

type InputProps = InputHTMLAttributes<HTMLInputElement>;

/** Transparent-background input with focus ring (DESIGN.md). */
function Input({ className, ...props }: InputProps) {
  return (
    <input
      className={cn(
        // モバイル (sm 未満) は 16px (text-base)。iOS Safari は 16px 未満の input に
        // フォーカスするとページ全体を強制ズームするため、それを防ぐ。
        "h-11 w-full rounded-md border border-border bg-transparent px-3 font-mono text-base text-foreground tabular-nums placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background disabled:opacity-50 sm:h-9 sm:text-sm",
        className,
      )}
      {...props}
    />
  );
}

export { Input };
