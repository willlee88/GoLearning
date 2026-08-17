---
lessonId: "H1"
title: "遊戲後端職責地圖"
description: "Gateway / Game / Match / Persist 如何切分。"
volume: "h"
order: 1
level: "l2"
status: "ready"
path_required: true
tags: ["architecture", "game-server"]
example: ""
prev: "H0"
next: "H2"
---

## 本章你會建立的心智模型

大型專案會拆服務；教學 demo 可先單進程，但**頭腦裡要分層**：

```text
Gateway(WS/TCP) → Session
       ↓
   Match/Lobby → Room Runtime (tick)
       ↓
   Rules (pure) + optional Persist
```

Gateway 不懂規則；Rules 不懂 socket。

## Python 對照

像把 FastAPI 路由、背景 worker、領域模型拆開——Go 用 package 邊界表達。

## L1 能用

Arena Mini 對應：

- `cmd/server`：組裝  
- `internal/hub`：連線與房間表  
- `internal/game`：規則與 tick 狀態  

## L2 機制

| 職責 | 做 | 不做 |
|------|----|------|
| Gateway | 握手、讀寫、心跳 | 判定勝負 |
| Room | 生命週期、廣播節奏 | DB schema |
| Rules | 合法移動、碰撞 | 持有 `net.Conn` |
| Persist | 結算、帳號 | 每 tick 寫庫 |

## 請丟掉的 Python 習慣

1. 一個 God class 又管 socket 又算物理。  
2. 請求執行緒直接改全域 dict 無隔離。  

## 遊戲 Server 連結

水平擴展時「房間親和」要求狀態有明確所在——先在單機把 Room 邊界劃清。

## 練習

### 必做

1. 畫 Arena Mini 三層對應表（你自己的理解）。  
2. 指出目前哪裡還「規則碰到 I/O」。  

## 延伸閱讀

- 規劃書 §6.4  
