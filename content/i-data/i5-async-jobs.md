---
lessonId: "I5"
title: "背景任務與結算：離開 tick，排好隊再寫"
description: "有界佇列、worker pool、至少一次投遞與冪等；關服時要排空，不能直接閃人。"
volume: "i"
order: 5
level: "l2"
status: "ready"
path_required: true
tags: ["async", "workers"]
example: "examples/i05-worker"
prev: "I4"
next: "J0"
---

## 這章你會搞懂什麼

對局結束後你常要做一串「重要但不必在這一幀完成」的事：

- 寫對局結果  
- 發獎勵、更新賽季分  
- 丟分析事件、打點  

這些若塞回 tick 熱路徑，房間模擬會跟 DB／網路 I/O 綁死——I1 講過的災難再上演一次。

這章的模式很樸素：

```text
Room.End → jobs <- MatchResult → worker → DB／郵件／log
```

關鍵有三：

1. **有界佇列**（bounded queue）：滿了要有背壓，不能無限 `go`  
2. **worker** 在旁邊慢慢消化  
3. **至少一次**投遞時，結算邏輯要 **冪等**（做兩次別發兩次獎）

範例：`examples/i05-worker`。

## 先跟 Python 對一下

| Python | Go |
|--------|-----|
| Celery／RQ／Huey | channel + worker pool，或外部 MQ |
| `asyncio.create_task` 丟了不管 | 明確的生命週期：`WaitGroup`、關閉、drain |
| broker 佇列長度監控 | 有界 channel 滿了立刻知道（或度量） |

Python 裡你可能習慣「丟到 Celery 就好」。Go 小服務常先用 **行程內 channel** 當佇列；流量大了再換成 Redis Stream／NATS／Kafka——但「有界、可關閉、可重試、要冪等」這組原則不變。

## 怎麼寫

最小可跑形狀（與範例同精神）：

```go
type Job struct {
	MatchID string
	Winner  string
}

func worker(id int, jobs <-chan Job, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		// 這裡假裝寫 DB／發獎
		fmt.Printf("worker=%d settled match=%s winner=%s\n", id, j.MatchID, j.Winner)
	}
}

func main() {
	jobs := make(chan Job, 8) // 有界：最多積 8 筆
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go worker(i, jobs, &wg)
	}

	// 模擬五局結束
	for i := 0; i < 5; i++ {
		jobs <- Job{MatchID: fmt.Sprintf("m%d", i), Winner: "Ada"}
	}
	close(jobs) // 不再派工，worker 會讀完離開
	wg.Wait()
}
```

要點：

- `make(chan Job, N)` 的 **N 就是背壓閥門**  
- 生產者（Room／結算點）在佇列滿時應 `select` 失敗路徑：回傳錯誤、拒開新局、或丟棄**非關鍵**分析事件——依業務選  
- `close(jobs)` 後 worker 用 `range` 自然結束；再 `Wait` 會合  

請自己跑 `examples/i05-worker`，看三個 worker 怎麼分掉五筆任務。

### 佇列滿了怎麼辦（L1 就該想）

```go
func Push(jobs chan<- Job, j Job) error {
	select {
	case jobs <- j:
		return nil
	default:
		return errQueueFull // 呼叫端決定：拒絕、改丟分析、或換策略
	}
}
```

「滿了還用無緩衝硬塞」會堵住 Room 迴圈；「滿了就開新 goroutine」則回到 C8 的無界扇出。有界＋明確錯誤，才是可控。

## 為什麼這樣設計

### 為什麼結算要冪等？

worker 可能：

- 寫 DB 成功了，但回報前進程掛了 → 重放  
- 網路逾時你以為失敗，其實對方已收到  

所以「至少一次」是常態。冪等手段例如：

- 以 `match_id` 當唯一鍵，`INSERT` 衝突就視為已結算  
- 狀態機：`pending → settled`，重複事件直接 no-op  

沒有冪等的重試＝玩家收兩次獎，客服收兩倍工單。

### 關服順序

長駐服務不要「殺行程了事」：

1. **停止接收**新的結算／新開局  
2. **drain**（排空）佇列，或把未完成落地到可重放存放  
3. 再結束 worker、關 DB  

J 卷的優雅關閉會跟這接起來；這章先建立態度：佇列裡的關鍵任務是「還沒入冊的帳」。

### Outbox（進階可先略過）

若你必須「DB 事務成功」與「訊息一定出去」綁在一起，常見作法是同一交易寫入 **Outbox 表**，再由獨立輸送把 Outbox 投遞出去。這比「先寫業務表再隨手 pubsub」不容易漏。新手先把行程內有界佇列＋冪等做熟即可。

## 遊戲 Server 會用在哪

Arena Mini 可在 `PhaseEnded`（或同等結束點）：

- 組一筆 JSON 結果  
- `Push` 進 jobs  
- worker 寫 log／之後換真 DB  

分析事件（「這局大家喜歡哪個模式」）在佇列滿時可以丟；**發獎與扣款**不行——滿了要背壓到「先別開新局／先別完成購買」這一層。

記得：Room 仍然**不在 tick 裡**做重 I/O；它只負責把 Job 丟過邊界。

## 請丟掉的舊習慣

1. **在請求／tick 裡同步寄信、寫庫 3 秒**——尾巴延遲直接餵給玩家。  
2. **無界 `go func() { settle() }` 狂丟**——洩漏與背壓全無（回顧 C1、C8）。  
3. **重試不談冪等**——等於隨機刷獎機器。  
4. **關服不 drain**——結算蒸發，還以為「重啟就好」。

## 動手練習

### 必做

1. 跑通 `examples/i05-worker`。  
2. 改成 queue 滿時 `Push` 回傳 error（`select` + `default`），並用測試或 `main` 演示滿佇列行為。  

### 選做

1. 幫 `Job` 加 `Attempt`，模擬失敗重試；用 `match_id` 集合確保同一局只印一次「settled」。  
2. 關機流程：再收一筆 Job 的開關關掉後，排空再印 `drained`。  

## 常見坑

- **緩衝開超大「先不當有界」**：延遲問題變成記憶體問題，崩得更晚更劇烈。  
- **close 後還往 channel 送**：panic；要先停生產者。  
- **worker 內再開無界任務**：問題只是被推下一層。  
- **把 pubsub 當成這章的隊列**（I3）：不可靠，不適合結算。  

## 延伸閱讀

- J3／J4：優雅關閉與負載——跟 drain、背壓同一條故事線  

## 本卷收束

I 卷要內化的不是某個 Redis 指令，而是：

| 問題 | 答案方向 |
|------|----------|
| tick 裡寫 DB 嗎？ | **不** |
| 權威對局狀態在哪？ | Room **記憶體** |
| 跨進程軟狀態／排行／TTL？ | 常是 **Redis** |
| 錢與庫存正本？ | **DB 交易** |
| 結束後的重活？ | **有界佇列 + worker**，並談冪等 |

下一卷 J：日誌、指標、關服、負載——讓這些資料路徑在生產環境跑得夠久。
