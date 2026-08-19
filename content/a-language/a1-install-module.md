---
lessonId: "A1"
title: "先把 Go 裝好，跑出第一個 module"
description: "安裝官方工具鏈、認識 go.mod，用最小程式確認「能建、能跑」。"
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

## 這章你會搞懂什麼

Go 的專案從 **module（模組）** 開始：旁邊那份 `go.mod` 會宣告「這個專案叫什麼、用哪個語言版本」。之後所有 `import`、下載依賴、編譯，都圍著它轉。

這章你只要練出肌肉記憶：**能建、能跑、能測**。版本細節之後慢慢補就好。

## Python 對照

| Python | Go |
|--------|-----|
| 裝 CPython／用 pyenv | 裝官方 toolchain；必要時設 `GOTOOLCHAIN` |
| `python -m venv` | 通常**不需要** venv；依賴由 module 管 |
| `pip install` + lock 檔 | `go get` + `go.sum`（校驗用） |
| `python main.py` | `go run .` 或 `go build` 產出執行檔 |

重點差異：Go 很常直接編譯成一個（或少數幾個）binary，部署時比較不像「整台機器要有對的直譯器」。

## 怎麼寫

1. 安裝：到 <https://go.dev/dl/>（Windows 用 MSI 就很穩）  
2. **關掉再開**一個終端機，確認指令找得到：

```bash
go version
go env GOROOT GOPATH GOMODCACHE
```

3. 自己建一個最小 module：

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
go test ./...   # 還沒有測試檔也會成功（等於「沒東西可測」）
```

本站現成範例在：`examples/a01-hello-module`。你可以先 `cd` 進去跑 `go run .`。

## 細節

### `go.mod` 在管什麼

- **module path**（例如 `example.com/hello`）是別人（或你自己別的 package）`import` 時的前綴。它**不一定**要等於資料夾名，但團隊最好約定一致，不然會很混亂。  
- **`go 1.22` 那一行**代表語言／標準庫基線。本站教材以 **Go 1.22+** 為準（後面路由模式、`slog` 等會用到）。  
- `go run`：編譯到暫存再執行，適合學習與小工具。  
- `go build`：產出 binary（Windows 上常是 `.exe`）。

### 為什麼要有 module，而不是「一個資料夾隨便跑」

因為依賴與版本要可重現。Python 你可能遇過「我機器跑得動、別人的 venv 不一樣」。Go 把模組邊界寫進 `go.mod`／`go.sum`，就是在減少這種運氣成分。

### 進階可先略過

- `GOTOOLCHAIN=auto` 時，碰到較新語法可能自動拉新版 toolchain。  
- 公司網路常要設 `GOPROXY`、`GOPRIVATE`。

## 遊戲 Server 會用在哪

每個服務（閘道 gateway、房間 worker、假 client）都應該是**清楚的 module**，或 monorepo 裡清楚的 package 樹。本站的 Arena Mini 在 `demo/arena-mini/server`——你之後會一直進出那個目錄。

## 請丟掉的舊習慣

1. 沒有專案邊界就到處用相對路徑亂 import。  
2. 依賴裝在「全域 site-packages」那種心智——在 Go 請看**當前目錄的 `go.mod`**。  
3. 用系統碰巧裝到的版本碰運氣——Go 版本要寫進／對齊 `go.mod`。

## 動手練習

### 必做

1. 本機 `go version` 成功。  
2. 跑通 `examples/a01-hello-module`。  

### 選做

1. `go build -o hello.exe .`，直接雙擊或在終端執行那個 binary。  

## 常見坑

- **裝完 PATH 還沒更新**：重開終端；Windows 可用 `where.exe go` 確認指到哪。  
- **在錯目錄跑 go**：旁邊沒有 `go.mod` 就會怪奇怪奇；先 `cd` 對。  
- **以為一定要 venv**：Go 日常開發通常靠 module，不必先複製一套 Python 流程。
