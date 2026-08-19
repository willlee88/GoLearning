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

## 這章你會搞懂什麼

**Middleware（中介層）**聽起來很框架，其實在 Go 裡就是一句話：

> 一個函式：吃進 `http.Handler`，吐出新的 `http.Handler`。

用它把「記 log、查 token、CORS、panic recover」疊在外層，讓業務 handler 專心做業務。先手寫理解生命週期，比一開始跳進巨型框架值得。

## 先跟 Python 對一下

| Python | Go | 注意 |
|--------|-----|------|
| Starlette／Django middleware | `func(http.Handler) http.Handler` | 型別上非常透明 |
| FastAPI `Depends` | 多半顯式在 handler 取 header／context | Go 較少「隱式注入魔法」 |
| `@app.middleware("http")` | `mux = withLog(withAuth(mux))` | 注意包裝順序 |

## 怎麼寫（能跑的最小例子）

```go
func withLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
```

認證示意：

```go
func withToken(secret string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Token") != secret {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return // 不呼叫 next
		}
		next.ServeHTTP(w, r)
	})
}
```

掛上去：

```go
var h http.Handler = mux
h = withToken("secret", h)
h = withLog(h) // 最外層
http.ListenAndServe(":8080", h)
```

可在 `examples/f04-http-api` 上改著玩。

## 為什麼這樣設計／底層在幹嘛

1. **洋蔥模型（onion）**  
   請求從外往內進，回應從內往外回。`withLog` 若包在最外，能量到含認證失敗在內的時間；若你想「只量通過認證的」，順序就要換。

2. **失敗可以短路**  
   認證失敗直接寫 401，不呼叫 `next`——業務碼零感覺。

3. **recover middleware**  
   單一請求 panic 不應弄垮整個行程（見 B3）。外層：

   ```go
   defer func() {
   	if rec := recover(); rec != nil {
   		log.Printf("panic: %v", rec)
   		http.Error(w, "internal", 500)
   	}
   }()
   next.ServeHTTP(w, r)
   ```

4. **跟框架的關係**  
   chi、echo、gin 等本質都是這套模式的糖衣與路由增強。你寫得動手動版，框架文件就看得懂。

## 遊戲 Server 會用在哪

- HTTP API：`Authorization: Bearer ...` 保護管理端、發房。  
- 公開的 `/healthz` 不要套強制登入。  
- WebSocket：常在 **upgrade 前** 查 token（query 或 header），或約定「連上後第一則認證訊息」（F7）。  
- 統一加 request id：middleware 產生，塞進 context（克制用 `WithValue`）。

## 請丟掉的舊習慣

1. **每個 handler 複製貼上取 header、查 token。**  
2. **middleware 偷偷改全域狀態**——除錯會哭。  
3. **以為「有框架就安全」**——auth 邏輯寫錯一樣裸奔。

## 動手練習

### 必做

1. 為 `f04-http-api` 加上 `withLog`。  
2. 若 `X-Token != secret` 回 401；用 curl 帶／不帶 header 各測一次。  

### 選做

1. 加 recover middleware，故意在 handler `panic("boom")` 看是否變 500。  

## 常見坑

- **包裝順序反了**：log 看不到 401、或 auth 跑在靜態檔伺服之外／之內搞錯。  
- **在 middleware 讀完 `r.Body` 卻不還原**：後面 handler 讀不到——要嘛別讀，要嘛緩衝重放。  
- **寫完回應後又寫**：`http.Error` 後不要再 `Encode`。  

## 延伸閱讀

- 搜尋「Go http middleware Handler」任何一篇標準庫風格教學即可；重點是型別簽名。  
