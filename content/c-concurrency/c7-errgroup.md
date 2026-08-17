---
lessonId: "C7"
title: "errgroup 與扇出"
description: "一組 goroutine 的錯誤與取消。"
volume: "c"
order: 7
level: "l2"
status: "ready"
path_required: false
tags: ["errgroup", "concurrency"]
example: "examples/c07-errgroup"
prev: "C6"
next: "C8"
---

## 本章你會建立的心智模型

`golang.org/x/sync/errgroup` 適合「並行做幾件事，任一錯就取消其他」。比手寫 WaitGroup+error mutex 少錯。

## Python 對照

| Python | Go |
|--------|-----|
| `asyncio.TaskGroup` | `errgroup.Group` |
| `concurrent.futures` | 手動或 errgroup |

## L1 能用

```go
g, ctx := errgroup.WithContext(context.Background())
g.Go(func() error { return taskA(ctx) })
g.Go(func() error { return taskB(ctx) })
err := g.Wait()
```

範例：`examples/c07-errgroup`。

## L2 機制

- `WithContext` 在第一個錯誤時 cancel。  
- 不要在 `Go` 裡 panicking 不 recover。  
- 限制併發：semaphore（加權 channel）。  

## 請丟掉的 Python 習慣

1. 開一堆 task 不管取消。  
2. 只拿最後一個 error。  

## 遊戲 Server 連結

並行拉設定／驗票／寫審計；**不要**在單房 tick 裡亂扇出無界。

## 練習

### 必做

1. 跑 `examples/c07-errgroup`。  
2. 讓 taskB 失敗，觀察 taskA 是否提早結束。  

## 延伸閱讀

- `x/sync/errgroup`  
