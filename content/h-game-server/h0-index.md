---
lessonId: "H0"
title: "H 卷導讀 · 遊戲 Server 核心"
description: "Session、Room、Tick、權威狀態——把網路層變成可玩的局。"
volume: "h"
order: 0
level: "l1"
status: "ready"
path_required: true
tags: ["game-server"]
example: ""
prev: "G3"
next: "H1"
---

## 本章你會建立的心智模型

F/G 卷解決「怎麼連、怎麼說」；H 卷解決「局怎麼跑」。核心物件：

| 概念 | 角色 |
|------|------|
| Session | 連線上的玩家身份 |
| Lobby / Room | 隔離的一局空間 |
| Phase | 大廳／進行中／結束 |
| Tick | 固定頻率更新 |
| Command | 客戶端意圖（非真相） |
| State | 伺服器權威快照 |

## 本卷地圖（M4）

| 章 | 主題 |
|----|------|
| H1 | 後端職責地圖 |
| H2 | Session 與連線 |
| H3 | Room 生命週期狀態機 |
| H4 | 權威 Server |
| H5 | Tick 與 fixed timestep |
| H6 | 輸入佇列與校驗 |
| H7 | 狀態同步策略 |
| H8 | 匹配入門 |
| H9 | 規則與 I/O 分離 |
| H10 | Lab：Arena Mini 權威對戰 |

## 檢查點

跑通 `demo/arena-mini`：**兩人 Ready → 開局 → 方向鍵/按鈕移動（伺服器算座標）→ 狀態廣播**。
