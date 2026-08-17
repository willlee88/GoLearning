---
lessonId: "H10"
title: "Lab：Arena Mini 權威對戰"
description: "H 卷檢查點：Ready、Tick、Input、State 廣播。"
volume: "h"
order: 10
level: "l2"
status: "ready"
path_required: true
tags: ["lab", "capstone", "game-server"]
example: "demo/arena-mini"
prev: "H9"
next: "I0"
---

## 本章你會建立的心智模型

把 H1–H9 落到 **Arena Mini M4**：

1. 進房（Lobby）  
2. 送 `ready`  
3. ≥2 人且皆 ready → **Playing** + 啟動 tick  
4. 送 `input`（dx,dy）  
5. Server 每 tick 更新座標、**碰撞得分**、廣播 `state`  
6. 先到 3 分或超時 → **Ended** → 數秒後回 **Lobby**  

## L1 能用

```powershell
cd F:\GoLearning\demo\arena-mini\server
go mod tidy
go test ./...
go run ./cmd/server
# http://localhost:8080
```

兩分頁：不同 name、同 room → 雙方 Ready → 用按鈕或方向移動。

## L2 檢查清單

- [ ] 客戶端不再直傳權威座標當真相  
- [ ] `internal/game` 可單測  
- [ ] Playing 前忽略 input  
- [ ] 斷線清理成員  
- [ ] `/healthz` `/rooms` 仍可用  

## 請丟掉的 Python 習慣

1. 開局邏輯寫在前端「人齊了就自己 start」。  
2. 無 tick 的事件驅動亂序狀態。  

## 練習

### 必做（H / M4 檢查點）

1. 雙人完成一局移動同步。  
2. 閱讀 `internal/game` 與 `internal/hub` 邊界。  
3. 加一個規則：速度或地圖大小改參數再測。  

### 選做

1. 改 `ScoreToWin` / `HitRadius` 參數再測手感。  
2. 加第三名玩家同房（capacity 已支援）。  

## 延伸閱讀

- `demo/arena-mini/README.md`  

## M4 收束

你已具備最小「可玩」遊戲後端閉環。下一階段 **M5（I/J）**：資料、觀測、壓測、優雅關閉。  
