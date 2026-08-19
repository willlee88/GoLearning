---
lessonId: "A5"
title: "函式與多回傳：（結果, error）文化"
description: "參數怎麼傳、多回傳怎麼寫；為什麼失敗用回傳值，而不是丟例外。"
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

## 這章你會搞懂什麼

函式是 Go 最基本的抽象單位。多回傳值讓 **`(result, error)`** 變成日常文化，而不是靠全域例外往上冒。

寫函式前先問清楚：輸入是什麼、失敗長什麼樣、錯誤時那個結果還能不能用。簽名寫清楚，呼叫端才不會猜。

## Python 對照

| Python | Go |
|--------|-----|
| 多回傳靠 tuple | 多回傳是一等公民 |
| 關鍵字參數 | **沒有**；常用 struct 或 options 模式 |
| 預設參數 | **沒有**；靠零值或 Option 函式 |
| `*args` / `**kwargs` | 可變參數 `...T`（有型別） |

## 怎麼寫

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

呼叫端幾乎總是這個節奏：

```go
v, err := div(10, 2)
if err != nil {
	// 處理失敗：記錄、回傳、轉成協定錯誤碼…
}
```

範例：`examples/a05-functions`。

## 細節

### 參數是值拷貝

傳進函式的是拷貝。對小整數沒感覺；對大 struct 常改傳指標（下一章）。  
slice／map 傳的是「小描述符／引用頭」，底層資料可能共享——細節在 A8。

### 命名回傳值

```go
func parse(s string) (n int, err error) {
	// 可以裸 return；初學建議還是寫明白
	return
}
```

錯誤路徑很複雜時，命名回傳有時比較好讀；但也容易跟陰影變數打架——**克制使用**，不要為用而用。

### 函式當值

函式可以當參數傳：中介層、策略、回呼。遊戲裡常見「訊息 type → handler」表。

### 為什麼用回傳 error，不用例外當主路徑

因為呼叫端被逼著看見失敗。你當然可以忽略（`result, _ := f()`），但那是**有意識的決定**，不是預設吞掉。B 卷會把 error 講得更深。

### 進階可先略過

- 閉包捕獲迴圈變數（Go 1.22 前後行為有差異，查 release notes）。  
- `defer` 與命名回傳的互動（下一章 A6）。

## 遊戲 Server 會用在哪

邊界 API 幾乎長這樣：

```go
func (r *Room) Join(p Player) error
func (r *Room) Apply(cmd Command) (Events, error)
```

簽名就是契約：呼叫端**永遠先看 `error`**。滿房、相位不對、非法指令，都走這條路，而不是丟例外碰運氣。

## 請丟掉的舊習慣

1. 用例外表達「找不到玩家」這種正常分支——改 `(Player, bool)` 或回 `error`。  
2. 一堆預設參數的巨函式——拆函式，或收成 config struct。  
3. 下意識 `_, _ =` 丢掉 error：除非你真的知道為什麼可以忽略。

## 動手練習

### 必做

1. 跑 `examples/a05-functions`。  
2. 寫 `func FindPlayer(ps []Player, id int64) (Player, bool)`。  

### 選做

1. 用 `...Command` 做批次套用，第一個錯誤就停。  

## 常見坑

- **`err != nil` 還把半成品結果當可用**：團隊要有約定——通常 **有錯就當結果不可用**，除非 API 文件明確說。  
- **回傳值順序混亂**：Go 慣例常是 `(value, error)` 或 `(value, bool)`，別自創難記順序。  
- **可變參數傳 slice**：要用 `sum(nums...)` 展開，不是直接塞一個 slice 變數當「多個參數」卻忘了 `...`。
