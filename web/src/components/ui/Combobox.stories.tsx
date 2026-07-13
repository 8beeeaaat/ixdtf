import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { expect, userEvent, within } from "storybook/test";
import { Combobox } from "./Combobox";

const meta = {
  component: Combobox,
} satisfies Meta<typeof Combobox>;

export default meta;
type Story = StoryObj<typeof meta>;

const OPTIONS = ["Asia/Tokyo", "UTC", "America/Los_Angeles", "Europe/London", "Australia/Sydney"];

function InteractiveCombobox() {
  const [value, setValue] = useState("Asia/Tokyo");
  return (
    <div className="w-72">
      <Combobox value={value} options={OPTIONS} onChange={setValue} aria-label="time zone" />
    </div>
  );
}

export const Default: Story = {
  args: { value: "Asia/Tokyo", options: OPTIONS, onChange: () => {} },
  render: () => <InteractiveCombobox />,
  // セマンティックトークン + font-mono を assert (DESIGN.md Storybook 規則)
  play: async ({ canvasElement }) => {
    const input = canvasElement.querySelector('[role="combobox"]');
    await expect(input?.className).toContain("border-border");
    await expect(input?.className).toContain("font-mono");
    await expect(input?.className).not.toMatch(/bg-(white|gray|slate)|text-\[#|dark:/);
  },
};

// フィルタ入力 → 候補クリックで値が確定する挙動 (F-1-3 / F-1-4)
export const FilterAndSelect: Story = {
  args: { value: "Asia/Tokyo", options: OPTIONS, onChange: () => {} },
  render: () => <InteractiveCombobox />,
  play: async ({ canvasElement }) => {
    const canvas = within(canvasElement);
    const input = canvas.getByRole("combobox");
    await userEvent.click(input); // フォーカスで listbox が開く
    await userEvent.type(input, "UTC"); // 候補を絞り込む
    const option = await canvas.findByRole("option", { name: "UTC" });
    await userEvent.click(option); // 選択で確定
    await expect(input).toHaveValue("UTC");
  },
};
