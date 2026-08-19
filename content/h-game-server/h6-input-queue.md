---
lessonId: "H6"
title: "輸入佇列與校驗：先排隊，tick 再套用"
description: "WS 讀到的是意圖：validate → enqueue；Tick 時再 apply。範圍檢查、相位檢查、合併策略。"
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

## 這章你會搞懂什麼

H4 說「送意圖」，H5 說「固定 tick」。這章把中間那段接起來：

```text
WS read
  → 解析信封（envelope）
  → 校驗（phase、範圍、是不是房裡的人）
  → enqueue（丟進該房的輸入佇列）
tick
  → 取出／合併輸入
  → apply（改權威狀態）
  → clear
  → snapshot → 交給 hub 廣播
```

你要記住一句話：

> **不要在讀 loop 裡做重物理。** 讀 loop 負責「收件與安检」，Tick 負責「世界往前走」。

## Python 對照

這就是經典的**生產者－消費者**：

| 角色 | Python 直覺 | Go 常見做法 |
|------|-------------|-------------|
| 生產者 | 每個 WS handler `queue.put(input)` | 讀 loop 呼叫 `room.PushInput` |
| 佇列 | `asyncio.Queue`／`queue.Queue` | channel，或 room 內 `[]Input` + mutex |
| 消費者 | 背景 task 定時 `get_nowait` 一批 | tick case 裡消費 inbox |

差別只是：遊戲還多了 **phase** 與 **數值校驗**，失敗要有一致政策（回 error 或靜默丟棄——選一個，別混）。

## 怎麼寫（能跑的最小例子）

範例：`examples/h06-apply-input`。

```go
type Input struct {
	Player string
	DX, DY int // 期望是 -1, 0, 1
}

func (r *Room) PushInput(in Input) error {
	if r.Phase != PhasePlaying {
		return ErrInvalidPhase
	}
	if abs(in.DX) > 1 || abs(in.DY) > 1 {
		return fmt.Errorf("%w: dx dy must be -1..1", ErrBadInput)
	}
	if _, ok := r.Players[in.Player]; !ok {
		return errors.New("unknown player")
	}
	r.inbox = append(r.inbox, in)
	return nil
}

// Tick：同一玩家多筆輸入時，最後一筆生效（範例策略）
func (r *Room) Tick() {
	last := map[string]Input{}
	for _, in := range r.inbox {
		last[in.Player] = in
	}
	r.inbox = r.inbox[:0]
	for name, in := range last {
		p := r.Players[name]
		p.X = clamp(p.X+in.DX*r.Speed, 0, r.Width)
		p.Y = clamp(p.Y+in.DY*r.Speed, 0, r.Height)
	}
}
```

跑測試：

```bash
cd examples/h06-apply-input
go test ./...
```

裡面已經有：

- 合法 input → 座標依 speed 前進  
- `|dx|>1` → `ErrBadInput`  
- Lobby 階段送 input → `ErrInvalidPhase`

## 校驗要檢什麼（最少集）

| 檢查 | 為什麼 |
|------|--------|
| Phase 是 Playing | 沒開局就移動會破壞流程 |
| 玩家在房內 | 防冒名／過期 session |
| dx,dy 範圍 | 防一幀飛全圖 |
| （建議）速率限制 | 防每秒上千包灌 inbox |
| （建議）數值 clamp 在 apply | 雙重保險，座標不穿牆出界 |

失敗時怎麼辦要**全專案一致**：

- 回 `type=error` 給該玩家——好除錯、好做 UI  
- 或靜默丟棄——省流量，但客戶端可能一頭霧水  

兩邊混用會讓你以為「偶發不同步」。

## 一個 tick 內多筆輸入怎麼辦

常見策略：

| 策略 | 意思 | 適合 |
|------|------|------|
| 最後輸入生效 | 同玩家只留最新方向 | 連續移動、搖桿 |
| 全部依序套用 | 每筆都模擬 | 需要精細輸入的玩法 |
| 合成向量 | 加總後再 clamp | 少數特例 |

Arena Mini／本範例用「最後輸入」——簡單、對方向鍵夠用。  
**重點是寫死策略並測試**，不要有時合併有時不合併。

### 進階可先略過

- **輸入序號**（sequence number）：防重放、便於和解與丟包偵測。  
- **每玩家每秒 N 個命令**的 token bucket。  
- channel 取代 slice：`select` 到 tick 時 `drain`；注意 goroutine 與關閉時機。  
- 讀 loop 持鎖時間要短：校驗完複製一份 input 再進佇列，別在鎖內做廣播。

## 遊戲 Server 會用在哪

Arena Mini：客戶端送 `type=input`，payload 類似 `"1,0"` 表示 dx,dy。  
Server：`PushInput` → 等 tick → 改座標 → 廣播 `state`。  
Lobby 時的 input 應被拒絕或忽略——H10 檢查清單會點名這件事。

## 請丟掉的舊習慣

1. **讀到封包立刻改 HP 並 `save` DB**——熱路徑又慢又難回滾。  
2. **信任客戶端帶來的 `dt`**（「我這邊過了 5 秒所以我該走很遠」）——請用 Server tick。  
3. **校驗只寫在前端**——前端校驗是體驗，不是安全。  
4. **inbox 無限成長**——攻擊者或 Bug 都能把記憶體打滿。

## 動手練習

### 必做

1. 跑通 `examples/h06-apply-input`。  
2. 確認 `|dx|>1` 會被拒；若你改範圍規則，補測試鎖住。  

### 選做

1. 改成「全部依序套用」，比較連按兩次右鍵在同 tick 內與「最後輸入」的座標差。  
2. 幫 `PushInput` 加「每玩家佇列最多 K 筆」，超過回 error。  

## 常見坑

- **在 `PushInput` 裡直接改 `Player.X`**：佇列形同虛設，tick 節奏被打穿。  
- **忘記清空 inbox**：每 tick 重複套用舊輸入，角色自己滑出去。  
- **錯誤被吃掉**：`PushInput` 失敗但 hub 沒回傳 `error` 給客戶端，兩邊認知分裂。  
- **用客戶端時間戳當權威順序**：時鐘可被改；以 Server 收到順序或 Server 發的序號為準較穩。

下一章 H7：算完狀態之後，怎麼同步給大家——全量快照還是 delta？
