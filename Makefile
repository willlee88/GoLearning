.PHONY: dev web-build content-check demo

dev:
	cd web && npm run dev

web-build:
	cd web && npm run build

content-check:
	node scripts/check-content.mjs

demo:
	cd demo/arena-mini/server && go run ./cmd/server
