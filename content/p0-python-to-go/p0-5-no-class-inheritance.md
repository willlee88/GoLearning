---
lessonId: "P0.5"
title: "沒有 class 繼承？用 struct、方法、介面組合"
description: "Go 不靠繼承樹。資料放 struct、行為綁方法、能力用 interface 描述；編譯期檢查你有沒有做到。"
volume: "p0"
order: 5
level: "l2"
status: "ready"
path_required: true
tags: ["python-bridge", "interfaces"]
example: ""
prev: "P0.4"
next: "P0.6"
---

## 這章你會搞懂什麼

Go **沒有** class，也**沒有**那種一層疊一層的繼承樹。

你改用三件事組合行為：

1. **struct** — 放資料  
2. **method（方法）** — 把行為綁在型別上  
3. **interface（介面）** — 描述「我需要什麼能力」  

型別只要擁有介面要求的方法，就**自動滿足**介面（結構化型別），不必寫 `implements`。讀完你要能寫出一個小介面、一個實作、一個只依賴介面的函式。

## 先跟 Python 對一下

| Python | Go |
|--------|-----|
| `class` + 繼承 | struct + 嵌入（embedding） |
| duck typing（執行期才知道像不像鴨子） | 方法集滿足 interface（**編譯期**就檢查） |
| ABC / Protocol | interface |
| `self` | receiver（接收者：值或指標） |

Python 的鴨子類型很自由，錯了常執行期才爆。Go 想保留「用能力說話」的味道，但把檢查拉到編譯期——少驚喜，也少線上爆炸。

## 怎麼寫（能跑的最小例子）

```go
type Player struct {
	Name string
	HP   int
}

func (p *Player) Damage(n int) {
	p.HP -= n
}

type Health interface {
	Damage(n int)
}
```

任何有 `Damage(int)` 方法的型別，都能當成 `Health` 來用——不用登記、不用繼承。

再來一個「只依賴介面」的函式：

```go
type Notifier interface {
	Notify(msg string)
}

type ConsoleNotifier struct{}

func (ConsoleNotifier) Notify(msg string) {
	fmt.Println(msg)
}

func Broadcast(n Notifier, msg string) {
	n.Notify(msg)
}
```

測試時你可以換一個假的 `Notifier`，不必連真實網路——這就是介面的甜頭。

## 為什麼這樣設計／底層在幹嘛

### 接受介面，回傳具體型別

常見 API 風格：函式**參數**用小介面（我需要什麼能力），**回傳**具體 struct（呼叫者拿到的是清楚的東西）。這樣依賴方向乾淨，也比較好測。

### 小介面比較好組合

1～3 個方法的介面，比「上帝介面」好用。介面太大 → 假實作很難寫、每個實作都臃腫、之後改一個方法全家炸。

寫錯會怎樣？一開始就為每個 struct 先抽介面 → 過度抽象，目錄裡全是只有一個實作的 interface，讀起來更累。等你有**第二個實作**或**測試需要替換**時再抽，通常比較甜。

### 嵌入（embedding）不是繼承 2.0

你可以嵌入 struct，外層會「提升」內層的方法，像委派。但它**不是**完整的子型別多型替代方案——要有意識使用，別暗中做成深繼承。

### 進階可先略過

- 介面值裡有「動態型別 + 動態值」；把 `nil` 指標放進 interface 後，`== nil` 可能不如你想（後面專章會打臉）。  
- 方法集：值接收者 vs 指標接收者，會影響「算不算實作了某個介面」。  

## 遊戲 Server 會用在哪

這些都適合作**小介面**，讓規則邏輯可單測、可替換：

- `Transport` — 怎麼送訊息（真實連線 vs 測試假傳輸）  
- `RoomStore` — 房間怎麼存  
- `Clock` — 可測時間（不要在規則裡直接 `time.Now()` 散落一地）  

規則套件依賴介面，不依賴「某個具體 WebSocket 函式庫」，之後換傳輸或寫單元測試才不會痛。

## 請丟掉的舊習慣

1. 一上來設計三層繼承：`BaseEntity → MovingEntity → Player`。  
2. 為了「以後可能會擴充」做巨大 base class——以後常常沒來，複雜度先來了。  
3. 把 interface 當成裝飾用註解，卻不在 API 邊界使用——那就白抽了。  

## 動手練習

### 必做

1. 定義 `type Notifier interface { Notify(msg string) }`，實作一個 `ConsoleNotifier`。  
2. 寫 `func Broadcast(n Notifier, msg string)`，函式本體只依賴介面。  

### 選做

1. 用嵌入做 `type Admin struct { Player; Level int }`，觀察 `Player` 的方法是否被提升。  

## 常見坑

- **每個 struct 先寫 interface**：過度抽象。等有第二個實作或測試替身再抽。  
- **介面方法太多**：假物件寫不完；拆小一點。  
- **為了「像 OOP」硬套繼承心智**：Go 的組合路徑不同，硬套會寫得很彆扭。  

## 延伸閱讀

- <https://go.dev/doc/effective_go#interfaces>  
