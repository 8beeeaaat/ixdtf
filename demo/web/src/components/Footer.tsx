import { Link } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { SPONSOR_LINKS } from "@/lib/sponsor";

/** App footer: honest free / open-source note plus support links (sponsorship model). */
export function Footer() {
  const { t } = useTranslation();

  return (
    <footer className="border-border border-t bg-background">
      <div className="mx-auto flex max-w-6xl flex-col items-center gap-3 px-6 py-6 text-center sm:flex-row sm:justify-between sm:text-left">
        <p className="font-sans text-muted-foreground text-xs">{t("footer.free")}</p>
        <nav className="flex flex-wrap items-center justify-center gap-4">
          <a
            href={SPONSOR_LINKS.repo}
            target="_blank"
            rel="noreferrer"
            className="font-sans text-muted-foreground text-xs hover:text-foreground"
          >
            {t("footer.repo")} ↗
          </a>
          <a
            href={SPONSOR_LINKS.buyMeACoffee}
            target="_blank"
            rel="noreferrer"
            className="font-sans text-muted-foreground text-xs hover:text-foreground"
          >
            {t("footer.coffee")} ↗
          </a>
          <Link
            to="/support"
            className="font-sans text-foreground text-xs underline underline-offset-4 hover:opacity-80"
          >
            {t("footer.sponsor")}
          </Link>
        </nav>
      </div>
    </footer>
  );
}
