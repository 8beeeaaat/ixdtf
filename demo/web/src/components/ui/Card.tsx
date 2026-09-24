import type { HTMLAttributes } from "react";
import { cn } from "@/lib/utils";

type DivProps = HTMLAttributes<HTMLDivElement>;

/** Flat card (DESIGN.md: hairline border, no shadow). */
function Card({ className, ...props }: DivProps) {
  return (
    <div
      className={cn("rounded-lg border border-border bg-card text-card-foreground", className)}
      {...props}
    />
  );
}

function CardHeader({ className, ...props }: DivProps) {
  return <div className={cn("flex flex-col gap-1 px-6 pt-5", className)} {...props} />;
}

function CardTitle({ className, ...props }: HTMLAttributes<HTMLHeadingElement>) {
  return (
    <h2
      className={cn(
        "font-sans font-semibold text-muted-foreground text-sm uppercase tracking-wider",
        className,
      )}
      {...props}
    />
  );
}

function CardContent({ className, ...props }: DivProps) {
  return <div className={cn("px-6 py-5", className)} {...props} />;
}

export { Card, CardContent, CardHeader, CardTitle };
