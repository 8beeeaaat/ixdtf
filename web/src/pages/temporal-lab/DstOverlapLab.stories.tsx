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
      native: { epochNanoseconds: "1541320200000000000", error: null },
      polyfill: { epochNanoseconds: "1541320200000000000", error: null },
    },
    later: {
      id: "later",
      input: laterSample.input,
      native: { epochNanoseconds: "1541323800000000000", error: null },
      polyfill: { epochNanoseconds: "1541323800000000000", error: null },
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
    // args の期待値がネイティブ Temporal の実測 (= polyfill と同値) に一致する
    await expect(temporal.ZonedDateTime.from(args.earlier.input).epochNanoseconds.toString()).toBe(
      args.earlier.polyfill.epochNanoseconds,
    );
    await expect(temporal.ZonedDateTime.from(args.later.input).epochNanoseconds.toString()).toBe(
      args.later.polyfill.epochNanoseconds,
    );
    // ネイティブ / polyfill のブロックが各観測に 1 つずつ (計 2) 併記される
    await expect(canvas.getAllByText(i18n.t("temporal.implNative"))).toHaveLength(2);
    await expect(canvas.getAllByText(i18n.t("temporal.implPolyfill"))).toHaveLength(2);
    // F-6-5: Go ixdtf の実測も併記され、同じ瞬間 (unix_nano) が Native+Polyfill+Go の
    // 3 箇所ずつ現れる
    await expect(canvas.getAllByText(i18n.t("temporalLab.goBadge"))).toHaveLength(2);
    await expect(canvas.getAllByText(args.earlier.polyfill.epochNanoseconds ?? "")).toHaveLength(3);
    await expect(canvas.getAllByText(args.later.polyfill.epochNanoseconds ?? "")).toHaveLength(3);
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
      native: { epochNanoseconds: null, error: "RangeError: offset is invalid for this time zone" },
      polyfill: {
        epochNanoseconds: null,
        error: "RangeError: offset is invalid for this time zone",
      },
    },
  },
  play: async ({ canvasElement }) => {
    const canvas = within(canvasElement);
    // later のネイティブ / polyfill 両ブロックが parseError を表示 (計 2)
    await expect(canvas.getAllByText(i18n.t("temporalLab.dstOverlap.parseError"))).toHaveLength(2);
    // polyfill の epoch が欠けるので経過時間は算出不能
    await expect(canvas.getByText(i18n.t("temporalLab.dstOverlap.unavailable"))).toBeVisible();
    const summaries = canvas.getAllByText(i18n.t("temporalLab.dstOverlap.errorDetails"));
    await expect(summaries).toHaveLength(2);
    await userEvent.click(summaries[0]);
    const detail = canvas.getAllByText("RangeError: offset is invalid for this time zone")[0];
    await expect(detail).toBeVisible();
    await expect(detail).toHaveClass("font-mono", "text-destructive");
  },
};

export const NativeUnsupported: Story = {
  args: {
    earlier: { ...meta.args.earlier, native: null },
    later: { ...meta.args.later, native: null },
  },
  play: async ({ canvasElement }) => {
    const canvas = within(canvasElement);
    // ネイティブブロックは案内メッセージに置き換わる (各観測に 1 つ、計 2)
    await expect(canvas.getAllByText(i18n.t("temporalLab.nativeUnsupported"))).toHaveLength(2);
    // polyfill ブロックは通常どおり epoch を表示し続ける
    await expect(canvas.getAllByText(i18n.t("temporal.implPolyfill"))).toHaveLength(2);
    await expect(canvas.getByText("+1 h")).toBeVisible();
  },
};
