---
lessonId: "P0.7"
title: "工具鏈一日通"
description: "go build/test/mod/vet/fmt 與日常工作流。"
volume: "p0"
order: 7
level: "l1"
status: "ready"
path_required: true
tags: ["toolchain"]
example: "examples/p0-config-stats"
prev: "P0.6"
next: "P0.8"
---

## 本章你會建立的心智模型

Go 的生產力很大一部分來自**官方工具鏈是預設文化**：格式化不爭論、測試在標準位置、競態偵測器內建、依賴有語意版本與 MVS。先把命令肌肉記住，後面深度章才有手感。

## Python 對照

| Python | Go |
|--------|-----|
| `black` / `ruff`（外加） | `gofmt` / `go fmt`（官方） |
| `pytest` | `go test` |
| `pip` / `uv` / `poetry` | `go mod` |
| `mypy` / linters | `go vet`、staticcheck 等 |

## L1 能用

```bash
go version
go env GOPATH GOMODCACHE

go mod init example.com/x
go get example.com/some/module@v1.2.3
go mod tidy

go fmt ./...
go vet ./...
go test ./...
go test -race ./...
go test -bench=. -benchmem ./...

go build -o server.exe .
go run .
```

## L2 機制

- **`gofmt` 是社交協定**：PR 不應討論空白。
- **`go test`**：`_test.go` 與被測包同目錄（或 `package foo_test` 黑箱）。
- **`-race`**：需要 cgo 工具鏈在某些平台；Windows 上通常可用，但較慢。
- **Modules**：語意化版本；不要手改 `go.sum` 除非你知道在做什麼。

## L3 深潛（可選）

- `GOTOOLCHAIN` 自動下載新版 toolchain 的行為。
- build cache 位置與 `go clean -cache`。

## 請丟掉的 Python 習慣

1. 每人一套格式化設定開戰。
2. 「測試框架選型」耗一週——先 `go test` 表驅動。
3. 全域亂裝依賴卻不鎖版本——看 `go.mod`/`go.sum`。

## 遊戲 Server 連結

CI 最低標：`go test ./...` + 併發包 `-race` + `go vet`。效能章再上 pprof 與 benchmark。

## 練習

### 必做

1. 在 `examples/p0-config-stats` 跑 `go test ./...` 與 `go fmt ./...`。
2. 故意破壞格式，再 `go fmt` 看差異。

### 選做

1. 安裝 `staticcheck`，對範例跑一次。

## 常見坑與如何看見

- 代理與防火牆：設定 `GOPROXY`（企業環境常見）。
- 忘記 module：在無 `go.mod` 目錄乱跑舊 GOPATH 模式。

## 延伸閱讀

- <https://go.dev/doc/diagnostic.html>
