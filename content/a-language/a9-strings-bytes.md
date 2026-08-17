---
lessonId: "A9"
title: "字串、rune 與 bytes"
description: "UTF-8、不可變字串、[]byte 與字串處理邊界。"
volume: "a"
order: 9
level: "l2"
status: "ready"
path_required: true
tags: ["strings", "unicode"]
example: "examples/a09-strings"
prev: "A8"
next: "A10"
---

## 本章你會建立的心智模型

Go 的 `string` 是**不可變的 UTF-8 byte 序列**。索引拿到的是 byte，不是「字元」；要正確處理多語文用 `rune`（Unicode code point）或 `utf8`/`unicode/utf8` 套件。遊戲裡聊天、暱稱、協定字串都踩在這條線上。

## Python 對照

| Python | Go |
|--------|-----|
| `str` Unicode 抽象 | `string` = UTF-8 bytes |
| `s[i]` 是字元 | `s[i]` 是 byte |
| `for c in s` 字元 | `for _, r := range s` 得 rune |
| `bytes` / `bytearray` | `[]byte`（可變） |

## L1 能用

```go
s := "你好Go"
fmt.Println(len(s)) // byte 長度，不是字元數

for i, r := range s {
	fmt.Printf("%d %c\n", i, r)
}

b := []byte(s)
s2 := string(b)
```

```go
import "strings"
strings.HasPrefix(name, "bot_")
strings.ToLower(name)
```

範例：`examples/a09-strings`。

## L2 機制

- 字串不可變 → 拼接多次用 `strings.Builder`。  
- `range` 字串：索引是 byte offset，值是 rune。  
- 非法 UTF-8 在 range 時會變成 `utf8.RuneError`。  
- 與 `[]byte` 轉換可能拷貝（實作細節／編譯器優化另說，語意上當拷貝想）。

## L3 深潛（可選）

- 正規化（NFC/NFD）與「看起來一樣的字」。  
- 效能：熱路徑避免無謂 `string([]byte)` 來回。

## 請丟掉的 Python 習慣

1. 用 `len(s)` 當「幾個字」。  
2. 用 byte 索引切開中文。  
3. 在協定裡假設「一字一碼」。

## 遊戲 Server 連結

暱稱長度限制應以 **rune 數或顯示寬度** 定義，並在進房前驗證；聊天過濾在 `[]byte` 或正規化後做。

## 練習

### 必做

1. 跑 `examples/a09-strings`。  
2. 寫 `func RuneCount(s string) int`（可用 `utf8.RuneCountInString`）。  

### 選做

1. 實作「暱稱最多 12 個 rune，且 trim 空白」。  

## 常見坑與如何看見

- 截斷 `s[:n]` 切在 UTF-8 中間 → 亂碼；用 `utf8.DecodeLastRuneInString` 等輔助。  

## 延伸閱讀

- <https://go.dev/blog/strings>  
