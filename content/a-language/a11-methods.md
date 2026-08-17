---
lessonId: "A11"
title: "method 與接收者"
description: "值／指標接收者、方法集，以及何時該用指針。"
volume: "a"
order: 11
level: "l2"
status: "ready"
path_required: true
tags: ["methods"]
example: "examples/a11-methods"
prev: "A10"
next: "A12"
---

## 本章你會建立的心智模型

Method 是「帶 receiver 的函式」。**值接收者**操作拷貝；**指標接收者**共享並可改狀態。方法集（method set）決定型別是否實作某 interface——這是後一章的橋。

## Python 對照

| Python | Go |
|--------|-----|
| `def f(self)` | `func (t T) F()` / `(t *T) F()` |
| 實例方法幾乎總是參考 | 可選值或指標語意 |
| `@classmethod` / `@staticmethod` | 包級函式或未匯出 helper |

## L1 能用

```go
type Counter struct{ n int }

func (c Counter) Value() int { return c.n }

func (c *Counter) Inc() {
	if c == nil {
		return
	}
	c.n++
}
```

範例：`examples/a11-methods`。

## L2 機制

| Receiver | 方法集含於 | 典型用途 |
|----------|------------|----------|
| `(T)` | `T` 與 `*T` | 只讀、小值 |
| `(*T)` | 僅 `*T` | 變更狀態、大結構 |

規則（簡化）：若有任一指標接收者方法，實作 interface 時通常用 `*T`。

一致性：同一型別的方法 receiver 風格盡量統一（多為指標）。

## L3 深潛（可選）

- 非 addressable 的暫時值無法取址呼叫指標方法。  
- 介面存放 `T` 與 `*T` 的動態類型差異。

## 請丟掉的 Python 習慣

1. 所有方法都假設 self 可變卻傳值拷貝。  
2. 在 method 裡隱藏大量全域狀態。  

## 遊戲 Server 連結

```go
func (r *Room) Apply(cmd Command) error
func (r Room) Snapshot() State // 只讀快照可值或回傳拷貝
```

## 練習

### 必做

1. 跑 `examples/a11-methods`。  
2. 說明為何 `var c Counter; c.Inc()` 能編譯（addressable）。  

### 選做

1. 嘗試對 map 中的 struct 值呼叫指標方法——觀察編譯錯誤並改存 `*T`。  

## 常見坑與如何看見

- 混用值／指標導致 interface 不滿足。  
- nil receiver：方法可防禦，但要文件化。  

## 延伸閱讀

- <https://go.dev/doc/effective_go#methods>  
