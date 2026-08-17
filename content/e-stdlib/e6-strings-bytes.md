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

## 本章你會建立的心智模型

`strings` / `bytes` 提供搜尋、切割、替換、`Builder` 高效拼接。與 A9 的 UTF-8 心智一起用：索引仍是 byte。

## L1 能用

```go
strings.Contains(s, "bot")
strings.Cut(s, ",")
var b strings.Builder
b.WriteString("x")
```

## L2 機制

- `Builder` 減少 `+` 造成的多次分配。  
- `bytes` 與 `strings` API 幾乎平行。  
- 正規表示式在 `regexp`（成本更高）。  

## 遊戲 Server 連結

解析 `"dx,dy"`、過濾暱稱、聊天敏感詞（簡易）。

## 練習

### 必做

1. 用 `Cut` 實作拆 `dx,dy`。  
2. 用 Builder 接 1000 段字串並 bench 對照 `+`（選）。  

## 延伸閱讀

- `strings` package  
