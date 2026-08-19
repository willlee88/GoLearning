---
lessonId: "F4"
title: "net/http 基礎"
description: "Handler、ServeMux、狀態碼與 JSON API。"
volume: "f"
order: 4
level: "l2"
status: "ready"
path_required: true
tags: ["http"]
example: "examples/f04-http-api"
prev: "F3"
next: "F5"
---

## 這章你會搞懂什麼

`net/http` 是 Go 做後端**控制面**的預設答案：健康檢查、登入、房間列表、管理 API、metrics 拉取。

核心概念：

- **`Handler`**：實作 `ServeHTTP(w, r)` 的東西  
- **`ServeMux`**：路由器；Go 1.22+ 支援 `"GET /path"` 這種方法＋路徑模式  
- 每個進來的請求，Server 預設用**一個 goroutine** 跑你的 handler  

遊戲的即時玩法通常走 WebSocket；HTTP 負責「連上遊戲之前／旁邊」的事。

## 先跟 Python 對一下

| Python | Go | 注意 |
|--------|-----|------|
| Flask／FastAPI 的 route | `http.HandleFunc`／自建 mux | 先標準庫，再決定要不要框架 |
| `return dict` → JSON | `json.NewEncoder(w).Encode(...)` | 記得設 `Content-Type` |
| `Depends`／middleware | 下一章 F5 手寫包裝 | |
| uvicorn／gunicorn worker | `http.Server` + goroutine／請求 | 模型不同，但都要注意超時 |

## 怎麼寫（能跑的最小例子）

```go
mux := http.NewServeMux()
mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
})
log.Fatal(http.ListenAndServe(":8080", mux))
```

JSON API 示意：

```go
mux.HandleFunc("POST /v1/echo", func(w http.ResponseWriter, r *http.Request) {
	var in map[string]any
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(in)
})
```

範例：`examples/f04-http-api`。

較穩的啟動方式是填 `http.Server`（可設超時），不要永遠只用便利函式：

```go
s := &http.Server{
	Addr:              ":8080",
	Handler:           mux,
	ReadHeaderTimeout: 5 * time.Second,
}
log.Fatal(s.ListenAndServe())
```

## 為什麼這樣設計／底層在幹嘛

1. **Handler 介面很小**  
   任何滿足 `ServeHTTP` 的型別都能掛上去；middleware 就是包一層 Handler（F5）。

2. **請求上下文**  
   `r.Context()` 在客戶端取消／超時時會取消。往下傳 DB 查詢才有意義。

3. **錯誤用狀態碼，不要 panic**  
   4xx＝客戶問題，5xx＝伺服器問題。panic 可能搞掉 worker；應用層用錯誤值＋明確回應。

4. **慢速攻擊與超時**  
   不設 `ReadHeaderTimeout` 等，惡意連線可以慢慢餵資料佔資源。生產環境務必填 Server 欄位。

5. **Go 1.22 路由**  
   `"GET /items/{id}"` 這類模式讓標準庫 mux 好用很多；舊教學若只寫 `"/healthz"` 要自己檢查 Method。

## 遊戲 Server 會用在哪

Arena Mini 典型：

- `GET /healthz`：探活  
- `GET /rooms`：房間列表  
- 靜態檔：伺服器網頁客戶端  
- `GET /ws`：升級成 WebSocket（F6）

## 請丟掉的舊習慣

1. **handler 裡用全域 map 不加鎖**——請求併發進來，直接 race。  
2. **一開始堆巨型框架，middleware 魔法看不懂**——先手寫 20 行。  
3. **成功／失敗都回 200 + 字串**——用好 HTTP 狀態碼與 JSON 錯誤體。

## 動手練習

### 必做

1. 跑 `examples/f04-http-api`，用瀏覽器或 curl 打 `/healthz` 與 `/v1/echo`。  
2. 加一個 `GET /v1/time`，回傳 UTC 時間（JSON 或純文字皆可）。  

### 選做

1. 把 `ListenAndServe` 改成自訂 `http.Server` 並設 `ReadHeaderTimeout`。  

## 常見坑

- **寫了 Header 之後才 `WriteHeader`，或 Encode 兩次**：狀態碼與 body 順序要清。  
- **忘記關閉／限制 `r.Body`**：可用 `http.MaxBytesReader`。  
- **在 handler 開 goroutine 卻還用 `w`／`r`**：請求結束後這些不能再用；要複製需要的資料。  
- **路由註冊順序／尾斜線**：測一下 `/rooms` vs `/rooms/`。

## 延伸閱讀

- <https://pkg.go.dev/net/http>  
