---
lessonId: "A4"
title: "控制流"
description: "if/for/switch 與短宣告作用域；狀態機寫法入門。"
volume: "a"
order: 4
level: "l2"
status: "ready"
path_required: true
tags: ["control-flow"]
example: ""
prev: "A3"
next: "A5"
---

## 本章你會建立的心智模型

Go 的控制流刻意精簡：只有一種 `for`、`switch` 很強大、條件必須是 `bool`。短宣告作用域會影響變數生命——這在處理 `err` 與連線分支時特別重要。

## Python 對照

| Python | Go |
|--------|-----|
| `if/elif/else` | `if/else if/else`（`else if` 常寫成 `else if` 或接 `else { if`） |
| `for` / `while` | 全用 `for` |
| `match`（3.10+） | `switch`（可無運算式） |
| 沒有三元 `?:` 常用式 | 無三元運算子，用 `if` |

## L1 能用

```go
if err != nil {
	return err
}

if p := lookup(id); p != nil {
	// p 只在這個 if 作用域
}

for i := 0; i < 3; i++ { }
for _, p := range players { }
for { /* 無限迴圈 */ break }

switch phase {
case PhaseLobby:
	// ...
case PhasePlaying, PhaseEnded:
	// ...
default:
	// ...
}
```

## L2 機制

### 只有 for

```go
for cond { }           // while
for i := 0; i < n; i++ { }
for i, v := range s { }
```

`range` 在 slice/map/string/channel 行為不同——字串是 rune 迭代（A10 再深）。

### switch

- 預設**不 fallthrough**（與 C 不同）；需要時顯式 `fallthrough`。  
- 無運算式的 `switch { case x > 0: }` 很適合取代長 if-else。  

### 短宣告作用域

```go
if err := do(); err != nil {
	return err
}
// 這裡的 err 不存在（若外層沒宣告）
```

## L3 深潛（可選）

- `goto` 存在但極少用於業務；多見於產生碼。  
- `break` 可帶 label 跳出外層。

## 請丟掉的 Python 習慣

1. `while True` 文化——Go 寫 `for { }`。  
2. 依賴隱式 truthy。  
3. 在 switch 假設 fallthrough。

## 遊戲 Server 連結

房間狀態機：

```go
switch r.phase {
case PhaseLobby:
	// 允許 Ready / Leave
case PhasePlaying:
	// 允許 Input，拒絕 Ready
default:
	return ErrInvalidPhase
}
```

先用 switch 表達**允許的轉移**，比散落的布林旗標清楚。

## 練習

### 必做

1. 用 `switch` 實作 `func canStart(phase Phase, ready int, need int) bool`。  
2. 用 `for range` 計算切片中 score > 10 的人數。  

### 選做

1. 畫一張房間狀態轉移圖，再寫成程式中的合法轉移表。  

## 常見坑與如何看見

- **`if v := f(); v != nil`** 與外層 `v` 陰影（shadow）——`go vet` 與人眼 code review。  
- range 時修改正在迭代的 slice 要小心。  

## 延伸閱讀

- <https://go.dev/doc/effective_go#control-structures>  
