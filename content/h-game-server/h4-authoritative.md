---
lessonId: "H4"
title: "權威 Server"
description: "客戶端送意圖；伺服器擁有真相。"
volume: "h"
order: 4
level: "l2"
status: "ready"
path_required: true
tags: ["authoritative", "game-server"]
example: "examples/h06-apply-input"
prev: "H3"
next: "H5"
---

## 本章你會建立的心智模型

**Authoritative server**：遊戲狀態以伺服器為準。客戶端可以預測顯示，但提交的是 **input/command**，不是「我的新 HP=999」。Server 校驗 → 套用 → 廣播狀態（或 delta）。

## Python 對照

不要做成「客戶端 POST 完整存檔」；像權威模擬一步。

## L1 能用

```text
壞：{ type:"state", payload:{ x:999, y:999 } }  // 客戶端自稱
好：{ type:"input", payload:{ dx:1, dy:0 } }   // 意圖
     Server: x = clamp(x+dx*speed)
```

## L2 機制

- 防作弊基線：邊界、速度、冷卻、視野。  
- 仍可被 reverse engineer——深度反作弊另層（J 卷）。  
- 單機 demo 也應走權威，養成肌肉。  

## 請丟掉的 Python 習慣

1. 信任前端算完的結果。  
2. 把同步當成「複製變數到所有人」。  

## 遊戲 Server 連結

M3 的 `pos` 點擊座標是**教學捷徑**；M4 改成 **input 方向**，座標只在 Server。

## 練習

### 必做

1. 說明為何 `pos` 直傳在正式遊戲危險。  
2. 設計 `input` payload 三個欄位。  

## 延伸閱讀

- G2 命令與事件  
