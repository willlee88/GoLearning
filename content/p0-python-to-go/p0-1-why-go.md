---
lessonId: "P0.1"
title: "為什麼學 Go（而不是把 Python 換皮）？"
description: "搞懂 Go 刻意做了哪些取捨：編譯、靜態型別、併發、好部署——以及什麼時候它特別適合遊戲 Server。"
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

## 這章你會搞懂什麼

Go **不是**「語法長得不一樣的 Python」。

它是一組**刻意的取捨**：語言特性少一點，換團隊好讀；編譯期把型別抓死，換重構與部署比較可預測；內建「一次做很多事」的原語（concurrency），服務網路服務與長連線這類工作。

學 Go 的第一課是：**接受這些取捨**，而不是把舊習慣硬塞進去。讀完你要能講出「我得到什麼、失去什麼」，以及為什麼遊戲 Server 常落在 Go 的舒適區。

## 先跟 Python 對一下

| 面向 | Python | Go |
|------|--------|-----|
| 怎麼跑 | 直譯／bytecode，常見 CPython + GIL（全域直譯器鎖） | 編譯成原生機器碼，也常做交叉編譯 |
| 型別 | 動態；可選 type hints | 靜態；編譯期就檢查 |
| 併發常見說法 | threading / multiprocessing / asyncio | goroutine + channel（另有 sync 套件） |
| 怎麼發佈 | 虛擬環境 + 直譯器／鎖檔 | 常打成單一（或少數）靜態連結的執行檔 |
| 語言有多大 | 大、多種寫法都通 | 刻意保持小 |

兩者都能做後端。Go 特別常被選在：**高連線數、運維想簡單、需要比較明確的效能與記憶體行為**——遊戲 Server 正好落在這個區間。

## 怎麼寫（能跑的最小例子）

你現在先記住三件事就好：

1. 原始碼 → `go build` → 可執行檔（Windows 上常是 `.exe`）。  
2. 一個目錄通常是一個 **package**（套件）；`main` package 裡的 `func main()` 是程式入口。  
3. 工具鏈內建：`fmt`（印東西／格式化）、`test`（測試）、`mod`（依賴）、`vet`（靜態檢查）。  

最小程式：

```go
package main

import "fmt"

func main() {
	fmt.Println("hello from GoLearning")
}
```

在該目錄執行：

```text
go run .
```

若還沒有 `go.mod`，先 `go mod init example.com/hello`，再跑——從第一天就養成「專案有邊界」的習慣。

## 為什麼這樣設計／底層在幹嘛

### Go 大概在承諾什麼

- **可讀性優先於花招**：少魔法，多寫清楚。別人半夜改你的 Server 才不會哭。  
- **組合優於繼承**：用小介面（interface）跟嵌入（embedding），而不是深層繼承樹。  
- **併發是一等公民**：語法跟 runtime 一起服務大量輕量執行單元（goroutine）。  
- **工程預設就齊**：格式化、測試、benchmark、profiling、modules 是文化，不是外掛。  

### 你得到什麼、變難什麼

| 你得到 | 你失去／變難 |
|--------|----------------|
| 編譯期抓掉大量錯誤 | 沒有 Python 那種「先跑再說」的彈性 |
| 部署相對單純 | 抽象要克制；不是每種 DSL 都好寫 |
| 併發模型清楚（用對的話） | channel／鎖用錯一樣會 race、洩漏 |
| 標準庫務實好用 | 生態不像 Python「每個領域一個巨框架」 |

### 和遊戲 Server 為什麼對得上

- **大量連線**：每個連線配一個（或一組）goroutine，心智上直覺。  
- **長生命週期服務**：明確的 `main`、設定、優雅關閉（graceful shutdown）文化。  
- **熱路徑可控**：profile、benchmark、少隱藏配置成為日常。  
- **規則與 I/O 分開**：靜態型別讓「純邏輯套件」好測。  

寫錯會怎樣？把 Go 當「比較快的腳本」亂寫、到處共享 map 又不加同步——線上會出現難重現的 race，關機也關不乾淨。後面章會教你怎麼避免。

### 進階可先略過

- 讀 [Go at Google: Language Design in the Service of Software Engineering](https://go.dev/talks/2012/splash.article)（歷史與動機）。  
- 跟 Rust（所有權更嚴）、Java/C#（runtime 與生態更重）、Node（單執行緒事件迴圈文化）比時，重點不是誰贏，而是**問題形狀**是否匹配。  

## 遊戲 Server 會用在哪

遊戲後端典型壓力是：**高連線數、狀態要一致、低延遲廣播、可觀測、能乾淨關機**。

Go 的工具鏈與併發原語，讓這些比較像「工程問題」而不是「語言一直打架」。後面卷會反覆回到：Session（連線工作階段）、Room（房間）、Tick（固定節奏更新）、廣播扇出。

## 請丟掉的舊習慣

1. **「先寫再說，型別以後補」** — 在 Go 裡型別是設計的一部分，不是裝飾。  
2. **「框架幫我搞定一切」** — 先懂 `net/http`、標準庫，再選框架。  
3. **「執行緒 = 併發的唯一心智」** — goroutine 更輕，但**不是**免費的；仍要思考誰擁有資料、何時結束。  
4. **「例外跳轉當控制流」** — 錯誤是值；控制流應該看得見（下一章會講專案形態，錯誤章在 P0.4）。  

## 動手練習

### 必做

1. 安裝 Go（若尚未安裝），執行 `go version`。本站以 **Go 1.22+** 為準。  
2. 在任意目錄建立小專案：`go mod init` + `hello.go`，`go run .` 成功印出字串。  
3. 用三句話寫下：你選 Go 學遊戲 Server 的理由（可以很務實，例如「好部署」「連線多」）。  

### 選做

1. 讀 Effective Go 的前言與「Formatting」小節，感受 `gofmt`「格式化不可爭論」的文化。  
2. 列出你過去 Python 專案裡「部署最痛」的三點，猜 Go 會怎麼改變它們。  

## 常見坑

- **裝了 Go 但 PATH 沒更新**：新開終端機再跑 `go version`。  
- **把 Go 當腳本亂放檔案**：從一開始用 module：`go mod init example.com/hello`。  
- **版本過舊**：泛型、`slog` 等後續內容以 **Go 1.22+** 為準；太舊會跟教材對不上。  

## 延伸閱讀

- <https://go.dev/doc/>  
- <https://go.dev/doc/effective_go>  
- 本站規劃書：`docs/規劃書.md` §1–§5  
