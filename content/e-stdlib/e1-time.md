---
lessonId: "E1"
title: "time 與計時器"
description: "Time、Duration、Ticker、Timer。"
volume: "e"
order: 1
level: "l2"
status: "ready"
path_required: false
tags: ["time"]
example: ""
prev: "E0"
next: "E2"
---

## 這章你會搞懂什麼

遊戲 Server 幾乎離不開「時間」：多久跑一次 tick、逾時踢人、心跳多久沒回就當斷線。

在 Go 裡，時間相關的核心型別都在 `time` 套件。你會分清：

- **`time.Duration`**：一段「長度」（例如 50 毫秒），不是牆上時鐘  
- **`time.Time`**：某一個時間點  
- **`Ticker`**：固定節奏一直滴答（很適合 tick loop）  
- **`Timer`**：響一次的鬧鐘（很適合單次超時）

## 先跟 Python 對一下

| Python | Go | 注意 |
|--------|-----|------|
| `time.sleep(0.1)` | `time.Sleep(100 * time.Millisecond)` | Go 的單位是型別安全的 `Duration`，不是裸 float 秒 |
| `datetime.datetime.now()` | `time.Now()` | 含時區資訊；比對／序列化時要小心 |
| 自己 `while` + sleep 模擬節拍 | `time.NewTicker(...)` | Ticker 把「節奏」變成 channel 事件 |
| `asyncio.wait_for` / timeout | `Timer` 或更常見 `context.WithTimeout` | 單次超時兩邊都能做；生產碼多半偏 context |

專有名詞先講人話：

- **Tick（滴答）**：固定間隔醒來做一次事，例如每秒 20 次更新房間。  
- **Deadline（截止時間）**：最晚等到哪個時間點，超過就放棄。

## 怎麼寫（能跑的最小例子）

### Sleep：暫停一下

```go
time.Sleep(100 * time.Millisecond)
```

`100 * time.Millisecond` 的型別是 `time.Duration`。底層用**納秒**當基數，所以你可以寫 `50*time.Millisecond`、`2*time.Second`，編譯器會幫你檢查別把「秒」跟「時間點」混在一起。

### Ticker：固定節奏

```go
t := time.NewTicker(50 * time.Millisecond) // 約 20Hz
defer t.Stop()                             // 用完一定要停，否則會漏資源

for i := 0; i < 3; i++ {
	<-t.C // 每次「滴」都會從 channel 收到一個 Time
	fmt.Println("tick", i)
}
```

### Timer：響一次

```go
timer := time.NewTimer(2 * time.Second)
defer timer.Stop()

select {
case <-timer.C:
	fmt.Println("timeout")
case <-done:
	fmt.Println("finished earlier")
}
```

實務上「函式最多跑多久」更常寫成 `context.WithTimeout`（見 E7／C4）；`Timer` 適合你明確在 `select` 裡跟其他事件搶誰先到。

## 為什麼這樣設計／底層在幹嘛

1. **`Duration` 是命名型別，不是 `int`**  
   避免你把「毫秒數字」跟「Unix timestamp」加來加去還編譯過。單位用常數乘：`time.Millisecond`。

2. **Ticker 透過 channel 發事件**  
   剛好跟 `select`、超時、取消同一個併發模型。遊戲房的主迴圈常是：

   ```text
   select {
     case <-ticker.C:   // 該模擬下一幀
     case cmd := <-inbox: // 有玩家指令
     case <-ctx.Done(): // 房間要拆了
   }
   ```

3. **記得 `Stop()`**  
   `NewTicker`／`NewTimer` 會在 runtime 登記計時器。不用了不 `Stop`，在長生命週期服務裡會慢慢漏。`defer t.Stop()` 是好習慣。

4. **進階可先略過：`Reset` 與已觸發的 channel**  
   Timer 觸發後要重用得小心；文件有說明何時該 `Stop` 再 `Reset`。新手宁可新建，也別抄半懂的 Reset 範式。

## 遊戲 Server 會用在哪

- Arena Mini **20Hz** ≈ `50 * time.Millisecond` 的 Ticker。  
- 心跳：每 15 秒送 ping（F7）。  
- 匹配／準備階段倒數：Deadline 或 Timeout context。  
- 日誌時間戳：`time.Now().UTC()`，前後端約定用 UTC，少踩時區坑。

## 請丟掉的舊習慣

1. **用「浮點秒」到處傳**——在 Go 請用 `time.Duration`。  
2. **`for { do(); sleep() }` 卻不管關閉**——改成 Ticker + `ctx.Done()`，關房才乾淨。  
3. **拿牆上時鐘做「間隔」算術卻忽略單調性**——量耗時用 `time.Since(start)`，別自己減兩個 `Now()` 還混時區。

## 動手練習

### 必做

1. 寫一支小程式：`NewTicker(100 * time.Millisecond)`，印 3 次後 `Stop` 結束。  
2. 改成 `select`：除了 tick，再加一個 `time.After(350 * time.Millisecond)` 當總超時，看誰先結束。

### 選做

1. 模擬「20Hz tick」跑 1 秒，統計實際收到幾次（感受排程與負載誤差）。  

## 常見坑

- **忘記 `Stop` Ticker**：房間開開關關時特別痛。  
- **在很緊的迴圈裡每次 `time.After`**：`After` 會產生 Timer；熱路徑用可複用的 `Timer`／Ticker，或算好 deadline。  
- **把 `Time` 當字串亂比**：要比先統一時區，或比 `UnixNano`。  
- **以為 Ticker 保證絕對準時**：負載高時會抖，遊戲邏輯要用「這一 tick 的 dt」或固定步進，不要假設每次剛好 50ms。

## 延伸閱讀

- <https://pkg.go.dev/time>  
