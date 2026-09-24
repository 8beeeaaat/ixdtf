import { useEffect, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { Button } from "@/components/ui/Button";

interface CopyButtonProps {
  /** Text to copy, or a producer evaluated at click time. */
  text: string | (() => string);
  label?: string;
  copiedLabel?: string;
  className?: string;
}

/** One-click clipboard copy with transient feedback (F-1-6). */
export function CopyButton({ text, label, copiedLabel, className }: CopyButtonProps) {
  const { t } = useTranslation();
  const [copied, setCopied] = useState(false);
  const timer = useRef<ReturnType<typeof setTimeout>>(undefined);

  useEffect(() => () => clearTimeout(timer.current), []);

  const onCopy = async () => {
    const value = typeof text === "function" ? text() : text;
    try {
      await navigator.clipboard.writeText(value);
      setCopied(true);
      clearTimeout(timer.current);
      timer.current = setTimeout(() => setCopied(false), 2000);
    } catch {
      // Clipboard permission denied — leave the label unchanged.
    }
  };

  return (
    <Button variant="secondary" size="sm" onClick={onCopy} className={className}>
      {copied ? (copiedLabel ?? t("common.copied")) : (label ?? t("common.copy"))}
    </Button>
  );
}
