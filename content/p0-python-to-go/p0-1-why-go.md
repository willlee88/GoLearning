---
lessonId: "P0.1"
title: "為什麼是 Go"
description: "理解 Go 的設計目標與取捨，而不是把它當成「比較快的 Python」。"
volume: "p0"
order: 1
level: "l2"
status: "ready"
path_required: true
tags: ["python-bridge", "mindset"]
example: ""
prev: "P0.0"
next: "P0.2"
---

## 本章你會建立的心智模型

Go 不是「語法換皮的 Python」，而是一組**刻意的取捨**：用較少的語言特性換團隊可讀性、用靜態型別與編譯換可預測的重構與部署、用 CSP 風格的併發原語服務**網路服務與長連線**這類工作負載。學 Go 的第一課是接受這些取捨，而不是把舊習慣硬塞進去。

## Python 對照

| 面向 | Python | Go |
|------|--------|-----|
| 執行模型 | 直譯／bytecode，常見 CPython + GIL | 編譯成原生機器碼（含交叉編譯） |
| 型別 | 動態；可選 type hints | 靜態；編譯期檢查 |
| 併發主流敘事 | threading / multiprocessing / asyncio | goroutine + channel（另有 sync） |
| 發佈 | 虛擬環境 + 直譯器/鎖檔 | 單一（或少數）靜態連結 binary 很常見 |
| 語言表面積 | 大、多典範 | 刻意小 |

兩者都能做後端；Go 特別常被選在：**高併發連線、運維簡單、需要明確效能與記憶體行為** 的服務——遊戲 Server 正好落在這個區間。

## L1 能用

你現在只需要記住三件事：

1. 原始碼 → `go build` → 可執行檔（Windows 上是 `.exe`）。
2. 一個目錄通常是一個 **package**；`main` package 的 `func main()` 是程式入口。
3. 工具鏈內建：`fmt`（格式化）、`test`（測試）、`mod`（依賴）、`vet`（靜態檢查）。

最小程式：

```go
package main

import "fmt"

func main() {
	fmt.Println("hello from GoLearning")
}
```

```text
go run .
```

## L2 機制

### 設計目標（簡化但實用的版本）

Go 的核心承諾大致是：

- **可讀性優先於表達力**：少魔法，多明確。
- **組合優於繼承**：用小介面與嵌入，而不是深繼承樹。
- **併發是一等公民**：語法與 runtime 一起服務大量 goroutine。
- **工程預設齊全**：格式化、測試、benchmark、profiling、modules 是文化而非外掛。

### 取捨（學深度必看）

| 你得到 | 你失去／變難 |
|--------|----------------|
| 編譯期抓大量錯誤 | 沒有 Python 那種「先跑再說」的彈性 |
| 部署簡單 | 泛型與抽象要克制；不是每種 DSL 都好寫 |
| 併發模型清晰（若用對） | 用錯 channel/鎖一樣會 race、洩漏 |
| 標準庫務實 | 生態不像 Python 那樣「每個領域一個巨框架」 |

### 和遊戲 Server 的契合點

- **大量連線**：每個連線一個（或一組）goroutine 的模型直覺。
- **長生命週期服務**：明確的 `main`、設定、優雅關閉文化。
- **熱路徑可控**：profile、benchmark、避免隱藏分配成為日常。
- **規則與 I/O 分離**：靜態型別讓「純邏輯套件」好測。

## L3 深潛（可選）

- 閱讀 [Go at Google: Language Design in the Service of Software Engineering](https://go.dev/talks/2012/splash.article)（歷史與動機）。
- 比較：Rust（更強所有權與零成本抽象）、Java/C#（重量級 runtime 與生態）、Node（單執行緒事件迴圈文化）。重點不是誰贏，而是**問題形狀**是否匹配。

## 請丟掉的 Python 習慣

1. **「先寫再說，型別以後補」** — 在 Go 裡型別是設計的一部分，不是裝飾。
2. **「框架幫我搞定一切」** — 先懂 `net/http`、標準庫，再選框架。
3. **「執行緒 = 併發的唯一心智」** — goroutine 更輕，但**不是**免費的；仍要思考所有權與生命週期。
4. **「例外跳轉當控制流」** — 錯誤是值；控制流應看得見。

## 遊戲 Server 連結

遊戲後端典型壓力是：**高連線數、狀態一致性、低延遲廣播、可觀測、可關機**。Go 的工具鏈與併發原語讓這些成為「工程問題」而非「語言打架」。後續卷會反覆回到：Session、Room、Tick、廣播扇出。

## 練習

### 必做

1. 安裝 Go（若尚未安裝），執行 `go version`。
2. 在任意目錄建立 `hello.go`，`go run .` 成功印出字串。
3. 用三句話寫下：你選 Go 學遊戲 Server 的理由（可以很務實，例如「好部署」）。

### 選做

1. 讀完 Effective Go 的前言與「Formatting」小節，對照 `gofmt` 不可爭論的格式化文化。
2. 列出你過去 Python 專案裡「部署最痛」的三點，猜測 Go 會如何改變它們。

## 常見坑與如何看見

- **裝了 Go 但 PATH 沒更新**：新開終端機再跑 `go version`。
- **把 Go 當腳本語言亂放檔案**：從一開始用 module：`go mod init example.com/hello`。
- **版本過舊**：遊戲與現代泛型、`slog` 等內容以 **Go 1.22+** 為準。

## 延伸閱讀

- <https://go.dev/doc/>
- <https://go.dev/doc/effective_go>
- 本站規劃書：`docs/規劃書.md` §1–§5
