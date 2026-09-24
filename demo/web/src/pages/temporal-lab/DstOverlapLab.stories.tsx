import type { Meta, StoryObj } from "@storybook/react-vite";
import { I18nextProvider } from "react-i18next";
import { expect, userEvent, within } from "storybook/test";
import i18n from "@/app/i18n";
import { TooltipProvider } from "@/components/ui/Tooltip";
import { requiredSample } from "@/lib/fixtures";
import { getTemporal } from "@/lib/temporal/detect";
import { DstOverlapLab } from "@/pages/temporal-lab/DstOverlapLab";

const earlierSample = requiredSample("dst-overlap-earlier");
const laterSample = requiredSample("dst-overlap-later");

const meta = {
  component: DstOverlapLab,
  decorators: [
    (Story) => (
      <I18nextProvider i18n={i18n}>
        <TooltipProvider>
          <div className="mx-auto max-w-6xl px-6">
            <Story />
          </div>
        </TooltipProvider>
      </I18nextProvider>
    ),
  ],
  args: {
    earlier: {
      id: "earlier",
      input: earlierSample.input,
      epochNanoseconds: "1541320200000000000",
      error: null,
    },
    later: {
      id: "later",
      input: laterSample.input,
      epochNanoseconds: "1541323800000000000",
      error: null,
    },
    earlierGo: { unixNano: "1541320200000000000", error: null, pending: false },
    laterGo: { unixNano: "1541323800000000000", error: null, pending: false },
  },
} satisfies Meta<typeof DstOverlapLab>;

export default meta;
type Story = StoryObj<typeof meta>;

export const ResolvedOverlap: Story = {
  play: async ({ args, canvasElement }) => {
    const canvas = within(canvasElement);
    await expect(canvas.getByText("+1 h")).toBeVisible();
    const temporal = getTemporal();
    await expect(temporal).not.toBeNull();
    if (!temporal) {
      throw new Error("Native Temporal is required by this story");
    }
    await expect(temporal.ZonedDateTime.from(args.earlier.input).epochNanoseconds.toString()).toBe(
      args.earlier.epochNanoseconds,
    );
    await expect(temporal.ZonedDateTime.from(args.later.input).epochNanoseconds.toString()).toBe(
      args.later.epochNanoseconds,
    );
    await expect(canvas.getAllByText(i18n.t("temporalLab.dstOverlap.nativeTemporal"))).toHaveLength(
      2,
    );
    // F-6-5: Go ixdtf の実測が併記され、同じ瞬間 (unix_nano) が 2 箇所ずつ現れる
    await expect(canvas.getAllByText(i18n.t("temporalLab.goBadge"))).toHaveLength(2);
    await expect(canvas.getAllByText(args.earlier.epochNanoseconds ?? "")).toHaveLength(2);
    await expect(canvas.getAllByText(args.later.epochNanoseconds ?? "")).toHaveLength(2);
    // D-12: match/mismatch の採点 UI は置かない
    await expect(canvas.queryByText(i18n.t("common.match"))).not.toBeInTheDocument();
    await expect(canvas.getByText(i18n.t("reproduce.summary"))).toBeVisible();
    await expect(canvas.getByText("+1 h")).toHaveClass("text-ixdtf-offset");
    await expect(canvasElement.innerHTML).not.toMatch(
      /(?:text|bg)-(?:gray|slate|zinc|red|white)|text-\[|bg-\[|dark:|glass|shadow-/,
    );
  },
};

export const TemporalError: Story = {
  args: {
    later: {
      ...meta.args.later,
      epochNanoseconds: null,
      error: "RangeError: offset is invalid for this time zone",
    },
  },
  play: async ({ canvasElement }) => {
    const canvas = within(canvasElement);
    await expect(canvas.getByText(i18n.t("common.error"))).toBeVisible();
    await expect(canvas.getByText("—", { exact: true })).toBeVisible();
    await expect(canvas.getByText(i18n.t("temporalLab.dstOverlap.unavailable"))).toBeVisible();
    await userEvent.click(canvas.getByText(i18n.t("temporalLab.dstOverlap.errorDetails")));
    const detail = canvas.getByText("RangeError: offset is invalid for this time zone");
    await expect(detail).toBeVisible();
    await expect(detail).toHaveClass("font-mono", "text-destructive");
  },
};
