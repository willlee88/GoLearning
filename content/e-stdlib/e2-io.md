---
lessonId: "E2"
title: "io.Reader 與 Writer"
description: "串流抽象；組合優於巨型 API。"
volume: "e"
order: 2
level: "l2"
status: "ready"
path_required: false
tags: ["io"]
example: "examples/e02-io"
prev: "E1"
next: "E3"
---

## 本章你會建立的心智模型

Go I/O 核心是小介面：`Reader`、`Writer`、`Closer`。檔案、網路、buffer 都實作它們，才能 `io.Copy`、層層包裝。

## Python 對照

file-like 物件；Go 用編譯期介面強制。

## L1 能用

```go
io.Copy(os.Stdout, strings.NewReader("hi\n"))
```

範例：`examples/e02-io`。

## 遊戲 Server 連結

長度前綴幀讀寫建立在 `io.Reader`（F3）。

## 練習

### 必做

1. 用 `io.LimitReader` 限制讀取。  

## 延伸閱讀

- `io` package  
