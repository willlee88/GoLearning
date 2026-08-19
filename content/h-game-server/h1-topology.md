---
lessonId: "H1"
title: "遊戲後端職責地圖：誰管連線、誰管規則"
description: "Gateway、Room、Rules、Persist 各做什麼、不該碰什麼；單進程也能先把腦裡分層畫清楚。"
volume: "h"
order: 1
level: "l2"
status: "ready"
path_required: true
tags: ["architecture", "game-server"]
example: ""
prev: "H0"
next: "H2"
---

## 這章你會搞懂什麼

大型遊戲後端常拆成很多服務：閘道、匹配、房間、結算、帳號……  
你現在做教學 demo，**不必一開始就微服務**；但若腦裡沒分層，程式會很快變成「一個檔案又讀 WebSocket 又判勝負又寫資料庫」。

這章目標：

- 能畫出一張**職責地圖**（誰做什麼、刻意不做什麼）  
- 知道 Arena Mini 的資料夾大概對應哪一層  
- 之後水平擴展時，知道「房間狀態該住在哪」——因為你早就劃清邊界了

## Python 對照

想像你用 FastAPI 寫即時遊戲（其實不太舒服，但比喻夠用）：

| 若混在一起 | 較乾淨的切法 |
|------------|--------------|
| 一個 route 裡讀 WS、改 `global players`、算碰撞、`INSERT` 分數 | 路由／WS 適配器只負責收發；領域函式純邏輯；背景 worker 跑 tick；DB 另包一層 |
| Django God Model 又管 HTTP 又管商業規則 | package 邊界：`hub` 不懂規則，`game` 不握 socket |

Go 沒有強制「框架幫你分層」——**邊界靠你用 package 與函式簽名表達**。這是好事：你想清楚了，結構就清楚。

## 怎麼想（先畫圖，再寫碼）

最小、夠用的心智圖：

```text
Client(s)
    │  WS / TCP
    ▼
Gateway ── 握手、讀寫、心跳、解信封
    │
    ▼
Session ──「這個玩家是誰、在哪一房」
    │
    ├──► Match / Lobby ── 湊人、進房
    │
    ▼
Room Runtime ── 生命週期 + Tick 迴圈 + 廣播節奏
    │
    ├──► Rules（純邏輯）── 合法移動、碰撞、勝負
    │
    └──► Persist（可選）── 結算、戰績、帳號
         （注意：通常不是每 tick 寫庫）
```

口訣：

- **Gateway 不懂規則**（不知道「撞到算不算得分」）  
- **Rules 不懂 socket**（函式參數裡不該出現 `*websocket.Conn`）  
- **Persist 不懂幀率**（不要在 20Hz 的熱路徑狂寫 DB）

## Arena Mini 怎麼對上這張圖

單進程沒關係，用資料夾表達分層就好：

| 職責 | Arena Mini 大概在哪 | 做什麼 |
|------|---------------------|--------|
| 組裝／啟動 | `demo/arena-mini/server/cmd/server` | `main`：聽 port、接上 hub |
| Gateway + 房間表 | `internal/hub` | 連線進出、轉訊息、廣播 |
| Rules + 房間狀態 | `internal/game` | Ready／Input／Tick／Snapshot |
| Persist | （M4 可幾乎沒有） | 之後 I 卷再補 |

你現在打開這些目錄，問自己一句：「這個 `.go` 檔如果拿掉網路，測試還能不能跑？」——答「能」的，多半是規則層；答「不能」的，多半是 I/O 層。

## 各層「做／不做」對照表

| 職責 | 做 | 不做 |
|------|----|------|
| Gateway | 握手、讀寫、心跳、基本解碼 | 判定勝負、算物理 |
| Session | 玩家身份、房號、連線握把 | 自己當權威世界 |
| Match / Lobby | 湊人、建房、拉人進房 | 每幀模擬 |
| Room Runtime | Phase、Tick、對房內廣播 | DB schema 設計 |
| Rules | 合法輸入、邊界、碰撞、得分 | 持有 `net.Conn` |
| Persist | 結算落庫、帳號、戰績查詢 | 每個 tick 寫一筆 |

「不做」跟「做」一樣重要。邊界被踩一次，之後測試與擴展都會痛。

## 為什麼要先在單機分層

之後若要水平擴展，常見痛點是**房間親和**（room affinity）：同一局的狀態必須落在同一台（或同一個 shard），因為 Tick 與權威狀態不能拆給互不相干的機器各改各的。

若你單機時期就已經：

- Room 有清楚 ID 與生命週期  
- 規則狀態只活在 Room／game 裡  
- Gateway 只是入口  

那之後把「整房」遷到某台 worker，只是部署問題；若狀態散落在一堆全域 map 與閉包，遷移會像考古。

### 進階可先略過

- 真的拆服務時，Gateway 與 Room Worker 之間常再走一層內部 RPC／訊息佇列。  
- 匹配服務可以偏無狀態；房間服務偏有狀態——H8 會再提。  
- 不是每個專案都要六角形架構全套名詞；先做到「規則可單測、I/O 可替換」就很賺。

## 遊戲 Server 會用在哪

幾乎所有即時玩法都會碰到這張地圖：

- 手遊大廳 → 匹配 → 進對戰房  
- 聊天室／觀戰：Gateway 仍在，Rules 可能極薄  
- 回合制卡牌：Tick 可能變「回合事件」，但分層一樣成立  

Arena Mini 是縮小版：手動輸入房號 ≈ 簡化匹配；hub ≈ Gateway + 房間索引；game ≈ Rules + Room 狀態。

## 請丟掉的舊習慣

1. **God class／God package**：一個型別既管連線又算物理又存檔。  
2. **請求執行緒／讀 loop 直接改全域 dict**，沒有房間隔離——兩個房的人會互相踩狀態。  
3. **「先能跑再說，分層以後重構」**——即時遊戲的 entangle 速度比 CRUD API 快得多，之後往往是重寫不是重構。  
4. **每接到一個封包就 `INSERT` 資料庫**——熱路徑會被磁碟／網路拖死。

## 動手練習

### 必做

1. 用紙或註解畫一張「Arena Mini 三層對應表」：Gateway／Room／Rules 各對到哪個目錄或型別。  
2. 打開 `internal/hub` 與 `internal/game`，指出**一處**「規則還碰到 I/O」或「已經分得很乾」的地方（寫下檔名與函式名）。  

### 選做

1. 假設要加「觀戰者只能收 state、不能送 input」——這規則該寫在 hub 還是 game？寫下你的理由。  

## 常見坑

- **用資料夾名稱假裝分層**：`game` 目錄裡卻 `import` 了 websocket——名稱很美，依賴已穿管。  
- **過早微服務**：兩個人的 demo 拆五個 repo，除錯成本會爆炸；先單進程、腦裡分層。  
- **Persist 黏在 Tick 上**：局內 20Hz 寫庫，人一多資料庫先掛，遊戲還在怪「網路不穩」。  
- **Room 邊界模糊**：所有人丟進同一個 `map[string]*Player` 當「全服世界」，之後想做多房／分流會很痛。

下一章 H2：Session——連線斷了，玩家還算不算「在線」？
