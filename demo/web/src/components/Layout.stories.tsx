import type { Meta, StoryObj } from "@storybook/react-vite";
import { createRootRoute, createRoute, createRouter, RouterProvider } from "@tanstack/react-router";
import { expect, userEvent, waitFor, within } from "storybook/test";
import { page } from "vitest/browser";
import i18n from "@/app/i18n";
import { Providers } from "@/app/providers";
import { Layout } from "@/components/Layout";

const meta = {
  component: Layout,
} satisfies Meta<typeof Layout>;

export default meta;
type Story = StoryObj<typeof meta>;

function StoryPage() {
  return (
    <section className="mx-auto max-w-6xl px-6">
      <h1 className="font-sans font-semibold text-2xl tracking-tight">Story route</h1>
      <p className="font-sans text-muted-foreground text-sm">Outlet content</p>
    </section>
  );
}

function createStoryRouter() {
  const rootRoute = createRootRoute({ component: Layout });
  const routes = ["/", "/about", "/playground", "/converter", "/guide", "/support"].map((path) =>
    createRoute({
      getParentRoute: () => rootRoute,
      path,
      component: StoryPage,
    }),
  );
  return createRouter({ routeTree: rootRoute.addChildren(routes), defaultPreload: false });
}

function StoryApp() {
  localStorage.setItem("theme", "auto");
  localStorage.setItem("lang", "en");
  void i18n.changeLanguage("en");
  return (
    <Providers>
      <RouterProvider router={createStoryRouter()} />
    </Providers>
  );
}

export const Default: Story = {
  render: () => <StoryApp />,
  play: async ({ canvasElement }) => {
    const canvas = within(canvasElement);
    // ロゴ兼ホームリンク: 可視ワードマークがアクセシブルネームを兼ねる
    await expect(canvas.getByText("IXDTF Demo")).toHaveClass("font-mono", "tabular-nums");
    await expect(canvas.getByRole("link", { name: "IXDTF Demo" })).toBeInTheDocument();
    await expect(canvas.getByRole("link", { name: "Workbench" })).toBeInTheDocument();
    await expect(canvas.getByRole("heading", { name: "Story route" })).toBeInTheDocument();
    await expect(canvasElement.querySelector("header")?.className).toContain("border-border");

    const headerNav = canvasElement.querySelector("header nav");
    if (!(headerNav instanceof HTMLElement)) {
      throw new Error("Header navigation not found");
    }
    const nav = within(headerNav);
    // ホームはロゴ側へ移したため、ナビのテキストリンクは 4 件 (IXDTF と Temporal / 変換 / ワークベンチ / ガイド)
    const mainLinks = nav.getAllByRole("link");
    await expect(mainLinks).toHaveLength(4);
    for (const name of ["IXDTF & Temporal", "Converter", "Workbench", "Guide"]) {
      const link = nav.getByRole("link", { name });
      await userEvent.click(link);
      await waitFor(() => expect(link).toHaveClass("text-foreground"));
    }
  },
};

export const MobileBottomTabBar: Story = {
  render: () => <StoryApp />,
  play: async ({ canvasElement }) => {
    // 375px 幅 (iPhone SE 相当) を再現する。修正前はヘッダーのテキストナビが
    // 実測 0 幅まで潰れ、ページ全体が横スクロールしていた不具合の回帰ガード (F-0-5)。
    await page.viewport(375, 812);
    const canvas = within(canvasElement);

    try {
      // md 未満ではデスクトップ用テキストナビ (hidden md:flex) が非表示になる。
      // "Workbench" 等はボトムタブバーとラベルを共有するため、テキストではなく
      // primary nav 要素自体の可視性を直接検証する。
      const primaryNav = canvasElement.querySelector('nav[aria-label="Primary navigation"]');
      if (!(primaryNav instanceof HTMLElement)) {
        throw new Error("Primary navigation not found");
      }
      await waitFor(() => expect(primaryNav).not.toBeVisible());

      // 代わりにボトムタブバーが現れ、Home を含む 5 画面すべてに到達できる
      const bottomNav = canvasElement.querySelector('nav[aria-label="Mobile navigation"]');
      if (!(bottomNav instanceof HTMLElement)) {
        throw new Error("Mobile bottom tab bar not found");
      }
      const tabs = within(bottomNav);
      for (const name of ["Home", "About", "Converter", "Workbench", "Guide"]) {
        await expect(tabs.getByRole("link", { name })).toBeVisible();
      }

      // タブをタップすると実際に遷移し、選択中タブが text-foreground になる
      const converterTab = tabs.getByRole("link", { name: "Converter" });
      await userEvent.click(converterTab);
      await waitFor(() => expect(converterTab).toHaveClass("text-foreground"));

      // ヘッダー右クラスタの Select はモバイルではラベルを隠しアイコンのみになる
      const langTrigger = canvas.getByRole("combobox", { name: "Language" });
      const langLabel = within(langTrigger).getByText("English");
      await expect(langLabel).not.toBeVisible();

      // ページ全体が横スクロールしない (修正前は scrollWidth > clientWidth だった)
      await waitFor(() => {
        expect(document.documentElement.scrollWidth).toBeLessThanOrEqual(
          document.documentElement.clientWidth,
        );
      });

      // md ブレークポイント (768px) の境界: 767 ではボトムタブバー、768 ではデスクトップナビ
      await page.viewport(767, 800);
      await waitFor(() => expect(bottomNav).toBeVisible());
      await expect(primaryNav).not.toBeVisible();

      await page.viewport(768, 800);
      await waitFor(() => expect(primaryNav).toBeVisible());
      await expect(bottomNav).not.toBeVisible();
    } finally {
      // 後続 story への影響を避けるためデスクトップ幅へ戻す
      await page.viewport(1280, 800);
    }
  },
};

export const SwitchesLanguageAndTheme: Story = {
  render: () => <StoryApp />,
  play: async ({ canvasElement }) => {
    const canvas = within(canvasElement);
    // Radix Select はオプションを document.body の Portal に描画する
    const body = within(document.body);

    // 言語 Select は現在値 (English) を表示する。プルダウンから日本語を選ぶ
    const langTrigger = canvas.getByRole("combobox", { name: "Language" });
    await expect(langTrigger).toHaveTextContent("English");
    await userEvent.click(langTrigger);
    await userEvent.click(await body.findByRole("option", { name: "日本語" }));
    await expect(
      await canvas.findByRole("link", { name: "IXDTF と Temporal" }),
    ).toBeInTheDocument();

    // テーマ Select は現在値 (自動) を表示する。ライトへ切り替えると <html> に light が付く
    const themeTrigger = canvas.getByRole("combobox", { name: "テーマ" });
    await expect(themeTrigger).toHaveTextContent("自動");
    await userEvent.click(themeTrigger);
    await userEvent.click(await body.findByRole("option", { name: "ライト" }));
    await expect(document.documentElement).toHaveClass("light");
    await expect(themeTrigger).toHaveTextContent("ライト");
  },
};
