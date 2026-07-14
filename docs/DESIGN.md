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

- 単一カラム: `mx-auto max-w-6xl px-6`
- ヘッダーナビ: 上部固定 `border-b border-border bg-background`。5 ページ (ホーム / はじめに / ガイド / ワークベンチ / 変換) のテキストリンク
  (アイコン過多にしない)。モバイルは横スクロール可能なタブ列
- セクション間隔: `space-y-12` 以上。「詰まった」印象を避ける

### Temporal ラボ モード — DST 重複時刻タイムライン

F-6-1 は「同じ壁時計時刻、異なる瞬間」を一目で理解できる横方向の時間軸として表現する。

- 左右に DST 終了前後の IXDTF 文字列を置き、中央に瞬間差 (`+1 h`) を大きな等幅数字で表示する
- 時間の進行方向はデスクトップで `→`、縦積みのモバイルで `↓` を併記し、配置だけに依存しない
- 両側ともローカル時刻とタイムゾーン名は同じにし、異なる UTC オフセットだけを
  `text-ixdtf-offset` で強調する。色だけに依存せず `DST` / `標準時` ラベルも併記する
- 各側にネイティブ Temporal が解決した `epochNanoseconds` を `font-mono tabular-nums` で表示する
- Go ixdtf との一致・不一致を示すバッジや Interop の詳細比較導線は置かない
- モバイルでは縦方向に積み、時間軸の順序と `+1 h` の関係を維持する

## ホーム・ヒーロー背景 (ドットマトリクス地球儀)

**ホーム (F-1) に限り**、選択中のタイムゾーンを直感的に見せる装飾背景として Three.js の
ドットマトリクス地球儀を敷く。これは「グラデ背景なし」原則の**意図的な例外**であり、
以下の制約で「文字が主役 / タイポグラフィ・ミニマル」を崩さない範囲に閉じる。

### 原則

- **文字が主役、地球儀は大気**: 時計・IXDTF 文字列の可読性を一切損なわない。ホームだけ
  コンテンツを狭幅・左寄せにし、空けた右側で地球儀をフィーチャーする。ホーム以外には出さない
- **インタラクティブ (ホーム、wide のみ)**: 地球儀はドラッグで回転、クリックで最寄り代表都市の
  タイムゾーンを選択できる。ヒット層はコンテンツ列 (z-10) の下 (z-[5]) に置き、コンテンツ外側を
  `pointer-events-none` にして右の空き領域だけ地球儀へ透過させる (カードは常に前面で操作可)。
  クリック→タイムゾーンも塗り分けと同じ境界ルックアップで解決する (ブラウザの `Intl` が
  解決できないゾーン ID のみ `nearestTimeZone` にフォールバック)。正当性判定には使わない
- **点は白黒基調**: 大陸・海のドットは無彩色トークン (`--muted-foreground` / `--foreground`) の
  明度差のみで表現する。地球儀自体に新しい彩度色を導入しない
- **彩度色は IXDTF トークンを流用**: 唯一の有彩色は「選択中タイムゾーンの地点マーカー」で、
  IXDTF ハイライトと同じ `--ixdtf-timezone` を使う (`[Asia/Tokyo]` の色 = 地球儀上のその地点の色)。
  UTC オフセットの弧 / ラベルは `--ixdtf-offset` を使う。これにより「彩度色は IXDTF 構成要素と
  状態表示にのみ」という規範を保つ
- **色はトークン駆動**: canvas は CSS クラスを使えないため、`getComputedStyle` で CSS 変数
  (`--muted-foreground` 等) を読み取り、テーマ切替 (F-0-3) に追従する。`dark:` は使わない

### 表現するもの (選択タイムゾーン)

1. 選択タイムゾーンの経度が正面に来るよう地球儀を回転し、その後はゆっくり自動回転する
   (実際の自転と同方向、約 120 秒 / 周)。ドラッグ中は停止し、`prefers-reduced-motion`
   では自動回転せず選択タイムゾーンを正面に向けた静止画のまま
2. 選択タイムゾーンの**陸地ドット**を `--ixdtf-timezone` でハイライトする (海のドットは
   通常色のまま)。塗り分けの「緯度経度 → IANA ゾーン」は
   [timezone-boundary-builder](https://github.com/evansiroky/timezone-boundary-builder)
   (OSM 由来の境界を IANA ゾーン ID に紐づけたデファクト標準データセット。IANA tzdb 自体は
   境界ポリゴンを持たない) を圧縮した `@photostructure/tz-lookup` で引くため、実際の
   タイムゾーン境界にほぼ一致する。あくまで表示専用の線引き (D-10) であり正当性判定には使わない
3. 代表都市座標に発光マーカー + パルス (`--ixdtf-timezone`)
4. 現在時刻に基づく昼夜テルミネータ (夜側のドットを減光)
5. UTC オフセットの弧 / ラベル (`--ixdtf-offset`)
6. 全代表都市を無彩色 (`--foreground`) の小マーカーとして常時表示し、ホバー時にその都市の
   タイムゾーン名 + UTC オフセットを選択中ラベルと同一スタイル (`text-ixdtf-timezone` /
   `text-ixdtf-offset` の 2 行) でラベル表示する (マーカー自体は選択中マーカーとの視覚的な
   区別のため彩度色にしない)

### アクセシビリティ / 性能 (N-3, N-4)

- **`prefers-reduced-motion` を尊重**: 自動回転とパルスを止め、選択タイムゾーンを正面に向けた
  静止画として描画する
- タブ非表示時 (`visibilitychange`) は描画ループを止める
- WebGL 非対応・初期化失敗時は**何も描画しない** (背景なしにフォールバック、機能は成立)
- Three.js は**遅延ロード** (dynamic import) し、ホームルート以外の初期バンドルを太らせない
- 座標マッピング (IANA タイムゾーン → 緯度経度) は主要都市の手書き表 → IANA tzdb `zone.tab`
  の代表座標 → 地域/オフセット近似の順で解決する。`tokenize.ts` と同じ**表示専用・近似**の
  線引きであり、正当性判定には一切使わない

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
| `Dialog` | HTML `dialog` ベース、モーダル・Escape・フォーカス復帰 | 公式参照ダイアログ (F-0-6) |

ドメイン部品 (`IxdtfHighlight`, `CopyButton`, `ComparisonTable` 等) は ui/ プリミティブの合成で作る。

### 公式参照ダイアログ

専門用語の横には小さな ghost ボタン「解説・参照」を置く。ボタンは本文や見出しより目立たせず、
`text-info` を参照導線にだけ用いる。開いたダイアログは以下の構造に統一する:

1. タイトル + 閉じるボタン
2. 1〜2 段落の短い解説
3. Temporal の実装差が関係する場合のみ、`border-warning` + `bg-warning/10` の注意 note
4. ヘアラインで区切った参照一覧。各行に種別ラベル (`仕様` / `ドキュメント` / `コード`)、
   リンク名、新規タブを示す記号を表示する

ダイアログは `max-w-2xl`、`bg-card`、`border-border` で、影やガラス表現を使わない。背景は
`bg-background/80` の backdrop で抑える。本文中に長い仕様解説を重複させず、ページの操作フローを
保ったまま一次情報へ段階的に掘り下げられる密度とする。

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
