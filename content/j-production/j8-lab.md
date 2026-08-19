---
lessonId: "J8"
title: "Lab：關得乾淨，再壓一下"
description: "把 metrics、優雅關閉、load-client 串成一套操作：看數字、加壓、Ctrl+C，留下你對瓶頸的第一個判斷。"
volume: "j"
order: 8
level: "l2"
status: "ready"
path_required: true
tags: ["lab", "production"]
example: "examples/j04-load-client"
prev: "J7"
next: "K0"
---

## 這章你會搞懂什麼

前面 J1～J7 是零件；這章是 **把操作練順**。

你要親手走完一輪：

1. 啟動 Arena Mini  
2. 看 `/healthz`、`/metrics`  
3. 用假 client 加壓  
4. `Ctrl+C` 看優雅關閉  
5. （選）自己加掛 pprof 或做一個小改進  

做完之後，J 卷的檢查點才算「身體會」，不是「眼睛看過」。

## Python 對照

像你在 Web 專案收尾時會做的：health check、metrics 頁、壓一下、看 graceful shutdown 日誌。  
遊戲版多了長連線與房間 tick——關服時特別要確認迴圈有停。

## 怎麼寫（操作步驟）

### 1. 啟動 Server

```powershell
cd F:\GoLearning\demo\arena-mini\server
go mod tidy
go test ./...
go run ./cmd/server
```

看到 listening 日誌後：

- 遊戲頁：<http://localhost:8080/>  
- 健康：<http://localhost:8080/healthz>  
- 指標：<http://localhost:8080/metrics>  

先記下一份「安靜時」的 metrics（connections／rooms 大概多少）。

### 2. 手動玩一下（30 秒）

兩個瀏覽器（或兩個分頁）同 room、不同 name，Ready 後動一下。  
再看 `/metrics`：`connections`、`messages_in`、`messages_out`、`inputs_applied` 應有變化。

### 3. 加壓

另開終端：

```powershell
cd F:\GoLearning\examples\j04-load-client
go run . -addr http://127.0.0.1:8080 -n 30 -seconds 15 -room load
```

壓的過程中重新整理 `/metrics`（或 curl），觀察：

- `connections` 是否上去  
- `messages_*` 是否暴衝  
- Server 終端 CPU 風扇／工作管理員體感  

壓完看 client 印出的 `clients_ok`／`fail`。

可再試一輪 `-n 50`（機器吃力就降），練習「漸進」而不是一步到位。

### 4. 優雅關閉

回到 Server 終端，`Ctrl+C`。

預期（對照 J3／`main.go`）：

1. 日誌出現 shutdown signal 一類訊息  
2. hub／房間迴圈先停  
3. HTTP `Shutdown`  
4. `bye`  
5. 行程結束（理想上乾淨離開）

若還有瀏覽器掛著，應感到連線結束，而不是半吊子一直轉。

### 5.（選）pprof

若你已熟悉 D4，可自行在 Server 加掛 pprof（注意：**別把開放的 pprof 暴露在公網**）。  
本機壓測時抓 CPU profile，驗證你在 J4／J6 的瓶頸猜測。

## 為什麼要整套做一遍

分開讀章節時，你以為自己會了。  
一連操作才會碰到真實摩擦：

- 開錯目錄、addr 填錯  
- metrics 看錯時點（壓測前就刷新了）  
- 關服順序與日誌對不上想像  
- 假人失敗其實是 Server 沒起來  

這些摩擦正是上線值班的縮小版。

## 遊戲 Server 會用在哪

把這套步驟當成你日後任何小專案的「最小上線儀式」：

| 步驟 | 問題 |
|------|------|
| healthz | 進程活著嗎？ |
| metrics | 現在負荷怎樣？ |
| load-client | 臨界點大概在哪？ |
| Ctrl+C | 關得乾不乾淨？ |

K 卷自評也會勾這些項。

## 請丟掉的舊習慣

1. Lab 用眼睛掃過指令就算完成。  
2. 只壓測不看 metrics。  
3. 關服直接關終端視窗，不看日誌順序。  
4. 一次改五個優化，卻說不出哪個有效。

## 動手練習

### 必做

1. 完整跑完：啟動 → metrics → load-client → Ctrl+C。  
2. 截一段（或複製）壓測中的 metrics JSON，保存下來。  
3. 寫下你觀察到的**第一個瓶頸猜測**（一句話即可，例如「同房廣播讓 messages_out 比 in 兇太多」）。

### 選做

1. 為 input 加每秒速率限制（呼應 J5），再用 load-client 看行為變化。  
2. 對局 `ended` 時丟一筆結算到背景 worker（呼應 I5）。  
3. 反覆「啟動→壓→關」三次，留意是否有 goroutine／房間數洩漏跡象（connections／rooms 關完後是否回到安靜基線）。

## 檢查清單（勾給自己看）

- [ ] SIGINT 後進程有明確結束（不是卡死）  
- [ ] metrics 能反映連線與訊息變化  
- [ ] 關閉後 tick／房間迴圈有停（日誌或行為上可觀察）  
- [ ] 壓測有記下數字與瓶頸猜測  
- [ ] （選）沒有明顯「越開關 metrics 基線越高」的洩漏感  

## 常見坑

- **在錯的 `server` 目錄跑**：請 `cd` 到 `demo/arena-mini/server`。  
- **load-client 的 `-addr` 要用 HTTP base**：例如 `http://127.0.0.1:8080`（見範例與 Arena Mini README）。  
- **壓測時開太多瀏覽器分頁**：自己也佔連線，數字會難讀。  
- **還沒看安靜基線就加壓**：失去對照，看不出變化。

## 延伸閱讀

- `demo/arena-mini/README.md`  
- `examples/j04-load-client`  
- J1～J7 任一章覺得虛的，回頭補那一塊即可  

## 往 K 卷

J 卷到這裡收工。  
K 卷是整條主路徑的**複習地圖、自評清單與延伸任務**——把「學過」變成「交得出作品」。
