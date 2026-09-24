// Package framework は HTTP ルーティングと静的ファイル配信を担う。
// vite build 成果物を go:embed し、単一バイナリで配布できるようにする (D-7)。
package framework

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/8beeeaaat/ixdtf/demo/server/controller"
	"github.com/8beeeaaat/ixdtf/demo/server/framework/api"
)

//go:embed all:dist
var distFS embed.FS

// NewMux は API ルートと SPA 静的配信を束ねた http.Handler を返す。
// API ルート定義は framework/api を参照 (Workers エントリポイントと共有)。
func NewMux(c *controller.Ixdtf) http.Handler {
	mux := http.NewServeMux()
	api.RegisterRoutes(mux, c)
	mux.Handle("/", spaHandler())
	return mux
}

// spaHandler は embed した dist/ を配信する。存在しないパスは index.html に
// フォールバックする (SPA)。index.html 自体が無ければ 404 を返す。
func spaHandler() http.Handler {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		panic(err) // ビルド時 embed により dist は必ず存在する
	}
	fileServer := http.FileServer(http.FS(sub))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if name == "" {
			name = "index.html"
		}
		if f, err := sub.Open(name); err == nil {
			_ = f.Close()
			fileServer.ServeHTTP(w, r)
			return
		}
		// 未知パスは index.html を返す (SPA フォールバック)。
		if f, err := sub.Open("index.html"); err == nil {
			_ = f.Close()
			r.URL.Path = "/"
			fileServer.ServeHTTP(w, r)
			return
		}
		http.NotFound(w, r)
	})
}
