---
lessonId: "E3"
title: "flag 與環境變數設定"
description: "行程設定的務實做法。"
volume: "e"
order: 3
level: "l2"
status: "ready"
path_required: false
tags: ["config"]
example: "examples/e03-flag"
prev: "E2"
next: "E4"
---

## 這章你會搞懂什麼

程式一啟動就要知道：聽哪個 port、靜態檔放哪、是不是 debug 模式。這些叫**設定（configuration）**。

本章你會用兩種最務實的來源：

1. **命令列 flag**（本機開發超方便）：`go run . -addr :8080`  
2. **環境變數**（部署／容器友善）：`ADDR=:8080`

並記住一句業界常講的人話版 **12-factor**：設定跟程式碼分開，**秘密（密碼、token）不要寫死在 repo**。

## 先跟 Python 對一下

| Python | Go | 注意 |
|--------|-----|------|
| `argparse` / `click` | 標準庫 `flag` | API 比較陽春，但夠小工具與服務啟動用 |
| `os.environ["ADDR"]` | `os.Getenv("ADDR")` | 不存在時 Go 回空字串，不會直接炸 |
| `.env` + python-dotenv | 可自行讀，或部署系統注入 env | 正式環境常由 systemd／K8s／compose 注入 |
| 設定類別 + pydantic | struct 自己填 | 先手動，夠清楚再考慮套件 |

## 怎麼寫（能跑的最小例子）

### flag

```go
package main

import (
	"flag"
	"fmt"
)

func main() {
	n := flag.Int("n", 1, "times to greet")
	name := flag.String("name", "world", "name")
	flag.Parse() // 記得呼叫，否則一直是預設值

	for i := 0; i < *n; i++ {
		fmt.Println("hello", *name)
	}
}
```

```bash
go run . -n 3 -name Ada
```

範例：`examples/e03-flag`。

注意：`flag.String` 回的是 `*string`（指標），所以要用 `*name` 取值。也可以先宣告變數再用 `flag.StringVar`。

### 環境變數

```go
addr := os.Getenv("ADDR")
if addr == "" {
	addr = ":8080" // 預設值
}
```

實務上常見組合：**env 當預設，flag 能覆寫**（本機方便、線上用 env）。

```go
defaultAddr := ":8080"
if v := os.Getenv("ADDR"); v != "" {
	defaultAddr = v
}
addr := flag.String("addr", defaultAddr, "listen address")
flag.Parse()
```

## 為什麼這樣設計／底層在幹嘛

1. **`flag.Parse` 會動到程式參數**  
   它吃 `os.Args`。解析後剩餘非 flag 參數可用 `flag.Args()`。若你在函式庫裡亂 Parse，可能跟別人搶——服務的 flag 通常集中在 `main`。

2. **環境變數適合「同一份 binary、不同環境」**  
   Dev／Staging／Prod 不改碼，只改 `ADDR`、`DATABASE_URL`。容器時代這几乎是預設。

3. **型別與驗證還是你的責任**  
   `Getenv` 只給字串。port、布林、數值要自己 parse，失敗要在啟動期就報錯（fail fast），別跑到一半才發現。

4. **秘密管理**  
   Token、金鑰：env 或密鑰服務，**不要** commit 進 git，也不要打進前端能下載的設定檔。

## 遊戲 Server 會用在哪

Arena Mini 這類 demo 常見：

- `ADDR`：HTTP／WS 聽哪裡  
- `WEB_DIR`：靜態網頁目錄  

本機用 flag 快速改 port；Docker Compose 用 environment 注入。之後加 DB、Redis 也是同一套路。

## 請丟掉的舊習慣

1. **把 port、密碼寫死在原始碼**——改設定還要重新編譯／還會漏進 git。  
2. **設定散落二十個套件各自 `Getenv`**——集中在 `main` 或小 `config` 載入一次，傳 struct 下去。  
3. **用「改預設值」代替文件**——flag 的 usage 字串、README 的 env 表要寫清楚。

## 動手練習

### 必做

1. 跑 `examples/e03-flag`，試 `-n` 與 `-name`。  
2. 幫它加一個 `-addr`（字串，預設 `:8080`）並印出來。  

### 選做

1. 讓 `-addr` 的預設值來自環境變數 `ADDR`。  

## 常見坑

- **忘了 `flag.Parse()`**：怎麼傳參數都像沒效。  
- **`Getenv` 空字串當「有設」**：要不要區分「未設定」與「設成空」要想清楚。  
- **在 init() 裡 Parse flag**：測試與多 main 時容易 spec 難搞。  
- **布林 flag**：`-debug` / `-debug=false` 的語法跟直覺不完全一樣，先看 `flag` 文件。

## 延伸閱讀

- <https://pkg.go.dev/flag>  
- <https://pkg.go.dev/os#Getenv>  
