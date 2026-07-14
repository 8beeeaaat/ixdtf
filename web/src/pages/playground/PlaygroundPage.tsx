import { useNavigate, useSearch } from "@tanstack/react-router";
import { type ReactNode, useEffect, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { CopyButton } from "@/components/CopyButton";
import { IxdtfHighlight, IxdtfLegend } from "@/components/IxdtfHighlight";
import { ReferenceDialog } from "@/components/ReferenceDialog";
import { ReproduceSection } from "@/components/ReproduceSection";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card";
import { FormField } from "@/components/ui/FormField";
import { Input } from "@/components/ui/Input";
import { TagBadge } from "@/components/ui/TagBadge";
import { ToggleSwitch } from "@/components/ui/ToggleSwitch";
import { useParseIxdtf } from "@/generated/api/endpoints";
import { goParseSample, jsParseSample } from "@/lib/reproduce";
import { type BrowserParseResult, browserParse } from "@/lib/temporal/browserParse";
import { getNativeTemporal, getPolyfillTemporal } from "@/lib/temporal/detect";
import { useDebounce } from "@/lib/useDebounce";
import { useNowIXDTF } from "../home/useNow";

function Row({ label, children }: { label: string; children: ReactNode }) {
  return (
    <div className="flex flex-wrap items-baseline justify-between gap-2 border-border border-t py-2 first:border-t-0 first:pt-0">
      <dt className="font-sans text-muted-foreground text-xs">{label}</dt>
      <dd className="break-all font-mono text-sm tabular-nums">{children}</dd>
    </div>
  );
}

/** F-2: real-time IXDTF parse / validate playground with a shareable URL (F-2-7). */
export function PlaygroundPage() {
  const { t } = useTranslation();
  const search = useSearch({ from: "/playground" });
  const navigate = useNavigate({ from: "/playground" });
  const nowIXDTF = useNowIXDTF();
  const [input, setInput] = useState(search.input ?? nowIXDTF);
  const strict = search.strict ?? false;
  const validateOnly = search.validateOnly ?? false;
  const debounced = useDebounce(input, 300);

  // 入力はデバウンス後に URL へ replace 同期する (F-2-7 共有 URL)
  useEffect(() => {
    void navigate({
      search: {
        input: debounced || undefined,
        strict: strict || undefined,
        validateOnly: validateOnly || undefined,
      },
      replace: true,
    });
  }, [debounced, strict, validateOnly, navigate]);

  const setFlag = (key: "strict" | "validateOnly") => (enabled: boolean) => {
    void navigate({
      search: {
        input: debounced || undefined,
        strict: (key === "strict" ? enabled : strict) || undefined,
        validateOnly: (key === "validateOnly" ? enabled : validateOnly) || undefined,
      },
      replace: true,
    });
  };

  const parseQuery = useParseIxdtf(
    { input: debounced, strict, validate_only: validateOnly || undefined },
    { query: { enabled: debounced.length > 0 } },
  );
  const server = parseQuery.data?.status === 200 ? parseQuery.data.data : null;

  // F-3 と同じ思想で 3 実装を明示指定して併記する (ヘッダーの選択に依存しない)。
  // native はブラウザ非対応なら null → そのカードに案内メッセージを出す。
  const nativeTemporal = getNativeTemporal();
  const nativeSupported = nativeTemporal !== null;
  const nativeBrowser = useMemo(
    () => (debounced && nativeTemporal ? browserParse(debounced, nativeTemporal) : null),
    [debounced, nativeTemporal],
  );
  const polyfillBrowser = useMemo(
    () => (debounced ? browserParse(debounced, getPolyfillTemporal()) : null),
    [debounced],
  );

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

  return (
    <div className="mx-auto max-w-6xl space-y-12 px-6">
      <header>
        <div className="flex flex-wrap items-center gap-2">
          <h1 className="font-sans font-semibold text-2xl tracking-tight">
            {t("playground.title")}
          </h1>
          <ReferenceDialog referenceId="ixdtf" />
        </div>
        <p className="mt-1 font-sans text-muted-foreground text-sm">{t("playground.tagline")}</p>
      </header>

      <section className="space-y-4">
        <FormField label={t("playground.inputLabel")} htmlFor="playground-input">
          <Input
            id="playground-input"
            value={input}
            onChange={(event) => setInput(event.target.value)}
            placeholder={t("playground.inputPlaceholder")}
            autoComplete="off"
            spellCheck={false}
          />
        </FormField>
        {debounced && (
          <div className="space-y-2">
            <IxdtfHighlight value={debounced} className="text-lg md:text-2xl" />
            <IxdtfLegend />
          </div>
        )}
        <div className="flex flex-wrap gap-x-8 gap-y-3">
          <div className="flex items-center gap-2">
            <ToggleSwitch
              enabled={strict}
              onChange={setFlag("strict")}
              label={t("playground.strict")}
            />
            <span className="font-sans text-sm">{t("playground.strict")}</span>
            <span className="font-sans text-muted-foreground text-xs">
              {t("playground.strictHint")}
            </span>
            <ReferenceDialog referenceId="offsetConsistency" />
          </div>
          <div className="flex items-center gap-2">
            <ToggleSwitch
              enabled={validateOnly}
              onChange={setFlag("validateOnly")}
              label={t("playground.validateOnly")}
            />
            <span className="font-sans text-sm">{t("playground.validateOnly")}</span>
            <span className="font-sans text-muted-foreground text-xs">
              {t("playground.validateOnlyHint")}
            </span>
            <ReferenceDialog referenceId="goParsing" />
          </div>
          <CopyButton
            text={() => window.location.href}
            label={t("playground.share")}
            copiedLabel={t("playground.shareCopied")}
          />
        </div>
      </section>

      {!debounced ? (
        <p className="font-sans text-muted-foreground text-sm">{t("playground.empty")}</p>
      ) : (
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
                        <Row label={t("playground.offsetSeconds")}>
                          {server.result.offset_seconds}
                        </Row>
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
                                <TagBadge
                                  key={tag.key}
                                  variant={tag.critical ? "critical" : "default"}
                                >
                                  {tag.key}={tag.value}
                                </TagBadge>
                              ))}
                            </span>
                          )}
                        </Row>
                      </dl>
                    ) : (
                      <p className="font-sans text-sm text-success">
                        {t("playground.validatedOk")}
                      </p>
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
            goSample={goParseSample({ input: debounced, strict, validateOnly })}
            jsSample={jsParseSample(debounced)}
          />
        </>
      )}
    </div>
  );
}
