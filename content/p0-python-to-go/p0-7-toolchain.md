---
lessonId: "P0.7"
title: "日常工具鏈：build、test、fmt、vet、mod"
description: "把 go build／test／fmt／vet／mod 變成肌肉記憶。格式化不爭論、測試在標準位置、競態偵測器內建。"
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

## 這章你會搞懂什麼

Go 的生產力，很大一部分來自：**官方工具鏈就是預設文化**。

- 格式化不爭論（`gofmt`）  
- 測試放標準位置（`go test`）  
- 競態偵測器內建（`go test -race`）  
- 依賴有語意版本與可預測的選版（modules + MVS）  

先把命令肌肉記住，後面深度章才有手感。讀完你要能在範例目錄跑通 `fmt`、`test`、`vet`，並知道 CI 最低標長什麼樣。

## 先跟 Python 對一下

| Python | Go |
|--------|-----|
| `black` / `ruff`（常要外加） | `gofmt` / `go fmt`（官方內建） |
| `pytest` 等測試框架 | `go test`（內建） |
| `pip` / `uv` / `poetry` | `go mod`／`go get` |
| `mypy`、各種 linter | `go vet`；社群還有 staticcheck 等 |

你少花時間在「選工具開戰」，多花時間在問題本身——這是刻意的。

## 怎麼寫（能跑的最小例子）

把這些當日常口訣：

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

請到 `examples/p0-config-stats` 實際跑：

```bash
go test ./...
go fmt ./...
go vet ./...
```

故意把縮排弄亂，再 `go fmt`，看它怎麼幫你收乾淨——這就是「格式化是社交協定」。

## 為什麼這樣設計／底層在幹嘛

### `gofmt` 是社交協定

PR 不該討論「要不要空格／換行風格」。機器統一格式，人腦留給邏輯與命名。不跑 `fmt` 會怎樣？code review 噪音變多，真正的 bug 被蓋住。

### `go test` 很直球

- 測試檔名：`*_test.go`  
- 通常跟被測套件同目錄；也可以 `package foo_test` 做黑箱測試  
- 表驅動測試（table-driven）是社群主流寫法——D 卷會練  

還沒寫測試時，`go test ./...` 仍常會成功（沒有測試檔不算失敗）。別誤會成「我測過了」。

### `-race` 為什麼重要

資料競爭是未定義行為。`-race` 會讓測試變慢，但能抓到很多「偶現、超難重現」的雷。Windows 上通常可用；某些平台需要 cgo 工具鏈——若跑不起來，先查官方說明，別直接放棄。

### Modules：別手賤改 `go.sum`

`go.mod`／`go.sum` 是依賴真相來源之一。不確定在做什麼就手改 `go.sum`，之後 `tidy` 與 CI 會跟你吵架。企業網路常要設 `GOPROXY`／`GOPRIVATE`。

### 進階可先略過

- `GOTOOLCHAIN`：自動下載／切換 toolchain 的行為。  
- build cache 位置與 `go clean -cache`。  

## 遊戲 Server 會用在哪

CI 最低標建議：

1. `go test ./...`  
2. 併發相關套件加 `-race`  
3. `go vet ./...`  

效能與瓶頸章再上 pprof、benchmark。優雅關機、metrics 也會依賴你「測得到、建得出 binary」這套肌肉。

## 請丟掉的舊習慣

1. 每人一套格式化設定開戰——交給 `gofmt`。  
2. 「測試框架選型」耗一週——先用內建 `go test` 表驅動。  
3. 全域亂裝依賴卻不鎖版本——看當前 module 的 `go.mod`／`go.sum`。  

## 動手練習

### 必做

1. 在 `examples/p0-config-stats` 跑 `go test ./...` 與 `go fmt ./...`。  
2. 故意破壞格式，再 `go fmt`，用 diff 看差在哪。  

### 選做

1. 安裝 `staticcheck`，對範例跑一次，讀它抱怨什麼。  

## 常見坑

- **代理與防火牆**：`go get` 失敗時檢查 `GOPROXY`（公司環境超常見）。  
- **忘記 module**：在沒有 `go.mod` 的目錄乱跑，掉進舊 GOPATH 心智。  
- **只跑 `go run` 從不 `go test`**：很多回歸要靠測試擋；Server 邏輯尤其是。  

## 延伸閱讀

- <https://go.dev/doc/diagnostic.html>  
- <https://go.dev/doc/modules/managing-dependencies>  
