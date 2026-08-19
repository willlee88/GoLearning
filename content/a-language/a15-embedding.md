---
lessonId: "A15"
title: "嵌入：少寫轉發，但別當成繼承"
description: "匿名嵌入會提升方法與欄位；適合 has-a，不適合深層 is-a 樹。可選深讀。"
volume: "a"
order: 15
level: "l2"
status: "ready"
path_required: false
tags: ["struct", "composition"]
example: "examples/a15-embedding"
prev: "A14"
next: "A16"
---

## 這章你會搞懂什麼

**嵌入（embedding）** 讓外層型別自動「提升」內層的方法與欄位，少寫一堆轉發函式。它是 Go 的組合工具，**不是繼承**。

適合：「我有一個 Logger／Mutex，想少打點字」。  
不適合：「我要蓋三層 BaseEntity 幻想 OOP 王國」。

這章是 A 卷加廣，主路徑做完 A14 再回來也行。

## Python 對照

| Python | Go |
|--------|-----|
| 多重繼承／mixin | 可嵌入多個 struct（請謹慎） |
| 繼承後覆寫 | 外層定義同名方法會遮罩內層 |

## 怎麼寫

```go
type Logger struct{}

func (Logger) Log(msg string) { fmt.Println(msg) }

type Server struct {
	Logger // 匿名嵌入
	Addr   string
}

s := Server{Addr: ":8080"}
s.Log("up") // 方法被提升，可直接叫
```

範例：`examples/a15-embedding`。

## 細節

### 提升時，receiver 仍是內層

你呼叫 `s.Log`，真正的接收者是裡面的 `Logger`。這對理解「狀態改到誰身上」很重要。

### JSON 形狀也會被影響

嵌入的欄位可能被提升到外層 JSON。先寫測試釘住輸出，別靠感覺。

### 名字衝突

兩邊都有同名方法／欄位時，要顯式走 `s.Logger.Log(...)`，或在外層定義同名方法遮罩。

### 嵌入 `sync.Mutex` 的常見寫法

很常見，但記住：

- 盡量**不要匯出**「含鎖且可能被值拷貝」的型別  
- 拷貝含 Mutex 的 struct = 拷貝鎖狀態，後面超難查的 race／死鎖味道

### 進階可先略過

- 外層因嵌入而「順便」滿足某個介面：有時方便，有時 API 意外變大。

## 遊戲 Server 會用在哪

很多時候你要的是**持有（has）**，不是嵌入：

```go
type roomRuntime struct {
	game *game.Room
	// conn map、其他執行期狀態…
}
```

規則物件保持獨立，通常比「房間嵌入一切」清楚。嵌入留給真正想提升的小能力（例如內部 helper），不要為嵌而嵌。

## 請丟掉的舊習慣

1. 為了「以後擴充」先嵌三層 Base。  
2. 把嵌入當子型別多型——不能把 `Server` 當成 `Logger` 傳進只要 `Logger` 的函式（除非它因方法集而滿足介面，那是另一回事）。  
3. 匯出含鎖 struct 又到處值傳遞。

## 動手練習

### 必做

1. 跑 `examples/a15-embedding`。  
2. 在外層覆寫 `Log`，觀察呼叫行為。  

### 選做

1. 試著嵌入兩個都有同名方法的型別，看衝突時怎麼寫才編譯。  

## 常見坑

- **以為嵌入 = 繼承多型**。  
- **JSON／方法集意外變大**：API 表面積失控。  
- **拷貝嵌了 Mutex 的值**：鎖語意崩潰。
