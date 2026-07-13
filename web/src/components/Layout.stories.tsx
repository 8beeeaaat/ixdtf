import type { Meta, StoryObj } from "@storybook/react-vite";
import { createRootRoute, createRoute, createRouter, RouterProvider } from "@tanstack/react-router";
import { expect, userEvent, waitFor, within } from "storybook/test";
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
  const routes = [
    "/",
    "/about",
    "/playground",
    "/interop",
    "/converter",
    "/temporal-lab",
    "/guide",
    "/support",
  ].map((path) =>
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
    await expect(canvas.getByText("IXDTF Demo")).toHaveClass("font-mono", "tabular-nums");
    await expect(canvas.getByRole("link", { name: "Home" })).toBeInTheDocument();
    await expect(canvas.getByRole("link", { name: "Playground" })).toBeInTheDocument();
    await expect(canvas.getByRole("link", { name: "Temporal Lab" })).toBeInTheDocument();
    await expect(canvas.getByRole("heading", { name: "Story route" })).toBeInTheDocument();
    await expect(canvasElement.querySelector("header")?.className).toContain("border-border");

    const headerNav = canvasElement.querySelector("header nav");
    if (!(headerNav instanceof HTMLElement)) {
      throw new Error("Header navigation not found");
    }
    const nav = within(headerNav);
    const mainLinks = nav.getAllByRole("link");
    await expect(mainLinks).toHaveLength(7);
    for (const name of [
      "Home",
      "About",
      "Playground",
      "Interop",
      "Converter",
      "Temporal Lab",
      "Guide",
    ]) {
      const link = nav.getByRole("link", { name });
      await userEvent.click(link);
      await waitFor(() => expect(link).toHaveClass("text-foreground"));
    }
  },
};

export const SwitchesLanguageAndTheme: Story = {
  render: () => <StoryApp />,
  play: async ({ canvasElement }) => {
    const canvas = within(canvasElement);
    await userEvent.click(canvas.getByRole("button", { name: "日本語" }));
    await expect(await canvas.findByRole("link", { name: "ホーム" })).toBeInTheDocument();

    // 言語ボタンも「切替先」を表示する (ja へ切替後の次は English)
    await expect(canvas.getByTitle("言語")).toHaveTextContent("English");

    // テーマ・言語ボタンはどちらも「切替先」を表示する (auto の次はライト)
    const themeButton = canvas.getByTitle("テーマ");
    await expect(themeButton).toHaveTextContent("ライト");
    await userEvent.click(themeButton);
    await expect(document.documentElement).toHaveClass("light");
    await expect(themeButton).toHaveTextContent("ダーク");
  },
};
