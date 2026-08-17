---
lessonId: "C3"
title: "select 與超時"
description: "多路等待、timeout、default 非阻塞。"
volume: "c"
order: 3
level: "l2"
status: "ready"
path_required: true
tags: ["select", "timeout"]
example: "examples/c03-select"
prev: "C2"
next: "C4"
---

## 本章你會建立的心智模型

`select` 同時等待多個 channel 操作，隨機選擇可進行的分支。搭配 `time.After` 或 `context` 做超時；`default` 做非阻塞嘗試。房間 tick 迴圈常是「收命令 / 收 tick / 收關閉」三路 select。

## Python 對照

| Python | Go |
|--------|-----|
| `asyncio.wait` | `select` |
| `queue.get(timeout=)` | `select` + timer |

## L1 能用

```go
select {
case cmd := <-inbox:
	handle(cmd)
case <-time.After(50 * time.Millisecond):
	tick()
case <-ctx.Done():
	return ctx.Err()
}
```

範例：`examples/c03-select`。

## L2 機制

- 多個 ready 時**伪随机**選一支。  
- `time.After` 在熱迴圈可能漏計時器——用 `time.NewTimer` 重用（進階）。  
- `default`：非阻塞 poll。

## 請丟掉的 Python 習慣

1. busy-loop sleep 輪詢替代事件。  

## 遊戲 Server 連結

固定 tick：`ticker.C`；同時吸收玩家命令；關房用 `ctx.Done()`。

## 練習

### 必做

1. 跑 `examples/c03-select`。  
2. 改成 100ms tick 印出計數，3 秒後結束。  

## 延伸閱讀

- <https://go.dev/ref/spec#Select_statements>  
