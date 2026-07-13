.PHONY: generate generate-server generate-web dev dev-server dev-web test test-server test-web lint lint-server lint-web build build-worker deploy

OAPI_CODEGEN := go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.7.2

generate: generate-web generate-server

# SSOT は OpenAPI 3.1 (api/openapi.yaml)。oapi-codegen は 3.1 未対応のため、
# openapi-down-convert で 3.0 へ変換した中間ファイル (git 管理外) を入力にする。
generate-server:
	cd web && npx openapi-down-convert --input ../api/openapi.yaml --output ../server/openapi-3.0.gen.yaml
	cd server && $(OAPI_CODEGEN) -config oapi-codegen.yaml openapi-3.0.gen.yaml

generate-web:
	cd web && npx orval --config orval.config.ts

dev:
	$(MAKE) -j2 dev-server dev-web

dev-server:
	cd server && go run ./cmd/server

dev-web:
	cd web && npm run dev

test: test-server test-web

test-server:
	cd server && go test ./...

test-web:
	cd web && npx vitest run

lint: lint-server lint-web

lint-server:
	cd server && golangci-lint run ./...

lint-web:
	cd web && npx biome check .

build:
	cd web && npm run build
	rm -rf server/framework/dist && cp -R web/dist server/framework/dist
	cd server && go build -o ../bin/ixdtf_demo ./cmd/server

# Cloudflare Workers 用 WASM ビルド (wrangler の build.command からも呼ばれる)
build-worker:
	cd server && go run github.com/syumai/workers/cmd/workers-assets-gen -mode=go -o ./build
	cd server && GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o ./build/app.wasm ./cmd/worker

# Cloudflare へデプロイ (要 wrangler login)。wrangler が build.command 経由で build-worker を実行する
deploy:
	cd web && npm run build
	npx -y wrangler@4 deploy
