---
lessonId: "E2"
title: "io.Reader 與 Writer"
description: "串流抽象；組合優於巨型 API。"
volume: "e"
order: 2
level: "l2"
status: "ready"
path_required: false
tags: ["io"]
example: "examples/e02-io"
prev: "E1"
next: "E3"
---

## 這章你會搞懂什麼

Go 處理「讀寫資料」的核心不是「檔案物件有一百個方法」，而是兩個超小的介面：

- **`io.Reader`**：你能從我這裡讀位元組  
- **`io.Writer`**：你能把位元組寫給我  

檔案、網路連線、記憶體 buffer、壓縮套件……很多東西都實作它們。所以你可以 `io.Copy`、一層包一層，不必為每種來源重寫邏輯。

## 先跟 Python 對一下

| Python | Go | 注意 |
|--------|-----|------|
| file-like：`.read()` / `.write()` | `io.Reader` / `io.Writer` | Go 在**編譯期**用介面約束，不是靠「長得像」 |
| `shutil.copyfileobj(a, b)` | `io.Copy(dst, src)` | 概念很像 |
| `io.BytesIO` | `bytes.Buffer` / `strings.Reader` | 記憶體裡的 Reader/Writer |
| 一次 `f.read()` 讀完 | 迴圈 `Read` 進 buffer | `Read` **不保證**一次讀滿；這點超重要 |

人話版jargon：

- **串流（stream）**：資料邊來邊處理，不一定整包進記憶體。  
- **組合（composition）**：小介面往上包（限速、計數、緩衝），而不是繼承一個巨型基底類別。

## 怎麼寫（能跑的最小例子）

```go
package main

import (
	"io"
	"os"
	"strings"
)

func main() {
	_, _ = io.Copy(os.Stdout, strings.NewReader("hello io.Reader\n"))
}
```

範例目錄：`examples/e02-io`（`go run .` 會印一行字）。

### 自己讀一點資料

```go
r := strings.NewReader("abcdef")
buf := make([]byte, 3)
n, err := r.Read(buf)
// n 可能是 3；再 Read 一次才會拿到剩下的
```

### 限制最多讀多少：`LimitReader`

```go
limited := io.LimitReader(r, 1024) // 外層看起來還是 Reader
n, err := io.Copy(dst, limited)
```

網路協定讀「長度前綴」時，常用 `io.ReadFull` 硬是凑滿 N 個 byte（F3 會用到）——因為單次 `Read` 可能只給你一部分。

```go
buf := make([]byte, 4)
_, err := io.ReadFull(r, buf) // 讀不滿 4 byte 就回錯
```

## 為什麼這樣設計／底層在幹嘛

1. **介面很小，實作門檻低**  
   你只要寫好 `Read([]byte) (int, error)`，就能接入整套 I/O 生態。這跟「必須繼承某個 AbstractStream」差很多。

2. **`Read` 的契約（一定要記住）**  
   - 回傳 `(n int, err error)`：這次實際讀了幾個 byte。  
   - `n > 0` 且 `err == io.EOF` 是合法的：最後一截資料 + 結束。  
   - **不要**假設「我要 100 byte，`Read` 就給 100」。要满 N byte 用 `io.ReadFull`，或自己迴圈。

3. **包裝器（wrapper）到處都是**  
   `bufio.Reader` 做緩衝、`io.LimitedReader` 做上限、自訂 Reader 做解密……外層函式只依賴介面，測起來也好假資料注入。

4. **`Closer` 是另一個小介面**  
   `io.ReadCloser` = 能讀也能關。網路連線、HTTP Body 常用；用完 `defer Close()`。

## 遊戲 Server 會用在哪

- TCP 連線是 `net.Conn`，本身就是 Reader/Writer。  
- **長度前綴幀**（F3）建立在「從 Reader 精確讀 N byte」。  
- 讀設定檔、寫日誌、把 JSON 編碼到 `http.ResponseWriter`，都是同一套抽象。  
- 測試時可用 `strings.NewReader`／`bytes.Buffer` 假裝一條連線。

## 請丟掉的舊習慣

1. **假設一次 `recv`/`Read` = 一則完整訊息**——TCP 是位元組流，訊息邊界要自己定（F3）。  
2. **為「檔案版／網路版」各寫一套幾乎一樣的函式**——改吃 `io.Reader`。  
3. **忽略 `n`，只看 `err`**——可能丟資料。

## 動手練習

### 必做

1. 跑 `examples/e02-io`。  
2. 用 `io.LimitReader` 包一個 `strings.NewReader`，確認最多只能拷貝你設的上限。  

### 選做

1. 實作一個計數 Writer：包住 `os.Stdout`，記錄總共寫了多少 byte。  

## 常見坑

- **把 EOF 當災難**：很多迴圈是「讀到 EOF 就正常結束」。  
- **`Read` 回 `n>0` 卻沒先處理資料就 return**：最後一塊可能丟掉。  
- **在熱路徑每次 `ioutil.ReadAll` 整包進記憶體**：連線或大檔要用串流／上限。  
- **（名稱）舊碼 `ioutil`**：新程式用 `io` / `os` 對應函式即可。

## 延伸閱讀

- <https://pkg.go.dev/io>  
