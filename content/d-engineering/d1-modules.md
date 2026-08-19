---
lessonId: "D1"
title: "go.mod：依賴要可重現，不要碰運氣"
description: "module path、語意化版本、go.sum，以及 Minimal Version Selection 的白話意思。"
volume: "d"
order: 1
level: "l2"
status: "ready"
path_required: true
tags: ["modules"]
example: ""
prev: "D0"
next: "D2"
---

## 這章你會搞懂什麼

`go.mod` 宣告兩件大事：**這個模組叫什麼**、**它依賴哪些版本**。  
旁邊的 `go.sum` 則像校驗清單，降低「下載到被竄改或不一致內容」的風險。

Go 解析版本時用 **Minimal Version Selection（MVS，最小版本選擇）**：在滿足所有約束的前提下，選比較「小」的版本，而不是永遠追最新。  
白話：比較可預期，比較不會因為某天傳遞依賴突然跳大版而全家爆炸——當然你仍要管自己的升級。

## Python 對照

| Python | Go |
|--------|-----|
| `requirements.txt`／Poetry／pip-tools | `go.mod` + `go.sum` |
| venv 隔離環境 | 通常**不需要** venv；以 module 為界 |
| 「全球 site-packages 裝一堆」 | 看當前模組的依賴圖 |
| 鎖檔格式百家爭鳴 | 工具鏈內建，文化統一 |

## 怎麼寫

```bash
go mod init example.com/x
go get example.com/y@v1.2.3
go mod tidy
```

日常節奏：

1. 改 import／加依賴  
2. `go mod tidy` 收斂 `go.mod`／`go.sum`  
3. 把兩個檔案一起進版控——**不要手改 go.sum**

讀真實範例：打開 `demo/arena-mini/server/go.mod`，看 `require` 長什麼樣。

## 細節

### module path 是什麼？

它是別人 `import` 你時的前綴，也是你自己套件路徑的根。  
本地資料夾名可以不同，但團隊最好約定一致，不然文件與 import 會對不上腦。

### 為什麼有 `go.sum`？

它記錄依賴內容的雜湊。CI／同事機器拉依賴時可以校驗。  
手動編輯幾乎一定是錯的；應透過 `go get`／`go mod tidy` 更新。

### `replace` 什麼時候用？

本機聯調、暫時指向 fork、或 monorepo 內替換路徑。  
很好用，但別把暫時的 `replace` 忘在正式發佈流程——文件化它。

### `internal/` 目錄

編譯器強制：只有「父樹」能 import `internal`。  
很適合放遊戲 Server 不該被外部專案當公共庫亂引用的碼。

### 進階可先略過

- MVS 的細節與 `go list -m all`  
- `GOPROXY`／`GOPRIVATE` 在公司網路的設定  

## 遊戲 Server 會用在哪

demo 依賴越少越好：少依賴＝少驚喜。  
Arena Mini 用 `x/net` 做 WebSocket 是刻意選擇——你要能在 `go.mod` 裡一眼看懂「我到底吃了什麼」。

## 請丟掉的舊習慣

1. 全域亂裝套件、不鎖版本。  
2. 手改 `go.sum`「對一下字」。  
3. 把 `replace` 當永久架構卻不說。

## 動手練習

### 必做

1. 打開 `demo/arena-mini/server/go.mod`，用自己的話解釋每一段 `require`。  
2. 在任意練習目錄跑 `go mod tidy`，看它改了什麼。  

### 選做

1. `go list -m all` 看完整模組圖（可能比你想的長）。  

## 常見坑

- **在錯目錄 init／tidy**：先確認旁邊就是你要的 `go.mod`。  
- **把 module path 當資料夾裝飾**：import 對不上時先查 `go.mod` 第一行。  
- **私有模組拉失敗**：常是 proxy／SSH／`GOPRIVATE` 問題，不是 Go「壞了」。

## 延伸閱讀

- <https://go.dev/ref/mod>  
