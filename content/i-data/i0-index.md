---
lessonId: "I0"
title: "I 卷導讀 · 資料與快取"
description: "什麼進記憶體、什麼進 DB、Redis 在遊戲裡的角色。"
volume: "i"
order: 0
level: "l1"
status: "ready"
path_required: true
tags: ["data"]
example: ""
prev: "H10"
next: "I1"
---

## 本章你會建立的心智模型

遊戲 Server 的狀態分層：

| 層 | 例子 | 延遲／一致性 |
|----|------|----------------|
| 記憶體（Room） | 座標、tick、準備狀態 | 最快，進程沒就沒 |
| 快取（Redis 等） | presence、排行榜、pubsub | 快，可跨進程 |
| 資料庫 | 帳號、庫存、賽季結算 | 慢但持久 |

**不是什麼都寫 DB。** Tick 熱路徑只碰記憶體。

## 本卷地圖（M5）

| 章 | 主題 |
|----|------|
| I1 | 持久化邊界 |
| I2 | database/sql 模式 |
| I3 | Redis 用途地圖 |
| I4 | 快取一致性直覺 |
| I5 | 背景任務與結算 |

## 練習

讀完後能回答：Arena Mini 哪些狀態該持久化、哪些絕對不該每 tick 寫庫。
