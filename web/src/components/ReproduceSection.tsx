import { useState } from "react";
import { useTranslation } from "react-i18next";
import { CopyButton } from "@/components/CopyButton";
import { cn } from "@/lib/utils";

const TABS = ["go", "js"] as const;
type Tab = (typeof TABS)[number];

interface ReproduceSectionProps {
  goSample: string;
  jsSample: string;
  className?: string;
}

/**
 * F-0-7: collapsible "reproduce this result" section shown under API-backed
 * results — Go (ixdtf) and JavaScript (Temporal) code samples with the
 * current inputs embedded.
 */
export function ReproduceSection({ goSample, jsSample, className }: ReproduceSectionProps) {
  const { t } = useTranslation();
  const [tab, setTab] = useState<Tab>("go");
  const samples: Record<Tab, string> = { go: goSample, js: jsSample };

  return (
    <details className={cn("group", className)}>
      <summary className="cursor-pointer list-none rounded-md py-2.5 font-sans text-muted-foreground text-sm transition-colors hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background sm:py-0 [&::-webkit-details-marker]:hidden">
        <span aria-hidden className="mr-1 inline-block transition-transform group-open:rotate-90">
          ▸
        </span>
        {t("reproduce.summary")}
      </summary>
      <div className="mt-3 space-y-3">
        <div className="flex flex-wrap items-center gap-1">
          {TABS.map((key) => (
            <button
              key={key}
              type="button"
              aria-pressed={tab === key}
              onClick={() => setTab(key)}
              className={cn(
                // モバイルはタップターゲットを 44px 相当へ (py-3.5)。デスクトップは従来の密度
                "whitespace-nowrap rounded-md px-3 py-3.5 font-sans text-xs transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background sm:px-2 sm:py-1",
                tab === key
                  ? "bg-muted font-medium text-foreground"
                  : "text-muted-foreground hover:text-foreground",
              )}
            >
              {t(`reproduce.tabs.${key}`)}
            </button>
          ))}
          <CopyButton text={samples[tab]} className="ml-auto" />
        </div>
        <pre className="overflow-x-auto rounded-md bg-muted p-3 font-mono text-xs tabular-nums leading-relaxed">
          <code>{samples[tab]}</code>
        </pre>
      </div>
    </details>
  );
}
