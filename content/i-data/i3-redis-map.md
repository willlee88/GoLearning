---
lessonId: "I3"
title: "Redis 用途地圖"
description: "presence、排行榜、pubsub、分散鎖——何時該上。"
volume: "i"
order: 3
level: "l2"
status: "ready"
path_required: true
tags: ["redis"]
example: ""
prev: "I2"
next: "I4"
---

## 本章你會建立的心智模型

Redis 常見於遊戲後端：

| 用途 | 結構 | 備註 |
|------|------|------|
| 在線 presence | SET / HASH + TTL | 心跳續期 |
| 排行榜 | ZSET | 分數排序 |
| 跨進程房間訊息 | Pub/Sub | 非可靠隊列 |
| 匹配佇列 | LIST / ZSET | 注意原子性 |
| 分散鎖 | SET NX EX | 小心死鎖 |

**單機 demo 可不裝 Redis**；多 Gateway 水平擴展時才痛。

## Python 對照

`redis-py` 同概念；Go 生態有 go-redis 等。

## L1 能用

決策樹：

1. 只在一個進程？→ 記憶體  
2. 要跨進程且可丟？→ pubsub  
3. 要排行／TTL？→ Redis  
4. 錢與庫存？→ DB 交易  

## L2 機制

- Pub/Sub **不保證**離線接收。  
- 鍵設計：`presence:{userId}`、`room:{id}:members`。  
- 熱 key 與大 key 會打爆單分片。  

## 請丟掉的 Python 習慣

1. 把 Redis 當永久 DB。  
2. 無 TTL 的 session 鍵堆到爆。  

## 遊戲 Server 連結

房間仍 sticky 到某 Game 進程；Redis 協助發現與旁路資料。

## 練習

### 必做

1. 設計 3 個 key 命名。  
2. 說明 pubsub 不能取代「可靠結算隊列」。  

## 延伸閱讀

- Redis 官方資料型別文件  
