---
lessonId: "J3"
title: "優雅關閉"
description: "signal → 停止接收 → drain → Shutdown。"
volume: "j"
order: 3
level: "l2"
status: "ready"
path_required: true
tags: ["shutdown"]
example: "examples/j03-shutdown"
prev: "J2"
next: "J4"
---

## 本章你會建立的心智模型

`SIGTERM` / `SIGINT`（Ctrl+C）時不應直接 `os.Exit`：

1. 停止 `Accept` / 對外宣告 not ready  
2. 關閉 listener / `Server.Shutdown`  
3. 停止 tick、廣播 `server_closing`  
4. 等 in-flight 或超時  
5. flush 日誌與任務佇列  

## Python 對照

| Python | Go |
|--------|-----|
| signal + uvicorn shutdown | `signal.Notify` + `http.Server.Shutdown` |
| atexit | defer + context cancel |

## L1 能用

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer stop()

go srv.ListenAndServe()
<-ctx.Done()
shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
_ = srv.Shutdown(shutdownCtx)
```

範例：`examples/j03-shutdown`。

## L2 機制

- `Shutdown` 等連線跑完；WS 需主動 Close。  
- 超時後 `Close` 強制。  
- K8s `preStop` + `terminationGracePeriodSeconds` 要對齊。  

## 請丟掉的 Python 習慣

1. kill -9 當日常。  
2. 關服時丟半寫入交易。  

## 遊戲 Server 連結

Arena Mini M5：Ctrl+C 可關。

## 練習

### 必做

1. 跑 j03 或 arena-mini，Ctrl+C 看日誌順序。  

## 延伸閱讀

- `net/http` Server.Shutdown  
