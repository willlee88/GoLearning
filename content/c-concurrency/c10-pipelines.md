---
lessonId: "C10"
title: "Pipeline：何時該拆 stage，何時別炫技"
description: "用 channel 串起產生→處理→輸出；搞清楚誰 close、怎麼取消，以及線上房間為何通常不用它。"
volume: "c"
order: 10
level: "l2"
status: "ready"
path_required: false
tags: ["pipeline"]
example: "examples/c10-pipeline"
prev: "C9"
next: "D0"
---

## 這章你會搞懂什麼

**Pipeline（管線）**：資料流過好幾個 stage，每個 stage 常是「一組 goroutine + 進／出 channel」。典型形狀：

```text
gen → parse → filter → sink
```

它很適合串流處理：日誌、回放轉檔、批次轉換。  

它**不適合**把本來 20 行就寫得完的函式，硬拆成五個 channel「看起來很併發」。複雜度上升、close／取消一漏就洩漏，收益卻可能是負的。

## Python 對照

| Python | Go |
|--------|-----|
| generator 串起來 | channel 連 stage |
| `multiprocessing.Pool` 管線 | fan-out／fan-in + 有界佇列 |
| 迭代器懶計算 | 真實並發；更要管生命週期 |

Python 的 generator 多半還是單線拉取；Go pipeline 常真的並發跑，所以 C1–C8 的紀律一件都不能少。

## 怎麼寫

每一段典型長這樣：

```go
func stage(in <-chan int, out chan<- int) {
	defer close(out)
	for v := range in {
		out <- v * 2
	}
}
```

精神：

1. 從 `in` 讀到關閉  
2. 處理後寫入 `out`  
3. **結束時由這段的擁有者 close `out`**（通常是唯一傳送者）  

扇出（fan-out）多個 worker 時，要有人 **fan-in** 合併，並想清楚誰負責最終 close。

範例：`examples/c10-pipeline`。

## 細節

### 誰 close？

規則仍是 C2：傳送側擁有 close。  
多個 worker 寫同一個 out 時，不要大家搶 close——用 `WaitGroup` 等 worker 結束，再由協調者 close，或每人寫自己的 channel 再 fan-in。

### 取消怎麼做？

整條鏈傳同一個 `ctx`。每個 stage 的 loop 用 `select` 聽 `ctx.Done()`，避免上游停了下游還堵在傳送上。  
取消＋close 規則沒想清楚，是 pipeline 洩漏的第一名原因。

### 為什麼線上 Room tick 通常不用 pipeline？

房間狀態需要**嚴格順序與單一擁有者**：input → 套用 → 廣播。  
拆成多 stage 並發改同一房間，瞬間回到 race／鎖／順序噩夢。  
離線分析回放可以 pipeline；線上狀態機用單線迴圈通常更清楚。

### 進階可先略過

- 官方 blog「Pipelines」裡的 md5／bounded 平行範例。  
- 錯誤傳遞：errgroup 或「錯誤 channel」——選一種風格貫穿全鏈。

## 遊戲 Server 會用在哪

| 場景 | 建議 |
|------|------|
| 回放檔轉統計、日誌清洗 | pipeline 很合適 |
| 線上 Room tick／權威狀態 | 單線狀態機 + inbox |
| 廣播編碼 | 可小範圍平行，但要有界、要測量 |

## 請丟掉的舊習慣

1. 為炫技把簡單邏輯拆一堆 channel。  
2. fan-out 了卻忘了 fan-in／close，留下永久阻塞的 G。  
3. 線上熱路徑無界平行「每個事件一個 stage」。

## 動手練習

### 必做

1. 跑 `examples/c10-pipeline`。  
2. 加一個 stage 統計通過的數量（結束時印出）。  

### 選做

1. 為整條鏈加上 `ctx` 取消；主程式兩秒後 cancel，確認 goroutine 會結束。  

## 常見坑

- **中間 stage panic／return 卻沒 close out**：下游永等。  
- **多個 sender 重複 close**：panic。  
- **用 pipeline「順便」共享可變狀態**：那又變回鎖問題——別混。

## 延伸閱讀

- <https://go.dev/blog/pipelines>  

## C 卷你帶走的三句話

1. 誰擁有資料，誰就能決定鎖還是 channel。  
2. 開得起就要關得了；取消要傳下去。  
3. `-race` 不是可選——共享寫入請證明你沒在裸奔。  
