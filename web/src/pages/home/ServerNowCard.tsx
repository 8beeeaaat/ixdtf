import { useTranslation } from "react-i18next";
import { IxdtfHighlight } from "@/components/IxdtfHighlight";
import { Button } from "@/components/ui/Button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card";
import { useGetNow } from "@/generated/api/endpoints";

interface ServerNowCardProps {
  timeZone: string;
  calendar: string;
}

/** Server-generated IXDTF now (`GET /api/now`, F-1-7) shown beside the browser clock. */
export function ServerNowCard({ timeZone, calendar }: ServerNowCardProps) {
  const { t } = useTranslation();
  const query = useGetNow({
    time_zone: timeZone,
    ...(calendar !== "iso8601" ? { calendar } : {}),
  });
  const now = query.data?.status === 200 ? query.data.data : null;

  return (
    <Card>
      <CardHeader>
        <CardTitle>{t("home.serverGenerated")}</CardTitle>
      </CardHeader>
      <CardContent className="space-y-3">
        {now ? (
          <>
            <IxdtfHighlight value={now.ixdtf} className="text-base md:text-xl" />
            <p className="font-mono text-muted-foreground text-xs tabular-nums">
              unix_nano: {now.unix_nano}
            </p>
          </>
        ) : (
          <p className="font-sans text-muted-foreground text-sm">
            {query.isFetching ? t("common.loading") : t("common.error")}
          </p>
        )}
        <div className="flex items-center justify-between gap-4">
          <p className="font-sans text-muted-foreground text-xs">{t("home.serverNote")}</p>
          <Button variant="secondary" size="sm" onClick={() => void query.refetch()}>
            {t("common.refresh")}
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}
