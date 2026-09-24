import { describe, expect, it } from "vitest";
import { samples } from "@/lib/fixtures";
import { type IxdtfToken, tokenize } from "./tokenize";

// Explicit expected decomposition for every fixture input (F-1-2 classification
// coverage). Keyed by sample input so the two identical offset-mismatch inputs
// (lenient / strict) share one expectation.
const EXPECTED: Record<string, IxdtfToken[]> = {
  "1996-12-19T16:39:57-08:00": [
    { type: "date", text: "1996-12-19" },
    { type: "time", text: "T16:39:57" },
    { type: "offset", text: "-08:00" },
  ],
  "2026-07-07T14:30:00Z": [
    { type: "date", text: "2026-07-07" },
    { type: "time", text: "T14:30:00" },
    { type: "offset", text: "Z" },
  ],
  "2026-07-07T14:30:00-00:00": [
    { type: "date", text: "2026-07-07" },
    { type: "time", text: "T14:30:00" },
    { type: "offset", text: "-00:00" },
  ],
  "2026-07-07T23:30:00.123456789+09:00[Asia/Tokyo]": [
    { type: "date", text: "2026-07-07" },
    { type: "time", text: "T23:30:00.123456789" },
    { type: "offset", text: "+09:00" },
    { type: "timezone", text: "[Asia/Tokyo]", critical: false },
  ],
  "2026-07-07T23:30:00+09:00[Asia/Tokyo]": [
    { type: "date", text: "2026-07-07" },
    { type: "time", text: "T23:30:00" },
    { type: "offset", text: "+09:00" },
    { type: "timezone", text: "[Asia/Tokyo]", critical: false },
  ],
  "2018-11-04T01:30:00-07:00[America/Los_Angeles]": [
    { type: "date", text: "2018-11-04" },
    { type: "time", text: "T01:30:00" },
    { type: "offset", text: "-07:00" },
    { type: "timezone", text: "[America/Los_Angeles]", critical: false },
  ],
  "2018-11-04T01:30:00-08:00[America/Los_Angeles]": [
    { type: "date", text: "2018-11-04" },
    { type: "time", text: "T01:30:00" },
    { type: "offset", text: "-08:00" },
    { type: "timezone", text: "[America/Los_Angeles]", critical: false },
  ],
  "2026-07-07T23:30:00+02:00[Asia/Tokyo]": [
    { type: "date", text: "2026-07-07" },
    { type: "time", text: "T23:30:00" },
    { type: "offset", text: "+02:00" },
    { type: "timezone", text: "[Asia/Tokyo]", critical: false },
  ],
  "2026-07-07T23:30:00+09:00[Asia/Tokyo][u-ca=japanese]": [
    { type: "date", text: "2026-07-07" },
    { type: "time", text: "T23:30:00" },
    { type: "offset", text: "+09:00" },
    { type: "timezone", text: "[Asia/Tokyo]", critical: false },
    { type: "extension", text: "[u-ca=japanese]", critical: false },
  ],
  "1996-12-19T16:39:57-08:00[America/Los_Angeles][u-ca=hebrew]": [
    { type: "date", text: "1996-12-19" },
    { type: "time", text: "T16:39:57" },
    { type: "offset", text: "-08:00" },
    { type: "timezone", text: "[America/Los_Angeles]", critical: false },
    { type: "extension", text: "[u-ca=hebrew]", critical: false },
  ],
  "2026-07-07T23:30:00+09:00[!Asia/Tokyo]": [
    { type: "date", text: "2026-07-07" },
    { type: "time", text: "T23:30:00" },
    { type: "offset", text: "+09:00" },
    { type: "timezone", text: "[!Asia/Tokyo]", critical: true },
  ],
  "2026-07-07T23:30:00+09:00[Asia/Tokyo][!u-ca=japanese]": [
    { type: "date", text: "2026-07-07" },
    { type: "time", text: "T23:30:00" },
    { type: "offset", text: "+09:00" },
    { type: "timezone", text: "[Asia/Tokyo]", critical: false },
    { type: "extension", text: "[!u-ca=japanese]", critical: true },
  ],
  "2026-07-07T23:30:00+09:00[Asia/Tokyo][!foo=bar]": [
    { type: "date", text: "2026-07-07" },
    { type: "time", text: "T23:30:00" },
    { type: "offset", text: "+09:00" },
    { type: "timezone", text: "[Asia/Tokyo]", critical: false },
    { type: "extension", text: "[!foo=bar]", critical: true },
  ],
  "2026-07-07T23:30:00+09:00[Asia/Tokyo][foo=bar]": [
    { type: "date", text: "2026-07-07" },
    { type: "time", text: "T23:30:00" },
    { type: "offset", text: "+09:00" },
    { type: "timezone", text: "[Asia/Tokyo]", critical: false },
    { type: "extension", text: "[foo=bar]", critical: false },
  ],
  "2026-07-07T23:30:00+09:00[Asia/Tokyo][x-private=demo]": [
    { type: "date", text: "2026-07-07" },
    { type: "time", text: "T23:30:00" },
    { type: "offset", text: "+09:00" },
    { type: "timezone", text: "[Asia/Tokyo]", critical: false },
    { type: "extension", text: "[x-private=demo]", critical: false },
  ],
  "2026-07-07T23:30:00+09:00[Asia/Tokyo][_exp=demo]": [
    { type: "date", text: "2026-07-07" },
    { type: "time", text: "T23:30:00" },
    { type: "offset", text: "+09:00" },
    { type: "timezone", text: "[Asia/Tokyo]", critical: false },
    { type: "extension", text: "[_exp=demo]", critical: false },
  ],
  "not-a-timestamp": [{ type: "invalid", text: "not-a-timestamp" }],
};

describe("tokenize", () => {
  it("covers every shared fixture input", () => {
    // Guard: ensure the fixtures did not drift away from our expectations.
    const inputs = new Set(samples.map((sample) => sample.input));
    for (const input of inputs) {
      expect(EXPECTED[input], `missing expectation for ${input}`).toBeDefined();
    }
  });

  for (const sample of samples) {
    it(`decomposes ${sample.id}`, () => {
      const tokens = tokenize(sample.input);
      expect(tokens).toEqual(EXPECTED[sample.input]);
      // Round-trip invariant: concatenated token text equals the input.
      expect(tokens.map((token) => token.text).join("")).toBe(sample.input);
    });
  }

  it("classifies an offset-only bare timestamp without brackets", () => {
    expect(tokenize("2026-07-07T00:00:00+05:30")).toEqual([
      { type: "date", text: "2026-07-07" },
      { type: "time", text: "T00:00:00" },
      { type: "offset", text: "+05:30" },
    ]);
  });
});
