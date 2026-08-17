---
lessonId: "D4"
title: "Benchmark 與 pprof 入門"
description: "先測量再優化。"
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

## 本章你會建立的心智模型

感覺不可靠。`testing.B` 做微基準；pprof 看 CPU/heap。優化順序：正確 → 測量 → 改熱點。

## Python 對照

| Python | Go |
|--------|-----|
| timeit / cProfile | bench / pprof |

## L1 能用

```go
func BenchmarkSnap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = snapshot()
	}
}
```

```bash
go test -bench=. -benchmem
go test -cpuprofile=cpu.out
go tool pprof cpu.out
```

範例：`examples/d04-bench`。

## 請丟掉的 Python 習慣

1. 無數據改「應該比較快」的寫法。  

## 遊戲 Server 連結

Marshal state、廣播迴圈是常見熱點。

## 練習

### 必做

1. 跑 bench 看 ns/op 與 B/op。  

## 延伸閱讀

- Diagnostics doc  
