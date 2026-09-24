import { useTranslation } from "react-i18next";
import { IxdtfHighlight, IxdtfLegend } from "@/components/IxdtfHighlight";
import { ReferenceDialog } from "@/components/ReferenceDialog";
import { ReproduceSection } from "@/components/ReproduceSection";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card";
import { TagBadge } from "@/components/ui/TagBadge";
import { goDstOverlapSample, jsDstOverlapSample } from "@/lib/reproduce";
import { cn } from "@/lib/utils";
import type { GoObservation } from "./observations";

interface DstOverlapObservation {
  id: "earlier" | "later";
  input: string;
  epochNanoseconds: string | null;
  error: string | null;
}

interface DstOverlapLabProps {
  earlier: DstOverlapObservation;
  later: DstOverlapObservation;
  earlierGo: GoObservation;
  laterGo: GoObservation;
}

function epochDeltaSeconds(earlier: string | null, later: string | null): number | null {
  if (earlier === null || later === null) {
    return null;
  }
  try {
    return Number((BigInt(later) - BigInt(earlier)) / 1_000_000_000n);
  } catch {
    return null;
  }
}

/** F-6-1: one repeated wall-clock time, made exact by two different offsets. */
function DstOverlapLab({ earlier, later, earlierGo, laterGo }: DstOverlapLabProps) {
  const { t } = useTranslation();
  const deltaSeconds = epochDeltaSeconds(earlier.epochNanoseconds, later.epochNanoseconds);

  const renderObservation = (observation: DstOverlapObservation, go: GoObservation) => {
    const titleId = `dst-overlap-${observation.id}`;

    return (
      <article aria-labelledby={titleId} className="flex min-w-0 flex-col gap-5 py-6 md:py-2">
        <div className="flex flex-wrap items-center justify-between gap-3">
          <div>
            <h3
              id={titleId}
              className="font-sans font-semibold text-muted-foreground text-xs uppercase tracking-wider"
            >
              {t(`temporalLab.dstOverlap.${observation.id}`)}
            </h3>
            <p className="mt-1 font-mono text-3xl tabular-nums">
              {t("temporalLab.dstOverlap.clockTime")}
            </p>
          </div>
          <TagBadge variant="default">
            {t(`temporalLab.dstOverlap.${observation.id}Offset`)}
          </TagBadge>
        </div>

        <IxdtfHighlight value={observation.input} className="text-sm" />

        <div className="space-y-2 border-border border-t pt-4" aria-live="polite">
          <div className="flex flex-wrap items-center justify-between gap-2">
            <span className="font-sans text-muted-foreground text-xs">
              {t("temporalLab.dstOverlap.epoch")}
            </span>
            <TagBadge variant={observation.error ? "mismatch" : "default"}>
              {observation.error ? t("common.error") : t("temporalLab.dstOverlap.nativeTemporal")}
            </TagBadge>
          </div>
          <p
            className={cn(
              "break-all font-mono text-xs",
              observation.error ? "text-destructive" : "text-muted-foreground tabular-nums",
            )}
          >
            {observation.error
              ? t("temporalLab.dstOverlap.parseError")
              : (observation.epochNanoseconds ?? "—")}
          </p>
          {observation.error && (
            <details className="font-sans text-muted-foreground text-xs">
              <summary className="cursor-pointer">
                {t("temporalLab.dstOverlap.errorDetails")}
              </summary>
              <code className="mt-2 block break-all font-mono text-destructive">
                {observation.error}
              </code>
            </details>
          )}
        </div>

        <div className="space-y-2 border-border border-t pt-4" aria-live="polite">
          <div className="flex flex-wrap items-center justify-between gap-2">
            <span className="font-sans text-muted-foreground text-xs">
              {t("temporalLab.goEpoch")}
            </span>
            <TagBadge variant="default">{t("temporalLab.goBadge")}</TagBadge>
          </div>
          <p className="break-all font-mono text-muted-foreground text-xs tabular-nums">
            {go.pending
              ? t("temporalLab.goPending")
              : (go.unixNano ?? t("temporalLab.goUnavailable"))}
          </p>
        </div>
      </article>
    );
  };

  return (
    <section className="space-y-4">
      <Card>
        <CardHeader>
          <div className="flex flex-wrap items-center gap-2">
            <CardTitle>{t("temporalLab.dstOverlap.title")}</CardTitle>
            <ReferenceDialog referenceId="dstAmbiguity" />
          </div>
          <p className="max-w-3xl font-sans text-muted-foreground text-sm">
            {t("temporalLab.dstOverlap.body")}
          </p>
        </CardHeader>
        <CardContent>
          <div className="grid items-stretch gap-0 md:grid-cols-[minmax(0,1fr)_9rem_minmax(0,1fr)] md:gap-6">
            {renderObservation(earlier, earlierGo)}

            <div
              className="flex flex-col items-center justify-center border-border border-y py-5 md:border-x md:border-y-0 md:py-0"
              aria-live="polite"
            >
              <span className="font-sans text-muted-foreground text-xs uppercase tracking-wider">
                {t("temporalLab.dstOverlap.elapsed")}
              </span>
              <span className="mt-1 font-mono text-muted-foreground md:hidden" aria-hidden="true">
                ↓
              </span>
              <span
                className="mt-1 hidden font-mono text-muted-foreground md:inline"
                aria-hidden="true"
              >
                →
              </span>
              <strong className="mt-1 font-mono font-normal text-4xl text-ixdtf-offset tabular-nums">
                {deltaSeconds === null
                  ? "—"
                  : t("temporalLab.dstOverlap.hours", { count: deltaSeconds / 3600 })}
              </strong>
              <span className="mt-1 font-sans text-muted-foreground text-xs">
                {deltaSeconds === null
                  ? t("temporalLab.dstOverlap.unavailable")
                  : t("temporalLab.dstOverlap.seconds", { count: deltaSeconds })}
              </span>
            </div>

            {renderObservation(later, laterGo)}
          </div>
          <IxdtfLegend className="mt-5 border-border border-t pt-4" />
          <p className="mt-5 border-border border-t pt-4 font-sans text-muted-foreground text-xs">
            {t("temporalLab.dstOverlap.conclusion")}
          </p>
        </CardContent>
      </Card>
      <ReproduceSection
        goSample={goDstOverlapSample({ earlier: earlier.input, later: later.input })}
        jsSample={jsDstOverlapSample({ earlier: earlier.input, later: later.input })}
      />
    </section>
  );
}

export type { DstOverlapLabProps, DstOverlapObservation };
export { DstOverlapLab, epochDeltaSeconds };
