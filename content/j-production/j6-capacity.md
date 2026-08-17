---
lessonId: "J6"
title: "容量與廣播放大"
description: "粗估頻寬與 CPU；O(N²) 的殺傷力。"
volume: "j"
order: 6
level: "l2"
status: "ready"
path_required: true
tags: ["capacity"]
example: ""
prev: "J5"
next: "J7"
---

## 本章你會建立的心智模型

粗估：

```text
出站 ≈ 房間數 × 每房人數 × 每 tick 快照大小 × tick Hz
```

全量互相廣播在人數上升時很痛。興趣管理、delta、降頻觀戰是優化方向。

## Python 對照

同公式；語言不是第一瓶頸時常是架構。

## L1 能用

例：4 人房、200B state、20Hz → 每房出站約 `4 × 200 × 20 = 16KB/s` 量級（示意，含 JSON 開銷會更高）。

## L2 機制

- JSON 編解碼 CPU。  
- 鎖競爭：廣播時持鎖送網路是大忌（先拷貝再送）。  
- GOMAXPROCS 與多核。  

## 請丟掉的 Python 習慣

1. 「先上線再看」完全無容量圖。  
2. 複製 AAA 架構卻只有 10 人。  

## 遊戲 Server 連結

Arena Mini 全量 state；用 load-client 感受曲線。

## 練習

### 必做

1. 估算 50 人同房 20Hz 200B 的出站。  
2. 提出一個降流量手段。  

## 延伸閱讀

- H7 狀態同步  
