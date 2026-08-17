---
lessonId: "C0"
title: "C 卷導讀 · 併發模型"
description: "goroutine、channel、sync 與 race——遊戲 Server 的中樞神經。"
volume: "c"
order: 0
level: "l1"
status: "ready"
path_required: true
tags: ["concurrency"]
example: ""
prev: "B4"
next: "C1"
---

## 本章你會建立的心智模型

C 卷是「深刻理解 Go」的分水嶺。你要能推理：**誰擁有資料、誰在什麼時候寫、如何停止、如何證明無 race**。遊戲 Server 的連線 loop、房間 tick、廣播，全部建立在這上面。

## 本卷地圖（M2 核心）

| 章 | 主題 |
|----|------|
| C1 | goroutine 生命週期與洩漏 |
| C2 | channel 語意 |
| C3 | select 與超時 |
| C4 | context 取消樹 |
| C5 | Mutex 與共享狀態 |
| C6 | race detector 實戰 lab |

（C7+ 背壓、scheduler 深潛等後續里程碑再補。）

## 工具

```bash
go test -race ./...
```

## 檢查點

完成 `examples/c06-race-lab`：先紅後綠，理解為何。
