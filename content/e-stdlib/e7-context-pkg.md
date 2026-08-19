---
lessonId: "E7"
title: "context 套件速查"
description: "Background、TODO、With* 家族與傳值紀律。"
volume: "e"
order: 7
level: "l2"
status: "ready"
path_required: false
tags: ["context", "stdlib"]
example: ""
prev: "E6"
next: "E8"
---

## 這章你會搞懂什麼

C4 已經用「取消樹」講過 **`context.Context`** 的心智。本章當**標準庫速查表**：有哪些建構函式、怎麼傳、值能不能亂塞。

人話：`context` 用來帶三類東西——

1. **取消信號**：上司說停，下屬都該停  
2. **截止時間**：最晚做到什麼時候  
3. **請求級 metadata**（克制）：例如 trace id  

慣例：函式第一個參數叫 `ctx`，一路往下傳。

## 先跟 Python 對一下

| Python | Go | 注意 |
|--------|-----|------|
| `asyncio.timeout`／`CancelledError` | `ctx.Done()`／`ctx.Err()` | 都要在可取消點檢查 |
| 請求級 thread-local／contextvar | **顯式**傳 `ctx` | Go 刻意不藏在隱式全域 |
| 沒有統一取消時自己設 `stopping` flag | `WithCancel` 樹狀擴散 | 關服特別好用 |

## 怎麼寫（能跑的最小例子）

```go
ctx := context.Background() // 根；通常在 main／測試頂層
ctx, cancel := context.WithTimeout(ctx, time.Second)
defer cancel() // 釋放計時器資源；別省

select {
case <-ctx.Done():
	return ctx.Err() // context.DeadlineExceeded 或 Canceled
case res := <-work:
	return res, nil
}
```

100ms 超時示範：

```go
ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
defer cancel()

select {
case <-time.After(500 * time.Millisecond):
	fmt.Println("work finished") // 正常不會走到
case <-ctx.Done():
	fmt.Println("gave up:", ctx.Err())
}
```

## 為什麼這樣設計／底層在幹嘛

### 建構家族

| 函式 | 人話 |
|------|------|
| `Background()` | 永遠的根，不會被取消 |
| `TODO()` | 「我暫時不知道該傳哪個 ctx」的佔位；別留到生產亂用 |
| `WithCancel` | 手動 `cancel()` 通知下游 |
| `WithTimeout`／`WithDeadline` | 時間到自動取消 |
| `WithValue` | 塞稀疏 metadata；**不是**可選參數袋 |

### 傳遞紀律

- **往下傳，不要塞進長命 struct 當「全域神器」**（除非你很清楚生命週期）。  
- 衍生：`child, cancel := context.WithCancel(parent)`——parent 取消，child 也取消。  
- 慢迴圈／每個請求步驟主動看 `ctx.Err()` 或 `<-ctx.Done()`，否則取消了還在傻忙。

### `WithValue` 為什麼常被罵

它方便，也容易變垃圾堆：把 DB 連線、使用者整包物件塞進去，呼叫端看不出依賴。社群共識是：**只放跨層追蹤用的請求級資料**（request id、log 欄位），業務依賴仍用參數傳。

## 遊戲 Server 會用在哪

- HTTP：`r.Context()`——客戶端斷線，handler 可感知。  
- 每個玩家連線一個 ctx；連線斷 → cancel → 讀寫 loop 收斂。  
- 房間 ctx：房主解散／局結束。  
- 關服：root cancel → 擴散到匹配、房間、DB 查詢。  
- 對外呼叫（HTTP client、DB）都把 ctx 傳下去，才有統一超時。

## 請丟掉的舊習慣

1. **全域 `var stopping bool`**——改 cancel 樹。  
2. **把可選參數、依賴注入塞进 `WithValue`**。  
3. **開了 `WithCancel`／`WithTimeout` 卻不 `defer cancel()`**——Timer 與子 goroutine 可能洩漏。  
4. **在函式庫深處呼叫 `Background()` 吃掉上層超時**——等於無視呼叫者的 deadline。

## 動手練習

### 必做

1. 寫一個 `select`：`WithTimeout(..., 100*time.Millisecond)`，另一邊是 `time.After(1*time.Second)`，確認會走超時。  
2. 複習 C4／`examples/c04-context`（若已做過可跳過）。  

### 選做

1. 父 `WithCancel`，子 `WithTimeout`；先 `cancel` 父，觀察子是否也 Done。  

## 常見坑

- **cancel 後還假設工作一定立刻停**：必須合作式檢查；狂佔 CPU 的迴圈不會自動中斷。  
- **把 ctx 存起來跨請求重用**：請求結束後 ctx 已死，下一個請求要新的。  
- **只 `cancel` 不關連線**：取消是信號，socket 仍要自己 `Close`。

## 延伸閱讀

- <https://go.dev/blog/context>  
- <https://pkg.go.dev/context>  
