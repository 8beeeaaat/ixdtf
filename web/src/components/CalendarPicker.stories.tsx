import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { I18nextProvider } from "react-i18next";
import { expect, userEvent, within } from "storybook/test";
import i18n from "@/app/i18n";
import { CalendarPicker } from "@/components/CalendarPicker";

const meta = {
  component: CalendarPicker,
  decorators: [
    (Story) => (
      <I18nextProvider i18n={i18n}>
        <div className="w-80">
          <Story />
        </div>
      </I18nextProvider>
    ),
  ],
} satisfies Meta<typeof CalendarPicker>;

export default meta;
type Story = StoryObj<typeof meta>;

function InteractiveCalendarPicker({ initial }: { initial: string }) {
  const [calendar, setCalendar] = useState(initial);
  return <CalendarPicker value={calendar} onChange={setCalendar} />;
}

export const Default: Story = {
  args: { value: "iso8601", onChange: () => {} },
  render: () => <InteractiveCalendarPicker initial="iso8601" />,
  play: async ({ canvasElement }) => {
    const input = within(canvasElement).getByRole("combobox");
    await expect(input).toHaveValue("iso8601");
    await expect(input.className).toContain("font-mono");
    await expect(input.className).toContain("border-border");
    await expect(input.className).not.toMatch(/bg-(white|gray|slate)|text-\[#|dark:/);
  },
};

export const SearchAndSelect: Story = {
  args: { value: "iso8601", onChange: () => {} },
  render: () => <InteractiveCalendarPicker initial="iso8601" />,
  play: async ({ canvasElement }) => {
    const canvas = within(canvasElement);
    const input = canvas.getByRole("combobox");
    await userEvent.click(input);
    await userEvent.type(input, "japanese");
    const option = await canvas.findByRole("option", { name: "japanese" });
    await userEvent.click(option);
    await expect(input).toHaveValue("japanese");
  },
};
