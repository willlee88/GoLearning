---
lessonId: "A7"
title: "指標：什麼時候要共享同一份資料"
description: "預設是值拷貝；需要改同一份或表達「可選」才用指標。沒有指標運算，但有 nil。"
volume: "a"
order: 7
level: "l2"
status: "ready"
path_required: true
tags: ["pointers"]
example: "examples/a07-pointers"
prev: "A6"
next: "A8"
---

## 這章你會搞懂什麼

Go 有指標（pointer），但**不能做指標運算**（沒有 `p++` 在記憶體上晃）。預設心智是**值語意**：賦值、傳參常常是拷貝。

什麼時候用指標？——你要**共享／修改同一份資料**，或明確表達「沒有」（`nil`）。  
對 `nil` 解引用會 panic，心態上類似小心 Python 的 `None`，但 Go 更常靠「零值 struct 就夠用」減少無謂指標。

## Python 對照

| Python | Go |
|--------|-----|
| 物件預設多是參考語意 | struct 賦值／傳參預設拷貝 |
| `None` | `nil`（要有型別上下文） |
| 不用寫 `*` | `*T` 型別、`&v` 取址、`*p` 解參 |

最大的落差：在 Python 你改 `p2.hp`，常連 `p1` 一起變；在 Go，若 `p2 := p1` 且是 struct 值，那是**另一份拷貝**。

## 怎麼寫

```go
x := 10
p := &x
*p = 20 // 現在 x == 20

var q *int
// *q // panic：nil 指標
```

```go
type Player struct{ HP int }

func heal(p *Player, n int) {
	if p == nil {
		return
	}
	p.HP += n
}
```

範例：`examples/a07-pointers`。

## 細節

### 什麼時候用指標

- 要讓呼叫端看得到修改（例如扣血、改相位）  
- 大結構不想一直拷貝  
- 表達可選：`nil` 代表沒有——但也可以用 `(T, bool)`，不一定要指標

### 什麼時候用值

- 小資料、偏不可變語意  
- 想清楚表達「這是拷貝，互不影响」  
- 當 map 的 key 時，型別要可比較；含指標的複合型別要特別小心

### 方法接收者預告

之後你會寫 `func (p *Player) Damage(n int)` 或 `(p Player)`：  
指標接收者才能改到同一個物件的狀態；這也會影響「算不算實作某個介面」（A11–A12）。

### 進階可先略過

- 逃逸分析（escape analysis）：區域變數被取址，可能搬到 heap（A18）。  
- 不用手動 free，但仍要避免無意義的大型物件圖一直被指著。

## 遊戲 Server 會用在哪

房間裡的玩家常是：

- `map[PlayerID]*Player`（共享修改、生命週期清楚）  
- 或 `map[PlayerID]Player`（較偏值拷貝／快照思維）  

連線 session、連線物件幾乎一定是指標或明確的參考型生命週期——因為你要在很多地方改同一份連線狀態。

## 請丟掉的舊習慣

1. 以為 `p2 := p1`（struct）還共享同一物件。  
2. 到處傳指標「以防萬一以後要改」——別名變多，生命週期更難想。  
3. 不檢查就解 `nil`。

## 動手練習

### 必做

1. 跑 `examples/a07-pointers`。  
2. 自己證明：值傳遞的 `heal(Player)` 改不到呼叫端 HP；指標可以。  

### 選做

1. 對 `return &Player{}` 跑 `go build -gcflags=-m`，看逃逸提示。  

## 常見坑

- **想對 map 裡的 struct 值取址**：`&m[k]` 不合法；需要改的話改存 `*Player`。  
- **回傳區域變數的指標**：通常 OK（編譯器會讓它逃到 heap），別用 C 的怕法來怕。  
- **指標／值混用導致「有時改得到有時改不到」**：先畫資料流，再決定接收者型別。
