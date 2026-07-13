import type { Meta, StoryObj } from "@storybook/react-vite";
import { expect } from "storybook/test";
import { FormField } from "./FormField";
import { Input } from "./Input";

const meta = {
  component: FormField,
} satisfies Meta<typeof FormField>;

export default meta;
type Story = StoryObj<typeof meta>;

export const WithHint: Story = {
  args: {
    label: "IXDTF 文字列",
    hint: "例: 2026-07-07T23:30:00+09:00[Asia/Tokyo]",
    children: <Input placeholder="2026-07-07T23:30:00+09:00[Asia/Tokyo]" />,
  },
  // ラベルがセマンティックトークンを使っていることを assert (DESIGN.md Storybook 規則)
  play: async ({ canvasElement }) => {
    const label = canvasElement.querySelector("label");
    await expect(label?.className).toContain("text-muted-foreground");
    await expect(label?.className).not.toMatch(/text-(gray|slate)-|text-\[#|dark:/);
  },
};

export const WithError: Story = {
  args: {
    label: "IXDTF 文字列",
    error: "基準時刻を解析できませんでした",
    children: <Input value="not-a-timestamp" readOnly />,
  },
};
