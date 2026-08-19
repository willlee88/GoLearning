---
lessonId: "P0.6"
title: "併發先建立直覺（詳情在 C 卷）"
description: "goroutine 不是 OS thread，也不是 asyncio 的直接翻譯。先搞懂：執行、通訊、多路等待，以及亂開會怎樣。"
volume: "p0"
order: 6
level: "l2"
status: "ready"
path_required: true
tags: ["python-bridge", "concurrency"]
example: ""
prev: "P0.5"
next: "P0.7"
---

## 這章你會搞懂什麼

Go 談「一次做很多事」（concurrency）時，你常會看到這組工具：

- **goroutine** — 輕量的執行單元（用 `go` 關鍵字啟動）  
- **channel** — 用來通訊、交接資料的管線  
- **select** — 同時等好幾條路，哪個先好用哪個  
- 再加上 **`sync`**（鎖等）與 **`context`**（取消、截止時間）  

你可以先把遊戲 Server 預設畫面想成：「每個連線一個讀寫迴圈」——但**搞懂生命週期與誰擁有資料之前，不要狂開 goroutine**。

這章是**預告**，細節在 **C 卷**。這裡只建立正確直覺，避免你用錯模型。

## 先跟 Python 對一下

| Python | Go（粗對照，**不要硬套**） |
|--------|---------------------------|
| `threading.Thread` | goroutine（更輕，模型也不同） |
| `asyncio` 協程 + event loop | 多個 goroutine；由 runtime 排程 |
| `queue.Queue` | channel |
| 取消／超時常見寫法 | `context.Context` |
| GIL 限制 CPU 平行 | 可以真正平行；**更要注意**共享記憶體 |

重點：把 `asyncio.create_task` 的感覺原封貼成到處 `go func`，之後會洩漏、會 race、會關不掉。

## 怎麼寫（能跑的最小例子）

啟動一個併發執行的函式：

```go
go func() {
	fmt.Println("runs concurrently")
}()
```

用 channel 送收資料（這裡用帶緩衝的 channel，比較好演示）：

```go
ch := make(chan string, 1)
ch <- "ping"
msg := <-ch
fmt.Println(msg)
```

注意：若 `main` 函式結束，整個程式結束，**其他 goroutine 會被直接砍掉**。所以你一定要設計「怎麼等待／怎麼結束」（`WaitGroup`、channel 會合、`context` 取消……C 卷會系統教）。

## 為什麼這樣設計／底層在幹嘛

### Goroutine 是 M:N 排程

很多個 goroutine，由 Go runtime 排到較少的作業系統執行緒上跑。所以它很輕，但**不是免費**：

- 每個都有堆疊與狀態  
- 開太多、又不結束 → 記憶體與排程壓力會炸  
- 共享 map／slice 卻不同步 → **data race（資料競爭）**，行為未定義，是 bug  

用 `go test -race` 可以抓很多這類問題——養成習慣。

### Channel 不只是 queue

無緩衝 channel 有**交接同步**的味道：送跟收要對上。關閉 channel 有「誰有權關閉」的規則；亂關會 panic。細節留給 C 卷，但你先記住：**通訊也是設計，不是隨手扔資料。**

### 口號別當教條

常聽到：「不要通過共享記憶體來通訊；通過通訊來共享記憶體。」這是好預設，但熱路徑上用鎖（mutex）仍然完全合理。重點是：**同一塊資料的所有權與同步策略要想清楚**。

寫錯會怎樣？（預告，後面會實驗室打臉）

- 兩個 goroutine 無鎖寫同一個 `map` → 可能 panic，或更糟：偶發壞資料  
- `main` 太早結束 → 背景工作根本沒跑完  
- 無界地 `go` + 無界佇列 → 背壓沒了，記憶體爆炸  

### 進階可先略過

- 為什麼鎖在熱路徑仍合理：通道不是萬能，拷貝與調度也有成本。  
- 排程器細節：C 卷與進階章再看。  

## 遊戲 Server 會用在哪

典型拆分（先混個臉熟）：

- **每連線**：讀取 goroutine + 寫入要串行化（避免交錯寫）  
- **每房間**：一個邏輯擁有者（專責 goroutine，或鎖保護的 tick）  
- **大廳／registry**：多人同時進房離房，要有明確的併發策略  

之後你會一直問：「這份狀態誰能寫？誰讀？何時結束？」——這比「我會不會寫 `go`」重要得多。

## 請丟掉的舊習慣

1. 把 `asyncio.create_task` 心智原封貼成到處 `go func`。  
2. 假設「不會平行就不會 race」——在 Go 預設就**可能**平行。  
3. 開了背景工作卻不定義如何停止 → goroutine 洩漏，關機關不乾淨。  

## 動手練習

### 必做

1. 寫一個小程式：主流程送 3 個訊息到 channel，另一個 goroutine 印出它們；想清楚 **main 如何等待** 對方跑完。  
2. 用文字描述：若兩個 goroutine 無鎖寫同一個 `map`，可能發生什麼。  

### 選做

1. 對任何含併發的小實驗跑 `go test -race`。  

## 常見坑

- **main 結束殺光所有 goroutine**：用 `WaitGroup` 或 channel 會合。  
- **無界 `go` + 無界工作佇列**：記憶體爆炸——後續會講背壓（backpressure）。  
- **以為 channel 取代所有鎖**：該鎖還是要鎖；混用更要規矩清楚。  

## 延伸閱讀

- <https://go.dev/blog/pipelines>  
- 卷 C 全卷（正式開打併發）  
