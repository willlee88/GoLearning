---
lessonId: "G2"
title: "命令與事件分離"
description: "client 下指令；server 廣播權威結果。"
volume: "g"
order: 2
level: "l2"
status: "ready"
path_required: true
tags: ["protocol", "game-server"]
example: "examples/g01-envelope"
prev: "G1"
next: "G3"
---

## 這章你會搞懂什麼

這章是遊戲網路裡最重要的心智之一：

- **命令（Command）**：客戶端表達**意圖**——「我想往左」「我點準備」  
- **事件／狀態（Event／State）**：伺服器宣布**發生了什麼**——「權威快照長這樣」「某人加入了」

千萬不要讓客戶端直接廣播「我的 HP=999」。正確流程是：

```text
Client → command
Server → 驗證 → 修改權威狀態 → 廣播 event/state
```

這叫**權威伺服器（authoritative server）**：規則在 server，客戶端顯示結果。

## 先跟 Python 對一下

若你聽過 **CQRS**（命令與查詢責任分離），這裡是超輕量版：命令進、狀態出。Python 遊戲後端一樣該這麼做；Go 只是用型別與 `room.Apply(cmd)` 把邊界畫得更清楚，方便單元測試。

「客戶端算完再告訴伺服器存檔」在單機或信任環境可以，**多人對戰不行**。

## 怎麼寫（訊息形狀示意）

```text
Client → { "v":1, "type":"move", "payload":{"dx":-1, "dy":0} }
Server → validate → apply
Server → { "v":1, "type":"state", "payload":{"players":[...]} }  // 廣播
```

系統事件也可以分開 type，例如 `join`／`leave`／`phase`，避免全塞進一個超大 `state`（看複雜度取捨）。

程式內（概念）：

```go
events, err := room.Apply(cmd)
if err != nil {
	// 回 error 給該玩家，或忽略非法命令
	return
}
broadcast(events) // 或定期廣播全量 state
```

## 為什麼這樣設計／底層在幹嘛

1. **反作弊與一致性**  
   所有人看到的世界以 server 為準。客戶端可以做預測（先動起來），但以 server 校正為準——細節在 H 卷。

2. **命令應可驗證**  
   `dx` 只能 -1..1、人必須活著、必須在對的 phase。驗證失敗 ≠ 伺服器崩潰。

3. **冪等與序號**  
   網路會重送。同一則「開火」處理兩次就爆。可帶 `cmd_id`／序號，或讓狀態轉移天然冪等。

4. **全量 state vs delta**  
   全量好懂、好重連；delta 省流量、複雜。教學先全量，優化以後再說。

5. **測試友善**  
   `Apply` 若是純邏輯（不碰網路），可用表格測試餵命令、断言狀態——這是 H 卷房間 FSM 的伏筆。

## 遊戲 Server 會用在哪

Arena Mini：

- 命令：移動、（之後）準備、開火…  
- 事件／狀態：位置快照、得分、階段變更、系統公告  

`room.Apply(cmd) (events, error)` 這種形狀，後面會一直出現。

## 請丟掉的舊習慣

1. **客戶端算完結果，只通知 server「存檔／轉發」。**  
2. **聊天、移動、系統訊息混成一條無結構字串。**  
3. **廣播「某玩家說他自己的狀態」而不經 server 認可。**

## 動手練習

### 必做

1. 列出 Arena Mini（或 f08）目前哪些訊息是命令、哪些是事件／狀態。  
2. 設計一則 `ready` 命令，以及狀態裡的 `phase` 欄位（紙上 JSON 即可）。  

### 選做

1. 想像重送同一則 `move`：你的 `Apply` 會發生什麼？要不要序號？  

## 常見坑

- **把 `state` 當雙向聊天**：客戶端不該發權威 state。  
- **驗證只做在前端**：前端檢查是體驗，後端檢查才是安全。  
- **錯誤時沉默**：非法命令至少 log；必要時回 `type=error` 讓客戶端 UI 有反應。  

## 延伸閱讀

- 搜尋「authoritative game server」任一概述文  
- 範例信封：`examples/g01-envelope`  
