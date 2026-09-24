# IXDTF デモ 収益化プラン

このデモは RFC 9557 (IXDTF) の**開発者向け教育デモ**であり、広告なし・トラッキングなしで
運用している。本ドキュメントは収益化の方針とその実装、および将来の発展経路を記録する。

## 採用モデル: スポンサー / サポート導線

### なぜこのモデルか

| 観点 | 判断 |
|---|---|
| 設計制約との整合 | **N-6 (DB なし・認証なし) / D-6 (model/ gateway/ を作らない) を一切壊さない**。変更は web/ のプレゼンテーション層に閉じ、Go サーバー (単一バイナリ / Workers WASM) は無改変 |
| プロダクトの性質 | ニッチな規格の教育デモ。有料ペイウォールより「無料で価値を出し、共感した人が支援する」導線が正直 |
| 実装コスト | 静的リンク + 1 ページ追加のみ。運用負荷ゼロ (課金・鍵管理・レート制限の基盤が不要) |
| 収益の現実性 | 後述のとおり規模は小さい。**事業ではなく、維持・改善の持続可能性を支える寄付**と位置づける |

### 検討したが採用しなかった案

| 案 | 不採用の理由 |
|---|---|
| 従量課金 API (API キー + 有料ティア) | 収益ポテンシャルは最大だが、鍵保存・レート制限・課金基盤が必要で N-6 / D-6 に正面から抵触。ニッチな規格 API に有料需要が立つ確度も低い。→ **将来の発展経路として保留** (下記) |
| リードジェン (ixdtf ライブラリ / コンサル誘導) | 各画面に CTA を撒くのは教育デモの体験を損なう。メール収集に外部サービス連携も要る |
| 広告 | クリーンな教育ツールの体験を大きく毀損し、この規模では収益も僅少 |

## 実装したもの

| 要素 | ファイル | 役割 |
|---|---|---|
| Sponsor ボタン | `web/src/components/Layout.tsx` (ヘッダー) | `/support` へ誘導。モノクロのハート (DESIGN.md: 彩度色は IXDTF ハイライト/状態表示に限定) |
| フッター導線 | `web/src/components/Footer.tsx` (+ `.stories.tsx`) | 「無料・OSS」の明示 + リポジトリ / Buy Me a Coffee / 支援ページへのリンク |
| サポートページ | `web/src/pages/support/SupportPage.tsx` | 「デモの意義」を先に説明し、その後に GitHub Sponsors / Buy Me a Coffee / スター & 拡散の 3 導線 (F-5 GuidePage のレイアウトを踏襲) |
| ルート | `web/src/app/router.tsx` | `/support` を追加 (6 画面のメインナビには含めず、ヘッダー/フッターから到達) |
| リンク一元管理 | `web/src/lib/sponsor.ts` | 外部 URL の単一情報源。差し替えはここだけ |
| i18n | `web/src/locales/{ja,en}/translation.json` | `sponsor` / `footer` / `support` namespace を両 locale 同期 (N-2) |
| リポジトリ Sponsor ボタン | `.github/FUNDING.yml` (ixdtf リポジトリルート) | GitHub リポジトリページに Sponsor ボタンを出すゼロコード施策 |

### 設計上の線引き

- **ヘッダーの Sponsor ボタンは `/support` (内部) へ誘導**し、外部課金導線へ直行させない。
  「価値を説明してから支援を募る」ほうが誠実で、コンバージョンにも資する
- ハートは慣習の赤ではなく `text-foreground` のモノクロ。既存の GitHub Octicon と同じ扱い
- 6 画面のメインナビへ Support は加えない (DESIGN.md: アイコン/リンク過多を避ける)

## 有効化チェックリスト (アカウント側の作業)

コードは投入済み。以下はリポジトリ所有者が行う:

- [ ] GitHub Sponsors を有効化する (`https://github.com/sponsors/8beeeaaat`)
- [ ] Buy Me a Coffee アカウントを作成する (未作成なら `web/src/lib/sponsor.ts` と
      `.github/FUNDING.yml` (ixdtf リポジトリルート) の URL を実在のものに差し替えるか、当面は該当リンクを外す)
- [ ] `.github/FUNDING.yml` (ixdtf リポジトリルート) をデフォルトブランチに置き、リポジトリの Sponsor ボタン表示を確認する
- [ ] `web/src/lib/sponsor.ts` の 3 URL が実在することを最終確認する

> リンクは `lib/sponsor.ts` に集約済み。存在しないアカウントへ誘導しないよう、
> 有効化が済むまでは未使用の導線を外す運用でもよい。

## 将来の発展経路 (トラフィックが増えた場合)

需要が確認できたら、従量課金 API へ発展できる。既存の 4 エンドポイント
(parse / format / roundtrip / now) はそのまま製品化の資産になる:

1. Cloudflare Workers KV に API キーと使用量を保存 (D-6 を「DB あり」に更新)
2. `framework/api` にレート制限ミドルウェアを追加 (無料 60 req/min 等)、超過は 429
3. `/pricing` `/account` を追加し、鍵発行と無料/有料ティアを提供
4. 課金は外部プロバイダ (Stripe 等) に委譲

この移行は本ドキュメントの「不採用案」を再評価するところから始める。
