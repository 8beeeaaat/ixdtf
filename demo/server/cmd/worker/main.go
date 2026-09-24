// Command worker は Cloudflare Workers 向けエントリポイント (GOOS=js GOARCH=wasm)。
// composition root として interactor を生成し controller へ注入、
// API ルートのみを syumai/workers 経由で配信する。
// 静的アセットは Workers Static Assets (wrangler.jsonc の assets) が担うため、
// dist を embed している framework パッケージは import しない。
package main

import (
	"net/http"

	// Workers 実行環境には OS の zoneinfo が無く、IANA タイムゾーン名
	// ([Asia/Tokyo] 等) の解決が失敗するため tzdata をバイナリに埋め込む。
	_ "time/tzdata"

	"github.com/8beeeaaat/ixdtf/demo/server/controller"
	"github.com/8beeeaaat/ixdtf/demo/server/framework/api"
	"github.com/8beeeaaat/ixdtf/demo/server/usecase/interactor"
	"github.com/syumai/workers"
)

func main() {
	usecase := interactor.NewIxdtf()
	ctrl := controller.NewIxdtf(usecase)
	mux := http.NewServeMux()
	api.RegisterRoutes(mux, ctrl)
	workers.Serve(mux)
}
