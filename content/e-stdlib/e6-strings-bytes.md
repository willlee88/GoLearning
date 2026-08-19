---
lessonId: "E6"
title: "strings 與 bytes 套件"
description: "Builder、Cut、Contains；[]byte 操作。"
volume: "e"
order: 6
level: "l2"
status: "ready"
path_required: false
tags: ["strings", "bytes"]
example: ""
prev: "E5"
next: "E7"
---

## 這章你會搞懂什麼

處理文字／位元組時，你會一直用到標準庫的 **`strings`** 與 **`bytes`**：找不找得到、怎麼切開、怎麼替換、怎麼**高效拼接**。

記得跟 A9 的心智一起用：Go 字串是 **UTF-8 位元組序列**；用數字索引時是**位元組下標**，不是「第幾個字元」。`strings`／`bytes` 的函式多半也是以位元組／子字串為單位（rune 相關 API 會另外標）。

## 先跟 Python 對一下

| Python | Go | 注意 |
|--------|-----|------|
| `"a" in s` | `strings.Contains(s, "a")` | |
| `s.split(",", 1)` | `strings.Cut(s, ",")` 或 `SplitN` | `Cut` 回 `(before, after, found)`，很好用 |
| `"".join(parts)` | `strings.Builder` 或 `strings.Join` | 熱路徑別狂用 `+` 拼很大串 |
| `s.replace(...)` | `strings.Replace` / `ReplaceAll` | |
| `bytes`／`bytearray` | `[]byte` + `bytes` 套件 | `bytes` 與 `strings` API 幾乎平行 |

## 怎麼寫（能跑的最小例子）

```go
strings.Contains(s, "bot")

before, after, ok := strings.Cut(s, ",")
if ok {
	// s == "1,-1" → before="1", after="-1"
}

var b strings.Builder
b.WriteString("hello")
b.WriteByte(' ')
b.WriteString("room")
out := b.String()
```

`[]byte` 版本長得很像：

```go
bytes.Contains(data, []byte("bot"))
bytes.Cut(data, []byte(","))
```

什麼時候用哪個？

- 已經是 `string`、偏文字邏輯 → `strings`  
- 網路緩衝、二進位、不想一直字串↔[]byte 轉換 → `bytes`

## 為什麼這樣設計／底層在幹嘛

1. **`string` 不可變；`[]byte` 可變**  
   改內容常用 `[]byte`；要當 map 鍵、印 log 再用 string。轉換會複製，熱路徑別無腦來回轉。

2. **`Builder` 減少分配**  
   `s = s + chunk` 在迴圈裡會一直配置新字串。`Builder` 內部緩衝可成長，最後一次 `String()`。小次數拼接其實沒差；次數多、片段長再在意。

3. **`Cut` vs `Split`**  
   只切一刀時 `Cut` 語意清楚，不必拿到 slice 再 `[0]` `[1]`。

4. **正規表示式在 `regexp`**  
   強大多了，也貴。暱稱過濾、簡單 token 能用 `strings` 就先別上 regex。

## 遊戲 Server 會用在哪

- 解析超簡單的輸入：`"dx,dy"` → 兩個整數。  
- 暱稱長度、是否含空白、簡易敏感詞：`Contains`／`TrimSpace`。  
- 組 log、組廣播文案：`Builder` 或 `fmt`（看情境）。  
- 二進位協定夾雜 ASCII 指令時用 `bytes`。

## 請丟掉的舊習慣

1. **把字串當「字元陣列」用整數 index 亂切中文**——可能切在 UTF-8 正中間。需要字元維度用 `range` rune 或 `utf8` 套件。  
2. **為了「看起來快」一開始就上怪招**——先正確，再用 benchmark 證明。  
3. **每次檢查都編譯新的 `regexp`**——要編譯一次存起來。

## 動手練習

### 必做

1. 用 `strings.Cut` 拆 `"1,-1"` 這種 `dx,dy`（再自己 `strconv.Atoi`）。  
2. 寫一個函式：暱稱去掉首尾空白後，長度是否在 1～12。  

### 選做

1. 用 `Builder` 接 1000 段字串，對照 `+` 做個超小 benchmark（呼應 D4）。  

## 常見坑

- **`Contains` 大小寫**：要無視大小寫用 `EqualFold`／先摺疊，不是直接 Contains。  
- **`Trim` 的 cutset 不是「字串前綴」**：是「集合裡的字元」都剝掉；剝前綴用 `TrimPrefix`。  
- **修改 `[]byte` 卻忘了它可能跟字串共用底層**（進階）：對字面值轉換來的行為要小心；新手分開配置較安心。

## 延伸閱讀

- <https://pkg.go.dev/strings>  
- <https://pkg.go.dev/bytes>  
