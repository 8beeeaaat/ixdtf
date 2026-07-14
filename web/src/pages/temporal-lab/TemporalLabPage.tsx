import { useMemo } from "react";
import { useTranslation } from "react-i18next";
import { TagBadge } from "@/components/ui/TagBadge";
import { useParseIxdtf, useRoundtripIxdtf } from "@/generated/api/endpoints";
import { requiredSample } from "@/lib/fixtures";
import { browserParse } from "@/lib/temporal/browserParse";
import { getNativeTemporal, getPolyfillTemporal, type TemporalApi } from "@/lib/temporal/detect";
import { CalendarProjectionLab } from "./CalendarProjectionLab";
import { DstOverlapLab, type DstOverlapObservation, type EpochReading } from "./DstOverlapLab";
import {
  type GoObservation,
  type GoRoundtripObservation,
  observeCalendarProjections,
  observeZoneArithmetic,
  PROJECTION_CALENDARS,
} from "./observations";
import { ZoneArithmeticLab } from "./ZoneArithmeticLab";

const DST_OVERLAP_EARLIER = requiredSample("dst-overlap-earlier");
const DST_OVERLAP_LATER = requiredSample("dst-overlap-later");
const CALENDAR_BASE = requiredSample("calendar-japanese");

function readEpoch(input: string, temporal: TemporalApi): EpochReading {
  const result = browserParse(input, temporal);
  return result.ok
    ? { epochNanoseconds: result.epochNanoseconds, error: null }
    : { epochNanoseconds: null, error: result.error };
}

/** One DST-overlap instant parsed by native (or `null`) and polyfill Temporal. */
function observe(
  id: DstOverlapObservation["id"],
  input: string,
  native: TemporalApi | null,
  polyfill: TemporalApi,
): DstOverlapObservation {
  return {
    id,
    input,
    native: native ? readEpoch(input, native) : null,
    polyfill: readEpoch(input, polyfill),
  };
}

function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : String(error);
}

/**
 * F-6-5: live Go ixdtf measurement of one IXDTF string via `POST /api/ixdtf/parse`.
 * `input === null` means the Temporal side produced nothing to hand over.
 */
function useGoParse(input: string | null, strict: boolean): GoObservation {
  const query = useParseIxdtf(
    { input: input ?? "", strict },
    { query: { enabled: input !== null } },
  );
  if (input === null) {
    return { unixNano: null, error: null, pending: false };
  }
  const response = query.data?.status === 200 ? query.data.data : null;
  return {
    unixNano: response?.result?.unix_nano ?? null,
    error: response?.error?.message ?? (query.isError ? errorMessage(query.error) : null),
    pending: query.isPending,
  };
}

/** F-6-4: Go carries the u-ca tag losslessly through `POST /api/ixdtf/roundtrip`. */
function useGoRoundtrip(input: string, strict: boolean): GoRoundtripObservation {
  const query = useRoundtripIxdtf({ input, strict });
  const response = query.data?.status === 200 ? query.data.data : null;
  return {
    unixNano: response?.parse.result?.unix_nano ?? null,
    formatted: response?.formatted ?? null,
    lossless: response?.lossless ?? null,
    error: response?.parse.error?.message ?? (query.isError ? errorMessage(query.error) : null),
    pending: query.isPending,
  };
}

/** F-6: focused demonstrations of IXDTF semantics through native Temporal. */
function TemporalLabPage() {
  const { t } = useTranslation();
  // F-3 と同じ思想で 3 実装を明示指定して併記する。native はブラウザ非対応なら
  // null → 各実験のネイティブブロックに案内メッセージを出す。
  const nativeTemporal = getNativeTemporal();
  const polyfillTemporal = getPolyfillTemporal();

  const earlier = observe("earlier", DST_OVERLAP_EARLIER.input, nativeTemporal, polyfillTemporal);
  const later = observe("later", DST_OVERLAP_LATER.input, nativeTemporal, polyfillTemporal);

  const nativeArithmetic = useMemo(
    () =>
      nativeTemporal ? observeZoneArithmetic(DST_OVERLAP_EARLIER.input, nativeTemporal) : null,
    [nativeTemporal],
  );
  const polyfillArithmetic = useMemo(
    () => observeZoneArithmetic(DST_OVERLAP_EARLIER.input, polyfillTemporal),
    [polyfillTemporal],
  );
  const nativeProjections = useMemo(
    () =>
      nativeTemporal
        ? observeCalendarProjections(CALENDAR_BASE.input, PROJECTION_CALENDARS, nativeTemporal)
        : null,
    [nativeTemporal],
  );
  const polyfillProjections = useMemo(
    () => observeCalendarProjections(CALENDAR_BASE.input, PROJECTION_CALENDARS, polyfillTemporal),
    [polyfillTemporal],
  );

  const earlierGo = useGoParse(DST_OVERLAP_EARLIER.input, DST_OVERLAP_EARLIER.strict);
  const laterGo = useGoParse(DST_OVERLAP_LATER.input, DST_OVERLAP_LATER.strict);
  const plusDayGo = useGoParse(polyfillArithmetic.plusDay.formatted, true);
  const plusHoursGo = useGoParse(polyfillArithmetic.plusHours.formatted, true);
  const calendarGo = useGoRoundtrip(CALENDAR_BASE.input, CALENDAR_BASE.strict);

  return (
    <div className="mx-auto max-w-6xl space-y-12 px-6">
      <header>
        <div className="flex flex-wrap items-center gap-3">
          <h1 className="font-sans font-semibold text-2xl tracking-tight">
            {t("temporalLab.title")}
          </h1>
          <TagBadge variant="default">{t("temporal.implNative")}</TagBadge>
          <TagBadge variant="default">{t("temporal.implPolyfill")}</TagBadge>
          <TagBadge variant="default">{t("temporalLab.goBadge")}</TagBadge>
        </div>
        <p className="mt-1 max-w-3xl font-sans text-muted-foreground text-sm">
          {t("temporalLab.tagline")}
        </p>
      </header>

      <DstOverlapLab earlier={earlier} later={later} earlierGo={earlierGo} laterGo={laterGo} />
      <ZoneArithmeticLab
        input={DST_OVERLAP_EARLIER.input}
        native={nativeArithmetic}
        polyfill={polyfillArithmetic}
        plusDayGo={plusDayGo}
        plusHoursGo={plusHoursGo}
      />
      <CalendarProjectionLab
        input={CALENDAR_BASE.input}
        native={nativeProjections}
        polyfill={polyfillProjections}
        go={calendarGo}
      />
    </div>
  );
}

export { TemporalLabPage };
