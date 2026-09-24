import type { Meta, StoryObj } from "@storybook/react-vite";
import { I18nextProvider } from "react-i18next";
import { expect, fn, userEvent, within } from "storybook/test";
import i18n from "@/app/i18n";
import { CopyButton } from "@/components/CopyButton";

const IXDTF = "2026-07-07T23:30:00+09:00[Asia/Tokyo][u-ca=japanese]";

const meta = {
  component: CopyButton,
  decorators: [
    (Story) => (
      <I18nextProvider i18n={i18n}>
        <Story />
      </I18nextProvider>
    ),
  ],
} satisfies Meta<typeof CopyButton>;

export default meta;
type Story = StoryObj<typeof meta>;

function installClipboard(writeText: (value: string) => Promise<void>): void {
  Object.defineProperty(navigator, "clipboard", {
    value: { writeText },
    configurable: true,
  });
}

export const Default: Story = {
  args: { text: IXDTF },
  play: async ({ canvasElement }) => {
    const button = within(canvasElement).getByRole("button");
    await expect(button.className).toContain("bg-secondary");
    await expect(button.className).not.toMatch(/bg-(white|gray|slate)|text-\[#|dark:/);
  },
};

export const CopySuccess: Story = {
  args: { text: IXDTF, label: "Copy IXDTF", copiedLabel: "Copied IXDTF" },
  play: async ({ canvasElement }) => {
    const writeText = fn(async (_value: string) => {});
    installClipboard(writeText);

    const button = within(canvasElement).getByRole("button", { name: "Copy IXDTF" });
    await userEvent.click(button);

    await expect(writeText).toHaveBeenCalledWith(IXDTF);
    await expect(button).toHaveTextContent("Copied IXDTF");
  },
};

export const CopyFailureKeepsLabel: Story = {
  args: { text: IXDTF, label: "Copy IXDTF", copiedLabel: "Copied IXDTF" },
  play: async ({ canvasElement }) => {
    const writeText = fn(async (_value: string) => {
      throw new DOMException("Clipboard denied", "NotAllowedError");
    });
    installClipboard(writeText);

    const button = within(canvasElement).getByRole("button", { name: "Copy IXDTF" });
    await userEvent.click(button);

    await expect(writeText).toHaveBeenCalledWith(IXDTF);
    await expect(button).toHaveTextContent("Copy IXDTF");
  },
};
