---
lessonId: "C5"
title: "Mutex 與共享狀態"
description: "鎖的粒度、RWMutex、與 channel 的取捨。"
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

## 本章你會建立的心智模型

當多個 goroutine 共享記憶體，需要同步：`sync.Mutex` / `RWMutex`。鎖的關鍵是**粒度**與**持鎖時不做 I/O**。房間 registry 常用 Mutex；也可改「單一 owner goroutine + inbox」避免共享。

## Python 對照

| Python | Go |
|--------|-----|
| `threading.Lock` | `sync.Mutex` |
| GIL 誤安全感 | **沒有 GIL**；race 真實存在 |

## L1 能用

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

範例：`examples/c05-mutex`。

## L2 機制

- `defer Unlock` 防早退。  
- `RWMutex`：多讀少寫。  
- 死鎖：鎖順序一致。  
- 複製含 Mutex 的 struct 是 bug。  

## 請丟掉的 Python 習慣

1. 「CPython 有 GIL 所以 dict 沒事」。  
2. 持鎖時 call 網路。  

## 遊戲 Server 連結

廣播：鎖內拷貝 peer 列表 → 解鎖 → 對拷貝發送。

## 練習

### 必做

1. 跑 `examples/c05-mutex` 含 `-race`。  
2. 寫一版無鎖並發 map 寫入，看 `-race` 報錯。  

## 延伸閱讀

- <https://pkg.go.dev/sync>  
