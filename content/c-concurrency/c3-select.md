---
lessonId: "C3"
title: "select：同時等很多路，含超時"
description: "多路等待命令／tick／取消；用 timer 做超時，用 default 做非阻塞嘗試。"
volume: "c"
order: 3
level: "l2"
status: "ready"
path_required: true
tags: ["select", "timeout"]
example: "examples/c03-select"
prev: "C2"
next: "C4"
---

## 這章你會搞懂什麼

房間迴圈很少只做一件事。它通常要同時盯：

- 玩家命令來了沒  
- tick 時間到了沒  
- 房間／行程要不要關了  

**`select`** 就是「同時等待多個 channel 操作，哪個就緒就走哪個」。  
多個都就緒時，Go 會**偽隨機**挑一條——別假設有固定優先順序。

搭配 `time.After`／`Ticker`／`context.Done()`，你就能寫出乾淨的超時與取消，而不用 busy-loop 一直 sleep 輪詢。

## Python 對照

| Python | Go |
|--------|-----|
| `asyncio.wait` / `wait_for` | `select` + timer／context |
| `queue.get(timeout=)` | `select`：收 queue 或收 timeout |
| 多個 future 競速 | 多個 case 競速 |

若你寫過 asyncio，select 的「等好幾件事之一」會很眼熟；差別是 Go 把它做成語言陳述句，而且跟 channel 綁很緊。

## 怎麼寫

```go
select {
case cmd := <-inbox:
	handle(cmd)
case <-time.After(50 * time.Millisecond):
	tick()
case <-ctx.Done():
	return ctx.Err()
}
```

固定節奏更常見的是 `time.NewTicker`，並 `defer ticker.Stop()`：

```go
ticker := time.NewTicker(50 * time.Millisecond)
defer ticker.Stop()

for {
	select {
	case cmd := <-inbox:
		handle(cmd)
	case <-ticker.C:
		tick()
	case <-ctx.Done():
		return ctx.Err()
	}
}
```

範例：`examples/c03-select`。

## 細節

### 為什麼多路就緒要「隨機」？

避免某一條 case 永遠餓著（starvation）的簡單策略。  
你若需要嚴格優先順序，要自己設計（例如先非阻塞檢查高優先 channel，再 select）——多數房間邏輯用「公平挑選」就夠。

### `default`：非阻塞

```go
select {
case ch <- job:
	// 送進去了
default:
	// 現在塞不進去，立刻返回
}
```

這在背壓（C8）很常用：佇列滿就拒絕／丟棄，而不是拖垮呼叫端。

### `time.After` 在熱迴圈的坑

每次 `time.After` 都會產生一個 Timer。若在極熱迴圈又常常「別的 case 先就緒」，可能累積待回收的 timer。  
進階作法：`time.NewTimer` + `Stop`／`Reset`。新手先用 `Ticker` 做固定 tick 通常更單純。

### 進階可先略過

- 空的 `select {}` 會永遠阻塞——有時用來停住 main，但可讀性一般。  
- 跟 nil channel 搭配可動態啟停分支（C2）。

## 遊戲 Server 會用在哪

Arena 風格房間迴圈幾乎就是：

1. 收 input  
2. 到點就 tick（例如 20Hz＝50ms）  
3. 聽到取消就收拾離開  

不要寫成「sleep 50ms → 再 poll 有沒有命令」——那會讓延遲變差，也比較難跟取消組在一起。

## 請丟掉的舊習慣

1. busy-loop + sleep 輪詢代替事件。  
2. 假設 select 的 case 順序＝優先順序。  
3. 超時了還不回傳／不取消下游（後面 C4 補齊）。

## 動手練習

### 必做

1. 跑 `examples/c03-select`。  
2. 改成 100ms tick，印計數，約 3 秒後結束。  

### 選做

1. 加一個 `default` 版本的「嘗試送入 inbox」，滿了印 `busy`。  

## 常見坑

- **每個 loop 都 `time.After` 又極熱**：改 Ticker／重用 Timer。  
- **忘記 `Stop` ticker**：小程式沒事，長駐服務會漏資源。  
- **在 case 裡做很重的事還持有不該持有的鎖**：先想清楚臨界區（C5）。

## 延伸閱讀

- <https://go.dev/ref/spec#Select_statements>  
