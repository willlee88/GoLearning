---
lessonId: "F2"
title: "TCP server 與 client"
description: "net.Listen、Accept、連線生命週期。"
volume: "f"
order: 2
level: "l2"
status: "ready"
path_required: true
tags: ["tcp", "net"]
example: "examples/f02-tcp-echo"
prev: "F1"
next: "F3"
---

## 本章你會建立的心智模型

TCP 提供可靠位元組流。Server：`Listen` → 迴圈 `Accept` → 每連線通常一個（或一組）goroutine 處理。連線是資源：要 `Close`、要處理半開、要設 deadline。

## Python 對照

| Python | Go |
|--------|-----|
| `socket.bind/listen/accept` | `net.Listen` / `Accept` |
| `conn.makefile` | `bufio` + `net.Conn` |

## L1 能用

```go
ln, err := net.Listen("tcp", ":9000")
if err != nil { log.Fatal(err) }
defer ln.Close()

for {
	conn, err := ln.Accept()
	if err != nil { continue }
	go handle(conn)
}
```

範例：`examples/f02-tcp-echo`。

## L2 機制

- `net.Conn` 同時是 `ReadWriter`。  
- `SetDeadline` / `SetReadDeadline` 防止永久阻塞。  
- Accept 迴圈與 handle 分離；關閉 listener 結束服務。  
- 優雅關閉：停止 Accept → 等待進行中連線（進階）。

## 請丟掉的 Python 習慣

1. 單執行緒 accept+handle 阻塞所有人。  
2. 忽略 `Close` 錯誤與洩漏 fd。  

## 遊戲 Server 連結

非瀏覽器客戶端常用 TCP + 長度前綴協定（F3）。瀏覽器則走 WebSocket（F6）。

## 練習

### 必做

1. 跑 echo server，用第二終端 `go run` client 或 `Test-NetConnection`／自寫 client。  
2. 為 `Read` 設 5s deadline。  

## 延伸閱讀

- <https://pkg.go.dev/net>  
