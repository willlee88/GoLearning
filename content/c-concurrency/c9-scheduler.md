---
lessonId: "C9"
title: "Scheduler 直觀：便宜，但不是免費"
description: "G／M／P 不用背；建立 goroutine、阻塞、GOMAXPROCS 的成本感就好。"
volume: "c"
order: 9
level: "l3"
status: "ready"
path_required: false
tags: ["scheduler", "runtime"]
example: ""
prev: "C8"
next: "C10"
---

## 這章你會搞懂什麼

Go runtime 把大量 **goroutine（G）** 多路複用到較少的作業系統執行緒（**M**），中間還有邏輯處理器 **P**（跟 `GOMAXPROCS` 有關）。  

你不必默寫源碼，但要有正確成本感：

- 建立 goroutine **便宜，不是免費**  
- 卡住在系統呼叫／鎖競爭，會影響吞吐  
- `GOMAXPROCS` 影響**平行**程度（真的多用幾個核心）  
- 超長 CPU 迴圈可能讓排程顯得「不公平」——通常應拆工作，而不是到處 `Gosched`

這章是 L3：幫你避免「因為便宜就開十萬個」的幻覺。

## Python 對照

| Python | Go |
|--------|-----|
| GIL 下 CPU 平行常受限 | 可真平行（也因此 race 更真實） |
| 執行緒很重，不敢亂開 | goroutine 輕，但濫開一樣死 |
| asyncio 單線事件迴圈心智 | 多 P 上跑許多 G；阻塞模型不同 |

## 怎麼寫

先用工具感受平行度，不必改業務碼：

```bash
go test -cpu=1,4 ./...
```

或跑自己的小基準，看不同 `GOMAXPROCS` 下耗時變化。  
真正要優化時，用 pprof（D4）看時間花在你的包還是 `runtime`／鎖上。

## 細節

### 你需要記住的直觀

1. **網路之類的阻塞** 很多會進 runtime 的 poller，不會傻傻佔死一個 OS thread 到絕望——但這不代表你可以無視生命週期。  
2. **鎖競爭、假共享** 會讓你「加核不加光速」。  
3. **每連線 1–2 個 goroutine** 通常可接受；**每個 tick 再炸一堆 goroutine** 要先有上限與測量。  

### 為什麼「開 10 萬個」聽起來酷、跑起來痛？

每個 G 都有棧與調度開銷；再加上 channel／鎖／記憶體配置，尖峰延遲會很難看。  
設計應從「誰擁有狀態、速率多少、背壓是什麼」出發，而不是從「語言允許多輕」出發。

### 進階可先略過

- 工作窃取（work stealing）、sysmon、hand-off 等調度細節。  
- `runtime.Gosched` 偶爾有用，但更常見的正解是把大CPU工作切開或丟有界 worker pool。

## 遊戲 Server 會用在哪

- 連線數 × 每連線 goroutine 數 ≈ 粗估調度壓力  
- 房間 tick 保持單線狀態機，通常比「每玩家每 tick 一個 G」穩  
- 發現延遲尖刺：先 `-race` 與 pprof，再懷疑 scheduler 魔法  

## 請丟掉的舊習慣

1. 因為便宜就無設計地開海量 goroutine。  
2. 持鎖做重 CPU。  
3. 不測量就調 `GOMAXPROCS`「應該會比較快」。

## 動手練習

### 必做

1. 讀一篇 Go scheduler 概述（官方 blog／設計文件皆可）。  
2. 估算：1 萬連線 × 2 goroutine，數量級你能不能接受？記憶體大概怎麼想？  

### 選做

1. 寫一個故意佔滿 CPU 的迴圈，觀察延遲；再改成分片處理比較體感。  

## 常見坑

- **把「邏輯 CPU」跟「goroutine 數」當同一件事**：不是。  
- **只看平均延遲、不看 P99**：遊戲更在意尖刺。  
- **過早微優化 runtime 參數**：先修算法、分配、鎖競爭。

## 延伸閱讀

- Go Blog / Scalable Go Scheduler Design Doc  
