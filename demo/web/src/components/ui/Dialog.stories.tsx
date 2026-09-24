import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { expect, fn, userEvent, waitFor, within } from "storybook/test";
import { Button } from "./Button";
import { Dialog } from "./Dialog";

function DialogExample({
  initialOpen = false,
  onOpenChange = fn(),
}: {
  initialOpen?: boolean;
  onOpenChange?: (open: boolean) => void;
}) {
  const [open, setOpen] = useState(initialOpen);
  const updateOpen = (nextOpen: boolean) => {
    setOpen(nextOpen);
    onOpenChange(nextOpen);
  };
  return (
    <>
      <Button onClick={() => updateOpen(true)}>Open dialog</Button>
      <Dialog open={open} onOpenChange={updateOpen} title="Reference" closeLabel="Close">
        <p className="font-sans text-sm">Dialog content</p>
      </Dialog>
    </>
  );
}

const meta = {
  component: DialogExample,
  args: { initialOpen: false },
} satisfies Meta<typeof DialogExample>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  play: async ({ canvasElement }) => {
    await userEvent.click(within(canvasElement).getByRole("button", { name: "Open dialog" }));
    const dialog = await within(document.body).findByRole("dialog");
    await expect(dialog).toHaveTextContent("Dialog content");
    await expect(dialog.className).toContain("bg-card");
    await expect(dialog.className).not.toMatch(/bg-(white|gray|slate)|text-\[#/);
    await userEvent.keyboard("{Escape}");
    await waitFor(() => expect(dialog).not.toBeVisible());
    await waitFor(() =>
      expect(within(canvasElement).getByRole("button", { name: "Open dialog" })).toHaveFocus(),
    );
  },
};

export const Open: Story = {
  args: { initialOpen: true },
};
