---
lessonId: "B1"
title: "error 介面與哨兵錯誤"
description: "error 是值；哨兵與 errors.Is。"
volume: "b"
order: 1
level: "l2"
status: "ready"
path_required: true
tags: ["errors"]
example: "examples/b02-room-errors"
prev: "B0"
next: "B2"
---

## 本章你會建立的心智模型

`error` 是單一方法的介面：`Error() string`。業務上常用 **sentinel error**（`var ErrX = errors.New(...)`）表達穩定條件，呼叫端用 `errors.Is` 判斷，而不是比字串。

## Python 對照

| Python | Go |
|--------|-----|
| 例外型別階層 | 哨兵值 / 自訂型別 + `Is/As` |
| `except RoomFull` | `errors.Is(err, ErrRoomFull)` |

## L1 能用

```go
var ErrRoomFull = errors.New("room full")

func Join(n, cap int) error {
	if n >= cap {
		return ErrRoomFull
	}
	return nil
}

if errors.Is(err, ErrRoomFull) {
	// 回傳協定碼 ROOM_FULL
}
```

## L2 機制

- 成功：`err == nil`。  
- 字串比對易碎；用 `Is`。  
- 哨兵應穩定、少而精；細節用 wrap 加上下文。

## 請丟掉的 Python 習慣

1. 用例外訊息字串當控制流。  
2. 裸 catch 一切。  

## 遊戲 Server 連結

客戶端可見錯誤 → 穩定錯誤碼；內部日誌 → wrap 上下文（房間 ID、玩家 ID）。

## 練習

### 必做

1. 閱讀 `examples/b02-room-errors` 的哨兵定義。  
2. 新增 `ErrNotReady`。  

## 延伸閱讀

- <https://pkg.go.dev/errors>  
