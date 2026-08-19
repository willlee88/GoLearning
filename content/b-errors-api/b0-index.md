---
lessonId: "B0"
title: "B 卷導讀：錯誤是值，API 要好維護"
description: "把 error 當值來設計；學會包裝、哨兵、panic 邊界，以及清楚的建構慣例。"
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

## 這章你會搞懂什麼

P0.4 跟你講過：**Go 不靠例外當日常控制流**，錯誤是回傳值。  
B 卷要把那個直覺變成工程肌肉：怎麼定義穩定條件、怎麼加上下文、什麼時候才能 panic、公開 API 長什麼樣比較好測。

為什麼這卷排在語言基礎後面、併發前面？因為遊戲 Server 裡「滿房、階段不對、票過期、斷線」都是**每天會發生的路**，不是災難。你若還用 Python 的 `raise` 心智亂甩，後面一上 goroutine 會更難收場。

讀完 B 卷，你要能把「房間加入失敗」做成：人看得懂訊息、程式用 `errors.Is` 分支、客戶端拿到穩定錯誤碼。

## 本卷地圖

| 章 | 主題 | 你會練到什麼 |
|----|------|----------------|
| B1 | error 基礎與哨兵 | `error` 介面、`ErrXxx`、`errors.Is` |
| B2 | wrap 與錯誤鏈 | `%w`、上下文、還能解開根因 |
| B3 | panic / recover 邊界 | 什麼能炸、在哪接住、接住後幹嘛 |
| B4 | API 慣例與建構 | 小介面注入、`NewX`、可測邊界 |

## 遊戲 Server 會用在哪

非法操作、滿房、斷線、逾時——全部應是**可預期的錯誤路徑**。  
客戶端看到的是協定錯誤碼；日誌裡看到的是「哪個房間、哪個玩家、哪一步失敗」。兩層資訊靠哨兵 + wrap 分工，不是靠 stack trace 碰運氣。

## 動手練習

完成 B1–B4，並把 `examples/b02-room-errors` 跑通、讀懂測試怎麼用 `errors.Is`。

## 常見坑

- **把 B 卷當「語法糖」略過**：後面網路卷一出現「壞封包」，你會不知道要回 error 還是 panic。  
- **只記 `if err != nil`**：沒哨兵、沒 wrap，日誌會變成一串沒上下文的字。  
