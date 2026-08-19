---
lessonId: "J3"
title: "優雅關閉：Ctrl+C 也不該直接蒸發"
description: "收到 SIGINT／SIGTERM 後：停接新連線、停 tick、排乾 in-flight，再離開——給玩家與資料一個體面結局。"
volume: "j"
order: 3
level: "l2"
status: "ready"
path_required: true
tags: ["shutdown"]
example: "examples/j03-shutdown"
prev: "J2"
next: "J4"
---

## 這章你會搞懂什麼

開發時按 `Ctrl+C`，行程死掉，好像沒關係。  
上線後，「關」常常等於：

- 滾動更新／換版  
- 機器縮容  
- 進程掛了被編排系統重啟  

若你直接 `os.Exit` 或被殺到一半：

- WebSocket 玩家瞬間斷，沒收到「伺服器維護中」  
- 房間 tick 寫到一半停在怪異 phase  
- 背景佇列（結算、寄信）丟半寫入  
- 負載平衡還在把新流量打進來

**優雅關閉（graceful shutdown）** 的目標不是「永遠不中斷」，而是：**按步驟收尾，並在超時後仍能離開**。

讀完你要能說出關閉順序，並看懂 `signal.NotifyContext` + `Server.Shutdown` 這組最小寫法。

## Python 對照

| Python | Go |
|--------|-----|
| 抓 `SIGINT`／`SIGTERM` | `signal.Notify`／`signal.NotifyContext` |
| uvicorn／gunicorn graceful | `http.Server.Shutdown` |
| `atexit` 註冊清理 | `defer` + context 取消 + 明確 Close |
| asyncio 取消任務樹 | 取消 root `context`，下游一起收 |

概念相同：先發信號，再限時清理。遊戲多了「停 tick、關房間、通知客戶端」這幾步。

## 怎麼寫

HTTP 服務的最小骨架：

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer stop()

go func() {
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}()

<-ctx.Done() // 等到 Ctrl+C 或 SIGTERM

shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
_ = srv.Shutdown(shutdownCtx)
```

要點：

1. **監聽**在背景 goroutine  
2. **主 goroutine** 卡在信號  
3. `Shutdown` 帶**超時**——不能無限等  
4. `http.ErrServerClosed` 是預期錯誤，別當成崩潰

範例：`examples/j03-shutdown`（`:8098`，有 `/healthz`）。

遊戲 Server 還要在 HTTP Shutdown 前後處理「有狀態」部分。Arena Mini 的順序大致是：

1. 收到信號，打日誌  
2. `hub.Close()`：停房間 tick 等遊戲迴圈  
3. `srv.Shutdown`：停止 HTTP 服務、等連線收尾（含超時）  
4. 打 `bye`

見 `demo/arena-mini/server/cmd/server/main.go`。

## 為什麼這樣設計／底層在幹嘛

### 建議關閉順序（人話版）

1. **對外說「別再來了」**  
   停止 `Accept` 新連線；若有就绪探針，讓編排系統不要再打流量過來。  
2. **停遊戲迴圈**  
   取消 tick、別再廣播新 state；必要時跟客戶端說 `server_closing`。  
3. **排乾 in-flight**  
   手上這幾個請求／寫庫／佇列項目，能做完就做完。  
4. **限時強制**  
   超時後 `Close`，避免卡死導致部署卡關。  
5. **flush**  
   日誌、metrics、最後一筆任務，能刷就刷。

為什麼要順序？反過來做——先硬關 listener、tick 還在跑——很容易出現「半死不活」的房間與送失敗洗版。

### `Shutdown` 跟 `Close` 差在哪？

- `Server.Shutdown`：停掉 listener，**盡量等**既有連線跑完（受你傳入的 context 超時限制）。  
- `Server.Close`：比較粗暴，立即關閉。  

WebSocket 長連線往往不會「自己結束」：你可能要在 hub 裡**主動 Close 連線**，或先取消連線用的 context，否則 `Shutdown` 會一直等到超時。

### 為什麼需要超時？

因為總有連線不走、鎖死、遠端沒回應。  
優雅是「盡力」；**可運維**要求「到點必須死掉，好讓新版本起來」。  
K8s 常見要對齊：`terminationGracePeriodSeconds`、`preStop` hook，跟你程式內超時同一數量級，否則平台先殺你、你的清理還沒跑完。

### 信號要懂哪些？

| 信號 | 常見來源 |
|------|----------|
| `SIGINT` | 本機 `Ctrl+C` |
| `SIGTERM` | Docker／K8s 停止容器 |

別把 `SIGKILL`（`kill -9`）當日常——它不能被捕捉，進程直接沒，沒有清理機會。

### 進階可先略過

- 就緒（ready）與存活（live）探針分離  
- 多模組關閉用 `errgroup` + 一個 root cancel  
- Windows 服務控制與信號差異

## 遊戲 Server 會用在哪

Arena Mini 練習劇本：

1. `go run ./cmd/server`  
2. 瀏覽器連進去開打  
3. 終端 `Ctrl+C`  
4. 看日誌：先有 shutdown signal，再 bye；遊戲迴圈應停  

若關服時還有人在房裡，理想體驗是客戶端很快知道連線結束，而不是半吊子 state 一直卡在「playing」。

跟 I 卷背景任務一起想：關服時佇列要嘛限時 drain，要嘛明確「這筆會重試／會丟」，不要默默吞。

## 請丟掉的舊習慣

1. 把 `kill -9` 當重啟手段。  
2. `main` 結束靠行程蒸發，不管 goroutine／tick。  
3. 關服中途還接受新房、新匹配。  
4. 沒有超時的「無限優雅」——部署會被你拖死。

## 動手練習

### 必做

1. 跑 `examples/j03-shutdown`，瀏覽器或 `curl` 打 `/healthz`，再 `Ctrl+C`，確認出現 `shutting down`／`bye` 這類日誌。  
2. 對 Arena Mini 做同樣的事；對照 `main.go` 裡 `h.Close()` 與 `srv.Shutdown` 的先後順序，用自己的話寫下來「為什麼先停 hub」。

### 選做

1. 把 shutdown 超時改成很短（例如 1ms）再試，觀察是否提早結束、有無 error 日誌。  
2. 設想：若要在關服前對所有連線廣播一則 `system` 訊息，應插在步驟的哪裡？

## 常見坑

- **只 `Shutdown` HTTP，忘了停 tick**：房間 goroutine 可能還在跑，或卡在送網上。  
- **在請求路徑裡呼叫會卡住的清理**：關閉邏輯要可預期、可超時。  
- **信號處理註冊太晚**：Listen 之後很久才註冊，中間的 Ctrl+C 行為會怪。  
- **把 `ListenAndServe` 的回傳錯誤一律 `Fatal`**：忘了排除 `ErrServerClosed`，關服反而看起來像崩潰。

## 延伸閱讀

- `net/http` 的 `Server.Shutdown`  
- 範例：`examples/j03-shutdown`  
- Arena Mini：`demo/arena-mini/server/cmd/server/main.go`
