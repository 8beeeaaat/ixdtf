import { useTranslation } from "react-i18next";

/** Fallback screen when the browser lacks native Temporal (F-0-4, no polyfill). */
export function UnsupportedBrowser() {
  const { t } = useTranslation();
  return (
    <main className="mx-auto flex min-h-screen max-w-6xl flex-col justify-center gap-4 px-6">
      <h1 className="font-sans font-semibold text-2xl tracking-tight">{t("unsupported.title")}</h1>
      <p className="font-sans text-foreground text-sm">{t("unsupported.body")}</p>
      <p className="font-mono text-sm tabular-nums">{t("unsupported.browsers")}</p>
      <p className="font-sans text-muted-foreground text-xs">{t("unsupported.note")}</p>
    </main>
  );
}
