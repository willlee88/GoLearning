---
lessonId: "C1"
title: "goroutine 生命週期"
description: "go 關鍵字、會合、洩漏模式。"
volume: "c"
order: 1
level: "l2"
status: "ready"
path_required: true
tags: ["goroutine"]
example: "examples/c01-goroutines"
prev: "C0"
next: "C2"
---

## 本章你會建立的心智模型

`go f()` 啟動一個並發執行單元。它很便宜，但不是免費：每個未結束的 goroutine 佔用記憶體與排程資源。**main 結束會砍掉一切**；背景工作必須定義停止條件與會合方式（`WaitGroup` / channel）。

## Python 對照

| Python | Go |
|--------|-----|
| `threading.Thread` | goroutine |
| `asyncio.create_task` | `go`（模型不同） |
| daemon thread | 無直接等價；靠退出與取消 |

## L1 能用

```go
var wg sync.WaitGroup
wg.Add(1)
go func() {
	defer wg.Done()
	work()
}()
wg.Wait()
```

範例：`examples/c01-goroutines`。

## L2 機制

常見洩漏：

- 阻塞在無人寫的 channel  
- 忘了 `Done`  
- 沒人監聽的的出口  

## 請丟掉的 Python 習慣

1. fire-and-forget 不管生命週期。  
2. 以為「行程還在就代表工作還健康」。  

## 遊戲 Server 連結

每連線讀 loop 一個 goroutine；斷線必須退出 loop 並 unregister。

## 練習

### 必做

1. 跑 `examples/c01-goroutines`。  
2. 去掉 `Wait` 觀察提早退出。  

## 延伸閱讀

- <https://go.dev/blog/pipelines>  
