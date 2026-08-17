---
lessonId: "A2"
title: "package 與匯出規則"
description: "編譯單位、可見性（大寫匯出）與目錄切分。"
volume: "a"
order: 2
level: "l2"
status: "ready"
path_required: true
tags: ["packages"]
example: "examples/a02-packages"
prev: "A1"
next: "A3"
---

## 本章你會建立的心智模型

**Package** 是 Go 的編譯與封裝邊界：同一目錄、同一 `package` 名。名稱**大寫開頭**的識別字可被其他 package 使用（匯出）；小寫則是套件私有。這比 Python 的底線約定更硬——編譯器替你守門。

## Python 對照

| Python | Go |
|--------|-----|
| 檔案／套件 + `__init__.py` | 目錄 = 一 package（慣例） |
| `_private` 約定 | 小寫 = 真·私有（跨 package） |
| `from x import *` | 顯式 `import`；沒有 star import 文化 |
| 循環 import 有時能跑 | **循環匯入直接編譯失敗** |

## L1 能用

```text
mymod/
  go.mod
  main.go          # package main
  player/
    player.go      # package player
```

```go
// player/player.go
package player

type Player struct {
	Name string // 匯出
	hp   int    // 未匯出
}

func New(name string, hp int) Player {
	return Player{Name: name, hp: hp}
}

func (p Player) HP() int { return p.hp }
```

```go
// main.go
package main

import (
	"fmt"
	"example.com/mymod/player"
)

func main() {
	p := player.New("Ada", 100)
	fmt.Println(p.Name, p.HP())
	// fmt.Println(p.hp) // 編譯錯誤
}
```

範例：`examples/a02-packages`。

## L2 機制

### 匯入路徑

`import "modulepath/player"` —— module path 來自 `go.mod`，後面接目錄。

### 可見性粒度

- 套用於：型別、函式、欄位、方法、常數、變數。  
- **同一 package 不同檔案** 可視為同一命名空間（可互相存取小寫）。  

### 避免循環

若 `room` 需要 `player`、`player` 又需要 `room`，把共用型別下沉到第三 package（例如 `game/types`），或改用介面在邊界解耦。

## L3 深潛（可選）

- `internal/` 目錄：只允許父樹匯入（編譯器強制）。  
- `package player_test` 黑箱測試 vs 白箱 `package player`。

## 請丟掉的 Python 習慣

1. 用檔名當唯一模組邊界、卻沒有清晰 API。  
2. 到處 `from .x import y` 造成隱性循環。  
3. 把「內部欄位」靠文件約定——在 Go 用小寫。

## 遊戲 Server 連結

建議切分：

- `internal/session`  
- `internal/room`  
- `internal/protocol`  

入口 `cmd/server` 只負責組裝與生命週期。

## 練習

### 必做

1. 跑 `examples/a02-packages` 的 `go test` / `go run`。  
2. 試著在 `main` 讀取未匯出欄位，閱讀編譯器訊息。  

### 選做

1. 新增 `internal/id` 套件產生玩家 ID。  

## 常見坑與如何看見

- **package 名與目錄名不一致**（可運行但混亂）——保持一致。  
- **main 太大**：邏輯應下沉 package 才好測。  

## 延伸閱讀

- <https://go.dev/doc/code#Organization>  
