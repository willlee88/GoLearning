---
lessonId: "F8"
title: "Lab：WebSocket 房間廣播"
description: "F 卷檢查點：多房間、JSON 信封、廣播。"
volume: "f"
order: 8
level: "l2"
status: "ready"
path_required: true
tags: ["lab", "websocket", "game-server"]
example: "examples/f08-ws-room"
prev: "F7"
next: "G0"
---

## 本章你會建立的心智模型

把 F1–F7 收成可跑系統：

- HTTP `/healthz`  
- WS `/ws?room=&name=`  
- JSON 信封 `{type, from, payload, room}`  
- Hub → Room → members 廣播  
- 連線退出時清理  

這是 Arena Mini 的縮小版；demo 會再加房間列表與狀態訊息。

## L1 能用

```bash
cd examples/f08-ws-room
go run .
# 瀏覽器開附帶的 static 或用兩個 WS 客戶端
```

## L2 機制複習清單

- [ ] 每連線讀 loop  
- [ ] 寫入串行或可接受的簡單模型  
- [ ] map 成員表有鎖  
- [ ] leave 時 delete  
- [ ] 應用訊息有 type  

## 請丟掉的 Python 習慣

1. 把所有人塞一個全局 list 無鎖。  
2. handler 裡 `time.sleep` 阻塞。  

## 遊戲 Server 連結

下一步 G 卷定義更好的命令／事件分離；H 卷加 tick 與權威狀態。

## 練習

### 必做（F 卷檢查點）

1. 跑通 `examples/f08-ws-room` 雙客戶端。  
2. 跑通 `demo/arena-mini`（`go mod tidy && go run ./cmd/server`）。  
3. 加一種 `type=pos` 廣播座標（字串即可）。  

### 選做

1. HTTP `GET /rooms` 回各房人數（見 arena-mini）。  

## 延伸閱讀

- `demo/arena-mini/README.md`  
