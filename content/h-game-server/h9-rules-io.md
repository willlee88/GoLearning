---
lessonId: "H9"
title: "規則與 I/O 分離：game 不該 import websocket"
description: "純邏輯可單測；hub 只負責收發。時間與隨機數可注入——讓 go test 不需開 port。"
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

## 這章你會搞懂什麼

前面章節幾乎都在暗示同一件事，這章把它說死：

```text
hub / ws     解碼信封、找房間、回寫連線
      ↓
room.PushInput / Ready / …
      ↓
game.Apply / Tick / Snapshot     ← 純邏輯
      ↓
hub          編碼、廣播
```

**規則層（`game`）不要 import `net`、不要 import `websocket`。**

好處非常實際：

- `go test` 不用聽 port、不用真瀏覽器  
- 規則改動有表驅動測試護著  
- 之後換傳輸（TCP、不同 WS 庫）不會撕開整個玩法

讀完你要能指出 Arena Mini 裡「純邏輯檔」的路徑，並為 `Apply`／`Tick` 寫一個不開 WS 的測試腦內草案。

## Python 對照

這就是小型的**六角形／洋蔥／清潔架構**精神——名字不重要：

| 層 | 責任 |
|----|------|
| 介面適配器 | FastAPI 路由、WS endpoint |
| 領域邏輯 | `apply_input(state, cmd) -> state` |
| 基礎設施 | DB、Redis、真實時鐘 |

Python 常見痛：測試裡 `TestClient` 起一整個 app 才測得到「往右走一步」。  
Go 若一開始就把 `Apply` 做成普通函式／方法，測試只是普通的 `go test`。

## 怎麼寫

### 規則層：只談資料結構

```go
// game/room.go — 純邏輯（示意）
package game

func (r *Room) PushInput(in Input) error { /* 校驗 + 入隊 */ }
func (r *Room) Tick()                   { /* 套用 inbox */ }
func (r *Room) Snapshot() State         { /* 回傳自洽快照 */ }
```

參考實作：`examples/h06-apply-input`（沒有網路依賴）。

### I/O 層：翻譯進出

```go
// hub — I/O（示意）
func (h *Hub) onMessage(sess *Session, env Envelope) {
	r := h.rooms[sess.RoomID]
	switch env.Type {
	case "input":
		in, err := parseInput(env.Payload)
		if err != nil { /* 回 error */ return }
		if err := r.PushInput(in); err != nil { /* 回 error */ return }
	case "ready":
		_ = r.Ready(sess.Name)
	}
}

// tick loop 裡：
snap := r.Snapshot()
h.broadcast(sess.RoomID, "state", snap)
```

hub 可以依賴 game；**game 不可以依賴 hub**。依賴方向只准單向。

## 為什麼這麼執著

### 1. 測試速度與穩定

開 port 的測試會：偶發佔用、CI 較慢、失敗訊息難讀。  
純邏輯測試可以一秒跑完「Lobby → Ready → Input → Tick → 斷言座標」。

### 2. 避免「規則簽名帶 websocket」

一旦出現：

```go
func ApplyDamage(conn *websocket.Conn, amount int)
```

你就已經把傳輸細節焊進玩法。之後：

- 無法在無網環境測  
- 觀戰／Bot／重放都難接  
- 併發模型被連線生命週期綁架  

### 3. 可注入時間與隨機

單測最恨真正的 `time.Now()` 與全球隨機：

```go
type Clock interface{ Now() time.Time }

// 測試塞假時鐘：開局超時、冷卻結束變得可斷言
```

隨機數同樣：注入 `*rand.Rand` 或 seed，讓「抽牌／暴擊」可重現。

### 進階可先略過

- 介面只為了測試而長得到處都是——先抽真的痛點（時間、隨機、ID 產生）。  
- 重放系統：把 input log 存下來，對同一 `game` 再跑一次 Tick——分層後幾乎免費得到。  
- 效能熱點若在 Snapshot 複製，仍屬規則／同步議題，不是讓你把 WS 寫回 game 的理由。

## 遊戲 Server 會用在哪

請打開：

- `demo/arena-mini/server/internal/game` — 規則與房間狀態  
- `demo/arena-mini/server/internal/hub` — 連線、路由訊息、廣播  

自問：

1. `game` 有沒有 import 網路套件？  
2. 不開 server，`go test ./internal/game` 能不能測移動與開局？  

H10 Lab 的檢查清單會再盯一次這條邊界。

## 請丟掉的舊習慣

1. **單元測試裡真的 `Listen` port** 才能斷言「往右一步」。  
2. **規則函式參數帶上 `websocket`／`http.ResponseWriter`**。  
3. **為了「方便」，在 game 裡直接 `conn.WriteJSON`**——方便是假的，耦合是真的。  
4. **測試只靠手動開兩個瀏覽器**——那叫煙霧測試，不能取代規則單測。

## 動手練習

### 必做

1. 指出 demo 中純邏輯相關檔案路徑（至少一個 `internal/game` 下的檔）。  
2. 為 `Apply`／`PushInput+Tick` 寫一則**不開 WS**的測試構想（Given／When／Then 三行即可）；有餘力就對 `examples/h06-apply-input` 再加一則。  

### 選做

1. 把「超時結束對戰」的時間改成可注入 Clock，寫測試：假時間一跳就 Ended。  
2. 檢查 `internal/game` 的 import 清單，確認沒有網路套件混入。  

## 常見坑

- **package 名叫 `game`，內容卻在寫 hub**：名稱無法拯救依賴。  
- **測試與正式碼複製兩份規則**：改一處漏一處——請直接測同一份 `game`。  
- **Snapshot 回傳內部 map 的引用**：測試或 hub 一改，權威狀態被外面踩壞；必要時回傳複本。  
- **循環 import**：`game` → `hub` → `game`；通常代表邊界畫錯。

下一章 H10：把 H1–H9 收成 Arena Mini 權威對戰 Lab——這卷的檢查點。
