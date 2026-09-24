import type { Meta, StoryObj } from "@storybook/react-vite";
import { expect } from "storybook/test";
import { Card, CardContent, CardHeader, CardTitle } from "./Card";

const meta = {
  component: Card,
} satisfies Meta<typeof Card>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  render: () => (
    <Card className="max-w-md">
      <CardHeader>
        <CardTitle>Section title</CardTitle>
      </CardHeader>
      <CardContent>
        <p className="font-mono text-sm tabular-nums">2026-07-07T23:30:00+09:00[Asia/Tokyo]</p>
      </CardContent>
    </Card>
  ),
  // セマンティックトークンのみ使用していることを assert (DESIGN.md Storybook 規則)
  play: async ({ canvasElement }) => {
    const card = canvasElement.querySelector("div");
    await expect(card?.className).toContain("bg-card");
    await expect(card?.className).toContain("border-border");
    await expect(card?.className).not.toMatch(/bg-(white|gray|slate)|text-\[#|dark:/);
  },
};
