import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { I18nextProvider } from "react-i18next";
import { expect, userEvent, within } from "storybook/test";
import i18n from "@/app/i18n";
import { TimeZonePicker } from "@/components/TimeZonePicker";

const meta = {
  component: TimeZonePicker,
  decorators: [
    (Story) => (
      <I18nextProvider i18n={i18n}>
        <div className="w-80">
          <Story />
        </div>
      </I18nextProvider>
    ),
  ],
} satisfies Meta<typeof TimeZonePicker>;

export default meta;
type Story = StoryObj<typeof meta>;

function InteractiveTimeZonePicker({ initial }: { initial: string }) {
  const [timeZone, setTimeZone] = useState(initial);
  return <TimeZonePicker value={timeZone} onChange={setTimeZone} />;
}

export const Default: Story = {
  args: { value: "Asia/Tokyo", onChange: () => {} },
  render: () => <InteractiveTimeZonePicker initial="Asia/Tokyo" />,
  play: async ({ canvasElement }) => {
    const input = within(canvasElement).getByRole("combobox");
    await expect(input).toHaveValue("Asia/Tokyo");
    await expect(input.className).toContain("font-mono");
    await expect(input.className).toContain("border-border");
    await expect(input.className).not.toMatch(/bg-(white|gray|slate)|text-\[#|dark:/);
  },
};

export const SearchAndSelect: Story = {
  args: { value: "Asia/Tokyo", onChange: () => {} },
  render: () => <InteractiveTimeZonePicker initial="Asia/Tokyo" />,
  play: async ({ canvasElement }) => {
    const canvas = within(canvasElement);
    const input = canvas.getByRole("combobox");
    await userEvent.click(input);
    await userEvent.type(input, "London");
    const option = await canvas.findByRole("option", { name: "Europe/London" });
    await userEvent.click(option);
    await expect(input).toHaveValue("Europe/London");
  },
};
