---
name: ixdtf-fixtures
description: IXDTF サンプル文字列の共有フィクスチャ管理。(1) parse / format / roundtrip のテスト追加・変更時、(2) F-5 学習ガイドのサンプル編集時、(3) Go と TypeScript の両方で同じ文字列を検証したい時に使用。Go テスト・Vitest・Guide UI (F-5) が単一カタログを共有し、両実装が同じ入力で検証されることを保証する。
---

# IXDTF 共有フィクスチャ

サンプル文字列はリポジトリルートの `testdata/ixdtf_samples.json` に一元管理する。
Go テスト・Vitest・F-5 SampleGallery の 3 消費者がこのカタログを読む。
**個々のテストにサンプル文字列を直書きしない** (ロジック固有の境界値のみ例外)。

## カタログ形式 (合成例)

```json
{
  "samples": [
    {
      "id": "critical-unknown-tag",
      "input": "2026-01-02T03:04:05+09:00[Asia/Tokyo][!u-xx=1]",
      "category": "critical",
      "strict": false,
      "expect_server": "error",
      "expect_browser": "error",
      "expect_error_substr": "critical",
      "note_key": "guide.samples.criticalUnknownTag"
    }
  ]
}
```

- `category`: `basic` | `timezone` | `calendar` | `critical` | `composite` | `rejected`
- `expect_server` / `expect_browser` は分離して持つ — 両実装の仕様差 (F-3 の展示対象) を
  フィクスチャで表現できるようにするため。`ok` | `error` の 2 値
- `note_key` は Guide 表示用の i18n キー (ja / en 両方に追加すること)

## 運用ルール

1. サンプル追加はカタログ → Go テスト → Vitest → (必要なら) Guide の順で反映する
2. F-5-2 の決定に従い、Guide に出すのは **10〜15 個に厳選** (`category: rejected` を必ず含める)。
   テスト専用サンプルは `guide: false` フィールドで除外できる
3. 出典を優先する: RFC 9557 本文の例 → ixdtf ライブラリのテストケース → 自作
4. 期待値を変える変更は ixdtf-domain-reviewer エージェントの検査対象 (D-4 ロスレス性に影響)
