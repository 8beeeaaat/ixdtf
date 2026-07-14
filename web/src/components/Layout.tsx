import { Link, Outlet } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { type Language, setLanguage } from "@/app/i18n";
import { type ThemeMode, useTheme } from "@/app/theme";
import { buttonVariants } from "@/components/ui/Button";
import { Select } from "@/components/ui/Select";
import { TooltipProvider } from "@/components/ui/Tooltip";
import { cn } from "@/lib/utils";

const NAV_ITEMS = [
  { to: "/", labelKey: "nav.home" },
  { to: "/about", labelKey: "nav.about" },
  { to: "/guide", labelKey: "nav.guide" },
  { to: "/playground", labelKey: "nav.playground" },
  { to: "/interop", labelKey: "nav.interop" },
  { to: "/converter", labelKey: "nav.converter" },
  { to: "/temporal-lab", labelKey: "nav.temporalLab" },
] as const;

/** 現在のテーマを表すモノクロ線画アイコン (DESIGN.md: 彩度色は導入しない)。 */
function ThemeModeIcon({ mode }: { mode: ThemeMode }) {
  const common = {
    viewBox: "0 0 16 16",
    "aria-hidden": true,
    fill: "none",
    stroke: "currentColor",
    strokeWidth: 1.3,
    strokeLinecap: "round",
    strokeLinejoin: "round",
    className: "h-4 w-4 shrink-0",
  } as const;
  if (mode === "light") {
    return (
      <svg {...common}>
        <circle cx="8" cy="8" r="3.2" />
        <path d="M8 1.2v1.6M8 13.2v1.6M1.2 8h1.6M13.2 8h1.6M3.2 3.2l1.1 1.1M11.7 11.7l1.1 1.1M12.8 3.2l-1.1 1.1M4.3 11.7l-1.1 1.1" />
      </svg>
    );
  }
  if (mode === "dark") {
    return (
      <svg {...common}>
        <path d="M13.2 9.8A5.8 5.8 0 1 1 6.2 2.8a4.6 4.6 0 0 0 7 7Z" />
      </svg>
    );
  }
  // auto = OS 設定に追従 (モニター)
  return (
    <svg {...common}>
      <rect x="1.7" y="2.7" width="12.6" height="8.6" rx="1.2" />
      <path d="M5.5 14h5M8 11.3V14" />
    </svg>
  );
}

/** 言語切替の地球儀アイコン。 */
function LanguageIcon() {
  return (
    <svg
      viewBox="0 0 16 16"
      aria-hidden="true"
      fill="none"
      stroke="currentColor"
      strokeWidth={1.3}
      strokeLinecap="round"
      className="h-4 w-4 shrink-0"
    >
      <circle cx="8" cy="8" r="6.3" />
      <path d="M1.7 8h12.6M8 1.7c1.9 1.8 1.9 10.8 0 12.6M8 1.7c-1.9 1.8-1.9 10.8 0 12.6" />
    </svg>
  );
}

/** App shell: fixed header nav (F-0-1), language (F-0-2) + theme (F-0-3) switches. */
export function Layout() {
  const { t, i18n } = useTranslation();
  const { mode, setMode } = useTheme();
  // ヘッダーの設定コントロールはすべて Select (プルダウン) で現在値を表示する。
  const currentLanguage: Language = i18n.language.startsWith("ja") ? "ja" : "en";
  const languageOptions = [
    { value: "ja", label: t("lang.ja") },
    { value: "en", label: t("lang.en") },
  ];
  const themeOptions = [
    { value: "auto", label: t("theme.auto") },
    { value: "light", label: t("theme.light") },
    { value: "dark", label: t("theme.dark") },
  ];

  return (
    <TooltipProvider>
      <div className="flex min-h-screen flex-col">
        <header className="sticky top-0 z-30 border-border border-b bg-background">
          <div className="mx-auto flex h-14 max-w-6xl items-center gap-6 px-6">
            <Link className="shrink-0 font-medium font-mono text-sm tabular-nums" to="/">
              {t("app.title")}
            </Link>
            <nav className="flex flex-1 gap-4 overflow-x-auto">
              {NAV_ITEMS.map((item) => (
                <Link
                  key={item.to}
                  to={item.to}
                  className="shrink-0 py-1 font-sans text-muted-foreground text-sm hover:text-foreground"
                  activeOptions={{ exact: item.to === "/" }}
                  activeProps={{ className: "text-foreground underline underline-offset-8" }}
                >
                  {t(item.labelKey)}
                </Link>
              ))}
            </nav>
            <div className="flex shrink-0 items-center gap-1">
              {/* ixdtf ライブラリへの導線。全幅バナーをやめ、GitHub マーク (foreground のみ着色) の
                  アイコンリンクとしてヘッダー右クラスタに常設する (DESIGN.md: 彩度色を導入しない) */}
              <a
                href="https://github.com/8beeeaaat/ixdtf"
                target="_blank"
                rel="noreferrer"
                className={cn(buttonVariants({ variant: "ghost", size: "sm" }), "shrink-0")}
              >
                <span className="sr-only">{t("banner.ixdtf")}</span>
                <svg
                  viewBox="0 0 16 16"
                  aria-hidden="true"
                  fill="currentColor"
                  className="h-4 w-4 shrink-0 text-foreground"
                >
                  <path d="M8 0c4.42 0 8 3.58 8 8a8.013 8.013 0 0 1-5.45 7.59c-.4.08-.55-.17-.55-.38 0-.27.01-1.13.01-2.2 0-.75-.25-1.23-.54-1.48 1.78-.2 3.65-.88 3.65-3.95 0-.88-.31-1.59-.82-2.15.08-.2.36-1.02-.08-2.12 0 0-.67-.22-2.2.82-.64-.18-1.32-.27-2-.27-.68 0-1.36.09-2 .27-1.53-1.03-2.2-.82-2.2-.82-.44 1.1-.16 1.92-.08 2.12-.51.56-.82 1.28-.82 2.15 0 3.06 1.86 3.75 3.64 3.95-.23.2-.44.55-.51 1.07-.46.21-1.61.55-2.33-.66-.15-.24-.6-.83-1.23-.82-.67.01-.27.38.01.53.34.19.73.9.82 1.13.16.45.68 1.31 2.69.94 0 .67.01 1.3.01 1.49 0 .21-.15.45-.55.38A7.995 7.995 0 0 1 0 8c0-4.42 3.58-8 8-8Z" />
                </svg>
              </a>
              <Link
                to="/support"
                aria-label={t("sponsor.ariaLabel")}
                className={cn(buttonVariants({ variant: "ghost", size: "sm" }), "shrink-0 gap-1.5")}
              >
                {/* モノクロのハート。GitHub マーク同様、foreground トークン以外に着色しない
                    (DESIGN.md: 彩度色は IXDTF ハイライトと状態表示に限定) */}
                <svg
                  viewBox="0 0 16 16"
                  aria-hidden="true"
                  fill="currentColor"
                  className="h-4 w-4 shrink-0 text-foreground"
                >
                  <path d="m8 14.25.345.666a.75.75 0 0 1-.69 0l-.008-.004-.018-.01a7.152 7.152 0 0 1-.31-.17 22.055 22.055 0 0 1-3.434-2.414C2.045 10.731 0 8.35 0 5.5 0 2.836 2.086 1 4.25 1 5.797 1 7.153 1.802 8 3.02 8.847 1.802 10.203 1 11.75 1 13.914 1 16 2.836 16 5.5c0 2.85-2.045 5.231-3.885 6.818a22.066 22.066 0 0 1-3.744 2.584l-.018.01-.006.003h-.002Z" />
                </svg>
                <span className="hidden sm:inline">{t("sponsor.cta")}</span>
              </Link>
              <Select
                aria-label={t("lang.label")}
                value={currentLanguage}
                onValueChange={(value) => setLanguage(value as Language)}
                options={languageOptions}
                icon={<LanguageIcon />}
              />
              <Select
                aria-label={t("theme.label")}
                value={mode}
                onValueChange={(value) => setMode(value as ThemeMode)}
                options={themeOptions}
                icon={<ThemeModeIcon mode={mode} />}
              />
            </div>
          </div>
        </header>
        <main className="w-full flex-1 py-12">
          <Outlet />
        </main>
      </div>
    </TooltipProvider>
  );
}
