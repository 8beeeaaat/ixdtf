import { describe, expect, it } from "vitest";
import { interopPresets, isInteropPreset, requiredSample, samples } from "@/lib/fixtures";

describe("shared IXDTF fixture routing", () => {
  it("keeps fixture IDs unique", () => {
    expect(new Set(samples.map((sample) => sample.id)).size).toBe(samples.length);
  });

  it("keeps Temporal Lab samples out of Interop presets", () => {
    const presetIds = interopPresets.map((sample) => sample.id);
    expect(presetIds).not.toContain("dst-overlap-earlier");
    expect(presetIds).not.toContain("dst-overlap-later");
  });

  it("honors the explicit Interop exclusion even if expectations later diverge", () => {
    const labSample = requiredSample("dst-overlap-earlier");
    expect(labSample.interop_preset).toBe(false);
    expect(
      isInteropPreset({
        ...labSample,
        category: "rejected",
        expect_browser: "error",
      }),
    ).toBe(false);
    expect(
      isInteropPreset({
        ...labSample,
        interop_preset: true,
        category: "rejected",
        expect_browser: "error",
      }),
    ).toBe(true);
  });
});
