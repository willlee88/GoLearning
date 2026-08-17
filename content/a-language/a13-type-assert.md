---
lessonId: "A13"
title: "型別斷言與 type switch"
description: "從 interface 取出具體型別；訊息分派模式。"
volume: "a"
order: 13
level: "l2"
status: "ready"
path_required: true
tags: ["interfaces", "dispatch"]
example: "examples/a13-type-switch"
prev: "A12"
next: "A14"
---

## 本章你會建立的心智模型

當你手裡只有 interface（或 `any`），需要依實際型別分支時，用**型別斷言**或 **type switch**。遊戲協定常先解成「信封」，再分派到具體命令——這是 Server 訊息路徑的基本功。

## Python 對照

| Python | Go |
|--------|-----|
| `isinstance(x, T)` | `x.(T)` / `x.(type)` |
| `match` 型別 | `switch x := v.(type)` |

## L1 能用

```go
var n Notifier = LogNotifier{}
if ln, ok := n.(LogNotifier); ok {
	_ = ln
}

switch v := msg.(type) {
case JoinCommand:
	// ...
case MoveCommand:
	// ...
default:
	// unknown
}
```

範例：`examples/a13-type-switch`。

## L2 機制

- `v.(T)` 失敗會 **panic**；用 `v.(T)` 的 comma-ok 形式。  
- type switch 每次分支綁定具體型別。  
- 過度依賴 type switch 可能是「該用介面方法卻沒用」的氣味。

## L3 深潛（可選）

- 與 visitor 模式取捨。  
- JSON 先解 `Type` 字串再解 payload 的兩段式（wire 更常見）。

## 請丟掉的 Python 習慣

1. 巨型 `if isinstance` 樹卻不文件化擴充點。  
2. 斷言失敗靠例外——Go 用 ok 或預先檢查。  

## 遊戲 Server 連結

```text
JSON { "type": "move", "payload": {...} }
  → 解 type
  → 解 payload 到 MoveCommand
  → room.Apply(cmd)
```

或內部已是 Go 型別時用 type switch。

## 練習

### 必做

1. 跑 `examples/a13-type-switch`。  
2. 為 `Command` 聯集加一種 `LeaveCommand`。  

### 選做

1. 比較「介面方法 Apply」vs「type switch」在規則模組的可測性。  

## 常見坑與如何看見

- 單值斷言 panic 難以定位——強制 comma-ok。  

## 延伸閱讀

- <https://go.dev/ref/spec#Type_assertions>  
