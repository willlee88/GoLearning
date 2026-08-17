---
lessonId: "C9"
title: "Scheduler 直觀"
description: "G/M/P 不用背，但要建立正確成本感。"
volume: "c"
order: 9
level: "l3"
status: "ready"
path_required: false
tags: ["scheduler", "runtime"]
example: ""
prev: "C8"
next: "C10"
---

## 本章你會建立的心智模型

Go runtime 把大量 goroutine（G）多路複用到作業系統執行緒（M），中間有邏輯處理器 P。你不必默寫源碼，但要知道：

- goroutine 建立便宜，**不是免費**  
- 阻塞在系統呼叫可能佔用執行緒  
- `GOMAXPROCS` 影響平行度  
- 長 CPU 迴圈可能讓調度不公平（可 `runtime.Gosched`，更常見是拆工作）  

## Python 對照

GIL 下 CPU 平行受限；Go 可真平行，race 也更真實。

## L1 能用

```bash
go test -cpu=1,4
```

觀察不同 GOMAXPROCS 下行為（邏輯上）。

## L2 機制

- 網路 poller 讓許多阻塞變高效。  
- 鎖競爭與假共享會限制擴展。  
- profile 看 `runtime` 與你的包。  

## 請丟掉的 Python 習慣

1. 開 10 萬 goroutine「因為便宜」無設計。  
2. 在持鎖時做重 CPU。  

## 遊戲 Server 連結

每連線 1–2 個 goroutine 通常 OK；每 tick 再炸 goroutine 要小心。

## 練習

### 必做

1. 讀一篇 Go scheduler 概述（官方 blog）。  
2. 估算：1 萬連線 × 2 goroutine 的數量級是否可接受。  

## 延伸閱讀

- Scalable Go Scheduler Design Doc / Go Blog  
