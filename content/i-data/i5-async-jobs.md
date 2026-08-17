---
lessonId: "I5"
title: "背景任務與結算"
description: "有界佇列、worker、至少一次投遞。"
volume: "i"
order: 5
level: "l2"
status: "ready"
path_required: true
tags: ["async", "workers"]
example: "examples/i05-worker"
prev: "I4"
next: "J0"
---

## 本章你會建立的心智模型

對局結束後：寫結果、發獎勵、分析事件——應**離開 tick 熱路徑**。模式：

```text
Room.End → jobs <- MatchResult → worker → DB / 郵件
```

有界 queue 滿了要背壓（拒絕開新局或丟棄非關鍵分析）。

## Python 對照

| Python | Go |
|--------|-----|
| Celery / RQ | channel + worker pool / 外部 MQ |
| `asyncio.create_task` 忘了管 | 明確 WaitGroup 與 shutdown |

## L1 能用

```go
jobs := make(chan Job, 1024)
for i := 0; i < n; i++ {
	go worker(jobs)
}
```

範例：`examples/i05-worker`。

## L2 機制

- 至少一次：worker 可能重試 → 結算要冪等。  
- 關服：停止接收 → drain queue → 退出。  
- 關鍵路徑用 Outbox 表（進階）。  

## 請丟掉的 Python 習慣

1. 在請求 handler 同步寄信 3 秒。  
2. 無界 `go func` 狂丟任務。  

## 遊戲 Server 連結

Arena Mini 可在 `PhaseEnded` 時非同步 log 一筆 JSON 結果（M5 選做）。

## 練習

### 必做

1. 跑 `examples/i05-worker`。  
2. 改成 queue 滿時 `Push` 回傳 error。  

## 延伸閱讀

- J4 優雅關閉  
