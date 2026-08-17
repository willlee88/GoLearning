---
lessonId: "P0.8"
title: "如何使用本站"
description: "主路徑、索引、L1/L2/L3 與檢查點怎麼走。"
volume: "p0"
order: 8
level: "l1"
status: "ready"
path_required: true
tags: ["meta"]
example: ""
prev: "P0.7"
next: "A0"
---

## 本章你會建立的心智模型

GoLearning 是**廣覆蓋手冊 + 可循序主路徑**。你可以像讀旅行指南一樣順走 P0→K，也可以用搜尋／卷冊當 API 手冊。深度用 **L1／L2／L3** 分層，避免「百科變水文」或「一上來就被 runtime 淹死」。

## Python 對照

把本站想成：

- **主路徑** ≈ 一門從 Python 轉 Go 的課程大綱  
- **卷冊索引** ≈ 文件站（像你查 Python docs）  
- **examples/** ≈ 可執行的 notebook／腳本倉庫  

## L1 能用 — 建議走法

1. **順完 P0**（你在這裡）：建立心智，不要急著抄語法。  
2. 進入 **A 卷** 打語言地基；遇併發名詞先記連結，細節在 C 卷。  
3. **C 卷** 認真做 race lab——這是「深刻」的分水嶺。  
4. **F→H** 接到遊戲 Server；邊做邊回鏈結語言章。  
5. **K Capstone** 收束。  

標記：

- **必讀**：主路徑 `path_required: true`  
- **選讀／深潛**：L3 或 `path_required: false` 的廣覆蓋章  

## L2 機制 — 內容標記

| 標記 | 意義 |
|------|------|
| L1 | 能用、最小例子 |
| L2 | 機制與契約（預設深度） |
| L3 | 深潛、可摺疊 |
| status: ready | 可學 |
| status: stub / draft | 佔位，之後補 |

章節 frontmatter 的 `example` 指向 repo 內路徑，請本機 `go test`／`go run`。

## L3 深潛（可選）

貢獻內容時使用 `content/_templates/lesson.md`，並跑 `scripts` 裡的連結檢查（陸續補齊）。

## 請丟掉的 Python 習慣

1. 只收藏文章不跑例子。  
2. 跳過錯誤處理章直接貼 Server 範例。  
3. 用「能跑」代替「理解 race 與生命週期」。  

## 遊戲 Server 連結

從現在起，每章尽量有「遊戲 Server 連結」框——即使是字串章，也會指出聊天過濾、協定欄位等落點。這是本站與一般 Go 教學的差別。

## 練習

### 必做

1. 在網站標完 P0 進度（若進度功能已啟用）。  
2. 執行 `examples/p0-config-stats`，用自己的話寫 5 條 Python vs Go 差異（檢查點 P0）。  

### 選做

1. 瀏覽 `docs/規劃書.md` 的卷冊地圖，標出你最想先深讀的 3 個主題。  

## 常見坑與如何看見

- 內容標 `stub` 仍硬讀：看 frontmatter `status`。  
- 範例失敗：先確認 Go ≥ 1.22 與在正確 module 目錄。  

## 延伸閱讀

- `docs/規劃書.md`  
- `README.md`  
