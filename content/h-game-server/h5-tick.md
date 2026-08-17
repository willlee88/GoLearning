---
lessonId: "H5"
title: "Tick 與 fixed timestep"
description: "固定頻率更新；輸入累積、狀態輸出。"
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

## 本章你會建立的心智模型

**Tick** 是 Server 的心跳：例如每 50ms（20Hz）處理輸入、推進模擬、廣播。Fixed timestep 讓規則可重現、好調參。可用 `time.Ticker` + `select` 收命令／tick／取消。

## Python 對照

| Python | Go |
|--------|-----|
| `asyncio` 迴圈 sleep | `time.NewTicker` + `select` |
| 變步長 dt | 固定 `dt`，餘數累積（進階） |

## L1 能用

```go
ticker := time.NewTicker(50 * time.Millisecond)
defer ticker.Stop()
for {
	select {
	case cmd := <-inbox:
		buf = append(buf, cmd)
	case <-ticker.C:
		applyAll(buf)
		buf = buf[:0]
		broadcast(snapshot())
	case <-ctx.Done():
		return
	}
}
```

範例：`examples/h05-tick`。

## L2 機制

- 輸入在 tick 之間排入佇列，tick 時批量套用。  
- 廣播頻率可 ≤ 模擬頻率。  
- 過載：有界 inbox，滿則丟或斷線。  

## 請丟掉的 Python 習慣

1. 每收到一包就立刻改狀態並廣播（無節奏，難重現）。  
2. 用 wall clock 直接當物理積分且不固定步。  

## 遊戲 Server 連結

Arena Mini M4：20Hz tick，廣播 `type=state`。

## 練習

### 必做

1. 跑 `examples/h05-tick`。  
2. 改 10Hz 觀察訊息量變化。  

## 延伸閱讀

- Fix Your Timestep（概念，Gaffer on Games）  
