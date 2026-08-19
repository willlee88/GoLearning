---
lessonId: "E8"
title: "sync 套件地圖"
description: "Mutex、WaitGroup、Once、Map、Pool 何時用。"
volume: "e"
order: 8
level: "l2"
status: "ready"
path_required: false
tags: ["sync", "stdlib"]
example: ""
prev: "E7"
next: "F0"
---

## 這章你會搞懂什麼

C 卷教過：多個 goroutine 一起碰同一塊記憶體，要同步。`sync` 套件就是這類**低階同步工具箱**。

本章當**地圖**：每個型別適合什麼場景、不適合什麼。重點不是背完 API，而是**選錯工具的後果**——例如無腦上 `sync.Map`，或複製含 Mutex 的 struct。

## 先跟 Python 對一下

| Python | Go | 注意 |
|--------|-----|------|
| `threading.Lock` | `sync.Mutex` | 觀念接近；Go 還有 `RWMutex` |
| 自己計數 + Event | `sync.WaitGroup` | 等一組 goroutine 結束很常見 |
| 沒有內建 Once 時自己寫 flag | `sync.Once` | 保證函式只執行一次 |
| 以為 `dict` 加鎖就好 | 通常 `map` + `Mutex` | `sync.Map` 是特例，不是預設 |

人話：

- **互斥鎖（Mutex）**：同一時間只准一個人進臨界區。  
- **WaitGroup**：像「還有幾個人沒回家」，歸零才開門。  
- **Once**：初始化只做一遍，之後呼叫都直接跳過。

## 地圖

| 型別 | 用途（人話） | 什麼時候想一想 |
|------|----------------|------------------|
| `Mutex`／`RWMutex` | 保護共用資料 | 讀很多寫很少可試 RWMutex；別重入（Go Mutex 不能鎖兩次） |
| `WaitGroup` | 等一組 worker 結束 | `Add` 要在 `Wait` 之前；`Add` 計數別弄錯 |
| `Once` | 懶初始化／只啟動一次 | `Do` 裡 panic 的語意要看文件 |
| `Map` | 特殊併發 map | key 穩定、寫少讀多或 key 不相交；**一般業務優先 map+Mutex** |
| `Pool` | 暫存物件重用，降 GC | `Get` 可能拿到髒物件，用前 Reset；別當快取 |
| `Cond` | 等條件成真 | 大多數情況 channel 更清晰；少用 |

## 怎麼寫（能跑的最小例子）

### Once

```go
var once sync.Once

func initAll() { /* 載入設定、註冊度量… */ }

func ensureInit() {
	once.Do(initAll)
}
```

### Mutex 保護 map（最常見）

```go
type Hub struct {
	mu    sync.Mutex
	rooms map[string]*Room
}

func (h *Hub) Get(id string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.rooms[id]
}
```

### WaitGroup

```go
var wg sync.WaitGroup
for i := 0; i < 3; i++ {
	wg.Add(1)
	go func() {
		defer wg.Done()
		work()
	}()
}
wg.Wait()
```

## 為什麼這樣設計／底層在幹嘛

1. **`sync` 是「共享記憶體＋鎖」路線**  
   Go 還有 channel 的「用通訊分享記憶體」。能單一 owner（一個 goroutine 管狀態）就少鎖；真要共享再鎖。

2. **為什麼不要複製含 Mutex 的 struct**  
   Mutex 複製會複製「鎖的狀態」，變成各看各的鎖，等於沒鎖。傳指標，或別把 Mutex 當值到處傳。

3. **`sync.Map` 不是 concurrent dict 萬靈丹**  
   文件寫了適用情境。API 也比較醜（存 `any`）。Hub 的 rooms：房間數不多、操作複雜，`map+Mutex` 通常更單純。

4. **`Pool` 不是快取**  
   runtime 可能隨時丟掉 Pool 裡的東西。適合「短命緩衝區」重用；Get 出來要當可能是任意舊物件。

## 遊戲 Server 會用在哪

- Hub：`rooms map` + `Mutex`。  
- 連線表、玩家 session 表：同樣模式。  
- `sync.Once`：啟動全域 ticker、註冊 metrics 只做一次。  
- 關服時 `WaitGroup` 等各房間／連線收尾（或改 errgroup，見 C7）。  
- 寫 WebSocket：多 goroutine 同寫一條連線 → 需要寫鎖（F6）。

## 請丟掉的舊習慣

1. **無腦 `sync.Map` 當「Python 的 concurrent dict」**。  
2. **鎖的粒度亂到整個世界一把大鎖**——或反過來鎖太碎死鎖。先粗後量。  
3. **在持鎖時做遠端 I/O**：鎖會被握很久，別的 goroutine 全堵。  
4. **以為有 GIL 心智就安全**——Go 沒有那種「單執行緒 bytecode」護身符，race 要用 `-race` 抓。

## 動手練習

### 必做

1. 用自己的話說明：Arena／教學 Hub 為什麼通常不用 `sync.Map`。  
2. 寫一個「計數器」：10 個 goroutine 各加 1000 次，比較「無鎖（故意錯）」與「Mutex」；有空加 `go test -race`。  

### 選做

1. 讀 `sync.Map` 文件的「何時適合」段落，對照你的 rooms 表是否符合。  

## 常見坑

- **`Unlock` 忘了**：用 `defer`。  
- **`WaitGroup` 複製傳進函式**：跟 Mutex 一樣，傳指標。  
- **`Add` 在 goroutine 裡面才呼叫**：可能 `Wait` 先看到 0。  
- **死鎖**：同 goroutine 鎖兩次、或鎖順序 A→B／B→A 交錯。  

## 延伸閱讀

- <https://pkg.go.dev/sync>  

## 接回主路徑

標準庫地圖先到這裡。若尚未讀網路卷 → **F0**。
