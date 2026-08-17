---
lessonId: "A10"
title: "struct、嵌入與 tag"
description: "實體建模、組合式嵌入、json tag 與零值 struct。"
volume: "a"
order: 10
level: "l2"
status: "ready"
path_required: true
tags: ["struct", "json"]
example: "examples/a10-struct-tags"
prev: "A9"
next: "A11"
---

## 本章你會建立的心智模型

`struct` 是 Go 描述「一筆資料長什麼樣」的主要方式。用**嵌入（embedding）**做組合，而不是繼承。`tag` 把中繼資料掛在欄位上（最常見是 `json`），讓編碼與領域模型對齊。

## Python 對照

| Python | Go |
|--------|-----|
| `@dataclass` / 普通 class | `struct` |
| 繼承 | 嵌入 + 介面 |
| `asdict` / pydantic | `encoding/json` + tag |

## L1 能用

```go
type Vec2 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Player struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Pos  Vec2   `json:"pos"`
}
```

```go
type Timed struct{ time.Time } // 嵌入；提升方法
```

範例：`examples/a10-struct-tags`。

## L2 機制

- 欄位大寫才會被 `encoding/json` 匯出。  
- 嵌入提升方法與欄位，但**不是**子型別：不能把外層當成內層型別傳入。  
- 匿名欄位衝突時需顯式選擇。  
- 零值 struct 合法且常用作預設狀態。

## L3 深潛（可選）

- 記憶體對齊（alignment）與欄位排序對 size 的影響。  
- `json:"-"`、`omitempty` 細節。

## 請丟掉的 Python 習慣

1. 深繼承樹 `BaseEntity`。  
2. 用 dict 當長期領域模型（除了真的動態 schema）。  

## 遊戲 Server 連結

`Player`、`RoomConfig`、`Command` 幾乎都是 struct；wire 格式用 tag 對齊。嵌入適合 `struct { mu sync.Mutex; /* 內部 */ }` 的小模式（注意匯出）。

## 練習

### 必做

1. 跑 `examples/a10-struct-tags`。  
2. 為 `Room` 加 `Capacity int \`json:"capacity"\`` 並 round-trip JSON。  

### 選做

1. 比較嵌入 `time.Time` 與欄位 `CreatedAt time.Time` 的 JSON 差異。  

## 常見坑與如何看見

- 小寫欄位靜默消失在 JSON。  
- 嵌入造成意外的方法提升與 JSON 形狀。  

## 延伸閱讀

- <https://go.dev/doc/effective_go#composite_literals>  
