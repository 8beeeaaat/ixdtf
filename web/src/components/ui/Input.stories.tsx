import type { Meta, StoryObj } from "@storybook/react-vite";
import { expect, fn, userEvent, within } from "storybook/test";
import { Input } from "./Input";

const meta = {
  component: Input,
} satisfies Meta<typeof Input>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: { placeholder: "2026-07-07T23:30:00+09:00[Asia/Tokyo]" },
  // セマンティックトークン + font-mono tabular-nums を assert (DESIGN.md)
  play: async ({ canvasElement }) => {
    const input = canvasElement.querySelector("input");
    await expect(input?.className).toContain("border-border");
    await expect(input?.className).toContain("font-mono");
    await expect(input?.className).not.toMatch(/bg-(white|gray|slate)|text-\[#|dark:/);
  },
};

export const Disabled: Story = {
  args: { disabled: true, value: "2026-07-07T23:30:00Z" },
};

// 入力で value が更新され onChange が発火する挙動
export const TypesValue: Story = {
  args: { placeholder: "type an IXDTF string", onChange: fn() },
  play: async ({ args, canvasElement }) => {
    const input = within(canvasElement).getByRole("textbox");
    await userEvent.type(input, "2026-07-07");
    await expect(input).toHaveValue("2026-07-07");
    await expect(args.onChange).toHaveBeenCalled();
  },
};
