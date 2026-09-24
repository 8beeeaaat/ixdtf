import { describe, expect, it } from "vitest";
import {
  goCalendarProjectionSample,
  goDstOverlapSample,
  goNowSample,
  goParseSample,
  goRoundtripSample,
  goZoneArithmeticSample,
  jsCalendarProjectionSample,
  jsDstOverlapSample,
  jsNowSample,
  jsParseSample,
  jsRoundtripSample,
  jsZoneArithmeticSample,
} from "@/lib/reproduce";

// ロジック固有の境界値 (文字列補間・エスケープ・ブラケット検出) のため直書きする
// — ixdtf-fixtures スキルの例外条項。パース挙動の検証はここでは行わない。
const INPUT = "2026-07-11T09:00:00+09:00[Asia/Tokyo][u-ca=japanese]";

describe("goParseSample", () => {
  it("embeds the input and the strict flag into ixdtf.Parse", () => {
    const sample = goParseSample({ input: INPUT, strict: true, validateOnly: false });
    expect(sample).toContain(`input := "${INPUT}"`);
    expect(sample).toContain("parsedTime, ext, err := ixdtf.Parse(input, true)");
    expect(sample).toContain('"github.com/8beeeaaat/ixdtf"');
  });

  it("uses ixdtf.Validate in validate-only mode", () => {
    const sample = goParseSample({ input: INPUT, strict: false, validateOnly: true });
    expect(sample).toContain("ixdtf.Validate(input, false)");
    expect(sample).not.toContain("ixdtf.Parse(");
  });

  it("escapes quotes and backslashes for the Go string literal", () => {
    const sample = goParseSample({ input: 'a"b\\c', strict: false, validateOnly: false });
    expect(sample).toContain('input := "a\\"b\\\\c"');
  });
});

describe("goRoundtripSample", () => {
  it("parses then formats and prints the lossless comparison", () => {
    const sample = goRoundtripSample({ input: INPUT, strict: false });
    expect(sample).toContain("ixdtf.Parse(input, false)");
    expect(sample).toContain("ixdtf.FormatNano(parsedTime, ext)");
    expect(sample).toContain('fmt.Println("lossless:", formatted == input)');
  });
});

describe("goNowSample", () => {
  it("resolves the location and adds the u-ca tag when a calendar is set", () => {
    const sample = goNowSample({ timeZone: "Asia/Tokyo", calendar: "japanese" });
    expect(sample).toContain('time.LoadLocation("Asia/Tokyo")');
    expect(sample).toContain("Location: loc,");
    expect(sample).toContain('map[string]string{ixdtf.ExtensionUnicodeCalendar: "japanese"}');
    expect(sample).toContain("time.Now().In(loc)");
  });

  it("skips LoadLocation and annotations for UTC without calendar", () => {
    const sample = goNowSample({ timeZone: "UTC" });
    expect(sample).toContain("ixdtf.NewIXDTFExtensions(nil)");
    expect(sample).toContain("time.Now().UTC()");
    expect(sample).not.toContain("LoadLocation");
  });
});

describe("jsParseSample", () => {
  it("uses ZonedDateTime.from when the input has a suffix annotation", () => {
    expect(jsParseSample(INPUT)).toContain("Temporal.ZonedDateTime.from(input)");
  });

  it("falls back to Instant.from for bare RFC 3339 input (browserParse rule)", () => {
    const sample = jsParseSample("2026-07-11T00:00:00Z");
    expect(sample).toContain("Temporal.Instant.from(input)");
    expect(sample).not.toContain("ZonedDateTime");
  });
});

describe("jsRoundtripSample", () => {
  it("round-trips through toString and compares with the input", () => {
    const sample = jsRoundtripSample(INPUT);
    expect(sample).toContain("Temporal.ZonedDateTime.from(input).toString()");
    expect(sample).toContain('console.log("lossless:", formatted === input);');
  });
});

describe("goDstOverlapSample", () => {
  it("parses both overlap strings strictly and prints the delta", () => {
    const sample = goDstOverlapSample({ earlier: "e-input", later: "l-input" });
    expect(sample).toContain('earlier, _, err := ixdtf.Parse("e-input", true)');
    expect(sample).toContain('later, _, err := ixdtf.Parse("l-input", true)');
    expect(sample).toContain('fmt.Println("delta:", later.Sub(earlier))');
  });
});

describe("goZoneArithmeticSample", () => {
  it("contrasts AddDate (wall clock) with Add(24h) (exact time)", () => {
    const sample = goZoneArithmeticSample({ input: INPUT });
    expect(sample).toContain(`ixdtf.Parse("${INPUT}", true)`);
    expect(sample).toContain("ixdtf.FormatNano(parsedTime.AddDate(0, 0, 1), ext)");
    expect(sample).toContain("ixdtf.FormatNano(parsedTime.Add(24*time.Hour), ext)");
    expect(sample).toContain('"time"');
  });
});

describe("goCalendarProjectionSample", () => {
  it("prints the instant, the carried u-ca tag, and the lossless roundtrip", () => {
    const sample = goCalendarProjectionSample({ input: INPUT });
    expect(sample).toContain("fmt.Println(parsedTime.UnixNano())");
    expect(sample).toContain('fmt.Println("u-ca:", ext.Tags[ixdtf.ExtensionUnicodeCalendar])');
    expect(sample).toContain('fmt.Println("lossless:", formatted == input)');
  });
});

describe("jsDstOverlapSample", () => {
  it("resolves both strings and prints the Duration between them", () => {
    const sample = jsDstOverlapSample({ earlier: "e-input", later: "l-input" });
    expect(sample).toContain('Temporal.ZonedDateTime.from("e-input")');
    expect(sample).toContain('Temporal.ZonedDateTime.from("l-input")');
    expect(sample).toContain('console.log("delta:", later.since(earlier).toString());');
  });
});

describe("jsZoneArithmeticSample", () => {
  it("adds one calendar day and 24 exact hours to the same ZonedDateTime", () => {
    const sample = jsZoneArithmeticSample(INPUT);
    expect(sample).toContain("zdt.add({ days: 1 }).toString()");
    expect(sample).toContain("zdt.add({ hours: 24 }).toString()");
  });
});

describe("jsCalendarProjectionSample", () => {
  it("iterates the given calendars with withCalendar", () => {
    const sample = jsCalendarProjectionSample(INPUT, ["iso8601", "japanese"]);
    expect(sample).toContain('for (const calendar of ["iso8601", "japanese"]) {');
    expect(sample).toContain("zdt.withCalendar(calendar)");
    expect(sample).toContain('console.log("instant:", zdt.epochNanoseconds);');
  });
});

describe("jsNowSample", () => {
  it("chains withCalendar only when a calendar is set", () => {
    expect(jsNowSample({ timeZone: "Asia/Tokyo", calendar: "japanese" })).toContain(
      'Temporal.Now.zonedDateTimeISO("Asia/Tokyo").withCalendar("japanese");',
    );
    expect(jsNowSample({ timeZone: "Asia/Tokyo" })).toContain(
      'Temporal.Now.zonedDateTimeISO("Asia/Tokyo");',
    );
  });
});
