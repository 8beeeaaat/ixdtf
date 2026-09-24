// Package api は API ルート定義。framework 本体から分離することで、
// dist を embed せずに Cloudflare Workers (WASM) の composition root
// からも同じルート定義を共有できるようにする。
package api

import (
	"net/http"

	"github.com/8beeeaaat/ixdtf/demo/server/controller"
)

// RegisterRoutes は API ルートを mux に登録する。
// Go 1.22+ の method 付きパターンでルーティングする。
func RegisterRoutes(mux *http.ServeMux, c *controller.Ixdtf) {
	mux.HandleFunc("POST /api/ixdtf/parse", c.Parse)
	mux.HandleFunc("POST /api/ixdtf/format", c.Format)
	mux.HandleFunc("POST /api/ixdtf/roundtrip", c.Roundtrip)
	mux.HandleFunc("GET /api/now", c.Now)
}
