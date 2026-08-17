---
lessonId: "A12"
title: "interface 本質"
description: "方法集滿足、小介面、nil interface 陷阱。"
volume: "a"
order: 12
level: "l2"
status: "ready"
path_required: true
tags: ["interfaces"]
example: "examples/a12-interfaces"
prev: "A11"
next: "A13"
---

## 本章你會建立的心智模型

Interface 描述**能力**（方法集），不是繼承樹。型別只要擁有所需方法即滿足介面（隱式）。介面值內部是「動態類型 + 動態值」——`nil` 指標裝進介面後，介面本身**不一定** `== nil`。這是 Go 經典深坑，也是深刻理解的考點。

## Python 對照

| Python | Go |
|--------|-----|
| duck typing（執行期） | 方法集（編譯期檢查賦值） |
| Protocol / ABC | `interface` |
| `None` | 小心 nil 具體值 vs nil 介面 |

## L1 能用

```go
type Notifier interface {
	Notify(msg string) error
}

type LogNotifier struct{}

func (LogNotifier) Notify(msg string) error {
	fmt.Println(msg)
	return nil
}

func Broadcast(n Notifier, msg string) error {
	return n.Notify(msg)
}
```

範例：`examples/a12-interfaces`。

## L2 機制

- **接受 interface，回傳具體型別** 是常見風格。  
- 介面越小越好組合（`io.Reader` 精神）。  
- 空介面 `any`（`interface{}`）等於放棄靜態資訊——邊界才用。  
- 賦值給介面時：動態類型與動態值一起存。

### nil 陷阱

```go
var p *LogNotifier = nil
var n Notifier = p
fmt.Println(n == nil) // false！動態類型是 *LogNotifier
```

## L3 深潛（可選）

- 介面轉換的 itable。  
- 效能：介面呼叫相對具體型別有代價（先正確再測）。

## 請丟掉的 Python 習慣

1. 為每個 struct 先寫巨大 interface。  
2. 用 `any` 當「動態語言回歸」。  
3. 只檢查 `err != nil` 卻回傳「typed nil」。

## 遊戲 Server 連結

```go
type Clock interface{ Now() time.Time }
type RoomStore interface{
	Get(id string) (*Room, error)
	Save(*Room) error
}
```

可測時間、可替換儲存；規則邏輯不 import 具體 DB。

## 練習

### 必做

1. 跑 `examples/a12-interfaces`（含 nil 示範輸出）。  
2. 定義 `Transporter` 介面：`Send([]byte) error`。  

### 選做

1. 寫一個回傳 `(Notifier, error)` 在錯誤時正確回 `nil, err`（不要 typed nil）。  

## 常見坑與如何看見

- 單元測試 mock 過度巨大 interface。  
- 日誌印 `%#v` 看動態類型。  

## 延伸閱讀

- <https://go.dev/blog/laws-of-reflection> （可選）  
- <https://go.dev/doc/effective_go#interfaces>  
