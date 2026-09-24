# AGENTS.md

This file provides shared guidance to coding agents when working with code in this repository.

## プロジェクト概要と現状

RFC 9557 (IXDTF) のデモアプリ。Go 製 [ixdtf ライブラリ](https://github.com/8beeeaaat/ixdtf) とブラウザの TC39 Temporal API (ネイティブ / temporal-polyfill を実行時に切替可能、`app/temporal.tsx` の `TemporalProvider` が選択) で同じ IXDTF 文字列を往復させ、相互運用性そのものを展示する。「Temporal」は常に TC39 Temporal API を指す (Temporal.io ではない)。

**配置**: このデモは ixdtf ライブラリのリポジトリの `demo/` 配下にある (リポジトリルートはライブラリ本体、別モジュール)。Claude Code / Codex は `demo/` で起動すること — agents・skills・hooks の設定は `demo/` 配下にあり、本書のパスとコマンドもすべて `demo/` 基準。ixdtf への依存は `server/go.mod` の公開タグで固定しているため、ルートのライブラリ変更はタグを打って依存を上げるまでデモに反映されない (architecture.md D-14)。

**現状: 設計フェーズ完了・実装未着手。** コードを書く前に必ず docs/ を読むこと。docs/ が唯一の規範であり、実装と docs が食い違ったら docs 更新を先行させる:

- `docs/requirements.md` — 採番済み要件 (F-0〜F-6 機能 / A-1〜A-3 API / N-1〜N-6 非機能)。受け入れ条件の SSOT
- `docs/architecture.md` — システム構成、OpenAPI スキーマ定義、Go レイヤー設計、設計判断 D-1〜D-8
- `docs/DESIGN.md` — UI デザインシステム (タイポグラフィ・ミニマル)

## コマンド

Makefile のターゲット (詳細は architecture.md「開発ワークフロー」):

```
make generate   # orval + oapi-codegen (api/openapi.yaml 変更時に必須)
make dev        # go run ./server/cmd/server + vite dev (/api は Vite proxy で :8080 へ)
make test       # go test ./server/... + vitest run
make lint       # golangci-lint + biome check
make build      # vite build → go:embed → 単一バイナリ
make deploy     # Cloudflare Workers へデプロイ (`cloudflare-deploy` スキル参照)
```

単一テスト実行:
- Go: `go test ./server/usecase/interactor/ -run TestXxx` (server/ ディレクトリで)
- Web: `npx vitest run src/lib/ixdtf/tokenize.test.ts` (web/ ディレクトリで)

## アーキテクチャ (詳細は architecture.md)

- **スキーマ駆動**: `api/openapi.yaml` が API 契約の SSOT。API 型の手書きは禁止 — orval → `web/src/generated/`、oapi-codegen → `server/generated/` で生成。generated/ の手動編集は PreToolUse hook がブロックする。変更手順は `api-contract-sync` スキル参照
- **server/** (module: `github.com/8beeeaaat/ixdtf/demo/server`): Clean Architecture ライト。依存方向は framework → controller → (presenter / usecase/inputport / generated)、interactor → (inputport 実装 / entity / ixdtf)。`model/` `gateway/` は作らない (DB なし、D-6)。composition root は cmd/server のみ
- **web/**: React + Vite + TS。TanStack Router (型付き search params が Playground 共有 URL = F-2-7) / Query (Orval 生成 hooks) / Form。i18n は react-i18next (ja/en)
- エンドポイントは 4 つのみ: parse / format / roundtrip / now

## 破ってはいけないドメインルール

- `unix_nano` は JSON 全経路で 10 進**文字列** (int64 は JS の MAX_SAFE_INTEGER 超、D-3)
- 拡張タグは**順序付き配列**で運ぶ (Go の map で順序を壊さない)。critical `!` フラグはタグ単位で保持 (D-4)
- 解析失敗は HTTP 200 + `ok: false` のドメイン結果 (D-1)。4xx はリクエスト不備のみ (A-2)。エラー文言は ixdtf ライブラリ原文のまま返す (A-3)
- IXDTF サンプル文字列はテストに直書きせず `testdata/ixdtf_samples.json` に一元管理 (`ixdtf-fixtures` スキル)
- `web/src/lib/ixdtf/tokenize.ts` は表示専用 — 正当性判定を実装しない (規格解釈を 3 つ目に増やさない)

## UI ルール (DESIGN.md の抜粋)

- className はセマンティックトークンのみ (生の色クラス禁止)。`glass` 系・`dark:` バリアント禁止 (テーマはトークン自動切替)
- `clsx` 直接 import 禁止 — `cn` (`@/lib/utils`)。バリアントは cva
- IXDTF 文字列・時計は `font-mono tabular-nums` 必須 (1 秒 tick のレイアウトシフト防止)、ハイライトは `text-ixdtf-*` トークン
- 表示文字列は `t()` 経由、ja / en の locale キー構造を同時に同期
- `src/components/ui/` の新規コンポーネントには `.stories.tsx` を同居させる

## 開発支援

- **共有 Skills**: 実体は `.agent/skills/`。Claude Code は `.claude/skills`、Codex は `.agents/skills` の symlink から同じ Skill を読む

- **Agents** (変更後に proactive 起動): `design-system-reviewer` (web/ の UI 変更) / `ixdtf-domain-reviewer` (ドメイン・API 変更) / `clean-arch-reviewer` (server/ 変更) / `spec-coverage-tracer` (要件 ID ↔ 実装の突き合わせ) / `storybook-steward` (web/ コンポーネント変更時の Storybook 製造・カバレッジ・interaction テスト検証)
- **Hooks**: Claude Code は `.claude/hooks/`、Codex は `.codex/hooks.json` + `.codex/hooks/`。generated/ 編集ブロック、openapi.yaml / locale / .tsx 変更時の追従リマインドを行う
- 実装は要件 ID 単位で進め、コミット・要件の完了判定はユーザーが行う
