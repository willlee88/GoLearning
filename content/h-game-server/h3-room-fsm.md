---
lessonId: "H3"
title: "Room 生命週期：用狀態機管准許的操作"
description: "Lobby → Playing → Ended。非法操作回 error，不要默默忽略或 panic；集中用 phase 判斷。"
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

## 這章你會搞懂什麼

房間（Room）不是一個無限自由的聊天室容器。  
它比較像一台**狀態機**（finite state machine，FSM）：現在處於哪個階段，就只允許某些命令。

例如：

- **Lobby（大廳）**：可以 Join、Ready；不該認真套用移動與得分  
- **Playing（對戰中）**：可以收 input；不該再讓路人 Join 進來改人數  
- **Ended（已結束）**：可以廣播結算、倒數回大廳或銷毀；不該繼續當正式對戰輸入

讀完這章，你看到任何 `ready`／`join`／`input`，都會先問：**現在 phase 允不允許？不允許時回什麼錯誤？**

## Python 對照

| Python | Go |
|--------|-----|
| `Enum` + 手寫轉移表 | `iota` 定義 phase + 方法開頭檢查 |
| `raise ValueError("bad state")` | 回傳 `error`，常用 `errors.Is(err, ErrInvalidPhase)` |
| 一堆 `if self.started and not self.finished` | 一個 `Phase` 欄位，互斥比較清楚 |

Python 裡布林旗標很容易變成：

```python
started = True
finished = False
closing = True  # …現在到底算哪種？
```

Go 教學裡我們偏向**單一 phase**，合法轉移寫在少數函式裡。

## 怎麼寫（能跑的最小例子）

範例在 `examples/h03-room-fsm`。核心長這樣：

```go
type Phase int

const (
	PhaseLobby Phase = iota
	PhasePlaying
	PhaseEnded
)

var ErrInvalidPhase = errors.New("invalid phase")

type Room struct {
	Capacity int
	Phase    Phase
	Members  map[string]bool // name -> ready
}

func (r *Room) Join(name string) error {
	if r.Phase != PhaseLobby {
		return fmt.Errorf("join: %w", ErrInvalidPhase)
	}
	// 容量、重複加入… 
	r.Members[name] = false
	return nil
}

func (r *Room) Ready(name string) error {
	if r.Phase != PhaseLobby {
		return fmt.Errorf("ready: %w", ErrInvalidPhase)
	}
	r.Members[name] = true
	return r.tryStart() // 人夠且全 ready → PhasePlaying
}
```

跑測試：

```bash
cd examples/h03-room-fsm
go test ./...
```

測試會覆蓋兩件很重要的事：

1. 兩人都 Ready → 進入 Playing  
2. Playing 時再 Join → `ErrInvalidPhase`

## 為什麼這樣設計

### 非法操作要「失敗得清楚」

三種壞做法：

| 做法 | 後果 |
|------|------|
| 默默 `return` 忽略 | 客戶端以為 Ready 成功，UI 與 Server 不一致 |
| `panic` | 一間房的壞訊息可能拖垮整個進程 |
| 在每個 handler 複製一長串 if | 漏改一處就出現「偶爾能作弊加入」 |

集中在 `Join`／`Ready`／`PushInput` 入口檢查 phase，再用 **sentinel error**（哨兵錯誤，例如 `ErrInvalidPhase`）讓上層決定：回 `type=error` 給客戶端，或記 log。

### 典型轉移長怎樣

```text
Lobby  --(人數 >= min 且全員 Ready)-->  Playing
Playing --(達成分／超時／裁判結束)-->  Ended
Ended   --(銷毀或重置)-->  GC 或回到 Lobby
```

你還要決定產品細節（寫下來，別放肚子裡）：

- **誰觸發開局**：最後一個 Ready？匹配服務喊 start？倒數 3 秒？  
- **Ended 之後**：觀戰？回 Lobby 再來一局？整房拆掉？  
- **中途有人斷線**：繼續、暫停、直接 Ended？（跟 H2 政策銜接）

### 進階可先略過

- 更嚴謹可用「事件 → 新狀態」表驅動，避免 `switch` 散落。  
- 有些遊戲有更多 phase：`Countdown`、`SuddenDeath`、`Settlement`——同一套思想，只是狀態比較多。  
- 轉移時順便啟動／停止 tick goroutine：進 Playing 才 `NewTicker`，進 Ended 要取消，避免空轉漏房。

## 遊戲 Server 會用在哪

Arena Mini：

1. 進房 → Lobby  
2. 兩人 Ready → Playing，並啟動 tick  
3. 先到分數或超時 → Ended  
4. 簡化版可數秒後回 Lobby，或直接銷毀房間  

H10 Lab 會叫你親手走完這條路。現在先確保：**相位判斷不是寫在前端「人齊了就自己 start」**——開局真相在 Server。

## 請丟掉的舊習慣

1. **布林旗標滿天飛**且互斥不清。  
2. **在每個 WS handler 複製 phase 判斷**，改一處漏五處。  
3. **前端決定開局**：人齊了客戶端廣播 `start`，Server 照單全收——這是對作弊與不同步開大門。  
4. **非法操作當成功**：UI 樂觀更新，Server 沒改，最後靠重連「修好」——其實是隱性 Bug。

## 動手練習

### 必做

1. 跑通 `examples/h03-room-fsm` 的測試。  
2. 讀 `tryStart`：確認「少於 2 人不會開局」。若你改規則，補一個測試鎖住它。  

### 選做

1. 加 `PhaseEnded`，並實作 `End()`；再加測試：Ended 時 `Ready` 必須失敗。  
2. 加「開局倒數」假 phase：`PhaseCountdown`——思考 tick 要不要在這階段就轉。  

## 常見坑

- **只在文件寫 Lobby／Playing，程式卻用兩個 bool**：文件與實作會漂。  
- **`tryStart` 回傳 error 表示「人還不夠」**：人不夠是正常等待，通常回 `nil` 即可；真錯誤（phase 不對）才回 error。看範例怎麼分。  
- **Playing 還允許 Join**：人數、配對、座位全亂，平衡與同步都會炸。  
- **轉移了 phase 卻沒啟動／停止 tick**：Lobby 空轉耗資源，或 Ended 後還在套用 input。

下一章 H4：權威 Server——客戶端只送意圖，真相只住在 Server。
