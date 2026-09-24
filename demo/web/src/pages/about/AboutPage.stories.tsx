import type { Meta, StoryObj } from "@storybook/react-vite";
import { createRootRoute, createRoute, createRouter, RouterProvider } from "@tanstack/react-router";
import { I18nextProvider } from "react-i18next";
import { expect, userEvent, waitFor, within } from "storybook/test";
import i18n from "@/app/i18n";
import { TooltipProvider } from "@/components/ui/Tooltip";
import { AboutPage } from "@/pages/about/AboutPage";

// AboutPage は Link を使うため、遷移先パスを持つ最小のルーターで包む
function createStoryRouter() {
  const rootRoute = createRootRoute({ component: AboutPage });
  const routes = ["/", "/playground", "/guide"].map((path) =>
    createRoute({
      getParentRoute: () => rootRoute,
      path,
      component: () => null,
    }),
  );
  return createRouter({ routeTree: rootRoute.addChildren(routes), defaultPreload: false });
}

function StoryApp() {
  void i18n.changeLanguage("en");
  return (
    <I18nextProvider i18n={i18n}>
      <TooltipProvider>
        <RouterProvider router={createStoryRouter()} />
      </TooltipProvider>
    </I18nextProvider>
  );
}

const meta = {
  component: AboutPage,
} satisfies Meta<typeof AboutPage>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  render: () => <StoryApp />,
  play: async ({ canvasElement }) => {
    const canvas = within(canvasElement);
    await expect(canvas.getByRole("heading", { name: i18n.t("about.title") })).toBeVisible();

    // F-7-1 / F-7-2 / F-7-3 の 3 セクションが揃う
    for (const key of ["about.ixdtf.heading", "about.temporal.heading", "about.relation.heading"]) {
      await expect(canvas.getByRole("heading", { name: i18n.t(key) })).toBeVisible();
    }

    // F-7-4: 一次情報への参照導線が各セクションに付く
    const triggers = canvas.getAllByRole("button", { name: i18n.t("references.trigger") });
    await expect(triggers).toHaveLength(3);

    // 各セクションに正しい参照 ID が配線されている (ダイアログタイトルで検証)
    const topicTitles = [
      "references.topics.ixdtf.title",
      "references.topics.temporalApi.title",
      "references.topics.temporalIxdtf.title",
    ];
    for (const [index, trigger] of triggers.entries()) {
      await userEvent.click(trigger);
      const dialog = await within(document.body).findByRole("dialog");
      await expect(within(dialog).getByText(i18n.t(topicTitles[index]))).toBeInTheDocument();
      await userEvent.keyboard("{Escape}");
      await waitFor(() => expect(dialog).not.toBeVisible());
    }

    // F-7-4: 関連画面への遷移導線
    for (const key of [
      "about.related.workbench",
      "about.related.guide",
      "about.related.temporalLab",
    ]) {
      await expect(canvas.getByRole("link", { name: i18n.t(key) })).toBeVisible();
    }
  },
};
