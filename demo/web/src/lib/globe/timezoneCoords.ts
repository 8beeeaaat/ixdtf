/**
 * IANA time zone → approximate latitude/longitude, plus its current UTC offset.
 *
 * Display-only (D-10): these coordinates place the home globe's marker and pick
 * the front-facing meridian. They are deliberately approximate and must never be
 * used for correctness — the same line drawn by `lib/ixdtf/tokenize.ts`.
 */

import { ZONE_TAB_COORDS } from "./zoneTab";

export interface LatLng {
  lat: number;
  lng: number;
}

/**
 * Representative coordinates (the zone's principal city) for common IANA zones.
 * Not exhaustive: the long tail is covered by {@link fallbackFromRegion}.
 */
const ZONE_COORDS: Record<string, LatLng> = {
  UTC: { lat: 0, lng: 0 },
  "Etc/UTC": { lat: 0, lng: 0 },
  // Asia
  "Asia/Tokyo": { lat: 35.68, lng: 139.69 },
  "Asia/Seoul": { lat: 37.57, lng: 126.98 },
  "Asia/Shanghai": { lat: 31.23, lng: 121.47 },
  "Asia/Hong_Kong": { lat: 22.32, lng: 114.17 },
  "Asia/Taipei": { lat: 25.03, lng: 121.57 },
  "Asia/Singapore": { lat: 1.35, lng: 103.82 },
  "Asia/Bangkok": { lat: 13.76, lng: 100.5 },
  "Asia/Jakarta": { lat: -6.21, lng: 106.85 },
  "Asia/Manila": { lat: 14.6, lng: 120.98 },
  "Asia/Kolkata": { lat: 22.57, lng: 88.36 },
  "Asia/Kathmandu": { lat: 27.71, lng: 85.32 },
  "Asia/Dhaka": { lat: 23.81, lng: 90.41 },
  "Asia/Karachi": { lat: 24.86, lng: 67.0 },
  "Asia/Dubai": { lat: 25.2, lng: 55.27 },
  "Asia/Tehran": { lat: 35.69, lng: 51.39 },
  "Asia/Jerusalem": { lat: 31.77, lng: 35.21 },
  "Asia/Riyadh": { lat: 24.71, lng: 46.68 },
  "Asia/Yekaterinburg": { lat: 56.84, lng: 60.61 },
  "Asia/Omsk": { lat: 54.99, lng: 73.37 },
  "Asia/Novosibirsk": { lat: 55.01, lng: 82.93 },
  "Asia/Krasnoyarsk": { lat: 56.01, lng: 92.87 },
  "Asia/Irkutsk": { lat: 52.29, lng: 104.3 },
  "Asia/Yakutsk": { lat: 62.03, lng: 129.73 },
  "Asia/Vladivostok": { lat: 43.12, lng: 131.89 },
  "Asia/Sakhalin": { lat: 46.96, lng: 142.74 },
  "Asia/Magadan": { lat: 59.56, lng: 150.8 },
  "Asia/Kamchatka": { lat: 53.02, lng: 158.65 },
  "Asia/Makassar": { lat: -5.15, lng: 119.43 },
  "Asia/Jayapura": { lat: -2.53, lng: 140.72 },
  "Asia/Pyongyang": { lat: 39.02, lng: 125.75 },
  "Asia/Ulaanbaatar": { lat: 47.92, lng: 106.92 },
  "Asia/Almaty": { lat: 43.24, lng: 76.95 },
  "Asia/Aqtobe": { lat: 50.28, lng: 57.17 },
  "Asia/Tashkent": { lat: 41.3, lng: 69.24 },
  "Asia/Bishkek": { lat: 42.87, lng: 74.59 },
  "Asia/Dushanbe": { lat: 38.56, lng: 68.79 },
  "Asia/Ashgabat": { lat: 37.95, lng: 58.38 },
  "Asia/Kabul": { lat: 34.53, lng: 69.17 },
  "Asia/Yangon": { lat: 16.87, lng: 96.2 },
  "Asia/Colombo": { lat: 6.93, lng: 79.85 },
  "Asia/Thimphu": { lat: 27.47, lng: 89.64 },
  "Asia/Baghdad": { lat: 33.31, lng: 44.36 },
  "Asia/Amman": { lat: 31.95, lng: 35.93 },
  "Asia/Beirut": { lat: 33.89, lng: 35.5 },
  "Asia/Kuwait": { lat: 29.38, lng: 47.98 },
  "Asia/Qatar": { lat: 25.29, lng: 51.53 },
  "Asia/Baku": { lat: 40.41, lng: 49.87 },
  "Asia/Yerevan": { lat: 40.18, lng: 44.51 },
  "Asia/Tbilisi": { lat: 41.72, lng: 44.79 },
  // Europe
  "Europe/London": { lat: 51.51, lng: -0.13 },
  "Europe/Dublin": { lat: 53.35, lng: -6.26 },
  "Europe/Lisbon": { lat: 38.72, lng: -9.14 },
  "Europe/Madrid": { lat: 40.42, lng: -3.7 },
  "Europe/Paris": { lat: 48.86, lng: 2.35 },
  "Europe/Berlin": { lat: 52.52, lng: 13.4 },
  "Europe/Rome": { lat: 41.9, lng: 12.5 },
  "Europe/Amsterdam": { lat: 52.37, lng: 4.9 },
  "Europe/Zurich": { lat: 47.38, lng: 8.54 },
  "Europe/Warsaw": { lat: 52.23, lng: 21.01 },
  "Europe/Athens": { lat: 37.98, lng: 23.73 },
  "Europe/Istanbul": { lat: 41.01, lng: 28.98 },
  "Europe/Moscow": { lat: 55.76, lng: 37.62 },
  "Europe/Kaliningrad": { lat: 54.71, lng: 20.51 },
  "Europe/Samara": { lat: 53.2, lng: 50.15 },
  "Europe/Kyiv": { lat: 50.45, lng: 30.52 },
  "Europe/Vienna": { lat: 48.21, lng: 16.37 },
  "Europe/Prague": { lat: 50.09, lng: 14.42 },
  "Europe/Budapest": { lat: 47.5, lng: 19.04 },
  "Europe/Brussels": { lat: 50.85, lng: 4.35 },
  "Europe/Stockholm": { lat: 59.33, lng: 18.06 },
  "Europe/Oslo": { lat: 59.91, lng: 10.75 },
  "Europe/Copenhagen": { lat: 55.68, lng: 12.57 },
  "Europe/Helsinki": { lat: 60.17, lng: 24.94 },
  "Europe/Riga": { lat: 56.95, lng: 24.11 },
  "Europe/Vilnius": { lat: 54.69, lng: 25.28 },
  "Europe/Tallinn": { lat: 59.44, lng: 24.75 },
  "Europe/Bucharest": { lat: 44.43, lng: 26.1 },
  "Europe/Sofia": { lat: 42.7, lng: 23.32 },
  "Europe/Belgrade": { lat: 44.79, lng: 20.45 },
  "Atlantic/Reykjavik": { lat: 64.15, lng: -21.94 },
  // Africa
  "Africa/Cairo": { lat: 30.04, lng: 31.24 },
  "Africa/Lagos": { lat: 6.52, lng: 3.38 },
  "Africa/Nairobi": { lat: -1.29, lng: 36.82 },
  "Africa/Johannesburg": { lat: -26.2, lng: 28.05 },
  "Africa/Casablanca": { lat: 33.57, lng: -7.59 },
  "Africa/Algiers": { lat: 36.75, lng: 3.06 },
  "Africa/Tunis": { lat: 36.81, lng: 10.18 },
  "Africa/Tripoli": { lat: 32.89, lng: 13.19 },
  "Africa/Khartoum": { lat: 15.5, lng: 32.56 },
  "Africa/Addis_Ababa": { lat: 9.03, lng: 38.74 },
  "Africa/Accra": { lat: 5.6, lng: -0.19 },
  "Africa/Abidjan": { lat: 5.35, lng: -4.03 },
  "Africa/Kinshasa": { lat: -4.44, lng: 15.27 },
  "Africa/Lubumbashi": { lat: -11.66, lng: 27.48 },
  "Africa/Luanda": { lat: -8.84, lng: 13.23 },
  "Africa/Windhoek": { lat: -22.56, lng: 17.08 },
  "Africa/Maputo": { lat: -25.97, lng: 32.57 },
  "Africa/Dar_es_Salaam": { lat: -6.79, lng: 39.21 },
  "Indian/Antananarivo": { lat: -18.88, lng: 47.51 },
  // Americas
  "America/New_York": { lat: 40.71, lng: -74.01 },
  "America/Toronto": { lat: 43.65, lng: -79.38 },
  "America/Chicago": { lat: 41.88, lng: -87.63 },
  "America/Winnipeg": { lat: 49.9, lng: -97.14 },
  "America/Regina": { lat: 50.45, lng: -104.62 },
  "America/Denver": { lat: 39.74, lng: -104.99 },
  "America/Edmonton": { lat: 53.55, lng: -113.49 },
  "America/Phoenix": { lat: 33.45, lng: -112.07 },
  "America/Los_Angeles": { lat: 34.05, lng: -118.24 },
  "America/Vancouver": { lat: 49.28, lng: -123.12 },
  "America/Anchorage": { lat: 61.22, lng: -149.9 },
  "America/Adak": { lat: 51.88, lng: -176.66 },
  "America/Halifax": { lat: 44.65, lng: -63.57 },
  "America/St_Johns": { lat: 47.56, lng: -52.71 },
  "America/Nuuk": { lat: 64.18, lng: -51.69 },
  "America/Mexico_City": { lat: 19.43, lng: -99.13 },
  "America/Guatemala": { lat: 14.63, lng: -90.51 },
  "America/Panama": { lat: 8.98, lng: -79.52 },
  "America/Havana": { lat: 23.13, lng: -82.38 },
  "America/Santo_Domingo": { lat: 18.49, lng: -69.94 },
  "America/Puerto_Rico": { lat: 18.47, lng: -66.11 },
  "America/Bogota": { lat: 4.71, lng: -74.07 },
  "America/Guayaquil": { lat: -2.17, lng: -79.92 },
  "America/Caracas": { lat: 10.48, lng: -66.9 },
  "America/Lima": { lat: -12.05, lng: -77.04 },
  "America/La_Paz": { lat: -16.5, lng: -68.15 },
  "America/Rio_Branco": { lat: -9.97, lng: -67.81 },
  "America/Manaus": { lat: -3.12, lng: -60.02 },
  "America/Sao_Paulo": { lat: -23.55, lng: -46.63 },
  "America/Noronha": { lat: -3.85, lng: -32.42 },
  "America/Asuncion": { lat: -25.3, lng: -57.64 },
  "America/Montevideo": { lat: -34.9, lng: -56.16 },
  "America/Argentina/Buenos_Aires": { lat: -34.6, lng: -58.38 },
  "America/Santiago": { lat: -33.45, lng: -70.67 },
  // Oceania
  "Australia/Sydney": { lat: -33.87, lng: 151.21 },
  "Australia/Melbourne": { lat: -37.81, lng: 144.96 },
  "Australia/Hobart": { lat: -42.88, lng: 147.33 },
  "Australia/Adelaide": { lat: -34.93, lng: 138.6 },
  "Australia/Darwin": { lat: -12.46, lng: 130.84 },
  "Australia/Perth": { lat: -31.95, lng: 115.86 },
  "Australia/Brisbane": { lat: -27.47, lng: 153.03 },
  "Pacific/Port_Moresby": { lat: -9.48, lng: 147.15 },
  "Pacific/Guadalcanal": { lat: -9.43, lng: 159.95 },
  "Pacific/Noumea": { lat: -22.28, lng: 166.46 },
  "Pacific/Guam": { lat: 13.47, lng: 144.75 },
  "Pacific/Majuro": { lat: 7.12, lng: 171.38 },
  "Pacific/Tarawa": { lat: 1.33, lng: 172.98 },
  "Pacific/Fiji": { lat: -18.14, lng: 178.44 },
  "Pacific/Tongatapu": { lat: -21.14, lng: -175.2 },
  "Pacific/Apia": { lat: -13.83, lng: -171.76 },
  "Pacific/Chatham": { lat: -43.95, lng: -176.56 },
  "Pacific/Auckland": { lat: -36.85, lng: 174.76 },
  "Pacific/Tahiti": { lat: -17.53, lng: -149.57 },
  "Pacific/Honolulu": { lat: 21.31, lng: -157.86 },
};

/** Central latitude band per IANA region prefix (used only by the fallback). */
const REGION_LAT: Record<string, number> = {
  Africa: 3,
  America: 25,
  Antarctica: -75,
  Arctic: 78,
  Asia: 33,
  Atlantic: 35,
  Australia: -27,
  Europe: 50,
  Indian: -12,
  Pacific: 0,
};

/**
 * Stable-ordered curated zone list, exported so the home globe can render
 * an always-visible marker (with hover label) per zone. Zone *fill* on the
 * globe no longer derives from these points — it uses real boundary data
 * (`@photostructure/tz-lookup`); these remain marker/label positions and the
 * {@link nearestTimeZone} click fallback.
 */
export const ZONE_ENTRIES: ReadonlyArray<readonly [string, LatLng]> = Object.entries(ZONE_COORDS);

const clampLng = (lng: number): number => Math.max(-180, Math.min(180, lng));

/**
 * Current UTC offset of a zone, via `Intl` (`GMT+09:00` → minutes + label).
 * Falls back to UTC when the runtime cannot resolve the zone.
 */
export function getUtcOffset(timeZone: string): { minutes: number; label: string } {
  try {
    const parts = new Intl.DateTimeFormat("en-US", {
      timeZone,
      timeZoneName: "longOffset",
    }).formatToParts(new Date());
    const name = parts.find((p) => p.type === "timeZoneName")?.value ?? "GMT";
    const match = name.match(/GMT([+-])(\d{1,2})(?::(\d{2}))?/);
    if (!match) return { minutes: 0, label: "UTC" };
    const sign = match[1] === "-" ? -1 : 1;
    const hours = Number(match[2]);
    const mins = Number(match[3] ?? "0");
    const minutes = sign * (hours * 60 + mins);
    const hh = String(hours).padStart(2, "0");
    const mm = String(mins).padStart(2, "0");
    return { minutes, label: minutes === 0 ? "UTC" : `UTC${match[1]}${hh}:${mm}` };
  } catch {
    return { minutes: 0, label: "UTC" };
  }
}

/**
 * Longitude from the zone's current UTC offset (15° per hour) paired with a
 * latitude guessed from the region prefix. A rough placement for zones absent
 * from {@link ZONE_COORDS}; good enough for a decorative background.
 */
function fallbackFromRegion(timeZone: string): LatLng {
  const region = timeZone.split("/")[0] ?? "";
  const lat = REGION_LAT[region] ?? 0;
  const lng = clampLng(getUtcOffset(timeZone).minutes / 4);
  return { lat, lng };
}

/**
 * Real (non-synthetic) coordinates for a zone: the hand-curated city table
 * first, then the IANA tzdb zone.tab representative point. Undefined when
 * only the rough region/offset fallback would remain.
 */
export function knownZoneCoords(timeZone: string): LatLng | undefined {
  return ZONE_COORDS[timeZone] ?? ZONE_TAB_COORDS[timeZone];
}

/** Approximate lat/lng for an IANA time zone (display-only, D-10). */
export function timeZoneToLatLng(timeZone: string): LatLng {
  return knownZoneCoords(timeZone) ?? fallbackFromRegion(timeZone);
}

/**
 * Reverse lookup: the curated zone whose city is nearest to a clicked lat/lng.
 * Used when the user clicks the globe to pick a zone (display-only, coarse — the
 * fine-grained picker still lists every IANA zone).
 */
export function nearestTimeZone(lat: number, lng: number): string {
  const rad = Math.PI / 180;
  const φ1 = lat * rad;
  const λ1 = lng * rad;
  let best = "UTC";
  let bestCos = -Infinity;
  for (const [zone, c] of ZONE_ENTRIES) {
    const φ2 = c.lat * rad;
    // cosine of the great-circle angle; larger = closer
    const cos =
      Math.sin(φ1) * Math.sin(φ2) + Math.cos(φ1) * Math.cos(φ2) * Math.cos(λ1 - c.lng * rad);
    if (cos > bestCos) {
      bestCos = cos;
      best = zone;
    }
  }
  return best;
}
