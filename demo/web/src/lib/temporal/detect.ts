/**
 * The runtime shape of the global `Temporal` namespace object.
 *
 * Expressed as a pure `typeof import(...)` type query so that NO runtime import
 * of `temporal-polyfill` is emitted — the polyfill is a type-only dependency
 * (N-1 / D-2: the app showcases the browser's NATIVE Temporal and must never
 * load a polyfill). All code reaches the API through {@link getTemporal}.
 */
export type TemporalApi = typeof import("temporal-polyfill").Temporal;

/** Read the native `globalThis.Temporal`, or `null` when unavailable. */
export function getTemporal(): TemporalApi | null {
  const scope = globalThis as typeof globalThis & { Temporal?: TemporalApi };
  return scope.Temporal ?? null;
}

/** Feature-detection gate used by `main.tsx` (F-0-4). */
export function isTemporalAvailable(): boolean {
  return getTemporal() !== null;
}

/**
 * Return the native Temporal API, throwing if absent. Safe to call from
 * components: the app only mounts them after {@link isTemporalAvailable}
 * passed in `main.tsx`, so this never throws in practice.
 */
export function requireTemporal(): TemporalApi {
  const temporal = getTemporal();
  if (!temporal) {
    throw new Error("Temporal API is not available in this browser.");
  }
  return temporal;
}
