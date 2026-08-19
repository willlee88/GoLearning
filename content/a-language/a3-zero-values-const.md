---
lessonId: "A3"
title: "變數、零值、常數與 iota"
description: "宣告了就有預設值；用 const／iota 把狀態與列舉寫清楚，少靠字串魔法。"
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

## 這章你會搞懂什麼

Go 很強調 **零值（zero value）可用**：變數宣告後就算你還沒賦值，也已經有型別規定的預設。  
常數用 `const`；`iota` 則是寫列舉時的省事工具。

你要練的不是背表格，而是：**把「安全的預設狀態」設計進型別**，而不是到處檢查 `None`。

## Python 對照

| Python | Go |
|--------|-----|
| 名稱沒綁定 → `NameError` | 變數必須宣告；零值已經在那 |
| `None` 表示缺席 | 基本型別沒有 None；指標／slice／map 等才可能 `nil` |
| `Enum`／模組級常數 | `const`，常搭配 `iota` |
| 很少用真正的常數 | `const` 是編譯期常數（數字、字串這類） |

## 怎麼寫

```go
var n int            // 0
var s string         // ""
var ok bool          // false
var p *int           // nil
name := "room-1"     // 短宣告（只能在函式內）

const MaxPlayers = 4

type Phase int
const (
	PhaseLobby Phase = iota // 0
	PhasePlaying            // 1
	PhaseEnded              // 2
)
```

範例：`examples/a03-zero-iota`。

## 細節

### 常見零值（精簡表）

| 類型 | 零值 |
|------|------|
| 數值 | 0 |
| bool | false |
| string | `""`（空字串，不是 nil） |
| pointer / slice / map / chan / func / interface | nil |
| struct | 每個欄位各自的零值 |

為什麼這很重要？因為 Go 鼓勵你把「剛宣告出來的東西」設計成**已經能安全用**。例如計數器從 0 開始、空字串當「尚未命名」——但前提是你文件化這些語意。

### `iota` 在幹嘛

每個 `const (` 區塊從 0 往上加。也可以拿來算：

```go
const (
	_  = iota // 丟掉 0
	KB = 1 << (10 * iota)
	MB
)
```

列舉很適合房間相位、錯誤碼分類、訊息 type。但若這些數字會上線到協定（wire），**插入新值導致序號漂移**就很痛——要版本化或改用穩定字串／明確指定數值。

### 短宣告 `:=`

至少要有一個**新變數**才行。作用域是當前區塊，包含 `if`／`for` 前面的短語句——這後面處理 `err` 時超常見。

### 進階可先略過

- 無型別常數（untyped constant）的精度與型別推論。  
- `nil` 沒有獨立型別，必須能從上下文推斷。

## 遊戲 Server 會用在哪

房間相位 `Lobby / Playing / Ended`、斷線原因、訊息 type，都很適合 `iota` 或集中常數。

特別注意：**零值 phase = 0 應該對應安全預設**（通常是 Lobby），千萬別讓 0 意外代表「已結束」——否則 `var Room{}` 一出來就像鬼打牆。

## 請丟掉的舊習慣

1. 用「未初始化」當錯誤狀態，卻不定義什麼叫合法預設。  
2. 字串魔法 `"lobby"` 散落全專案——改集中常數或 `iota`。  
3. 依賴 truthy：`if n` 在 Go 不合法，要寫 `if n != 0`。

## 動手練習

### 必做

1. 跑 `examples/a03-zero-iota`。  
2. 幫 `Phase` 寫 `String() string`，印人話而不是裸數字。  

### 選做

1. 設計 `DisconnectReason` 列舉，並註明「零值會不會出現在網路上」。  

## 常見坑

- **iota 中間插入新成員**：舊客戶端／舊存檔可能解錯——協定整數枚舉要謹慎。  
- **以為空字串是 nil**：`string` 的零值是 `""`，跟指標／slice 的 nil 不同。  
- **struct 零值算不算「合法房間」**：請在建構函式或註解寫清楚。
