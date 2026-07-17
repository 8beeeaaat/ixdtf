import { useState } from "react";
import { useTranslation } from "react-i18next";
import { CalendarPicker } from "@/components/CalendarPicker";
import { ReproduceSection } from "@/components/ReproduceSection";
import { TimeZonePicker } from "@/components/TimeZonePicker";
import { FormField } from "@/components/ui/FormField";
import { goNowSample, jsNowSample } from "@/lib/reproduce";
import { ClockCard } from "@/pages/home/ClockCard";
import { GlobeBackground } from "@/pages/home/GlobeBackground";
import { MonthCalendar } from "@/pages/home/MonthCalendar";
import { ServerNowCard } from "@/pages/home/ServerNowCard";

/** F-1: IXDTF clock / calendar — the face of the app. */
export function HomePage() {
  const { t } = useTranslation();
  const [timeZone, setTimeZone] = useState(
    () => new Intl.DateTimeFormat().resolvedOptions().timeZone,
  );
  const [calendar, setCalendar] = useState("iso8601");
  const sampleArgs = { timeZone, calendar: calendar !== "iso8601" ? calendar : undefined };

  return (
    <>
      <GlobeBackground timeZone={timeZone} onTimeZoneChange={setTimeZone} />
      {/* Outer is click-through so the open right side reaches the globe; the inner
          column re-enables pointer events so the cards stay interactive (D-10). */}
      <div className="pointer-events-none relative z-10 mx-auto max-w-6xl px-6">
        {/* Home only: narrow, left-aligned column so the globe owns the right (D-10). */}
        <div className="pointer-events-auto max-w-full space-y-12 md:max-w-2xl">
          <header>
            <h1 className="font-sans font-semibold text-2xl tracking-tight">{t("home.title")}</h1>
            <p className="mt-1 font-sans text-muted-foreground text-sm">{t("home.tagline")}</p>
          </header>
          <ClockCard timeZone={timeZone} calendar={calendar} />
          <div className="grid gap-4 md:grid-cols-2">
            <FormField label={t("home.timeZoneLabel")} htmlFor="home-time-zone">
              <TimeZonePicker id="home-time-zone" value={timeZone} onChange={setTimeZone} />
            </FormField>
            <FormField label={t("home.calendarLabel")} htmlFor="home-calendar">
              <CalendarPicker id="home-calendar" value={calendar} onChange={setCalendar} />
            </FormField>
          </div>
          <ServerNowCard timeZone={timeZone} calendar={calendar} />
          <ReproduceSection goSample={goNowSample(sampleArgs)} jsSample={jsNowSample(sampleArgs)} />
          <MonthCalendar timeZone={timeZone} calendar={calendar} />
        </div>
      </div>
    </>
  );
}
