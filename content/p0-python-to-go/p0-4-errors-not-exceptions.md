---
lessonId: "P0.4"
title: "錯誤不是例外"
description: "用值傳遞錯誤，建立看得見的控制流。"
volume: "p0"
order: 4
level: "l2"
status: "ready"
path_required: true
tags: ["python-bridge", "errors"]
example: "examples/p0-config-stats"
prev: "P0.3"
next: "P0.5"
---

## 本章你會建立的心智模型

Go 把錯誤當成 **`error` 介面值** 來回傳，而不是預設用例外往上跳。控制流留在函式裡：呼叫者必須決定要處理、包裝，還是回傳。這在遊戲 Server 裡特別重要——非法操作、斷線、逾時都是**日常路徑**，不該長得像災難。

## Python 對照

| Python | Go |
|--------|-----|
| `raise` / `try` / `except` | `return ..., err` |
| 例外階層 | `errors.Is` / `errors.As` + 包裝 |
| 常忘了接的例外 | 未處理的 `err` 會被 linter／習慣抓住 |
| `else`/`finally` | `defer` 做清理 |

## L1 能用

```go
func load(path string) ([]byte, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return b, nil
}
```

慣例：成功時 `err == nil`。

## L2 機制

- `error` 是介面：`Error() string`。
- 用 `fmt.Errorf("load config: %w", err)` 包裝以保留鏈（`%w`）。
- `panic` 存在，但**不是日常業務錯誤**；類似「程式不變量被破壞」。

## L3 深潛（可選）

- 哨兵錯誤（sentinel）vs 自訂型別錯誤。
- 為何早期 Go 沒有內建 `?` 運算子：明確性 vs 雜訊的社群辯論。

## 請丟掉的 Python 習慣

1. 用例外表達「使用者輸入不合法」的正常分支（在 Go 回傳 error）。
2. 裸 `except:` 吞掉一切——在 Go 同樣禁止 `_ = do()` 無腦丟棄（除了有意識的情況）。
3. 依賴 stack trace 當唯一除錯資訊——應包裝上下文（房間 ID、玩家 ID）。

## 遊戲 Server 連結

「加入已滿房」「ticket 過期」「狀態不允許開始」都應是**可預期的 error**，方便轉成協定錯誤碼給客戶端，而不是把 stack 甩到連線上。

## 練習

### 必做

1. 閱讀 `examples/p0-config-stats` 如何回傳與包裝檔案／JSON 錯誤。
2. 故意給錯路徑，觀察錯誤訊息是否含上下文。

### 選做

1. 定義 `var ErrRoomFull = errors.New("room full")`，在測試裡用 `errors.Is`。

## 常見坑與如何看見

- `if err != nil` 疲勞：用小函式、一致包裝；工具 `errcheck` / 編譯器未來方向不論，**習慣先建立**。
- `panic` 在 library 裡亂甩：code review 紅燈。

## 延伸閱讀

- <https://go.dev/blog/error-handling-and-go>
- <https://pkg.go.dev/errors>
