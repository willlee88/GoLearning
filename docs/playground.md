# 線上 Go Playground

教學站可在瀏覽器輸入 Go 程式，送到後端 **runner** 執行並回傳結果。

## 使用

1. 開啟 http://localhost:8088/playground/
2. 編輯 `package main` 程式
3. 按 **執行**（或 Ctrl+Enter）
4. 右側顯示 stdout / stderr 與耗時

## 架構

```text
Browser  --POST /api/run-->  nginx (web)
                               |
                               v
                            runner:8091
                               |
                          go run (timeout)
                               |
                          stdout/stderr JSON
```

## 安全限制（學習用）

| 項目 | 限制 |
|------|------|
| 逾時 | 預設 5s，上限 10s |
| 程式大小 | 64KB |
| 輸出大小 | 256KB（超出截斷） |
| 併發 | 最多 4 個同時執行 |
| 速率 | 每 IP 約 30 次／分鐘 |
| 套件 | 須 `package main`；禁止 `os/exec`、`net/http`、`unsafe`、外部 module 等 |
| 網路 | runner 在 Docker `internal` 網路，無外網 |

**不是**多租戶線上 Judge；勿直接暴露到公網而不加認證與更強沙箱。

## 本機開發（不用 Docker 全套）

終端 1 — runner：

```powershell
cd F:\GoLearning\services\runner
go run .
```

終端 2 — 網站：

```powershell
cd F:\GoLearning\web
npm run dev
```

瀏覽器若 CORS／代理不同，可在 console 設：

```js
window.__RUNNER__ = 'http://127.0.0.1:8091/run'
```

Docker 部署時透過 nginx 同源 `/api/run`，不需手動設定。
