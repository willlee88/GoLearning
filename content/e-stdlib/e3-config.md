---
lessonId: "E3"
title: "flag 與環境變數設定"
description: "行程設定的務實做法。"
volume: "e"
order: 3
level: "l2"
status: "ready"
path_required: false
tags: ["config"]
example: "examples/e03-flag"
prev: "E2"
next: "E4"
---

## 本章你會建立的心智模型

12-factor：設定走環境變數；本機開發可用 flag。不要把秘密寫死在 repo。

## L1 能用

```go
addr := flag.String("addr", ":8080", "listen")
flag.Parse()
```

```go
os.Getenv("ADDR")
```

範例：`examples/e03-flag`。

## 遊戲 Server 連結

Arena Mini 的 `ADDR`、`WEB_DIR`。

## 練習

### 必做

1. 為小工具加 `-n` 次數 flag。  

## 延伸閱讀

- `flag` package  
