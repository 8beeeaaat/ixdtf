import { useTranslation } from "react-i18next";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card";
import { requireTemporal } from "@/lib/temporal/detect";
import { cn } from "@/lib/utils";

interface MonthCalendarProps {
  timeZone: string;
  calendar: string;
}

/**
 * Month grid rendered in the selected calendar system (F-1-5,
 * `Temporal.PlainDate` / `daysInMonth`). 和暦なら「令和8年7月」等の表記になる。
 */
export function MonthCalendar({ timeZone, calendar }: MonthCalendarProps) {
  const { t, i18n } = useTranslation();
  const temporal = requireTemporal();

  const today = temporal.Now.plainDateISO(timeZone).withCalendar(calendar);
  const first = today.with({ day: 1 });
  const leading = (first.dayOfWeek - 1) % 7;
  const days = Array.from({ length: today.daysInMonth }, (_, i) => i + 1);
  const weekdays = t("home.weekdays").split(",");

  const monthLabel = new Intl.DateTimeFormat(
    i18n.language,
    calendar === "iso8601" || calendar === "gregory"
      ? { year: "numeric", month: "long", calendar }
      : { era: "short", year: "numeric", month: "long", calendar },
  ).format(new Date(temporal.Now.zonedDateTimeISO(timeZone).epochMilliseconds));

  return (
    <Card>
      <CardHeader>
        <CardTitle>{t("home.calendarHeading")}</CardTitle>
        <p className="font-mono text-lg tabular-nums">{monthLabel}</p>
      </CardHeader>
      <CardContent>
        <div className="grid grid-cols-7 gap-y-1 text-center">
          {weekdays.map((weekday) => (
            <span key={weekday} className="pb-2 font-sans text-muted-foreground text-xs">
              {weekday}
            </span>
          ))}
          {Array.from({ length: leading }, (_, i) => (
            // biome-ignore lint/suspicious/noArrayIndexKey: leading blanks are positional
            <span key={`blank-${i}`} />
          ))}
          {days.map((day) => (
            <span
              key={day}
              className={cn(
                "mx-auto flex h-8 w-8 items-center justify-center rounded-md font-mono text-sm tabular-nums",
                day === today.day && "border border-foreground font-medium",
              )}
            >
              {day}
            </span>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
