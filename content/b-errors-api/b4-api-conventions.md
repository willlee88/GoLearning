---
lessonId: "B4"
title: "API 慣例：好測、好組、少驚喜"
description: "接受小介面、回傳具體型別；用 NewX 保證不變式，錯誤可被 Is/As。"
volume: "b"
order: 4
level: "l2"
status: "ready"
path_required: true
tags: ["api", "design"]
example: "examples/b02-room-errors"
prev: "B3"
next: "C0"
---

## 這章你會搞懂什麼

好的 Go API 讀起來常常有點「無聊」——但這是優點。你半夜改 Server，最想看到的是：

1. 依賴用**小介面（interface）**注入，方便假雙／替換  
2. 回傳**具體型別**（struct），呼叫端不猜  
3. 建構函式 `NewX(...) (*X, error)` 在門口就把不變式檢查完  
4. 錯誤走哨兵／wrap，測試能用 `errors.Is` 鎖定行為  

這章把 B1–B3 收到「怎麼對外說話」：不是多學一個關鍵字，而是讓模組邊界站得住。

## Python 對照

| Python | Go |
|--------|-----|
| DI 框架／很大的建構子注入 | 手寫 `New` + struct 裡放小介面欄位 |
| 屬性事後亂設、半初始化物件 | 未匯出欄位 + `New` 保證「生出來就合法」 |
| 全域單例到處 import | 顯式傳依賴；少藏在 `init`／套件全域 |
| 例外當 API 合約的一部分 | `error` 值 + 文件化哨兵 |

## 怎麼寫

```go
type Store interface {
	Save(*Room) error
}

type Server struct {
	store Store
}

func NewServer(store Store) (*Server, error) {
	if store == nil {
		return nil, errors.New("store nil")
	}
	return &Server{store: store}, nil
}
```

房間也一樣：`examples/b02-room-errors` 的 `New(id, capacity)` 會拒絕 `capacity <= 0`。  
非法設定寧願在建構時失敗，也不要讓 `var Room{}` 以「零值看起來像能跑」的狀態流進戰局。

## 細節

### 為什麼「接受介面、回傳結構」？

- **接受介面**：你只要求「我會 Save」，測試可塞假 Store，正式環境塞真 DB。  
- **回傳結構**：呼叫端拿到的是完整能力；不會被迫依賴一個你隨手抽的大介面。  

介面要**小**。一個方法的介面最好組；上帝介面（十幾個方法）會讓 mock 跟實作一起痛苦（呼應 A12）。

### `New` 在保證什麼？

「物件存在 ⇒ 不變式成立」。例如：

- capacity > 0  
- map 已 `make`（不是 nil map）  
- 必要依賴非 nil  

為什麼不靠呼叫端自己設欄位？因為你無法審計每一個 `Room{...}` 字面值。未匯出欄位 + 唯一建構入口，會少很多「半成品房間」。

### 關於 `init()` 與 functional options

- `init()` 裡做可能失敗的重邏輯（讀檔、連網）很難測、也難控啟動順序——能避則避。  
- functional options（`WithTimeout(...)` 這種）是進階糖，小專案用清楚的 `Config` struct 通常更夠用。**勿為炫技上 options。**

### 進階可先略過

- 回傳介面的少數合理場景（例如要隱藏多種實作）。預設仍回具體型別。  
- 相容性：公開函式簽名一改就是 breaking change——遊戲協定與 Go API 都要有版本意識。

## 遊戲 Server 會用在哪

```go
func NewRoom(cfg Config) (*Room, error)  // 驗證 capacity、tick rate
func NewHub(log Logger) *Hub             // 注入日誌／metrics
```

- Registry、Hub、Room worker 的依賴在 `main`／`cmd` 組裝  
- 規則邏輯套件不直接 import 具體 DB 驅動，只認小介面  

這樣 CI 才能對 `internal/game` 開一堆表驅動測試，不必先起 Redis。

## 請丟掉的舊習慣

1. 全域單例到處 `import`，測試時無法替換。  
2. 半初始化物件（nil map、capacity=0）流進正式路徑。  
3. 為每個 struct 先寫巨大 interface「以後也許用得到」。

## 動手練習

### 必做

1. 為 registry 或 room 加 `New`，拒絕非法 capacity。  
2. 抽一個 `Clock` 介面（`Now() time.Time`），讓超時邏輯可測。  

### 選做

1. 把某個套件級全域變數改成 `New` 注入，寫一則測試塞假依賴。  

## 常見坑

- **`New` 回了物件，不變式卻還能被外部改壞**：欄位該小寫就小寫。  
- **介面定義在「實作方」而不是「使用方」**：誰用誰定義小介面，通常比較靈活。  
- **錯誤不可 `Is`**：公開 API 的預期失敗請給哨兵或穩定型別，不要只回裸字串每次 `New`。

## 延伸閱讀

- <https://go.dev/doc/effective_go>  

## B 卷檢查點

你現在應該能把「滿房／不存在／非法階段」做成**哨兵 + wrap**，並在測試用 `errors.Is` 鎖住行為；公開入口用 `New` 擋下非法設定。  

下一卷 C：goroutine、channel、race——錯誤處理的肌肉會一直用到。  
