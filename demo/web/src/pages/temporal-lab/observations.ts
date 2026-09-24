import { getTemporal } from "@/lib/temporal/detect";

/**
 * F-6-3 / F-6-4: pure observation helpers for the Temporal Lab experiments.
 *
 * Like `lib/temporal/browserParse`, each helper depends only on a minimal
 * structural slice of the Temporal API so unit tests can inject a tiny mock —
 * real-value verification stays in the Storybook play tests (native Temporal).
 */

interface ArithmeticZonedDateTime {
  add(duration: { days?: number; hours?: number }): ArithmeticZonedDateTime;
  toString(): string;
  epochNanoseconds: bigint;
}

export interface ZoneArithmeticTemporal {
  ZonedDateTime: {
    from(item: string): ArithmeticZonedDateTime;
  };
}

export interface ArithmeticObservation {
  /** IXDTF string produced by Temporal's own serialization. */
  formatted: string | null;
  epochNanoseconds: string | null;
  error: string | null;
}

export interface ZoneArithmeticObservations {
  base: ArithmeticObservation;
  plusDay: ArithmeticObservation;
  plusHours: ArithmeticObservation;
}

/**
 * F-6-5: a live Go ixdtf measurement fetched through the server API and shown
 * beside the Temporal observation — no match/mismatch scoring (D-12).
 */
export interface GoObservation {
  unixNano: string | null;
  error: string | null;
  pending: boolean;
}

/** F-6-4: the Go side of the calendar demo — a roundtrip carrying the u-ca tag. */
export interface GoRoundtripObservation {
  unixNano: string | null;
  formatted: string | null;
  lossless: boolean | null;
  error: string | null;
  pending: boolean;
}

function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : String(error);
}

function arithmeticError(error: string): ArithmeticObservation {
  return { formatted: null, epochNanoseconds: null, error };
}

function observeZonedDateTime(zdt: ArithmeticZonedDateTime): ArithmeticObservation {
  return {
    formatted: zdt.toString(),
    epochNanoseconds: zdt.epochNanoseconds.toString(),
    error: null,
  };
}

/**
 * F-6-3: `+1 day` (calendar arithmetic, wall-clock preserving) vs `+24 hours`
 * (exact-time arithmetic) from the same IXDTF string across a DST boundary.
 */
export function observeZoneArithmetic(
  input: string,
  temporal: ZoneArithmeticTemporal | null = getTemporal(),
): ZoneArithmeticObservations {
  if (!temporal) {
    const unavailable = arithmeticError("Temporal API is not available in this browser.");
    return { base: unavailable, plusDay: unavailable, plusHours: unavailable };
  }
  try {
    const base = temporal.ZonedDateTime.from(input);
    return {
      base: observeZonedDateTime(base),
      plusDay: observeZonedDateTime(base.add({ days: 1 })),
      plusHours: observeZonedDateTime(base.add({ hours: 24 })),
    };
  } catch (error) {
    const failed = arithmeticError(errorMessage(error));
    return { base: failed, plusDay: failed, plusHours: failed };
  }
}

interface CalendarZonedDateTime {
  withCalendar(calendar: string): CalendarZonedDateTime;
  toString(): string;
  epochNanoseconds: bigint;
  era?: string;
  eraYear?: number;
  year: number;
  monthCode: string;
  day: number;
}

export interface CalendarProjectionTemporal {
  ZonedDateTime: {
    from(item: string): CalendarZonedDateTime;
  };
}

export interface CalendarProjection {
  calendarId: string;
  /** IXDTF string with the `[u-ca=...]` annotation Temporal serializes. */
  formatted: string | null;
  epochNanoseconds: string | null;
  era: string | null;
  eraYear: number | null;
  year: number | null;
  monthCode: string | null;
  day: number | null;
  error: string | null;
}

/** F-6-4: the calendars one instant is projected into. */
export const PROJECTION_CALENDARS = ["iso8601", "japanese", "hebrew", "islamic-umalqura"] as const;

function projectionError(calendarId: string, error: string): CalendarProjection {
  return {
    calendarId,
    formatted: null,
    epochNanoseconds: null,
    era: null,
    eraYear: null,
    year: null,
    monthCode: null,
    day: null,
    error,
  };
}

/**
 * F-6-4: project one instant into several calendar systems. Failures are kept
 * per calendar (a browser may lack an ICU calendar without invalidating the rest).
 */
export function observeCalendarProjections(
  input: string,
  calendars: readonly string[],
  temporal: CalendarProjectionTemporal | null = getTemporal(),
): CalendarProjection[] {
  if (!temporal) {
    const error = "Temporal API is not available in this browser.";
    return calendars.map((calendarId) => projectionError(calendarId, error));
  }
  let base: CalendarZonedDateTime;
  try {
    base = temporal.ZonedDateTime.from(input);
  } catch (error) {
    const message = errorMessage(error);
    return calendars.map((calendarId) => projectionError(calendarId, message));
  }
  return calendars.map((calendarId) => {
    try {
      const projected = base.withCalendar(calendarId);
      return {
        calendarId,
        formatted: projected.toString(),
        epochNanoseconds: projected.epochNanoseconds.toString(),
        era: projected.era ?? null,
        eraYear: projected.eraYear ?? null,
        year: projected.year,
        monthCode: projected.monthCode,
        day: projected.day,
        error: null,
      };
    } catch (error) {
      return projectionError(calendarId, errorMessage(error));
    }
  });
}
