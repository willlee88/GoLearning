---
lessonId: "A2"
title: "package：誰能看見誰（大寫才匯出）"
description: "同一個目錄就是編譯邊界；名稱大寫才能給別的 package 用——比 Python 底線約定硬很多。"
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

## 這章你會搞懂什麼

**Package（套件）** 是 Go 的編譯與封裝邊界：通常「一個目錄、同一個 `package` 名稱」。  
名稱**大寫開頭**的識別字才能被其他 package 使用，這叫**匯出（export）**；小寫就是套件私有。

這比 Python 的 `_private` 約定硬——不是君子協定，是**編譯器直接擋**。你少猜很多「到底能不能碰內部欄位」。

## Python 對照

| Python | Go |
|--------|-----|
| 檔案／套件 + `__init__.py` | 目錄 ≈ 一個 package（慣例） |
| `_private` 靠約定 | 小寫 = 跨 package 真的看不到 |
| `from x import *` | 顯式 `import`；沒有 star import 文化 |
| 循環 import 有時還能跑 | **循環匯入直接編譯失敗** |

## 怎麼寫

目錄長這樣：

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
	Name string // 匯出：別的 package 看得到
	hp   int    // 未匯出：只有同 package 碰得到
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
	// fmt.Println(p.hp) // 編譯錯誤：hp 未匯出
}
```

範例：`examples/a02-packages`。建議你故意打開 `p.hp` 看編譯器怎麼罵。

## 細節

### 匯入路徑怎麼拼

`import "modulepath/player"` —— 前面的 module path 來自 `go.mod`，後面接**目錄路徑**。  
所以資料夾怎麼切，幾乎就決定 import 怎麼寫。

### 可見性套用在什麼上

型別、函式、欄位、方法、常數、變數——只要是識別字，都吃同一套大小寫規則。

同 package、不同 `.go` 檔，可視為同一個命名空間：小寫彼此還是看得到。這讓你可以拆檔，卻不必把內部 API 公開。

### 為什麼循環匯入要擋

因為編譯與初始化順序會變得難推理。若 `room` 需要 `player`、`player` 又需要 `room`，常見解法是：

- 把共用型別下沉到第三 package（例如 `game/types`）  
- 或在邊界改用 **interface（介面）** 解耦（後面 A12 會講）

### 進階可先略過

- `internal/` 目錄：只允許「父樹」匯入，編譯器強制。很適合放不該被外部專案當庫用的碼。  
- `package player_test` 是黑箱測試；`package player` 是白箱（可測未匯出符號）。

## 遊戲 Server 會用在哪

實務上常這樣切（名稱可調，精神是「邊界清楚」）：

- `internal/session` — 連線／玩家工作階段  
- `internal/room` — 房間狀態與規則  
- `internal/protocol` — 封包／訊息形狀  

入口 `cmd/server` 只負責組裝與生命週期，不要把所有邏輯塞進 `main.go`——不然後面超難測。

## 請丟掉的舊習慣

1. 用「檔名」當唯一模組邊界，卻沒有清楚 API。  
2. 到處相對 import，搞出隱性循環，還指望執行期運氣。  
3. 內部欄位只靠文件寫「請勿使用」——在 Go 請直接小寫。

## 動手練習

### 必做

1. 跑 `examples/a02-packages` 的 `go test`／`go run`。  
2. 在 `main` 讀取未匯出欄位，認真讀編譯器訊息。  

### 選做

1. 新增 `internal/id` 套件，負責產生玩家 ID。  

## 常見坑

- **package 名跟目錄名不一致**：有時能跑，但團隊會崩潰；盡量一致。  
- **`main` 肥到爆炸**：邏輯下沉到 package 才好單測。  
- **想 import 別 package 的小寫符號**：這不是設定問題，是語言規則。
