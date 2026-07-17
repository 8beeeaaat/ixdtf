import * as SelectPrimitive from "@radix-ui/react-select";
import type { ReactNode } from "react";
import { cn } from "@/lib/utils";

export interface SelectOption {
  value: string;
  label: string;
}

interface SelectProps {
  value: string;
  onValueChange: (value: string) => void;
  options: SelectOption[];
  "aria-label": string;
  /** Optional leading glyph reflecting the current value. */
  icon?: ReactNode;
  className?: string;
}

function ChevronDown() {
  return (
    <svg
      viewBox="0 0 16 16"
      aria-hidden="true"
      fill="none"
      stroke="currentColor"
      strokeWidth={1.3}
      strokeLinecap="round"
      strokeLinejoin="round"
      className="h-3.5 w-3.5 shrink-0 text-muted-foreground"
    >
      <path d="m4 6 4 4 4-4" />
    </svg>
  );
}

function Check() {
  return (
    <svg
      viewBox="0 0 16 16"
      aria-hidden="true"
      fill="none"
      stroke="currentColor"
      strokeWidth={1.5}
      strokeLinecap="round"
      strokeLinejoin="round"
      className="h-3.5 w-3.5 shrink-0"
    >
      <path d="m3.5 8.5 3 3 6-7" />
    </svg>
  );
}

/**
 * Compact dropdown that shows the CURRENT value in the header settings row
 * (Radix Select for behaviour; visual style is ours per DESIGN.md). Unlike a
 * cycle button that displays the next state, this keeps every header control
 * consistent: the trigger always reflects the selected value.
 */
function Select({
  value,
  onValueChange,
  options,
  "aria-label": ariaLabel,
  icon,
  className,
}: SelectProps) {
  const current = options.find((option) => option.value === value);
  return (
    <SelectPrimitive.Root value={value} onValueChange={onValueChange}>
      <SelectPrimitive.Trigger
        aria-label={ariaLabel}
        className={cn(
          "inline-flex h-8 shrink-0 items-center gap-1.5 whitespace-nowrap rounded-md px-3 font-medium font-sans text-foreground text-xs transition-colors hover:bg-accent hover:text-accent-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background",
          className,
        )}
      >
        {icon}
        <span className={cn(icon && "hidden sm:inline")}>{current?.label ?? value}</span>
        <SelectPrimitive.Icon asChild>
          <ChevronDown />
        </SelectPrimitive.Icon>
      </SelectPrimitive.Trigger>
      <SelectPrimitive.Portal>
        <SelectPrimitive.Content
          position="popper"
          sideOffset={4}
          align="end"
          className="z-50 min-w-[8rem] overflow-hidden rounded-md border border-border bg-card py-1 text-card-foreground"
        >
          <SelectPrimitive.Viewport>
            {options.map((option) => (
              <SelectPrimitive.Item
                key={option.value}
                value={option.value}
                className="relative flex cursor-default items-center py-1.5 pr-3 pl-7 font-sans text-sm outline-none data-[highlighted]:bg-accent data-[highlighted]:text-accent-foreground"
              >
                <SelectPrimitive.ItemIndicator className="absolute left-2 inline-flex items-center">
                  <Check />
                </SelectPrimitive.ItemIndicator>
                <SelectPrimitive.ItemText>{option.label}</SelectPrimitive.ItemText>
              </SelectPrimitive.Item>
            ))}
          </SelectPrimitive.Viewport>
        </SelectPrimitive.Content>
      </SelectPrimitive.Portal>
    </SelectPrimitive.Root>
  );
}

export { Select };
