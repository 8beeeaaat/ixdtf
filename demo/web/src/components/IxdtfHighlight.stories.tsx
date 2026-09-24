import type { Meta, StoryObj } from "@storybook/react-vite";
import { I18nextProvider } from "react-i18next";
import { expect, userEvent, within } from "storybook/test";
import i18n from "@/app/i18n";
import { IxdtfHighlight, IxdtfLegend } from "@/components/IxdtfHighlight";
import { TooltipProvider } from "@/components/ui/Tooltip";

const COMPLETE = "2026-07-07T23:30:00+09:00[Asia/Tokyo][u-ca=japanese]";
const CRITICAL = "2026-07-07T23:30:00+09:00[!Asia/Tokyo][!u-ca=japanese]";

const meta = {
  component: IxdtfHighlight,
  decorators: [
    (Story) => (
      <I18nextProvider i18n={i18n}>
        <TooltipProvider>
          <Story />
        </TooltipProvider>
      </I18nextProvider>
    ),
  ],
} satisfies Meta<typeof IxdtfHighlight>;

export default meta;
type Story = StoryObj<typeof meta>;

export const CompleteIxdtf: Story = {
  args: { value: COMPLETE },
  play: async ({ canvasElement }) => {
    await expect(canvasElement.querySelector(".text-ixdtf-date")).toHaveTextContent("2026-07-07");
    await expect(canvasElement.querySelector(".text-ixdtf-time")).toHaveTextContent("T23:30:00");
    await expect(canvasElement.querySelector(".text-ixdtf-offset")).toHaveTextContent("+09:00");
    await expect(canvasElement.querySelector(".text-ixdtf-timezone")).toHaveTextContent(
      "[Asia/Tokyo]",
    );
    await expect(canvasElement.querySelector(".text-ixdtf-extension")).toHaveTextContent(
      "[u-ca=japanese]",
    );
  },
};

export const CriticalAnnotations: Story = {
  args: { value: CRITICAL },
  play: async ({ canvasElement }) => {
    const criticalTokens = canvasElement.querySelectorAll(".decoration-warning");
    await expect(criticalTokens.length).toBe(2);
    for (const token of criticalTokens) {
      await expect(token.className).toContain("decoration-wavy");
    }
  },
};

export const InvalidInput: Story = {
  args: { value: "not-an-ixdtf-string" },
  play: async ({ canvasElement }) => {
    await expect(canvasElement.querySelector(".text-destructive")).toHaveTextContent(
      "not-an-ixdtf-string",
    );
  },
};

export const Legend: Story = {
  args: { value: COMPLETE },
  render: () => <IxdtfLegend />,
  play: async ({ canvasElement }) => {
    await expect(canvasElement.querySelector(".text-ixdtf-date")).toBeInTheDocument();
    await expect(canvasElement.querySelector(".text-ixdtf-time")).toBeInTheDocument();
    await expect(canvasElement.querySelector(".text-ixdtf-offset")).toBeInTheDocument();
    await expect(canvasElement.querySelector(".text-ixdtf-timezone")).toBeInTheDocument();
    await expect(canvasElement.querySelector(".text-ixdtf-extension")).toBeInTheDocument();
    await expect(canvasElement.querySelector(".text-warning")).toHaveClass("decoration-wavy");
  },
};

export const TooltipOnHover: Story = {
  args: { value: COMPLETE },
  play: async ({ canvasElement }) => {
    const token = within(canvasElement).getByText("[Asia/Tokyo]");
    await userEvent.hover(token);
    const matches = await within(document.body).findAllByText(/time-zone|タイムゾーン/i);
    await expect(matches.length).toBeGreaterThan(0);
  },
};
