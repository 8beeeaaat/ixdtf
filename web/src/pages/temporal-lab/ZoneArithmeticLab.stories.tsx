import type { Meta, StoryObj } from "@storybook/react-vite";
import { I18nextProvider } from "react-i18next";
import { expect, within } from "storybook/test";
import i18n from "@/app/i18n";
import { TooltipProvider } from "@/components/ui/Tooltip";
import { requiredSample } from "@/lib/fixtures";
import { getTemporal } from "@/lib/temporal/detect";
import { ZoneArithmeticLab } from "@/pages/temporal-lab/ZoneArithmeticLab";

const baseSample = requiredSample("dst-overlap-earlier");

const meta = {
  component: ZoneArithmeticLab,
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
    input: baseSample.input,
    native: {
      base: {
        formatted: baseSample.input,
        epochNanoseconds: "1541320200000000000",
        error: null,
      },
      plusDay: {
        formatted: "2018-11-05T01:30:00-08:00[America/Los_Angeles]",
        epochNanoseconds: "1541410200000000000",
        error: null,
      },
      plusHours: {
        formatted: "2018-11-05T00:30:00-08:00[America/Los_Angeles]",
        epochNanoseconds: "1541406600000000000",
        error: null,
      },
    },
    polyfill: {
      base: {
        formatted: baseSample.input,
        epochNanoseconds: "1541320200000000000",
        error: null,
      },
      plusDay: {
        formatted: "2018-11-05T01:30:00-08:00[America/Los_Angeles]",
        epochNanoseconds: "1541410200000000000",
        error: null,
      },
      plusHours: {
        formatted: "2018-11-05T00:30:00-08:00[America/Los_Angeles]",
        epochNanoseconds: "1541406600000000000",
        error: null,
      },
    },
    plusDayGo: { unixNano: "1541410200000000000", error: null, pending: false },
    plusHoursGo: { unixNano: "1541406600000000000", error: null, pending: false },
  },
} satisfies Meta<typeof ZoneArithmeticLab>;

export default meta;
type Story = StoryObj<typeof meta>;

export const CalendarVersusExactTime: Story = {
  play: async ({ args, canvasElement }) => {
    const canvas = within(canvasElement);
    const temporal = getTemporal();
    await expect(temporal).not.toBeNull();
    if (!temporal) {
      throw new Error("Native Temporal is required by this story");
    }
    // args の期待値がネイティブ Temporal の実測と一致することを検証する
    const zdt = temporal.ZonedDateTime.from(args.input);
    await expect(zdt.add({ days: 1 }).toString()).toBe(args.polyfill.plusDay.formatted);
    await expect(zdt.add({ days: 1 }).epochNanoseconds.toString()).toBe(
      args.polyfill.plusDay.epochNanoseconds,
    );
    await expect(zdt.add({ hours: 24 }).toString()).toBe(args.polyfill.plusHours.formatted);
    await expect(zdt.add({ hours: 24 }).epochNanoseconds.toString()).toBe(
      args.polyfill.plusHours.epochNanoseconds,
    );

    await expect(canvas.getByText("+25 h")).toBeVisible();
    await expect(canvas.getByText("+24 h")).toBeVisible();
    // F-6-5: Native / Polyfill の epoch と Go の unix_nano が同値で 3 箇所ずつ現れる
    await expect(canvas.getAllByText(args.polyfill.plusDay.epochNanoseconds ?? "")).toHaveLength(3);
    await expect(canvas.getAllByText(args.polyfill.plusHours.epochNanoseconds ?? "")).toHaveLength(
      3,
    );
    // ネイティブ / polyfill のブロックが各結果に 1 つずつ (計 2) 併記される
    await expect(canvas.getAllByText(i18n.t("temporal.implNative"))).toHaveLength(2);
    await expect(canvas.getAllByText(i18n.t("temporal.implPolyfill"))).toHaveLength(2);
    await expect(canvas.getAllByText(i18n.t("temporalLab.goBadge"))).toHaveLength(2);
    await expect(canvas.queryByText(i18n.t("common.match"))).not.toBeInTheDocument();
    await expect(canvas.getByText(i18n.t("reproduce.summary"))).toBeVisible();
    await expect(canvasElement.innerHTML).not.toMatch(
      /(?:text|bg)-(?:gray|slate|zinc|red|white)|text-\[|bg-\[|dark:|glass|shadow-/,
    );
  },
};

export const ServerUnavailable: Story = {
  args: {
    plusDayGo: { unixNano: null, error: "fetch failed", pending: false },
    plusHoursGo: { unixNano: null, error: null, pending: true },
  },
  play: async ({ canvasElement }) => {
    const canvas = within(canvasElement);
    await expect(canvas.getByText(i18n.t("temporalLab.goUnavailable"))).toBeVisible();
    await expect(canvas.getByText(i18n.t("temporalLab.goPending"))).toBeVisible();
    await expect(canvas.getByText("+25 h")).toBeVisible();
  },
};

export const NativeUnsupported: Story = {
  args: {
    native: null,
  },
  play: async ({ canvasElement }) => {
    const canvas = within(canvasElement);
    // ネイティブブロックは案内メッセージに置き換わる (各結果に 1 つ、計 2)
    await expect(canvas.getAllByText(i18n.t("temporalLab.nativeUnsupported"))).toHaveLength(2);
    // polyfill ブロックと経過時間は通常どおり表示され続ける
    await expect(canvas.getAllByText(i18n.t("temporal.implPolyfill"))).toHaveLength(2);
    await expect(canvas.getByText("+25 h")).toBeVisible();
  },
};
