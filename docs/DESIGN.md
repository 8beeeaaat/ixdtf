# IXDTF Demo — UI デザインシステム (DESIGN.md)

タイポグラフィ・ミニマルを視覚言語とする、このプロジェクト専用のデザインシステム。
グローバルの Liquid Glass + shadcn/ui 規約から**工学的規約 (トークン設計・コンポーネントパターン・
ダークモード方式・禁止パターン・Storybook 規則) を継承**し、ガラス表現は採用しない。

要件との対応: ビジュアル方針は [architecture.md](./architecture.md)「スタイリング」節、
テーマ切替は requirements.md F-0-3、a11y は N-3。

## デザインコンセプト

**IXDTF 文字列そのものが主役。UI は書体と余白で語る。**

1. **白/黒基調** — 彩度のある色は IXDTF 構成要素のハイライトと状態表示 (成功/エラー/critical) にのみ使う
2. **等幅フォント中心** — 時計・IXDTF 文字列・コード類はすべて `font-mono` + `tabular-nums`
3. **大きな余白** — 詰めて飾るより、空けて絞る。1 画面 1 主役
4. **面より線** — 影やぼかしではなくヘアラインボーダー (`border-border`) で領域を区切る

## Stack

- Tailwind CSS v4 — `web/src/index.css` に `@import "tailwindcss"` + `@theme inline`
- shadcn/ui パターン (CLI 不使用) — `src/components/ui/` に自作コンポーネント
- `cn()` = tailwind-merge + clsx (`@/lib/utils`) — `clsx` 直接 import は禁止
- `cva` (class-variance-authority) — バリアントを持つコンポーネントに使用
- Radix UI primitives — **挙動 (a11y) のみ**利用可 (Tooltip / Select / Switch)。スタイルは常に自作
- 色空間は OKLCH で統一

## カラートークン

### セマンティックトークン (shadcn 命名)

`index.css` で CSS 変数として定義し、`@theme inline` で Tailwind ユーティリティに接続する。
**className に生の色 (`text-gray-500` 等) を書くことは禁止。**

```css
:root {
  --background: oklch(0.99 0 0);          /* ほぼ白 */
  --foreground: oklch(0.16 0 0);          /* ほぼ黒 */
  --card: oklch(1 0 0);
  --card-foreground: var(--foreground);
  --primary: oklch(0.2 0 0);              /* モノクロの主ボタン */
  --primary-foreground: oklch(0.98 0 0);
  --secondary: oklch(0.96 0 0);
  --secondary-foreground: oklch(0.3 0 0);
  --muted: oklch(0.96 0 0);
  --muted-foreground: oklch(0.48 0 0);
  --accent: oklch(0.95 0 0);              /* ホバー面 */
  --accent-foreground: var(--foreground);
  --border: oklch(0.9 0 0);               /* ヘアライン */
  --ring: oklch(0.6 0 0);
  --destructive: oklch(0.55 0.19 25);     /* エラー / 不一致 */
  --success: oklch(0.55 0.14 150);        /* 成功 / 一致 */
  --warning: oklch(0.6 0.14 85);          /* critical フラグ注意 */
  --info: oklch(0.55 0.12 240);
}
```

ダーク値は同名変数を `prefers-color-scheme: dark` (`:root:not(.light)`) と `.dark` クラスで上書きする
(下記「ダークモード」)。ダークは「黒地に白」の反転 + 各色の明度を 0.7 前後へ引き上げる。

### トークン早見表

className にはセマンティックトークンのみ使う:

| 用途 | Foreground | Background | Border |
|---|---|---|---|
| ページ | `text-foreground` | `bg-background` | — |
| カード | `text-card-foreground` | `bg-card` | `border-border` |
| 主アクション | `text-primary-foreground` | `bg-primary` | — |
| 副アクション | `text-secondary-foreground` | `bg-secondary` | `border-border` |
| ラベル / 補足 | `text-muted-foreground` | `bg-muted` | — |
| ホバー面 | `text-accent-foreground` | `bg-accent` | — |
| エラー / 不一致 | `text-destructive` | `bg-destructive/10` | — |
| 成功 / 一致 | `text-success` | `bg-success/10` | — |
| critical 注意 | `text-warning` | `bg-warning/10` | — |
| 情報 | `text-info` | `bg-info/10` | — |

### IXDTF ハイライトトークン (プロジェクト固有)

`lib/ixdtf/tokenize.ts` のトークン種別と **1:1 対応**。低彩度 (chroma ≈ 0.09) で「控えめな色分け」を保つ。

```css
:root {
  --ixdtf-date:      oklch(0.45 0.09 250);  /* 2026-07-07 */
  --ixdtf-time:      oklch(0.45 0.09 160);  /* 23:45:12 */
  --ixdtf-offset:    oklch(0.5 0.09 60);    /* +09:00 */
  --ixdtf-timezone:  oklch(0.45 0.09 320);  /* [Asia/Tokyo] */
  --ixdtf-extension: oklch(0.5 0.09 20);    /* [u-ca=japanese] */
}
```

| tokenize 種別 | クラス | 色以外の手掛かり (N-3) |
|---|---|---|
| `date` | `text-ixdtf-date` | ホバー/フォーカスでツールチップ + 凡例ラベル |
| `time` | `text-ixdtf-time` | 同上 |
| `offset` | `text-ixdtf-offset` | 同上 |
| `timezone` | `text-ixdtf-timezone` | 同上 |
| `extension` | `text-ixdtf-extension` | 同上 |
| critical 修飾 (`!`) | 上記色 + `underline decoration-wavy decoration-warning` | 波下線 + ツールチップ |

ハイライト表示には必ず凡例 (トークン名ラベル) を併置し、色だけに意味を持たせない。

## タイポグラフィ

```css
@theme inline {
  --font-sans: system-ui, -apple-system, "Segoe UI", sans-serif;
  --font-mono: "IBM Plex Mono", ui-monospace, "SF Mono", "Cascadia Code", monospace;
}
```

- `IBM Plex Mono` は @fontsource でセルフホスト (weight 400 / 500 のみ)。読み込み失敗時は
  `ui-monospace` へフォールバック
- **IXDTF 文字列・時計・コードは常に `font-mono tabular-nums`** — `tabular-nums` は
  1 秒 tick での桁幅変動によるレイアウトシフト防止 (N-4) に必須
- 本文・ナビ・ラベルは `font-sans`

| 役割 | スタイル |
|---|---|
| 時計 (F-1) | `font-mono text-6xl md:text-8xl font-normal tabular-nums tracking-tight` |
| IXDTF 文字列 (主役表示) | `font-mono text-lg md:text-2xl tabular-nums break-all` |
| ページタイトル | `font-sans text-2xl font-semibold tracking-tight` |
| セクション見出し | `font-sans text-sm font-semibold uppercase tracking-wider text-muted-foreground` |
| 本文 | `font-sans text-sm text-foreground` |
| 補足 / ヒント | `font-sans text-xs text-muted-foreground` |

## サーフェス (glass の置き換え)

**`glass` / `glass-dense` / `glass-subtle` はこのプロジェクトでは使用しない** (タイポグラフィ・
ミニマル方針)。面の階層はボーダーと背景明度のみで表現する:

| 面 | スタイル | 用途 |
|---|---|---|
| ページ | `bg-background` | 全画面。グラデーション背景なし |
| カード | `rounded-lg border border-border bg-card` | `<Card>` コンポーネント経由でのみ使用 |
| 強調行 | `bg-muted` | テーブルの差分行、コード背景 |
| 区切り | `border-t border-border` | セクション区切り (影は使わない) |

- 角丸は `--radius: 0.375rem` (控えめ)。影 (`shadow-*`) は原則使用しない
- フォーカスは `focus-visible:ring-2 ring-ring ring-offset-2 ring-offset-background`

## レイアウト

- 単一カラム: `mx-auto max-w-3xl px-6` (Interop の比較テーブルのみ `max-w-5xl`)
- ヘッダーナビ: 上部固定 `border-b border-border bg-background`。5 ページのテキストリンク
  (アイコン過多にしない)。モバイルは横スクロール可能なタブ列
- セクション間隔: `space-y-12` 以上。「詰まった」印象を避ける

## コンポーネントパターン

### 新規コンポーネントのチェックリスト

1. `cn` を `@/lib/utils` から import (`clsx` 直接 import 禁止)
2. `className` prop を受け取り `cn()` でマージ
3. `class` ではなく関数宣言。named export
4. プロジェクト内 import は `@/` エイリアス
5. 汎用プリミティブは `src/components/ui/`、ドメイン部品は feature フォルダへ

### バリアントコンポーネント (CVA パターン)

例: 拡張タグ表示用 `TagBadge` (F-2-3 / F-3-2 で使用):

```tsx
import { cn } from "@/lib/utils";
import { cva, type VariantProps } from "class-variance-authority";
import type { HTMLAttributes } from "react";

const tagBadgeVariants = cva(
  "rounded-full px-2 py-0.5 font-mono text-xs font-medium",
  {
    variants: {
      variant: {
        default: "bg-muted text-muted-foreground",
        critical: "bg-warning/10 text-warning underline decoration-wavy",
        match: "bg-success/10 text-success",
        mismatch: "bg-destructive/10 text-destructive",
      },
    },
    defaultVariants: { variant: "default" },
  },
);

type TagBadgeProps = HTMLAttributes<HTMLSpanElement> &
  VariantProps<typeof tagBadgeVariants>;

function TagBadge({ className, variant, ...props }: TagBadgeProps) {
  return <span className={cn(tagBadgeVariants({ variant, className }))} {...props} />;
}

export { TagBadge, tagBadgeVariants };
```

コンポーネントと variants 関数の両方を export する。

### このプロジェクトで作る ui/ プリミティブ

| コンポーネント | バリアント / 備考 | 主な使用箇所 |
|---|---|---|
| `Button` | primary / secondary / ghost (glass バリアントなし) | 全画面 |
| `Card` | flat 固定 (`border border-border bg-card`) + CardHeader/Title/Content | F-1, F-2, F-3 |
| `Input` | 透明背景、focus ring | F-2, F-3, F-4 |
| `Combobox` | Radix ベース、検索可能 | TimeZonePicker / CalendarPicker |
| `ToggleSwitch` | Radix Switch ベース、`enabled`/`onChange` | StrictToggle, テーマ切替 |
| `TagBadge` | 上記 CVA 例 | 拡張タグ、比較結果 |
| `Tooltip` | Radix ベース | IXDTF ハイライトの説明 (F-1-2) |
| `FormField` | ラベル + エラー/ヒント (TanStack Form と接続) | F-2, F-3 |

ドメイン部品 (`IxdtfHighlight`, `CopyButton`, `ComparisonTable` 等) は ui/ プリミティブの合成で作る。

## ダークモード

- 自動: `@media (prefers-color-scheme: dark)` 内で `:root:not(.light)` に対しダーク値を定義
- 手動: `<html>` に `.dark` / `.light` クラス (3 状態: auto / light / dark を localStorage に保存)
- 両方式が同時に成立し、トークンが自動で切り替わる
- **`dark:` Tailwind バリアントは使わない** — 切替はすべてセマンティックトークン側で行う

## 禁止パターン

- `clsx` の直接 import — 常に `cn` (`@/lib/utils`)
- className への生の色値 (`text-gray-500`, `bg-white`, `text-[#333]` 等)
- **`glass` / `glass-dense` / `glass-subtle`** — このプロジェクトでは不使用
- **`dark:` バリアント** — トークン自動切替に一本化
- IXDTF ハイライト色の直書き — 必ず `text-ixdtf-*` トークンを使う
- `any` / `unknown` 型
- `class` キーワードによるコンポーネント定義
- `<Card>` を通さないカード風 div (`rounded-lg border ...` の手書き散在)
- `shadow-*` の多用 (面の階層はボーダーで表現する)
- ハードコードされた表示文字列 — 常に `t()` (react-i18next)

## Storybook 規則

`src/components/ui/` の全コンポーネントに `.stories.tsx` を必須とする:

- 全バリアントの visual story
- play function による interaction test — className がセマンティック/ixdtf トークンを
  使っていることを assert (生の色クラスが混入したら fail)
- preview に `index.css` を読み込み、ライト/ダーク両テーマの toolbar 切替を用意する
