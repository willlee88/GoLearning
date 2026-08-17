---
lessonId: "B3"
title: "panic 與 recover 邊界"
description: "何時可以 panic；在邊界 recover 避免進程被拖垮。"
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

## 本章你會建立的心智模型

`panic` 是**不可繼續的程式狀態**（不變量破壞），不是日常業務錯誤。`recover` 只能在 **defer** 中使用，常見於：goroutine 邊界、外掛式 handler、避免單一連線邏輯拖垮行程。

## Python 對照

| Python | Go |
|--------|-----|
| 未捕捉例外殺執行緒 | 未 recover 的 panic 殺**整個行程** |
| 框架頂層 try/except | 中介層 `defer recover` |

## L1 能用

```go
func safeHandle(h func()) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("panic: %v", rec)
		}
	}()
	h()
}
```

## L2 機制

- library 應回傳 error，少 panic。  
- 對程式員錯誤（索引越界）panic 可接受；對網路輸入應驗證後回 error。  
- recover 後要記錄堆疊（`debug.Stack`）並決定是否關閉連線。

## 請丟掉的 Python 習慣

1. 用 panic 模擬 throw 業務例外。  
2. 吞掉 recover 不記日誌。  

## 遊戲 Server 連結

每個 WS 命令 handler 外層可套 protect，避免一個壞訊息殺服；但**房間狀態損壞**時應考慮隔離或重啟該 room。

## 練習

### 必做

1. 寫一個會 panic 的函式，用 defer recover 捕獲並回傳 error。  
2. 說明為何在 `room.Apply` 內部對「非法但不該發生」的狀態選 error 而非 panic。  

## 延伸閱讀

- <https://go.dev/blog/defer-panic-and-recover>  
