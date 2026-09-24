// Cloudflare Workers エントリラッパー (手書き)。
//
// 目的:
//  1. 旧 URL (ixdtf-demo.8beeeaaat.workers.dev) へのアクセスを正規ドメイン
//     (ixdtf.8beeeaaat.com) へ 301 恒久リダイレクトする。
//  2. ワークベンチへ統合した旧パス (/interop /temporal-lab) を、モード付きの
//     /playground へ 301 恒久リダイレクトする (N-7)。クライアント側の redirect と同じ対応。
//  3. それ以外は従来どおり /api/* を Go WASM worker へ、
//     残りを Static Assets へ委譲する。
//
// wrangler.jsonc の run_worker_first: true により全リクエストが本 worker を
// 先に通る。build/worker.mjs は make build-worker が再生成する生成物のため、
// リダイレクト処理は再生成の影響を受けない手書きの外側ラッパーに置く。
import goWorker from "./build/worker.mjs";

// 旧 workers.dev サブドメイン → 正規カスタムドメイン。
const LEGACY_HOST = "ixdtf-demo.8beeeaaat.workers.dev";
const CANONICAL_HOST = "ixdtf.8beeeaaat.com";

// 旧パス → ワークベンチのモード (web/src/app/router.tsx の redirect と対応)。
const LEGACY_WORKBENCH_PATHS = new Map([
  ["/interop", "roundtrip"],
  ["/temporal-lab", "lab"],
]);

export default {
  async fetch(req, env, ctx) {
    const url = new URL(req.url);

    // 1. 旧 URL は path/query を保ったまま正規ドメインへ 301 恒久リダイレクト。
    if (url.hostname === LEGACY_HOST) {
      url.hostname = CANONICAL_HOST;
      return Response.redirect(url.toString(), 301);
    }

    // 2. 旧パスは他のクエリを保ったまま mode 付きの /playground へ 301 恒久リダイレクト。
    const legacyMode = LEGACY_WORKBENCH_PATHS.get(url.pathname.replace(/\/+$/, ""));
    if (legacyMode) {
      url.pathname = "/playground";
      url.searchParams.set("mode", legacyMode);
      return Response.redirect(url.toString(), 301);
    }

    // 3. API は Go WASM worker (composition root) へ委譲。
    //    WASM はこの分岐に入ったときだけ遅延ロードされる。
    if (url.pathname.startsWith("/api/")) {
      return goWorker.fetch(req, env, ctx);
    }

    // 4. 残りは Static Assets へ。not_found_handling(SPA フォールバック) は
    //    ASSETS binding の fetch が尊重する。
    return env.ASSETS.fetch(req);
  },
};
