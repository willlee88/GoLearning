---
lessonId: "C6"
title: "race detector 實戰 Lab"
description: "製造 race、讀報告、修好它——C 卷檢查點。"
volume: "c"
order: 6
level: "l2"
status: "ready"
path_required: true
tags: ["race", "lab", "game-server"]
example: "examples/c06-race-lab"
prev: "C5"
next: "C7"
---

## 本章你會建立的心智模型

**Data race 是 bug**，不是「機率問題」。`go test -race` / `go run -race` 插入偵測，報告衝突的讀寫堆疊。本 lab 先給你一個會 race 的大廳計數，再給正確版本（鎖或 channel 擁有者）。

## Python 對照

Python 多執行緒對 dict 的問題常被 GIL 與運氣掩蓋；在 Go 請預設：**共享寫入 = 要同步**。

## L1 能用

```bash
cd examples/c06-race-lab
go test -race ./...
go test -race -run TestBroken ./broken/   # 預期失敗
go test -race ./safe/
```

## L2 機制

報告會指出：

- 前一次寫  
- 這一次讀/寫  
- goroutine 創建位置  

修復策略：

1. Mutex 保護 map  
2. 單一 goroutine 擁有 map + 請求 channel  
3. 分片減少鎖競爭（進階）

## L3 深潛（可選）

- happens-before 與記憶體模型。  
- race detector 的成本（開發／CI 開，產線 binary 通常關）。

## 請丟掉的 Python 習慣

1. 「測很多次沒炸就是沒 race」。  
2. 修 race 靠 sleep。  

## 遊戲 Server 連結

大廳人數、房間列表、session 索引都是 race 熱點。CI 應對 `internal/...` 開 `-race`。

## 練習

### 必做（C 卷檢查點）

1. 閱讀 `broken` 為何 race。  
2. 確認 `safe` 通過 `-race`。  
3. 自己實作第三種修法（例如 owner goroutine）於 `safe` 旁新檔。  

### 選做

1. 對 `a08-player-registry` 跑 `-race` 並開多 goroutine 壓測 Add。  

## 常見坑與如何看見

- Windows 需支援 cgo 的 toolchain 才能 -race（官方安裝通常可）。  
- 測試要**真的並發**才會觸發。  

## 延伸閱讀

- <https://go.dev/doc/articles/race_detector>  
- <https://go.dev/ref/mem>  

## M2 里程碑收束

你已具備：

- A 卷核心語言模型（至 interface／generics）  
- B 卷錯誤與 API  
- C 卷併發核心 + race lab  

下一階段（M3）：網路卷 F 與 WebSocket 房間範例。  
