---
lessonId: "J0"
title: "J 卷導讀：能玩還不夠，要能關、能看、能守"
description: "上線前補齊觀測、優雅關閉、壓測、安全與擴展邊界——讓 Arena Mini 從 demo 變成可操作的服務。"
volume: "j"
order: 0
level: "l1"
status: "ready"
path_required: true
tags: ["production"]
example: ""
prev: "I5"
next: "J1"
---

## 這章你會搞懂什麼

前面幾卷你已經能把房間轉起來、tick 跑起來、狀態同步出去。  
但「本機能玩」跟「敢放給別人連」中間，還差一截很實際的工程：

- **出事時你看不看得懂？**（日誌、metrics）  
- **要關機時會不會把玩家狀態砍一半？**（優雅關閉）  
- **人一多會先爆在哪？**（壓測、容量）  
- **壞人／壞客戶端會不會把你打掛？**（安全基線）  
- **加機器能不能救？還是架構先卡死？**（水平擴展）

J 卷就是補這一段。讀完你要有一句話當標準：

> **能玩 ≠ 能上線。能關、能看、能壓、能守，才算接近生產。**

## Python 對照

| 你在 Python 常見做法 | Go／本卷對應 |
|----------------------|--------------|
| `logging` + 偶爾 print | `log/slog` 結構化欄位 |
| 自己拼 JSON 狀態頁 | `/metrics`（教學先 JSON，之後可換 Prometheus） |
| uvicorn／gunicorn 的 graceful | `signal` + `http.Server.Shutdown` + 停 tick |
| locust／自寫 asyncio client | `examples/j04-load-client` |
| Django middleware／Pydantic 驗證 | 輸入驗證、速率限制、信任邊界 |

概念大多「你聽過」，差別在：**遊戲長連線 + 有狀態房間**，關服與擴展比純 HTTP API 更挑剔。

## 本卷地圖

| 章 | 主題 | 為什麼要讀 |
|----|------|------------|
| J1 | 結構化日誌 `slog` | 出事靠欄位搜，不是靠人肉 parse |
| J2 | Metrics 與 `/metrics` | 回答「現在系統怎樣」 |
| J3 | 優雅關閉 | Ctrl+C／SIGTERM 時關得乾淨 |
| J4 | 壓測方法 | 有假設、漸進加壓、看誰先爆 |
| J5 | 安全基線 | 客戶端不可信 |
| J6 | 容量與廣播放大 | 粗估頻寬／CPU，怕 O(N²) |
| J7 | 水平擴展提醒 | 有狀態 Room 不能亂 round-robin |
| J8 | Lab | 把 metrics、shutdown、load-client 串起來 |

## 這卷跟 Arena Mini 的關係

Arena Mini（`demo/arena-mini`）在 M5 已經接上：

- JSON 風格的 `slog`  
- `GET /metrics`、`GET /healthz`  
- Ctrl+C 時先停房間、再 `Shutdown` HTTP  

J 卷不是叫你重寫整套，而是讓你**看懂為什麼要這樣、缺了會怎樣**，再用 load-client 自己加壓驗證。

## 遊戲 Server 會用在哪

想像上線當天：

1. 玩家說「卡死了」→ 你打開 metrics 看連線數／錯誤數有沒有異常  
2. 要滾動更新 → 發 SIGTERM，期望舊連線收乾淨、新流量切走  
3. 活動前 → 用假 client 把連線數拉上去，先找瓶頸  

這些都是 J 卷的日常劇本。

## 請丟掉的舊習慣

1. 「能連上就先上，觀測以後再說。」  
2. 關服直接殺行程，當日常操作。  
3. 一次開滿假人，只看有沒有 crash。  

## 動手練習

讀完本卷地圖後，先到 J1。你不需要一次背完所有術語；每章都會用白話再講一次。

## 常見坑

- **把 demo 的「能跑」直接當成容量保證**：沒壓過，數字都是猜的。  
- **只打日誌、不打 metrics**：高頻 tick 用 Info 日誌會把自己淹死。  
- **以為多開幾個進程就能水平擴展**：同一房的狀態若沒親和／同步策略，會更慘（J7）。

## 檢查點（本卷結束時）

Arena Mini：

1. `Ctrl+C` 能優雅退出（看得到關閉日誌順序）  
2. `GET /metrics` 有合理數字  
3. 用 `examples/j04-load-client` 打一輪假人  

細節在 J8 lab 一次做完。
