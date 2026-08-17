---
lessonId: "A18"
title: "逃逸分析與配置直觀"
description: "stack vs heap；用編譯器旗標看逃逸。"
volume: "a"
order: 18
level: "l3"
status: "ready"
path_required: false
tags: ["performance", "gc"]
example: "examples/a18-escape"
prev: "A17"
next: "B0"
---

## 本章你會建立的心智模型

Go 自動管理記憶體。編譯器做**逃逸分析**：變數若只在函式內用，可能留在 stack；否則上 heap 給 GC 管。深刻理解不是背規則，而是：**熱路徑減少不必要分配**，並用工具驗證。

## Python 對照

CPython 物件幾乎都在堆上；Go 有機會更省，但寫錯一樣狂配。

## L1 能用

```bash
go build -gcflags="-m" .
```

看 `escapes to heap` 提示。

範例：`examples/a18-escape`。

## L2 機制

常見上堆原因：

- 回傳局部變數指標  
- 介面裝箱  
- 閉包捕獲  
- 太大無法 stack  

優化順序：算法 → 減少分配 → 再微觀。

## L3 深潛

- `sync.Pool` 重用 buffer（注意生命週期）。  
- pprof heap 對照。  

## 請丟掉的 Python 習慣

1. 過早微優化。  
2. 無 pprof 只靠感覺。  

## 遊戲 Server 連結

每 tick 對 N 人 `json.Marshal` 會配；可重用 buffer、降頻、delta（見 H7/J6）。

## 練習

### 必做

1. 對範例跑 `-gcflags=-m` 讀兩行輸出。  
2. 猜哪種寫法更易逃逸，再驗證。  

## 延伸閱讀

- Go Blog: Profiling Go Programs  
