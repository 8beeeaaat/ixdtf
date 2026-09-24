// `temporal-polyfill/full` bundles the non-ISO Intl calendars (japanese, hebrew,
// islamic-umalqura, …). The slim default entry only supports ISO 8601 and throws
// "Unknown calendar japanese" on `[u-ca=japanese]` inputs — which surfaces on
// browsers without native Temporal (e.g. Safari), where this polyfill runs.
import { Temporal as PolyfillTemporal } from "temporal-polyfill/full";

/**
 * The runtime shape of the `Temporal` namespace object.
 *
 * The bundled `temporal-polyfill` provides the canonical type; the browser's
 * native `globalThis.Temporal` is assignable to it. Both implementations are
 * reached through {@link getTemporal}, which prefers native and falls back to
 * the polyfill.
 */
export type TemporalApi = typeof PolyfillTemporal;

/** Read the native `globalThis.Temporal`, or `null` when unavailable. */
export function getNativeTemporal(): TemporalApi | null {
  const scope = globalThis as typeof globalThis & { Temporal?: TemporalApi };
  return scope.Temporal ?? null;
}

/** The bundled `temporal-polyfill` implementation (always available). */
export function getPolyfillTemporal(): TemporalApi {
  return PolyfillTemporal;
}

/**
 * The Temporal implementation the live UI runs on (Home clock, Converter):
 * the browser's native Temporal when present, otherwise the bundled
 * `temporal-polyfill`. The comparison screens (Playground / Interop / Temporal
 * Lab) bypass this and drive both implementations explicitly, so there is no
 * user-selectable mode.
 */
export function getTemporal(): TemporalApi {
  return getNativeTemporal() ?? getPolyfillTemporal();
}

/**
 * Alias of {@link getTemporal} for callers that want an explicit "must exist"
 * contract. The polyfill is always available, so this never fails.
 */
export function requireTemporal(): TemporalApi {
  return getTemporal();
}
