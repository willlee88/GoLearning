---
lessonId: "F6"
title: "WebSocket 原理與實作"
description: "HTTP 升級、幀、與瀏覽器互通。"
volume: "f"
order: 6
level: "l2"
status: "ready"
path_required: true
tags: ["websocket", "game-server"]
example: "examples/f08-ws-room"
prev: "F5"
next: "F7"
---

## 本章你會建立的心智模型

WebSocket 在 HTTP(S) 上**升級**成長連線，之後雙向傳訊息幀。瀏覽器原生支援，是網頁遊戲與本站 demo 的首選。底層仍是 TCP（或 TLS）——F3 的「應用訊息邊界」變成「一則 WS message」，但你仍要定義 JSON 信封。

## Python 對照

| Python | Go |
|--------|-----|
| `websockets` / Starlette | `golang.org/x/net/websocket` 或 `gorilla/websocket` / `nhooyr` |
| async for message | 讀 loop goroutine |

## L1 能用（概念）

```text
Client: GET /ws  Connection: Upgrade
Server: 101 Switching Protocols
之後: 雙向 text/binary frames
```

本站範例與 Arena Mini 使用 `x/net/websocket`（教學依賴少）。

## L2 機制

- 每個連線：讀 loop；寫入需串行（mutex 或 write pump）。  
- 跨源：瀏覽器 CORS／Origin 檢查策略。  
- 代理：nginx 要開 WebSocket upgrade。  
- 與 HTTP 共享 port 很常見。  

## L3 深潛（可選）

- 控制幀 ping/pong。  
- 壓縮擴展 permessage-deflate 成本。

## 請丟掉的 Python 習慣

1. 在讀 callback 裡做重運算阻塞讀 loop。  
2. 多 goroutine 無鎖同寫一條 conn。  

## 遊戲 Server 連結

`/ws?room=&name=` → session → room inbox。狀態廣播用 JSON `type=state`。

## 練習

### 必做

1. 閱讀 `examples/f08-ws-room` 與 `demo/arena-mini`。  
2. 雙分頁連同一 room 互傳。  

## 延伸閱讀

- RFC 6455  
