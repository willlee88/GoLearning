---
lessonId: "F5"
title: "中介層與路由模式"
description: "logging、auth、recover 的 Handler 包裝。"
volume: "f"
order: 5
level: "l2"
status: "ready"
path_required: true
tags: ["http", "middleware"]
example: "examples/f04-http-api"
prev: "F4"
next: "F6"
---

## 本章你會建立的心智模型

Middleware 就是 `func(http.Handler) http.Handler`。把 logging、CORS、token 檢查、panic recover 疊在外層，業務 handler 保持純淨。這比一開始上巨型框架更利於理解請求生命週期。

## Python 對照

| Python | Go |
|--------|-----|
| FastAPI Depends / Starlette middleware | Handler 包裝 |
| `@app.middleware` | `chain(h)` |

## L1 能用

```go
func withLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
```

## L2 機制

- 順序：外層先入後出（像洋蔥）。  
- 認證失敗直接寫 401，不呼叫 next。  
- recover middleware 避免單請求 panic 殺進程（見 B3）。  

## 請丟掉的 Python 習慣

1. 業務裡到處 copy-paste 取 header。  
2. 中介層隱式改全域狀態。  

## 遊戲 Server 連結

`Authorization: Bearer ...` 只保護 HTTP API；WS 可在 query/token 或首訊驗證（F7）。

## 練習

### 必做

1. 為 `f04-http-api` 加 withLog。  
2. 若 `X-Token != secret` 回 401。  

## 延伸閱讀

- Just for Func / 標準庫 middleware 模式文章  
