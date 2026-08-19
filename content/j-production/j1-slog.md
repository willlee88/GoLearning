---
lessonId: "J1"
title: "結構化日誌 slog：出事時要找得到人"
description: "用 log/slog 打帶欄位的日誌；房間、玩家、請求能關聯，別再靠人肉 parse 字串。"
volume: "j"
order: 1
level: "l2"
status: "ready"
path_required: true
tags: ["logging"]
example: ""
prev: "J0"
next: "J2"
---

## 這章你會搞懂什麼

平常 `fmt.Println` 或 `log.Printf("player %s joined room %s", ...)` 看起來夠用。  
一到線上，你會想問：

- 這個錯誤是哪個房間？哪個玩家？哪次請求？  
- 同一秒幾千行，怎麼過濾？  
- 告警系統能不能靠「欄位」判斷，而不是用正則硬刮字串？

Go 1.21+ 標準庫的 **`log/slog`**（structured logging，結構化日誌）就是為這個設計的：每筆日誌除了訊息文字，還帶 **鍵值屬性**（例如 `room`、`player`、`tick`）。

讀完這章你要能：

1. 把一句普通 log 改成帶欄位的 `slog`  
2. 知道等級（Debug／Info／Warn／Error）大概怎麼分  
3. 知道「高頻路徑」為什麼不該狂打 Info（那是 metrics 的活）

## Python 對照

| Python | Go |
|--------|-----|
| `logging` 模組 | `log`／`log/slog` |
| `structlog`、`extra={...}` | `slog.String("room", id)` 這類屬性 |
| JSON formatter 給 ELK／CloudWatch | `slog.NewJSONHandler(...)` |
| 習慣 `print` 除錯 | 生產請別留；開發也建議早點用 slog |

心智很像：你不是在「印一句話」，你是在「寫一筆可查的事件」。

## 怎麼寫

最常見寫法：

```go
package main

import "log/slog"

func playerJoined(roomID, name string) {
	slog.Info("player joined",
		"room", roomID,
		"player", name,
	)
}
```

訊息是人類讀的短句；後面的鍵值才是機器與搜尋用的。

若要 JSON 輸出（收集系統比較好吞），可以在 `main` 設一次預設 logger：

```go
log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelInfo,
}))
slog.SetDefault(log)

log.Info("arena-mini listening", "addr", ":8080", "version", "m5.1")
```

Arena Mini 的 `demo/arena-mini/server/cmd/server` 就是這種接法：stdout 丟 JSON，關服、聽 port 都會帶欄位。

也可以綁固定屬性，之後每行自動帶上：

```go
roomLog := slog.Default().With("room", roomID)
roomLog.Info("tick started", "hz", 20)
roomLog.Warn("player timeout", "player", name)
```

## 為什麼這樣設計／底層在幹嘛

### 為什麼要「結構化」，不要只靠字串插值？

字串長這樣時：

```text
player bob joined room X7
error: write failed for bob
```

人還看得懂；機器要可靠抽出 `player`／`room` 就痛苦。欄位一分離：

```json
{"msg":"player joined","room":"X7","player":"bob","level":"INFO"}
```

過濾 `room=X7`、統計某玩家錯誤率，都變簡單。遊戲 Server 一忙起來，**可搜尋**比「好看」重要。

### 日誌等級怎麼用（新手版）

| 等級 | 什麼時候 |
|------|----------|
| Debug | 開發細節、極吵的東西；生產預設通常關掉 |
| Info | 正常但重要的生命週期：啟動、進房、關服開始 |
| Warn | 可恢復的怪事：重試、接近上限、非法但已拒絕的輸入 |
| Error | 真的失敗、需要人看：寫庫失敗、非預期 panic 路徑 |

原則：**Info 不該刷到你來不及捲動。** 每 tick 打一次 Info，20Hz 的房間會把自己淹死。

### 日誌 vs Metrics

| | 日誌 | Metrics |
|--|------|---------|
| 回答 | 「發生了什麼事／給誰」 | 「現在多快、多少、多錯」 |
| 適合 | 事件、排查單次問題 | 連線數、訊息速率、延遲分位 |
| 高頻 | 容易爆量、爆錢、爆磁碟 | 設計上就是高頻聚合 |

tick、訊息速率這類東西，優先看 J2 的 metrics；日誌留給「進房失敗原因」「關服步驟」這種敘事。

### 安全：千萬別打秘密全文

Token、密碼、cookie、完整授權標頭，**不要**進日誌。  
就算內部系統，日誌常會進集中平台、被更多人看見。打錯誤時用「驗證失敗」「token 無效」就夠；必要時只留後幾碼或雜湊前綴。

### 進階可先略過

- 自訂 `Handler`、把 trace id 從 `context` 取出來自動附加  
- 抽樣（sampling）：極高頻 Warn 只留 1%  
- 多 handler 分流（檔案 + stdout）

主路徑先會：`Info`／`Error` + 關鍵欄位 + JSON handler。

## 遊戲 Server 會用在哪

典型欄位（依你的協定調整）：

- `room`／`room_id`  
- `player`／`session`  
- `tick`（出錯當下的 tick）  
- `phase`（lobby／playing／ended）  
- `err`（錯誤物件；slog 會處理 `error` 型別）

劇本舉例：

1. 玩家回報「進不去房」→ 用 `room` + `Warn`／`Error` 找拒絕原因  
2. 關服後有人說資料怪異 → 對時間線看 `shutdown signal` 與最後幾筆 `player` 事件  
3. 某房訊息暴增 → 先用 metrics 確認速率，再用該 `room` 的日誌看是否有人狂送非法 input

Arena Mini M5 已用 slog；你在 J8 lab 關服時會看到 JSON 行。

## 請丟掉的舊習慣

1. 生產殘留 `print`／臨時 `log.Printf` 除錯。  
2. 只有「失敗了」幾個字，沒有 room／player／err。  
3. 每個 tick、每筆 input 都打 Info。  
4. 為了方便把 token 整串印出來。

## 動手練習

### 必做

1. 打開 `demo/arena-mini/server/cmd/server/main.go`，找出 `slog.NewJSONHandler` 與關服相關的 `log.Info`／`log.Error`。  
2. 在你自己的小程式（或房間進房路徑）把一句 `log.Printf` 改成 `slog.Info`，並帶上 `room`（或同等欄位）。跑起來確認輸出看得到該鍵。

### 選做

1. 用 `.With("room", id)` 做一個 room 專用 logger，進房／退房各打一筆。  
2. 暫時把 `Level` 調成 `slog.LevelDebug`，感受一下「多吵」；再改回 Info。

## 常見坑

- **鍵名不統一**：有時 `room`、有時 `room_id`、有時中文鍵——搜尋會哭。團隊先約一小套英文鍵。  
- **把很大的 payload 整包打進日誌**：快照、整房 state 塞進去，磁碟與個資風險一起爆。打摘要或 id 即可。  
- **用日誌代替監控**：你以為「有打 log 就有觀測」；真要看曲線與告警，請接 J2。  
- **Error 裡吞掉 `err`**：只寫 `"failed"` 不帶 `err`，之後幾乎無法除錯。

## 延伸閱讀

- <https://pkg.go.dev/log/slog>  
- Arena Mini：`demo/arena-mini/server/cmd/server/main.go`
