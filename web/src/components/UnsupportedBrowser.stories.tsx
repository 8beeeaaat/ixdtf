import type { Meta, StoryObj } from "@storybook/react-vite";
import { I18nextProvider } from "react-i18next";
import { expect, within } from "storybook/test";
import i18n from "@/app/i18n";
import { UnsupportedBrowser } from "@/components/UnsupportedBrowser";

const meta = {
  component: UnsupportedBrowser,
  decorators: [
    (Story) => (
      <I18nextProvider i18n={i18n}>
        <Story />
      </I18nextProvider>
    ),
  ],
} satisfies Meta<typeof UnsupportedBrowser>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  play: async ({ canvasElement }) => {
    const canvas = within(canvasElement);
    await expect(canvas.getByRole("heading", { name: /Temporal|未対応/i })).toBeInTheDocument();
    const browserList = canvas.getByText(/Chrome 144\+.*Edge 144\+.*Firefox 139\+/);
    await expect(browserList).toHaveClass("font-mono", "tabular-nums");
    await expect(canvasElement.querySelector("main")?.className).toContain("max-w-6xl");
  },
};
