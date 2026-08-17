---
lessonId: "G0"
title: "G 卷導讀 · 序列化與協定"
description: "訊息形狀、版本化、命令與事件分離。"
volume: "g"
order: 0
level: "l1"
status: "ready"
path_required: true
tags: ["protocol"]
example: ""
prev: "F8"
next: "G1"
---

## 本章你會建立的心智模型

網路只是水管；**協定**決定能不能演進、除錯、防呆。G 卷建立：

1. JSON 信封的成本與邊界  
2. 命令（client→server）vs 事件／狀態（server→client）  
3. 版本欄位與相容  
4. Protobuf 等二進位方向（概念 + 小例）

## 章節

| 章 | 主題 |
|----|------|
| G1 | JSON 信封實務 |
| G2 | 命令與事件分離 |
| G3 | 版本化與相容策略 |

## 練習

完成 G1–G3，並對照 Arena Mini 的 message 型別。
