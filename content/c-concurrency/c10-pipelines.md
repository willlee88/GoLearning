---
lessonId: "C10"
title: "管道模式與扇入扇出"
description: "階段化處理；何時用 pipeline 何時別用。"
volume: "c"
order: 10
level: "l2"
status: "ready"
path_required: false
tags: ["pipeline"]
example: "examples/c10-pipeline"
prev: "C9"
next: "D0"
---

## 本章你會建立的心智模型

Pipeline：資料經多個 stage，每 stage 一組 goroutine + channel。適合串流處理（日誌、回放轉檔）；**不適合**把簡單函式呼叫硬拆成管道增加複雜度。

## Python 對照

generator 管線；Go 用 channel 連 stage。

## L1 能用

```text
gen → parse → filter → sink
```

每段 `for v := range in { out <- f(v) }`，結束時 close out。

範例：`examples/c10-pipeline`。

## L2 機制

- 誰 close channel 要清楚。  
- 扇出 fan-out 要 fan-in 合併。  
- 用 context 取消整條鏈。  

## 請丟掉的 Python 習慣

1. 為炫技把 20 行邏輯拆 5 個 channel。  

## 遊戲 Server 連結

離線分析回放可用 pipeline；線上 Room tick 通常是單線狀態機更清楚。

## 練習

### 必做

1. 跑 `examples/c10-pipeline`。  
2. 加一個 stage 統計數量。  

## 延伸閱讀

- Go Blog: Pipelines  
