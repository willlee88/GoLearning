---
lessonId: "A6"
title: "defer：離開函式前一定做的事"
description: "延遲執行、後進先出；關檔、解鎖、取消 context 的標準寫法。"
volume: "a"
order: 6
level: "l2"
status: "ready"
path_required: true
tags: ["defer"]
example: "examples/a06-defer"
prev: "A5"
next: "A7"
---

## 這章你會搞懂什麼

`defer` 把呼叫**推遲到函式快返回之前**再執行。多個 defer 是**堆疊（stack）**，也就是後進先出（LIFO）。

它是資源清理的主力：解鎖（Unlock）、關檔／關連線（Close）、打點結束時間。肌肉記憶就一句話——**拿到資源的下一行，立刻 defer 釋放**。

## Python 對照

| Python | Go |
|--------|-----|
| `with`／context manager | `defer`（有時還是要自己管生命週期） |
| `try` / `finally` | `defer` |
| 物件解構時 `__exit__` | 沒有內建 RAII；靠 defer 紀律 |

## 怎麼寫

```go
f, err := os.Open(path)
if err != nil {
	return err
}
defer f.Close()
```

```go
mu.Lock()
defer mu.Unlock()
```

範例：`examples/a06-defer`。先猜輸出順序，再跑對答案。

## 細節

### 參數什麼時候求值？

`defer f(x)` 裡的 **`x` 在執行到 defer 那一行時就先求值**，不是函式返回那一刻。  
若你要的是「返回當下」的值，用閉包：

```go
defer func() { log.Println(result) }()
```

這是超常見坑：以為 defer 會「記住變數名字的最新值」，結果印到的是舊的。

### 後進先出

```go
defer fmt.Println(1)
defer fmt.Println(2)
// 先印 2，再印 1
```

多層清理時，順序通常剛好是你要的（先關內層、再關外層）。

### 跟錯誤處理一起用

進階一點可以在 defer 裡看命名回傳的 `err` 再補日誌。初學先把「鎖／檔案／cancel 一定清掉」練穩就好。

### 進階可先略過

- 熱路徑上大量 defer 的成本（多數情況可接受；先測再優化）。  
- **在 `for` 裡 defer**：每次迭代都註冊一個——開很多檔卻等函式結束才關，容易把描述子耗光。解法是包一層小函式，每檔開閉一次。

## 遊戲 Server 會用在哪

- 請求結束取消：`defer cancel()`  
- 房間鎖：`defer room.mu.Unlock()`  
- 連線讀 loop 結束：`defer unregister(conn)`  

長連線服務最怕的就是「成功路徑有清、失敗路徑忘記」。defer 就是在消這種分岐漏清。

## 請丟掉的舊習慣

1. 只在成功路徑 `close`，錯誤分支靠人腦記得——改成拿到就 defer。  
2. 把 defer 想成「丟到背景非同步執行」——它**不會**開新 goroutine，只是延後到返回前。

## 動手練習

### 必做

1. 跑 `examples/a06-defer`，預測輸出順序再對答案。  
2. 寫一個函式：量測 `work()` 耗時，用 defer 印出來。  

### 選做

1. 做一組對照：defer 參數立刻求值 vs 閉包讀最新值。  

## 常見坑

- **迴圈裡 defer Close 一堆檔案**：改成每次開檔包進小函式。  
- **Lock 了卻沒 Unlock**：死鎖或 `-race` 會來找你；defer Unlock 是預設解。  
- **defer 裡再 panic／再犯錯**：清理邏輯要保持單純，別在裡面塞業務大劇情。
