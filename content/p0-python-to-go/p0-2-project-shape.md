---
lessonId: "P0.2"
title: "專案長什麼樣？module、package、目錄"
description: "從 Python 的腳本／套件心智，換成 Go 的 module（依賴邊界）與 package（編譯與可見性邊界）。"
volume: "p0"
order: 2
level: "l2"
status: "ready"
path_required: true
tags: ["python-bridge", "modules"]
example: "examples/p0-config-stats"
prev: "P0.1"
next: "P0.3"
---

## 這章你會搞懂什麼

Python 專案常是「一堆 `.py` + venv + 某個進入點」。

Go 專案你要先搞懂兩個「盒子」：

- **module（模組）**：依賴與版本的邊界，根目錄有 `go.mod`  
- **package（套件）**：編譯與「別人能不能看到我的名字」的邊界，通常一個目錄一個 package  

盒子放對了，後面的測試、範例、遊戲 Server 才放得下。放錯會怎樣？循環匯入編不過、該藏的 API 全被外面乱用、CI 也不知道從哪測起。

## 先跟 Python 對一下

| Python | Go |
|--------|-----|
| 專案根 + `pyproject.toml` / `requirements.txt` | 根目錄 `go.mod` |
| package ≈ 目錄 + `__init__.py`（概念上） | package = 同一目錄下相同的 `package` 名 |
| `if __name__ == "__main__"` | `package main` + `func main()` |
| venv 隔離依賴 | module 路徑 + 版本；toolchain 也可釘選 |
| `src/` 佈局可選 | 常見慣例：`cmd/`、`internal/`、`pkg/`（慣例，不是語法強制） |

## 怎麼寫（能跑的最小例子）

最小形狀大概長這樣：

```text
myapp/
  go.mod
  main.go          # package main
  player/          # package player
    player.go
```

```go
// go.mod
module example.com/myapp

go 1.22
```

常用指令：

```bash
go mod init example.com/myapp
go run .
go test ./...
```

本站第一個可跑範例：`examples/p0-config-stats`（讀 JSON 設定、統計假玩家）。請真的進目錄跑一遍，不要只看文字。

## 為什麼這樣設計／底層在幹嘛

### Module 在幹嘛

- `go.mod` 宣告 **module path**（匯入時的前綴）與 Go 版本。  
- 依賴怎麼選版本：採 **Minimal Version Selection（MVS，最小版本選擇）**——後續 D 卷會深講。你現在只要知道：不是「永遠裝最新」，而是有一套可預測規則。  

注意：module path **不必**等於資料夾名，但 `import` 路徑必須跟 `go.mod` 一致，亂改會編不過。

### Package 與可見性（誰看得到誰）

- 名稱**大寫開頭**的識別字，可被**外部 package** 匯入使用（匯出／exported）。  
- **小寫開頭** = 這個套件私有。  

這比 Python 靠「底線 `_` 約定」更硬：編譯器直接擋，不是靠大家自覺。好處是 API 邊界清楚；壞處是你不能再「先亂公開再說」。

### 建議目錄（對 Server 友善）

```text
cmd/server/main.go     # 入口薄薄一層：組裝、讀設定、啟動
internal/room/         # 外部其他 module 匯不進來（編譯器強制）
internal/session/
pkg/protocol/          # 若真要給別人匯入再用（可選）
```

遊戲 Server 習慣口訣：**入口薄、規則純、I/O 放在邊緣**。

為什麼？規則套件不依賴網路細節，才能單獨 `go test`；入口只負責接线，重啟／換傳輸層時比較不痛。

### 進階可先略過

- `internal/` 目錄的「編譯器強制：外面 module 不能 import」規則。  
- `go.work` 多模組工作區（本 monorepo 之後可選用）。  

## 遊戲 Server 會用在哪

之後 Arena Mini 會類似這種切法：

```text
demo/arena-mini/server/
  cmd/server/
  internal/gateway/
  internal/room/
  internal/session/
```

你現在學的不是「目錄潔癖」，而是：**之後房間規則、連線閘道、入口**各自有家，才好測、好關機、好換實作。

## 請丟掉的舊習慣

1. 到處 `import *` 式的扁平亂放——Go 沒這招，也不該想這招。  
2. 只用檔名組織，卻沒有清晰 package 邊界——改檔時會牽一髮動全身。  
3. 在業務套件裡塞 `print` 當永久日誌——之後觀測章會用結構化日誌；先別把除錯印成架構。  

## 動手練習

### 必做

1. 閱讀並執行 `examples/p0-config-stats`（需要已安裝 Go）。  
2. 指出哪個檔是 `main`、哪個符號是可匯出 API（看大寫開頭）。  

### 選做

1. 把範例拆成 `cmd/` + `internal/stats/` 兩個 package，再讓它跑起來。  

## 常見坑

- **module path 跟資料夾名無關，但 import 必須對**：編不過時先對一下 `go.mod` 與 `import` 字串。  
- **循環匯入**：A import B、B 又 import A → 編譯器直接拒絕。解法通常是抽介面，或把共用型別下沉到更低層 package。  
- **在錯的目錄跑 `go run`**：確認你在有 `go.mod` 的 module 根（或子目錄）裡。  

## 延伸閱讀

- <https://go.dev/doc/modules/layout>  
- <https://go.dev/doc/code>  
