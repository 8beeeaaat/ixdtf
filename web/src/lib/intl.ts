type SupportedValuesKey =
  | "calendar"
  | "collation"
  | "currency"
  | "numberingSystem"
  | "timeZone"
  | "unit";

type SupportedValuesOf = (key: SupportedValuesKey) => string[];

const intlWithSupported = Intl as typeof Intl & {
  supportedValuesOf?: SupportedValuesOf;
};

/** Enumerate IANA time zones (`Intl.supportedValuesOf`), with a small fallback. */
export function supportedTimeZones(): string[] {
  return (
    intlWithSupported.supportedValuesOf?.("timeZone") ?? [
      "UTC",
      "Asia/Tokyo",
      "America/Los_Angeles",
      "Europe/London",
    ]
  );
}

/** Enumerate calendar systems (`Intl.supportedValuesOf`), with a small fallback. */
export function supportedCalendars(): string[] {
  return (
    intlWithSupported.supportedValuesOf?.("calendar") ?? [
      "iso8601",
      "gregory",
      "japanese",
      "hebrew",
      "islamic",
      "chinese",
    ]
  );
}
