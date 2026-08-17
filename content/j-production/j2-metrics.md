---
lessonId: "J2"
title: "Metrics 與 /metrics"
description: "計數器、房間數、連線數；先 JSON 後 Prometheus。"
volume: "j"
order: 2
level: "l2"
status: "ready"
path_required: true
tags: ["metrics"]
example: "examples/j02-metrics"
prev: "J1"
next: "J3"
---

## 本章你會建立的心智模型

Metrics 回答「現在系統怎樣」：在線連線、房間數、訊息速率、錯誤率、tick 耗時。教學先用 **原子計數 + JSON `/metrics`**；上線可換 Prometheus 格式。

## Python 對照

prometheus_client 同概念；先有正確標籤再建儀表板。

## L1 能用

```go
type Metrics struct {
	Connections atomic.Int64
	Rooms       atomic.Int64
	MessagesIn  atomic.Int64
}
```

範例：`examples/j02-metrics`。

## L2 機制

- 計數器只增；ゲージ可上下。  
- 直方圖看延遲分位（進階）。  
- 標籤基數爆炸會打爆時序庫。  

## 請丟掉的 Python 習慣

1. 只靠日誌當監控。  
2. 每個玩家一個 metric 標籤。  

## 遊戲 Server 連結

Arena Mini：`GET /metrics`。

## 練習

### 必做

1. 打開 `/metrics` 對照開房前後。  
2. 為 `messages_out` 加計數。  

## 延伸閱讀

- USE / RED 方法論（概念）  
