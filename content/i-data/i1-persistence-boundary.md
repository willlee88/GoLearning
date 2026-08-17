---
lessonId: "I1"
title: "持久化邊界"
description: "熱路徑 vs 冷路徑；何時 fsync 思維。"
volume: "i"
order: 1
level: "l2"
status: "ready"
path_required: true
tags: ["persistence", "game-server"]
example: ""
prev: "I0"
next: "I2"
---

## 本章你會建立的心智模型

寫入成本從低到高：記憶體 → 本機檔案 → Redis → DB → 跨區複製。權威對局狀態活在 **Room 記憶體**；持久化發生在**事件邊界**：

- 登入／登出  
- 購買、郵件  
- 對局結束結算  
- 斷線保留 session（短 TTL 可用 Redis）

## Python 對照

| Python | Go |
|--------|-----|
| 請求內 `session.commit()` | 同樣避免在 tick 裡 commit |
| Celery 非同步任務 | worker goroutine / queue |

## L1 能用

規則：

1. Tick 內：只改記憶體  
2. 命令成功且需持久：丟到有界 channel 給 writer  
3. 關房／關服：flush 未完成寫入  

## L2 機制

- **至少一次 vs 剛好一次** 結算語意。  
- 崩潰丟失「進行中局」是否可接受（多數 PvP 可接受）。  
- 帳號資產不可丟 → 獨立交易邊界。  

## 請丟掉的 Python 習慣

1. ORM 在每個小更新自動 flush。  
2. 把 Redis 當唯一真相且從不備份關鍵資產。  

## 遊戲 Server 連結

Arena Mini 教學可不接 DB；概念上 `Ended` 時可寫一筆 match result。

## 練習

### 必做

1. 列出 5 個遊戲狀態，標「記憶體／Redis／DB」。  
2. 說明為何玩家座標不該每 tick 寫 Postgres。  

## 延伸閱讀

- 規劃書卷 I  
