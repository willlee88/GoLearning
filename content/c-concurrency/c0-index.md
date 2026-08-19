---
lessonId: "C0"
title: "C 卷導讀：併發——遊戲 Server 的中樞神經"
description: "搞懂誰擁有資料、怎麼停、怎麼證明沒 race。goroutine／channel／sync 都為這三件事服務。"
volume: "c"
order: 0
level: "l1"
status: "ready"
path_required: true
tags: ["concurrency"]
example: ""
prev: "B4"
next: "C1"
---

## 這章你會搞懂什麼

C 卷是「認真搞懂 Go」的分水嶺。

併發（concurrency，一次結構好多件進行中的事）不是「多開幾個 thread 就好」。你要能推理這三句話：

1. **誰擁有這份資料？**  
2. **誰在什麼時候讀／寫？**  
3. **怎麼停止、怎麼會合、怎麼證明沒有 data race？**

遊戲 Server 的連線 loop、房間 tick、廣播扇出，全部建立在這上面。寫錯的代價不是單元測試紅一下，而是線上偶發錯亂、關不乾淨、記憶體洩漏。

## 本卷地圖

| 章 | 主題 | 為什麼要讀 |
|----|------|------------|
| C1 | goroutine 生命週期 | 開了就要能停、能等 |
| C2 | channel 語意 | 通訊與擁有者規則 |
| C3 | select 與超時 | 房間迴圈的多路等待 |
| C4 | context 取消樹 | 關服／逾時往下傳 |
| C5 | Mutex 與共享狀態 | 沒有 GIL，鎖是真的 |
| C6 | race detector lab | **本卷檢查點**：先紅後綠 |
| C7 | errgroup | 一組任務的錯誤與取消 |
| C8 | 背壓 | 佇列有界，不然 OOM |
| C9 | Scheduler 直觀 | 成本感，不必背源碼 |
| C10 | Pipeline | 何時該拆 stage、何時別炫技 |

## 你手上的關鍵工具

```bash
go test -race ./...
```

這行會反覆出現。把它當成安全帶，不是可選裝飾。

## 遊戲 Server 會用在哪

- 每條連線：讀 loop goroutine  
- 每個房間：tick + 收命令的 select 迴圈  
- 大廳／房間表：Mutex 或單一 owner  
- 關服：root `context` cancel 往下剪枝  

## 動手練習

把 C1–C6 當主路徑；C6 的 `examples/c06-race-lab` 必須「看得懂為什麼紅、為什麼綠」。C7–C10 幫你補齊實務模式。

## 常見坑

- **還沒想清楚擁有者就開 goroutine**：後面只能加 sleep「碰運氣修好」。  
- **略過 `-race`**：在 Go 裡 data race 是 bug，不是「機率問題」。  
