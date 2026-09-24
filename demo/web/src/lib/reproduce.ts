/**
 * F-0-7: reproduction snippets for the collapsible "Reproduce" section.
 *
 * Pure string builders that embed the current inputs into Go
 * (github.com/8beeeaaat/ixdtf) and JavaScript (TC39 Temporal) code samples,
 * mirroring how server/usecase/interactor and lib/temporal/browserParse
 * actually call the two implementations. No validity judgement happens here.
 */

/**
 * Escape a value for a Go / JavaScript double-quoted string literal.
 * JSON escaping is a valid subset of both languages' literal syntax.
 */
const quote = (value: string): string => JSON.stringify(value);

export interface ParseSampleArgs {
  input: string;
  strict: boolean;
  validateOnly: boolean;
}

export interface RoundtripSampleArgs {
  input: string;
  strict: boolean;
}

export interface NowSampleArgs {
  timeZone: string;
  /** Omitted when the UI is on the default iso8601 calendar (no u-ca tag). */
  calendar?: string;
}

function goProgram(imports: readonly string[], body: readonly string[]): string {
  return [
    "package main",
    "",
    "import (",
    ...imports.map((line) => (line === "" ? "" : `\t${line}`)),
    ")",
    "",
    "func main() {",
    ...body.map((line) => (line === "" ? "" : `\t${line}`)),
    "}",
    "",
  ].join("\n");
}

/** Go sample for `POST /api/ixdtf/parse` — ixdtf.Parse / ixdtf.Validate. */
export function goParseSample({ input, strict, validateOnly }: ParseSampleArgs): string {
  const body = validateOnly
    ? [
        `if err := ixdtf.Validate(input, ${strict}); err != nil {`,
        `\tfmt.Println("validation error:", err)`,
        "\treturn",
        "}",
        `fmt.Println("valid")`,
      ]
    : [
        `parsedTime, ext, err := ixdtf.Parse(input, ${strict})`,
        "if err != nil {",
        `\tfmt.Println("parse error:", err)`,
        "\treturn",
        "}",
        "fmt.Println(parsedTime.UnixNano())",
        `fmt.Printf("%+v\\n", ext)`,
      ];
  return goProgram(
    ['"fmt"', "", '"github.com/8beeeaaat/ixdtf"'],
    [`input := ${quote(input)}`, "", ...body],
  );
}

/** Go sample for `POST /api/ixdtf/roundtrip` — Parse then FormatNano. */
export function goRoundtripSample({ input, strict }: RoundtripSampleArgs): string {
  return goProgram(
    ['"fmt"', "", '"github.com/8beeeaaat/ixdtf"'],
    [
      `input := ${quote(input)}`,
      "",
      `parsedTime, ext, err := ixdtf.Parse(input, ${strict})`,
      "if err != nil {",
      `\tfmt.Println("parse error:", err)`,
      "\treturn",
      "}",
      "formatted, err := ixdtf.FormatNano(parsedTime, ext)",
      "if err != nil {",
      `\tfmt.Println("format error:", err)`,
      "\treturn",
      "}",
      "fmt.Println(formatted)",
      `fmt.Println("lossless:", formatted == input)`,
    ],
  );
}

/** Go sample for `GET /api/now` — FormatNano(time.Now(), ext). */
export function goNowSample({ timeZone, calendar }: NowSampleArgs): string {
  // gofmt aligns composite-literal values; pad "Tags:" only when both fields exist.
  const tagsField = (pad: string) =>
    `Tags:${pad}map[string]string{ixdtf.ExtensionUnicodeCalendar: ${quote(calendar ?? "")}},`;
  const utc = timeZone === "" || timeZone === "UTC";
  const extLines = utc
    ? calendar
      ? [
          "ext := ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{",
          `\t${tagsField(" ")}`,
          "})",
        ]
      : ["ext := ixdtf.NewIXDTFExtensions(nil)"]
    : [
        "ext := ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{",
        "\tLocation: loc,",
        ...(calendar ? [`\t${tagsField("     ")}`] : []),
        "})",
      ];
  const body = utc
    ? [...extLines, "formatted, err := ixdtf.FormatNano(time.Now().UTC(), ext)"]
    : [
        `loc, err := time.LoadLocation(${quote(timeZone)})`,
        "if err != nil {",
        `\tfmt.Println("error:", err)`,
        "\treturn",
        "}",
        ...extLines,
        "formatted, err := ixdtf.FormatNano(time.Now().In(loc), ext)",
      ];
  return goProgram(
    ['"fmt"', '"time"', "", '"github.com/8beeeaaat/ixdtf"'],
    [
      ...body,
      "if err != nil {",
      `\tfmt.Println("format error:", err)`,
      "\treturn",
      "}",
      "fmt.Println(formatted)",
    ],
  );
}

export interface DstOverlapSampleArgs {
  earlier: string;
  later: string;
}

/** Go sample for F-6-1 — parse both overlap strings and print the real-time delta. */
export function goDstOverlapSample({ earlier, later }: DstOverlapSampleArgs): string {
  return goProgram(
    ['"fmt"', "", '"github.com/8beeeaaat/ixdtf"'],
    [
      `earlier, _, err := ixdtf.Parse(${quote(earlier)}, true)`,
      "if err != nil {",
      `\tfmt.Println("parse error:", err)`,
      "\treturn",
      "}",
      `later, _, err := ixdtf.Parse(${quote(later)}, true)`,
      "if err != nil {",
      `\tfmt.Println("parse error:", err)`,
      "\treturn",
      "}",
      "fmt.Println(earlier.UnixNano())",
      "fmt.Println(later.UnixNano())",
      `fmt.Println("delta:", later.Sub(earlier))`,
    ],
  );
}

export interface ZoneArithmeticSampleArgs {
  input: string;
}

/**
 * Go sample for F-6-3 — `Parse` applies the annotated IANA location to the
 * returned time, so `AddDate` is wall-clock arithmetic in that zone while
 * `Add(24 * time.Hour)` stays exact-time arithmetic.
 */
export function goZoneArithmeticSample({ input }: ZoneArithmeticSampleArgs): string {
  return goProgram(
    ['"fmt"', '"time"', "", '"github.com/8beeeaaat/ixdtf"'],
    [
      `parsedTime, ext, err := ixdtf.Parse(${quote(input)}, true)`,
      "if err != nil {",
      `\tfmt.Println("parse error:", err)`,
      "\treturn",
      "}",
      "",
      "plusDay, err := ixdtf.FormatNano(parsedTime.AddDate(0, 0, 1), ext)",
      "if err != nil {",
      `\tfmt.Println("format error:", err)`,
      "\treturn",
      "}",
      "plusHours, err := ixdtf.FormatNano(parsedTime.Add(24*time.Hour), ext)",
      "if err != nil {",
      `\tfmt.Println("format error:", err)`,
      "\treturn",
      "}",
      `fmt.Println("+1 day:  ", plusDay)`,
      `fmt.Println("+24 hours:", plusHours)`,
    ],
  );
}

export interface CalendarProjectionSampleArgs {
  input: string;
}

/**
 * Go sample for F-6-4 — Go does not interpret `u-ca`; it carries the tag
 * losslessly through Parse → FormatNano while the instant stays unchanged.
 */
export function goCalendarProjectionSample({ input }: CalendarProjectionSampleArgs): string {
  return goProgram(
    ['"fmt"', "", '"github.com/8beeeaaat/ixdtf"'],
    [
      `input := ${quote(input)}`,
      "",
      "parsedTime, ext, err := ixdtf.Parse(input, false)",
      "if err != nil {",
      `\tfmt.Println("parse error:", err)`,
      "\treturn",
      "}",
      "fmt.Println(parsedTime.UnixNano())",
      `fmt.Println("u-ca:", ext.Tags[ixdtf.ExtensionUnicodeCalendar])`,
      "",
      "formatted, err := ixdtf.FormatNano(parsedTime, ext)",
      "if err != nil {",
      `\tfmt.Println("format error:", err)`,
      "\treturn",
      "}",
      "fmt.Println(formatted)",
      `fmt.Println("lossless:", formatted == input)`,
    ],
  );
}

/**
 * JavaScript sample mirroring browserParse: `ZonedDateTime.from` when the
 * input carries a `[...]` annotation, otherwise `Instant.from`.
 */
export function jsParseSample(input: string): string {
  const lines = input.includes("[")
    ? [
        "const zdt = Temporal.ZonedDateTime.from(input);",
        "console.log(zdt.epochNanoseconds);",
        "console.log(zdt.timeZoneId, zdt.calendarId);",
      ]
    : ["const instant = Temporal.Instant.from(input);", "console.log(instant.epochNanoseconds);"];
  return [`const input = ${quote(input)};`, "", ...lines, ""].join("\n");
}

/** JavaScript sample for the roundtrip comparison: parse then `.toString()`. */
export function jsRoundtripSample(input: string): string {
  const parseLine = input.includes("[")
    ? "const formatted = Temporal.ZonedDateTime.from(input).toString();"
    : "const formatted = Temporal.Instant.from(input).toString();";
  return [
    `const input = ${quote(input)};`,
    "",
    parseLine,
    "console.log(formatted);",
    'console.log("lossless:", formatted === input);',
    "",
  ].join("\n");
}

/** JavaScript sample for the browser-side "now" counterpart. */
export function jsNowSample({ timeZone, calendar }: NowSampleArgs): string {
  const chain = calendar ? `.withCalendar(${quote(calendar)})` : "";
  return [
    `const zdt = Temporal.Now.zonedDateTimeISO(${quote(timeZone)})${chain};`,
    "console.log(zdt.toString());",
    "",
  ].join("\n");
}

/** JavaScript sample for F-6-1 — resolve both overlap strings and diff them. */
export function jsDstOverlapSample({ earlier, later }: DstOverlapSampleArgs): string {
  return [
    `const earlier = Temporal.ZonedDateTime.from(${quote(earlier)});`,
    `const later = Temporal.ZonedDateTime.from(${quote(later)});`,
    "",
    "console.log(earlier.epochNanoseconds);",
    "console.log(later.epochNanoseconds);",
    'console.log("delta:", later.since(earlier).toString());',
    "",
  ].join("\n");
}

/** JavaScript sample for F-6-3 — calendar vs exact-time arithmetic on one ZonedDateTime. */
export function jsZoneArithmeticSample(input: string): string {
  return [
    `const zdt = Temporal.ZonedDateTime.from(${quote(input)});`,
    "",
    "console.log(zdt.add({ days: 1 }).toString());",
    "console.log(zdt.add({ hours: 24 }).toString());",
    "",
  ].join("\n");
}

/** JavaScript sample for F-6-4 — one instant projected into several calendars. */
export function jsCalendarProjectionSample(input: string, calendars: readonly string[]): string {
  return [
    `const zdt = Temporal.ZonedDateTime.from(${quote(input)});`,
    "",
    `for (const calendar of [${calendars.map(quote).join(", ")}]) {`,
    "  const projected = zdt.withCalendar(calendar);",
    "  console.log(projected.toString());",
    "  console.log(calendar, projected.era, projected.eraYear ?? projected.year);",
    "}",
    'console.log("instant:", zdt.epochNanoseconds);',
    "",
  ].join("\n");
}
