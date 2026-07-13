import { useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { CalendarPicker } from "@/components/CalendarPicker";
import { IxdtfHighlight } from "@/components/IxdtfHighlight";
import { TimeZonePicker } from "@/components/TimeZonePicker";
import { Button } from "@/components/ui/Button";
import { Card, CardContent } from "@/components/ui/Card";
import { FormField } from "@/components/ui/FormField";
import { Input } from "@/components/ui/Input";
import { requireTemporal } from "@/lib/temporal/detect";
import { useDebounce } from "@/lib/useDebounce";
import { useLocalStorage } from "@/lib/useLocalStorage";
import { useNowIXDTF } from "../home/useNow";

const DEFAULT_ZONES = ["Asia/Tokyo", "UTC", "America/Los_Angeles", "Europe/London"];

/** F-4: world clock — one instant across time zones (`withTimeZone`) and calendars (`withCalendar`). */
export function ConverterPage() {
  const { t, i18n } = useTranslation();
  const temporal = requireTemporal();
  const nowIXDTF = useNowIXDTF();
  const [baseInput, setBaseInput] = useState(nowIXDTF);
  const [calendar, setCalendar] = useState("iso8601");
  const [zones, setZones] = useLocalStorage<string[]>("converter.timeZones", DEFAULT_ZONES);
  const [pendingZone, setPendingZone] = useState("UTC");
  const [nowValue, setNowValue] = useState(() => temporal.Now.zonedDateTimeISO("UTC"));
  const debounced = useDebounce(baseInput, 300);

  // F-4-4: 基準時刻は IXDTF 文字列で指定できる (空欄なら現在時刻)
  const parsedBase = useMemo(() => {
    if (!debounced) {
      return null;
    }
    try {
      return temporal.ZonedDateTime.from(debounced);
    } catch {
      try {
        return temporal.Instant.from(debounced).toZonedDateTimeISO("UTC");
      } catch {
        return null;
      }
    }
  }, [temporal, debounced]);

  const base = debounced ? parsedBase : nowValue;
  const invalidBase = debounced.length > 0 && parsedBase === null;

  const addZone = () => {
    if (pendingZone && !zones.includes(pendingZone)) {
      setZones([...zones, pendingZone]);
    }
  };

  return (
    <div className="mx-auto max-w-6xl space-y-12 px-6">
      <header>
        <h1 className="font-sans font-semibold text-2xl tracking-tight">{t("converter.title")}</h1>
        <p className="mt-1 font-sans text-muted-foreground text-sm">{t("converter.tagline")}</p>
      </header>

      <section className="space-y-4">
        <FormField
          label={t("converter.baseTime")}
          htmlFor="converter-base"
          hint={t("converter.baseTimeHint")}
          error={invalidBase ? t("converter.invalidBase") : undefined}
        >
          <div className="flex gap-2">
            <Input
              id="converter-base"
              value={baseInput}
              onChange={(event) => setBaseInput(event.target.value)}
              placeholder={t("converter.baseTimePlaceholder")}
              autoComplete="off"
              spellCheck={false}
            />
            <Button
              variant="secondary"
              onClick={() => {
                setBaseInput(nowIXDTF);
                setNowValue(temporal.Now.zonedDateTimeISO("UTC"));
              }}
            >
              {t("converter.useNow")}
            </Button>
          </div>
        </FormField>
        <div className="grid gap-4 md:grid-cols-2">
          <FormField label={t("converter.calendarLabel")}>
            <CalendarPicker value={calendar} onChange={setCalendar} />
          </FormField>
          <FormField label={t("converter.addTimeZone")}>
            <div className="flex gap-2">
              <TimeZonePicker value={pendingZone} onChange={setPendingZone} className="flex-1" />
              <Button variant="secondary" onClick={addZone}>
                {t("common.add")}
              </Button>
            </div>
          </FormField>
        </div>
      </section>

      {zones.length === 0 ? (
        <p className="font-sans text-muted-foreground text-sm">{t("converter.empty")}</p>
      ) : (
        base && (
          <Card>
            <CardContent className="divide-y divide-border px-6 py-2">
              {zones.map((zone) => {
                const zoned = base.withTimeZone(zone);
                const ixdtf =
                  calendar === "iso8601"
                    ? zoned.toString()
                    : zoned.withCalendar(calendar).toString();
                // iso8601 を Intl に渡すと CLDR の汎用パターン ("2026 7月 12, 日曜日" 等) に
                // なるため、既定カレンダーのときはロケール標準の暦・書式で表示する
                const human = new Intl.DateTimeFormat(i18n.language, {
                  timeZone: zone,
                  ...(calendar !== "iso8601" && { calendar }),
                  dateStyle: "full",
                  timeStyle: "medium",
                }).format(new Date(zoned.epochMilliseconds));
                return (
                  <div key={zone} className="space-y-1 py-4">
                    <div className="flex items-baseline justify-between gap-4">
                      <p className="font-medium font-mono text-sm tabular-nums">{zone}</p>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => setZones(zones.filter((z) => z !== zone))}
                      >
                        {t("common.remove")}
                      </Button>
                    </div>
                    <p className="font-sans text-sm">{human}</p>
                    <IxdtfHighlight value={ixdtf} className="text-sm" />
                  </div>
                );
              })}
            </CardContent>
          </Card>
        )
      )}
    </div>
  );
}
