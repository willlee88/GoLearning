---
lessonId: "A10"
title: "struct 與 tag：把一筆資料描述清楚"
description: "用 struct 建模；tag 掛 json 之類中繼資料。組合靠嵌入，不是繼承樹。"
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

## 這章你會搞懂什麼

`struct`（結構體）是 Go 描述「一筆資料長什麼樣」的主力。玩家、房間設定、指令 payload，幾乎都是它。

**嵌入（embedding）** 用來做組合（composition），不是 class 繼承。  
**tag** 則把中繼資料掛在欄位上——最常見是 `json:"name"`，讓編碼格式跟領域模型對齊。

（嵌入的完整取捨在 A15 會再深講；這章先會用、知道它不是繼承就好。）

## Python 對照

| Python | Go |
|--------|-----|
| `@dataclass`／普通 class | `struct` |
| 繼承 | 嵌入 + 介面（後面） |
| `asdict`／pydantic | `encoding/json` + tag |

## 怎麼寫

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
type Timed struct{ time.Time } // 匿名嵌入；方法可能被提升
```

範例：`examples/a10-struct-tags`。

## 細節

### 大寫欄位才會進 JSON

`encoding/json` 只匯出**大寫開頭**欄位。小寫欄位會「安靜地消失」——超常見坑，除錯時先看這個。

### 嵌入 ≠ 子型別

嵌入可以提升方法與欄位，讓你少寫轉發。但外層型別**不能**當成內層型別傳進函式——沒有「Server is-a Logger」這種繼承多型。

欄位名衝突時，要顯式寫 `s.Logger.Log(...)`。

### 零值 struct 很好用

`var p Player` 立刻有合法的欄位零值。設計時想想：全零的 Player／Room 算不算安全預設？

### 進階可先略過

- 記憶體對齊（alignment）與欄位排序影響 struct 大小。  
- `json:"-"`、`omitempty` 的細節與空值語意。

## 遊戲 Server 會用在哪

`Player`、`RoomConfig`、`Command` 幾乎都是 struct；對外協定用 tag 對齊欄位名。  
小模式像「裡面藏一把鎖」：

```go
type room struct {
	mu sync.Mutex
	// 其他未匯出狀態…
}
```

注意：含鎖的 struct **不要隨便值拷貝**出去，否則會拷到一把「另一把鎖」。

## 請丟掉的舊習慣

1. 深繼承樹 `BaseEntity → MovingEntity → Player`。  
2. 長期用 `dict`／`map[string]any` 當領域模型（真的動態 schema 另說）。  
3. 以為 tag 是註解——它會被反射讀走，寫錯就影響編碼。

## 動手練習

### 必做

1. 跑 `examples/a10-struct-tags`。  
2. 幫 `Room` 加上 `Capacity int \`json:"capacity"\``，做 JSON round-trip。  

### 選做

1. 比較嵌入 `time.Time` 與欄位 `CreatedAt time.Time` 的 JSON 形狀差異。  

## 常見坑

- **小寫欄位在 JSON 裡蒸發**。  
- **嵌入讓 JSON 形狀變扁平或怪掉**：先寫測試釘住輸出。  
- **把嵌入當繼承多型**：型別系統不會讓你這樣混用。
