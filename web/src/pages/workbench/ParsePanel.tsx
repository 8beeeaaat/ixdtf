import { type ReactNode, useMemo } from "react";
import { useTranslation } from "react-i18next";
import { ReferenceDialog } from "@/components/ReferenceDialog";
import { ReproduceSection } from "@/components/ReproduceSection";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card";
import { TagBadge } from "@/components/ui/TagBadge";
import { useParseIxdtf } from "@/generated/api/endpoints";
import { goParseSample, jsParseSample } from "@/lib/reproduce";
import { type BrowserParseResult, browserParse } from "@/lib/temporal/browserParse";
import { getNativeTemporal, getPolyfillTemporal } from "@/lib/temporal/detect";

function Row({ label, children }: { label: string; children: ReactNode }) {
  return (
    <div className="flex flex-wrap items-baseline justify-between gap-2 border-border border-t py-2 first:border-t-0 first:pt-0">
      <dt className="font-sans text-muted-foreground text-xs">{label}</dt>
      <dd className="break-all font-mono text-sm tabular-nums">{children}</dd>
    </div>
  );
}

interface ParsePanelProps {
  /** デバウンス後の入力 (空でないことを呼び出し側が保証する) */
  input: string;
  strict: boolean;
  validateOnly: boolean;
}

/** ワークベンチ「解析・検証」モード (旧 F-2 Playground)。3 実装の解析結果を併記する。 */
export function ParsePanel({ input, strict, validateOnly }: ParsePanelProps) {
  const { t } = useTranslation();

  const parseQuery = useParseIxdtf(
    { input, strict, validate_only: validateOnly || undefined },
    { query: { enabled: input.length > 0 } },
  );
  const server = parseQuery.data?.status === 200 ? parseQuery.data.data : null;

  // 3 実装を明示指定して併記する (ヘッダーの選択に依存しない)。
  // native はブラウザ非対応なら null → そのカードに案内メッセージを出す。
  const nativeTemporal = getNativeTemporal();
  const nativeSupported = nativeTemporal !== null;
  const nativeBrowser = useMemo(
    () => (nativeTemporal ? browserParse(input, nativeTemporal) : null),
    [input, nativeTemporal],
  );
  const polyfillBrowser = useMemo(() => browserParse(input, getPolyfillTemporal()), [input]);

  // ブラウザ側 (native / polyfill) の結果カード本体。両カードで共通に使う。
  const renderBrowserResult = (parse: BrowserParseResult | null) => {
    if (!parse) {
      return <p className="font-sans text-muted-foreground text-sm">{t("common.loading")}</p>;
    }
    return parse.ok ? (
      <dl>
        <Row label={t("playground.via")}>
          {parse.via === "zonedDateTime" ? "ZonedDateTime.from" : "Instant.from"}
        </Row>
        <Row label={t("playground.parsedTimeBrowser")}>{parse.formatted}</Row>
        <Row label={t("playground.unixNano")}>{parse.epochNanoseconds}</Row>
        <Row label={t("playground.timeZone")}>{parse.timeZone ?? t("common.none")}</Row>
      </dl>
    ) : (
      <div className="space-y-2">
        <p className="font-sans text-muted-foreground text-xs">
          {t("playground.errorLabelBrowser")}
        </p>
        <p className="break-all rounded-md bg-destructive/10 p-3 font-mono text-destructive text-sm">
          {parse.error}
        </p>
      </div>
    );
  };

  if (!input) {
    return <p className="font-sans text-muted-foreground text-sm">{t("playground.empty")}</p>;
  }

  return (
    <>
      <div className="grid items-start gap-6 lg:grid-cols-3">
        <Card>
          <CardHeader>
            <CardTitle className="flex flex-wrap items-center gap-1">
              {t("playground.browserResult", { impl: t("temporal.implNative") })}
              <ReferenceDialog referenceId="temporalParsing" />
            </CardTitle>
          </CardHeader>
          <CardContent>
            {nativeSupported ? (
              renderBrowserResult(nativeBrowser)
            ) : (
              <p className="font-sans text-muted-foreground text-sm">
                {t("playground.nativeUnsupported")}
              </p>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex flex-wrap items-center gap-1">
              {t("playground.browserResult", { impl: t("temporal.implPolyfill") })}
              <ReferenceDialog referenceId="temporalParsing" />
            </CardTitle>
          </CardHeader>
          <CardContent>{renderBrowserResult(polyfillBrowser)}</CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex flex-wrap items-center gap-1">
              {t("playground.serverResult")}
              <ReferenceDialog referenceId="goParsing" />
            </CardTitle>
          </CardHeader>
          <CardContent>
            {server ? (
              server.ok ? (
                server.result ? (
                  <dl>
                    <Row label={t("playground.parsedTime")}>{server.result.rfc3339}</Row>
                    <Row label={t("playground.unixNano")}>{server.result.unix_nano}</Row>
                    <Row label={t("playground.offsetSeconds")}>{server.result.offset_seconds}</Row>
                    <Row label={t("playground.timeZone")}>
                      {server.result.time_zone ?? t("common.none")}
                      {server.result.time_zone_critical && (
                        <TagBadge variant="critical" className="ml-2">
                          !
                        </TagBadge>
                      )}
                    </Row>
                    <Row label={t("playground.tags")}>
                      {server.result.tags.length === 0 ? (
                        t("playground.noTags")
                      ) : (
                        <span className="flex flex-wrap justify-end gap-1">
                          {server.result.tags.map((tag) => (
                            <TagBadge key={tag.key} variant={tag.critical ? "critical" : "default"}>
                              {tag.key}={tag.value}
                            </TagBadge>
                          ))}
                        </span>
                      )}
                    </Row>
                  </dl>
                ) : (
                  <p className="font-sans text-sm text-success">{t("playground.validatedOk")}</p>
                )
              ) : (
                <div className="space-y-2">
                  <p className="font-sans text-muted-foreground text-xs">
                    {t("playground.errorLabel")}
                  </p>
                  <p className="break-all rounded-md bg-destructive/10 p-3 font-mono text-destructive text-sm">
                    {server.error?.message}
                  </p>
                </div>
              )
            ) : (
              <p className="font-sans text-muted-foreground text-sm">{t("common.loading")}</p>
            )}
          </CardContent>
        </Card>
      </div>
      <ReproduceSection
        goSample={goParseSample({ input, strict, validateOnly })}
        jsSample={jsParseSample(input)}
      />
    </>
  );
}
