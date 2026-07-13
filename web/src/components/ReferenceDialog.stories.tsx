import type { Meta, StoryObj } from "@storybook/react-vite";
import { I18nextProvider } from "react-i18next";
import { expect, userEvent, waitFor, within } from "storybook/test";
import i18n from "@/app/i18n";
import { ReferenceDialog } from "@/components/ReferenceDialog";

const meta = {
  component: ReferenceDialog,
  args: { referenceId: "ixdtf" },
  decorators: [
    (Story) => (
      <I18nextProvider i18n={i18n}>
        <Story />
      </I18nextProvider>
    ),
  ],
} satisfies Meta<typeof ReferenceDialog>;

export default meta;
type Story = StoryObj<typeof meta>;

export const IxdtfReference: Story = {
  play: async ({ canvasElement }) => {
    const trigger = within(canvasElement).getByRole("button");
    trigger.focus();
    await userEvent.keyboard("{Enter}");
    const dialog = await within(document.body).findByRole("dialog");
    const links = within(dialog).getAllByRole("link");
    await expect(links.length).toBeGreaterThan(1);
    for (const link of links) {
      await expect(link).toHaveAttribute("target", "_blank");
      await expect(link).toHaveAttribute("rel", "noreferrer");
      await expect(link).toHaveAccessibleName(/new tab|新しいタブ/i);
    }
    await userEvent.keyboard("{Escape}");
    await waitFor(() => expect(dialog).not.toBeVisible());
    await waitFor(() => expect(trigger).toHaveFocus());
  },
};

export const OffsetConsistencyReference: Story = {
  args: { referenceId: "offsetConsistency" },
  play: async ({ canvasElement }) => {
    await userEvent.click(within(canvasElement).getByRole("button"));
    const dialog = await within(document.body).findByRole("dialog");
    await expect(within(dialog).getByText(/timezone\.go/)).toBeInTheDocument();
    await expect(within(dialog).getByText(/Code|コード/)).toBeInTheDocument();
  },
};

export const SuffixSyntaxReference: Story = { args: { referenceId: "suffixSyntax" } };
export const CriticalAnnotationsReference: Story = {
  args: { referenceId: "criticalAnnotations" },
};
export const RejectedExtensionsReference: Story = {
  args: { referenceId: "rejectedExtensions" },
};
export const GoParsingReference: Story = { args: { referenceId: "goParsing" } };
export const TemporalParsingReference: Story = {
  args: { referenceId: "temporalParsing" },
  play: async ({ canvasElement }) => {
    await userEvent.click(within(canvasElement).getByRole("button"));
    const dialog = await within(document.body).findByRole("dialog");
    const note = within(dialog).getByRole("note");
    await expect(note).toHaveTextContent(/browser|ブラウザ/i);
    await expect(note).toHaveTextContent(/version|バージョン/i);
    await expect(note.className).toContain("bg-warning/10");
  },
};
export const RoundtripReference: Story = { args: { referenceId: "roundtrip" } };
export const DstAmbiguityReference: Story = { args: { referenceId: "dstAmbiguity" } };
