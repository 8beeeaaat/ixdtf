import type { Meta, StoryObj } from "@storybook/react-vite";
import { createRootRoute, createRoute, createRouter, RouterProvider } from "@tanstack/react-router";
import { expect, within } from "storybook/test";
import i18n from "@/app/i18n";
import { Providers } from "@/app/providers";
import { Footer } from "@/components/Footer";

const meta = {
  component: Footer,
} satisfies Meta<typeof Footer>;

export default meta;
type Story = StoryObj<typeof meta>;

/** Footer uses TanStack Router <Link>, so the story needs a router context. */
function createStoryRouter() {
  const rootRoute = createRootRoute({ component: Footer });
  const routes = ["/", "/support"].map((path) =>
    createRoute({ getParentRoute: () => rootRoute, path, component: () => null }),
  );
  return createRouter({ routeTree: rootRoute.addChildren(routes), defaultPreload: false });
}

function StoryApp() {
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
    await expect(await canvas.findByText("This demo is free and open source")).toBeInTheDocument();

    const support = canvas.getByRole("link", { name: "Support" });
    await expect(support).toHaveClass("text-foreground");

    const footer = canvasElement.querySelector("footer");
    await expect(footer?.className).toContain("border-border");
    // DESIGN.md: no raw color classes, no dark: variant.
    await expect(footer?.className).not.toMatch(/bg-(white|gray|slate)|text-\[#|dark:/);
  },
};
