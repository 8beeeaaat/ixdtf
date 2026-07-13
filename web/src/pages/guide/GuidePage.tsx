import { Link } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { IxdtfHighlight } from "@/components/IxdtfHighlight";
import { ReferenceDialog } from "@/components/ReferenceDialog";
import { buttonVariants } from "@/components/ui/Button";
import { Card, CardContent } from "@/components/ui/Card";
import { TagBadge } from "@/components/ui/TagBadge";
import { guideSamplesByCategory } from "@/lib/fixtures";
import type { ReferenceId } from "@/lib/references";
import { cn } from "@/lib/utils";

const GUIDE_SECTIONS = [
  { key: "rfc3339", referenceId: "ixdtf" },
  { key: "suffix", referenceId: "suffixSyntax" },
  { key: "critical", referenceId: "criticalAnnotations" },
  { key: "rejected", referenceId: "rejectedExtensions" },
] as const satisfies readonly { key: string; referenceId: ReferenceId }[];

/** F-5: RFC 9557 learning guide with the curated sample gallery (fixtures). */
export function GuidePage() {
  const { t } = useTranslation();
  const groups = guideSamplesByCategory();

  return (
    <div className="mx-auto max-w-6xl space-y-12 px-6">
      <header>
        <h1 className="font-sans font-semibold text-2xl tracking-tight">{t("guide.title")}</h1>
        <p className="mt-1 font-sans text-muted-foreground text-sm">{t("guide.tagline")}</p>
      </header>

      <section className="space-y-8">
        {GUIDE_SECTIONS.map(({ key, referenceId }) => (
          <div key={key} className="space-y-2">
            <div className="flex flex-wrap items-center gap-1">
              <h2 className="font-sans font-semibold text-muted-foreground text-sm uppercase tracking-wider">
                {t(`guide.sections.${key}.heading`)}
              </h2>
              <ReferenceDialog referenceId={referenceId} />
            </div>
            <p className="max-w-3xl font-sans text-sm leading-relaxed">
              {t(`guide.sections.${key}.body`)}
            </p>
          </div>
        ))}
      </section>

      <section className="space-y-8">
        <h2 className="font-sans font-semibold text-muted-foreground text-sm uppercase tracking-wider">
          {t("guide.gallery")}
        </h2>
        {groups.map(({ category, items }) => (
          <div key={category} className="space-y-4">
            <h3 className="font-sans font-semibold text-muted-foreground text-xs uppercase tracking-wider">
              {t(`guide.categories.${category}`)}
            </h3>
            {items.map((sample) => (
              <Card key={sample.id}>
                <CardContent className="space-y-3">
                  <IxdtfHighlight value={sample.input} className="text-base md:text-lg" />
                  <p className="font-sans text-muted-foreground text-sm">{t(sample.note_key)}</p>
                  <div className="flex flex-wrap items-center gap-2">
                    <TagBadge variant={sample.expect_server === "ok" ? "match" : "mismatch"}>
                      {t("guide.expectServer")}:{" "}
                      {t(sample.expect_server === "ok" ? "common.ok" : "common.error")}
                    </TagBadge>
                    <TagBadge variant={sample.expect_browser === "ok" ? "match" : "mismatch"}>
                      {t("guide.expectBrowser")}:{" "}
                      {t(sample.expect_browser === "ok" ? "common.ok" : "common.error")}
                    </TagBadge>
                    {sample.strict && <TagBadge variant="default">{t("common.strict")}</TagBadge>}
                    <Link
                      to="/playground"
                      search={{
                        input: sample.input,
                        strict: sample.strict || undefined,
                      }}
                      className={cn(
                        buttonVariants({ variant: "secondary", size: "sm" }),
                        "ml-auto",
                      )}
                    >
                      {t("common.tryInPlayground")}
                    </Link>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        ))}
      </section>
    </div>
  );
}
