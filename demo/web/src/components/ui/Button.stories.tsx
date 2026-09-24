import type { Meta, StoryObj } from "@storybook/react-vite";
import { expect } from "storybook/test";
import { Button } from "./Button";

const meta = {
  component: Button,
  args: { children: "Button" },
} satisfies Meta<typeof Button>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Primary: Story = {
  args: { variant: "primary" },
  // セマンティックトークンのみ使用していることを assert (DESIGN.md Storybook 規則)
  play: async ({ canvasElement }) => {
    const button = canvasElement.querySelector("button");
    await expect(button?.className).toContain("bg-primary");
    await expect(button?.className).not.toMatch(/bg-(white|gray|slate)|text-\[#/);
  },
};

export const Secondary: Story = { args: { variant: "secondary" } };
export const Ghost: Story = { args: { variant: "ghost" } };
export const Small: Story = { args: { size: "sm" } };
export const Disabled: Story = { args: { disabled: true } };
