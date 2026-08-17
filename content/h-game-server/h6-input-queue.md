---
lessonId: "H6"
title: "輸入佇列與校驗"
description: "命令入隊、範圍檢查、每 tick 套用。"
volume: "h"
order: 6
level: "l2"
status: "ready"
path_required: true
tags: ["input", "validation"]
example: "examples/h06-apply-input"
prev: "H5"
next: "H7"
---

## 本章你會建立的心智模型

輸入路徑：

```text
WS read → parse envelope → validate → enqueue(player, input)
tick → for each input: apply → clear/consume → snapshot
```

校驗失敗回 `type=error` 或靜默丟棄（要一致）。**不要**在讀 loop 做重物理。

## Python 對照

生產者-消費者佇列；Go channel 或 room 內 slice+mutex。

## L1 能用

```go
type Input struct {
	PlayerID string
	DX, DY   int // -1,0,1
}

func (r *Room) PushInput(in Input) error {
	if r.Phase != PhasePlaying {
		return ErrInvalidPhase
	}
	if abs(in.DX) > 1 || abs(in.DY) > 1 {
		return ErrBadInput
	}
	r.inbox = append(r.inbox, in)
	return nil
}
```

範例：`examples/h06-apply-input`。

## L2 機制

- 每玩家每 tick 可合併為「最後輸入」或「全部套用」。  
- 速率限制：每秒最多 N 個命令。  
- 序列號防重放（進階）。  

## 請丟掉的 Python 習慣

1. 讀到封包立刻改 HP 並 save DB。  
2. 信任客戶端 dt。  

## 遊戲 Server 連結

Arena Mini：`type=input` payload `"1,0"` 表示 dx,dy。

## 練習

### 必做

1. 跑 `examples/h06-apply-input`。  
2. 拒絕 |dx|>1。  

## 延伸閱讀

- H4 權威  
