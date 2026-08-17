---
lessonId: "A8"
title: "slice 與 map"
description: "len/cap、底層陣列共享、append 陷阱；map 語意與非併發安全。"
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


## 本章你會建立的心智模型

**slice** 是指向底層陣列的描述符（指標、len、cap）；**map** 是雜湊表引用。兩者都是 Go 日常資料結構，也是 bug 熱區：誤判共享、append 重分配、並發寫 map。這一章要「深刻」到你能畫出記憶體關係。

## Python 對照

| Python | Go |
|--------|-----|
| `list` | slice（更貼近「陣列視圖」） |
| `dict` | map |
| list 賦值共享 + 可變 | slice 頭拷貝但可共享底層陣列 |
| dict 非執行緒安全（GIL 掩蓋） | map **明確**非併發安全 |

## L1 能用

```go
s := []int{1, 2, 3}
s = append(s, 4)
fmt.Println(len(s), cap(s))

m := map[string]int{"a": 1}
m["b"] = 2
v, ok := m["c"] // ok == false
delete(m, "a")
```

```go
// 正確建立
m := make(map[PlayerID]*Player)
// var m map[PlayerID]*Player // nil map，讀 OK、寫 panic
```

範例（檢查點）：`examples/a08-player-registry`。

## L2 機制

### slice 三元組

| 欄位 | 意義 |
|------|------|
| ptr | 底層陣列起點 |
| len | 可見長度 |
| cap | 從起點到陣列末的容量 |

```go
a := []int{1, 2, 3, 4}
b := a[1:3] // 共享底層；改 b[0] 影響 a[1]
```

### append

容量不足時**分配新陣列**並拷貝；可能切斷與舊底層的共享。規則：

- 若你需要獨立拷貝：`copy` 到新 slice。  
- 預分配：`make([]T, 0, n)` 減少分配。  

### map

- 迭代順序**隨機**（刻意）。  
- 鍵必須可比較。  
- **併發讀寫 = data race / panic**：加鎖或分片，或改由單一 goroutine 擁有。  

## L3 深潛（可選）

- slice grow 策略與記憶體浪費。  
- map 在刪除與增長時的桶行為（概念級）。  
- `slices` / `maps` 標準庫輔助套件（Go 1.21+）。

## 請丟掉的 Python 習慣

1. 以為 `b = a[1:3]` 是深拷貝。  
2. 在多執行緒／多 goroutine 無鎖共用一個 dict/map。  
3. 用 `list` 線性掃當唯一索引——房間查找需要 map。

## 遊戲 Server 連結

典型結構：

```go
type Registry struct {
	mu      sync.RWMutex
	byID    map[PlayerID]*Player
	order   []PlayerID // 若需要穩定順序
}
```

廣播時：先在鎖內複製 id 列表或 snapshot，再鎖外寫網路，避免持鎖做 I/O。

## 練習

### 必做（A 卷檢查點）

1. 完成並理解 `examples/a08-player-registry`：  
   - `go test .`  
   - `go run ./cmd/demo`  
   - 新增／刪除玩家、依 ID 查詢、JSON 快照  
2. 書面或註解說明：為何 `Snapshot()` 要拷貝。  

### 選做

1. 寫出一個「append 共享踩坑」的最小重現，再修掉。  
2. 用 `-race` 跑一個**故意錯誤**的並發 map 寫入 demo（獨立檔，勿提交為正確範例）。  

## 常見坑與如何看見

- **nil map 寫入 panic**：`make` 或字面量初始化。  
- **range 時刪 key**：允許但要懂語意。  
- **slice 洩漏**：大陣列被小 slice 引用導致無法回收——必要時 copy。  

## 延伸閱讀

- <https://go.dev/blog/slices-intro>  
- <https://go.dev/blog/maps>  
