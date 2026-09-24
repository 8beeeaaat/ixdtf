import { useTranslation } from "react-i18next";
import { CopyButton } from "@/components/CopyButton";
import { IxdtfHighlight, IxdtfLegend } from "@/components/IxdtfHighlight";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card";
import { useNowTimer } from "@/pages/home/useNow";

interface ClockCardProps {
  timeZone: string;
  calendar: string;
}

/**
 * The 1-second clock (F-1-1) and the full IXDTF string with per-component
 * highlighting (F-1-2). The tick lives entirely inside this card (N-4).
 */
export function ClockCard({ timeZone, calendar }: ClockCardProps) {
  const { t } = useTranslation();
  const now = useNowTimer(timeZone);
  const display = calendar === "iso8601" ? now : now.withCalendar(calendar);
  const ixdtf = display.toString({ fractionalSecondDigits: 9 });

  return (
    <Card>
      <CardHeader>
        <CardTitle>{t("home.browserGenerated")}</CardTitle>
      </CardHeader>
      <CardContent className="space-y-6">
        <div>
          <p className="font-mono font-normal text-6xl tabular-nums tracking-tight md:text-8xl">
            {now.toPlainTime().toString({ smallestUnit: "second" })}
          </p>
          <p className="mt-2 font-mono text-muted-foreground text-sm tabular-nums">
            {now.toPlainDate().toString()}
          </p>
        </div>
        <div className="space-y-3 border-border border-t pt-4">
          <IxdtfHighlight value={ixdtf} className="text-base md:text-xl" />
          <div className="flex items-center justify-between gap-4">
            <IxdtfLegend />
            <CopyButton text={() => ixdtf} />
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
