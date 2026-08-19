---
lessonId: "C1"
title: "goroutine：開得起，也要關得了"
description: "go 關鍵字很便宜，但不是免費；學會 WaitGroup、停止條件，以及常見洩漏長什麼樣。"
volume: "c"
order: 1
level: "l2"
status: "ready"
path_required: true
tags: ["goroutine"]
example: "examples/c01-goroutines"
prev: "C0"
next: "C2"
---

## 這章你會搞懂什麼

`go f()` 會啟動一個 **goroutine**——Go runtime 排程的輕量執行單元。  
它比作業系統執行緒便宜得多，所以網路服務常「一連線一（或一組）goroutine」。

但便宜 ≠ 免費：每個沒結束的 goroutine 都佔記憶體與排程資源。更重要的是生命週期：

- **`main` 結束，整個程式結束**，背景工作會被直接砍掉  
- 你必須定義：**什麼條件停止、誰負責 Wait 會合**

讀完這章，看到 `go func()` 你要反射性問：「它什麼時候結束？誰在等它？取消信號從哪來？」

## Python 對照

| Python | Go |
|--------|-----|
| `threading.Thread` | goroutine（更輕，模型也不同） |
| `asyncio.create_task` | 也是「丟去跑」，但取消／排程語意不同 |
| daemon thread | 沒有 1:1 等價；靠程序退出與明確取消 |
| `thread.join()` | `WaitGroup.Wait`／從 channel 收「做完了」 |

Python 裡你可能 fire-and-forget 一個 daemon thread。在 Go 長駐 Server 裡，這種習慣很容易變成：**洩漏的 goroutine 越積越多，直到記憶體與調度一起垮**。

## 怎麼寫

最常見的會合方式是 `sync.WaitGroup`：

```go
var wg sync.WaitGroup
wg.Add(1)
go func() {
	defer wg.Done()
	work()
}()
wg.Wait()
```

要點：

- `Add` 在開 goroutine **之前**  
- `Done` 用 `defer`，避免早退路徑忘記  
- `Wait` 在確定要會合的地方  

範例：`examples/c01-goroutines`。試著拿掉 `Wait`，看程式為什麼「好像沒做完就結束」。

## 細節

### 為什麼 main 一結束全部蒸發？

因為行程都沒了。單元測試或小 demo 裡這很常見：你以為背景有在跑，其實 `main` 已經 return。  
長駐服務的 `main` 通常會卡在：聽訊號、`Wait`、或伺服器的 `Listen` 迴圈。

### 常見洩漏長怎樣

| 症狀來源 | 典型原因 |
|----------|----------|
| 永遠卡在 channel 接收 | 沒人再寫、也沒 close |
| 永遠卡在傳送 | 沒人接收，又沒有緩衝／退出路徑 |
| `WaitGroup` 少 `Done` | Wait 永久卡住（或反過來 Add 不對） |
| 開了 worker 沒有取消 | 連線斷了 loop 還在 |

洩漏不一定立刻炸；它比較像慢性病。

### 停止條件從哪來？

後面幾章會系統講：

- channel 關閉／哨兵值  
- `context` 取消（C4）  
- select 多路退出（C3）  

這章先建立態度：**開之前就要想停。**

### 進階可先略過

- goroutine 啟動時會捕獲當前變數；迴圈裡開 goroutine 要注意閉包（Go 1.22 前後細節不同）。  
- 堆疊可成長，但仍有上限；無限遞迴一樣會炸。

## 遊戲 Server 會用在哪

典型形狀：

- 每條連線一個讀 loop goroutine  
- 斷線 → 讀失敗或 ctx 取消 → **退出 loop** → `unregister`  
- 房間 worker 一個（或少數）長壽命 goroutine，而不是每個 tick 開新的  

「連線還在不代表工作健康」——要以 loop 是否退出、資源是否釋放為準。

## 請丟掉的舊習慣

1. fire-and-forget，不管生命週期。  
2. 用 `time.Sleep` 假裝同步。  
3. 以為「行程還活著＝背景任務都正常」。

## 動手練習

### 必做

1. 跑 `examples/c01-goroutines`。  
2. 故意拿掉 `Wait`，觀察提早退出；再改回來。  

### 選做

1. 開 3 個 worker，用一個 `WaitGroup` 等全部結束，最後印一句 `all done`。  

## 常見坑

- **`Add` 寫在 goroutine 裡面**：可能 `Wait` 先跑到，計數還是 0。  
- **對同一個 WaitGroup 複製後再 Wait／Done**：別複製 `WaitGroup`（含它的 struct）。  
- **在請求路徑無限 `go`**：沒有上限的扇出是背壓與洩漏的溫床（C8）。

## 延伸閱讀

- <https://go.dev/blog/pipelines>  
