---
lessonId: "H0"
title: "H 卷導讀 · 遊戲 Server 核心"
description: "連線會了之後，局怎麼跑？Session、Room、Tick、權威狀態——把網路層變成可玩的一局。"
volume: "h"
order: 0
level: "l1"
status: "ready"
path_required: true
tags: ["game-server"]
example: ""
prev: "G3"
next: "H1"
---

## 這章你會搞懂什麼

前面 F 卷教你「怎麼連」（TCP／WebSocket），G 卷教你「怎麼說」（信封、命令、事件）。  
連得上、訊息也解析得出來之後，還差一塊：**局到底怎麼跑**。

H 卷就是在補這塊。讀完你應該能用白話講清楚：

- 誰是玩家身份（Session）  
- 一局關在哪裡（Room）  
- 什麼時候能 Ready、什麼時候能移動（Phase／狀態機）  
- 為什麼座標不能信客戶端（權威 Server）  
- 為什麼不是「收到一包就立刻改世界」（Tick）  

這卷的終點是 Lab：把這些概念落到 `demo/arena-mini`，兩個人真的能開一局、移動、看到對方位置。

## 先建立心智模型

把一局遊戲想成「有節奏的模擬」，不要想成「一堆 HTTP 請求互改全域變數」。

| 概念 | 白話 | 英文常見寫法 |
|------|------|----------------|
| Session | 這個玩家在線的上下文（名字、座位、連線握把） | session |
| Lobby / Room | 隔離的一局空間；別人房裡發生的事不影響你 | room / lobby |
| Phase | 大廳／進行中／結束——合法操作跟著階段變 | phase / FSM |
| Tick | Server 固定頻率的心跳：套用輸入、推進模擬、廣播 | tick / fixed timestep |
| Command / Input | 客戶端送來的**意圖**（我想往右），不是真相 | command / input |
| State | 伺服器算出來的權威快照（你現在在哪、幾分） | authoritative state |

一句話串起來：

> 連線進來 → 綁成 Session → 進 Room → Lobby 等人 Ready → 進入 Playing → 每個 Tick 吃輸入、改 State → 廣播給房內所有人 → 達成分結束 Ended。

## 跟前面幾卷怎麼接

| 卷 | 解決的問題 | H 卷怎麼用到 |
|----|------------|--------------|
| F 網路 | socket、心跳、斷線 | Session 掛在連線上；斷線要清座位 |
| G 協定 | 信封 `type` + `payload` | `ready` / `input` / `state` / `error` |
| C 併發 | goroutine、channel、select | 每房一個 tick loop 很常見 |
| B 錯誤 | sentinel error | `ErrInvalidPhase`、壞輸入要可辨認 |

你不用把前面全背熟，但寫 Room 時若開始「在 handler 裡亂改 map」，請回頭看 C／H9：規則跟 I/O 要分開。

## 本卷地圖

| 章 | 主題 | 你讀完要能回答 |
|----|------|----------------|
| H1 | 後端職責地圖 | Gateway 跟 Rules 各自不該碰什麼 |
| H2 | Session 與連線 | 連線斷了，玩家身份還在不在 |
| H3 | Room 狀態機 | Playing 時為什麼不能 Join |
| H4 | 權威 Server | 為什麼不能讓客戶端直傳座標當真相 |
| H5 | Tick | 為什麼要固定頻率更新 |
| H6 | 輸入佇列與校驗 | 輸入何時入隊、何時套用 |
| H7 | 狀態同步 | 全量快照夠不夠、何時才要 delta |
| H8 | 匹配入門 | 自動配對跟「輸入房號」差在哪 |
| H9 | 規則與 I/O 分離 | 為什麼 `game` 不該 import websocket |
| H10 | Arena Mini Lab | 兩人真的跑完一局權威對戰 |

## 這卷會一直出現的 demo

本站的對照實作是 **Arena Mini**（`demo/arena-mini`）：

- 瀏覽器開兩分頁、不同名字、同一個房號  
- 雙方 Ready → Server 開局  
- 用方向鍵／按鈕送 **input**（意圖）  
- Server 算座標、碰撞得分，再廣播 **state**  

後面幾乎每章都會指回它。你可以先把 server 跑起來玩一把，再回來讀細節——有畫面會比較有感覺。

```powershell
cd F:\GoLearning\demo\arena-mini\server
go mod tidy
go run ./cmd/server
# 瀏覽器開 http://localhost:8080
```

## 請丟掉的舊習慣（整卷通用）

1. **把前端算完的結果當答案**（HP、座標、得分）。前端可以預測顯示，真相在 Server。  
2. **收到一包就立刻改世界並廣播**，沒有 tick、沒有相位。之後難重現、難測、難防作弊。  
3. **一個 God 物件**：又握著 `net.Conn`，又在算碰撞，又在寫資料庫。  
4. **布林旗標滿天飛**（`started` / `finished` / `closing`）互斥不清——用明確的 Phase。

## 檢查點（讀完整卷後）

你應該能親手驗證：

1. 跑通 `demo/arena-mini`：兩人 Ready → 開局 → 移動同步 → 狀態廣播。  
2. 用自己的話解釋：為什麼 input 是意圖、state 是真相。  
3. 指出 `internal/game`（規則）跟 `internal/hub`（連線／廣播）的邊界。  

下一章 H1：先畫後端職責地圖——腦裡分層，程式才分得開。
