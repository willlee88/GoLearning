---
lessonId: "P0.3"
title: "型別思維重建"
description: "從動態型別切到靜態型別、零值與明確轉換。"
volume: "p0"
order: 3
level: "l2"
status: "ready"
path_required: true
tags: ["python-bridge", "types"]
example: ""
prev: "P0.2"
next: "P0.4"
---

## 本章你會建立的心智模型

在 Go 裡，**型別是編譯期契約**：什麼能賦值、什麼能呼叫、介面是否滿足，多數在 `go build` 時決定。零值（zero value）是設計的一部分——「沒初始化」常常仍是可用的預設狀態，而不是 `None` 地雷（但仍有 `nil` 指標／介面陷阱，後續會專講）。

## Python 對照

| Python | Go |
|--------|-----|
| 執行期 `TypeError` | 編譯期型別錯誤 |
| `None` 處處可能 | 指標／slice／map／channel／interface 可 nil；基本型別有零值 |
| 隱式 bool 真理值很寬鬆 | 條件必須是 `bool` |
| `int` 任意精度 | `int` 是固定寬度（平台相關）；大數用 `math/big` |
| 註解可選 | 型別大多可推斷，但 API 邊界要清晰 |

## L1 能用

```go
var n int          // 0
var s string       // ""
var ok bool        // false
name := "player"   // 短宣告，函式內

// 無隱式數值轉換
var x int32 = 1
var y int64 = int64(x)
```

## L2 機制

### 零值哲學

結構體欄位若未設定，會是各型別零值。這讓 `var room Room` 常常能直接用（例如計數器從 0 開始），但也要求你**想清楚零值是否「安全預設」**。

### 轉換必須明確

Python 有時靠運算隱式推進型別；Go 要求 `T(v)` 形式的轉換，避免安靜的精度／語意流失。

### 型別不只是「標籤」

方法集、介面滿足、可否取址，都跟型別定義綁在一起（A、B 卷展開）。

## L3 深潛（可選）

- 底層型別（underlying type）與自定義型別的轉換規則。
- `int` 在 32/64 位元平台的差異；協定中偏好 `int32`/`int64`。

## 請丟掉的 Python 習慣

1. 用同一個變數先裝 `int` 再裝 `str`。
2. 依賴 truthy/falsy 的「空容器算 False」——在 Go 要寫清楚 `len(s) == 0` 或 `s == ""`。
3. 把 `None` 檢查文化原封不動搬來，卻忽略「零值已是合法狀態」的設計。

## 遊戲 Server 連結

封包欄位、實體 ID、tick 序號都需要**穩定、明確的寬度與符號性**。動態型別在協定邊界容易變成線上炸裂；靜態型別把問題前移到編譯期。

## 練習

### 必做

1. 寫一個 `Player` struct（`ID int64`, `Name string`, `Score int`），印出零值實例。
2. 試著把 `int32` 直接賦給 `int64` 變數，閱讀編譯器錯誤訊息。

### 選做

1. 比較 Python `int` 與 Go `int64` 在「超大整數」上的行為差異。

## 常見坑與如何看見

- 混用 `int` 與 `int64`：讓編譯器教你；協定層統一別名如 `type EntityID int64`。

## 延伸閱讀

- <https://go.dev/ref/spec#Types>
- <https://go.dev/blog/declaration-syntax>
