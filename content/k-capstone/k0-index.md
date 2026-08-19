---
lessonId: "K0"
title: "K 卷導讀：收束——你走到哪了，下一步去哪"
description: "不灌新理論：用複習地圖把主路徑串起來，再進自評清單與延伸任務，把理解變成作品。"
volume: "k"
order: 0
level: "l1"
status: "ready"
path_required: true
tags: ["capstone"]
example: "demo/arena-mini"
prev: "J8"
next: "K1"
---

## 這章你會搞懂什麼

若你一路從 P0 走到 J8，你已經不是「只會印 Hello」的狀態了。  
你手上有一條完整主路徑：

- 用 Go 的方式想型別、錯誤、套件  
- 用併發把連線與房間跑起來，並知道 race 可怕  
- 用網路與協定把客戶端／Server 接上  
- 用權威 tick 做可玩的小對局  
- 用日誌、metrics、關閉與壓測讓服務可操作  

**K 卷不灌新大理論。**  
這卷要幫你：**收束、自評、選一塊做成作品集。**

讀完 K0，你要知道該回顧哪些章；下一章 K1 給你可勾的清單與延伸任務說明。

## Python 對照

像你學完一套 Django／FastAPI 教程之後的「總複習 + 自己做個專案」階段：  
不是再塞第十個套件，而是問「我真的能獨立做一個小後端嗎？」

Go 這邊的錨點就是 **Arena Mini**（`demo/arena-mini`）：規則不大，但主路徑概念都壓進去了。

## 建議複習地圖

時間有限就按這個順序「掃」，覺得虛的再深讀：

| 順序 | 回去哪 | 你要能回答 |
|------|--------|------------|
| 1 | P0 心智差異 | 為什麼遊戲後端常選 Go？錯誤／併發跟 Python 差在哪？ |
| 2 | A 卷基礎 | module、slice／map、指標、介面——寫小程式不卡死 |
| 3 | B 卷錯誤 | `error`、`errors.Is`、wrap；別用例外心智硬套 |
| 4 | C 卷併發 | goroutine 生命週期、channel、context、Mutex；**能看懂 `-race`** |
| 5 | F／G 網路與信封 | TCP／HTTP／WS 心智；JSON 信封怎麼演進 |
| 6 | H 權威 tick | 房間狀態機、input≠狀態、20Hz 在幹嘛 |
| 7 | I 資料邊界 | 什麼進記憶體、什麼進 DB、背景結算 |
| 8 | J 關閉與 metrics | 能關、能看、能壓、能守 |

實作錨點：

```powershell
cd F:\GoLearning\demo\arena-mini\server
go test ./...
go run ./cmd/server
```

再配 J8 的 load-client 走一輪，記憶會牢很多。

## 怎麼使用 K 卷

1. 用上面地圖快速自問，卡關就回對應章。  
2. 打開 K1，把清單逐項勾。  
3. 從延伸任務挑 **一個** 做完（比挑五個都做一半有用）。  

## 遊戲 Server 會用在哪

Capstone 的意思是：你能把「課」收成「可展示的後端行為」。  
例如錄一段：雙人進房 → 碰撞得分 → 結束回大廳 → `/metrics` → Ctrl+C 關乾淨。  
這段影片或 README，比空口說「我學過 Go」有力。

## 請丟掉的舊習慣

1. 複習時從頭精讀每一字——會累死；用問題導向。  
2. 自評全勾但從來沒跑過 Arena Mini。  
3. 延伸任務一次開太多，最後交不出來。

## 動手練習

### 必做

1. 依複習地圖，標出你最不熟的兩塊。  
2. 確認本機能啟動 Arena Mini，打開遊戲頁與 `/metrics`。  

### 選做

1. 用一頁紙畫：Client → Hub／Room → tick → state 廣播；旁註 metrics 與 shutdown 插在哪。

## 常見坑

- **以為 K 卷會教全新架構**：不會；這裡是收束。  
- **跳過 J8 直接勾 J 相關清單**：操作手感會假。  
- **只讀不跑**：Go 的手感建立在 `go run`／`go test` 上。

## 下一章

K1：自評清單（含每項在檢查什麼）與延伸任務說明。
