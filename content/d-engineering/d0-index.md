---
lessonId: "D0"
title: "D 卷導讀：能跑之後，要能維護"
description: "把程式變成可依賴、可測、可進 CI 的工程：modules、表驅動測試、fuzz、bench／pprof。"
volume: "d"
order: 0
level: "l1"
status: "ready"
path_required: true
tags: ["engineering"]
example: ""
prev: "C10"
next: "D1"
---

## 這章你會搞懂什麼

C 卷讓你寫得出併發；D 卷問的是另一個問題：

> 這東西換一台機器、換一個人、放進 CI，還能不能穩？

工程化不是「多裝套件」。對 Go 新手來說，主路徑很務實：

- 依賴版本說得清楚（`go.mod`／`go.sum`）  
- 測試好加案例（表驅動）  
- 會用 fuzz 砸壞輸入  
- 先測量再優化（bench／pprof）  

讀完 D0–D4，你應能把 `internal/game` 的規則放進 CI，並對 `-race`、測試、粗略效能有固定節奏。

## 本卷地圖

| 章 | 主題 | 你會得到什麼 |
|----|------|----------------|
| D1 | go.mod 與依賴心智 | 版本可重現，少「我機器可以」 |
| D2 | 表驅動測試 | 加案例像填表，不是複製貼上 |
| D3 | fuzz 入門 | 專治「我以為不會發生」的輸入 |
| D4 | bench 與 pprof | 感覺讓位給數據 |

## 遊戲 Server 會用在哪

- `Ready`／`PushInput`／相位轉移：表驅動測爆  
- 封包／字串解析：fuzz  
- 狀態快照、廣播：bench 找分配  
- CI：`go test ./...` 與關鍵包 `-race`  

Arena Mini 的 `demo/arena-mini/server` 就是這些習慣的落點。

## 動手練習

依序完成 D1–D4；每章的範例目錄都很小，請真的跑，不要只讀。

## 常見坑

- **覺得「小專案不必測試」**：遊戲狀態機一旦回歸，沒有表驅動會很痛。  
- **先微優化再測**：D4 會糾正你，但最好從現在就養成「先對、再快」。  
