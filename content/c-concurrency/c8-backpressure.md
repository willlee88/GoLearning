---
lessonId: "C8"
title: "背壓：佇列有界，壓力要回得去"
description: "生產者快、消費者慢時怎麼辦：阻塞、拒絕、丟棄、斷線——別讓無界 channel 變成延遲的 OOM。"
volume: "c"
order: 8
level: "l2"
status: "ready"
path_required: false
tags: ["backpressure"]
example: "examples/c08-backpressure"
prev: "C7"
next: "C9"
---

## 這章你會搞懂什麼

**背壓（backpressure）** 的意思很白話：後面消化不完時，要把壓力傳回前面——讓前面慢下來、被拒絕、或被丟掉——而不是默默在中間堆積到爆記憶體。

在 Go 裡最常見的陷阱是：

```go
ch := make(chan Job) // 或超大緩衝
go func(){ for j := range ch { ... } }()
// 然後無腦 ch <- job
```

生產者永遠比較快時，無界（或「心智上無界」）佇列＝**延遲發作的 OOM**。遊戲 Server 被惡意／故障客戶端灌 input 時特別真實。

## Python 對照

| Python | Go |
|--------|-----|
| `asyncio.Queue(maxsize=N)` 滿了就 await | 有界 channel；滿了阻塞或 `select`+`default` 拒絕 |
| 無限 list append | 等同無界緩衝——一樣危險 |
| 丟棄／降級常自己寫 | 一樣要**明確策略**，語言不會幫你選 |

## 怎麼寫

有界＋非阻塞拒絕：

```go
ch := make(chan Job, 1024)
select {
case ch <- job:
	// 受理
default:
	return ErrBusy // 滿了：拒絕／改丟棄／叫客戶端減速
}
```

有界＋阻塞：直接 `ch <- job`，讓生產者變慢——適合「寧可延遲、不可丟」的路徑。

範例：`examples/c08-backpressure`。

## 細節

### 策略怎麼選？

| 策略 | 什麼時候 |
|------|----------|
| 阻塞 | 可接受延遲，且生產者也是你能控的內部元件 |
| 拒絕（回錯／忙） | API／命令入口，讓呼叫端決定重試 |
| 丟最舊／最新 | 遙測、可丟的狀態快照 |
| 斷開客戶端 | 明顯輸入洪水、惡意連線 |
| 降級 | 關特效、降 tick、少廣播欄位 |

沒有萬能策略。要問：**丟了會不會破壞公平／權威狀態？延遲高一點可不可以？**

### 為什麼「開超大 buffer」常常只是自我安慰？

1024 跟 10240000 都只是延後爆點。若長期生產速率 > 消費速率，再大的緩衝也會滿。  
緩衝用來吸收**短暫尖峰**，不是用來否認產能不足。

### 跟遊戲輸入的關係

每玩家每秒輸入應有上限（例如 40）。超了就斷／拒／合併——這是玩法與安全的交界，不只是效能題。

### 進階可先略過

- Reactive Streams 那套背壓術語；概念相同。  
- 多級佇列（連線 → 房間 → worker）每一級都要有界。

## 遊戲 Server 會用在哪

- 房間 input inbox  
- 結算／寫庫的 async worker queue（後面 I5／J6 也會談）  
- 廣播：peer 太慢時，是阻塞整房、踢人，還是丟更新？要有政策  

Arena Mini 相關設計時，把「滿了怎麼辦」寫進註解或文件，比事後救火便宜。

## 請丟掉的舊習慣

1. `go func` + 無界 channel 當萬靈丹。  
2. 無限 list／slice 當緩衝。  
3. 出問題才加 buffer，從不設拒絕路徑。

## 動手練習

### 必做

1. 跑 `examples/c08-backpressure`，觀察滿佇列行為。  
2. 為 Arena 設計一條規則：「每玩家每秒最多 40 input」，寫下超限時你要阻塞、拒絕還是斷線。  

### 選做

1. 實作「丟最舊」：緩衝滿時先抽一個再塞新的（注意併發與正確性）。  

## 常見坑

- **緩衝滿導致整個房間 loop 卡在發送**：關鍵路徑用非阻塞＋政策，或隔離慢消費者。  
- **以為 reject 就要 panic**：回 `error`／關連線即可。  
- **度量缺失**：沒有 queue 長度／拒絕次數，你看不見背壓有沒有在運作。

## 延伸閱讀

- 搜尋 “backpressure” 與有界佇列設計；搭配本卷 C2／C3。  
