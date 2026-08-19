---
lessonId: "H5"
title: "Tick：固定心跳推進模擬"
description: "用 fixed timestep 累積輸入、推進狀態、再廣播。time.Ticker + select，而不是收到一包就改世界。"
volume: "h"
order: 5
level: "l2"
status: "ready"
path_required: true
tags: ["tick"]
example: "examples/h05-tick"
prev: "H4"
next: "H6"
---

## 這章你會搞懂什麼

**Tick** 就是 Server 的心跳：每隔固定時間（例如 50ms → 20Hz）做一輪：

1. 把這段時間累積的輸入拿出來  
2. 依規則推進模擬（移動、碰撞、得分…）  
3. （可選）廣播快照給房裡的人  

為什麼不「收到一包立刻改、立刻廣播」？

- 節奏不穩：網路抖、客戶端狂送，世界推進忽快忽慢  
- 難重現：同樣操作序列，因到達時間不同結果不同，除錯像抓煙  
- 難調參：速度、冷卻、物理步長沒有共同時鐘  

這章你要會用 `time.Ticker` + `select` 寫出最小 tick loop，並知道輸入是「先排隊、tick 時再套用」。

## Python 對照

| Python | Go |
|--------|-----|
| `asyncio` 裡 `await asyncio.sleep(dt)` 迴圈 | `time.NewTicker` + `for { select { ... } }` |
| 每次用真實經過時間當 `dt` 做積分 | 教學與多數玩法先用**固定 dt**（fixed timestep） |
| 在每個 socket handler 直接改全局狀態 | handler 只 enqueue；tick 才 apply |

變步長（variable timestep）不是不能用，但新手很容易把「幀率不穩」寫進遊戲規則裡。先固定步長，人生較幸福。

## 怎麼寫（能跑的最小例子）

範例：`examples/h05-tick`。概念版：

```go
ticker := time.NewTicker(50 * time.Millisecond)
defer ticker.Stop()

var buf []Cmd
for {
	select {
	case cmd := <-inbox:
		buf = append(buf, cmd) // 只累積，不立刻當物理世界
	case <-ticker.C:
		applyAll(buf)
		buf = buf[:0]
		broadcast(snapshot())
	case <-ctx.Done():
		return
	}
}
```

自己跑看看：

```bash
cd examples/h05-tick
go run .
```

你會看到：輸入在 tick 之間進來，真正 `applied` 印在 tick 當下。試著改成 `100 * time.Millisecond`（10Hz），感受訊息節奏變慢。

### 寫的時候記得

- **`defer ticker.Stop()`**：房間結束要停，避免洩漏。  
- **`ctx.Done()`**：取消信號進來要離開迴圈（H2／C4 的生命週期）。  
- **inbox 最好有界**：滿了怎麼辦要有政策（丟指令、斷線、回壓）——C8 背壓同一類問題。

## 為什麼是 fixed timestep

### 可重現、好調參

規則若假設「每 tick 移動 `speed` 像素」，那 20Hz 與 10Hz 的手感差異是**可預期的參數問題**，不是「看誰 CPU 忙」。  
測試也可以：塞一串 input → 呼叫 N 次 `Tick()` → 斷言座標——不必真的睡 50ms。

### 輸入與模擬解耦

```text
WS 讀 loop  ：高速、不規則到達
inbox       ：緩衝
Tick        ：固定節奏消費
broadcast   ：可以跟模擬同頻，或較低頻
```

讀 loop 裡做重物理，會讓慢客戶端拖住讀取，也讓鎖與延遲更難估。

### 廣播頻率可以 ≤ 模擬頻率

例如模擬 30Hz、廣播 15Hz：省頻寬，中間幀靠客戶端插值顯示。  
教學 demo（Arena Mini）常直接 20Hz 模擬 + 20Hz 全量廣播——先求正確。

### 進階可先略過

- **餘數累積**（spiral of death 相關）：真實時間落下很多時，一次追很多 tick，要有上限，避免越追越忙。經典文章：*Fix Your Timestep*（Gaffer on Games）。  
- 模擬 tick 與渲染幀分離——Server 通常沒有渲染，但「廣播頻率」仍可分開。  
- 房間空轉：Lobby 不一定要跑玩法 tick；進 Playing 再啟動，省資源。

## 遊戲 Server 會用在哪

Arena Mini M4：大約 **20Hz tick**，每 tick 套用 input、更新座標／碰撞，廣播 `type=state`。  
你之後改 `ScoreToWin`、地圖大小、速度時，會感謝自己用的是固定步——手感調整是參數，不是玄學。

## 請丟掉的舊習慣

1. **每收到一包就改狀態並廣播**——無節奏、難測、易被刷封包拖死。  
2. **用牆鐘時間直接當物理積分且不固定步**，又沒有任何 clamp。  
3. **每個 tick `go func() { simulate() }`**——短命 goroutine 海，競態與調度開銷齊飛；一房一個長壽命 loop 通常較乾淨。  
4. **房間銷毀後 ticker 還在跑**——典型洩漏。

## 動手練習

### 必做

1. 跑 `examples/h05-tick`，看 `applied` 與 `sum` 的輸出節奏。  
2. 改成 10Hz（100ms），觀察同一批輸入被批量套用的方式有何不同。  

### 選做

1. 給 inbox 設小容量（例如 2），刻意塞爆，決定要丟棄還是退出，並印 log。  
2. 把「apply + snapshot」抽成函式，寫一個**不 Sleep**的單元測試：連續呼叫 `Tick()` 斷言結果。  

## 常見坑

- **忘記 `Stop` ticker**：房結了還在 tick。  
- **在 tick case 裡做很重的同步 I/O**（寫 DB、同步 HTTP）：一卡就整房時間膨脹。  
- **無界 slice 當 inbox**：客戶端狂送，記憶體被灌爆。  
- **用 `time.Sleep` 假裝 ticker**：不好取消、不好跟多路事件 `select`；長駐服務請用 `Ticker`／`timer` + context。

下一章 H6：輸入怎麼入隊、怎麼校驗、一個 tick 內多筆輸入怎麼合併。
