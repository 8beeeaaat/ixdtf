import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { expect, userEvent, within } from "storybook/test";
import { ToggleSwitch } from "./ToggleSwitch";

const meta = {
  component: ToggleSwitch,
} satisfies Meta<typeof ToggleSwitch>;

export default meta;
type Story = StoryObj<typeof meta>;

function InteractiveToggle({ initial }: { initial: boolean }) {
  const [enabled, setEnabled] = useState(initial);
  return <ToggleSwitch enabled={enabled} onChange={setEnabled} label="strict" />;
}

export const Off: Story = {
  args: { enabled: false, onChange: () => {}, label: "strict" },
  render: () => <InteractiveToggle initial={false} />,
  // セマンティックトークンのみ使用していることを assert (DESIGN.md Storybook 規則)
  play: async ({ canvasElement }) => {
    const root = canvasElement.querySelector('button[role="switch"]');
    await expect(root?.className).toContain("bg-muted");
    await expect(root?.className).not.toMatch(/bg-(white|gray|slate)|text-\[#|dark:/);
  },
};

export const On: Story = {
  args: { enabled: true, onChange: () => {}, label: "strict" },
  render: () => <InteractiveToggle initial={true} />,
};

// クリックで off → on にトグルする挙動 (aria-checked が反転)
export const TogglesOnClick: Story = {
  args: { enabled: false, onChange: () => {}, label: "strict" },
  render: () => <InteractiveToggle initial={false} />,
  play: async ({ canvasElement }) => {
    const sw = within(canvasElement).getByRole("switch");
    await expect(sw).toHaveAttribute("aria-checked", "false");
    await userEvent.click(sw);
    await expect(sw).toHaveAttribute("aria-checked", "true");
  },
};
