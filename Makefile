# Trimora dev shortcuts. Everything is runnable from the repo root.

.PHONY: test test-fast test-api test-ui lint typecheck build dev-api dev-web clean

test:        ## full stack: go + web + api smoke + playwright
	./scripts/test-all.sh

test-fast:   ## same as `test` but skip Playwright UI checks
	./scripts/test-all.sh --skip-ui

test-api:    ## curl-driven API endpoint matrix only (assumes API on :8080)
	./scripts/smoke-api.sh

test-ui:     ## Playwright UI smoke only (assumes web on :4173 + api on :8080)
	node ./scripts/smoke-ui.mjs

lint:        ## eslint over web
	cd web && npm run lint

typecheck:   ## tsc --noEmit over web
	cd web && npm run typecheck

build:       ## build api binary + web bundle
	cd api && go build ./...
	cd web && npm run build

dev-api:     ## run api with `go run`
	cd api && go run ./cmd/server

dev-web:     ## run vite dev server
	cd web && npm run dev
