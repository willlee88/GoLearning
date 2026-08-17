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

## 本章你會建立的心智模型

`time.Duration` 是納秒基數的命名型別。遊戲用 `Ticker` 做 tick；超時用 `Timer` 或 `context.WithTimeout`。

## L1 能用

```go
time.Sleep(100 * time.Millisecond)
t := time.NewTicker(50 * time.Millisecond)
defer t.Stop()
```

## 遊戲 Server 連結

Arena Mini 20Hz = `50 * time.Millisecond`。

## 練習

### 必做

1. 寫一個 3 次 tick 後停止的小程式。  

## 延伸閱讀

- `time` package  
