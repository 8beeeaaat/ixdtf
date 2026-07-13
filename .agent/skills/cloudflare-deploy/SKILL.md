---
name: cloudflare-deploy
description: Cloudflare Workers への本番デプロイワークフロー (Go WASM + Static Assets)。(1) デプロイ・本番反映・リリースの依頼時、(2) wrangler / make deploy / workers.dev に言及がある時、(3) API エンドポイント追加後に worker への反映を確認する時、(4) 本番の error 1042 / 404 / time_zone null 等のデプロイ起因の不具合調査時に使用。
---

# Cloudflare デプロイワークフロー

本番 URL (正規): **https://ixdtf.8beeeaaat.com** (カスタムドメイン、2026-07-13〜)
副系: **https://ixdtf-demo.8beeeaaat.workers.dev** (workers.dev、併存)

> `wrangler.jsonc` の `routes` に `custom_domain:true` で公開。routes を書くと `workers_dev` が
> デフォルト無効化されるため、workers.dev 併存には `workers_dev: true` の明示が必須。

API は `server/cmd/worker` を `GOOS=js GOARCH=wasm` でビルドした WASM (syumai/workers)、
フロントエンドは Workers Static Assets (`web/dist`) で配信する。`/api/*` のみ Worker が処理し、
他パスはアセット直配信 + SPA フォールバック。設計判断は `docs/architecture.md`「デプロイ」節と D-9 を参照。

## 前提

- 裸の `wrangler` はリポジトリルートの Node (nodenv) に無く失敗する — 常に `npx -y wrangler@4` を使う
- 認証が必要 (`npx -y wrangler@4 whoami` で確認)。期限切れならユーザーに
  `! npx -y wrangler@4 login` の実行を依頼する (ブラウザ OAuth のためエージェントからは不可)

## デプロイ手順

```bash
make deploy   # リポジトリルートで実行 (vite build → make build-worker → wrangler deploy)
```

- 必ず**リポジトリルート**から実行する (wrangler の build.command が `make build-worker` を呼ぶため、CWD がずれると "No rule to make target" で失敗する)
- デプロイ後は本番スモークを行う:

```bash
BASE=https://ixdtf-demo.8beeeaaat.workers.dev
curl -s $BASE/api/now
curl -s -X POST $BASE/api/ixdtf/parse -H 'Content-Type: application/json' \
  -d '{"input":"2026-07-07T23:30:00.123456789+09:00[Asia/Tokyo]","strict":false}'
# → "time_zone":"Asia/Tokyo" が返ること (null なら tzdata 欠落を疑う)
curl -s -o /dev/null -w "%{http_code}\n" $BASE/playground   # SPA フォールバック 200
```

## ローカル検証 (デプロイ前)

```bash
npx -y wrangler@4 dev --port 8787   # リポジトリルートから。未ログインでも動く
```

## 破ってはいけない制約

- `cmd/worker` は `framework` 本体を import しない — dist の `go:embed` が WASM に混入して
  サイズ制限を圧迫する。API ルートは `framework/api` を共有する (server/worker 両対応)
- API エンドポイントの追加は `framework/api/api.go` に登録すれば両エントリポイントに反映される
- `cmd/worker` の `_ "time/tzdata"` blank import を消さない — Workers には OS の zoneinfo が
  無く、`[Asia/Tokyo]` 等の解決が **ok:true のまま time_zone だけ null になる silent 劣化**を起こす
- WASM サイズは gzip 後 3MB (無料プラン上限) 以内に保つ。計測:
  `gzip -9 -c server/build/app.wasm | wc -c` (現状 約 1.8MB)

## トラブルシューティング

- デプロイ直後に error 1042 / 404 が混在 → Cloudflare 側の伝播遅延。数十秒待って再確認する
- `make build-worker` が "No rule to make target" → CWD がリポジトリルートでない
- `nodenv: wrangler: command not found` → 裸の `wrangler` を使っている。`npx -y wrangler@4` に置き換える
- 本番で `time_zone: null` (ネイティブでは値が返る) → worker バイナリに tzdata が入っていない
