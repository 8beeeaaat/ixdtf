/**
 * Subsolar point — the lat/lng where the sun is directly overhead at an instant.
 * Used to shade the globe's night hemisphere (the day/night terminator).
 *
 * Display-only (D-10): a low-order approximation (equation of time ignored, so
 * up to ~4° longitude error). Fine for a decorative terminator; never used for
 * correctness.
 */

const DEG = Math.PI / 180;

/** Day of the year (1–366) in UTC. */
function dayOfYearUTC(date: Date): number {
  const start = Date.UTC(date.getUTCFullYear(), 0, 0);
  return Math.floor((date.getTime() - start) / 86_400_000);
}

/**
 * Subsolar point for `date`.
 *
 * - latitude ≈ solar declination (−23.44° at Dec solstice … +23.44° at Jun)
 * - longitude = the meridian at local solar noon: 15° × (12 − UTC hours)
 */
export function subsolarPoint(date: Date): { lat: number; lng: number } {
  const n = dayOfYearUTC(date);
  const lat = -23.44 * Math.cos(DEG * ((360 / 365) * (n + 10)));

  const utcHours = date.getUTCHours() + date.getUTCMinutes() / 60 + date.getUTCSeconds() / 3600;
  let lng = 15 * (12 - utcHours);
  if (lng > 180) lng -= 360;
  if (lng < -180) lng += 360;

  return { lat, lng };
}
