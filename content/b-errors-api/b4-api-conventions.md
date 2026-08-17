---
lessonId: "B4"
title: "API 慣例與建構"
description: "接受 interface、回傳 struct；建構函式與可測邊界。"
volume: "b"
order: 4
level: "l2"
status: "ready"
path_required: true
tags: ["api", "design"]
example: "examples/b02-room-errors"
prev: "B3"
next: "C0"
---

## 本章你會建立的心智模型

好的 Go API 讀起來無聊而清楚：

1. 依賴用**小 interface** 注入  
2. 回傳**具體型別**  
3. 建構函式 `NewX(...) (*X, error)` 保證不變式  
4. 錯誤可 `Is/As`  

## Python 對照

| Python | Go |
|--------|-----|
| DI 框架 / 建構子注入 | 手寫 New + 介面欄位 |
| 屬性可事後亂設 | 未匯出欄位 + 建構保證 |

## L1 能用

```go
type Store interface {
	Save(*Room) error
}

type Server struct {
	store Store
}

func NewServer(store Store) (*Server, error) {
	if store == nil {
		return nil, errors.New("store nil")
	}
	return &Server{store: store}, nil
}
```

## L2 機制

- 可測：mock `Store`。  
- 避免在 `init()` 做有失敗可能的重邏輯。  
- 設定用 struct，必要時 functional options（進階，勿濫用）。

## 請丟掉的 Python 習慣

1. 全域單例到處 import。  
2. 半初始化物件。  

## 遊戲 Server 連結

`NewRoom(cfg Config) (*Room, error)` 驗證 capacity>0；`NewHub(log Logger)` 注入日誌。

## 練習

### 必做

1. 為 registry 或 room 加 `New` 並拒絕非法 cap。  
2. 用 interface 抽出「時鐘」以便測試超時。  

## 延伸閱讀

- <https://go.dev/doc/effective_go>  

## B 卷檢查點

你應能把「滿房／不存在／非法階段」做成哨兵 + wrap，並在測試用 `errors.Is` 鎖定行為。  
