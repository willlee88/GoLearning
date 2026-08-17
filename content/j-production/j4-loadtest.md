---
lessonId: "J4"
title: "壓測方法"
description: "假 client、漸進加壓、看哪裡先爆。"
volume: "j"
order: 4
level: "l2"
status: "ready"
path_required: true
tags: ["loadtest"]
example: "examples/j04-load-client"
prev: "J3"
next: "J5"
---

## 本章你會建立的心智模型

壓測要有**假設**：例如「單進程 1k 連線、20Hz 廣播是否 CPU 打滿」。步驟：

1. 定場景（連線數、房人數、訊息率）  
2. 假 client 漸進加壓  
3. 看 metrics + pprof  
4. 找到瓶頸（廣播 O(N²)、鎖、JSON）  

## Python 對照

locust / 自寫 asyncio client；Go 假 client 更輕。

## L1 能用

```bash
go run ./examples/j04-load-client -addr ws://127.0.0.1:8080/ws -n 50 -room load
```

範例：`examples/j04-load-client`。

## L2 機制

- 先測本機 loopback，再測跨網。  
- 區分「連得上」與「狀態仍正確」。  
- `-race` 不開在壓測熱路徑（太慢）。  

## 請丟掉的 Python 習慣

1. 一次開滿然後只看有沒有 crash。  
2. 沒有基線就優化。  

## 遊戲 Server 連結

對 Arena Mini 加壓時觀察 `/metrics`。

## 練習

### 必做

1. 啟動 arena-mini，跑 load-client 20 連線。  
2. 記錄 connections 與 messages。  

## 延伸閱讀

- pprof 官方部落格  
