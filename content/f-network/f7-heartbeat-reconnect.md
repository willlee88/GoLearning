---
lessonId: "F7"
title: "心跳、重連與會話"
description: "連線≠玩家；閒置檢測與重連綁定概念。"
volume: "f"
order: 7
level: "l2"
status: "ready"
path_required: true
tags: ["session", "websocket"]
example: "examples/f08-ws-room"
prev: "F6"
next: "F8"
---

## 本章你會建立的心智模型

行動網路會斷。設計上要分開：

1. **連線（Connection）** — 當下 socket  
2. **會話（Session）** — 玩家身份與可恢復狀態  

心跳（ping/應用層 heartbeat）用來偵測假死；重連時用 token/session id 綁回同一玩家，避免變成「另一個匿名」。

## Python 對照

概念相同；實作上用 cancel scope / task 做 ping 迴圈。Go 用 ticker + context。

## L1 能用

```go
// 概念：讀 loop 與 ping loop 共享 conn 寫鎖
ticker := time.NewTicker(15 * time.Second)
defer ticker.Stop()
for {
	select {
	case <-ticker.C:
		_ = send(conn, pingMsg)
	case <-ctx.Done():
		return
	}
}
```

## L2 機制

- 讀 deadline：多久沒訊息就踢。  
- 應用層 `{"type":"ping"}` / `pong` 簡單可觀測。  
- 重連：伺服器保留 session 一段 TTL；局內狀態是否可恢復看遊戲。  
- 冪等：重送「開火」不能射兩次（G/H 卷）。  

## 請丟掉的 Python 習慣

1. 斷線就當玩家毀滅、不留任何恢復窗口。  
2. 無心跳靠 TCP keepalive 對付所有中間盒（不夠）。  

## 遊戲 Server 連結

Arena Mini 進階分支：session id、短暫停牌、重連回房。M3 先把概念與訊息型別預留。

## 練習

### 必做

1. 在 f08 或 arena-mini 加 `type=ping/pong` 日誌。  
2. 設計一張表：哪些狀態可重連恢復、哪些必須重開局。  

## 延伸閱讀

- 雲廠商 WebSocket 閒置逾時文件  
