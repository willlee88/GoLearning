---
lessonId: "H3"
title: "Room 生命週期狀態機"
description: "Lobby → Playing → Ended 的合法轉移。"
volume: "h"
order: 3
level: "l2"
status: "ready"
path_required: true
tags: ["room", "state-machine"]
example: "examples/h03-room-fsm"
prev: "H2"
next: "H4"
---

## 本章你會建立的心智模型

房間是**狀態機**。非法操作（Playing 時 Join、Lobby 時 Move）應回 error，而不是默默忽略或 panic。用 `switch phase` 集中表達允許的命令。

## Python 對照

| Python | Go |
|--------|-----|
| enum + 手寫轉移表 | `iota` phase + 方法校驗 |
| 例外 | `errors.Is(ErrInvalidPhase)` |

## L1 能用

```go
const (
	PhaseLobby Phase = iota
	PhasePlaying
	PhaseEnded
)

func (r *Room) Ready(id string) error {
	if r.Phase != PhaseLobby {
		return ErrInvalidPhase
	}
	// ...
}
```

範例：`examples/h03-room-fsm`。

## L2 機制

典型轉移：

```text
Lobby --(all ready & count>=min)--> Playing
Playing --(goal/timeout)--> Ended
Ended --(destroy)--> GC
```

- 誰觸發開局：最後一個 Ready、或匹配服務、或倒數。  
- Ended 後是否允許觀戰。  

## 請丟掉的 Python 習慣

1. 布林旗標滿天飛（`started`/`finished`/`closing`）互斥不清。  
2. 在 handler 各處複製 phase 判斷。  

## 遊戲 Server 連結

Arena Mini：兩人 Ready → Playing；Ended 可回 Lobby（簡化可直接銷毀）。

## 練習

### 必做

1. 跑 `examples/h03-room-fsm` 測試。  
2. 加「最少 2 人才能開局」規則。  

## 延伸閱讀

- B 卷錯誤哨兵  
