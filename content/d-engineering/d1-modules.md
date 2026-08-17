---
lessonId: "D1"
title: "go.mod 與依賴心智"
description: "module path、語意化版本、MVS 直觀。"
volume: "d"
order: 1
level: "l2"
status: "ready"
path_required: true
tags: ["modules"]
example: ""
prev: "D0"
next: "D2"
---

## 本章你會建立的心智模型

`go.mod` 定義模組路徑與依賴版本。Go 用 **Minimal Version Selection**：選滿足所有約束的最小版本，行為比「永遠最新」可預期。

## Python 對照

| Python | Go |
|--------|-----|
| requirements / lock | go.mod + go.sum |
| venv | 通常不需要 |

## L1 能用

```bash
go mod init example.com/x
go get example.com/y@v1.2.3
go mod tidy
```

## L2 機制

- `go.sum` 校驗完整性，勿手改。  
- `replace` 本機聯調。  
- `internal/` 防止外部匯入。  

## 請丟掉的 Python 習慣

1. 全局亂裝套件不鎖版本。  

## 遊戲 Server 連結

demo 依賴越少越好；WS 只用 `x/net` 是刻意選擇。

## 練習

### 必做

1. 打開 `demo/arena-mini/server/go.mod` 讀懂 require。  

## 延伸閱讀

- Modules reference  
