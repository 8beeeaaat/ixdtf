import { useNavigate, useSearch } from "@tanstack/react-router";
import { useCallback, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import type { PlaygroundSearch, WorkbenchMode } from "@/app/router";
import { CopyButton } from "@/components/CopyButton";
import { IxdtfHighlight, IxdtfLegend } from "@/components/IxdtfHighlight";
import { ReferenceDialog } from "@/components/ReferenceDialog";
import { FormField } from "@/components/ui/FormField";
import { Input } from "@/components/ui/Input";
import { ToggleSwitch } from "@/components/ui/ToggleSwitch";
import { useDebounce } from "@/lib/useDebounce";
import { cn } from "@/lib/utils";
import { useNowIXDTF } from "../home/useNow";
import { LabPanel } from "../temporal-lab/LabPanel";
import { ParsePanel } from "./ParsePanel";
import { RoundtripPanel } from "./RoundtripPanel";

const MODES: readonly WorkbenchMode[] = ["parse", "roundtrip", "lab"];

// モードごとのラベル / hint キー。Record にして網羅性を型で担保する (モード追加漏れを検出)。
const MODE_LABEL_KEY: Record<WorkbenchMode, string> = {
  parse: "workbench.modeParse",
  roundtrip: "workbench.modeRoundtrip",
  lab: "workbench.modeLab",
};
const MODE_HINT_KEY: Record<WorkbenchMode, string> = {
  parse: "workbench.modeParseHint",
  roundtrip: "workbench.modeRoundtripHint",
  lab: "workbench.modeLabHint",
};

/**
 * F-2 IXDTF ワークベンチ: 旧 Playground (解析・検証)・旧 Interop (往復・実装差比較)・
 * 旧 Temporal Lab (F-6 日時モデル観察) を 1 画面 3 モードに統合する。解析・往復は
 * 共有入力・共有ハイライト・共有 3 実装を使い、Lab は固定 fixture 駆動 (LabPanel)。
 * URL クエリ (input / strict / validateOnly / mode) で状態を共有する (F-2-7)。
 */
export function WorkbenchPage() {
  const { t } = useTranslation();
  const search = useSearch({ from: "/playground" });
  const navigate = useNavigate({ from: "/playground" });
  const nowIXDTF = useNowIXDTF();
  const [input, setInput] = useState(search.input ?? nowIXDTF);
  const strict = search.strict ?? false;
  const validateOnly = search.validateOnly ?? false;
  const mode: WorkbenchMode =
    search.mode === "roundtrip" || search.mode === "lab" ? search.mode : "parse";
  const debounced = useDebounce(input, 300);

  // URL に載せる search を組み立てる。解析モードは既定なので mode を書かず、
  // validateOnly は解析モードでのみ意味を持つ (旧 /playground?input=&strict= を後方互換で受理)。
  const buildSearch = useCallback(
    (overrides: Partial<PlaygroundSearch>): PlaygroundSearch => {
      const nextMode = overrides.mode ?? mode;
      const nextInput = overrides.input ?? debounced;
      const nextStrict = overrides.strict ?? strict;
      const nextValidateOnly = overrides.validateOnly ?? validateOnly;
      return {
        input: nextInput || undefined,
        strict: nextStrict || undefined,
        validateOnly: (nextMode === "parse" && nextValidateOnly) || undefined,
        // parse は既定なので mode を書かない (旧 deep link 後方互換)。roundtrip / lab は明示。
        // input / strict は lab でも URL に残す — strict は URL 由来の state なので、
        // lab 進入時に落とすと parse へ戻った際に失われるため (往復整合性を優先)。
        mode: nextMode === "parse" ? undefined : nextMode,
      };
    },
    [mode, debounced, strict, validateOnly],
  );

  // 入力はデバウンス後に URL へ replace 同期する (F-2-7 共有 URL)
  useEffect(() => {
    void navigate({ search: buildSearch({}), replace: true });
  }, [buildSearch, navigate]);

  const setStrict = (enabled: boolean) =>
    void navigate({ search: buildSearch({ strict: enabled }), replace: true });
  const setValidateOnly = (enabled: boolean) =>
    void navigate({ search: buildSearch({ validateOnly: enabled }), replace: true });
  const setMode = (next: WorkbenchMode) =>
    void navigate({ search: buildSearch({ mode: next }), replace: true });
  const applyPreset = (presetInput: string, presetStrict: boolean) => {
    setInput(presetInput);
    void navigate({
      search: buildSearch({ input: presetInput, strict: presetStrict, mode: "roundtrip" }),
      replace: true,
    });
  };

  return (
    <div className="mx-auto max-w-6xl space-y-12 px-6">
      <header>
        <div className="flex flex-wrap items-center gap-2">
          <h1 className="font-sans font-semibold text-2xl tracking-tight">
            {t("workbench.title")}
          </h1>
          <ReferenceDialog referenceId="ixdtf" />
        </div>
        <p className="mt-1 font-sans text-muted-foreground text-sm">{t("workbench.tagline")}</p>
      </header>

      <section className="space-y-4">
        {/* fieldset + sr-only legend でモード群に名前を付ける (biome useSemanticElements 準拠)。
            fieldset のデフォルト min-inline-size: min-content は min-w-0 で解除する。 */}
        <fieldset
          className="inline-flex min-w-0 flex-wrap gap-0.5 rounded-md border border-border p-0.5"
          aria-label={t("workbench.modeLabel")}
        >
          {MODES.map((m) => (
            <button
              key={m}
              type="button"
              aria-pressed={mode === m}
              onClick={() => setMode(m)}
              className={cn(
                // モバイルはタップターゲット 44px 相当 (py-3)。デスクトップは従来の密度
                "shrink-0 whitespace-nowrap rounded px-3 py-3 font-sans text-sm transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background sm:py-1.5",
                mode === m
                  ? "bg-accent text-accent-foreground"
                  : "text-muted-foreground hover:text-foreground",
              )}
            >
              {t(MODE_LABEL_KEY[m])}
            </button>
          ))}
        </fieldset>
        <p className="font-sans text-muted-foreground text-sm">{t(MODE_HINT_KEY[mode])}</p>

        {/* Lab モードは固定 fixture 駆動なので共有入力欄・strict・ハイライトを持たない (F-6 / F-2-0)。
            モード切替タブと hint のみ共有し、入力系コントロールは解析・往復モードでのみ表示する。 */}
        {mode !== "lab" && (
          <>
            <FormField label={t("playground.inputLabel")} htmlFor="workbench-input">
              <Input
                id="workbench-input"
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
            <div className="flex flex-wrap items-center gap-x-8 gap-y-3">
              <div className="flex items-center gap-2">
                <ToggleSwitch
                  enabled={strict}
                  onChange={setStrict}
                  label={t("playground.strict")}
                />
                <span className="font-sans text-sm">{t("playground.strict")}</span>
                <span className="font-sans text-muted-foreground text-xs">
                  {t("playground.strictHint")}
                </span>
                <ReferenceDialog referenceId="offsetConsistency" />
              </div>
              {mode === "parse" && (
                <div className="flex items-center gap-2">
                  <ToggleSwitch
                    enabled={validateOnly}
                    onChange={setValidateOnly}
                    label={t("playground.validateOnly")}
                  />
                  <span className="font-sans text-sm">{t("playground.validateOnly")}</span>
                  <span className="font-sans text-muted-foreground text-xs">
                    {t("playground.validateOnlyHint")}
                  </span>
                  <ReferenceDialog referenceId="goParsing" />
                </div>
              )}
              <CopyButton
                text={() => window.location.href}
                label={t("playground.share")}
                copiedLabel={t("playground.shareCopied")}
              />
            </div>
          </>
        )}
      </section>

      {mode === "parse" ? (
        <ParsePanel input={debounced} strict={strict} validateOnly={validateOnly} />
      ) : mode === "roundtrip" ? (
        <RoundtripPanel input={debounced} strict={strict} onApplyPreset={applyPreset} />
      ) : (
        <LabPanel />
      )}
    </div>
  );
}
