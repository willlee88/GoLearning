---
lessonId: "P0.6"
title: "併發預告"
description: "先建立正確直覺：goroutine 不是 thread，也不是 asyncio 的直接翻譯。"
volume: "p0"
order: 6
level: "l2"
status: "ready"
path_required: true
tags: ["python-bridge", "concurrency"]
example: ""
prev: "P0.5"
next: "P0.7"
---

## 本章你會建立的心智模型

Go 的併發三件套是 **goroutine（執行）**、**channel（通訊）**、**select（多路等待）**，再加上 `sync` 與 `context`。你可以先把「每個連線一個讀寫迴圈」當成遊戲 Server 的預設畫面，但**不要**在搞懂生命週期與所有權前狂開 goroutine。

## Python 對照

| Python | Go（粗對照，勿硬套） |
|--------|----------------------|
| `threading.Thread` | goroutine（更輕，但模型不同） |
| `asyncio` 協程 + event loop | 多 goroutine；runtime 排程 |
| `queue.Queue` | channel |
| `async with` 取消 | `context.Context` |
| GIL 限制 CPU 平行 | 可平行；仍要注意共享記憶體 |

## L1 能用

```go
go func() {
	fmt.Println("runs concurrently")
}()
```

```go
ch := make(chan string, 1)
ch <- "ping"
msg := <-ch
```

## L2 機制

- Goroutine 由 Go runtime **M:N 排程**到作業系統執行緒。
- Channel 傳送有**同步語意**（無緩衝時交接）；關閉 channel 有擁有權規則。
- 共享記憶體靠同步（鎖或 CSP）；**data race 是 bug**，用 `go test -race` 抓。

詳細在 **卷 C**。本章只建立：併發是設計問題，不是加 `go` 關鍵字就結束。

## L3 深潛（可選）

- 為何「不要通過共享記憶體來通訊；通過通訊來共享記憶體」是口號而非教條——鎖在熱路徑仍然合理。

## 請丟掉的 Python 習慣

1. 把 `asyncio.create_task` 的心智原封貼成到處 `go func`。
2. 假設「不會平行就不會 race」——在 Go 預設就可能平行。
3. 開了背景工作卻不定義如何停止（洩漏）。

## 遊戲 Server 連結

典型拆分：

- 每連線：讀取 goroutine + 寫入串行化  
- 每房間：一個邏輯擁有者（goroutine 或鎖保護的 tick）  
- 大廳：registry 的併發存取策略  

## 練習

### 必做

1. 寫一個程式：主 goroutine 送 3 個訊息到 channel，另一個印出後結束（注意如何讓 main 等待）。
2. 用文字描述：若兩個 goroutine 無鎖寫同一個 `map` 會怎樣。

### 選做

1. 跑 `go test -race` 於任何含併發的小實驗。

## 常見坑與如何看見

- main 結束殺光所有 goroutine：用 `WaitGroup` 或 channel 會合。
- 無界 `go` + 無界工作佇列：記憶體爆炸——後續講背壓。

## 延伸閱讀

- <https://go.dev/blog/pipelines>
- 卷 C 全卷
