---
lessonId: "D2"
title: "表驅動測試：加案例像填表"
description: "用 cases 表格 + t.Run 變成肌肉記憶；遊戲狀態機特別吃這套。"
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

## 這章你會搞懂什麼

**表驅動測試（table-driven tests）** 的精神是：把「輸入 → 期望輸出」列成一張表，迴圈跑；差異只在資料，不在測試骨架。

好處很實際：

- 加案例通常一行（或一個 struct）  
- 失敗時有名字（`t.Run`）好找  
- 逼你把行為想成「規格」，而不是「偶然跑過的路徑」

Go 社群極常用這招。遊戲的相位轉移、輸入校驗、滿房邏輯——簡直是為它量身做的。

## Python 對照

| Python | Go |
|--------|-----|
| `pytest.mark.parametrize` | `for _, tc := range tests { t.Run(...) }` |
| 一堆複製貼上的 `test_xxx` | 一張表 + 一個迴圈 |
| fixture 很大很重 | 表驅動偏單元；重依賴再搭配介面注入（B4） |

## 怎麼寫

```go
tests := []struct {
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
			t.Fatalf("got %d want %d", got, tc.want)
		}
	})
}
```

重點：

- `name` 寫人話（`full room rejects join`）  
- 失敗訊息包含 got／want  
- 錯誤路徑也要入表：`wantErr: ErrRoomFull` 之類  

範例：`examples/d02-table-test`。

## 細節

### 為什麼比「十個 TestXxx」好？

複製貼上會漂：有的 asserte 詳細、有的忘記測錯誤、改簽名要改十處。  
表格把規格集中；review 時眼睛掃資料就行。

### 怎麼測 error？

常見欄位：`wantErr error`，然後：

```go
if !errors.Is(err, tc.wantErr) {
	t.Fatalf("err=%v want %v", err, tc.wantErr)
}
```

這跟 B 卷哨兵直接接起來——API 可測性從設計那一刻就開始。

### 子測試與平行

`t.Run` 出的子測試可用 `t.Parallel()`（進階）。  
一開始先求正確與可讀；真要加速再平行，並注意共享狀態。

### 進階可先略過

- golden file（大輸出快照）  
- `cmp.Diff` 之類的可讀 diff 工具  

## 遊戲 Server 會用在哪

超適合入表的例子：

- `Join`：空房／滿房／錯誤階段  
- `Ready`：人數不足／重複 Ready  
- `PushInput`：非法方向／非 Playing  
- 相位：Lobby→Playing→Ended 合法與非法轉移  

這些測過，重構房間碼才敢下手。

## 請丟掉的舊習慣

1. 複製貼上十個 Test 只改數字。  
2. 只測 happy path。  
3. 斷言寫 `err != nil` 就結束——不確認是哪一種錯。

## 動手練習

### 必做

1. 跑 `examples/d02-table-test`。  
2. 加一個會失敗的 case，看失敗訊息；再修到綠。  

### 選做

1. 幫 `examples/b02-room-errors` 的 `Join` 補一張更完整的表（滿房、錯誤階段、成功）。  

## 常見坑

- **案例名都叫 `case1`**：失敗時完全不記得測了什麼。  
- **表格裡放會被突變的共享 slice／map**：小心測試互相污染。  
- **把整合測試硬塞進表驅動單元測試**：可以，但 setup 太重時考慮拆檔。

## 延伸閱讀

- <https://go.dev/wiki/TableDrivenTests>  
