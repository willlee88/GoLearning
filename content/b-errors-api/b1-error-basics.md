---
lessonId: "B1"
title: "error 是值：哨兵錯誤怎麼用"
description: "搞懂 error 介面、哨兵 ErrXxx，以及為什麼要比對值而不是比對字串。"
volume: "b"
order: 1
level: "l2"
status: "ready"
path_required: true
tags: ["errors"]
example: "examples/b02-room-errors"
prev: "B0"
next: "B2"
---

## 這章你會搞懂什麼

在 Go 裡，**錯誤是值**。`error` 只是一個很小的介面：有 `Error() string` 就算。  
業務上常會定義 **哨兵錯誤（sentinel error）**——例如 `var ErrRoomFull = errors.New("room full")`——代表一個**穩定、可分支的條件**。

呼叫端不要去比字串（字串一改就全炸），而是用 `errors.Is(err, ErrRoomFull)` 問：「這是不是滿房那種錯？」

讀完這章你要能：定義少數幾個哨兵、在函式裡回傳它們、在上層用 `Is` 轉成給客戶端的錯誤碼。

## Python 對照

| Python | Go | 白話差在哪 |
|--------|-----|------------|
| `raise RoomFull()` | `return ErrRoomFull` | 錯誤走回傳，不跳控制流 |
| `except RoomFull:` | `errors.Is(err, ErrRoomFull)` | 比的是**值／鏈**，不是例外型別階層 |
| 例外訊息當分支條件 | 幾乎永遠別這麼做 | 字串一改，行為就 silently 壞掉 |
| 很深的例外繼承樹 | 少數哨兵 + 必要時自訂型別 | Go 偏好「少而穩」 |

你在 Python 可能習慣「丟一個型別，上面 catch」。Go 比較像：**把一個約定好的值傳回去，上面用 `Is` 認人。**

## 怎麼寫

```go
var ErrRoomFull = errors.New("room full")

func Join(n, cap int) error {
	if n >= cap {
		return ErrRoomFull
	}
	return nil
}

if errors.Is(err, ErrRoomFull) {
	// 回傳協定碼 ROOM_FULL 給客戶端
}
```

成功時慣例是 **`err == nil`**。失敗時結果值通常當作「不可用」（除非 API 文件明確說半成品也有意義）。

本站範例：`examples/b02-room-errors`（裡面還有 wrap，下一章會細講）。

## 細節

### 為什麼要有哨兵，而不是每次 `errors.New("room full")`？

因為每次 `errors.New` 都是**新的值**。字串一樣，但 `==` / `errors.Is` 對不上同一個哨兵。  
哨兵的重點是：**全專案共用同一個變數**，行為才穩定。

### 為什麼不要比對 `err.Error()` 字串？

- 訊息常會被 wrap 加上上下文（下一章）  
- 翻譯、標點、大小寫一改，分支就錯  
- 測試會變脆  

人看的是字串；**程式分支看的是哨兵或型別**。

### 哨兵要少而精

不是每個小失敗都要一個 `ErrXxx`。問自己：

> 上層會不會需要**不同處理**？（不同錯誤碼、要不要重試、要不要踢人）

要 → 值得哨兵。  
只是日誌多寫一句 → wrap 上下文就夠，不必新哨兵。

### 進階可先略過

- 自訂型別實作 `Error()`，再用 `errors.As` 取出欄位（B2 會碰到）。  
- 有些庫用「哨兵 + 包裝型別」同時帶穩定碼與細節。

## 遊戲 Server 會用在哪

| 情況 | 建議 |
|------|------|
| 滿房、玩家不存在、階段不允許加入 | 哨兵 → 轉協定錯誤碼 |
| 磁碟／DB／網路暫時失敗 | 多半 wrap 底層 error，必要時另訂可重試策略 |
| 內部除錯「哪個房間」 | 訊息裡帶 ID；分支仍靠哨兵 |

客戶端可見錯誤要**穩定**；內部日誌要**有上下文**。兩件事別混成「只丟一句英文」。

## 請丟掉的舊習慣

1. 用例外訊息字串當控制流（`if "room full" in str(e)`）。  
2. 裸 `except:` / 無腦 `_ = do()` 吞掉錯誤。  
3. 為每個小失敗發明巨大例外階層——在 Go 先問「要不要不同分支」。

## 動手練習

### 必做

1. 閱讀 `examples/b02-room-errors` 的哨兵定義（`ErrRoomFull` 等）。  
2. 新增 `ErrNotReady`，在某個不合法操作回傳它，並用 `errors.Is` 測到。  

### 選做

1. 列一張表：你的遊戲裡哪些錯誤「客戶端要分碼」、哪些「只記日誌」。  

## 常見坑

- **每次函式裡 `return errors.New("room full")`**：字串相同但不是同一個哨兵，`Is` 會失敗。  
- **成功時回傳非 nil error、或失敗時還依賴結果值**：先約定「`err != nil` 則結果別用」。  
- **哨兵訊息寫得很口語又常改**：訊息給人看可以好讀，但別靠它做邏輯。

## 延伸閱讀

- <https://pkg.go.dev/errors>  
- <https://go.dev/blog/error-handling-and-go>  
