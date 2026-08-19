---
lessonId: "P0.3"
title: "型別怎麼想？靜態檢查與零值"
description: "從 Python 動態型別切到 Go：編譯期契約、零值（zero value）是設計的一部分、轉換必須寫清楚。"
volume: "p0"
order: 3
level: "l2"
status: "ready"
path_required: true
tags: ["python-bridge", "types"]
example: ""
prev: "P0.2"
next: "P0.4"
---

## 這章你會搞懂什麼

在 Go 裡，**型別是編譯期契約**：什麼能賦值、什麼能呼叫、介面有沒有被滿足，多數在 `go build` 時就定生死。

還有一個 Python 腦要改的點：**零值（zero value）**。變數「沒手動初始化」常常仍是可用的預設（`0`、`""`、`false`），不是滿地 `None` 地雷——但指標、slice、map、channel、interface 仍可能是 `nil`，後面會專講那些坑。

讀完你要能：寫出帶零值的 struct、看懂「為什麼不能把 `int32` 直接塞給 `int64`」、知道協定欄位為什麼要定死寬度。

## 先跟 Python 對一下

| Python | Go |
|--------|-----|
| 執行到才爆 `TypeError` | 編譯期型別錯誤，檔案都跑不起來 |
| `None` 到處可能出現 | 基本型別有零值；指標／slice／map／channel／interface 可 `nil` |
| if 條件對 truthy/falsy 很寬鬆 | 條件必須是真正的 `bool` |
| `int` 任意精度 | `int` 是固定寬度（跟平台有關）；超大數用 `math/big` |
| 型別註解可選 | 函式內常可推斷，但 API 邊界要寫清楚 |

寫錯會怎樣？在 Python 你可能「線上某條路徑才 TypeError」；在 Go 很多問題直接編不過——這是好事，但你要習慣「先讓編譯器開心」。

## 怎麼寫（能跑的最小例子）

```go
var n int          // 0
var s string       // ""
var ok bool        // false
name := "player"   // 短宣告，只能在函式內

// 沒有隱式數值轉換：寬度不同就要你親手轉
var x int32 = 1
var y int64 = int64(x)
```

試著拿掉 `int64(x)`、直接 `var y int64 = x`，讀編譯器怎麼罵你——這就是「契約」在說話。

再看 struct 零值：

```go
type Player struct {
	ID    int64
	Name  string
	Score int
}

func main() {
	var p Player // ID=0, Name="", Score=0
	fmt.Printf("%+v\n", p)
}
```

## 為什麼這樣設計／底層在幹嘛

### 零值哲學：沒設值 ≠ 不能用

結構體欄位若沒設定，會是各型別的零值。所以 `var room Room` 常常能直接用（例如計數器從 0 開始）。

但你必須問自己：**這個零值是不是「安全預設」？**

- 分數從 0 開始 → 通常 OK  
- 「房間狀態」若 0 剛好等於某個合法 phase → 要小心是不是意外進入錯誤階段  
- 指標／map 的零值是 `nil` → 還沒 `make` 就寫入會 panic  

「詳細」的重點：零值是設計工具，不是偷懶免初始化的藉口。

### 轉換必須明確

Python 有時靠運算默默推進型別；Go 要求你寫 `T(v)` 這種轉換，避免安靜的精度流失或語意變了你卻不知道。

什麼時候用？協定、檔案格式、資料庫欄位——任何「進出系統邊界」的地方，寬度與符號性都要清醒。

### 型別不只是貼紙

之後你還會碰到：方法集（method set）、介面是否滿足、值接收者還是指標接收者……都跟型別定義綁在一起。A、B 卷會展開；這裡先建立「型別會影響能不能呼叫」的直覺。

### 進階可先略過

- 底層型別（underlying type）與自定義型別之間的轉換規則。  
- `int` 在 32／64 位元平台的差異；對外協定通常偏好寫死 `int32`／`int64`。  

## 遊戲 Server 會用在哪

封包欄位、實體 ID、tick 序號，都需要**穩定、明確的寬度與有無符號**。

動態型別在協定邊界很容易變成「某天客戶端多送一種形狀，線上才炸」。靜態型別把問題前移到編譯期；再配合後面的 JSON／自訂協定章，你會越來越感謝這層契約。

## 請丟掉的舊習慣

1. 同一個變數先裝 `int` 再裝 `str`——Go 不允許，也別想繞。  
2. 依賴 truthy／falsy（空字串、空 list 算 False）——在 Go 要寫清楚：`len(s) == 0` 或 `s == ""`。  
3. 把 `None` 檢查文化原封搬來，卻忽略「零值常常已是合法狀態」——該檢查 `nil` 的地方查 `nil`，不該把每個 `0` 都當錯誤。  

## 動手練習

### 必做

1. 寫一個 `Player` struct（`ID int64`, `Name string`, `Score int`），印出零值實例。  
2. 試著把 `int32` 直接賦給 `int64` 變數，把編譯器錯誤訊息讀完（不要只看紅字閃過）。  

### 選做

1. 比較 Python `int` 與 Go `int64`：塞一個超大整數時各自會怎樣。  

## 常見坑

- **混用 `int` 與 `int64`**：讓編譯器教你；協定層可統一別名，例如 `type EntityID int64`。  
- **以為零值 struct「還沒初始化所以不能用」**：常常能用；真正要小心的是 `nil` 的那幾種型別。  
- **在 if 裡寫 `if score`**：必須是 `bool`，改成 `if score != 0` 之類。  

## 延伸閱讀

- <https://go.dev/ref/spec#Types>  
- <https://go.dev/blog/declaration-syntax>  
