---
lessonId: "H9"
title: "規則與 I/O 分離"
description: "純邏輯可測；hub 只負責收發。"
volume: "h"
order: 9
level: "l2"
status: "ready"
path_required: true
tags: ["architecture", "testing"]
example: "examples/h06-apply-input"
prev: "H8"
next: "H10"
---

## 本章你會建立的心智模型

```text
hub/ws  →  decode  →  room.PushInput / Ready
room tick → game.Apply  →  Snapshot
hub      →  encode broadcast
```

`game` package **不 import** `net` / `websocket`。這樣 `go test` 不需開 port。

## Python 對照

領域層 vs 介面適配器（六角形架構的小規模版）。

## L1 能用

```go
// game/room.go — 純邏輯
func (r *Room) Apply(in Input) error
func (r *Room) Snapshot() State

// hub — I/O
func (h *Hub) onMessage(...) { r.PushInput(...) }
```

## L2 機制

- 時間用 `Clock` 介面可測。  
- 隨機數注入 seed。  
- 單測：表驅動走完整 lobby→play。  

## 請丟掉的 Python 習慣

1. 單元測試裡真的 listen port。  
2. 規則函式簽名帶 `websocket`。  

## 遊戲 Server 連結

看 `demo/arena-mini/server/internal/game`。

## 練習

### 必做

1. 指出 demo 中純邏輯檔案路徑。  
2. 為 `Apply` 寫一個不開 WS 的測試。  

## 延伸閱讀

- B4 API 慣例  
