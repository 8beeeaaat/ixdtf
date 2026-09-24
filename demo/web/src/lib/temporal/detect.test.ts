import { describe, expect, it } from "vitest";
import { getNativeTemporal, getPolyfillTemporal, getTemporal, requireTemporal } from "./detect";

describe("temporal detect", () => {
  it("exposes the Temporal surface from the bundled polyfill", () => {
    const t = getPolyfillTemporal();
    expect(typeof t.ZonedDateTime.from).toBe("function");
    expect(typeof t.Instant.from).toBe("function");
    expect(typeof t.Now.zonedDateTimeISO).toBe("function");
  });

  it("prefers native Temporal and falls back to the polyfill", () => {
    // The jsdom unit environment has no native Temporal; browsers with it are
    // covered by the browser-project stories.
    if (getNativeTemporal() === null) {
      expect(getTemporal()).toBe(getPolyfillTemporal());
    } else {
      expect(getTemporal()).toBe(getNativeTemporal());
    }
    // The live-UI helper never returns null (the polyfill is always available).
    expect(requireTemporal()).toBe(getTemporal());
  });
});
