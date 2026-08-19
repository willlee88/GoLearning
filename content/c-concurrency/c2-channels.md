---
lessonId: "C2"
title: "channel：通訊、同步，以及誰能 close"
description: "無緩衝是交接、有緩衝是佇列；搞清楚擁有者，才不會對已關閉的 channel 傳送而 panic。"
volume: "c"
order: 2
level: "l2"
status: "ready"
path_required: true
tags: ["channel"]
example: "examples/c02-channels"
prev: "C1"
next: "C3"
---

## 這章你會搞懂什麼

**Channel** 是 Go 裡很有名的通訊原語（常跟 CSP 風格放在一起談）：用傳遞訊息來同步，而不只是共享記憶體。

兩種你一定要分清：

- **無緩衝** `make(chan T)`：送與收必須「握手」——一邊等到另一邊，才完成交接  
- **有緩衝** `make(chan T, N)`：中間有佇列；沒滿可先送、沒空可先收  

還有一條鐵律：**誰寫、誰 close，要有清楚擁有者**。對已經關閉的 channel 再傳送 → **panic**。

## Python 對照

| Python | Go |
|--------|-----|
| `queue.Queue` | 有緩衝 channel 近似（但仍有 close／select 語意） |
| `asyncio.Queue` | 同上，再加上取消通常靠 context |
| 條件變數通知 | 常用 channel 當訊號（空 struct 訊號很常見） |
| 多執行緒共享 list＋鎖 | 也可以；channel 不是唯一解 |

channel 很強，但不是萬靈丹。共享結構 + `Mutex` 有時更清楚（C5）。選工具看「資料誰擁有」。

## 怎麼寫

```go
ch := make(chan int)      // 無緩衝
ch := make(chan int, 16)  // 有緩衝

ch <- 1       // 送
v := <-ch     // 收
v, ok := <-ch // ok == false：已關閉且空了

close(ch)
```

`ok == false` 只代表「關閉且讀完」，不要跟「讀到零值」搞混——零值也可能是合法資料，所以需要 `ok` 或明確協定。

範例：`examples/c02-channels`。

## 細節

### 為什麼無緩衝叫「同步交接」？

因為沒有中間盒子可暫放。送方卡到收方出現（或反向）。  
這很好用在：確保對方「真的拿到了」才繼續；或拿 channel 當事件通知。

有緩衝則像有限長度的佇列：解耦速度，但**緩衝大小就是你願意積壓多少**。無界心智很危險（C8）。

### close 的規則（為什麼容易 panic）

- **只由傳送側（或約定的唯一擁有者）close**  
- 關閉後：還可以繼續接收剩餘值；收完後 `ok=false`  
- **對已關閉 channel 傳送 → panic**  
- 重複 close → panic  
- 接收關閉且空的 channel → 立刻拿到零值 + `ok=false`（不阻塞）

實務上常「生產者負責 close」，消費者用 `for v := range ch` 讀到關閉為止。

### `nil` channel 的怪癖（超有用）

在 **nil channel** 上送或收會**永遠阻塞**。  
在 `select` 裡，把某個 case 的 channel 設成 nil，等於**暫時禁用那條分支**——這是進階技巧，房間狀態機有時會用到。

### 進階可先略過

- 規格裡 channel 的 happens-before：送入被接收，同步記憶體可見性。  
- 方向限制：`chan<- T` 只能送、`<-chan T` 只能收，適合函式簽名表達擁有權。

## 遊戲 Server 會用在哪

| 用途 | 常見做法 |
|------|----------|
| 玩家輸入 | 有界 `chan Command` |
| 房間事件 | room 擁有的 inbox，只有 room loop 讀 |
| 關房 | close inbox，或改用 context cancel（C4） |
| 廣播任務 | 小心：別對每個 peer 無界塞 |

原則：**房間狀態由單一 goroutine 擁有**時，外面只准「丟訊息進 inbox」，不准大家一起改 map。

## 請丟掉的舊習慣

1. 好幾個地方都能 `close`——最後一定有人踩雷。  
2. 無緩衝 + 狂 `go` 送，或超大／「先開一百萬」當無界佇列。  
3. 「全部同步都改成 channel」——有時一把鎖更可读。

## 動手練習

### 必做

1. 跑 `examples/c02-channels`。  
2. 單獨寫個小程式：對已關閉 channel 傳送，親眼看 panic。  

### 選做

1. 用 `for v := range ch` 改寫消費者；確認生產者結束時有 close。  

## 常見坑

- **死人發送／沒人接收**：goroutine 洩漏（C1）。  
- **close 後還有人送**：panic；用擁有者規則＋碼審避免。  
- **用緩衝大小「賭」不會滿**：滿了還是會阻塞——要麼背壓策略，要麼別裝死（C8）。

## 延伸閱讀

- <https://go.dev/ref/spec#Channel_types>  
