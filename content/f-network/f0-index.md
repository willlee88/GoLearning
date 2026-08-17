---
lessonId: "F0"
title: "F 卷導讀 · 網路"
description: "從 TCP/HTTP 到 WebSocket：遊戲 Server 的連線層。"
volume: "f"
order: 0
level: "l1"
status: "ready"
path_required: true
tags: ["network"]
example: ""
prev: "C6"
next: "F1"
---

## 本章你會建立的心智模型

F 卷把併發能力接到**真實網路 I/O**。遊戲後端幾乎永遠是長連線 + 請求控制面：

- **TCP / 自訂幀**：傳統 socket 遊戲  
- **HTTP**：健康檢查、登入、管理 API  
- **WebSocket**：瀏覽器可連的即時通道（本站 Arena Mini 主軸）

## 建議路徑

| 章 | 主題 |
|----|------|
| F1 | 網路心智：延遲、吞吐、可靠性 |
| F2 | TCP server/client |
| F3 | 封包邊界與長度前綴 |
| F4 | net/http 基礎 |
| F5 | 中介層與路由模式 |
| F6 | WebSocket 原理與實作 |
| F7 | 心跳、重連、會話 |
| F8 | Lab：WS 房間廣播 |

做完 F8 後接 **G 卷（協定）** 與 Arena Mini。

## Python 對照

| Python | Go |
|--------|-----|
| `socket` / `asyncio.open_connection` | `net` |
| Flask/FastAPI | `net/http`（先標準庫） |
| `websockets` / Starlette WS | `x/net/websocket` 或第三方 |

## 練習

順讀 F1→F8，跑通 `examples/f08-ws-room` 與 `demo/arena-mini`。
