---
lessonId: "H8"
title: "匹配入門"
description: "佇列、人數門檻、超時開房。"
volume: "h"
order: 8
level: "l2"
status: "ready"
path_required: true
tags: ["matchmaking"]
example: ""
prev: "H7"
next: "H9"
---

## 本章你會建立的心智模型

匹配把玩家從「在線」變成「同房」。最小版：

1. 入隊  
2. 湊滿 N 人 → 建 Room → 拉進房  
3. 超時 → 打 Bot 或取消  

Arena Mini 可用手動建房／房號代替完整匹配。

## Python 對照

Redis 佇列或記憶體 list；概念相同。

## L1 能用

```go
type Queue struct {
	mu   sync.Mutex
	wait []string
}

func (q *Queue) Enqueue(id string) (roomID string, started bool) {
	// if len>=2 create room...
}
```

## L2 機制

- 分池（段位、地區）。  
- 避免重複匹配。  
- 匹配服務無狀態 vs 房間有狀態。  

## 請丟掉的 Python 習慣

1. 雙層 for 暴力掃全服玩家每次請求。  
2. 匹配成功卻不處理拒絕進房。  

## 遊戲 Server 連結

本站：玩家輸入相同 room id = 手動匹配。自動匹配列二期。

## 練習

### 必做

1. 設計 2 人匹配的偽碼。  
2. 指出與「玩家自建房」的差異。  

## 延伸閱讀

- 規劃書 D8  
