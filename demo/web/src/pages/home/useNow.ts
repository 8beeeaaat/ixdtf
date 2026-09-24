import { useEffect, useState } from "react";
import { requireTemporal, type TemporalApi } from "@/lib/temporal/detect";

type ZonedDateTime = ReturnType<TemporalApi["Now"]["zonedDateTimeISO"]>;

/**
 * Current time ticking at 100 ms (F-1-1). Only mount this inside the component
 * that displays the clock so the re-render stays contained (N-4).
 */
export function useNowTimer(timeZone: string): ZonedDateTime {
  const temporal = requireTemporal();
  const [now, setNow] = useState(() => temporal.Now.zonedDateTimeISO(timeZone));

  useEffect(() => {
    setNow(temporal.Now.zonedDateTimeISO(timeZone));
    const id = setInterval(() => setNow(temporal.Now.zonedDateTimeISO(timeZone)), 100);
    return () => clearInterval(id);
  }, [temporal, timeZone]);

  return now;
}

export function useNowIXDTF(): string {
  const temporal = requireTemporal();
  return temporal.Now.zonedDateTimeISO().toString();
}
