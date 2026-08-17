---
lessonId: "C8"
title: "背壓與有界佇列"
description: "無界 channel 的危險；滿了怎麼辦。"
volume: "c"
order: 8
level: "l2"
status: "ready"
path_required: false
tags: ["backpressure"]
example: "examples/c08-backpressure"
prev: "C7"
next: "C9"
---

## 本章你會建立的心智模型

生產者快、消費者慢時，**無界佇列 = 延遲的 OOM**。背壓（backpressure）讓壓力回傳：阻塞、丟棄、斷線、降級。

## Python 對照

`asyncio.Queue(maxsize=N)` 滿了 await；Go 用有界 channel 或拒絕策略。

## L1 能用

```go
ch := make(chan Job, 1024)
select {
case ch <- job:
default:
	return ErrBusy // 有界且非阻塞丟棄/拒絕
}
```

範例：`examples/c08-backpressure`。

## L2 機制

策略：

| 策略 | 場景 |
|------|------|
| 阻塞 | 可接受延遲 |
| 丟最舊/最新 | 遙測、可丟狀態 |
| 斷開客戶端 | 輸入洪水 |
| 降級 | 關特效、降 tick |

## 請丟掉的 Python 習慣

1. `go func` + 無界 channel 當萬靈丹。  
2. 無限 list 緩衝。  

## 遊戲 Server 連結

input inbox、結算 worker queue 都要有界（I5/J6）。

## 練習

### 必做

1. 跑範例看滿佇列行為。  
2. 為 Arena 設計「每玩家每秒最多 40 input」。  

## 延伸閱讀

- Reactive 背壓概念  
