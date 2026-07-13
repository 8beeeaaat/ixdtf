import { type DialogHTMLAttributes, type ReactNode, useEffect, useId, useRef } from "react";
import { Button } from "@/components/ui/Button";
import { cn } from "@/lib/utils";

interface DialogProps extends Omit<DialogHTMLAttributes<HTMLDialogElement>, "open" | "title"> {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title: ReactNode;
  closeLabel: string;
  children: ReactNode;
}

/** Native modal dialog with browser-managed focus trapping and Escape handling. */
function Dialog({
  open,
  onOpenChange,
  title,
  closeLabel,
  children,
  className,
  ...props
}: DialogProps) {
  const dialogRef = useRef<HTMLDialogElement>(null);
  const titleId = useId();

  useEffect(() => {
    const dialog = dialogRef.current;
    if (!dialog) {
      return;
    }
    if (open && !dialog.open) {
      dialog.showModal();
    } else if (!open && dialog.open) {
      dialog.close();
    }
  }, [open]);

  return (
    <dialog
      ref={dialogRef}
      aria-labelledby={titleId}
      className={cn(
        "m-auto max-h-[calc(100dvh-3rem)] w-[calc(100%-3rem)] max-w-2xl overflow-hidden rounded-lg border border-border bg-card p-0 text-card-foreground backdrop:bg-background/80",
        className,
      )}
      onCancel={() => onOpenChange(false)}
      onKeyDown={(event) => {
        if (event.key === "Escape") {
          event.preventDefault();
          onOpenChange(false);
        }
      }}
      {...props}
    >
      <div className="flex items-start justify-between gap-6 border-border border-b px-6 py-5">
        <h2 id={titleId} className="font-sans font-semibold text-lg tracking-tight">
          {title}
        </h2>
        <Button
          variant="ghost"
          size="sm"
          className="-mr-2 shrink-0 px-2 text-muted-foreground text-xl leading-none"
          aria-label={closeLabel}
          onClick={() => onOpenChange(false)}
        >
          ×
        </Button>
      </div>
      <div className="max-h-[calc(100dvh-8rem)] overflow-y-auto px-6 py-5">{children}</div>
    </dialog>
  );
}

export { Dialog };
