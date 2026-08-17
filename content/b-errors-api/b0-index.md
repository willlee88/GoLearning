---
lessonId: "B0"
title: "B 卷導讀 · 錯誤與 API 設計"
description: "把 error 當值、設計可維護的邊界 API。"
volume: "b"
order: 0
level: "l1"
status: "ready"
path_required: true
tags: ["errors", "api"]
example: ""
prev: "A14"
next: "B1"
---

## 本章你會建立的心智模型

B 卷把 P0.4 的直覺做成工程能力：包裝、哨兵、`Is/As`、panic 邊界、以及「接受 interface、回傳具體型別」的 API 味道。遊戲 Server 的非法操作、滿房、斷線都應是**可預期的錯誤路徑**。

## 本卷地圖

| 章 | 主題 |
|----|------|
| B1 | error 基礎與哨兵 |
| B2 | wrap 與錯誤鏈 |
| B3 | panic / recover 邊界 |
| B4 | API 慣例與建構 |

## 練習

完成 B1–B4 與 `examples/b02-room-errors`。
