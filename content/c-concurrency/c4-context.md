---
lessonId: "C4"
title: "context：取消是一棵樹"
description: "WithCancel／Timeout 怎麼往下傳；Done／Err 怎麼用；Value 為什麼要克制。"
volume: "c"
order: 4
level: "l2"
status: "ready"
path_required: true
tags: ["context"]
example: "examples/c04-context"
prev: "C3"
next: "C5"
---

## 這章你會搞懂什麼

`context.Context` 承載三類東西（實務上你主要用前兩個）：

1. **取消信號** — 「別做了」  
2. **截止時間** — 「最晚到這裡」  
3. **請求範圍的值** — 極克制，例如 trace id  

取消是一棵**樹**：父取消，子也會取消。所以關服時 cancel root，連線、房間、進行中的匹配可以一起收到信號——前提是你真的把 `ctx` 傳下去，並在迴圈裡檢查。

慣例：函式第一個參數叫 `ctx`，**顯式傳遞**。別學 thread-local 那種隱形通道。

## Python 對照

| Python | Go |
|--------|-----|
| `asyncio.timeout`／`CancelledError` | `ctx.Done()`／`ctx.Err()` |
| 手動 `Event`／全域 `stopping=True` | 往下傳的 context 樹 |
| 請求級 thread-local／contextvars | 顯式 `ctx` 參數（更囉唆，也更清楚） |

囉嗦是代價；換來的是：測試好塞 fake deadline，也不會有「不知道誰設的全域旗標」。

## 怎麼寫

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel() // 重要：盡早釋放 timer／通知子樹

select {
case <-ctx.Done():
	return ctx.Err()
case res := <-work:
	return res, nil
}
```

衍生方式你會一直看到：

- `WithCancel`：手動 cancel  
- `WithTimeout`／`WithDeadline`：時間到自動 cancel  
- `WithValue`：塞請求級 metadata（慎用）

範例：`examples/c04-context`。

## 細節

### 為什麼一定要 `defer cancel()`？

即使 timeout 會自動觸發，`cancel` 仍應呼叫：釋放資源、讓子 context 儘快結束。  
漏 cancel 在大量短請求下會累積 timer／記憶體壓力。

### 誰該檢查 `ctx`？

- 慢迴圈、poll、長時間 wait 的入口  
- 會阻塞在 channel／IO 的地方：優先選「可被取消的等待」（select 加上 `<-ctx.Done()`）  
- 純 CPU 小函式未必每次都查，但長任務要查  

取消是**協作式**的：沒人看 `Done`，就誰也不會停。

### `WithValue` 為什麼要克制？

它看起來像「隱式參數袋」。塞可選參數、塞依賴、塞一整個 User 物件——最後 API 會變得不可讀、也難靜態檢查。  

建議：只放跨層都需要的請求級 metadata（trace id、request id）。業務依賴仍靠參數／struct 注入（B4）。

### 進階可先略過

- context 不該塞進長壽命 struct 亂存（容易生命周期錯亂）；存 cancel 函式或派生規則時要非常清醒。  
- `Err()` 常見 `Canceled`、`DeadlineExceeded`——上層可翻譯成協定碼。

## 遊戲 Server 會用在哪

常見一棵樹：

```text
process root ctx
  ├─ 連線 ctx（斷線／踢人就 cancel）
  ├─ 房間 ctx（關房）
  └─ 單次匹配／單一 RPC 的 timeout ctx
```

關服：root cancel → 連線停止讀寫 → 房間迴圈從 `select` 離開 → WaitGroup 會合 → 行程退出。  
這條路比到處設 `stopping=True` 好推理得多。

## 請丟掉的舊習慣

1. 全域旗標 `stopping=True` 當取消協定。  
2. 把 context 塞進無所不在的單例。  
3. 開了 `WithTimeout` 卻不 `cancel`，也不檢查 `Done`。

## 動手練習

### 必做

1. 跑 `examples/c04-context`。  
2. 用 timeout 包一個「最多等 200ms 的假 IO」；超時應回 `context.DeadlineExceeded`。  

### 選做

1. 父 `WithCancel`，開兩個子工作聽 `Done`；父 cancel 後確認兩邊都結束。  

## 常見坑

- **把 `context.Background()` 往業務深處亂傳當預設**：邊界用 Background／TODO，內部應用呼叫端傳來的 ctx。  
- **超時後仍繼續寫共享狀態**：聽到取消要停手，別寫到一半當沒看到。  
- **用 Value 傳必填依賴**：改回函式參數。

## 延伸閱讀

- <https://go.dev/blog/context>  
