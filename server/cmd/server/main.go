// Command server は IXDTF デモの HTTP サーバー。
// composition root として interactor を生成し controller へ注入、
// framework の mux を :8080 で起動する。
package main

import (
	"log"
	"net/http"

	"github.com/8beeeaaat/ixdtf_demo/server/controller"
	"github.com/8beeeaaat/ixdtf_demo/server/framework"
	"github.com/8beeeaaat/ixdtf_demo/server/usecase/interactor"
)

func main() {
	usecase := interactor.NewIxdtf()
	ctrl := controller.NewIxdtf(usecase)
	mux := framework.NewMux(ctrl)

	const addr = ":8080"
	log.Printf("ixdtf_demo server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil { //nolint:gosec // demo server, no timeouts needed
		log.Fatal(err)
	}
}
