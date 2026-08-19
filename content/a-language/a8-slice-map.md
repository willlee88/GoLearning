---
lessonId: "A8"
title: "slice 與 map（本卷超重要檢查點）"
description: "slice 是底層陣列的視圖；map 是雜湊表引用。共享、append、併發寫入，都是實戰雷區。"
volume: "a"
order: 8
level: "l2"
status: "ready"
path_required: true
tags: ["slice", "map", "game-server"]
example: "examples/a08-player-registry"
prev: "A7"
next: "A9"
---

## 這章你會搞懂什麼

**slice（切片）** 是指向底層陣列的小描述符：指標、長度（len）、容量（cap）。  
**map** 比較像「雜湊表的引用」。

兩者是 Go 日常資料結構，也是 bug 熱區：誤判有沒有共享、`append` 何時重分配、多 goroutine 一起寫 map。  
這章要深刻到你能**在紙上畫出記憶體關係**——這是 A 卷檢查點。

## Python 對照

| Python | Go |
|--------|-----|
| `list` | slice（更像「陣列視圖」） |
| `dict` | map |
| 切片賦值／可變共享 | slice 頭會拷貝，但常共享底層陣列 |
| dict 非執行緒安全（有時被 GIL 掩蓋） | map **明確**非併發安全 |

## 怎麼寫

```go
s := []int{1, 2, 3}
s = append(s, 4)
fmt.Println(len(s), cap(s))

m := map[string]int{"a": 1}
m["b"] = 2
v, ok := m["c"] // ok == false，表示沒這個 key
delete(m, "a")
```

```go
// 正確建立可寫的 map
m := make(map[PlayerID]*Player)
// var m map[PlayerID]*Player // nil map：讀通常 OK，寫會 panic
```

檢查點範例：`examples/a08-player-registry`。

## 細節

### slice 三元組

| 欄位 | 意義 |
|------|------|
| ptr | 底層陣列從哪裡開始看 |
| len | 你看得到多長 |
| cap | 從起點算，底層還能長到哪 |

```go
a := []int{1, 2, 3, 4}
b := a[1:3] // 共享底層；改 b[0] 會動到 a[1]
```

所以 `b := a[1:3]` **不是深拷貝**。若要獨立副本，請 `copy` 到新 slice。

### `append` 何時「切斷共享」

容量不夠時，`append` 會**配新陣列並拷貝**，之後就可能不再跟舊底層共享。  
規則實用版：

- 要獨立：自己 `make` + `copy`  
- 知道大概長度：`make([]T, 0, n)` 預留 cap，少配幾次記憶體  

### map 必記三件事

1. 迭代順序**刻意隨機**——別依賴「插入順序」。  
2. 鍵必須可比較（comparable）。  
3. **併發讀寫 = data race，甚至直接 panic**：加鎖、分片，或讓單一 goroutine 擁有那份 map。

### 進階可先略過

- slice 成長策略與短暫浪費的容量。  
- `slices`／`maps` 標準庫輔助（Go 1.21+）。

## 遊戲 Server 會用在哪

典型玩家登記：

```go
type Registry struct {
	mu    sync.RWMutex
	byID  map[PlayerID]*Player
	order []PlayerID // 若你需要穩定順序
}
```

廣播時常見手法：鎖內先複製 ID 列表或做 snapshot，**鎖外**再寫網路。持鎖做 I/O 會把延遲與死鎖風險一起拉高。

`Snapshot()` 為什麼常要拷貝？——避免呼叫端還在讀，你房間邏輯又在改同一塊底層資料。

## 請丟掉的舊習慣

1. 以為 `b = a[1:3]` 是深拷貝。  
2. 多執行緒／多 goroutine 無鎖共用一個 dict／map。  
3. 只用 list 線性掃當唯一索引——房間查人請用 map。

## 動手練習

### 必做（A 卷檢查點）

1. 完成並理解 `examples/a08-player-registry`：  
   - `go test .`  
   - `go run ./cmd/demo`  
   - 新增／刪除玩家、依 ID 查詢、JSON 快照  
2. 用註解或筆記說明：為什麼 `Snapshot()` 要拷貝。  

### 選做

1. 寫一個「append 共享踩坑」最小重現，再修掉。  
2. 單獨做一個故意錯的並發 map 寫入，用 `-race` 看它爆炸（別當成正確範例提交）。  

## 常見坑

- **nil map 寫入 panic**：記得 `make` 或用字面量初始化。  
- **小 slice 卡住大陣列**：殘留引用可能讓 GC 收不回記憶體——必要時 copy 出來。  
- **range 時刪 key**：語意允許但要懂在幹嘛；不熟就先收集再刪。  
- **併發 map**：不是「偶爾怪怪的」，是明確不允許。
