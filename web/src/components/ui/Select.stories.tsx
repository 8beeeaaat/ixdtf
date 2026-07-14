import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { expect, userEvent, within } from "storybook/test";
import { Select } from "@/components/ui/Select";

const OPTIONS = [
  { value: "auto", label: "Auto" },
  { value: "light", label: "Light" },
  { value: "dark", label: "Dark" },
];

const meta = {
  component: Select,
  args: {
    "aria-label": "Theme",
    value: "auto",
    onValueChange: () => {},
    options: OPTIONS,
  },
} satisfies Meta<typeof Select>;

export default meta;
type Story = StoryObj<typeof meta>;

function Controlled() {
  const [value, setValue] = useState("auto");
  return <Select aria-label="Theme" value={value} onValueChange={setValue} options={OPTIONS} />;
}

/** トリガーは常に現在値を表示し、プルダウンから別の値を選ぶと即座に反映される。 */
export const Default: Story = {
  render: () => <Controlled />,
  play: async ({ canvasElement }) => {
    const canvas = within(canvasElement);
    // Radix Select はオプションを document.body の Portal に描画する
    const body = within(document.body);
    const trigger = canvas.getByRole("combobox", { name: "Theme" });
    await expect(trigger).toHaveTextContent("Auto");
    await userEvent.click(trigger);
    await userEvent.click(await body.findByRole("option", { name: "Dark" }));
    await expect(trigger).toHaveTextContent("Dark");
  },
};
