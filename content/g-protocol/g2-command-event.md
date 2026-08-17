---
lessonId: "G2"
title: "命令與事件分離"
description: "client 下指令；server 廣播權威結果。"
volume: "g"
order: 2
level: "l2"
status: "ready"
path_required: true
tags: ["protocol", "game-server"]
example: "examples/g01-envelope"
prev: "G1"
next: "G3"
---

## 本章你會建立的心智模型

**命令（Command）**：客戶端意圖（我想往左、我點準備）。  
**事件／狀態（Event/State）**：伺服器認定發生了什麼（權威快照、某人加入）。

千萬不要讓客戶端直接廣播「我的 HP=999」。權威 Server 校驗命令 → 改狀態 → 再廣播。

## Python 對照

CQRS 輕量版：相同概念，Go 用型別與 room.Apply 表達。

## L1 能用

```text
Client → {type:"move", payload:{dx:-1}}
Server → validate → apply
Server → {type:"state", payload:{players:[...]}}  (廣播)
```

## L2 機制

- 命令應冪等或帶序號（重連重送）。  
- 狀態可全量或 delta（H 卷）。  
- 系統事件（join/leave）與遊戲狀態可分開 type。  

## 請丟掉的 Python 習慣

1. 客戶端算完結果只通知 Server「存檔」。  
2. 把聊天與移動混成無結構字串。  

## 遊戲 Server 連結

`room.Apply(cmd) (events, error)` 純邏輯可測（接 H 卷）。

## 練習

### 必做

1. 列出 Arena Mini 目前哪些是命令、哪些是事件。  
2. 設計 `ready` 命令與 `phase` 狀態欄位。  

## 延伸閱讀

- 權威伺服器概念文（任選）  
