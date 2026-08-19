---
lessonId: "J4"
title: "壓測：先有假設，再漸進加壓"
description: "用假 client 拉高連線與訊息量；對照 metrics，找出是廣播、鎖還是編碼先爆——別一次開滿只看 crash。"
volume: "j"
order: 4
level: "l2"
status: "ready"
path_required: true
tags: ["loadtest"]
example: "examples/j04-load-client"
prev: "J3"
next: "J5"
---

## 這章你會搞懂什麼

「感覺應該沒問題」在容量上幾乎沒用。  
**壓測（load test）** 是用可重複的方式問：

> 在某某場景下，系統會先在哪裡壞？壞的時候長什麼樣？

好的壓測有三件事：

1. **假設**（例如：單進程 200 連線、同房、20Hz 廣播，CPU 會不會打滿？）  
2. **漸進加壓**（20 → 50 → 100，而不是瞬間 5000）  
3. **觀察**（metrics、延遲、錯誤率、必要時 pprof）

讀完你要能跑 `examples/j04-load-client`，並說出你在 `/metrics` 看到了什麼。

## Python 對照

| Python | Go |
|--------|-----|
| locust、k6、自寫 asyncio client | 自寫假 client（本站範例）／k6 等 |
| 常測 HTTP QPS | 遊戲還要測 **長連線 + 訊息速率 + 狀態是否仍正確** |
| 一次打滿看紅字 | 一樣是壞習慣 |

Go 寫假 client 很輕：幾千個 goroutine 模擬連線，本機就能玩出有感曲線。

## 怎麼寫

本站假 client：

```powershell
cd F:\GoLearning\examples\j04-load-client
go run . -addr http://127.0.0.1:8080 -n 20 -seconds 10 -room load
```

參數意思（對照範例原始碼）：

| 參數 | 意義 |
|------|------|
| `-addr` | Arena Mini 的 HTTP 基底 URL（會轉成 `ws://.../ws`） |
| `-n` | 假人數 |
| `-seconds` | 持續秒數 |
| `-room` | 進哪個房間 |

它會連上 WS、送 `ready`、再狂送 `input`，最後印大概的成功／失敗／訊息操作數。

建議流程：

1. 終端 A：啟動 Arena Mini（`demo/arena-mini/server` → `go run ./cmd/server`）  
2. 瀏覽器打開 `http://localhost:8080/metrics`（或 curl）  
3. 終端 B：跑 load-client，從 `-n 10` 起跳  
4. 記錄 `connections`、`messages_in`／`messages_out`、CPU 體感  
5. 再加大 `-n`，看哪裡先歪

## 為什麼這樣設計／底層在幹嘛

### 為什麼要有「假設」？

沒有假設的壓測，最後只會得到：「哦，炸了。」  
有假設，你才能驗證或推翻，例如：

- 「瓶頸在 JSON 編碼」→ 看 CPU／pprof 是否卡在 encode  
- 「瓶頸在廣播 O(N²)」→ `messages_out` 隨人數二次感上升（J6）  
- 「只是本機 loopback 太快，不代表公網」→ 換環境再測

### 「連得上」≠「狀態正確」

假人可能：

- TCP／WS 握手成功，但進房失敗  
- 一直連著，但 phase 卡死  
- 伺服器還活著，可是 tick 嚴重落後  

所以要同時看：錯誤日誌、metrics、抽樣進房玩一下手感。只看「行程沒崩潰」會過度樂觀。

### 漸進加壓為什麼重要？

一次開滿：

- 分不清是「連線風暴」還是「穩態負載」問題  
- 客戶端自己先報錯，污染結果  
- 找不到臨界點（從哪一個 N 開始崩）

漸進才能畫出「還好 → 勉強 → 壞掉」的曲線。

### 本機 loopback 的限制

先在本機測，適合找 CPU、鎖、分配、廣播演算法問題。  
它**不代表**：真實 RTT、無線網路、跨區、TLS 終結、代理緩衝。  
所以結論要寫清楚：「本機 200 連線 OK」≠「公網 200 沒問題」。

### 跟 pprof／`-race` 的關係

- 找熱點：用 D4 學的 pprof（CPU／heap）  
- **`-race` 不要開在壓測熱路徑**：它很慢，會讓你量到的是「偵測器開銷」而不是真實容量  
- race 偵測放在測試與較小場景；壓測放「代表性負載」

### 進階可先略過

- 分散式壓測（多台產生器）  
- 錄放流量、混沌工程  
- SLO：例如 99 分位延遲、斷線率上限

## 遊戲 Server 會用在哪

對 Arena Mini：

- 同房很多人 → 觀察廣播與 tick 是否跟得上  
- 多房各一小撮人 → 另一種資源形狀（比較多房間物件、比較少超大廣播）  

上線前活動：用「預期同時在線 × 安全係數」當場景，至少跑過一輪，並留下 metrics 截圖或數字。

## 請丟掉的舊習慣

1. 一次開滿，只看有沒有 crash。  
2. 沒有基線數字就開始「優化」。  
3. 只測握手，不測穩態訊息。  
4. 把本機結果直接寫進對外承諾。

## 動手練習

### 必做

1. 啟動 Arena Mini，跑 load-client：`-n 20 -seconds 10`。  
2. 壓測前後各看一次 `/metrics`，記下 `connections`、`messages_in`、`messages_out`。  
3. 用一句話寫你的第一個瓶頸猜測（CPU？廣播？還是 client 自己先炸？）。

### 選做

1. 同一組參數跑兩次，看數字是否大致可重複。  
2. 把人都塞進同一個 `-room`，再改成每人不同 room（若你改 client），比較 `messages_out` 曲線直覺。

## 常見坑

- **Server 沒開就壓**：client 全 fail，你卻以為 Server 性能差。  
- **看錯 addr**：範例預設是 HTTP base（`http://127.0.0.1:8080`），不是隨便填 `ws://` 路徑就完事——以 `examples/j04-load-client` 的旗標為準。  
- **終端被日誌刷爆**：access log 太吵會拖慢；Arena Mini 對 `/metrics`、`/healthz` 有略過，但 WS 量大時仍要注意觀測成本。  
- **把假 client 與 Server 綁在同一顆過熱 CPU 上搶資源**：解讀結果時要知道產生器自己也吃 CPU。

## 延伸閱讀

- D4：benchmark 與 pprof  
- 範例：`examples/j04-load-client`  
- Arena Mini README 的壓測段落
