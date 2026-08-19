---
lessonId: "A13"
title: "型別斷言與 type switch：從介面拿回具體型別"
description: "手裡只有介面時怎麼分支；comma-ok 才安全。遊戲訊息分派很常用。"
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

## 這章你會搞懂什麼

當你手上只有 interface（或 `any`），卻需要依「實際是什麼型別」分支時，用**型別斷言（type assertion）**或 **type switch**。

遊戲協定常先解成信封（envelope），再分派到具體命令——這是 Server 訊息路徑的基本功。

## Python 對照

| Python | Go |
|--------|-----|
| `isinstance(x, T)` | `x.(T)` 或 type switch |
| `match` 依型別 | `switch x := v.(type)` |

## 怎麼寫

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

## 細節

### 一定要用 comma-ok

```go
v := x.(T)      // 失敗 → panic
v, ok := x.(T)  // 失敗 → ok == false，安全
```

初學請當預設規則：**能寫 comma-ok 就寫**。單值斷言的 panic 常常難追。

### type switch 適合「封閉的一組型別」

每次分支會把 `v` 綁成該案的具體型別，讀起來清楚。  
但若你發現自己寫出巨型 type switch，也可能是氣味：也許該讓各型別實作同一個方法（例如 `Apply`），用多型取代手動分支。

### Wire 上更常見的兩段式

網路上你常先看到字串／數字的 `type`，再解 payload：

```text
JSON { "type": "move", "payload": {...} }
  → 解 type
  → 解 payload 到 MoveCommand
  → room.Apply(cmd)
```

內部已經是 Go 型別時，再用 type switch 也很合理。G 卷會再談信封設計。

### 進階可先略過

- 跟 visitor 模式的取捨。  
- 反射派發（通常太重，業務別先走這條）。

## 遊戲 Server 會用在哪

進房、移動、Ready、離開……每種指令一個 struct，集中分派到 `Room.Apply`。  
未知 type 要有穩定錯誤（別默默吞），否則客戶端協議一歪你全場默劇。

## 請丟掉的舊習慣

1. 巨型 `if isinstance` 樹又不文件化擴充點。  
2. 斷言失敗靠例外碰運氣——Go 用 `ok` 或先檢查。  
3. default 什麼都不做：至少記 log／回錯誤碼。

## 動手練習

### 必做

1. 跑 `examples/a13-type-switch`。  
2. 幫命令聯集加一種 `LeaveCommand`。  

### 選做

1. 比較「介面方法 Apply」vs「type switch」在規則模組的可測性。  

## 常見坑

- **單值斷言 panic**：強制 comma-ok。  
- **漏案卻沒 default**：新指令上線後行為像「沒發生」。  
- **在熱路徑用反射取代單純 type switch**：先求清楚與可測，再談炫技。
