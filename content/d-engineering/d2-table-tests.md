---
lessonId: "D2"
title: "表驅動測試"
description: "table-driven tests 成為肌肉記憶。"
volume: "d"
order: 2
level: "l2"
status: "ready"
path_required: true
tags: ["testing"]
example: "examples/d02-table-test"
prev: "D1"
next: "D3"
---

## 本章你會建立的心智模型

表驅動：一組 cases，迴圈跑，差異只在輸入輸出。易加案例、報告清楚。

## Python 對照

| Python | Go |
|--------|-----|
| pytest parametrize | `for _, tc := range tests` |

## L1 能用

```go
tests := []struct{
	name string
	in   int
	want int
}{
	{"zero", 0, 0},
	{"pos", 2, 4},
}
for _, tc := range tests {
	t.Run(tc.name, func(t *testing.T) {
		if got := double(tc.in); got != tc.want {
			t.Fatalf("got %d", got)
		}
	})
}
```

範例：`examples/d02-table-test`。

## 請丟掉的 Python 習慣

1. 複製貼上十個 TestXxx 只改數字。  

## 遊戲 Server 連結

`Ready`/`PushInput`/`phase` 轉移都適合表驅動。

## 練習

### 必做

1. 跑範例；再加一個 failing case 再修。  

## 延伸閱讀

- Go Wiki: TableDrivenTests  
