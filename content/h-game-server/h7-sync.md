---
lessonId: "H7"
title: "狀態同步策略"
description: "全量快照 vs delta；廣播放大。"
volume: "h"
order: 7
level: "l2"
status: "ready"
path_required: true
tags: ["sync"]
example: ""
prev: "H6"
next: "H8"
---

## 本章你會建立的心智模型

| 策略 | 優點 | 代價 |
|------|------|------|
| 全量 State | 簡單、易重連 | 頻寬大 |
| Delta | 省流量 | 複雜、要基準快照 |
| 興趣管理 | 大世界必要 | 實作重 |

M4 用**全量 JSON 快照**足夠；先正確再優化。

## Python 對照

同概念；注意序列化成本。

## L1 能用

```json
{ "v":1, "type":"state", "payload":"{\"phase\":\"playing\",\"tick\":12,\"players\":[...]}" }
```

或 payload 直接嵌物件（信封進化）。

## L2 機制

- 廣播 O(玩家 × 觀測者)。  
- 快照必須自洽（同一 tick 產出）。  
- 客戶端顯示可插值，但不改權威。  

## 請丟掉的 Python 習慣

1. 每人單獨不同步的狀態複本長期分叉。  
2. 過早上 delta 壓縮。  

## 遊戲 Server 連結

Arena Mini 每 tick 全量 `players[{name,x,y}]`。

## 練習

### 必做

1. 估算 4 人、20Hz、每快照 200B 的出站流量。  
2. 寫下何時需要 delta。  

## 延伸閱讀

- 規劃書 D7/H 卷  
