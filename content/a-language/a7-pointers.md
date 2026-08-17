---
lessonId: "A7"
title: "指標"
description: "值語意、指標語意、nil 與方法接收者預告。"
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

## 本章你會建立的心智模型

Go 有指標，但**沒有指標運算**（不能 `p++` 遊走記憶體）。預設是值語意；需要共享或修改同一份資料時用指標。`nil` 指標解引用會 panic——要像對待 Python `None` 一樣清醒，但更常搭配「零值 struct 就夠用」的設計。

## Python 對照

| Python | Go |
|--------|-----|
| 物件預設共享（參考語意） | 結構體賦值是拷貝 |
| `None` | `nil`（有型別上下文） |
| 無顯式 `*` | `*T` 型別、`&v` 取址、`*p` 解參 |

## L1 能用

```go
x := 10
p := &x
*p = 20 // x == 20

var q *int
// *q // panic
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

## L2 機制

### 何時用指標

- 要修改呼叫端看得到的 struct  
- 避免大結構拷貝  
- 表達「可選」（nil）——但可選也可以用 `(T, bool)`  

### 何時用值

- 小資料、不可變語意、map 的 key（不能是含指標的某些型別組合要小心）  
- 清楚表達拷貝  

### 方法接收者預告

`func (p *Player) Damage(n int)` vs `(p Player)` —— 指標接收者可改狀態；方法集規則影響 interface（A12–A13）。

## L3 深潛（可選）

- 逃逸分析：區域變數取址可能讓物件上 heap（`go build -gcflags=-m`）。  
- 指標與 GC：不必 free，但仍避免無用大型物件圖。

## 請丟掉的 Python 習慣

1. 以為 `p2 := p1`（struct）共享同一物件。  
2. 到處傳指標「以後可能要改」——導致別名與生命週期混亂。  
3. 不檢查就解 `nil`。

## 遊戲 Server 連結

房間內玩家實體通常以 `map[PlayerID]*Player` 或 `map[PlayerID]Player` 管理；選哪種取決於是否共享修改、序列化成本。連線 session 幾乎一定是指標／參考型生命週期物件。

## 練習

### 必做

1. 跑 `examples/a07-pointers`。  
2. 證明：值傳遞的 `heal(Player)` 無法改呼叫端 HP；指標可以。  

### 選做

1. 觀察 `go build -gcflags=-m` 對 `return &Player{}` 的逃逸提示。  

## 常見坑與如何看見

- **map 取元素無法取址**（舊規則／細節）：`p := &m[k]` 非法；存 `*Player`。  
- 回傳區域變數指標通常 OK（逃逸到 heap）。  

## 延伸閱讀

- <https://go.dev/doc/effective_go#pointers_vs_values>  
