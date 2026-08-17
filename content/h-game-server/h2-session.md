---
lessonId: "H2"
title: "Session 與連線模型"
description: "連線≠玩家；進房、斷線、重連綁定。"
volume: "h"
order: 2
level: "l2"
status: "ready"
path_required: true
tags: ["session"]
example: ""
prev: "H1"
next: "H3"
---

## 本章你會建立的心智模型

**Connection** 是 socket；**Session** 是「這個玩家在線上下文」。斷線可以銷毀連線但短暫保留 Session（重連窗口）。進房後 Session 綁到 Room seat。

## Python 對照

類似 websocket 連線物件 vs 使用者 session cookie——即時遊戲要更顯式。

## L1 能用

```go
type Session struct {
	ID     string
	Name   string
	RoomID string
	Conn   Conn // 可抽象
}
```

## L2 機制

- 讀 loop 屬於連線；業務狀態屬 session/room。  
- 斷線：從 room 移除或標 `Disconnected` 等重連。  
- 同名佔據：拒絕或踢舊連線（要文件化）。  

## 請丟掉的 Python 習慣

1. 只用 `websocket` 物件當唯一身份。  
2. 斷線立即抹掉一切進度（總是）。  

## 遊戲 Server 連結

M4 Arena Mini：`name` query 當簡化 session；進階再加 token。

## 練習

### 必做

1. 寫下：斷線 5 秒內重連，局內座標應如何處理。  
2. 列出 session 至少 4 個欄位。  

## 延伸閱讀

- F7 心跳與重連  
