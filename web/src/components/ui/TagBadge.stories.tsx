import type { Meta, StoryObj } from "@storybook/react-vite";
import { expect } from "storybook/test";
import { TagBadge } from "./TagBadge";

const meta = {
  component: TagBadge,
  args: { children: "u-ca=japanese" },
} satisfies Meta<typeof TagBadge>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const Critical: Story = {
  args: { variant: "critical", children: "!foo=bar" },
  // critical バリアントが warning トークン + 波下線 (N-3) を使うことを assert
  play: async ({ canvasElement }) => {
    const badge = canvasElement.querySelector("span");
    await expect(badge?.className).toContain("text-warning");
    await expect(badge?.className).toContain("decoration-wavy");
    await expect(badge?.className).not.toMatch(/text-(gray|slate)-|text-\[#/);
  },
};

export const Match: Story = { args: { variant: "match", children: "match" } };
export const Mismatch: Story = { args: { variant: "mismatch", children: "mismatch" } };
