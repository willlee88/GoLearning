---
lessonId: "E4"
title: "filepath 與 embed 概念"
description: "路徑可攜；把靜態檔編進 binary。"
volume: "e"
order: 4
level: "l2"
status: "ready"
path_required: false
tags: ["filepath", "embed"]
example: ""
prev: "E3"
next: "E5"
---

## 本章你會建立的心智模型

`filepath` 處理 OS 路徑差異；`path` 用於 URL。`embed` 可把靜態前端打進 binary，部署變單檔。

## L1 能用

```go
//go:embed web/*
var webFS embed.FS
```

## 遊戲 Server 連結

可把 Arena 的 `web/` embed 進 server（選做重構）。

## 練習

### 必做

1. 讀 `embed` 文件示例。  
2. 說明為何不要用字串 `+ "/"` 拼 Windows 路徑。  

## 延伸閱讀

- `embed` package  

## 接回主路徑

下一站 **F0 網路卷**（若尚未讀）。  
