---
lessonId: "A6"
title: "defer"
description: "延遲執行的堆疊語意：解鎖、關連線、量測。"
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

## 本章你會建立的心智模型

`defer` 把呼叫**推遲到函式返回前**，多個 defer 以**堆疊（LIFO）**執行。它是資源清理的主力：`Unlock`、`Close`、計時結束。先建立「取得資源的下一行就 defer 釋放」的肌肉。

## Python 對照

| Python | Go |
|--------|-----|
| `with` / context manager | `defer`（+ 有時手動生命週期） |
| `try/finally` | `defer` |
| 解構時 `__exit__` | 無內建 RAII；靠 defer 紀律 |

## L1 能用

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

範例：`examples/a06-defer`。

## L2 機制

### 參數何時求值

`defer f(x)` 的 **`x` 在 defer 語句執行時就求值**，不是返回時。若要返回時的值，傳指標或閉包：

```go
defer func() { log.Println(result) }()
```

### LIFO

```go
defer fmt.Println(1)
defer fmt.Println(2)
// 輸出 2 然後 1
```

### 與錯誤包裝

常在 defer 裡根據 `err` 做額外日誌（搭配命名回傳，進階模式）。

## L3 深潛（可選）

- 熱路徑上大量 defer 的成本（通常可接受；先測再優化）。  
- `defer` 在 `for` 迴圈內：每次迭代註冊——檔案迴圈要小心。

## 請丟掉的 Python 習慣

1. 只在成功路徑 `close`，忘記失敗分支——defer 一次搞定。  
2. 把 defer 當「任意非同步排程」——它不開新 goroutine。  

## 遊戲 Server 連結

- 請求結束取消：`defer cancel()`  
- 房間鎖：`defer room.mu.Unlock()`  
- 連線：讀 loop 結束 `defer unregister(conn)`  

## 練習

### 必做

1. 跑 `examples/a06-defer`，預測輸出順序再對答案。  
2. 寫函式：計時執行 `work()`，用 defer 印出耗時。  

### 選做

1. 示範 defer 參數立刻求值 vs 閉包的差異。  

## 常見坑與如何看見

- **在迴圈 defer Close 一堆檔案**：改為函式包一層每次開閉。  
- 忘記 Unlock：`-race` 與死鎖會找上你。  

## 延伸閱讀

- <https://go.dev/blog/defer-panic-and-recover>  
