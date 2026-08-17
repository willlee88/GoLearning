---
lessonId: "J7"
title: "水平擴展提醒"
description: "房間親和、為何不能隨便 round-robin WS。"
volume: "j"
order: 7
level: "l2"
status: "ready"
path_required: true
tags: ["scale"]
example: ""
prev: "J6"
next: "J8"
---

## 本章你會建立的心智模型

有狀態的 Room **不能**在無共享記憶體下隨便把同一房的封包打到不同機器。策略：

1. **Sticky**：同一 room 到同一 Game 進程  
2. **遷移**：極少用、複雜  
3. **Redis 同步狀態**：延遲與複雜度上升  

Gateway 無狀態 + Game 有狀態 是常見切分。

## Python 對照

Django 多 worker 的 session sticky 問題同源。

## L1 能用

教學結論：先垂直優化單機 Room；擴展時按 room id 一致性雜湊。

## L2 機制

- 匹配服務與對局服務分離。  
- 全局廣播用 pubsub 需接受最終一致。  
- 觀測：每實例 metrics + 彙總。  

## 請丟掉的 Python 習慣

1. 無狀態 HTTP 思維硬套即時對局。  
2. 多副本搶寫同一記憶體地圖。  

## 遊戲 Server 連結

Arena Mini 單進程；文件標明擴展邊界。

## 練習

### 必做

1. 畫「兩 Gateway + 兩 Game」房間路由草圖。  

## 延伸閱讀

- 規劃書 E1  
