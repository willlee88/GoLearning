---
lessonId: "C4"
title: "context 取消樹"
description: "WithCancel/Timeout、傳遞規範、値的限制。"
volume: "c"
order: 4
level: "l2"
status: "ready"
path_required: true
tags: ["context"]
example: "examples/c04-context"
prev: "C3"
next: "C5"
---

## 本章你會建立的心智模型

`context.Context` 承載**截止時間、取消信號、請求範圍的值（克制）**。取消是一棵樹：父取消則子取消。第一個參數慣例為 `ctx`，向下傳、不存 struct 長生命（或只存衍生規則要小心）。

## Python 對照

| Python | Go |
|--------|-----|
| `asyncio.timeout` / CancelledError | `ctx.Done()` / `ctx.Err()` |
| 請求級 thread-local | context 顯式傳遞 |

## L1 能用

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

select {
case <-ctx.Done():
	return ctx.Err()
case res := <-work:
	return res, nil
}
```

範例：`examples/c04-context`。

## L2 機制

- `WithCancel` / `WithTimeout` / `WithDeadline` / `WithValue`。  
- `WithValue` 只放請求級 metadata（trace id），不要塞可選參數。  
- 函式入口應檢查 `ctx.Err()` 於慢迴圈。  

## 請丟掉的 Python 習慣

1. 全域旗標 `stopping=True`。  
2. 把 context 塞進無所不在的單例。  

## 遊戲 Server 連結

連線 ctx、房間 ctx、單次匹配 ctx；關服時 root cancel → 擴散。

## 練習

### 必做

1. 跑 `examples/c04-context`。  
2. 用 timeout 包一個「最多等 200ms 的假 IO」。  

## 延伸閱讀

- <https://go.dev/blog/context>  
