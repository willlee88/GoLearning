---
lessonId: "A5"
title: "函式與多回傳值"
description: "參數、命名回傳、多回傳與錯誤慣例起點。"
volume: "a"
order: 5
level: "l2"
status: "ready"
path_required: true
tags: ["functions", "errors"]
example: "examples/a05-functions"
prev: "A4"
next: "A6"
---

## 本章你會建立的心智模型

函式是 Go 的基本抽象單位。多回傳值讓 **`(result, error)`** 成為文化，而不是靠全域例外。先把簽名寫清楚：輸入是什麼、失敗長什麼樣、零值結果在錯誤時是否仍有意義。

## Python 對照

| Python | Go |
|--------|-----|
| 多回傳用 tuple | 一等公民多回傳 |
| 關鍵字參數 | 無；用 struct 或 options 模式 |
| 預設參數 | 無；零值或 Option 函式 |
| `*args/**kwargs` | 可變參數 `...T` |

## L1 能用

```go
func add(a, b int) int {
	return a + b
}

func div(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}
```

```go
v, err := div(10, 2)
if err != nil {
	// handle
}
```

範例：`examples/a05-functions`。

## L2 機制

### 傳值

參數是**值拷貝**（slice/map 頭是小描述符，另講）。大 struct 常傳指標。

### 命名回傳值

```go
func parse(s string) (n int, err error) {
	// 可裸 return；初學建議仍明確
	return
}
```

錯誤路徑較複雜時命名回傳可讀，但易與陰影變數打架——**克制使用**。

### 一等函式

函式可當值傳：中介層、策略、回呼。遊戲裡可用於「訊息 handler 表」。

## L3 深潛（可選）

- 閉包捕獲迴圈變數（Go 1.22 前後行為差異值得查 release notes）。  
- `defer` 與命名回傳的互動（A6）。

## 請丟掉的 Python 習慣

1. 用例外當「找不到玩家」的正常路徑——改 `(Player, bool)` 或 `error`。  
2. 一堆預設參數的巨函式——拆函式或 config struct。  
3. 忽略回傳的 error：`result, _ := f()` 必須是有意識決策。

## 遊戲 Server 連結

```go
func (r *Room) Join(p Player) error
func (r *Room) Apply(cmd Command) (Events, error)
```

邊界函式簽名即協定：呼叫端永遠先看 `error`。

## 練習

### 必做

1. 跑 `examples/a05-functions`。  
2. 寫 `func FindPlayer(ps []Player, id int64) (Player, bool)`。  

### 選做

1. 用 `...Command` 實作批次套用並在第一個錯誤停止。  

## 常見坑與如何看見

- 錯誤時仍回傳「半成品」卻未文件化——約定：**err != nil 時結果不可用**，除非 API 明確說。  

## 延伸閱讀

- <https://go.dev/doc/effective_go#functions>  
