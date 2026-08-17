---
lessonId: "I2"
title: "database/sql 模式"
description: "標準庫資料庫介面、連線池、context。"
volume: "i"
order: 2
level: "l2"
status: "ready"
path_required: true
tags: ["sql"]
example: "examples/i02-sql-memory"
prev: "I1"
next: "I3"
---

## 本章你會建立的心智模型

`database/sql` 是驅動無關的介面。你拿到 `*sql.DB`（連線池），用 `QueryContext` / `ExecContext` 並傳 `context`。教學可用 **sqlite 或純介面 mock**；本站範例用記憶體 store 模擬 repository 模式，避免強制裝 DB。

## Python 對照

| Python | Go |
|--------|-----|
| SQLAlchemy session | `*sql.DB` + 顯式 SQL 或輕量包 |
| `with conn` | `defer rows.Close()` + context |

## L1 能用

```go
type PlayerRepo interface {
	Upsert(ctx context.Context, p Player) error
	Get(ctx context.Context, id string) (Player, error)
}
```

實作可換：memory / postgres / sqlite。

範例：`examples/i02-sql-memory`。

## L2 機制

- `DB` 應長期共用，不要 per-request `sql.Open`。  
- 設 `SetMaxOpenConns` 等池參數。  
- 掃描用 `rows.Scan`；注意 null。  
- 遷移（migrate）與應用程式版本一起管。  

## 請丟掉的 Python 習慣

1. 字串拼接 SQL（注入）。用參數佔位符。  
2. 忽略 `rows.Err()`。  

## 遊戲 Server 連結

帳號、庫存、賽季分走 Repo；Room 不持有 `*sql.DB` 在熱路徑。

## 練習

### 必做

1. 跑 `examples/i02-sql-memory`。  
2. 為 repo 加 `ListByScore`。  

## 延伸閱讀

- <https://go.dev/doc/database/>  
