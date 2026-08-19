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

## 這章你會搞懂什麼

**WebSocket** 讓瀏覽器跟伺服器維持一條**雙向長連線**。它不是憑空變出來的協定魔術，而是：

1. 先用普通 HTTP 打去（例如 `GET /ws`）  
2. 雙方同意 **Upgrade**  
3. 之後改走 WebSocket 的**訊息幀**（text／binary／控制幀）

底層通常還是 TCP（或 TLS）。F3 說的「應用訊息邊界」，在 WS 上常常變成「一則 WebSocket message」；但你**仍要**定義 JSON 信封（type／payload），否則只是換成另一種粘在一起的字串。

## 先跟 Python 對一下

| Python | Go | 注意 |
|--------|-----|------|
| `websockets`／Starlette WS | `golang.org/x/net/websocket`、gorilla、nhooyr… | 本站 demo 用 `x/net/websocket`，依賴少、適合教學 |
| `async for message in ws` | 讀 loop 放在 goroutine | |
| 自己要小心同時寫 | 一樣：寫入要串行（mutex 或專責 write pump） | |

## 怎麼寫（概念）

```text
Client: GET /ws  Connection: Upgrade  Upgrade: websocket ...
Server: 101 Switching Protocols
之後: 雙向 text/binary frames
```

瀏覽器端大致是：

```js
const ws = new WebSocket("ws://localhost:8080/ws?room=r1&name=Ada");
ws.onmessage = (ev) => console.log(ev.data);
ws.send(JSON.stringify({ type: "chat", payload: { text: "hi" } }));
```

伺服器端：HTTP handler 裡完成升級，然後：

- 一個 **讀 loop**：收訊息 → 解析信封 → 丟進房間 inbox  
- **寫入**：廣播時加鎖，或固定由單一 goroutine 寫  

本站範例：`examples/f08-ws-room`；完整一點：`demo/arena-mini`。

## 為什麼這樣設計／底層在幹嘛

1. **為什麼要 Upgrade**  
   沿用 80／443、穿透企業代理與防火牆的路徑較友善，也讓「同一 port 同時服務網頁 + WS」很自然。

2. **為什麼寫入要串行**  
   兩個 goroutine 同時 `Write` 同一條連線，幀會交錯損壞。常見模式：`writeMu sync.Mutex`，或 channel 交給唯一 writer。

3. **讀 loop 不要做重活**  
   解析後把命令丟進房間佇列，讓房間 tick／邏輯 goroutine 處理；否則一個慢 handler 卡住收封包。

4. **Origin／跨源**  
   瀏覽器會帶 Origin；伺服器要決定要不要接受。教學本機常放寬，生產要白名單。

5. **反代**  
   前面有 nginx／Caddy，要開啟 WebSocket upgrade（連線升級標頭與長逾時）。否則「本機可以、上雲不行」。

6. **進階可先略過**  
   - 控制幀 ping／pong（連線保活）  
   - `permessage-deflate` 壓縮：省頻寬但吃 CPU  

## 遊戲 Server 會用在哪

```text
/ws?room=&name= → Session → 加入 Room → 收命令／廣狀態
```

常見訊息：`type=state` 廣播快照、`type=chat`、`type=move`。HTTP 仍負責靜態頁與 `/healthz`。

## 請丟掉的舊習慣

1. **在讀 callback／讀 loop 裡做重運算或 `sleep`。**  
2. **多 goroutine 無鎖同寫一條 conn。**  
3. **以為有了 WS 就不必定義應用協定**——沒有 `type` 會變大型字串地獄。

## 動手練習

### 必做

1. 閱讀 `examples/f08-ws-room` 與 `demo/arena-mini` 的連線進入點。  
2. 開兩個瀏覽器分頁連同一 room，確認能互傳／看到廣播。  

### 選做

1. 在伺服器 log 每則訊息的 `type`。  

## 常見坑

- **混用 `ws://` 與頁面 `https://`**：瀏覽器會擋（mixed content）——要用 `wss://`。  
- **忘了處理連線關閉**：成員表要刪除，否則廣播送到死連線。  
- **把巨大狀態每幀塞 text JSON 卻不量測**：先跑通，再優化（G／H／J）。  
- **依賴函式庫預設的跨源策略不查文件**：上線才發現全拒或全開。

## 延伸閱讀

- RFC 6455（WebSocket 協定；可略讀概念）  
- `examples/f08-ws-room`  
