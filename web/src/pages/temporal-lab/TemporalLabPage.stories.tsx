import type { Meta, StoryObj } from "@storybook/react-vite";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { I18nextProvider } from "react-i18next";
import { expect, within } from "storybook/test";
import i18n from "@/app/i18n";
import { TooltipProvider } from "@/components/ui/Tooltip";
import { requiredSample } from "@/lib/fixtures";
import { getTemporal } from "@/lib/temporal/detect";
import { TemporalLabPage } from "@/pages/temporal-lab/TemporalLabPage";

const earlierSample = requiredSample("dst-overlap-earlier");
const laterSample = requiredSample("dst-overlap-later");

// Storybook では Go サーバーが居ないことがある。Go 併記セルは pending →
// 取得不可のフォールバックに落ち、Temporal 側の実測はそのまま表示される。
const queryClient = new QueryClient({
  defaultOptions: { queries: { retry: false } },
});

const meta = {
  component: TemporalLabPage,
  decorators: [
    (Story) => (
      <I18nextProvider i18n={i18n}>
        <QueryClientProvider client={queryClient}>
          <TooltipProvider>
            <Story />
          </TooltipProvider>
        </QueryClientProvider>
      </I18nextProvider>
    ),
  ],
} satisfies Meta<typeof TemporalLabPage>;

export default meta;
type Story = StoryObj<typeof meta>;

export const NativeTemporalPage: Story = {
  play: async ({ canvasElement }) => {
    const canvas = within(canvasElement);
    const temporal = getTemporal();
    await expect(temporal).not.toBeNull();
    if (!temporal) {
      throw new Error("Native Temporal is required by this story");
    }
    const earlierEpoch = temporal.ZonedDateTime.from(
      earlierSample.input,
    ).epochNanoseconds.toString();
    const laterEpoch = temporal.ZonedDateTime.from(laterSample.input).epochNanoseconds.toString();

    await expect(canvas.getByRole("heading", { name: i18n.t("temporalLab.title") })).toBeVisible();
    // ヘッダーと各実験に Native / Polyfill / Go の 3 実装が併記される
    await expect(canvas.getAllByText(i18n.t("temporal.implNative")).length).toBeGreaterThan(0);
    await expect(canvas.getAllByText(i18n.t("temporal.implPolyfill")).length).toBeGreaterThan(0);
    // Native / Polyfill / Go 併記で同じ値が複数箇所に現れるため getAllByText で検証する
    await expect(canvas.getAllByText(earlierEpoch, { exact: true }).length).toBeGreaterThan(0);
    await expect(canvas.getAllByText(laterEpoch, { exact: true }).length).toBeGreaterThan(0);
    await expect(canvas.getByText("+1 h", { exact: true })).toBeVisible();

    // F-6-3 / F-6-4 の実験が並ぶ
    await expect(canvas.getByText(i18n.t("temporalLab.zoneArithmetic.title"))).toBeVisible();
    await expect(canvas.getByText("+25 h", { exact: true })).toBeVisible();
    await expect(canvas.getByText("+24 h", { exact: true })).toBeVisible();
    await expect(canvas.getByText(i18n.t("temporalLab.calendarProjection.title"))).toBeVisible();
    await expect(canvas.getByText("japanese", { exact: true })).toBeVisible();

    // F-6-5: Go ixdtf の併記ラベルはあるが、match/mismatch の採点 UI は無い (D-12)
    await expect(canvas.getAllByText(i18n.t("temporalLab.goBadge")).length).toBeGreaterThan(0);
    await expect(canvas.queryByText(i18n.t("common.match"))).not.toBeInTheDocument();
    await expect(canvas.queryByText(i18n.t("common.mismatch"))).not.toBeInTheDocument();
    // F-0-7: 再現セクションが各実験に付く
    await expect(canvas.getAllByText(i18n.t("reproduce.summary"))).toHaveLength(3);
  },
};
