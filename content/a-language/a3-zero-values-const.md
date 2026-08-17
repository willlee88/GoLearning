---
lessonId: "A3"
title: "變數、零值、常數與 iota"
description: "零值可用哲學、const 與 iota 列舉模式。"
volume: "a"
order: 3
level: "l2"
status: "ready"
path_required: true
tags: ["types", "zero-value"]
example: "examples/a03-zero-iota"
prev: "A2"
next: "A4"
---

## 本章你會建立的心智模型

Go 強調 **零值（zero value）可用**：變數宣告後即使沒賦值，也有型別定義的預設。常數用 `const`；`iota` 讓列舉變得簡潔。你要學會把「安全的預設狀態」設計進型別，而不是處處 `None` 檢查。

## Python 對照

| Python | Go |
|--------|-----|
| 未綁定名稱 → `NameError` | 變數必須宣告；零值已存在 |
| `None` 表示缺席 | 基本型別無 None；指標／slice／map 等可 nil |
| `Enum` / 模組級常數 | `const` + 常 `iota` |
| `FINAL` 少用 | `const` 編譯期常數（限基本／字串等） |

## L1 能用

```go
var n int            // 0
var s string         // ""
var ok bool          // false
var p *int           // nil
name := "room-1"     // 短宣告（函式內）

const MaxPlayers = 4

type Phase int
const (
	PhaseLobby Phase = iota
	PhasePlaying
	PhaseEnded
)
```

範例：`examples/a03-zero-iota`。

## L2 機制

### 零值表（精簡）

| 類型 | 零值 |
|------|------|
| 數值 | 0 |
| bool | false |
| string | `""` |
| pointer / slice / map / chan / func / interface | nil |
| struct | 欄位各自零值 |

### iota

每個 `const` 區塊從 0 遞增；可運算：

```go
const (
	_ = iota
	KB = 1 << (10 * iota)
	MB
)
```

### 短宣告 `:=`

至少要有一個**新變數**；作用域是當前區塊（含 `if`/`for` 簡短語句）。

## L3 深潛（可選）

- 無型別常數（untyped constant）的精度與推論。  
- 為何 `nil` 沒有獨立型別——必須能推斷上下文型別。

## 請丟掉的 Python 習慣

1. 用「未初始化」表示錯誤狀態卻不定義合法預設。  
2. 字串魔法狀態 `"lobby"` 散落——改 `iota` 或字串常數集中。  
3. 依賴 truthy：`if n` 非法，要 `if n != 0`。

## 遊戲 Server 連結

房間相位 `Lobby/Playing/Ended`、錯誤碼、訊息 type 枚舉，都適合 `iota`。零值 phase=0 應對應**安全預設**（例如 Lobby），避免零值變成「已結束」。

## 練習

### 必做

1. 跑 `examples/a03-zero-iota`。  
2. 為 `Phase` 寫 `String() string` 方法。  

### 選做

1. 設計 `DisconnectReason` 列舉並文件化每個零值是否可出現在 wire 上。  

## 常見坑與如何看見

- **iota 中插入新值導致序號漂移**：協定中的整數枚舉變更要版本化。  
- **struct 零值是否「合法房間」**：在建構函式文件化。  

## 延伸閱讀

- <https://go.dev/ref/spec#The_zero_value>  
