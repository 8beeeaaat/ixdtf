# IXDTF デモアプリケーション 機能要件

[github.com/8beeeaaat/ixdtf](https://pkg.go.dev/github.com/8beeeaaat/ixdtf) (RFC 9557 の Go 実装) と
ブラウザネイティブの [TC39 Temporal API](https://tc39.es/proposal-temporal/) を組み合わせ、
IXDTF (Internet Extended Date/Time Format) 規格を体験できるデモアプリケーションの機能要件を定義する。

## 目的

- RFC 9557 (IXDTF) が RFC 3339 に何を加えた規格なのかを、触って理解できるようにする
- 同一の IXDTF 文字列を「ブラウザの Temporal」と「Go の ixdtf」の両実装で処理し、相互運用性を実演する
- ixdtf ライブラリの API (`Parse` / `Format` / `FormatNano` / `Validate`) のショーケースとする

## 対象ユーザー

- IXDTF / Temporal に興味のある開発者
- ixdtf ライブラリの利用を検討している Go 開発者

## 画面構成

SPA として以下の 5 画面をヘッダーナビゲーションで切り替える。

| ID | 画面 | 概要 |
|----|------|------|
| F-1 | Home (Clock) | 現在時刻を IXDTF 表記のモダンな時計・カレンダーで表示 |
| F-2 | Playground | IXDTF 文字列の Parse / Validate をリアルタイム可視化 |
| F-3 | Interop | ブラウザ Temporal と Go ixdtf のラウンドトリップ比較 |
| F-4 | Converter | タイムゾーン / カレンダー変換 (ワールドクロック) |
| F-5 | Guide | RFC 9557 学習ガイド (サンプル文字列集) |

## 機能要件

### F-0: アプリケーションシェル (共通)

- **F-0-1** ヘッダーナビゲーションで 5 画面を切り替えられる SPA であること
- **F-0-2** UI 言語は日本語 / 英語を切り替えられること (react-i18next、全キーは両 locale で同期)
- **F-0-3** ライト / ダークテーマに対応すること (`prefers-color-scheme` 追従 + 手動切替)
- **F-0-4** 起動時に `globalThis.Temporal` の存在を検出し、未対応ブラウザでは対応ブラウザ一覧
  (Chrome 144+ / Edge 144+ / Firefox 139+) を示す案内画面を表示すること (ポリフィルは使用しない)
- **F-0-5** レスポンシブ対応 (モバイル〜デスクトップ)

### F-1: Home — IXDTF Clock / Calendar

現在時刻を IXDTF 規格の表記そのものを主役にして見せる、アプリの顔となる画面。

- **F-1-1** 現在時刻を 1 秒間隔で更新する時計として表示すること
  - `Temporal.Now.zonedDateTimeISO(timeZone)` を使用
- **F-1-2** 時計の下に完全な IXDTF 文字列 (例: `2026-07-07T23:30:00+09:00[Asia/Tokyo][u-ca=japanese]`)
  を表示し、**構成要素ごとに色分けハイライト**すること
  - 日付 / 時刻 / UTC オフセット / `[タイムゾーン名]` / `[u-ca=カレンダー]` を区別
  - 各要素のホバー (タップ) で規格上の意味をツールチップ表示
- **F-1-3** タイムゾーンを切り替えられること (`Intl.supportedValuesOf("timeZone")` から検索選択)
- **F-1-4** カレンダー体系を切り替えられること (`Intl.supportedValuesOf("calendar")` から選択。
  `iso8601` / `gregory` / `japanese` / `hebrew` / `islamic` / `chinese` 等)
  - 切替は IXDTF 文字列の `[u-ca=...]` サフィックスに即時反映
- **F-1-5** 選択中のカレンダー体系で描画した月間カレンダーを表示すること
  (`Temporal.PlainYearMonth` / `daysInMonth` / `monthsInYear` を使用。和暦なら「令和8年」表記等)
- **F-1-6** IXDTF 文字列をワンクリックでクリップボードにコピーできること
- **F-1-7** サーバー現在時刻 (`GET /api/now`、Go 側 `ixdtf.FormatNano` で生成) を併記し、
  「ブラウザ生成」と「サーバー生成」の IXDTF 文字列を見比べられること

### F-2: IXDTF Playground

- **F-2-1** IXDTF 文字列を入力すると、デバウンス付きでリアルタイムに解析結果を表示すること
- **F-2-2** strict モードをトグルで切り替えられること
  (strict 時はオフセットとタイムゾーンの整合性を検証 — ixdtf の `strict` 引数に対応)
- **F-2-3** サーバー解析結果 (`POST /api/ixdtf/parse`) として以下を表示すること
  - 解析された時刻 (RFC 3339 表記)
  - タイムゾーン名と critical フラグ (`[!Asia/Tokyo]` の `!`)
  - 拡張タグの一覧 (キー / 値 / critical) をテーブル表示
  - エラー時はエラーメッセージ全文 (ixdtf のエラーは規格上の理由を含むため教材になる)
- **F-2-4** 同じ入力に対するブラウザ側 `Temporal.ZonedDateTime.from()` の結果 (成功値または例外) を併記すること
- **F-2-5** 入力文字列をトークン分解し、F-1-2 と同じ色分けルールでハイライト表示すること
- **F-2-6** 「検証のみ」モード (ixdtf の `Validate` 相当) を選択できること
- **F-2-7** URL クエリで入力値と strict 設定を共有できること (Guide からの遷移にも使用)

### F-3: Interop — ラウンドトリップ比較

同一入力に対する 2 実装の挙動差そのものをコンテンツにする画面。

- **F-3-1** 入力した IXDTF 文字列を次の 2 経路で Parse → Format 往復させ、結果を並列表示すること
  - ブラウザ経路: `Temporal.ZonedDateTime.from(input).toString()`
  - サーバー経路: `POST /api/ixdtf/roundtrip` (ixdtf の `Parse` → `FormatNano`)
- **F-3-2** 比較テーブルで以下の観点を行ごとに表示し、両者が一致するか差分ハイライトすること
  - 解析成否 / エラーメッセージ
  - 解析された瞬間 (epoch)、タイムゾーン、カレンダー
  - 往復後の文字列 (元の文字列と一致するか = ロスレス性)
  - 不明な拡張タグ・critical フラグの扱い
- **F-3-3** 挙動差の出るプリセットサンプルを用意し、ワンクリックで試せること
  (例: strict でのオフセット不整合、critical な未知タグ、`x-` プライベート拡張)

### F-4: Converter — タイムゾーン / カレンダー変換

- **F-4-1** 基準時刻 (デフォルト: 現在時刻) を複数のタイムゾーンで一覧表示するワールドクロックであること
  - 各行に IXDTF 文字列と人間可読表記 (`Intl.DateTimeFormat`) を併記
  - `Temporal.ZonedDateTime.withTimeZone()` を使用
- **F-4-2** タイムゾーン行を追加・削除でき、選択は localStorage に永続化されること
- **F-4-3** カレンダー体系を切り替えると全行の表記と `[u-ca=...]` が一括で変わること
  (`withCalendar()` を使用)
- **F-4-4** 基準時刻を任意の IXDTF 文字列で指定できること (Playground からの引き継ぎ含む)

### F-5: Guide — RFC 9557 学習ガイド

- **F-5-1** RFC 9557 の概要 (RFC 3339 との差分、サフィックス構文、critical フラグ、
  private / experimental 拡張の拒否ルール) を解説するコンテンツを表示すること
- **F-5-2** カテゴリ別のサンプル文字列集を提供すること
  - 基本 (RFC 3339 互換) / タイムゾーン付き / カレンダー付き / critical フラグ / 複合
  - **拒否されるケース**: `x-` プライベート拡張、`_` 実験的拡張、strict でのオフセット不整合
  - サンプルは RFC 9557 本文の例と代表的な落とし穴から **10〜15 個に厳選**する
    (網羅より読み切れる学習体験を優先)
- **F-5-3** 各サンプルに「Playground で試す」ボタンを設け、F-2-7 のクエリ共有で遷移すること
- **F-5-4** 各サンプルに期待される結果 (成功 / 特定エラー) の注記を付けること

## API 要件 (Go サーバー)

すべて JSON。詳細スキーマは [architecture.md](./architecture.md) の OpenAPI 定義を参照。

| エンドポイント | 用途 | 対応する ixdtf API |
|---|---|---|
| `POST /api/ixdtf/parse` | 解析 + 検証 | `Parse(s, strict)` / `Validate(s, strict)` |
| `POST /api/ixdtf/format` | 時刻 + 拡張 → IXDTF 文字列 | `FormatNano(t, ext)` |
| `POST /api/ixdtf/roundtrip` | 解析→整形の往復を 1 リクエストで | `Parse` → `FormatNano` |
| `GET /api/now` | サーバー現在時刻の IXDTF 表記 | `FormatNano(time.Now(), ext)` |

- **A-1** 入力文字列の解析失敗は HTTP 200 の正常応答 (`ok: false` + エラー詳細) として返すこと。
  解析失敗はこのアプリのドメイン上の正常系 (デモの見せ場) であり、例外ではない
- **A-2** リクエスト自体の不備 (JSON 不正、必須フィールド欠落) のみ 4xx を返すこと
- **A-3** エラーメッセージは ixdtf ライブラリの原文をそのまま返すこと (加工しない)

## 非機能要件

- **N-1 ブラウザ対応**: ポリフィル不使用。ネイティブ Temporal 実装を持つ
  Chrome 144+ / Edge 144+ / Firefox 139+ を対象とする。Safari は安定版が未対応 (2026-07 時点、
  Technology Preview のみ) のため対象外とし、F-0-4 の案内画面でフォールバックする
- **N-2 i18n**: 全ユーザー向け文字列は `t()` 経由。ja / en の locale ファイルはキー構造を完全同期
- **N-3 アクセシビリティ**: キーボード操作可能、カラーハイライトは色以外の手掛かり (下線・ラベル) を併用
- **N-4 パフォーマンス**: 時計更新 (1 秒間隔) で再レンダリングが時計コンポーネント内に閉じること
- **N-5 品質ゲート**:
  - web: Biome (lint + format)、Vitest (unit)、Storybook (コンポーネントカタログ + interaction test)
  - server: golangci-lint、table-driven test
  - OpenAPI スキーマと生成コード (Orval / oapi-codegen) の乖離を CI で検出すること
- **N-6 依存最小**: Go サーバーは ixdtf + 標準ライブラリを基本とする (DB なし、認証なし)

## スコープ外 (Non-goals)

- ユーザーアカウント、データ永続化 (localStorage を除く)
- デプロイ / ホスティング構成 (ローカル実行のみを想定)
- Temporal API 全機能の網羅 (IXDTF 文字列に関係する範囲に限定)
- Safari 安定版のサポート (ポリフィル不使用の帰結)
