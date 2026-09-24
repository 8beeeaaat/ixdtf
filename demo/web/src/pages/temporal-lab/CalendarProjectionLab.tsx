import { useTranslation } from "react-i18next";
import { IxdtfHighlight } from "@/components/IxdtfHighlight";
import { ReferenceDialog } from "@/components/ReferenceDialog";
import { ReproduceSection } from "@/components/ReproduceSection";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card";
import { TagBadge } from "@/components/ui/TagBadge";
import { goCalendarProjectionSample, jsCalendarProjectionSample } from "@/lib/reproduce";
import type { CalendarProjection, GoRoundtripObservation } from "./observations";

interface CalendarProjectionLabProps {
  input: string;
  projections: CalendarProjection[];
  go: GoRoundtripObservation;
}

/**
 * F-6-4: one instant projected into several calendars via the `u-ca` tag —
 * Temporal interprets the annotation, Go ixdtf carries it losslessly, and the
 * instant (unix_nano / epochNanoseconds) never changes.
 */
function CalendarProjectionLab({ input, projections, go }: CalendarProjectionLabProps) {
  const { t } = useTranslation();
  const instant =
    projections.find((projection) => projection.epochNanoseconds !== null)?.epochNanoseconds ??
    null;

  return (
    <section className="space-y-4">
      <Card>
        <CardHeader>
          <div className="flex flex-wrap items-center gap-2">
            <CardTitle>{t("temporalLab.calendarProjection.title")}</CardTitle>
            <ReferenceDialog referenceId="calendarAnnotation" />
          </div>
          <p className="max-w-3xl font-sans text-muted-foreground text-sm">
            {t("temporalLab.calendarProjection.body")}
          </p>
        </CardHeader>
        <CardContent className="space-y-5">
          <IxdtfHighlight value={input} className="text-sm" />

          <div className="space-y-2 border-border border-t pt-4">
            <div className="flex flex-wrap items-center justify-between gap-2">
              <span className="font-sans text-muted-foreground text-xs">
                {t("temporalLab.calendarProjection.instant")}
              </span>
              <TagBadge variant="default">{t("temporalLab.nativeBadge")}</TagBadge>
            </div>
            <p className="break-all font-mono text-muted-foreground text-xs tabular-nums">
              {instant ?? "—"}
            </p>
          </div>

          <div className="overflow-x-auto border-border border-t pt-4">
            <table className="w-full border-collapse text-left">
              <thead>
                <tr className="border-border border-b">
                  <th className="py-2 pr-4 font-medium font-sans text-muted-foreground text-xs">
                    {t("temporalLab.calendarProjection.calendar")}
                  </th>
                  <th className="py-2 pr-4 font-medium font-sans text-muted-foreground text-xs">
                    {t("temporalLab.calendarProjection.year")}
                  </th>
                  <th className="py-2 pr-4 font-medium font-sans text-muted-foreground text-xs">
                    {t("temporalLab.calendarProjection.month")}
                  </th>
                  <th className="py-2 pr-4 font-medium font-sans text-muted-foreground text-xs">
                    {t("temporalLab.calendarProjection.day")}
                  </th>
                  <th className="py-2 font-medium font-sans text-muted-foreground text-xs">
                    {t("temporalLab.calendarProjection.ixdtf")}
                  </th>
                </tr>
              </thead>
              <tbody>
                {projections.map((projection) => (
                  <tr
                    key={projection.calendarId}
                    className="border-border border-b align-top last:border-b-0"
                  >
                    <td className="whitespace-nowrap py-2 pr-4 font-mono text-sm">
                      {projection.calendarId}
                    </td>
                    {projection.error ? (
                      <td colSpan={4} className="py-2">
                        <span className="break-all font-mono text-destructive text-xs">
                          {t("temporalLab.calendarProjection.projectionError")}: {projection.error}
                        </span>
                      </td>
                    ) : (
                      <>
                        <td className="whitespace-nowrap py-2 pr-4 font-mono text-sm tabular-nums">
                          {projection.era !== null && projection.eraYear !== null
                            ? `${projection.era} ${projection.eraYear}`
                            : (projection.year ?? "—")}
                        </td>
                        <td className="whitespace-nowrap py-2 pr-4 font-mono text-sm tabular-nums">
                          {projection.monthCode ?? "—"}
                        </td>
                        <td className="whitespace-nowrap py-2 pr-4 font-mono text-sm tabular-nums">
                          {projection.day ?? "—"}
                        </td>
                        <td className="max-w-96 py-2">
                          <span className="break-all font-mono text-sm">
                            {projection.formatted ?? "—"}
                          </span>
                        </td>
                      </>
                    )}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <div className="space-y-2 border-border border-t pt-4" aria-live="polite">
            <div className="flex flex-wrap items-center justify-between gap-2">
              <span className="font-sans text-muted-foreground text-xs">
                {t("temporalLab.calendarProjection.goRoundtrip")}
              </span>
              <TagBadge variant="default">{t("temporalLab.goBadge")}</TagBadge>
            </div>
            {go.pending ? (
              <p className="font-sans text-muted-foreground text-xs">
                {t("temporalLab.goPending")}
              </p>
            ) : go.formatted ? (
              <>
                <p className="break-all font-mono text-sm">{go.formatted}</p>
                <p className="break-all font-mono text-muted-foreground text-xs tabular-nums">
                  unix_nano: {go.unixNano ?? "—"}
                  {go.lossless === null ? "" : ` · lossless: ${String(go.lossless)}`}
                </p>
              </>
            ) : (
              <p className="font-sans text-muted-foreground text-xs">
                {go.error ?? t("temporalLab.goUnavailable")}
              </p>
            )}
          </div>

          <p className="border-border border-t pt-4 font-sans text-muted-foreground text-xs">
            {t("temporalLab.calendarProjection.conclusion")}
          </p>
        </CardContent>
      </Card>
      <ReproduceSection
        goSample={goCalendarProjectionSample({ input })}
        jsSample={jsCalendarProjectionSample(
          input,
          projections.map((projection) => projection.calendarId),
        )}
      />
    </section>
  );
}

export type { CalendarProjectionLabProps };
export { CalendarProjectionLab };
