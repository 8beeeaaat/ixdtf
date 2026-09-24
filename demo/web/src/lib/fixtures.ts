import samplesJson from "../../../testdata/ixdtf_samples.json";

/** Expectation flags recorded in the shared fixture (testdata/ixdtf_samples.json). */
export type SampleOutcome = "ok" | "error";

export interface IxdtfSample {
  id: string;
  input: string;
  category: string;
  strict: boolean;
  expect_server: string;
  expect_browser: string;
  expect_error_substr?: string;
  note_key: string;
  interop_preset?: boolean;
  guide: boolean;
}

/** All shared samples (F-5 Guide + Vitest read the same source of truth). */
export const samples: IxdtfSample[] = samplesJson.samples as IxdtfSample[];

/** Resolve a fixture used as a named demo building block; fail loudly on catalog drift. */
export function requiredSample(id: string): IxdtfSample {
  const sample = samples.find((candidate) => candidate.id === id);
  if (!sample) {
    throw new Error(`Missing IXDTF fixture: ${id}`);
  }
  return sample;
}

/** Samples flagged for the learning guide (F-5-2: the curated 10–15). */
export const guideSamples: IxdtfSample[] = samples.filter((sample) => sample.guide);

/**
 * Samples whose two implementations disagree, plus the outright-rejected ones —
 * the Interop presets that make the behavioural differences the content (F-3-3).
 */
export function isInteropPreset(sample: IxdtfSample): boolean {
  return (
    sample.interop_preset !== false &&
    (sample.expect_server !== sample.expect_browser || sample.category === "rejected")
  );
}

export const interopPresets: IxdtfSample[] = samples.filter(isInteropPreset);

const CATEGORY_ORDER = [
  "basic",
  "timezone",
  "calendar",
  "critical",
  "composite",
  "rejected",
] as const;

/** Guide samples grouped by category in a stable presentation order. */
export function guideSamplesByCategory(): { category: string; items: IxdtfSample[] }[] {
  return CATEGORY_ORDER.map((category) => ({
    category,
    items: guideSamples.filter((sample) => sample.category === category),
  })).filter((group) => group.items.length > 0);
}
