import type { Meta, StoryObj } from "@storybook/react-vite";
import { I18nextProvider } from "react-i18next";
import { expect, within } from "storybook/test";
import i18n from "@/app/i18n";
import { TooltipProvider } from "@/components/ui/Tooltip";
import { requiredSample } from "@/lib/fixtures";
import { getNativeTemporal, getPolyfillTemporal } from "@/lib/temporal/detect";
import { CalendarProjectionLab } from "@/pages/temporal-lab/CalendarProjectionLab";
import {
  observeCalendarProjections,
  PROJECTION_CALENDARS,
} from "@/pages/temporal-lab/observations";

const calendarSample = requiredSample("calendar-japanese");
// 暦データは ICU 依存のため、args はレンダリング環境の実測から組み立てる
const nativeTemporal = getNativeTemporal();
const nativeProjections = nativeTemporal
  ? observeCalendarProjections(calendarSample.input, PROJECTION_CALENDARS, nativeTemporal)
  : null;
const polyfillProjections = observeCalendarProjections(
  calendarSample.input,
  PROJECTION_CALENDARS,
  getPolyfillTemporal(),
);
// Go の unix_nano は polyfill の実測瞬間に合わせ、常に一致させる
const polyfillInstant =
  polyfillProjections.find((projection) => projection.epochNanoseconds !== null)
    ?.epochNanoseconds ?? null;

const meta = {
  component: CalendarProjectionLab,
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
    input: calendarSample.input,
    native: nativeProjections,
    polyfill: polyfillProjections,
    go: {
      unixNano: polyfillInstant,
      formatted: calendarSample.input,
      lossless: true,
      error: null,
      pending: false,
    },
  },
} satisfies Meta<typeof CalendarProjectionLab>;

export default meta;
type Story = StoryObj<typeof meta>;

export const OneInstantManyCalendars: Story = {
  play: async ({ args, canvasElement }) => {
    const canvas = within(canvasElement);
    const resolved = args.polyfill.filter((projection) => projection.error === null);
    await expect(resolved.length).toBeGreaterThan(0);
    // F-6-4: 全投影で瞬間 (epochNanoseconds) は同一で、Go の unix_nano とも一致する
    const instants = new Set(resolved.map((projection) => projection.epochNanoseconds));
    await expect(instants.size).toBe(1);
    await expect([...instants][0]).toBe(args.go.unixNano);

    for (const projection of args.polyfill) {
      await expect(canvas.getByText(projection.calendarId)).toBeVisible();
    }
    const japanese = args.polyfill.find(
      (projection) => projection.calendarId === "japanese" && projection.error === null,
    );
    if (japanese && japanese.era !== null) {
      await expect(canvas.getByText(`${japanese.era} ${japanese.eraYear}`)).toBeVisible();
    }
    // ネイティブ / polyfill の瞬間ブロックが 1 つずつ併記される
    await expect(canvas.getByText(i18n.t("temporal.implNative"))).toBeVisible();
    await expect(canvas.getByText(i18n.t("temporal.implPolyfill"))).toBeVisible();
    await expect(canvas.getByText(i18n.t("temporalLab.goBadge"))).toBeVisible();
    await expect(canvas.getByText(/lossless: true/)).toBeVisible();
    await expect(canvas.queryByText(i18n.t("common.match"))).not.toBeInTheDocument();
    await expect(canvas.getByText(i18n.t("reproduce.summary"))).toBeVisible();
    await expect(canvasElement.innerHTML).not.toMatch(
      /(?:text|bg)-(?:gray|slate|zinc|red|white)|text-\[|bg-\[|dark:|glass|shadow-/,
    );
  },
};

export const ServerPending: Story = {
  args: {
    go: { unixNano: null, formatted: null, lossless: null, error: null, pending: true },
  },
  play: async ({ canvasElement }) => {
    const canvas = within(canvasElement);
    await expect(canvas.getByText(i18n.t("temporalLab.goPending"))).toBeVisible();
  },
};

export const NativeUnsupported: Story = {
  args: {
    native: null,
  },
  play: async ({ canvasElement }) => {
    const canvas = within(canvasElement);
    // ネイティブの瞬間ブロックは案内メッセージに置き換わる (1 箇所)
    await expect(canvas.getByText(i18n.t("temporalLab.nativeUnsupported"))).toBeVisible();
    // polyfill の瞬間ブロックと投影テーブルは通常どおり表示され続ける
    await expect(canvas.getByText(i18n.t("temporal.implPolyfill"))).toBeVisible();
    await expect(canvas.getByText("japanese")).toBeVisible();
  },
};
