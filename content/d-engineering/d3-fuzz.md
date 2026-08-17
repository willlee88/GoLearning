---
lessonId: "D3"
title: "Fuzz 測試入門"
description: "亂餵輸入找 panic 與錯誤假設。"
volume: "d"
order: 3
level: "l2"
status: "ready"
path_required: true
tags: ["fuzz"]
example: "examples/d03-fuzz"
prev: "D2"
next: "D4"
---

## 本章你會建立的心智模型

Fuzzing 自動產生輸入，專治「我以為不會發生」的封包與字串。協定 parse、input 解析很適合。

## Python 對照

hypothesis 類似；Go 1.18+ 內建 fuzz。

## L1 能用

```go
func FuzzParse(f *testing.F) {
	f.Add("1,0")
	f.Fuzz(func(t *testing.T, s string) {
		_, _, _ = parseInput(s) // 不該 panic
	})
}
```

```bash
go test -fuzz=FuzzParse -fuzztime=5s
```

範例：`examples/d03-fuzz`。

## 請丟掉的 Python 習慣

1. 只測 happy path。  

## 遊戲 Server 連結

對 `dx,dy` 字串與 JSON envelope fuzz。

## 練習

### 必做

1. 跑 5 秒 fuzz。  

## 延伸閱讀

- Go fuzzing tutorial  
