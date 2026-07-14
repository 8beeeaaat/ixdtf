import { useTranslation } from "react-i18next";
import { buttonVariants } from "@/components/ui/Button";
import { Card, CardContent } from "@/components/ui/Card";
import { SPONSOR_LINKS } from "@/lib/sponsor";
import { cn } from "@/lib/utils";

interface SupportOption {
  key: "sponsor" | "coffee" | "star";
  href: string;
  variant: "primary" | "secondary";
}

const OPTIONS: SupportOption[] = [
  { key: "sponsor", href: SPONSOR_LINKS.githubSponsors, variant: "primary" },
  // { key: "coffee", href: SPONSOR_LINKS.buyMeACoffee, variant: "secondary" },
  // { key: "star", href: SPONSOR_LINKS.repo, variant: "secondary" },
];

/** Support page: explains the demo's value first, then offers the sponsorship links. */
export function SupportPage() {
  const { t } = useTranslation();

  return (
    <div className="mx-auto max-w-6xl space-y-12 px-6">
      <header>
        <h1 className="font-sans font-semibold text-2xl tracking-tight">{t("support.title")}</h1>
        <p className="mt-1 font-sans text-muted-foreground text-sm">{t("support.tagline")}</p>
      </header>

      <section className="space-y-2">
        <h2 className="font-sans font-semibold text-muted-foreground text-sm uppercase tracking-wider">
          {t("support.why.heading")}
        </h2>
        <p className="max-w-3xl font-sans text-sm leading-relaxed">{t("support.why.body")}</p>
      </section>

      <section className="space-y-4">
        <h2 className="font-sans font-semibold text-muted-foreground text-sm uppercase tracking-wider">
          {t("support.how.heading")}
        </h2>
        <div className="mx-auto grid w-full max-w-3xl gap-4 sm:grid-cols-[repeat(auto-fit,minmax(15rem,1fr))]">
          {OPTIONS.map((opt) => (
            <Card key={opt.key}>
              <CardContent className="flex h-full flex-col gap-3">
                <h3 className="font-sans font-semibold text-base">
                  {t(`support.options.${opt.key}Title`)}
                </h3>
                <p className="flex-1 font-sans text-muted-foreground text-sm leading-relaxed">
                  {t(`support.options.${opt.key}Body`)}
                </p>
                <a
                  href={opt.href}
                  target="_blank"
                  rel="noreferrer"
                  className={cn(buttonVariants({ variant: opt.variant, size: "sm" }), "w-full")}
                >
                  {t(`support.options.${opt.key}Cta`)} ↗
                </a>
              </CardContent>
            </Card>
          ))}
        </div>
      </section>

      <p className="font-sans text-muted-foreground text-xs">{t("support.note")}</p>
    </div>
  );
}
