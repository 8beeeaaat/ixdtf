import { describe, expect, it } from "vitest";
import { browserParse, type TemporalParser } from "./browserParse";

// Minimal Temporal mock: ZonedDateTime.from only accepts strings whose offset
// matches the annotated zone; Instant.from accepts any offset-bearing string.
const mockTemporal: TemporalParser = {
  ZonedDateTime: {
    from(item) {
      if (item.includes("[Asia/Tokyo]") && item.includes("+09:00")) {
        return {
          toString: () => item,
          epochNanoseconds: 1_800_000_000_000_000_000n,
          timeZoneId: "Asia/Tokyo",
          calendarId: "iso8601",
        };
      }
      throw new RangeError("ZonedDateTime: could not parse");
    },
  },
  Instant: {
    from(item) {
      if (/^\d{4}-\d{2}-\d{2}T[\d:.]+(Z|[+-]\d{2}:\d{2})$/.test(item)) {
        return {
          toString: () => `${item}::instant`,
          epochNanoseconds: 900_000_000_000_000_000n,
        };
      }
      throw new RangeError("Instant: could not parse");
    },
  },
};

describe("browserParse", () => {
  it("returns unavailable when Temporal is null", () => {
    const result = browserParse("2026-07-07T14:30:00Z", null);
    expect(result.ok).toBe(false);
  });

  it("parses a full IXDTF string via ZonedDateTime", () => {
    const result = browserParse("2026-07-07T23:30:00+09:00[Asia/Tokyo]", mockTemporal);
    expect(result).toEqual({
      ok: true,
      via: "zonedDateTime",
      formatted: "2026-07-07T23:30:00+09:00[Asia/Tokyo]",
      epochNanoseconds: "1800000000000000000",
      timeZone: "Asia/Tokyo",
      calendar: "iso8601",
    });
  });

  it("falls back to Instant for a bare RFC 3339 string (no bracket)", () => {
    const result = browserParse("2026-07-07T14:30:00Z", mockTemporal);
    expect(result.ok).toBe(true);
    if (result.ok) {
      expect(result.via).toBe("instant");
      expect(result.timeZone).toBeNull();
    }
  });

  it("does NOT fall back to Instant when a bracket is present (offset mismatch)", () => {
    // Has a [Asia/Tokyo] bracket but +02:00 — genuine error, no silent success.
    const result = browserParse("2026-07-07T23:30:00+02:00[Asia/Tokyo]", mockTemporal);
    expect(result.ok).toBe(false);
    if (!result.ok) {
      expect(result.error).toContain("ZonedDateTime");
    }
  });

  it("reports an error for unparseable garbage", () => {
    const result = browserParse("not-a-timestamp", mockTemporal);
    expect(result.ok).toBe(false);
    if (!result.ok) {
      expect(result.error).toContain("Instant");
    }
  });
});
