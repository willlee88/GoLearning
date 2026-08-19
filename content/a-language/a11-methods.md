---
lessonId: "A11"
title: "method：值接收者還是指標接收者？"
description: "方法就是帶 receiver 的函式。要改狀態多用指標；方法集會影響能不能塞進介面。"
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

## 這章你會搞懂什麼

Method（方法）就是「左邊多一個 **receiver（接收者）** 的函式」。  
**值接收者**操作的是拷貝；**指標接收者**才能共享並改到同一份狀態。

方法還會組成**方法集（method set）**——這決定型別算不算實作某介面。下一章 interface 會接著打這橋。

## Python 對照

| Python | Go |
|--------|-----|
| `def f(self)` | `func (t T) F()` 或 `(t *T) F()` |
| 實例方法幾乎總是參考語意 | 你可以選值或指標語意 |
| `@classmethod`／`@staticmethod` | 多半寫成 package 級函式 |

## 怎麼寫

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

## 細節

### 怎麼選 receiver

| Receiver | 方法集大致落在 | 典型用途 |
|----------|----------------|----------|
| `(T)` | `T` 與 `*T` 都能用這些方法 | 只讀、小值 |
| `(*T)` | 主要算在 `*T` | 改狀態、大結構 |

實用口訣：

- 要改欄位 → 指標  
- 同一型別的方法風格盡量統一（實務上房間／玩家多半整組用指標）  
- 若已經有任一指標方法，對介面賦值時通常準備好用 `*T`

### 為什麼 `var c Counter; c.Inc()` 常常能編譯？

因為 `c` 是**可取址（addressable）**變數，編譯器可以幫你變成 `(&c).Inc()`。  
但不是所有東西都可取址（例如 map 裡的 struct 值）——那就会編譯失敗，解法常是改存 `*T`。

### nil receiver

指標方法被用 `nil` 叫到時可以防禦（像上面 `Inc`），但要文件化：這算不算合法用法？

### 進階可先略過

- 暫存值不可取址，無法直接叫指標方法。  
- 介面裡裝的是 `T` 還是 `*T`，動態類型不同，下一章的 nil 陷阱會相關。

## 遊戲 Server 會用在哪

```go
func (r *Room) Apply(cmd Command) error
func (r Room) Snapshot() State // 只讀快照：值語意或回傳拷貝都可以想清楚
```

規則寫在方法上很直覺；但記得——**方法不自動等於執行緒安全**，鎖還是要你自己設計。

## 請丟掉的舊習慣

1. 滿脑子 `self` 可變，結果用值接收者，改了等於改副本。  
2. 在 method 裡偷摸一堆全域狀態，測試瞬間變地獄。  
3. 同一型別值／指標接收者混到自己也說不清。

## 動手練習

### 必做

1. 跑 `examples/a11-methods`。  
2. 用自己的話說明：為什麼可取址的 `Counter` 能叫 `Inc`。  

### 選做

1. 對 map 裡的 struct 值叫指標方法——看編譯錯誤，再改成存 `*T`。  

## 常見坑

- **混用值／指標 → 介面「好像實作了又沒有」**。  
- **map 元素不能取址**：這不是 bug，是規則；存指標。  
- **nil receiver 沒防禦又沒文件**：偶發 panic，超難追。
