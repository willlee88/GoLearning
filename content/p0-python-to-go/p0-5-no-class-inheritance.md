---
lessonId: "P0.5"
title: "沒有 class 繼承的世界"
description: "用 struct、method 與 interface 組合行為。"
volume: "p0"
order: 5
level: "l2"
status: "ready"
path_required: true
tags: ["python-bridge", "interfaces"]
example: ""
prev: "P0.4"
next: "P0.6"
---

## 本章你會建立的心智模型

Go 沒有 class / 繼承樹。你用 **struct 保存資料**、**method 綁行為**、**interface 描述能力**。型別只要擁有介面要求的方法集，就自動滿足介面（結構化型別），不必宣告 `implements`。

## Python 對照

| Python | Go |
|--------|-----|
| `class` + 繼承 | struct + 嵌入（embedding） |
| duck typing（執行期） | 方法集滿足 interface（編譯期） |
| ABC / Protocol | interface |
| `self` | receiver（值或指標） |

## L1 能用

```go
type Player struct {
	Name string
	HP   int
}

func (p *Player) Damage(n int) {
	p.HP -= n
}

type Health interface {
	Damage(n int)
}
```

任何有 `Damage(int)` 方法的型別都能當成 `Health`。

## L2 機制

- **接受 interface，回傳具體型別** 是常見 API 風格。
- 小介面（1–3 個方法）比巨大「上帝介面」好組合。
- 嵌入 struct 可委派方法，但**不是**子型別多型的完整替代——要有意識使用。

## L3 深潛（可選）

- 介面值的動態類型與動態值；（`nil` 指標放進 interface 不是 `== nil`）陷阱預告。
- 方法集：值接收者 vs 指標接收者如何影響是否實作介面。

## 請丟掉的 Python 習慣

1. 一上來設計三層繼承的 `BaseEntity → MovingEntity → Player`。
2. 為了「以後擴充」做巨大 base class。
3. 把 interface 當成「註解」，卻不在 API 邊界使用。

## 遊戲 Server 連結

`Transport`、`RoomStore`、`Clock`（可測時間）都適合作小介面，讓規則邏輯可單測、可替換。

## 練習

### 必做

1. 定義 `type Notifier interface { Notify(msg string) }`，實作一個 `ConsoleNotifier`。
2. 寫函式 `func Broadcast(n Notifier, msg string)`，只依賴介面。

### 選做

1. 用嵌入做 `type Admin struct { Player; Level int }`，觀察方法提升。

## 常見坑與如何看見

- 為每個 struct 先寫 interface：過度抽象。等有第二個實作或測試需要再抽。

## 延伸閱讀

- <https://go.dev/doc/effective_go#interfaces>
