---
lessonId: "C7"
title: "errgroup：一組任務，錯誤與取消一起收"
description: "並行做幾件事，任一失敗就取消其他；比手寫 WaitGroup + error 鎖少踩坑。"
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

## 這章你會搞懂什麼

你常會碰到：「同時做 A、B、C，任一失敗就別做了，順便把錯誤帶回主線。」  
手寫的話要：`WaitGroup` + 保護 error 的 Mutex + 自己 cancel context——很容易漏。

`golang.org/x/sync/errgroup` 把這套打包好：`Go` 丢任務，`Wait` 等結束並回第一個錯誤；`WithContext` 還能在首錯時取消其他任務。

它不是魔法，但能少寫一堆易錯膠水碼。

## Python 對照

| Python | Go |
|--------|-----|
| `asyncio.TaskGroup` | `errgroup.Group`／`WithContext` |
| `concurrent.futures.wait` | 手動或 errgroup |
| 一個 task 炸、其他還在跑 | `WithContext` 可把取消傳下去（任務要肯理 ctx） |

## 怎麼寫

```go
g, ctx := errgroup.WithContext(context.Background())
g.Go(func() error { return taskA(ctx) })
g.Go(func() error { return taskB(ctx) })
err := g.Wait()
```

注意：`Go` 裡的函式要回 `error`；成功回 `nil`。  
子任務必須**使用傳入的 `ctx`**，取消才有意義——否則只是「其他人停了，你還在埋頭做」。

範例：`examples/c07-errgroup`。

## 細節

### `WithContext` 做了什麼？

第一個回傳非 nil error 的任務，會觸發派生 context 的 cancel。  
其他任務若卡在 `select`／IO 且有聽 `ctx.Done()`，就能提早收工。`Wait` 仍會等所有已啟動任務結束（或結束路徑跑完）。

### 併發上限

errgroup 預設不會幫你限流。開 10000 個 `Go` 仍可能把自己打爆。  
需要限流時：加 semaphore（常見寫法是有緩衝 channel 當權證），或分批。

### 跟 panic 的關係

`Go` 裡若 panic 且沒 recover，一樣可能幹掉行程（B3）。  
errgroup 管的是 **error 值**，不是 panic 防護罩。

### 進階可先略過

- `SetLimit`（較新版本）可限制同時執行數。  
- 只要「等全部完成、收集多個錯誤」而不是「首錯取消」，可能要別的結構（自己累積 errors）。

## 遊戲 Server 會用在哪

適合：

- 啟動時並行拉設定／驗票／寫審計  
- 一次請求內扇出幾個獨立 IO  

不適合：

- 單房 tick 裡無界扇出「每個玩家一個 goroutine 做重計算」——狀態與背壓都會變難  

線上熱路徑先問：有沒有上限？失敗要不要取消兄弟任務？

## 請丟掉的舊習慣

1. 開一堆 task 不管取消。  
2. 只隨手拿「最後一個 error」，前面的失敗蒸發。  
3. 以為用了 errgroup 就自動限流。

## 動手練習

### 必做

1. 跑 `examples/c07-errgroup`。  
2. 讓 taskB 失敗，觀察 taskA 是否因 ctx 取消而提早結束（在 A 裡打日誌最清楚）。  

### 選做

1. 幫扇出加上「最多同時 4 個」的 semaphore。  

## 常見坑

- **子任務忽略 ctx**：表面上用了 `WithContext`，實則取消失靈。  
- **在 `Go` 裡改共享狀態卻無同步**：errgroup 不替你解決 race。  
- **把業務 panic 當 error**：先轉成 error，或在邊界 recover。

## 延伸閱讀

- <https://pkg.go.dev/golang.org/x/sync/errgroup>  
