import type { ReactNode } from "react";
import { cn } from "@/lib/utils";

interface FormFieldProps {
  label: string;
  htmlFor?: string;
  hint?: string;
  error?: string;
  children: ReactNode;
  className?: string;
}

/** Label + control + hint/error wrapper (TanStack Form friendly). */
function FormField({ label, htmlFor, hint, error, children, className }: FormFieldProps) {
  return (
    <div className={cn("flex flex-col gap-1.5", className)}>
      <label htmlFor={htmlFor} className="font-medium font-sans text-muted-foreground text-xs">
        {label}
      </label>
      {children}
      {error ? (
        <p className="font-sans text-destructive text-xs">{error}</p>
      ) : hint ? (
        <p className="font-sans text-muted-foreground text-xs">{hint}</p>
      ) : null}
    </div>
  );
}

export { FormField };
