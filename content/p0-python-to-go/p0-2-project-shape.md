---
lessonId: "P0.2"
title: "專案形態對照"
description: "從 Python 腳本／套件心智切到 Go 的 package、module 與目錄習慣。"
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

## 本章你會建立的心智模型

Python 專案常是「一堆 `.py` + venv + 進入點」；Go 專案是 **module（依賴與版本邊界）+ package（編譯與可見性邊界）**。你要先會用正確的「盒子」放程式，後面的測試、範例與 Server 才放得下。

## Python 對照

| Python | Go |
|--------|-----|
| 專案根 + `pyproject.toml` / `requirements.txt` | 根目錄 `go.mod` |
| package = 目錄 + `__init__.py`（概念上） | package = 同一目錄下相同 `package` 名 |
| `if __name__ == "__main__"` | `package main` + `func main()` |
| venv 隔離依賴 | module 路徑 + 版本；toolchain 可釘選 |
| `src/` 佈局可選 | 常見：`cmd/`、`internal/`、`pkg/`（慣例非強制） |

## L1 能用

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

```bash
go mod init example.com/myapp
go run .
go test ./...
```

本站第一個可跑範例：`examples/p0-config-stats`（讀 JSON 設定、統計假玩家）。

## L2 機制

### Module

- `go.mod` 宣告 **module path**（匯入前綴）與 Go 版本。
- 依賴解析採 **Minimal Version Selection (MVS)**（後續 D 卷深講）。

### Package 與可見性

- 名稱**大寫開頭**的識別字可被外部 package 匯入（匯出）。
- 小寫 = 套件私有。這比 Python 的「底線約定」更硬。

### 建議目錄（Server 友善）

```text
cmd/server/main.go     # 入口薄薄一層
internal/room/         # 不讓外部 module 匯入
internal/session/
pkg/protocol/          # 若真要給別人匯入（可選）
```

遊戲 Server 習慣：**入口薄、規則純、I/O 在邊緣**。

## L3 深潛（可選）

- `internal/` 目錄的編譯器強制可見性規則。
- `go.work` 多模組工作區（本 monorepo 後續可選）。

## 請丟掉的 Python 習慣

1. 到處 `import *` 式的扁平亂放。
2. 用檔名當唯一組織手段，卻沒有清晰 package 邊界。
3. 在業務套件裡塞 `print` 除錯當永久日誌。

## 遊戲 Server 連結

之後 Arena Mini 會類似：

```text
demo/arena-mini/server/
  cmd/server/
  internal/gateway/
  internal/room/
  internal/session/
```

## 練習

### 必做

1. 閱讀並執行 `examples/p0-config-stats`（需 Go）。
2. 指出哪個檔是 `main`、哪個是可匯出 API。

### 選做

1. 把範例拆成 `cmd/` + `internal/stats/` 兩個 package。

## 常見坑與如何看見

- **module path 與資料夾名無關**，但匯入路徑必須一致。
- 循環匯入：編譯器直接拒絕——用介面或把共用型別下沉到更低層 package。

## 延伸閱讀

- <https://go.dev/doc/modules/layout>
- <https://go.dev/doc/code>
