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

function DotIcon() {
  return (
    <svg viewBox="0 0 16 16" aria-hidden="true" className="h-4 w-4 shrink-0">
      <circle cx="8" cy="8" r="3" fill="currentColor" />
    </svg>
  );
}

function ControlledWithIcon() {
  const [value, setValue] = useState("auto");
  return (
    <Select
      aria-label="Theme"
      value={value}
      onValueChange={setValue}
      options={OPTIONS}
      icon={<DotIcon />}
    />
  );
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

/** icon 指定時は sm 未満でラベルを隠しアイコンのみにする (ヘッダー右クラスタのモバイル縮退、F-0-5)。 */
export const WithIcon: Story = {
  render: () => <ControlledWithIcon />,
  play: async ({ canvasElement }) => {
    const canvas = within(canvasElement);
    const trigger = canvas.getByRole("combobox", { name: "Theme" });
    const label = within(trigger).getByText("Auto");
    await expect(label).toHaveClass("hidden", "sm:inline");
  },
};
