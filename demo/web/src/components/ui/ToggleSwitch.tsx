import * as SwitchPrimitive from "@radix-ui/react-switch";
import { cn } from "@/lib/utils";

interface ToggleSwitchProps {
  enabled: boolean;
  onChange: (enabled: boolean) => void;
  /** Accessible name; also used as the action title. */
  label: string;
  id?: string;
  className?: string;
}

/** Radix Switch with project styling (behaviour from Radix, style our own). */
function ToggleSwitch({ enabled, onChange, label, id, className }: ToggleSwitchProps) {
  return (
    <SwitchPrimitive.Root
      id={id}
      checked={enabled}
      onCheckedChange={onChange}
      aria-label={label}
      title={label}
      className={cn(
        // 見た目は 36x20 のまま、透明な擬似要素でタップ可能領域を約 52x44 に拡大する
        "relative h-5 w-9 shrink-0 rounded-full border border-border bg-muted transition-colors after:absolute after:-inset-x-2 after:-inset-y-3 after:content-[''] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background data-[state=checked]:bg-primary",
        className,
      )}
    >
      <SwitchPrimitive.Thumb className="block h-4 w-4 translate-x-0.5 rounded-full bg-card transition-transform data-[state=checked]:translate-x-4" />
    </SwitchPrimitive.Root>
  );
}

export { ToggleSwitch };
