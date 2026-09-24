import { readFileSync } from "node:fs";
import { resolve } from "node:path";
import { afterEach, describe, expect, it } from "vitest";
import i18n from "@/app/i18n";
import { router } from "@/app/router";
import { applyDocumentHead, PAGE_META, resolveDocumentHead, SITE_ORIGIN } from "@/lib/seo";

const t = (key: string) => i18n.t(key, { lng: "en" });

function readPublic(file: string): string {
  return readFileSync(resolve(process.cwd(), "public", file), "utf8");
}

describe("SEO page catalog (N-7)", () => {
  it("lists exactly the catalog pages in sitemap.xml", () => {
    const locs = [...readPublic("sitemap.xml").matchAll(/<loc>([^<]+)<\/loc>/g)].map((m) => m[1]);
    expect(locs).toEqual(Object.keys(PAGE_META).map((path) => `${SITE_ORIGIN}${path}`));
  });

  it("points robots.txt at the sitemap", () => {
    expect(readPublic("robots.txt")).toContain(`Sitemap: ${SITE_ORIGIN}/sitemap.xml`);
  });

  it("only catalogs routes that render a page", () => {
    const routePaths = Object.keys(router.routesByPath);
    for (const path of Object.keys(PAGE_META)) {
      expect(routePaths, path).toContain(path);
    }
    // リダイレクト専用の旧パスはインデックス対象にしない。
    expect(Object.keys(PAGE_META)).not.toContain("/interop");
    expect(Object.keys(PAGE_META)).not.toContain("/temporal-lab");
  });

  it("keeps every title / description key available in both locales", () => {
    for (const meta of Object.values(PAGE_META)) {
      for (const language of ["ja", "en"]) {
        expect(i18n.exists(meta.titleKey, { lng: language }), meta.titleKey).toBe(true);
        expect(i18n.exists(meta.descriptionKey, { lng: language }), meta.descriptionKey).toBe(true);
      }
    }
  });
});

describe("resolveDocumentHead", () => {
  it("uses the site title alone for the home page", () => {
    const head = resolveDocumentHead("/", t);
    expect(head.title).toBe(t("app.title"));
    expect(head.canonical).toBe(`${SITE_ORIGIN}/`);
    expect(head.noindex).toBe(false);
  });

  it("titles the workbench from workbench.* and drops the trailing slash from canonical", () => {
    const head = resolveDocumentHead("/playground/", t);
    expect(head.title).toBe(`${t("workbench.title")} · ${t("app.title")}`);
    expect(head.description).toBe(t("workbench.tagline"));
    expect(head.canonical).toBe(`${SITE_ORIGIN}/playground`);
  });

  it("marks unknown paths noindex without a canonical", () => {
    const head = resolveDocumentHead("/no-such-page", t);
    expect(head.noindex).toBe(true);
    expect(head.canonical).toBeUndefined();
  });
});

describe("applyDocumentHead", () => {
  afterEach(() => {
    document.head.replaceChildren();
  });

  it("writes title, description and canonical, then swaps to noindex for unknown paths", () => {
    applyDocumentHead(document, resolveDocumentHead("/guide", t));
    expect(document.title).toBe(`${t("guide.title")} · ${t("app.title")}`);
    expect(document.head.querySelector('meta[name="description"]')?.getAttribute("content")).toBe(
      t("guide.tagline"),
    );
    expect(document.head.querySelector('link[rel="canonical"]')?.getAttribute("href")).toBe(
      `${SITE_ORIGIN}/guide`,
    );
    expect(document.head.querySelector('meta[name="robots"]')).toBeNull();

    applyDocumentHead(document, resolveDocumentHead("/missing", t));
    expect(document.head.querySelector('link[rel="canonical"]')).toBeNull();
    expect(document.head.querySelector('meta[name="robots"]')?.getAttribute("content")).toBe(
      "noindex",
    );
    expect(document.head.querySelectorAll('meta[name="description"]')).toHaveLength(1);
  });
});
