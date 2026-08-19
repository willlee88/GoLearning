---
lessonId: "H2"
title: "Session 與連線：斷線不等於玩家蒸發"
description: "Connection 是 socket；Session 是玩家在線上下文。進房、斷線、重連要分開想。"
volume: "h"
order: 2
level: "l2"
status: "ready"
path_required: true
tags: ["session"]
example: ""
prev: "H1"
next: "H3"
---

## 這章你會搞懂什麼

很多新手會把「WebSocket 物件」直接當「玩家」。  
斷線的一瞬間，就把名字、分數、座位、房間成員全刪光——玩到一半網路抖一下，整個人從局裡蒸發，體驗很糟。

這章要你分開兩件事：

- **Connection（連線）**：底層的 socket／WebSocket，會斷、會重連、可以換一條新的  
- **Session（工作階段）**：這個玩家在線的上下文——是誰、在哪一房、座位、暫存狀態  

讀完你要能回答：「斷線 5 秒內重連，局內座標該怎麼辦？」——而且答案不該是「反正砍掉重來」這唯一選項。

## Python 對照

| Python 裡你可能看過 | 遊戲 Server 要更顯式的地方 |
|--------------------|---------------------------|
| Flask／Django 的 session cookie | 即時局不能只靠 cookie；要有記憶體中的 Session 物件 |
| `websocket` 連線物件本身 | 連線只是運輸；身份與局內狀態另存 |
| 斷線 = `finally: del users[ws]` | 可能先標 `Disconnected`，留重連窗口 |

Web 應用裡「關掉分頁 = 登出」有時沒關係。對戰中途不行——你至少要決定：**保留多久、保留什麼、重連後怎麼綁回座位**。

## 怎麼寫（最小模型）

先別追求完美框架，用一個清楚的 struct 就夠建立感覺：

```go
// Conn 先抽象成介面：方便之後換實作或寫測試假物件。
type Conn interface {
	WriteJSON(v any) error
	Close() error
}

type Session struct {
	ID       string // 穩定一點的身份（之後可改成 token 解出的 user id）
	Name     string // 顯示名稱
	RoomID   string // 空字串 = 還在大廳／未進房
	Conn     Conn   // 目前這條連線；斷線時可先變 nil
	Connected bool
}

func (s *Session) BindConn(c Conn) {
	s.Conn = c
	s.Connected = true
}

func (s *Session) DropConn() {
	s.Conn = nil
	s.Connected = false
	// 注意：這裡「先不」清 RoomID——留給重連窗口決策
}
```

重點不在欄位多寡，而在**生命週期不同**：

- 連線沒了 → `DropConn`  
- Session 要不要銷毀 → 看你的重連／退房政策  
- 進房後 → Session 綁到 Room 的某個 seat（座位）

## 為什麼要分開：幾種真實狀況

### 1. 讀 loop 屬於連線；業務狀態屬於 Session／Room

典型形狀（你在 F／C 卷看過類似的）：

```text
每條連線一個讀 goroutine
    → 解出 envelope
    → 找到 Session
    → 把命令丟給 Room（入隊），而不是在讀 loop 裡算物理
斷線
    → 讀失敗或 ctx 取消
    → 退出讀 loop
    → unregister／標斷線
```

若你把「HP、座標、Inventory」只掛在連線物件的閉包區域變數裡，重連等於重生——因為舊閉包已經沒了。

### 2. 斷線策略要「寫成規則」，不要靠當下感覺

常見三檔（選一個寫進設計，別混用）：

| 策略 | 行為 | 適合 |
|------|------|------|
| 立刻移除 | 斷線馬上退房、座位空出 | 休閒房、人多容易補 |
| 短窗重連 | 標 `Disconnected`，N 秒內同身份可綁新連線 | 對戰中途抖網路 |
| 託管／AI | 斷線後 Server 代打或凍結操作 | 競技、長局 |

Arena Mini 教學版可以先做「斷線就從房間移除」——但你**要知道這是簡化**，不是唯一真理。

### 3. 同名／同帳號重複登入

兩個人都叫 `Ada`，或同一帳號開兩個分頁，你必須文件化政策：

- 拒絕第二個  
- 踢掉舊連線（頂號）  
- 允許觀戰分身但只有一個能操作  

沒寫清楚時，Bug 會以「為什麼我突然被踢／為什麼房間裡有兩個我」出現。

### 進階可先略過

- 正式一點會用登入 token（JWT 或 session id）換 `UserID`，名字只是顯示欄位。  
- 重連時帶 `roomID + seat + reconnectToken`，Server 驗證後把新 Conn 綁回舊 Session。  
- F7 心跳：先發現「連線假死」，再走斷線流程，比乾等 TCP 超時友善。

## 遊戲 Server 會用在哪

Arena Mini（M4）目前常用 **`name` query** 當簡化 Session——夠教學，不夠安全（誰都能冒名）。  
進階路徑：

1. 先分清 Conn vs Session（這章）  
2. 再加 token／user id  
3. 再加重連窗口  

只要模型對，後面加欄位是增量；模型不對，會整段重寫。

## 請丟掉的舊習慣

1. **只用 websocket 物件當唯一身份**——它一斷，你的「玩家」就沒了。  
2. **斷線一律立刻抹掉一切進度**——連考慮重連窗口的機會都不留。  
3. **在讀 loop 裡堆全部業務狀態**——生命週期綁死在連線上。  
4. **同名政策靠運氣**——偶爾覆蓋、偶爾重複，除錯到哭。

## 動手練習

### 必做

1. 寫下你的政策（兩三句即可）：**斷線 5 秒內重連**，局內座標／分數應保留、重算、還是凍結？為什麼？  
2. 列出 Session **至少 4 個欄位**，並標註哪個跟連線走、哪個跟房間走。  

### 選做

1. 設計「頂號」流程：舊連線收到什麼訊息？房間廣播什麼？新連線何時能送 input？  

## 常見坑

- **Session ID 用短暫連線 id**：重連後對不上，等於沒 Session。  
- **Room 只存 `map[*Conn]Player`**：鍵是連線指標，一重連鍵就變了。應以玩家 id／座位當鍵。  
- **清連線時誤清權威狀態**：座標是 Room／Rules 的真相，不是 Conn 的私有變數。  
- **忘記 unregister**：連線 goroutine 退出了，房間表還指著舊 Session → 廣播寫到已關閉連線、洩漏、奇怪錯誤。

下一章 H3：Room 狀態機——Lobby 能做的事，Playing 就不該能做。
