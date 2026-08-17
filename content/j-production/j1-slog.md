---
lessonId: "J1"
title: "結構化日誌 slog"
description: "log/slog 鍵值日誌與請求／房間關聯。"
volume: "j"
order: 1
level: "l2"
status: "ready"
path_required: true
tags: ["logging"]
example: ""
prev: "J0"
next: "J2"
---

## 本章你會建立的心智模型

`log/slog`（Go 1.21+）讓日誌帶 **屬性**：`room_id`、`player`、`tick`。搜尋與告警靠欄位，不是靠人肉 parse 字串。

## Python 對照

| Python | Go |
|--------|-----|
| `logging` + structlog | `log/slog` |
| `extra={}` | `slog.String("room", id)` |

## L1 能用

```go
slog.Info("player joined",
	"room", roomID,
	"player", name,
)
```

## L2 機制

- JSON handler 利於收集。  
- 日誌等級：Debug/Info/Warn/Error。  
- **勿打密碼、token 全文**。  
- 高頻 tick 不要每 tick Info（用 metrics）。  

## 請丟掉的 Python 習慣

1. `print` 除錯留在生產。  
2. 無關聯 id 的散落字串。  

## 遊戲 Server 連結

Arena Mini M5 使用 slog。

## 練習

### 必做

1. 把一句 `log.Printf` 改成 slog 並帶 room。  

## 延伸閱讀

- <https://pkg.go.dev/log/slog>  
