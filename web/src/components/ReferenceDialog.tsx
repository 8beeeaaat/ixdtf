import { useId, useState } from "react";
import { useTranslation } from "react-i18next";
import { Button } from "@/components/ui/Button";
import { Dialog } from "@/components/ui/Dialog";
import { TagBadge } from "@/components/ui/TagBadge";
import { type ReferenceId, references } from "@/lib/references";
import { cn } from "@/lib/utils";

interface ReferenceDialogProps {
  referenceId: ReferenceId;
  className?: string;
}

/** Contextual explanation and primary-source links for IXDTF concepts (F-0-6). */
function ReferenceDialog({ referenceId, className }: ReferenceDialogProps) {
  const { t } = useTranslation();
  const [open, setOpen] = useState(false);
  const sourcesHeadingId = useId();
  const reference = references[referenceId];

  return (
    <>
      <Button
        variant="ghost"
        size="sm"
        className={cn(
          // py-3.5 + 負マージンで、レイアウトを動かさずタップ可能領域だけ 44px 相当に広げる
          "-my-2.5 h-auto px-2 py-3.5 font-medium text-info text-xs normal-case tracking-normal sm:my-0 sm:py-1",
          className,
        )}
        aria-haspopup="dialog"
        onClick={() => setOpen(true)}
      >
        <span aria-hidden="true">ⓘ</span>
        {t("references.trigger")}
      </Button>
      <Dialog
        open={open}
        onOpenChange={setOpen}
        title={t(reference.titleKey)}
        closeLabel={t("references.close")}
      >
        <div className="space-y-6">
          <p className="font-sans text-sm leading-relaxed">{t(reference.bodyKey)}</p>
          {"noticeKey" in reference && (
            <div role="note" className="border-warning border-l-2 bg-warning/10 px-4 py-3">
              <p className="font-sans font-semibold text-warning text-xs uppercase tracking-wider">
                {t("references.implementationNote")}
              </p>
              <p className="mt-1 font-sans text-sm leading-relaxed">{t(reference.noticeKey)}</p>
            </div>
          )}
          <section aria-labelledby={sourcesHeadingId} className="space-y-3">
            <h3
              id={sourcesHeadingId}
              className="border-border border-t pt-5 font-sans font-semibold text-muted-foreground text-xs uppercase tracking-wider"
            >
              {t("references.sourcesHeading")}
            </h3>
            <ul className="space-y-2">
              {reference.sources.map((source) => (
                <li key={source.href}>
                  <a
                    href={source.href}
                    target="_blank"
                    rel="noreferrer"
                    className="flex items-start gap-3 rounded-md px-2 py-2 text-foreground transition-colors hover:bg-accent hover:text-accent-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                  >
                    <TagBadge className="mt-0.5 shrink-0">
                      {t(`references.kinds.${source.kind}`)}
                    </TagBadge>
                    <span className="min-w-0 flex-1 font-sans text-sm">{t(source.labelKey)}</span>
                    <span className="shrink-0 text-muted-foreground" aria-hidden="true">
                      ↗
                    </span>
                    <span className="sr-only">({t("references.newTab")})</span>
                  </a>
                </li>
              ))}
            </ul>
          </section>
        </div>
      </Dialog>
    </>
  );
}

export { ReferenceDialog };
