---
lessonId: "C5"
title: "Mutex：沒有 GIL，共享就要同步"
description: "鎖的粒度、defer Unlock、RWMutex；以及何時該改用「單一擁有者 + channel」。"
volume: "c"
order: 5
level: "l2"
status: "ready"
path_required: true
tags: ["sync", "mutex"]
example: "examples/c05-mutex"
prev: "C4"
next: "C6"
---

## 這章你會搞懂什麼

多個 goroutine 碰同一塊記憶體（例如同一個 `map`），就需要同步。最常見是 `sync.Mutex`（互斥鎖）：同一時間只有一個持有者能進臨界區。

關鍵不是「會不會 Lock」，而是：

- **粒度**：鎖太大 → 變慢；太小 → 仍 race 或難懂  
- **持鎖時別做慢 IO**  
- **解鎖一定要發生**（所以常 `defer Unlock`）  

Go **沒有 Python 那種 GIL 給你假安全感**。data race 是真 bug。下一章 C6 會用偵測器釘死這件事。

另一條路：不要共享——讓**單一 goroutine 擁有資料**，外人只透過 channel 送請求。兩種都合理，選可讀的那個。

## Python 對照

| Python | Go |
|--------|-----|
| `threading.Lock` | `sync.Mutex` |
| `RLock`／讀寫鎖配方 | `sync.RWMutex` |
| GIL 下「偶發沒事」 | **沒有 GIL**；並發寫 map 會炸或被 `-race` 抓 |
| 常靠運氣／單執行緒 asyncio | 必須顯式同步或避免共享 |

「我在 CPython 對 dict 這樣寫都沒事」——到 Go 請直接假設：**會出事**。

## 怎麼寫

```go
type Registry struct {
	mu sync.Mutex
	m  map[string]int
}

func (r *Registry) Inc(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.m[id]++
}
```

`defer Unlock` 的好處：中間有 `return`／錯誤路徑也不容易漏解鎖。

範例：`examples/c05-mutex`。務必加 `-race` 跑。

## 細節

### 鎖保護的是「不變量」，不是「那一行字」

例如 map 的「讀改寫」要在同一段臨界區完成。  
只鎖 `m[id] = x` 卻在鎖外先讀舊值——仍可能 race。

### 廣播時的經典手法

```text
鎖內：把 peer 清單拷貝出來
解鎖
對拷貝發送（可能很慢、可能阻塞）
```

為什麼？因為持鎖時做網路發送，別的人都卡住，延遲會互相傳染，也更容易死鎖。

### `RWMutex`

多讀少寫時，讀鎖可並存。但別預設「RW 一定比較快」——要測量；錯誤使用（讀鎖升級寫鎖等）也容易死鎖。

### 絕對別複製含 Mutex 的 struct

複製會複製鎖的狀態，變成「各解各的鎖」——非常難查。傳指標，或讓鎖欄位不被複製（設計上別值傳整個含鎖物件）。

### 進階可先略過

- 鎖順序：多把鎖時全局約定順序，避免 AB/BA 死鎖。  
- `sync.Map` 有特定場景；不是一般 map 的免費加速器。

## 遊戲 Server 會用在哪

| 資料 | 常見策略 |
|------|----------|
| 大廳人數、房間列表 | Mutex，或單一 lobby owner |
| 單一房間狀態 | **最好**單一 room goroutine 擁有；外人丟命令 |
| Session 索引 | Mutex／分片（進階） |

房間內規則邏輯能單線跑就單線跑——比到處加鎖簡單太多。

## 請丟掉的舊習慣

1. 「有 GIL／我測很多次沒炸」當正確性證明。  
2. 持鎖時 call 網路、DB、別人的回呼（回呼可能再搶同一把鎖）。  
3. 用 `sleep` 降低撞車機率來「修」race。

## 動手練習

### 必做

1. 跑 `examples/c05-mutex`，加上 `go test -race`。  
2. 寫一版無鎖並發寫 map，看 `-race` 怎麼報。  

### 選做

1. 把同一份資料改成「owner goroutine + 請求 channel」，比較可讀性。  

## 常見坑

- **鎖住了但保護範圍不夠**：資料結構的整次更新要在臨界區。  
- **死鎖**：同 goroutine 重複 Lock 普通 Mutex；或鎖順序交錯。  
- **以為讀不用鎖**：並發讀寫 map 在 Go 非法；純讀多執行緒也要有 happens-before 保證。

## 延伸閱讀

- <https://pkg.go.dev/sync>  
