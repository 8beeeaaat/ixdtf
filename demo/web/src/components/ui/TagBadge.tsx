import { cva, type VariantProps } from "class-variance-authority";
import type { HTMLAttributes } from "react";
import { cn } from "@/lib/utils";

const tagBadgeVariants = cva("rounded-full px-2 py-0.5 font-medium font-mono text-xs", {
  variants: {
    variant: {
      default: "bg-muted text-muted-foreground",
      critical: "bg-warning/10 text-warning underline decoration-wavy",
      match: "bg-success/10 text-success",
      mismatch: "bg-destructive/10 text-destructive",
    },
  },
  defaultVariants: { variant: "default" },
});

type TagBadgeProps = HTMLAttributes<HTMLSpanElement> & VariantProps<typeof tagBadgeVariants>;

function TagBadge({ className, variant, ...props }: TagBadgeProps) {
  return <span className={cn(tagBadgeVariants({ variant, className }))} {...props} />;
}

export { TagBadge, tagBadgeVariants };
