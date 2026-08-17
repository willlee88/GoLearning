---
lessonId: "B2"
title: "wrap 與錯誤鏈"
description: "fmt.Errorf %w、errors.Is/As 與上下文。"
volume: "b"
order: 2
level: "l2"
status: "ready"
path_required: true
tags: ["errors"]
example: "examples/b02-room-errors"
prev: "B1"
next: "B3"
---

## 本章你會建立的心智模型

錯誤向上返回時應**加上上下文**並保留根因：`fmt.Errorf("join room %s: %w", id, err)`。`errors.Is` / `errors.As` 能沿鏈查找哨兵或型別。

## Python 對照

| Python | Go |
|--------|-----|
| `raise X from e` | `%w` wrap |
| `__cause__` | `errors.Unwrap` 鏈 |

## L1 能用

```go
if err := store.Save(r); err != nil {
	return fmt.Errorf("save room %s: %w", r.ID, err)
}
```

## L2 機制

- `%w` 才能被 `Is/As` 解開；`%v` 只變字串。  
- 自訂型別實作 `Is` / `Unwrap` 可進階控制。  
- 日誌印 `%+v` 與否依你是否用支援堆疊的套件（標準庫無堆疊）。

## 請丟掉的 Python 習慣

1. 只 `return err` 一路無上下文。  
2. wrap 時丟掉哨兵（應 `%w`）。  

## 遊戲 Server 連結

```text
apply move: validate: illegal cell
```

外層給人看，內層哨兵給程式分支。

## 練習

### 必做

1. 跑 `go test` 於 `examples/b02-room-errors`。  
2. 斷言 `errors.Is` 在 wrap 後仍成立。  

## 延伸閱讀

- <https://go.dev/blog/go1.13-errors>  
