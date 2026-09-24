import type { Meta, StoryObj } from "@storybook/react-vite";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { I18nextProvider } from "react-i18next";
import { expect, within } from "storybook/test";
import i18n from "@/app/i18n";
import { TooltipProvider } from "@/components/ui/Tooltip";
import { requiredSample } from "@/lib/fixtures";
import { getTemporal } from "@/lib/temporal/detect";
import { LabPanel } from "@/pages/temporal-lab/LabPanel";

const earlierSample = requiredSample("dst-overlap-earlier");
const laterSample = requiredSample("dst-overlap-later");

// Storybook では Go サーバーが居ないことがある。Go 併記セルは pending →
// 取得不可のフォールバックに落ち、Temporal 側の実測はそのまま表示される。
const queryClient = new QueryClient({
  defaultOptions: { queries: { retry: false } },
});

const meta = {
  component: LabPanel,
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
} satisfies Meta<typeof LabPanel>;

export default meta;
type Story = StoryObj<typeof meta>;

// F-6: ワークベンチ Lab モードの本体。画面シェル (h1・モード切替) は WorkbenchPage が
// 提供するため、この Panel 単体では tagline + 3 実装バッジ + 3 実験を検証する。
export const NativeTemporalLab: Story = {
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

    // F-6-2: 実測であることを明示する tagline を Panel が持つ (h1 は WorkbenchPage 側)
    await expect(canvas.getByText(i18n.t("temporalLab.tagline"))).toBeVisible();
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
