---
lessonId: "F8"
title: "Lab：WebSocket 房間廣播"
description: "F 卷檢查點：多房間、JSON 信封、廣播。"
volume: "f"
order: 8
level: "l2"
status: "ready"
path_required: true
tags: ["lab", "websocket", "game-server"]
example: "examples/f08-ws-room"
prev: "F7"
next: "G0"
---

## 這章你會搞懂什麼

這是 F 卷的**檢查點 lab**：把 F1–F7 收成一個可跑的小系統，而不是再學一個孤立 API。

你要親手摸到：

- HTTP `/healthz`（活著沒）  
- WS `/ws?room=&name=`（進房）  
- JSON 信封（至少有 `type`，常見還有 `from`／`payload`／`room`）  
- **Hub → Room → members** 廣播  
- 連線退出時**清理**成員表  

這是 Arena Mini 的縮小版；正式 demo 會再加房間列表、狀態同步等。

## 先跟 Python 對一下

你若寫過 aiohttp／Starlette 的 WS 聊天室，架構會長得很像：連線集合、房間 dict、廣播迴圈。差別在 Go 要用鎖或單一 writer 把併發講清楚，並習慣「每連線一個讀 loop goroutine」。

## 怎麼寫（先跑起來）

```bash
cd examples/f08-ws-room
go run .
# 瀏覽器開附帶的 static，或開兩個客戶端連同一 room
```

也請跑：

```bash
cd demo/arena-mini/server
go mod tidy
go run ./cmd/server
```

照 README 用瀏覽器連。目標是「雙客戶端看得到彼此訊息／狀態」。

## 機制複習清單（請一邊看程式一邊勾）

- [ ] 每連線有讀 loop  
- [ ] 寫入有串行（mutex 或同等模型）  
- [ ] 房間成員 `map` 有鎖  
- [ ] leave／斷線時 `delete`，不會廣播到幽靈連線  
- [ ] 應用訊息有 `type` 可分派  
- [ ] 至少一種錯誤路徑：壞 JSON 怎麼辦（斷線／回 error 事件）  

## 為什麼這樣設計／底層在幹嘛

1. **Hub／Room 分層**  
   Hub 管「房間表」；Room 管「這個房的人與廣播」。之後加 tick、局狀態，邏輯會落在 Room，不會全塞進 WS handler。

2. **清理跟加入一樣重要**  
   只實現 join、不實現 leave，過幾小時記憶體與廣播成本會爆。`defer` 取消註冊是好習慣。

3. **信封先簡單**  
   教學用 JSON text frame 足够。字段怎麼演進交給 G 卷；權威狀態與 tick 交給 H 卷。

4. **Lab 的意義**  
   網路卷最怕「都看過、沒串過」。這裡強迫你串一次，後面協定與遊戲邏輯才有掛載點。

## 遊戲 Server 會用在哪

這條路幾乎就是多人線上後端的胚胎：

```text
連線進入 → 認證／命名 → 進房 → 收命令 → 廣播事件／狀態 → 離房清理
```

Arena Mini 只是在上面加規則、tick、分數。

## 請丟掉的舊習慣

1. **所有人塞一個全局 list 又不加鎖。**  
2. **在 handler 裡 `time.Sleep` 或重計算卡住讀 loop。**  
3. **斷線不清理還假裝廣播成功。**

## 動手練習

### 必做（F 卷檢查點）

1. 跑通 `examples/f08-ws-room` 雙客戶端。  
2. 跑通 `demo/arena-mini`。  
3. 加一種 `type=pos`（或類似），payload 帶座標字串並廣播給同房。  

### 選做

1. HTTP `GET /rooms` 回各房人數（可對照 arena-mini）。  
2. 加 `ping`／`pong` log（呼應 F7）。  

## 常見坑

- **開了 server 卻連錯 port／路徑。**  
- **兩個分頁 room 名稱打不同**：以為 bug，其實各聊各的。  
- **改完程式沒重跑**：還在看舊 binary。  
- **race**：加 `-race` 跑測（若有測試）或壓一點進出房。

## 延伸閱讀

- `demo/arena-mini/README.md`  

下一站 **G0**：把「訊息長什麼樣、命令跟事件怎麼分」講清楚。  
