import { useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { IxdtfHighlight, IxdtfLegend } from "@/components/IxdtfHighlight";
import { ReferenceDialog } from "@/components/ReferenceDialog";
import { ReproduceSection } from "@/components/ReproduceSection";
import { Button } from "@/components/ui/Button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card";
import { FormField } from "@/components/ui/FormField";
import { Input } from "@/components/ui/Input";
import { TagBadge } from "@/components/ui/TagBadge";
import { ToggleSwitch } from "@/components/ui/ToggleSwitch";
import { useRoundtripIxdtf } from "@/generated/api/endpoints";
import { interopPresets } from "@/lib/fixtures";
import { goRoundtripSample, jsRoundtripSample } from "@/lib/reproduce";
import { browserParse } from "@/lib/temporal/browserParse";
import { useDebounce } from "@/lib/useDebounce";
import { cn } from "@/lib/utils";
import { useNowIXDTF } from "../home/useNow";

interface ComparisonRow {
  key: string;
  browser: string | null;
  server: string | null;
  /** undefined = 比較対象外 (ハイライトしない) */
  match?: boolean;
}

/** F-3: the behavioural differences between the two implementations ARE the content. */
export function InteropPage() {
  const { t } = useTranslation();
  const nowIXDTF = useNowIXDTF();
  const [input, setInput] = useState(nowIXDTF);
  const [strict, setStrict] = useState(false);
  const debounced = useDebounce(input, 300);

  const roundtripQuery = useRoundtripIxdtf(
    { input: debounced, strict },
    { query: { enabled: debounced.length > 0 } },
  );
  const server = roundtripQuery.data?.status === 200 ? roundtripQuery.data.data : null;
  const browser = useMemo(() => (debounced ? browserParse(debounced) : null), [debounced]);

  // 選択中プリセット (入力と strict の一致から導出)。fixtures の note_key で
  // 「このプリセットが何を実証するか」を説明する (F-3-3 の学習導線)
  const activePreset =
    interopPresets.find((preset) => preset.input === debounced && preset.strict === strict) ?? null;

  const rows: ComparisonRow[] = useMemo(() => {
    if (!browser || !server) {
      return [];
    }
    const parse = server.parse;
    const serverResult = parse.result ?? null;
    const browserOk = browser.ok;
    const serverCalendarTag = serverResult?.tags.find((tag) => tag.key === "u-ca") ?? null;
    const serverCalendar = serverCalendarTag?.value ?? (serverResult ? "iso8601" : null);
    const browserCalendar = browser.ok ? browser.calendar : null;
    const browserLossless = browser.ok ? String(browser.formatted === debounced) : null;
    const okLabel = (ok: boolean) => t(ok ? "common.ok" : "common.error");
    return [
      {
        key: "status",
        browser: okLabel(browserOk),
        server: okLabel(parse.ok),
        match: browserOk === parse.ok,
      },
      {
        key: "error",
        browser: browser.ok ? null : browser.error,
        server: parse.error?.message ?? null,
      },
      {
        key: "epoch",
        browser: browser.ok ? browser.epochNanoseconds : null,
        server: serverResult?.unix_nano ?? null,
        match:
          browser.ok && serverResult
            ? browser.epochNanoseconds === serverResult.unix_nano
            : undefined,
      },
      {
        key: "timeZone",
        browser: browser.ok ? browser.timeZone : null,
        server: serverResult?.time_zone ?? null,
        match:
          browser.ok && serverResult
            ? (browser.timeZone ?? "") === (serverResult.time_zone ?? "")
            : undefined,
      },
      {
        key: "calendar",
        browser: browserCalendar,
        server: serverCalendar,
        match:
          browserCalendar !== null && serverCalendar !== null
            ? browserCalendar === serverCalendar
            : undefined,
      },
      {
        key: "roundtrip",
        browser: browser.ok ? browser.formatted : null,
        server: server.formatted,
        match:
          browser.ok && server.formatted !== null
            ? browser.formatted === server.formatted
            : undefined,
      },
      {
        key: "lossless",
        browser: browserLossless,
        server: server.lossless === null ? null : String(server.lossless),
        match:
          browserLossless !== null && server.lossless !== null
            ? browserLossless === String(server.lossless)
            : undefined,
      },
    ];
  }, [browser, server, debounced, t]);

  return (
    <div className="mx-auto max-w-6xl space-y-12 px-6">
      <header>
        <div className="flex flex-wrap items-center gap-2">
          <h1 className="font-sans font-semibold text-2xl tracking-tight">{t("interop.title")}</h1>
          <ReferenceDialog referenceId="roundtrip" />
        </div>
        <p className="mt-1 font-sans text-muted-foreground text-sm">{t("interop.tagline")}</p>
      </header>

      <section className="space-y-4">
        <FormField label={t("interop.inputLabel")} htmlFor="interop-input">
          <Input
            id="interop-input"
            value={input}
            onChange={(event) => setInput(event.target.value)}
            placeholder={t("interop.inputPlaceholder")}
            autoComplete="off"
            spellCheck={false}
          />
        </FormField>
        <div className="flex flex-wrap items-center gap-x-8 gap-y-3">
          <div className="flex items-center gap-2">
            <ToggleSwitch enabled={strict} onChange={setStrict} label={t("common.strict")} />
            <span className="font-sans text-sm">{t("common.strict")}</span>
            <ReferenceDialog referenceId="offsetConsistency" />
          </div>
        </div>
        {debounced && (
          <div className="space-y-2">
            <IxdtfHighlight value={debounced} className="text-lg md:text-2xl" />
            <IxdtfLegend />
          </div>
        )}
      </section>

      <section className="space-y-3">
        <h2 className="font-sans font-semibold text-muted-foreground text-sm uppercase tracking-wider">
          {t("interop.presets")}
        </h2>
        <div className="flex flex-wrap gap-2">
          {interopPresets.map((preset) => (
            <Button
              key={preset.id}
              variant="secondary"
              size="sm"
              aria-pressed={activePreset?.id === preset.id}
              className={cn("font-mono", activePreset?.id === preset.id && "border-foreground")}
              onClick={() => {
                setInput(preset.input);
                setStrict(preset.strict);
              }}
            >
              {preset.id}
              {preset.strict && (
                <TagBadge variant="default" className="ml-1">
                  {t("common.strict")}
                </TagBadge>
              )}
            </Button>
          ))}
        </div>
        {activePreset && (
          <p className="max-w-3xl font-sans text-muted-foreground text-sm">
            {t(activePreset.note_key)}
          </p>
        )}
      </section>

      {!debounced || rows.length === 0 ? (
        <p className="font-sans text-muted-foreground text-sm">{t("interop.empty")}</p>
      ) : (
        <>
          <Card>
            <CardHeader>
              <CardTitle className="flex flex-wrap items-center gap-1">
                {t("interop.resultTitle")}
                <ReferenceDialog referenceId="roundtrip" />
              </CardTitle>
            </CardHeader>
            <CardContent className="overflow-x-auto">
              <table className="w-full border-collapse text-left">
                <thead>
                  <tr className="border-border border-b">
                    <th className="py-2 pr-4 font-medium font-sans text-muted-foreground text-xs">
                      {t("interop.aspect")}
                    </th>
                    <th className="py-2 pr-4 font-medium font-sans text-muted-foreground text-xs">
                      {t("interop.browserCol")}
                    </th>
                    <th className="py-2 font-medium font-sans text-muted-foreground text-xs">
                      {t("interop.serverCol")}
                    </th>
                  </tr>
                </thead>
                <tbody>
                  {rows.map((row) => {
                    const cell = (value: string | null) => (
                      <span
                        className={cn(
                          "break-all font-mono text-sm tabular-nums",
                          row.match === true && "text-success",
                          row.match === false && "text-destructive",
                        )}
                      >
                        {value ?? "—"}
                      </span>
                    );
                    return (
                      <tr
                        key={row.key}
                        className="border-border border-b align-top last:border-b-0"
                      >
                        <td className="whitespace-nowrap py-2 pr-4 font-sans text-muted-foreground text-xs">
                          {t(`interop.rows.${row.key}`)}
                          {row.match === true && (
                            <TagBadge variant="match" className="ml-2">
                              {t("common.match")}
                            </TagBadge>
                          )}
                          {row.match === false && (
                            <TagBadge variant="mismatch" className="ml-2">
                              {t("common.mismatch")}
                            </TagBadge>
                          )}
                        </td>
                        <td className="max-w-96 py-2 pr-4">{cell(row.browser)}</td>
                        <td className="max-w-96 py-2">{cell(row.server)}</td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </CardContent>
          </Card>
          <ReproduceSection
            goSample={goRoundtripSample({ input: debounced, strict })}
            jsSample={jsRoundtripSample(debounced)}
          />
        </>
      )}
    </div>
  );
}
