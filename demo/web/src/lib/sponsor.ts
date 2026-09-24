/**
 * サポート導線の外部リンク一元管理 (収益化: スポンサー/サポートモデル)。
 * アカウントを有効化・変更したらここだけを直せばよい (docs/monetization.md)。
 * DB・認証を持たないこのデモの制約 (N-6 / D-6) を壊さず、静的リンクのみで完結させる。
 */
export const SPONSOR_LINKS = {
  /** GitHub Sponsors (ユーザー名から導出可能な主導線)。要: Sponsors の有効化。 */
  githubSponsors: "https://github.com/sponsors/8beeeaaat",
  /** Buy Me a Coffee (単発の寄付)。要: アカウント作成。 */
  buyMeACoffee: "https://www.buymeacoffee.com/8beeeaaat",
  /** このデモのソース (star / 拡散導線)。ixdtf ライブラリのリポジトリの demo/ 配下にある。 */
  repo: "https://github.com/8beeeaaat/ixdtf/tree/main/demo",
} as const;
