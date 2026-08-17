---
lessonId: "E7"
title: "context 套件速查"
description: "Background、TODO、With* 家族與傳值紀律。"
volume: "e"
order: 7
level: "l2"
status: "ready"
path_required: false
tags: ["context", "stdlib"]
example: ""
prev: "E6"
next: "E8"
---

## 本章你會建立的心智模型

C4 已講取消樹；本章當**標準庫速查**。API 邊界第一參數 `ctx context.Context`。值只放請求級 metadata。

## L1 能用

```go
ctx := context.Background()
ctx, cancel := context.WithTimeout(ctx, time.Second)
defer cancel()
```

## L2 機制

| 函式 | 用途 |
|------|------|
| Background | root |
| TODO | 暫不知用哪個 |
| WithCancel | 手動取消 |
| WithTimeout/Deadline | 截止 |
| WithValue | 稀疏 metadata |

## 請丟掉的 Python 習慣

1. 把可選參數塞進 context value。  

## 遊戲 Server 連結

HTTP `r.Context()`、關服 root cancel、DB 呼叫帶 ctx。

## 練習

### 必做

1. 寫一個 100ms timeout 的 select 示範。  

## 延伸閱讀

- Go Blog: Context  
