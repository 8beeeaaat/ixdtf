import type { Meta, StoryObj } from "@storybook/react-vite";
import { I18nextProvider } from "react-i18next";
import { expect, fn, userEvent, within } from "storybook/test";
import i18n from "@/app/i18n";
import { ReproduceSection } from "@/components/ReproduceSection";
import { goParseSample, jsParseSample } from "@/lib/reproduce";

const INPUT = "2026-07-07T23:30:00+09:00[Asia/Tokyo][u-ca=japanese]";

const meta = {
  component: ReproduceSection,
  decorators: [
    (Story) => (
      <I18nextProvider i18n={i18n}>
        <Story />
      </I18nextProvider>
    ),
  ],
  args: {
    goSample: goParseSample({ input: INPUT, strict: false, validateOnly: false }),
    jsSample: jsParseSample(INPUT),
  },
} satisfies Meta<typeof ReproduceSection>;

export default meta;
type Story = StoryObj<typeof meta>;

function installClipboard(writeText: (value: string) => Promise<void>): void {
  Object.defineProperty(navigator, "clipboard", {
    value: { writeText },
    configurable: true,
  });
}

export const ClosedByDefault: Story = {
  play: async ({ canvasElement }) => {
    const details = canvasElement.querySelector("details");
    await expect(details?.open).toBe(false);
    // 閉じた <details> の中身は DOM には残るが不可視 (F-0-7: 初期状態は閉)
    await expect(canvasElement.querySelector("pre")).not.toBeVisible();
  },
};

export const OpenShowsGoSample: Story = {
  play: async ({ canvasElement }) => {
    const canvas = within(canvasElement);
    await userEvent.click(canvas.getByText(i18n.t("reproduce.summary")));
    const pre = canvasElement.querySelector("pre");
    await expect(pre?.textContent).toContain("ixdtf.Parse(input, false)");
    await expect(pre?.textContent).toContain(`input := "${INPUT}"`);
    await expect(pre?.className).not.toMatch(/bg-(white|gray|slate)|text-\[#|dark:/);
  },
};

export const SwitchTabs: Story = {
  play: async ({ canvasElement }) => {
    const canvas = within(canvasElement);
    await userEvent.click(canvas.getByText(i18n.t("reproduce.summary")));

    const jsButton = canvas.getByRole("button", { name: i18n.t("reproduce.tabs.js") });
    await userEvent.click(jsButton);
    await expect(jsButton).toHaveAttribute("aria-pressed", "true");
    await expect(canvas.getByRole("button", { name: i18n.t("reproduce.tabs.go") })).toHaveAttribute(
      "aria-pressed",
      "false",
    );
    await expect(canvasElement.querySelector("pre")?.textContent).toContain(
      "Temporal.ZonedDateTime.from(input)",
    );
  },
};

export const CopyReflectsActiveTab: Story = {
  play: async ({ canvasElement }) => {
    const writeText = fn(async (_value: string) => {});
    installClipboard(writeText);

    const canvas = within(canvasElement);
    await userEvent.click(canvas.getByText(i18n.t("reproduce.summary")));
    await userEvent.click(canvas.getByRole("button", { name: i18n.t("reproduce.tabs.js") }));
    await userEvent.click(canvas.getByRole("button", { name: i18n.t("common.copy") }));

    await expect(writeText).toHaveBeenCalledWith(
      expect.stringContaining("Temporal.ZonedDateTime.from(input)"),
    );
  },
};
