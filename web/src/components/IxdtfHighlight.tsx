import { useTranslation } from "react-i18next";
import { Tooltip } from "@/components/ui/Tooltip";
import { type IxdtfTokenType, tokenize } from "@/lib/ixdtf/tokenize";
import { cn } from "@/lib/utils";

const TOKEN_CLASS: Record<IxdtfTokenType, string> = {
  date: "text-ixdtf-date",
  time: "text-ixdtf-time",
  offset: "text-ixdtf-offset",
  timezone: "text-ixdtf-timezone",
  extension: "text-ixdtf-extension",
  invalid: "text-destructive",
};

const LEGEND_TYPES: IxdtfTokenType[] = ["date", "time", "offset", "timezone", "extension"];

interface IxdtfHighlightProps {
  value: string;
  className?: string;
}

/**
 * Colour-coded IXDTF string (F-1-2 / F-2-5). Tooltips explain each component;
 * pair with {@link IxdtfLegend} so colour is never the only cue (N-3).
 */
export function IxdtfHighlight({ value, className }: IxdtfHighlightProps) {
  const { t } = useTranslation();
  const tokens = tokenize(value);
  return (
    <span className={cn("break-all font-mono tabular-nums", className)}>
      {tokens.map((token, index) => {
        const tokenClass = cn(
          TOKEN_CLASS[token.type],
          token.critical && "underline decoration-warning decoration-wavy",
        );
        if (token.type === "invalid") {
          return (
            // biome-ignore lint/suspicious/noArrayIndexKey: tokens are positional by nature
            <span key={`${index}-${token.text}`} className={tokenClass}>
              {token.text}
            </span>
          );
        }
        const label = token.critical
          ? `${t(`tooltip.${token.type}`)} — ${t("tooltip.critical")}`
          : t(`tooltip.${token.type}`);
        return (
          // biome-ignore lint/suspicious/noArrayIndexKey: tokens are positional by nature
          <Tooltip key={`${index}-${token.text}`} content={label}>
            {/* tabIndex: ツールチップは「ホバー/フォーカスで表示」(DESIGN.md N-3) —
                キーボードでも各トークンの説明へ到達できるようフォーカス可能にする */}
            <span
              // biome-ignore lint/a11y/noNoninteractiveTabindex: Radix Tooltip.Trigger (asChild) のトリガーで、フォーカスがツールチップを開く操作になる (N-3)
              tabIndex={0}
              className={cn(
                tokenClass,
                "rounded-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring",
              )}
            >
              {token.text}
            </span>
          </Tooltip>
        );
      })}
    </span>
  );
}

/** Colour legend for the highlight (N-3: label + colour, not colour alone). */
export function IxdtfLegend({ className }: { className?: string }) {
  const { t } = useTranslation();
  return (
    <div className={cn("flex flex-wrap gap-x-4 gap-y-1", className)}>
      {LEGEND_TYPES.map((type) => (
        <span key={type} className="font-sans text-muted-foreground text-xs">
          <span className={cn("font-mono", TOKEN_CLASS[type])}>■</span> {t(`legend.${type}`)}
        </span>
      ))}
      <span className="font-sans text-muted-foreground text-xs">
        <span className="font-mono text-warning underline decoration-wavy">!</span>{" "}
        {t("legend.critical")}
      </span>
    </div>
  );
}
