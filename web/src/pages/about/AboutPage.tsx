import { Link } from "@tanstack/react-router";
import { Trans, useTranslation } from "react-i18next";
import { IxdtfHighlight, IxdtfLegend } from "@/components/IxdtfHighlight";
import { ReferenceDialog } from "@/components/ReferenceDialog";
import { buttonVariants } from "@/components/ui/Button";
import { Card, CardContent } from "@/components/ui/Card";
import { requiredSample } from "@/lib/fixtures";
import type { ReferenceId } from "@/lib/references";
import { cn } from "@/lib/utils";

// F-7-1 の構造例: タイムゾーン注釈とカレンダー注釈の両方を持つ共有 fixture
const ANATOMY_SAMPLE = requiredSample("calendar-japanese");

interface SectionHeadingProps {
  headingKey: string;
  referenceId: ReferenceId;
}

function SectionHeading({ headingKey, referenceId }: SectionHeadingProps) {
  const { t } = useTranslation();
  return (
    <div className="flex flex-wrap items-center gap-1">
      <h2 className="font-sans font-semibold text-muted-foreground text-sm uppercase tracking-wider">
        {t(headingKey)}
      </h2>
      <ReferenceDialog referenceId={referenceId} />
    </div>
  );
}

/** F-7: what Temporal (the API) and IXDTF (the standard) each are, and how they relate. */
export function AboutPage() {
  const { t } = useTranslation();

  return (
    <div className="mx-auto max-w-6xl space-y-12 px-6">
      <header>
        <h1 className="font-sans font-semibold text-2xl tracking-tight">{t("about.title")}</h1>
        <p className="mt-1 font-sans text-muted-foreground text-sm">{t("about.tagline")}</p>
      </header>

      {/* F-7-1: IXDTF はデータ形式の規格 */}
      <section className="space-y-4">
        <SectionHeading headingKey="about.ixdtf.heading" referenceId="ixdtf" />
        <p className="max-w-3xl font-sans text-sm leading-relaxed">{t("about.ixdtf.body1")}</p>
        <p className="max-w-3xl font-sans text-sm leading-relaxed">{t("about.ixdtf.body2")}</p>
        <Card>
          <CardContent className="space-y-3">
            <IxdtfHighlight value={ANATOMY_SAMPLE.input} className="text-base md:text-lg" />
            <IxdtfLegend />
            <p className="font-sans text-muted-foreground text-sm">
              {t("about.ixdtf.exampleCaption")}
            </p>
          </CardContent>
        </Card>
      </section>

      {/* F-7-2: Temporal は ECMAScript の日時 API */}
      <section className="space-y-4">
        <SectionHeading headingKey="about.temporal.heading" referenceId="temporalApi" />
        <p className="max-w-3xl font-sans text-sm leading-relaxed">{t("about.temporal.body1")}</p>
        <p className="max-w-3xl font-sans text-sm leading-relaxed">
          <Trans
            i18nKey="about.temporal.body2"
            components={{
              polyfillLink: (
                // biome-ignore lint/a11y/useAnchorContent: <Trans> injects the link text from the translation string
                <a
                  href="https://github.com/fullcalendar/temporal-polyfill"
                  target="_blank"
                  rel="noreferrer"
                  className="text-info underline underline-offset-4 hover:opacity-80"
                />
              ),
            }}
          />
        </p>
      </section>

      {/* F-7-3: 両者の関係 — API と交換形式 */}
      <section className="space-y-4">
        <SectionHeading headingKey="about.relation.heading" referenceId="temporalIxdtf" />
        <p className="max-w-3xl font-sans text-sm leading-relaxed">{t("about.relation.body1")}</p>
        <p className="max-w-3xl font-sans text-sm leading-relaxed">{t("about.relation.body2")}</p>
        <p className="max-w-3xl font-sans text-sm leading-relaxed">{t("about.relation.body3")}</p>
        <Card>
          <CardContent className="space-y-2">
            <div className="flex flex-wrap items-center justify-center gap-x-3 gap-y-1 text-center font-sans text-sm">
              <span>{t("about.relation.diagramBrowser")}</span>
              <span aria-hidden="true">⇄</span>
              <span className="font-semibold">{t("about.relation.diagramString")}</span>
              <span aria-hidden="true">⇄</span>
              <span>{t("about.relation.diagramGo")}</span>
            </div>
            <p className="text-center font-sans text-muted-foreground text-sm">
              {t("about.relation.diagramCaption")}
            </p>
          </CardContent>
        </Card>
      </section>

      {/* F-7-4: 関連画面への導線 */}
      <section className="space-y-4">
        <h2 className="font-sans font-semibold text-muted-foreground text-sm uppercase tracking-wider">
          {t("about.related.heading")}
        </h2>
        <div className="flex flex-wrap gap-2">
          <Link
            to="/playground"
            className={cn(buttonVariants({ variant: "secondary", size: "sm" }))}
          >
            {t("about.related.workbench")}
          </Link>
          <Link to="/guide" className={cn(buttonVariants({ variant: "secondary", size: "sm" }))}>
            {t("about.related.guide")}
          </Link>
          {/* F-6 は F-2 ワークベンチの Lab モードになったため mode=lab で遷移する (旧 /temporal-lab) */}
          <Link
            to="/playground"
            search={{ mode: "lab" }}
            className={cn(buttonVariants({ variant: "secondary", size: "sm" }))}
          >
            {t("about.related.temporalLab")}
          </Link>
        </div>
      </section>
    </div>
  );
}
