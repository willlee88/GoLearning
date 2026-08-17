---
lessonId: "A1"
title: "安裝、版本與第一個 module"
description: "裝好 Go、理解 go.mod，跑通最小程式。"
volume: "a"
order: 1
level: "l1"
status: "ready"
path_required: true
tags: ["toolchain", "modules"]
example: "examples/a01-hello-module"
prev: "A0"
next: "A2"
---

## 本章你會建立的心智模型

Go 的專案從 **module** 開始：`go.mod` 宣告模組路徑與語言版本，之後所有 `import`、依賴與建置都繞著它轉。先把「能建、能跑、能測」的肌肉記住。

## Python 對照

| Python | Go |
|--------|-----|
| 安裝 CPython / 用 pyenv | 安裝官方 toolchain；可設 `GOTOOLCHAIN` |
| `python -m venv` | 通常不需要 venv；module 管依賴 |
| `pip install` + lock | `go get` + `go.sum` |
| `python main.py` | `go run .` / `go build` |

## L1 能用

1. 安裝：<https://go.dev/dl/>（Windows 用 MSI 即可）  
2. 新開終端機：

```bash
go version
go env GOROOT GOPATH GOMODCACHE
```

3. 建立 module：

```bash
mkdir hello && cd hello
go mod init example.com/hello
```

```go
// main.go
package main

import "fmt"

func main() {
	fmt.Println("hello module")
}
```

```bash
go run .
go test ./...   # 尚無測試也會成功（無測試檔時）
```

本站範例：`examples/a01-hello-module`。

## L2 機制

- **`go.mod` 的 module path** 是匯入前綴，不必等於資料夾名，但團隊應約定一致。  
- **`go` 行**（如 `go 1.22`）表示語言／標準庫基線。  
- 本站教材以 **Go 1.22+** 為準（路由模式、`slog` 等後續會用到）。  
- `go run` 編譯到暫存再執行；`go build` 產出 binary。

## L3 深潛（可選）

- `GOTOOLCHAIN=auto` 時，新語法可能自動拉新版 toolchain。  
- `GOPROXY`、`GOPRIVATE` 在公司網路很常見。

## 請丟掉的 Python 習慣

1. 沒有「專案邊界」就到處 `import` 相對路徑亂檔。  
2. 依賴裝在全域 site-packages 心智——改看當前目錄的 `go.mod`。  
3. 用系統 Python 版本碰運氣——Go 版本寫進 `go.mod`。

## 遊戲 Server 連結

每個服務（gateway、room worker、假 client）都應是清晰 module 或 monorepo 中的清晰 package 樹。Arena Mini 在 `demo/arena-mini/server`。

## 練習

### 必做

1. 本機 `go version` 成功。  
2. 跑 `examples/a01-hello-module`。  

### 選做

1. `go build -o hello.exe .` 並直接執行 binary。  

## 常見坑與如何看見

- **PATH 未更新**：重開終端；`where.exe go`（Windows）。  
- **在錯目錄跑 go**：確認旁邊有 `go.mod`。  

## 延伸閱讀

- <https://go.dev/doc/install>  
- <https://go.dev/doc/tutorial/getting-started>  
