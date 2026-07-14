import { useTranslation } from "react-i18next";
import { IxdtfHighlight } from "@/components/IxdtfHighlight";
import { ReferenceDialog } from "@/components/ReferenceDialog";
import { ReproduceSection } from "@/components/ReproduceSection";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card";
import { TagBadge } from "@/components/ui/TagBadge";
import { goZoneArithmeticSample, jsZoneArithmeticSample } from "@/lib/reproduce";
import { cn } from "@/lib/utils";
import { epochDeltaSeconds } from "./DstOverlapLab";
import type {
  ArithmeticObservation,
  GoObservation,
  ZoneArithmeticObservations,
} from "./observations";

interface ZoneArithmeticLabProps {
  input: string;
  /** `null` when this browser has no native Temporal. */
  native: ZoneArithmeticObservations | null;
  polyfill: ZoneArithmeticObservations;
  plusDayGo: GoObservation;
  plusHoursGo: GoObservation;
}

/**
 * F-6-3: `+1 day` (calendar arithmetic) vs `+24 hours` (exact-time arithmetic)
 * across a DST boundary — the time-zone annotation keeps both meaningful after
 * serialization, and Go ixdtf re-parses Temporal's results to the same instants.
 */
function ZoneArithmeticLab({
  input,
  native,
  polyfill,
  plusDayGo,
  plusHoursGo,
}: ZoneArithmeticLabProps) {
  const { t } = useTranslation();

  // One implementation's epoch for a result (`observation === null` = no native Temporal).
  const renderEpochBlock = (label: string, observation: ArithmeticObservation | null) => (
    <div className="space-y-2 border-border border-t pt-4" aria-live="polite">
      <div className="flex flex-wrap items-center justify-between gap-2">
        <span className="font-sans text-muted-foreground text-xs">
          {t("temporalLab.zoneArithmetic.temporalEpoch")}
        </span>
        <TagBadge variant="default">{label}</TagBadge>
      </div>
      {observation === null ? (
        <p className="font-sans text-muted-foreground text-xs">
          {t("temporalLab.nativeUnsupported")}
        </p>
      ) : (
        <p
          className={cn(
            "break-all font-mono text-xs",
            observation.error ? "text-destructive" : "text-muted-foreground tabular-nums",
          )}
        >
          {observation.epochNanoseconds ?? "—"}
        </p>
      )}
    </div>
  );

  const renderResult = (
    id: "plusDay" | "plusHours",
    nativeObservation: ArithmeticObservation | null,
    polyfillObservation: ArithmeticObservation,
    go: GoObservation,
  ) => {
    const titleId = `zone-arithmetic-${id}`;
    // polyfill is always present, so the elapsed delta is measured against it.
    const deltaSeconds = epochDeltaSeconds(
      polyfill.base.epochNanoseconds,
      polyfillObservation.epochNanoseconds,
    );

    return (
      <article aria-labelledby={titleId} className="flex min-w-0 flex-col gap-5 py-6 md:py-2">
        <div className="flex flex-wrap items-end justify-between gap-3">
          <div>
            <h3
              id={titleId}
              className="font-sans font-semibold text-muted-foreground text-xs uppercase tracking-wider"
            >
              {t(`temporalLab.zoneArithmetic.${id}`)}
            </h3>
            <p className="mt-1 font-mono text-3xl text-ixdtf-offset tabular-nums">
              {deltaSeconds === null
                ? "—"
                : t("temporalLab.zoneArithmetic.hours", { count: deltaSeconds / 3600 })}
            </p>
          </div>
          <span className="font-sans text-muted-foreground text-xs">
            {t("temporalLab.zoneArithmetic.realDelta")}
          </span>
        </div>

        {polyfillObservation.formatted ? (
          <IxdtfHighlight value={polyfillObservation.formatted} className="text-sm" />
        ) : (
          <p className="break-all font-mono text-destructive text-xs">
            {polyfillObservation.error ?? "—"}
          </p>
        )}

        {renderEpochBlock(t("temporal.implNative"), nativeObservation)}
        {renderEpochBlock(t("temporal.implPolyfill"), polyfillObservation)}

        <div className="space-y-2 border-border border-t pt-4" aria-live="polite">
          <div className="flex flex-wrap items-center justify-between gap-2">
            <span className="font-sans text-muted-foreground text-xs">
              {t("temporalLab.zoneArithmetic.goParsed")}
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
            <CardTitle>{t("temporalLab.zoneArithmetic.title")}</CardTitle>
            <ReferenceDialog referenceId="zoneArithmetic" />
          </div>
          <p className="max-w-3xl font-sans text-muted-foreground text-sm">
            {t("temporalLab.zoneArithmetic.body")}
          </p>
        </CardHeader>
        <CardContent>
          <div className="space-y-2 border-border border-b pb-5">
            <span className="font-sans font-semibold text-muted-foreground text-xs uppercase tracking-wider">
              {t("temporalLab.zoneArithmetic.base")}
            </span>
            <IxdtfHighlight value={input} className="text-sm" />
            <p className="break-all font-mono text-muted-foreground text-xs tabular-nums">
              {polyfill.base.epochNanoseconds ?? polyfill.base.error ?? "—"}
            </p>
          </div>
          <div className="grid items-stretch gap-0 md:grid-cols-2 md:gap-6">
            {renderResult("plusDay", native?.plusDay ?? null, polyfill.plusDay, plusDayGo)}
            <div className="md:border-border md:border-l md:pl-6">
              {renderResult(
                "plusHours",
                native?.plusHours ?? null,
                polyfill.plusHours,
                plusHoursGo,
              )}
            </div>
          </div>
          <p className="mt-5 border-border border-t pt-4 font-sans text-muted-foreground text-xs">
            {t("temporalLab.zoneArithmetic.conclusion")}
          </p>
        </CardContent>
      </Card>
      <ReproduceSection
        goSample={goZoneArithmeticSample({ input })}
        jsSample={jsZoneArithmeticSample(input)}
      />
    </section>
  );
}

export type { ZoneArithmeticLabProps };
export { ZoneArithmeticLab };
