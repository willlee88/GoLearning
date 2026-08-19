---
lessonId: "A4"
title: "控制流：if、for、switch（就這麼精簡）"
description: "只有一種 for、條件一定是 bool、短宣告有作用域；房間狀態機很適合用 switch 寫。"
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

## 這章你會搞懂什麼

Go 的控制流刻意很少：`if`、一種 `for`、很能打的 `switch`。條件必須是真正的 `bool`，沒有 Python 那種「非空就算真」。

另外，短宣告的**作用域（scope）**很常影響 `err` 與分支變數——看起來像小語法，寫連線／錯誤處理時卻天天踩。

## Python 對照

| Python | Go |
|--------|-----|
| `if/elif/else` | `if` / `else if` / `else` |
| `for` 與 `while` | **全部**用 `for` |
| `match`（3.10+） | `switch`（還可以沒有運算式） |
| 三元運算 `a if c else b` | **沒有**，老老實實寫 `if` |

## 怎麼寫

```go
if err != nil {
	return err
}

if p := lookup(id); p != nil {
	// p 只活在這個 if 裡
}

for i := 0; i < 3; i++ {
}
for _, p := range players {
}
for { // 無限迴圈
	break
}

switch phase {
case PhaseLobby:
	// ...
case PhasePlaying, PhaseEnded:
	// ...
default:
	// ...
}
```

## 細節

### 為什麼只有一種 `for`

Go 不想讓你記三套迴圈語法。型態靠寫法區分：

```go
for cond { }              // 像 while
for i := 0; i < n; i++ { }
for i, v := range s { }   // 走訪 slice／map／字串／channel
```

`range` 碰到不同型別行為不一樣：字串走的是 rune（後面 A9），map 順序還是刻意打亂的（A8）。

### `switch` 的兩個重點

1. **預設不 fallthrough**（跟 C 不一樣）。真要穿透就寫顯式 `fallthrough`。  
2. 無運算式的 `switch { case x > 0: }` 很適合取代長長的 if-else。

### 短宣告作用域

```go
if err := do(); err != nil {
	return err
}
// 這裡外層若沒宣告過 err，這個 err 已經不存在了
```

好處是錯誤處理變數不會污染外層；壞處是你可能「以為還有 err」或跟外層同名變數**陰影（shadow）**——`go vet` 與 code review 要盯。

### 進階可先略過

- `goto` 存在，業務碼幾乎不用，多見於產生碼。  
- `break` 可以帶 label，一次跳出外層迴圈。

## 遊戲 Server 會用在哪

房間狀態機超適合 `switch`：

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

先把「這個相位允許哪些動作」寫清楚，比散落一堆布林旗標好維護太多。後面 H 卷的 FSM 也是同一精神。

## 請丟掉的舊習慣

1. 滿手 `while True`——在 Go 寫 `for { }`。  
2. 依賴隱式 truthy（空字串、0、空 list）。  
3. 以為 `switch` 會像 C 一樣自動往下落。

## 動手練習

### 必做

1. 用 `switch` 實作 `func canStart(phase Phase, ready int, need int) bool`。  
2. 用 `for range` 算切片裡 `score > 10` 的人數。  

### 選做

1. 先畫房間狀態轉移圖，再寫成程式裡的合法轉移表。  

## 常見坑

- **`if v := f(); …` 跟外層 `v` 陰影**：讀碼時很常看錯。  
- **一邊 range 一邊改正在迭代的 slice**：行為容易超乎直覺，必要時先拷貝。  
- **條件寫成非 bool**：直接編譯失敗——這是好事，別想繞。
