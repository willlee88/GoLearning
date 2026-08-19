---
lessonId: "A9"
title: "字串、rune 與 []byte：別把「字」當 byte"
description: "string 是不可變 UTF-8；索引拿到的是 byte。處理中文暱稱與聊天要用 rune。"
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

## 這章你會搞懂什麼

Go 的 `string` 是**不可變的 UTF-8 byte 序列**。  
`s[i]` 拿到的是 **byte**，不是 Python 那種「字元」。要正確處理多語文，請用 **rune**（大致可當 Unicode code point）或 `unicode/utf8` 套件。

遊戲裡的暱稱、聊天、協定字串，全都踩在這條線上——用錯 `len` 會直接變成亂碼或長度限制不公平。

## Python 對照

| Python | Go |
|--------|-----|
| `str` 偏 Unicode 抽象 | `string` = UTF-8 bytes |
| `s[i]` 是字元 | `s[i]` 是 byte |
| `for c in s` 走字元 | `for _, r := range s` 得到 rune |
| `bytes` / `bytearray` | `[]byte`（可變） |

## 怎麼寫

```go
s := "你好Go"
fmt.Println(len(s)) // byte 長度，不是「幾個字」

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

## 細節

### 為什麼 `len("你好")` 不是 2

因為 `len` 數的是 byte。中文常佔 3 個 byte（UTF-8），所以直觀上的「字數」要用：

- `utf8.RuneCountInString(s)`  
- 或自己 `range` 數 rune  

### `range` 字串在幹嘛

索引是 **byte offset**，值是 rune。遇到非法 UTF-8，會變成 `utf8.RuneError`——協定邊界要決定要拒絕還是替換。

### 不可變與拼接

字串不可變，迴圈裡狂 `s += x` 容易配很多次。多次拼接請用 `strings.Builder`。

### 跟 `[]byte` 互換

語意上請當「可能拷貝」來想。熱路徑別無意義地 `string` ↔ `[]byte` 來回甩。

### 進階可先略過

- Unicode 正規化（NFC／NFD）：「看起來一樣」的字不一定 byte 一樣。  
- 顯示寬度（全形／半形）跟 rune 數又不完全相同。

## 遊戲 Server 會用在哪

暱稱長度限制應以 **rune 數**（或顯示寬度）定義，並在進房前驗證。  
聊天過濾常在正規化之後、或直接在 `[]byte` 層做——重點是：**別用 byte 下刀切中文**。

## 請丟掉的舊習慣

1. 用 `len(s)` 當「幾個字」。  
2. 用 `s[:n]` 當截斷，結果切在 UTF-8 正中間變亂碼。  
3. 在協定裡假設「一個字元一個 byte」（ASCII 幻想）。

## 動手練習

### 必做

1. 跑 `examples/a09-strings`。  
2. 寫 `func RuneCount(s string) int`（可以直接包 `utf8.RuneCountInString`）。  

### 選做

1. 實作「暱稱最多 12 個 rune，且先 trim 空白」。  

## 常見坑

- **截斷切壞 UTF-8**：用 `utf8.DecodeLastRuneInString` 等輔助，或按 rune 重建。  
- **大小寫／摺疊**：多語文的大小寫規則比英文複雜，別以為 `ToLower` 能解決一切比對。  
- **字串當可變 buffer**：要改內容請用 `[]byte` 或 Builder，不要幻想原地改 string。
