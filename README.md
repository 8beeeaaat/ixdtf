# ixdtf_demo
A demo for https://pkg.go.dev/github.com/8beeeaaat/ixdtf

## Docs

- [機能要件](./docs/requirements.md)
- [設計書](./docs/architecture.md)
- [UI デザインシステム](./docs/DESIGN.md)

## 開発支援

実装フェーズの規約を自動で守らせるプロジェクト設定。
Skills は `.agent/skills/` を実体とし、Claude Code は `.claude/skills`、Codex は `.agents/skills`
の symlink から同じ内容を参照する。

### Agents (自動レビュー)

Claude Code 用のプロジェクト agents は `.claude/agents/` に置く。

| Agent | 検査対象 | 規範 |
|---|---|---|
| `design-system-reviewer` | `web/` の .tsx / .css / stories | docs/DESIGN.md |
| `ixdtf-domain-reviewer` | IXDTF セマンティクス (server/, web/src/lib/ixdtf, openapi.yaml) | requirements.md A-1〜A-3, architecture.md D-1〜D-4, RFC 9557 |
| `clean-arch-reviewer` | `server/` のレイヤー依存方向 | architecture.md 依存図 |
| `spec-coverage-tracer` | 要件 ID ↔ 実装・テストの充足 | requirements.md F/A/N 採番 |

### Skills

- `api-contract-sync` — `api/openapi.yaml` (SSOT) 変更時のスキーマ駆動ワークフロー
- `ixdtf-fixtures` — `testdata/ixdtf_samples.json` による Go / TS 共有テストフィクスチャ規約

### Hooks

- Claude Code: `.claude/hooks/`
- Codex: `.codex/hooks.json` + `.codex/hooks/`。初回または変更後は Codex の `/hooks` で内容を確認して trust する
- `guard-generated.sh` (PreToolUse) — codegen 生成物 (`*/generated/`) の手動編集をブロック
- `post-edit-checks.sh` (PostToolUse) — openapi / locale 変更の追従リマインド、DESIGN.md 禁止パターン警告
