---
lessonId: "H10"
title: "Lab：Arena Mini 權威對戰檢查點"
description: "把 H1–H9 落到可玩的一局：Ready、Tick、Input、State 廣播。雙人同步過關，才算 H 卷收束。"
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

## 這章你會搞懂什麼

這是 H 卷的**實作檢查點**，不是再塞新理論。  
你要把前面學的名詞，全部對到 `demo/arena-mini` 的行為上：

1. 進房（Lobby）  
2. 送 `ready`  
3. ≥2 人且皆 ready → **Playing**，並啟動 tick  
4. 送 `input`（dx, dy）——意圖，不是自稱座標  
5. Server 每 tick 更新座標、處理碰撞得分、廣播 `state`  
6. 先到分數或超時 → **Ended** → 數秒後可回 **Lobby**  

讀完並動手後，你應該能指著程式說：「這段是 hub，那一段是 game；開局真相在 Server。」

## 這局在驗證什麼（對照前面各章）

| 章 | Lab 裡你要看到的證據 |
|----|----------------------|
| H1 職責地圖 | `cmd/server` 組裝；`hub` 收發；`game` 規則 |
| H2 Session | 不同 name 進同房；斷線會清理成員 |
| H3 Room FSM | Playing 前的 input 被忽略／拒絕；Ready 才能開局 |
| H4 權威 | 客戶端不靠直傳座標當真相 |
| H5 Tick | 移動是一格一格依節奏更新，不是「封包一到瞬移無序」 |
| H6 輸入佇列 | `type=input` → 入隊 → tick 套用 |
| H7 同步 | 定期收到全量 `state` |
| H8 匹配 | 相同 room id ≈ 手動湊局 |
| H9 規則／I/O | `go test ./internal/game` 可不開瀏覽器 |

若某一列對不上，回去補那一章，別只靠「畫面好像會動」。

## 怎麼跑

```powershell
cd F:\GoLearning\demo\arena-mini\server
go mod tidy
go test ./...
go run ./cmd/server
```

瀏覽器開：`http://localhost:8080`

建議流程：

1. 開**兩個分頁**（或兩個視窗）  
2. 填**不同 name**、**相同 room**  
3. 雙方按 Ready  
4. 開局後用按鈕或方向鍵移動  
5. 確認你看得到對方位置隨 tick 更新；碰撞／得分若有實作也一併確認  

更完整的說明見：`demo/arena-mini/README.md`。

## 檢查清單（請逐項打勾）

做 Lab 時把下面當驗收單：

- [ ] 客戶端**不再**把權威座標當真相直傳給 Server 採信  
- [ ] `internal/game` 可以單測（`go test` 通過）  
- [ ] Playing **之前**的 input 不會亂改正式對戰狀態  
- [ ] 斷線會清理成員（不會變成幽靈座位一直佔著）  
- [ ] `/healthz`、`/rooms` 這類探針／除錯端點仍可用（若 demo 有提供）  
- [ ] 兩人真的完成至少一局：Ready → 移動同步 → 看到 state  

全部勾完，H 卷才算「可玩閉環」落地，而不只是讀過名詞。

## 建議閱讀順序（進 repo 時）

1. `internal/game`：Phase、Ready、Input、Tick、Snapshot——規則真相  
2. `internal/hub`：連線註冊、訊息分派、廣播——I/O  
3. `cmd/server`：怎麼組起來聽 port  
4. 前端頁面（若有）：它送的是 `ready`／`input` 還是舊的 `pos`？  

一邊看一邊問 H9 那句：依賴方向有沒有被破壞？

## 請丟掉的舊習慣

1. **開局邏輯寫在前端**：「我看人齊了就自己 start」——兩端時間差會讓世界分裂，也無法防惡意客戶端。  
2. **沒有 tick 的純事件堆疊**：誰先到誰先改，局很難重現。  
3. **只測快樂路徑**：從不試「一人 Ready、Lobby 就狂送 input、斷線重進」。  
4. **改手感卻改到協定**：想調速度時，應改 Server 參數（如 speed、地圖大小），不是讓客戶端多傳倍速。

## 動手練習

### 必做（H／M4 檢查點）

1. 雙人完成一局移動同步（建議錄一下自己的操作步驟，之後回歸用）。  
2. 閱讀並用自己的話寫三句：`internal/game` 與 `internal/hub` 的邊界。  
3. 改一個規則參數（速度或地圖大小），再測手感——確認改的是 Server 側。  

### 選做

1. 調整 `ScoreToWin`／`HitRadius`（若專案有這些常數）感受對局長度與碰撞寬鬆度。  
2. 加第三名玩家進同房（capacity 允許的話），觀察廣播與開局條件。  
3. 故意在 Lobby 送 input，確認 Server 行為符合你的 FSM 設計。  

## 常見坑

- **兩個分頁用了同一個 name**：頂號／覆蓋／拒絕政策若沒注意，會以為是同步 Bug。  
- **以為 `go run` 過了就不用 `go test`**：規則回歸要靠測試；瀏覽器是體驗驗證。  
- **改前端預測當成功**：本機看起來順，另一個分頁對不上——回去查權威 state。  
- **房號不一致**：一人在 `room1`、一人在 `room2`，永遠等不到開局。  

## M4／H 卷收束

你現在具備最小的「可玩」遊戲後端閉環：

> 連線 → Session → Room／Phase → 權威 Input → Tick → State 廣播 → 結束  

下一階段 **M5（I／J 卷）** 會往外長：資料落庫、觀測指標、壓測、優雅關閉——讓這間小 Arena 比較像能值班的服務，而不只是 demo。

若時間有限：先保證本 Lab 雙人同步與 `internal/game` 單測穩綠，再進 I0。
