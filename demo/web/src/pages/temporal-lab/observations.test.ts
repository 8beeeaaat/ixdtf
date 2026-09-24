import { describe, expect, it } from "vitest";
import {
  type CalendarProjectionTemporal,
  observeCalendarProjections,
  observeZoneArithmetic,
  type ZoneArithmeticTemporal,
} from "./observations";

// 観測ヘルパーの配線とエラー分離のみを検証する。実値 (DST 境界の演算結果や
// 暦投影) の検証は Storybook play テストがネイティブ Temporal で行う。

function arithmeticMock(): ZoneArithmeticTemporal {
  const make = (formatted: string, epochNanoseconds: bigint) => ({
    epochNanoseconds,
    toString: () => formatted,
    add(duration: { days?: number; hours?: number }) {
      if (duration.days === 1) {
        return make("plus-day", 90_000_000_000_000n);
      }
      if (duration.hours === 24) {
        return make("plus-hours", 86_400_000_000_000n);
      }
      throw new Error(`unexpected duration: ${JSON.stringify(duration)}`);
    },
  });
  return {
    ZonedDateTime: {
      from: (item: string) => {
        if (item !== "valid") {
          throw new RangeError("invalid IXDTF");
        }
        return make("base", 0n);
      },
    },
  };
}

describe("observeZoneArithmetic", () => {
  it("observes the base and both additions as decimal strings", () => {
    const result = observeZoneArithmetic("valid", arithmeticMock());
    expect(result.base).toEqual({ formatted: "base", epochNanoseconds: "0", error: null });
    expect(result.plusDay.formatted).toBe("plus-day");
    expect(result.plusDay.epochNanoseconds).toBe("90000000000000");
    expect(result.plusHours.formatted).toBe("plus-hours");
    expect(result.plusHours.epochNanoseconds).toBe("86400000000000");
  });

  it("propagates a parse failure to all three observations", () => {
    const result = observeZoneArithmetic("broken", arithmeticMock());
    expect(result.base.error).toBe("invalid IXDTF");
    expect(result.plusDay).toEqual(result.base);
    expect(result.plusHours.formatted).toBeNull();
  });

  it("reports the missing Temporal API", () => {
    const result = observeZoneArithmetic("valid", null);
    expect(result.base.error).toContain("Temporal API");
  });
});

function calendarMock(): CalendarProjectionTemporal {
  const projected = (calendar: string) => ({
    withCalendar: (next: string) => {
      if (next === "broken") {
        throw new RangeError("unknown calendar");
      }
      return projected(next);
    },
    toString: () => `formatted-${calendar}`,
    epochNanoseconds: 42n,
    era: calendar === "japanese" ? "reiwa" : undefined,
    eraYear: calendar === "japanese" ? 8 : undefined,
    year: 2026,
    monthCode: "M07",
    day: 7,
  });
  return {
    ZonedDateTime: {
      from: (item: string) => {
        if (item !== "valid") {
          throw new RangeError("invalid IXDTF");
        }
        return projected("iso8601");
      },
    },
  };
}

describe("observeCalendarProjections", () => {
  it("keeps the instant identical across projections and surfaces era fields", () => {
    const [iso, japanese] = observeCalendarProjections(
      "valid",
      ["iso8601", "japanese"],
      calendarMock(),
    );
    expect(iso.epochNanoseconds).toBe("42");
    expect(japanese.epochNanoseconds).toBe("42");
    expect(iso.era).toBeNull();
    expect(japanese.era).toBe("reiwa");
    expect(japanese.eraYear).toBe(8);
    expect(japanese.formatted).toBe("formatted-japanese");
  });

  it("isolates a failing calendar without invalidating the rest", () => {
    const [ok, broken] = observeCalendarProjections("valid", ["iso8601", "broken"], calendarMock());
    expect(ok.error).toBeNull();
    expect(broken.error).toBe("unknown calendar");
    expect(broken.calendarId).toBe("broken");
    expect(broken.formatted).toBeNull();
  });

  it("propagates a parse failure to every projection", () => {
    const projections = observeCalendarProjections("nope", ["iso8601", "hebrew"], calendarMock());
    expect(projections).toHaveLength(2);
    for (const projection of projections) {
      expect(projection.error).toBe("invalid IXDTF");
    }
  });
});
