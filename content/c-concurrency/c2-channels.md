---
lessonId: "C2"
title: "channel 語意"
description: "無緩衝／有緩衝、close、擁有者規則。"
volume: "c"
order: 2
level: "l2"
status: "ready"
path_required: true
tags: ["channel"]
example: "examples/c02-channels"
prev: "C1"
next: "C3"
---

## 本章你會建立的心智模型

Channel 是 CSP 風格的通訊原語。無緩衝 channel 的傳送與接收**同步交接**；有緩衝則多了佇列。**誰 close、誰寫**必須有清楚擁有者——對已關閉 channel 傳送會 panic。

## Python 對照

| Python | Go |
|--------|-----|
| `queue.Queue` | buffered channel 近似 |
| 條件變數通知 | channel 訊號 |

## L1 能用

```go
ch := make(chan int)      // 無緩衝
ch := make(chan int, 16)  // 有緩衝

ch <- 1
v := <-ch
v, ok := <-ch // ok=false 表示已關閉且空

close(ch)
```

範例：`examples/c02-channels`。

## L2 機制

- 只有傳送側（或約定擁有者）應 close。  
- 關閉後可繼續接收剩餘值。  
- `nil` channel 上的操作會永遠阻塞——可用於 select 禁用分支。  
- 不要用 channel「解決所有同步」；共享結構 + Mutex 有時更清楚。

## 請丟掉的 Python 習慣

1. 多處 close。  
2. 無界成長的工作佇列（無緩衝+狂發 or 超大 buffer）。  

## 遊戲 Server 連結

- 輸入命令：`chan Command` 有界  
- 結構事件：room 擁有的 inbox  
- 關房：close 或 context cancel  

## 練習

### 必做

1. 跑 `examples/c02-channels`。  
2. 示範對 closed channel 傳送的 panic（獨立小程式）。  

## 延伸閱讀

- <https://go.dev/ref/spec#Channel_types>  
