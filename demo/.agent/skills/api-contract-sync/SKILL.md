---
name: api-contract-sync
description: api/openapi.yaml (API 契約の SSOT) を変更する際のスキーマ駆動ワークフロー。(1) API のエンドポイント・型の追加/変更時、(2) server/generated・web/src/generated とスキーマの乖離が疑われる時、(3) orval / oapi-codegen / make generate / codegen に言及がある時に使用。手書き API 型の禁止と生成物 drift の解消を保証する。
---

# API 契約同期ワークフロー

`api/openapi.yaml` が唯一の契約。生成物 (`server/generated/`, `web/src/generated/`) の手動編集は
禁止 (PreToolUse hook でブロックされる)。詳細な設計判断は `docs/architecture.md`「スキーマ駆動開発」節を参照。

## 手順

1. `api/openapi.yaml` を編集する。記法チェック:
   - OpenAPI 3.1: nullable は `type: ["string", "null"]` (`nullable: true` は 3.0 記法で禁止)
   - 全レスポンスに `description: "OK"` (Orval validator が要求)
   - enum は snake_case、`operationId` は camelCase
   - int64 値 (unix_nano 等) は `type: string` で運ぶ (D-3)
2. `make generate` を実行 (orval + oapi-codegen の両方が走る)
3. `git diff` で生成物の変化をレビューする (意図しない breaking change の検出)
4. server 側: `presenter/` の変換と `controller/` のハンドラを新スキーマに追従させる
5. web 側: Orval 生成 hooks (TanStack Query) の利用側を更新する
6. `make test` で両側のテストを通す

## トラブルシューティング

- oapi-codegen が 3.1 記法でエラー → `make generate-server` が `openapi-down-convert` で
  3.0 中間ファイルへ自動変換済み。down-convert が扱えない `oneOf: [$ref, null]` は使わず、
  nullable なオブジェクト参照は optional な `$ref` で表現する (architecture.md リスク表参照)
- 生成物の drift が CI で検出された → ローカルで `make generate` を再実行してコミット
- スキーマにない型が必要になった → コードに手書きせず、必ずスキーマへ追加してから生成する
