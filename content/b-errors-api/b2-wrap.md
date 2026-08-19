---
lessonId: "B2"
title: "wrap：幫錯誤加上上下文，還留得住根因"
description: "用 fmt.Errorf 與 %w 把「哪一步失敗」寫進錯誤鏈，並用 Is/As 往下找。"
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

## 這章你會搞懂什麼

錯誤往上返回時，最怕變成一句冷冰冰的 `room full`，卻不知道**哪個房間、哪一步**。  
解法是 **wrap（包裝）**：外面加上下文，裡面仍保留原本的錯誤，讓 `errors.Is` / `errors.As` 還能沿著鏈找到哨兵或型別。

關鍵寫法：

```go
fmt.Errorf("join room %s: %w", id, err)
```

注意是 **`%w`**（wrap），不是 `%v`。  
`%v` 只是把錯誤印進字串，鏈就斷了，上層 `Is` 會找不到哨兵。

## Python 對照

| Python | Go |
|--------|-----|
| `raise X from e` | `fmt.Errorf("...: %w", err)` |
| `__cause__` / exception chaining | `errors.Unwrap` 形成的鏈 |
| 只 `raise` 新例外、丟掉舊的 | 等同用 `%v` 或只回新字串——**別這樣** |

白話：Python 的 `from` 是在說「因為這個才raise那個」；Go 的 `%w` 是在說「我加一句話，但根因還在鏈上」。

## 怎麼寫

```go
if err := store.Save(r); err != nil {
	return fmt.Errorf("save room %s: %w", r.ID, err)
}
```

房間範例裡常見形狀：

```go
if len(r.Players) >= r.Capacity {
	return fmt.Errorf("join %s: %w", r.ID, ErrRoomFull)
}
```

上面這行同時做了兩件事：

1. **給人看**：`join room-42: room full`  
2. **給程式看**：`errors.Is(err, ErrRoomFull)` 仍然是 true  

跑起來看：`examples/b02-room-errors`。

## 細節

### `%w` 跟 `%v` 差在哪（為什麼常踩坑）

| 動詞 | 結果 |
|------|------|
| `%w` | 包進錯誤鏈，`Is`/`As`/`Unwrap` 有效 |
| `%v` / `%s` | 變成普通字串的一部分，**哨兵丟了** |

所以：要保留可分支的根因 → 一定 `%w`。  
若你只是寫進 log 一行字，那是日誌問題，不是回傳問題。

### `errors.Is` / `errors.As` 在找什麼

- **`Is(err, target)`**：沿著鏈問「有沒有等於／匹配這個哨兵？」  
- **`As(err, &ptr)`**：沿著鏈找「有沒有這個型別」，找到就填進指標  

你上層通常這樣分工：

- 穩定條件 → `Is`（滿房、未找到）  
- 需要欄位（例如 HTTP 狀態、偏移量）→ `As`

### 為什麼「只 `return err`」常常不夠？

底層可能只說 `EOF` 或 `connection reset`。沒有「save room」「apply move」這層，半夜看日誌會想哭。  
每跨一個**有意義的邊界**（儲存、房間、撮合），就加一點上下文——但也不要每一行都包到噪音。

### 進階可先略過

- 自訂型別可實作 `Unwrap() error` 或 `Is(error) bool` 來控制匹配。  
- 標準庫錯誤**沒有**內建堆疊；若要堆疊需額外套件或自己在邊界記錄 `debug.Stack`（通常在 panic 路徑，見 B3）。

## 遊戲 Server 會用在哪

理想日誌／錯誤字串長得像：

```text
apply move: validate: illegal cell
join room-7: room full
```

- **外層**：人很快定位「哪條路徑」  
- **內層哨兵**：程式決定回哪個協定碼、要不要重試  

別把 stack 直接甩給客戶端；客戶端要的是穩定碼，不是你的函式名稱。

## 請丟掉的舊習慣

1. 一路 `return err`，到 API 邊界才發現訊息完全沒上下文。  
2. wrap 時用 `%v`，結果哨兵蒸發，測試裡 `Is` 莫名失敗。  
3. 把所有細節塞進一個超長字串，卻沒有可 `Is` 的穩定條件。

## 動手練習

### 必做

1. 在 `examples/b02-room-errors` 跑 `go test`。  
2. 確認 wrap 之後 `errors.Is(err, ErrRoomFull)` 仍為 true。  

### 選做

1. 故意改成 `%v` 再跑同一測試，親眼看 `Is` 失敗——然後改回 `%w`。  

## 常見坑

- **`fmt.Errorf("...%w", nil)`**：別把 nil 拿去 wrap；成功路徑直接 `return nil`。  
- **包太深、訊息重複三次「failed to」**：上下文要有資訊量（ID、階段），不是同義詞堆疊。  
- **在熱迴圈配字串**：通常沒關係；真熱再測。正確性優先。

## 延伸閱讀

- <https://go.dev/blog/go1.13-errors>  
