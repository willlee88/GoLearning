---
lessonId: "B3"
title: "panic 跟 recover：什麼能炸、在哪接住"
description: "panic 不是業務例外；學會只在邊界 recover，避免一條連線拖垮整個行程。"
volume: "b"
order: 3
level: "l2"
status: "ready"
path_required: true
tags: ["panic", "reliability"]
example: ""
prev: "B2"
next: "B4"
---

## 這章你會搞懂什麼

`panic` 代表：**程式認為自己已經處於不該繼續的狀態**（不變量被破壞），不是「使用者又送了壞招」那種日常錯。  
日常錯請回傳 `error`（B1、B2）。`panic` 是另一條路。

更關鍵的是：在 Go 裡，**沒被 recover 的 panic 會炸掉整個行程**，不是只死一條執行緒。  
`recover` 只能寫在 **`defer` 裡面**，常見於：goroutine 邊界、請求／連線 handler 外層，避免單點邏輯拖垮整台 Server。

讀完你要能講清楚：什麼情況回 error、什麼情況允許 panic、recover 之後還要做什麼（記日誌、關連線、要不要隔離房間）。

## Python 對照

| Python | Go |
|--------|-----|
| 未捕捉例外通常死掉**該執行緒** | 未 recover 的 panic 預設死**整個行程** |
| 框架頂層 `try/except Exception` | 邊界上 `defer` + `recover` |
| `raise` 當業務控制流很常見 | **不要**用 panic 模擬業務 throw |
| `finally` 清理 | `defer`（跟 recover 常一起出現） |

這點對寫過 Python 的人特別重要：你以為「頂多死一個 worker thread」，在 Go 可能是**整台遊戲 Server 進程沒了**。

## 怎麼寫

```go
func safeHandle(h func()) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("panic: %v", rec)
			// 實務還要：堆疊、metrics、關閉有問題的連線
		}
	}()
	h()
}
```

把「可能被壞輸入／第三方外掛搞炸」的程式，關在有保護的邊界裡。  
但請記住：recover **不是**把業務錯誤改写成 panic 的許可證。

## 細節

### 什麼時候可以用 panic？

比較說得通的例子：

- 程式員保證不會發生的事發生了（內部不變量）  
- 啟動期設定錯到無法運行（有人偏好直接 FATAL，也行）  
- 真的不知如何繼續，且「繼續跑」比重啟更危險  

什麼時候**不該** panic：

- 客戶端亂送封包  
- 房間滿了、狀態不對、權限不足  
- 網路暫時失敗  

對網路輸入：先驗證，再回 `error`。把壞人輸入變成 panic，等於把 DoS 變成「一包資料殺服」。

### recover 之後要幹嘛？

只 `recover` 不做事＝把火災煙霧警報拆掉。至少想這幾件事：

1. **記日誌**（含 `debug.Stack()`）  
2. **metrics** 打一筆 panic 計數  
3. **隔離**：關掉該連線；若房間狀態可能已壞，考慮踢房或重啟該 room worker  
4. 能轉成 error 就轉，讓上層走正常關閉流程  

### 為什麼 library 少 panic？

函式庫不知道呼叫端有沒有 defer recover，也不知道行程能不能死。  
公開函式優先 `return error`；把「炸不炸」留給 application 邊界決定。

### 進階可先略過

- `Goexit` 與 panic 不同（別混用）。  
- 測試裡用 `httptest`／自寫 helper 斷言「必須 panic」——少見，多半是測不變量。

## 遊戲 Server 會用在哪

每個 WebSocket 命令 handler 外層可以套一層 protect：一個壞訊息不該殺服。  

但若 panic 發生在**改房間狀態的半途**，記憶體裡的房間可能已經髒了。這時「只 log 然後繼續 tick」可能更糟——寧可隔離該房，也別默默帶傷開打。

## 請丟掉的舊習慣

1. 用 panic 模擬 `throw` 業務例外。  
2. recover 後吞掉、不記堆疊、不當回事。  
3. 在深層 library 對「輸入不合法」直接 panic。

## 動手練習

### 必做

1. 寫一個會 panic 的小函式，用 `defer` + `recover` 接住，並改回傳 `error`。  
2. 用自己的話回答：`room.Apply` 碰到「非法但不該發生」的狀態時，為何通常選 `error` 而不是 `panic`？  

### 選做

1. 在 recover 裡印 `debug.Stack()`，看堆疊長什麼樣。  

## 常見坑

- **以為 panic 只死 goroutine**：沒 recover 就死行程。  
- **在別的 goroutine panic，這邊的 recover 接不住**：recover 只對**同一個 goroutine** 有效。  
- **用 recover 當通用錯誤處理**：會讓控制流變隱形，debug 更痛。

## 延伸閱讀

- <https://go.dev/blog/defer-panic-and-recover>  
