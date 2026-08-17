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

## 本章你會建立的心智模型

`net/http` 是 Go 後端控制面的預設答案。`Handler` 是 `ServeHTTP(w,r)`；Go 1.22+ 的 `ServeMux` 支援方法與路徑模式。遊戲 Server 用 HTTP 做 **healthz、登入發 token、管理、metrics**，即時玩法走 WS。

## Python 對照

| Python | Go |
|--------|-----|
| Flask/FastAPI route | `http.HandleFunc` / mux |
| `jsonable` response | `json.NewEncoder(w).Encode` |

## L1 能用

```go
mux := http.NewServeMux()
mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
})
log.Fatal(http.ListenAndServe(":8080", mux))
```

範例：`examples/f04-http-api`。

## L2 機制

- 每個請求一個 goroutine（Server 預設模型）。  
- 要設 `ReadHeaderTimeout` 等防慢速攻擊（`http.Server` 結構體）。  
- 錯誤用狀態碼 + JSON body；別 panic。  
- Context：`r.Context()` 在客戶端斷線時取消。

## 請丟掉的 Python 習慣

1. 全域可變狀態無鎖。  
2. 框架魔法中介層不懂就堆——先手寫 20 行 middleware。  

## 遊戲 Server 連結

Arena Mini：`GET /healthz`、`GET /rooms`、靜態客戶端、`GET /ws` 升級。

## 練習

### 必做

1. 跑 `examples/f04-http-api`，用瀏覽器或 curl 打 `/healthz` 與 `/v1/echo`。  
2. 加一個 `GET /v1/time` 回 UTC。  

## 延伸閱讀

- <https://pkg.go.dev/net/http>  
