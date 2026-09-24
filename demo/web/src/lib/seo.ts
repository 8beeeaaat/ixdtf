/**
 * SEO 用の画面メタ情報カタログ (N-7, D-15)。
 * 文言は各画面の locale キーを流用し (N-2)、ここではキーと正規 URL の組み立てだけを持つ。
 * public/sitemap.xml はこのカタログと一致させる (seo.test.ts が検証する)。
 */

/** canonical / sitemap に使う正規オリジン。workers.dev は worker-entry.mjs がここへ 301 する。 */
export const SITE_ORIGIN = "https://ixdtf.8beeeaaat.com";

interface PageMeta {
  /** 画面名の locale キー。`/` 以外は「画面名 · サイト名」を title にする。 */
  titleKey: string;
  /** meta description の locale キー。 */
  descriptionKey: string;
}

/**
 * インデックス対象の画面。リダイレクト専用の旧パス (/interop /temporal-lab) は含めない。
 * /playground はワークベンチ (F-2) を描画するため workbench.* を使う。
 */
export const PAGE_META = {
  "/": { titleKey: "app.title", descriptionKey: "app.tagline" },
  "/about": { titleKey: "about.title", descriptionKey: "about.tagline" },
  "/playground": { titleKey: "workbench.title", descriptionKey: "workbench.tagline" },
  "/converter": { titleKey: "converter.title", descriptionKey: "converter.tagline" },
  "/guide": { titleKey: "guide.title", descriptionKey: "guide.tagline" },
  "/support": { titleKey: "support.title", descriptionKey: "support.tagline" },
} as const satisfies Record<string, PageMeta>;

type IndexedPath = keyof typeof PAGE_META;

export interface DocumentHead {
  title: string;
  description: string;
  /** 未定義パスでは付けない。 */
  canonical?: string;
  /** 未定義パス (ソフト 404) を検索結果に載せないための noindex。 */
  noindex: boolean;
}

function isIndexedPath(path: string): path is IndexedPath {
  return Object.hasOwn(PAGE_META, path);
}

/** 末尾スラッシュを除いたパス。`/` はそのまま。 */
function normalizePath(pathname: string): string {
  return pathname.replace(/\/+$/, "") || "/";
}

/**
 * パスと翻訳関数から head の内容を決める。canonical は共有 URL のクエリ (F-2-7) を含めない。
 */
export function resolveDocumentHead(pathname: string, t: (key: string) => string): DocumentHead {
  const path = normalizePath(pathname);
  const siteTitle = t("app.title");
  if (!isIndexedPath(path)) {
    return { title: siteTitle, description: t("app.tagline"), noindex: true };
  }
  const meta = PAGE_META[path];
  return {
    title: path === "/" ? siteTitle : `${t(meta.titleKey)} · ${siteTitle}`,
    description: t(meta.descriptionKey),
    canonical: `${SITE_ORIGIN}${path}`,
    noindex: false,
  };
}

function upsertMeta(doc: Document, name: string, content: string): void {
  let meta = doc.head.querySelector<HTMLMetaElement>(`meta[name="${name}"]`);
  if (!meta) {
    meta = doc.createElement("meta");
    meta.name = name;
    doc.head.append(meta);
  }
  meta.content = content;
}

/** DocumentHead を実際の <head> に反映する。canonical / robots は不要なら取り除く。 */
export function applyDocumentHead(doc: Document, head: DocumentHead): void {
  doc.title = head.title;
  upsertMeta(doc, "description", head.description);

  const canonical = doc.head.querySelector<HTMLLinkElement>('link[rel="canonical"]');
  if (head.canonical) {
    const link = canonical ?? doc.createElement("link");
    link.rel = "canonical";
    link.href = head.canonical;
    if (!canonical) {
      doc.head.append(link);
    }
  } else {
    canonical?.remove();
  }

  if (head.noindex) {
    upsertMeta(doc, "robots", "noindex");
  } else {
    doc.head.querySelector('meta[name="robots"]')?.remove();
  }
}
