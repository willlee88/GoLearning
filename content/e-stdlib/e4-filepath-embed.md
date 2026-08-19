---
lessonId: "E4"
title: "filepath 與 embed 概念"
description: "路徑可攜；把靜態檔編進 binary。"
volume: "e"
order: 4
level: "l2"
status: "ready"
path_required: false
tags: ["filepath", "embed"]
example: ""
prev: "E3"
next: "E5"
---

## 這章你會搞懂什麼

兩個很實際的問題：

1. **路徑**：Windows 是 `\ `，Linux／macOS 是 `/`。字串硬拼很容易在某一邊炸掉。  
2. **靜態檔**：前端小頁、`index.html`、設定範本——能不能**打進同一顆執行檔**，部署時少搬一個資料夾？

答案分別是：

- 作業系統路徑用 **`path/filepath`**  
- URL 路徑（web path）用 **`path`**（正斜線語意）  
- 把檔案編進 binary 用 **`embed`**（Go 1.16+）

## 先跟 Python 對一下

| Python | Go | 注意 |
|--------|-----|------|
| `pathlib` / `os.path` | `path/filepath` | `Join`、`Clean`、分隔符依 OS |
| `urllib.parse` 的 path | `path` 套件 | 給 URL 用，不是給本機磁碟 |
| 把靜態檔跟程式一起打包（各種土法） | `//go:embed` | 編譯期嵌進去 |
| `importlib.resources` | `embed.FS` | 執行期用檔案系統介面讀 |

## 怎麼寫（能跑的最小例子）

### 可攜路徑

```go
import "path/filepath"

p := filepath.Join("web", "index.html") // Windows / Unix 都會用對的分隔符
abs, err := filepath.Abs(p)
clean := filepath.Clean(userInputPath)  // 去掉多餘的 . / ..
```

**不要**自己 `"web" + "/" + name` 拼磁碟路徑——在 Windows 上又醜又容易錯。  
反過來說：組 **URL**（例如 `/static/` + 檔名）才用 `path.Join`。

### embed：把資料夾嵌進程式

```go
package main

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed web/*
var webFS embed.FS

func main() {
	// 直接當 HTTP 靜態檔服務（示意）
	sub, _ := fs.Sub(webFS, "web")
	http.Handle("/", http.FileServer(http.FS(sub)))
	_ = http.ListenAndServe(":8080", nil)
}
```

重點：

- 註解 **`//go:embed`** 必須緊挨變數，且套件要 `import "embed"`（即使只為了副作用／型別）。  
- 可 embed 單一檔、多檔、或目錄 pattern。  
- 得到的是編譯進 binary 的唯讀樹；部署常常變成「丟一個 `.exe` 就能跑」。

## 為什麼這樣設計／底層在幹嘛

1. **`filepath` vs `path`**  
   磁碟路徑會碰到磁碟機代號、分隔符、大小寫規則；URL path 語意不同。混用是經典 bug 來源。

2. **`Clean`／`Join` 防目錄逃逸還不夠**  
   若路徑來自使用者，光 `Clean` 不夠當安全邊界；靜態檔服務還要確保解析後仍落在允許的根目錄內（`filepath.Rel` 檢查等）。教學 demo 可以先覺察這件事。

3. **`embed` 換來的是部署簡單，不是「熱更新前端」**  
   改 HTML 要重新編譯。開發期常：用 `WEB_DIR` 指到本機資料夾；發佈再 embed。兩種模式可以並存。

4. **binary 會變大**  
   嵌一整個 `node_modules` 產物要想清楚；只嵌需要的 `web/` 建置結果。

## 遊戲 Server 會用在哪

- Arena Mini 現在多半用目錄伺服靜態頁；你可以**選做**改成 embed，變成單檔 demo。  
- 關卡設定、預設暱稱列表、小圖示，都能 embed。  
- 日誌路徑、資料目錄仍應用 `filepath` 處理，方便 Windows 開發者。

## 請丟掉的舊習慣

1. **字串 `+ "\\"` 或 `+ "/"` 拼本機路徑**。  
2. **假設所有人的工作目錄都跟你一樣**——重要檔案用相對「可設定的根」或 embed。  
3. **把祕密設定檔 embed 進公開發佈的 binary**——有人會直接字串搜出來。

## 動手練習

### 必做

1. 讀官方 `embed` 套件文件的範例，看懂 `//go:embed`。  
2. 用口語解釋：為什麼 `"dir" + "/" + file` 在 Windows 上是壞習慣。  

### 選做

1. 寫一支小程式 embed 一個 `hello.txt` 並印出內容。  
2. 思考 Arena 的 `web/` 改 embed 時，`WEB_DIR` 還要不要留作開發開關。  

## 常見坑

- **`go:embed` 放錯位置**：不能在函式裡；要套件層變數。  
- **pattern 配不到檔**：編譯失敗或嵌到空的——檢查路徑相對 package 目錄。  
- **用 `filepath` 去處理 URL**：反了。  
- **embed 後還用錯的工作目錄去 `os.Open("web/...")`**：嵌進去的要用 `embed.FS`／`fs.FS` 開。

## 延伸閱讀

- <https://pkg.go.dev/path/filepath>  
- <https://pkg.go.dev/embed>  

## 接回主路徑

下一站可以繼續 **E5（JSON）**，或若你趕網路主線 → **F0**。
