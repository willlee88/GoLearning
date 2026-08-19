---
lessonId: "C6"
title: "Race Lab：先紅後綠，親眼看見 data race"
description: "用 go test -race 讀報告、修好共享寫入——這是 C 卷檢查點。"
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

## 這章你會搞懂什麼

**Data race（資料競爭）是 bug**，不是「機率問題」「測一百次沒炸就好」。  
兩個 goroutine 同時碰同一變數，至少一個是寫，又沒有同步——這就是 race。

Go 提供 **race detector**：`go test -race`／`go run -race` 會插入偵測碼，告訴你衝突的讀寫堆疊。

本 lab 先給你一個會 race 的大廳計數（`broken`），再給正確版本（`safe`）。你的任務不是背 API，而是：**看懂報告 → 選修法 → 證明變綠**。

## Python 對照

Python 多執行緒對 dict／共用變數的問題，常被 GIL 跟運氣蓋住；asyncio 單執行緒又是另一個模型。  

到 Go 請預設一句話：**共享寫入 = 要同步，或根本不要共享。**

## 怎麼寫

```bash
cd examples/c06-race-lab
go test -race ./...
go test -race -run TestBroken ./broken/   # 預期失敗（偵測到 race）
go test -race ./safe/                     # 預期通過
```

先讀 `broken` 的程式：多個 goroutine 同時改同一份 map／計數，沒有鎖、也沒有單一擁有者。  
再讀 `safe`：要嘛 Mutex 保護，要嘛把寫入收斂到安全路徑。

## 細節

### 報告在說什麼

典型 race 報告會指出：

1. 前一次寫在哪裡  
2. 這一次讀／寫在哪裡  
3. goroutine 是在哪裡被 `go` 出來的  

別被牆原文嚇到：你只要抓住「這兩個堆疊在搶同一塊記憶體」。

### 三種常見修法

| 策略 | 什麼時候適合 |
|------|----------------|
| `Mutex` 保護 map | 資料結構共享、API 是同步方法 |
| 單一 owner goroutine + 請求 channel | 想要狀態機清晰、少鎖 |
| 分片／減少共享 | 進階：降鎖競爭 |

沒有永遠正確的一種。lab 的重點是你能**解釋為什麼這樣就沒 race**。

### 為什麼「多測幾次」不可靠？

Race 跟排程時機有關。今天沒炸，明天負載一變就炸。  
Detector 也不是數學證明萬無一失，但它抓實務中的問題很兇——CI 開著它，CP 值極高。

### 進階可先略過

- happens-before 與 Go 記憶體模型（`go.dev/ref/mem`）。  
- race detector 有成本：開發／CI 開；產線 binary 通常關。

## 遊戲 Server 會用在哪

大廳人數、房間列表、session 索引、連線註冊表——全是 race 熱點。  

建議：CI 對 `internal/...` 跑 `go test -race`。遊戲狀態一旦偶發錯亂，除錯成本遠高於「一開始就同步」。

## 請丟掉的舊習慣

1. 「測很多次沒炸＝沒 race」。  
2. 修 race 靠 `sleep`／亂加隨機延遲。  
3. 看到報告就隨手加一把大鎖包山包海，卻說不清不變量是什麼。

## 動手練習

### 必做（C 卷檢查點）

1. 閱讀 `broken`：指出誰在並發讀寫什麼。  
2. 確認 `safe` 通過 `-race`。  
3. 自己實作第三種修法（例如 owner goroutine）放在 `safe` 旁新檔，並用 `-race` 證明。  

### 選做

1. 對 `examples/a08-player-registry` 開多 goroutine 壓測 `Add`，跑 `-race`。  

## 常見坑

- **Windows 上 `-race` 需要可用的 CGO 工具鏈**：官方安裝通常 OK；若失敗先查環境。  
- **測試不夠並發**：只有單線跑，detector 沒事可抓——測試要真的並行。  
- **修了 A 變數、忘了 B 變數**：報告可能一次只秀一組；修完再跑直到乾淨。

## 延伸閱讀

- <https://go.dev/doc/articles/race_detector>  
- <https://go.dev/ref/mem>  

## 檢查點收束

走到這裡，你應具備：

- 錯誤當值（B）  
- goroutine／channel／select／context／mutex（C1–C5）  
- 用 `-race` 證明共享記憶體沒在裸奔（C6）  

C7–C10 是實務加分：errgroup、背壓、排程成本感、pipeline。主路徑也可以先往 D／F 走，但建議至少讀完 C8 的背壓直覺。  
