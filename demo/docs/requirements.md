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

SPA として以下の 5 画面をヘッダーナビゲーションで切り替える。ワークベンチ (F-2) は「解析・検証」「往復・実装差比較」「Temporal ラボ」の 3 モードを内包する。

| ID | 画面 | 概要 |
|----|------|------|
| F-1 | Home (Clock) | 現在時刻を IXDTF 表記のモダンな時計・カレンダーで表示 |
| F-2 | IXDTF ワークベンチ | 解析・検証 (F-2) / 往復・実装差比較 (F-3) / Temporal ラボ (F-6) の 3 モード。旧 Playground + Interop + Temporal Lab を統合 |
| F-4 | Converter | タイムゾーン / カレンダー変換 (ワールドクロック) |
| F-5 | Guide | RFC 9557 学習ガイド (サンプル文字列集) |
| F-7 | About | IXDTF (規格) と Temporal (API) それぞれの解説と、両者の関係性の説明 |

## 機能要件

### F-0: アプリケーションシェル (共通)

- **F-0-1** ヘッダーで 5 画面を切り替えられる SPA であること (ホームへはロゴ兼ホームリンクで遷移し、残り 4 画面 IXDTF と Temporal / 変換 / ワークベンチ / ガイド はテキストナビで遷移。旧 Playground+Interop+Temporal Lab はワークベンチの 3 モードに統合)
  - md 未満 (モバイル) では、テキストナビの代わりに Home を含む 5 項目の固定ボトムタブバーで遷移する (F-0-5 参照)
- **F-0-2** UI 言語は日本語 / 英語を切り替えられること (react-i18next、全キーは両 locale で同期)
- **F-0-3** ライト / ダークテーマに対応すること (`prefers-color-scheme` 追従 + 手動切替)
- **F-0-4** ブラウザの TC39 Temporal 実装をユーザー選択なしで扱うこと (実行時セレクタは撤去)
  - `globalThis.Temporal` を検出し、単一 live-UI (Home / Converter) はネイティブがあればネイティブ、
    無ければ temporal-polyfill に自動フォールバックする (`web/src/lib/temporal/detect.ts`)
  - 比較系 (ワークベンチの往復・実装差比較モード / Temporal ラボ モード) はネイティブと polyfill を両方明示併記する
  - ネイティブ非対応ブラウザでは、該当カラム / カードにインラインで「ネイティブ Temporal なし」と
    示して縮退する (旧: 全画面の案内画面・localStorage 選択・ヘッダー切替コントロールは撤去済み)
- **F-0-5** レスポンシブ対応 (モバイル〜デスクトップ)
  - ヘッダーのテキストナビ (F-0-1) は md 未満で非表示にし、Home を含む 5 項目のアイコン+短ラベルの
    固定ボトムタブバーに切り替えること (DESIGN.md「モバイル ボトムタブバー」参照)
  - md 未満でページ全体が横スクロールしないこと (ヘッダー右クラスタや一覧行の長い文字列で
    行/ページが画面幅を超えないよう、折り返し・省略・アイコン化のいずれかで縮退させる)
- **F-0-6** ワークベンチ / Guide を中心に、専門用語・処理系固有の挙動には文脈に応じた
  「解説・参照」導線を表示できること
  - 導線から開くダイアログに、要点を短くまとめた解説と参照リンクを表示する
  - 参照リンクは一次情報 (RFC Editor、TC39 Temporal 公式仕様・ドキュメント、Go ixdtf の
    pkg.go.dev・バージョン固定ソースコード) に限定し、「仕様」「ドキュメント」「コード」を区別する
  - 外部リンクは新しいタブで開き、リンク先であることを表示上・支援技術の双方から判別できること
  - ダイアログはキーボードで開閉でき、Escape で閉じた後は起点へフォーカスが戻ること
  - 解説とリンク名は日本語 / 英語の locale を同期すること (URL は共通カタログで一元管理)
  - 初期配置はワークベンチ / Guide とし、同じ参照 ID を使って他画面へ拡張できること
  - Temporal の挙動を説明する解説には、実際の結果はブラウザ内蔵実装に依存し、ブラウザや
    バージョンによって未実装・実装差・仕様との差が生じ得ることを明記すること
- **F-0-7** サーバー API を呼び出す画面 (F-1-7 / F-2 / F-6) では、結果表示の直下に折りたたみ式の
  「再現」セクションを設け、閲覧したエンジニアが同じ結果を手元で再現できるようにすること
  - 現在の入力値を埋め込んだコードサンプルを Go (ixdtf ライブラリ直接利用) と
    JavaScript (Temporal API) の 2 言語で表示する
  - 各コードサンプルはワンクリックでコピーできること (F-1-6 と同じ挙動)
  - 初期状態は閉じており、結果そのものの閲覧を妨げないこと

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

### F-2: IXDTF ワークベンチ — 解析・検証 / 往復比較 / Temporal ラボ (旧 Playground + Interop + Temporal Lab 統合)

同一の IXDTF 入力を 3 実装 (ネイティブ Temporal / temporal-polyfill / Go サーバー) に通す単一画面。「解析・検証」「往復・実装差比較」「Temporal ラボ」の 3 モードを切り替える。解析・検証と往復・実装差比較の 2 モードは入力欄・ハイライト (F-2-5)・3 実装を共有する。Temporal ラボ モードは自由入力ではなく固定 fixture による日時モデル観察であり (F-6)、共有入力欄を持たず各実験内に 3 実装を併記する。

- **F-2-0** モードを「解析・検証」「往復・実装差比較」「Temporal ラボ」で切り替えられること
  - 解析・検証 / 往復・実装差比較は入力・ハイライト・3 実装を共有する (F-2-1〜F-2-6 / F-3)
  - Temporal ラボ は固定 fixture 駆動の観察モードで共有入力欄を持たない (F-6)
- **F-2-1** IXDTF 文字列を入力すると、デバウンス付きでリアルタイムに解析結果を表示すること
- **F-2-2** strict モードをトグルで切り替えられること
  (strict 時はオフセットとタイムゾーンの整合性を検証 — ixdtf の `strict` 引数に対応)
- **F-2-3** サーバー解析結果 (`POST /api/ixdtf/parse`) として以下を表示すること
  - 解析された時刻 (RFC 3339 表記)
  - タイムゾーン名と critical フラグ (`[!Asia/Tokyo]` の `!`)
  - 拡張タグの一覧 (キー / 値 / critical) をテーブル表示
  - エラー時はエラーメッセージ全文 (ixdtf のエラーは規格上の理由を含むため教材になる)
- **F-2-4** 同じ入力に対するブラウザ側 `Temporal.ZonedDateTime.from()` の結果を、ネイティブ Temporal と temporal-polyfill の 2 実装ぶん併記すること (サーバーと合わせ 3 実装)
- **F-2-5** 入力文字列をトークン分解し、F-1-2 と同じ色分けルールでハイライト表示すること
- **F-2-6** 「検証のみ」モード (ixdtf の `Validate` 相当) を選択できること
- **F-2-7** URL クエリで入力値・strict・表示モードを共有できること (Guide=F-5-3 / Converter=F-4-4 からの遷移に使用。旧 `/playground?input=&strict=` を後方互換で受理)

### F-3: 往復・実装差比較モード (F-2 ワークベンチ内)

F-2 ワークベンチの 1 モード。同一入力の 3 実装の挙動差そのものをコンテンツにする。

- **F-3-1** 入力した IXDTF 文字列を次の 3 経路で Parse → Format 往復させ、結果を並列表示すること
  - ネイティブ Temporal: `Temporal.ZonedDateTime.from(input).toString()`
  - temporal-polyfill: 同上 API を polyfill 実装で
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
- **F-4-4** 基準時刻を任意の IXDTF 文字列で指定できること (ワークベンチ=F-2 からの引き継ぎ含む)

### F-5: Guide — RFC 9557 学習ガイド

- **F-5-1** RFC 9557 の概要 (RFC 3339 との差分、サフィックス構文、critical フラグ、
  private / experimental 拡張の拒否ルール) を解説するコンテンツを表示すること
- **F-5-2** カテゴリ別のサンプル文字列集を提供すること
  - 基本 (RFC 3339 互換) / タイムゾーン付き / カレンダー付き / critical フラグ / 複合
  - **拒否されるケース**: `x-` プライベート拡張、`_` 実験的拡張、strict でのオフセット不整合
  - サンプルは RFC 9557 本文の例と代表的な落とし穴から **10〜15 個に厳選**する
    (網羅より読み切れる学習体験を優先)
- **F-5-3** 各サンプルに「ワークベンチで試す」ボタンを設け、F-2-7 のクエリ共有で遷移すること
- **F-5-4** 各サンプルに期待される結果 (成功 / 特定エラー) の注記を付けること

### F-6: Temporal ラボ モード — IXDTF 日時モデル (F-2 ワークベンチ内)

F-2 ワークベンチの第 3 モード。往復・実装差比較モード (F-3) とは目的を分け、IXDTF と
Temporal が持つ日時モデル上の性質そのものを小さな実験群として観察する。各実験はブラウザ
ネイティブ Temporal の実測値を主役に、Go ixdtf のライブ実測値と再現コード (F-0-7) を
併記する。ただし往復・実装差比較モードのような一致・不一致の採点は行わない — 両実装が
同じ日時モデルへ到達することを示す教材である。自由入力ではなく専用の固定 fixture で駆動し、
挙動差プリセット (`interop_preset: false` で分離) とは関連付けない。

- **F-6-1** DST 終了時に繰り返される同一ローカル時刻を使った「重複時刻」デモを表示すること
  - 同じ日付・時刻・IANA タイムゾーン名に対し、DST 前後の 2 つの有効な UTC オフセットを
    持つ IXDTF 文字列を並べる
  - 各文字列をブラウザネイティブの `Temporal.ZonedDateTime.from()` で解析し、解決された
    `epochNanoseconds` を表示する
  - 壁時計表示が同じでも 2 つの瞬間が 1 時間離れていることを時間軸で可視化する
  - 往復・実装差比較モードの挙動差プリセットとは関連付けず、match / mismatch の採点 UI も置かないこと
    (F-6-5 に基づく Go ixdtf 実測値の併記は行う)
- **F-6-2** 画面内の Temporal 実行結果はブラウザ内蔵実装の実測値であることを明示し、
  一次情報への解説・参照導線を表示すること
- **F-6-3** DST 境界をまたぐ「+1 日 (カレンダー演算)」と「+24 時間 (実時間演算)」の差を示す
  タイムゾーン演算デモを表示すること
  - F-6-1 と同じ DST 終了前夜の IXDTF 文字列を基点に `Temporal.ZonedDateTime.add()` で
    2 通りの加算を行い、結果の IXDTF 文字列と `epochNanoseconds`、基点からの実時間差を並べる
  - タイムゾーン注釈があるからこそ直列化後もカレンダー演算が正しく行える
    (数値オフセットだけの RFC 3339 では +1 日がゾーンの規則を反映できない) ことを説明する
  - Temporal が生成した演算結果の IXDTF 文字列を `POST /api/ixdtf/parse` (strict) で
    Go ixdtf に解析させ、両実装が同じ瞬間に到達することを `unix_nano` で示す
    (IXDTF が実装間の交換形式として機能する見せ場)
- **F-6-4** `u-ca` 拡張タグによるカレンダー投影デモを表示すること
  - 同一の瞬間を `Temporal.ZonedDateTime.withCalendar()` で iso8601 / japanese / hebrew /
    islamic-umalqura に投影し、紀年法・年月日と IXDTF 表記 (u-ca タグ) を並べる
  - `epochNanoseconds` が全投影で不変であること (注釈は瞬間を変えない) を明示する
  - 同じ文字列を `POST /api/ixdtf/roundtrip` で Go ixdtf に往復させ、Go はタグを解釈せず
    losslessly 運搬することを実測で示す (注釈の解釈は処理系の裁量という RFC 9557 の設計)
- **F-6-5** 各実験に Go ixdtf のライブ実測値 (既存 4 エンドポイント経由) と、
  F-0-7 準拠の再現セクション (Go / JavaScript サンプルコード) を併記すること

### F-7: About — IXDTF と Temporal の解説

Guide (F-5) が RFC 9557 の構文詳細を扱うのに対し、About は「Temporal とは何か」
「IXDTF とは何か」「両者はどう関わるのか」という前提概念を静的コンテンツで解説する
入口ページ。API 呼び出しは行わない。

- **F-7-1** IXDTF (RFC 9557) が「データ形式の規格」であることを解説すること
  - RFC 3339 を拡張した IETF の規格であること、角括弧サフィックス (タイムゾーン注釈・
    拡張タグ・critical フラグ) の概要、発行の経緯 (IETF SEDATE WG、2024 年) に触れる
- **F-7-2** Temporal が「ECMAScript の日時 API (日時モデル)」であることを解説すること
  - TC39 で標準化が進む `Date` の後継 API であること、`Temporal.ZonedDateTime` が
    タイムゾーンとカレンダーを値として保持すること、ネイティブ実装のブラウザ対応状況
    (N-1 と同期) に触れる
- **F-7-3** 両者の関係性を解説すること
  - Temporal の直列化 (`ZonedDateTime.toString()` / `from()`) には従来の RFC 3339 で
    運べない情報 (IANA タイムゾーン名・カレンダー) が必要であり、その文字列形式が
    IETF で RFC 9557 として標準化されたという経緯
  - IXDTF は実装間 (ブラウザ Temporal ⇄ Go ixdtf 等) の交換形式として機能すること
  - Temporal は RFC 9557 を自身のプロファイルで解釈するため実装間で挙動差が
    あり得ること (詳細は Interop / Guide へ誘導)
- **F-7-4** F-0-6 の解説・参照導線 (ReferenceDialog) で一次情報へ誘導し、
  関連画面 (ワークベンチ / Guide) への遷移導線を設けること

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

- **N-1 ブラウザ対応**: ネイティブ Temporal を持つ Chrome 144+ / Edge 144+ / Firefox 139+ を主対象とし、
  ネイティブとの挙動比較のために temporal-polyfill を同梱して実行時に切り替えられる (F-0-4)。
  Safari など安定版が未対応 (2026-07 時点) のブラウザは polyfill をデフォルトとして動作する
- **N-2 i18n**: 全ユーザー向け文字列は `t()` 経由。ja / en の locale ファイルはキー構造を完全同期
- **N-3 アクセシビリティ**: キーボード操作可能、カラーハイライトは色以外の手掛かり (下線・ラベル) を併用
- **N-4 パフォーマンス**: 時計更新 (1 秒間隔) で再レンダリングが時計コンポーネント内に閉じること
- **N-5 品質ゲート**:
  - web: Biome (lint + format)、Vitest (unit)、Storybook (コンポーネントカタログ + interaction test)
  - server: golangci-lint、table-driven test
  - OpenAPI スキーマと生成コード (Orval / oapi-codegen) の乖離を CI で検出すること
- **N-6 依存最小**: Go サーバーは ixdtf + 標準ライブラリを基本とする (DB なし、認証なし)
- **N-7 SEO**: 検索エンジンが各画面を個別のページとして扱えること
  - 画面ごとの `<title>` / `<meta name="description">` / `<link rel="canonical">` を、ルート遷移と
    UI 言語 (F-0-2) の切り替えに追従させる。文言は各画面の locale (`<page>.title` / `<page>.tagline`) を流用する (N-2)
  - canonical は正規ドメインの `origin + pathname` とし、共有 URL のクエリ (F-2-7) は含めない。
    未定義パスは `noindex` にする
  - `robots.txt` と `sitemap.xml` を配信する。sitemap には実在する画面のみを載せ、リダイレクト専用の旧パスは載せない
  - 旧パス (`/interop` `/temporal-lab`) はエッジで 301 リダイレクトする
- **N-8 アクセス解析**: 広告は置かない。アクセス解析は Cookie を使わず個人データを収集しない
  Cloudflare Web Analytics のみとし、その旨をサポート画面 (`support.why.body`) で明示する

## スコープ外 (Non-goals)

- ユーザーアカウント、データ永続化 (localStorage を除く)
- デプロイ / ホスティング構成 (ローカル実行のみを想定)
- Temporal API 全機能の網羅 (IXDTF 文字列に関係する範囲に限定)
