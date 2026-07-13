import { getTemporal } from "@/lib/temporal/detect";

/**
 * Minimal structural view of the Temporal API surface `browserParse` needs.
 * `TemporalApi` (from detect.ts) is assignable to this, and a tiny mock in the
 * unit test is too — so we avoid depending on the full native namespace.
 */
export interface TemporalParser {
  ZonedDateTime: {
    from(item: string): {
      toString(): string;
      epochNanoseconds: bigint;
      timeZoneId: string;
      calendarId: string;
    };
  };
  Instant: {
    from(item: string): {
      toString(): string;
      epochNanoseconds: bigint;
    };
  };
}

export interface BrowserParseSuccess {
  ok: true;
  /** Which native API accepted the string. */
  via: "zonedDateTime" | "instant";
  /** `.toString()` of the parsed value (round-tripped IXDTF / RFC 3339). */
  formatted: string;
  /** Epoch nanoseconds as a decimal string (D-3: int64 exceeds JS safe range). */
  epochNanoseconds: string;
  timeZone: string | null;
  calendar: string | null;
}

export interface BrowserParseFailure {
  ok: false;
  error: string;
}

export type BrowserParseResult = BrowserParseSuccess | BrowserParseFailure;

function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : String(error);
}

/**
 * Parse an IXDTF / RFC 3339 string with the browser's native Temporal.
 *
 * Strategy (architecture.md「ブラウザ側 Temporal の利用ポイント」):
 * - Try `ZonedDateTime.from` first.
 * - Only when the input has NO `[...]` time-zone annotation (a bare RFC 3339
 *   string, which `ZonedDateTime.from` cannot accept) fall back to
 *   `Instant.from`. Strings that DO carry a bracket but still fail (offset
 *   mismatch, unknown critical tag) are genuine errors and must not silently
 *   succeed via Instant — this is exactly the browser/server behaviour the
 *   fixtures (`expect_browser`) contrast on the Interop screen.
 */
export function browserParse(
  input: string,
  temporal: TemporalParser | null = getTemporal(),
): BrowserParseResult {
  if (!temporal) {
    return { ok: false, error: "Temporal API is not available in this browser." };
  }
  const trimmed = input.trim();
  if (trimmed.length === 0) {
    return { ok: false, error: "Empty input." };
  }

  try {
    const zdt = temporal.ZonedDateTime.from(trimmed);
    return {
      ok: true,
      via: "zonedDateTime",
      formatted: zdt.toString(),
      epochNanoseconds: zdt.epochNanoseconds.toString(),
      timeZone: zdt.timeZoneId,
      calendar: zdt.calendarId,
    };
  } catch (zonedError) {
    if (!trimmed.includes("[")) {
      try {
        const instant = temporal.Instant.from(trimmed);
        return {
          ok: true,
          via: "instant",
          formatted: instant.toString(),
          epochNanoseconds: instant.epochNanoseconds.toString(),
          timeZone: null,
          calendar: null,
        };
      } catch (instantError) {
        return { ok: false, error: errorMessage(instantError) };
      }
    }
    return { ok: false, error: errorMessage(zonedError) };
  }
}
