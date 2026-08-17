---
lessonId: "A0"
title: "A 卷導讀 · 語言與資料模型"
description: "從安裝到 slice/map/struct：建立 Go 的資料與型別心智模型。"
volume: "a"
order: 0
level: "l1"
status: "ready"
path_required: true
tags: ["language"]
example: ""
prev: "P0.8"
next: "A1"
---

## 本章你會建立的心智模型

A 卷是語言地基。目標不是背 API，而是能回答：

1. **資料放哪、怎麼共享**（值、指標、slice 底層陣列）  
2. **誰能看見誰**（package 匯出）  
3. **編譯期契約長什麼樣**（型別、方法集預告）  

**A1–A14** 為語言核心主路徑；**A15–A18** 為廣覆蓋／深潛（嵌入、any、unsafe、逃逸）。

## 建議走法

| 順序 | 章 | 重點 |
|------|----|------|
| 1 | A1–A3 | 安裝、package、零值 |
| 2 | A4–A7 | 控制流、函式、defer、指標 |
| 3 | A8 | slice / map 檢查點 |
| 4 | A9–A11 | 字串、struct、method |
| 5 | A12–A14 | interface、type switch、generics |

每章都有 **Python 對照** 與 **遊戲 Server 連結**。做完 A 卷後進入 **B（錯誤）→ C（併發）**。

## Python 對照（本卷總覽）

| 你熟悉的 | 在 A 卷會變成 |
|----------|----------------|
| 動態型別 + `None` | 靜態型別 + 零值 / `nil` |
| list / dict | slice / map（語意更「底層」） |
| class 方法 | struct + method（值或指標 receiver） |
| 模組與 `__all__` 約定 | 大寫匯出，編譯器強制 |

## 練習

### 必做

1. 順讀 A1→A8，每章至少跑一個 `go run` 或 `go test`。  
2. 完成 `examples/a08-player-registry`（A8 檢查點）。  

### 選做

1. 把 P0 的 `p0-config-stats` 改寫成多 package（呼應 A2）。  

## 遊戲 Server 連結

A 卷結束後，你應能乾淨地表達：玩家實體、房間索引 map、快照用的 slice buffer——這些是 Room / Session 的細胞。
