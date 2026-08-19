---
lessonId: "P0.4"
title: "錯誤是回傳值，不是例外跳轉"
description: "Go 用 error 當值往上傳，控制流留在函式裡。遊戲 Server 的滿房、逾時、非法操作都是日常路徑。"
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

## 這章你會搞懂什麼

Go 把錯誤當成 **`error` 介面值** 回傳，而不是預設用例外（exception）往上跳。

控制流留在函式裡：呼叫者必須決定——處理、包裝後再傳、還是在邊界變成給客戶端的錯誤碼。

這在遊戲 Server 特別重要：非法操作、斷線、逾時、房間已滿，都是**每天都會發生的路徑**，不該長得像「程式災難」。讀完你要能寫出 `if err != nil`，並知道什麼時候該用 `%w` 包裝上下文。

## 先跟 Python 對一下

| Python | Go |
|--------|-----|
| `raise` / `try` / `except` | `return ..., err` |
| 例外階層（繼承一堆 Error 類別） | `errors.Is` / `errors.As` + 包裝鏈 |
| 常忘了接的例外，執行期才爆 | 未處理的 `err` 靠習慣與 linter 抓住 |
| `finally` 做清理 | `defer` 做清理（保證函式離開前跑） |

Python 裡「用例外當正常分支」很常見；在 Go 社群這通常被看成味道不對——正常失敗請回傳 `error`。

## 怎麼寫（能跑的最小例子）

```go
func load(path string) ([]byte, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return b, nil
}
```

慣例：

- 成功 → `err == nil`  
- 失敗 → 回傳有意義的 `error`（常搭配零值／`nil` 資料）  

包裝上下文（之後查 log 才知道是哪一步炸的）：

```go
b, err := os.ReadFile(path)
if err != nil {
	return nil, fmt.Errorf("load config %s: %w", path, err)
}
```

`%w` 會把底層錯誤包進鏈裡，之後可用 `errors.Is`／`errors.As` 辨認。

## 為什麼這樣設計／底層在幹嘛

### `error` 是介面，不是特殊語法糖

只要實作 `Error() string`，就是一個 `error`。標準庫很多函式回傳它；你自己的業務錯誤也可以是哨兵變數或自訂型別。

### 為什麼到處都是 `if err != nil`？

因為設計者想讓**失敗路徑看得見**。代價是囉嗦；好處是讀程式時不太會漏掉「這裡可能失敗」。

寫錯會怎樣？

- `_ = doSomething()` 無腦丟掉錯誤 → 線上靜默壞掉，超難查  
- 每個錯誤都 `panic` → 一個壞請求可能整支 Server 掛掉  
- 不包裝上下文 → log 只剩 `EOF`，不知道哪個房間、哪個檔  

### `panic` 還在，但不是日常武器

`panic` 比較像「程式不變量被破壞」（例如演算法保證不會走到的分支卻走到了）。函式庫裡亂 `panic`，呼叫者很難優雅恢復——code review 通常直接紅燈。

遊戲業務：「房間滿了」「ticket 過期」→ **回傳 error**，不要 panic。

### 進階可先略過

- 哨兵錯誤（sentinel，例如 `var ErrRoomFull = errors.New("room full")`）vs 自訂型別錯誤。  
- 為何早期 Go 沒有內建 `?` 運算子：明確性 vs 雜訊的社群辯論——你先把習慣練熟比等語法糖重要。  

## 遊戲 Server 會用在哪

這些都應該是**可預期的 error**，方便轉成協定錯誤碼給客戶端：

- 加入已滿房  
- ticket 過期  
- 目前階段不允許開始遊戲  
- 輸入非法（移動到牆外、重複 Ready…）  

不要把 stack trace 直接甩到連線上——客戶端要的是穩定錯誤碼與人話訊息；你自己的 log 才留詳細上下文（房間 ID、玩家 ID）。

範例 `examples/p0-config-stats` 就有檔案／JSON 錯誤的回傳與包裝，值得對照。

## 請丟掉的舊習慣

1. 用例外表達「使用者輸入不合法」的正常分支——在 Go 回傳 `error`。  
2. 裸 `except:` 吞掉一切——在 Go 同樣禁止無腦 `_ = do()`（除非你真的有意識、並註明為什麼）。  
3. 只靠 stack trace 除錯——應包裝上下文（哪個房間、哪個玩家、哪一步）。  

## 動手練習

### 必做

1. 閱讀 `examples/p0-config-stats` 如何回傳與包裝檔案／JSON 錯誤。  
2. 故意給錯路徑跑一次，看錯誤訊息有沒有足夠上下文。  

### 選做

1. 定義 `var ErrRoomFull = errors.New("room full")`，在測試裡用 `errors.Is` 判斷。  

## 常見坑

- **`if err != nil` 疲勞**：用小函式、一致的包裝格式；先把習慣練起來，再談怎麼減少樣板。  
- **library 裡亂 `panic`**：業務可預期失敗請回傳 error。  
- **包了卻沒用 `%w`**：之後 `errors.Is` 對不到底層錯誤，測試與處理會很痛苦。  

## 延伸閱讀

- <https://go.dev/blog/error-handling-and-go>  
- <https://pkg.go.dev/errors>  
