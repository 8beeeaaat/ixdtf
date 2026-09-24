import type { Meta, StoryObj } from "@storybook/react-vite";
import { expect, userEvent, within } from "storybook/test";
import { Tooltip, TooltipProvider } from "./Tooltip";

const meta = {
  component: Tooltip,
} satisfies Meta<typeof Tooltip>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    content: "IANA タイムゾーン名の注釈 [Area/Location]",
    children: <span className="font-mono text-ixdtf-timezone">[Asia/Tokyo]</span>,
  },
  render: (args) => (
    <TooltipProvider>
      <Tooltip {...args} />
    </TooltipProvider>
  ),
  // トリガーが text-ixdtf-* トークンを使っていることを assert (DESIGN.md Storybook 規則)
  play: async ({ canvasElement }) => {
    const trigger = canvasElement.querySelector("span");
    await expect(trigger?.className).toContain("text-ixdtf-timezone");
  },
};

// ホバーで注釈が表示される挙動 (Radix Portal に描画されるため document.body を探索)
export const ShowsOnHover: Story = {
  args: {
    content: "annotation-shown",
    children: (
      <button type="button" className="font-mono text-ixdtf-timezone">
        [Asia/Tokyo]
      </button>
    ),
  },
  render: (args) => (
    <TooltipProvider>
      <Tooltip {...args} />
    </TooltipProvider>
  ),
  play: async ({ canvasElement }) => {
    const trigger = within(canvasElement).getByRole("button");
    await userEvent.hover(trigger);
    const matches = await within(document.body).findAllByText("annotation-shown");
    await expect(matches.length).toBeGreaterThan(0);
  },
};
