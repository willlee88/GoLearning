---
lessonId: "E8"
title: "sync 套件地圖"
description: "Mutex、WaitGroup、Once、Map、Pool 何時用。"
volume: "e"
order: 8
level: "l2"
status: "ready"
path_required: false
tags: ["sync", "stdlib"]
example: ""
prev: "E7"
next: "F0"
---

## 本章你會建立的心智模型

C5 講過 Mutex；本章是 **sync 工具箱索引**。

| 型別 | 用途 |
|------|------|
| Mutex / RWMutex | 互斥 |
| WaitGroup | 等一組 goroutine |
| Once | 只跑一次 |
| Map | 特例併發 map（有取捨） |
| Pool | 重用 transient 物件 |
| Cond | 條件等待（少用） |

## L1 能用

```go
var once sync.Once
once.Do(func() { initAll() })
```

## L2 機制

- `sync.Map` 適合 key 穩定、寫少讀多或不相交 key；一般業務更常用自己的 `map+Mutex`。  
- `Pool` 的 Get 可能拿到「任意」物件，要 Reset。  
- 不要複製含 Mutex 的 struct。  

## 請丟掉的 Python 習慣

1. 無腦 `sync.Map` 當 concurrent dict。  

## 遊戲 Server 連結

Hub 的 rooms map + Mutex；`tickOnce sync.Once` 啟動 ticker。

## 練習

### 必做

1. 說明 Arena Hub 為何不用 `sync.Map`。  

## 延伸閱讀

- `sync` package  

## 接回主路徑

若尚未讀網路卷 → **F0**。  
