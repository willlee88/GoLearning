---
lessonId: "D4"
title: "Benchmark 與 pprof：先測量，再談優化"
description: "用 testing.B 看 ns／op 與記憶體分配；用 pprof 找 CPU／heap 熱點——別靠感覺改碼。"
volume: "d"
order: 4
level: "l2"
status: "ready"
path_required: true
tags: ["benchmark", "pprof"]
example: "examples/d04-bench"
prev: "D3"
next: "E0"
---

## 這章你會搞懂什麼

優化順序只有一句話：

> **正確 → 測量 → 改熱點 → 再測量**

「我覺得這樣比較快」在遊戲 Server 裡非常容易錯：分配、鎖、編碼、廣播放大，常跟直覺打架。

Go 內建：

- **`testing.B`**：微基準（benchmark）  
- **pprof**：看 CPU／記憶體等剖面  

這章帶你進門：會跑、會看 `ns/op` 與 `B/op`，知道熱點大概在哪。

## Python 對照

| Python | Go |
|--------|-----|
| `timeit` | `go test -bench` |
| `cProfile`／`py-spy` | `pprof`（CPU／heap 等） |
| 常先改寫法再猜 | 文化上更強調先量 |

## 怎麼寫

```go
func BenchmarkSnap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = snapshot()
	}
}
```

```bash
go test -bench=. -benchmem
go test -cpuprofile=cpu.out -bench=.
go tool pprof cpu.out
```

`-benchmem` 會顯示分配相關指標（例如 `B/op`、`allocs/op`）——遊戲熱路徑很愛看這個。

範例：`examples/d04-bench`。

## 細節

### `b.N` 是什麼？

測試框架會調整迴圈次數，讓測量穩定到一定時間。  
你只要把「要量的工作」放進迴圈；不要自己傻 sleep。

### 常見噪音

- 基準裡含網路／磁碟：結果會抖，也難解釋  
- 第一次分配／全域 init 混進測量：必要時先在迴圈外暖機  
- 編譯器把沒用到的結果優化掉：用 `_ =` 或累加避免「空轉被消掉」

### pprof 看什麼？

先找：**累積時間最多的函式**、**分配最多的地方**。  
遊戲裡常出現的熱點：序列化、廣播迴圈、不斷 `make`／字串拼接、鎖競爭。

### 進階可先略過

- `benchstat` 比較兩次基準是否真有顯著差異  
- mutex／block profile  
- 火焰圖閱讀細節  

## 遊戲 Server 會用在哪

優先量這些，而不是先糾結微語法：

- 狀態 snapshot／marshal  
- 廣播給 N 個連線  
- 每 tick 的分配是否隨人數線性爆掉  

人數一多，分配跟扇出比「某個 if 少寫一行」重要太多。

## 請丟掉的舊習慣

1. 無數據改「應該比較快」的寫法。  
2. 正確性還沒測完就開始微優化。  
3. 只看平均、不看分配與尖刺（結合負載測試是後面 J 卷）。

## 動手練習

### 必做

1. 跑 `examples/d04-bench`，讀懂 `ns/op`、`B/op`、`allocs/op`。  
2. 故意讓被測函式多分配一點（例如每圈 `make`），再比一次數字。  

### 選做

1. 產出 `cpu.out`，用 `go tool pprof` 的 `top` 看前幾名。  

## 常見坑

- **基準寫錯，量到的是 fmt 印字或測試框架本身**。  
- **一次改十處再量**：不知道哪改有效。  
- **把 debug 日誌留在熱路徑**：量到的是 IO。

## 延伸閱讀

- <https://go.dev/doc/diagnostics>  
- <https://pkg.go.dev/runtime/pprof>  

## D 卷檢查點

你現在應該能：

1. 讀懂並維護 `go.mod`／`go.sum`  
2. 用表驅動鎖住房間規則  
3. 對解析路徑跑 fuzz  
4. 用 bench／pprof 回答「慢在哪」而不是「我覺得」  

下一卷 E 會把標準庫常用零件補齊；網路與遊戲房間則在 F／H。  
