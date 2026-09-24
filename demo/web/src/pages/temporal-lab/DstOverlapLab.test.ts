import { describe, expect, it } from "vitest";
import { epochDeltaSeconds } from "./DstOverlapLab";

describe("epochDeltaSeconds", () => {
  it("reports the one-hour difference between the two overlap instants", () => {
    expect(epochDeltaSeconds("1541320200000000000", "1541323800000000000")).toBe(3600);
  });

  it("requires both Temporal observations", () => {
    expect(epochDeltaSeconds(null, "1541323800000000000")).toBeNull();
  });
});
