---
lessonId: "D3"
title: "Fuzz：亂餵輸入，抓你以為不會發生的事"
description: "Go 內建 fuzz：給種子、自動長出怪字串，專治封包解析與邊界假設。"
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

## 這章你會搞懂什麼

單元測試很會測「你想得到的案例」。  
**Fuzz（模糊測試）** 則是：引擎自動長出大量輸入，專找「我以為不可能」的崩潰與錯誤假設。

對遊戲 Server 超對味的目標：

- 解析玩家輸入字串  
- JSON envelope  
- 長度／符號／奇怪 Unicode  

Go 1.18+ 工具鏈內建 fuzz，不必先上很重的框架。目標常常是：**不要 panic、不要無限迴圈、錯誤要乾淨**。

## Python 對照

| Python | Go |
|--------|-----|
| Hypothesis 等庫 | `testing.F`／`FuzzXxx` 內建 |
| 手動想邊界案例 | 種子（seed）+ 引擎突變 |
| 多半抓例外 | 抓 panic、錯誤不變量、錯誤的「成功」 |

## 怎麼寫

```go
func FuzzParse(f *testing.F) {
	f.Add("1,0") // 種子：合理例子先放
	f.Fuzz(func(t *testing.T, s string) {
		_, _, _ = parseInput(s) // 重點：不該 panic
	})
}
```

跑一小段時間：

```bash
go test -fuzz=FuzzParse -fuzztime=5s
```

範例：`examples/d03-fuzz`。

## 細節

### 種子為什麼重要？

引擎從種子附近突變，比較容易碰到「像真的、但又怪」的輸入。  
只丟空測試、不給種子，也能跑，但效率常較差。把真實封包樣例當種子很賺。

### 你在断言什麼？

常見不變量：

- 任何輸入都**不 panic**  
- 失敗就回 error，不要半改全域狀態  
- 若回成功，輸出必須滿足範圍（例如 dx∈{-1,0,1}）  

「解析失敗」本身往往是 OK——客戶端本來就會亂送。可怕的是解析器炸掉整服（接 B3 的教訓）。

### 失敗會留下 corpus

Fuzz 找到的問題輸入可被保存，之後當成回歸案例——這點非常香。修好後再跑，應維持綠。

### 進階可先略過

- 多參數 fuzz、自訂類型  
- 與 CI 的時間預算（例如 PR 跑 15s、夜跑更久）

## 遊戲 Server 會用在哪

- `dx,dy` 字串  
- JSON envelope 的 type／payload  
- ticket／token 字串  

凡是「從網路線進來的位元組」，都值得問一句：fuzz 過沒？

## 請丟掉的舊習慣

1. 只測 happy path。  
2. 用 panic 表示「壞輸入」。  
3. 以為「協定寫了客戶端就不會亂送」。

## 動手練習

### 必做

1. 在 `examples/d03-fuzz` 跑大約 5 秒 fuzz。  
2. 若有失敗輸入，讀懂它、修好、再跑。  

### 選做

1. 幫解析函式加不變量：成功時座標必在合法集合；用 fuzz 驗證。  

## 常見坑

- **Fuzz 函式裡用到外部隨機／時間**：難重現；保持純。  
- **把 fuzz 當效能測試**：它不是 bench；目的是找錯。  
- **忽略留下的失敗語料**：修 bug 時把那筆變成固定測試。

## 延伸閱讀

- <https://go.dev/doc/fuzz/>  
