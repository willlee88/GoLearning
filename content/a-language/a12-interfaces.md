---
lessonId: "A12"
title: "interface：描述「會做什麼」，不是繼承"
description: "方法集對了就算滿足介面；越小越好用。另外：nil 指標塞進介面後，介面不一定是 nil。"
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

## 這章你會搞懂什麼

Interface（介面）描述的是**能力**——一組方法——不是繼承樹。  
某個型別只要擁有需要的方法，就算滿足介面（隱式，不用寫 `implements`）。

介面值內部可以想成：「動態類型 + 動態值」。經典深坑是：把一個 **nil 指標**裝進介面後，介面本身**不一定** `== nil`。這章一定要搞懂它。

## Python 對照

| Python | Go |
|--------|-----|
| duck typing（多半執行期才爆） | 方法集在編譯期檢查「賦值／傳參」 |
| Protocol／ABC | `interface` |
| `None` | 要分「nil 介面」vs「介面裡裝著 typed nil」 |

## 怎麼寫

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

範例：`examples/a12-interfaces`（含 nil 示範輸出）。

## 細節

### 小介面比較香

Go 文化很喜歡小介面：一個方法、兩個方法，好組合、好喝假。精神上接近 `io.Reader`——你要的是「能讀」，不是「整個宇宙的檔案抽象」。

常見風格：**函式參數收介面，回傳具體型別**。這樣呼叫端好用，實作者也好換。

### `any` 是什麼

`any` 就是 `interface{}`：等於「我放棄靜態資訊」。除錯、框架邊界偶爾用；業務核心狂用 `any`，等於把 Python 動態痛點請回來。

### nil 陷阱（必看）

```go
var p *LogNotifier = nil
var n Notifier = p
fmt.Println(n == nil) // false！因為動態類型是 *LogNotifier
```

為什麼？介面值要「類型與值都空」才是 nil。你把 typed nil 塞進去，類型欄位還在，所以 `n != nil`，接著呼叫方法可能又踩到 nil receiver。

回傳 `(Notifier, error)` 時，錯誤路徑請回 **`nil, err`**（真的 nil 介面），不要回「某個 nil 指標當介面」。

### 進階可先略過

- 介面轉換用的 itable。  
- 介面呼叫相對具體型別的一點點開銷——先正確，熱路徑再測。

## 遊戲 Server 會用在哪

```go
type Clock interface{ Now() time.Time }
type RoomStore interface {
	Get(id string) (*Room, error)
	Save(*Room) error
}
```

可替換時間（測試好控）、可替換儲存；規則邏輯不必 import 具體資料庫套件。

## 請丟掉的舊習慣

1. 每個 struct 先寫巨大 interface「以後也許用得到」。  
2. 用 `any` 當「動態語言回歸」。  
3. 回傳 typed nil 卻只檢查 `err != nil`／`n != nil`，結果邏輯全歪。

## 動手練習

### 必做

1. 跑 `examples/a12-interfaces`，盯著 nil 示範輸出。  
2. 定義 `Transporter` 介面：`Send([]byte) error`。  

### 選做

1. 寫 `(Notifier, error)`：錯誤時正確回 `nil, err`（不要 typed nil）。  

## 常見坑

- **介面不是 nil，裡面卻是 nil 指標**：用 `%#v`／看動態類型除錯。  
- **為了 mock 做出 20 個方法的介面**：測試痛苦是設計味道，不是框架不夠強。  
- **過早抽象**：先有兩個實作再抽介面，通常比較準。
