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

## 這章你會搞懂什麼

**TCP** 提供的是**可靠的位元組流**（bytes 會到、有序），不是「訊息佇列」。Server 典型骨架是：

1. `Listen` 聽一個位址  
2. 迴圈裡 `Accept` 接到新連線  
3. 每個連線丟給一個（或一組）goroutine 處理  
4. 用完 `Close`；讀寫要設 **deadline**，避免永久卡住

## 先跟 Python 對一下

| Python | Go | 注意 |
|--------|-----|------|
| `socket.bind`／`listen`／`accept` | `net.Listen`／`Accept` | 流程幾乎一一對應 |
| `conn.recv`／`sendall` | `conn.Read`／`Write` | 一樣可能短讀短寫 |
| `conn.makefile` | `bufio.NewReader(conn)` | 緩衝讀較好寫 |
| `asyncio` 一堆 Stream | 每連線 goroutine 很常見 | 模型直覺，仍要管生命週期 |

## 怎麼寫（能跑的最小例子）

```go
ln, err := net.Listen("tcp", ":9000")
if err != nil {
	log.Fatal(err)
}
defer ln.Close()

for {
	conn, err := ln.Accept()
	if err != nil {
		continue // 實務要分辨「listener 關閉」與暫時錯誤
	}
	go handle(conn)
}

func handle(conn net.Conn) {
	defer conn.Close()
	_ = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	// 讀寫…
}
```

範例：`examples/f02-tcp-echo`（echo：你打什麼回什麼）。

Client 端概念：

```go
conn, err := net.Dial("tcp", "127.0.0.1:9000")
defer conn.Close()
_, _ = conn.Write([]byte("hi\n"))
```

## 為什麼這樣設計／底層在幹嘛

1. **`net.Conn` 就是 Reader＋Writer**  
   後面所有 framing、`io.Copy`、JSON Decoder 都能接上去。

2. **Accept 迴圈與業務分離**  
   接受連線要快；慢邏輯放 `handle`。關服時先 `ln.Close()` 讓 Accept 失敗跳出，再等進行中的連線收尾（優雅關閉進階題，J 卷也會談）。

3. **Deadline 是必需品**  
   `SetReadDeadline`／`SetWriteDeadline`／`SetDeadline`：時間到 `Read`/`Write` 回 deadline 錯誤。沒設的話，對端不說話，你的 goroutine 可能永遠卡在 `Read`。

4. **短讀短寫**  
   `Write` 也可能一次沒寫完（雖然對 TCP 在很多情況下你會包迴圈）。讀訊息邊界更是應用層責任 → F3。

5. **連線是 OS 資源**  
   檔案描述符有限。洩漏 `Close`＝慢慢把伺服器卡死。`defer conn.Close()` 是基本禮貌。

## 遊戲 Server 會用在哪

- 非瀏覽器客戶端（手機原生、桌面）常用 **TCP + 長度前綴協定**（F3）。  
- 瀏覽器不能亂開裸 TCP 到任意 port（受限），網頁遊戲多走 **WebSocket**（F6），底層其實仍是 TCP（或 TLS 上的 TCP）。  
- 健康檢查、匹配服務之間的內部 RPC，有時也是 TCP／HTTP2。

## 請丟掉的舊習慣

1. **單執行緒 accept 完又自己同步處理**，下一個玩家卡在門外。  
2. **忽略 `Close`、忽略錯誤**，fd 與 goroutine 洩漏。  
3. **讀到底都不設 timeout**，假死連線堆成山。

## 動手練習

### 必做

1. 跑 echo server；另開終端當 client（範例若有 client 就用；或自己 `net.Dial`）。  
2. 為 `Read` 設 5 秒 deadline，超時印出錯誤並關閉。  

### 選做

1. 用 `bufio.Scanner` 讀一行一行（這是「換行 framing」的一種，與 F3 長度前綴對照）。  

## 常見坑

- **Listen 位址**：`:9000` 聽所有介面；`127.0.0.1:9000` 只本機。  
- **Accept 錯誤一律 `continue`**：listener 關閉時應結束迴圈，否則空轉。  
- **在 handle 裡 panic**：會弄掉單一 goroutine；重要服務外層可 recover，但更好是別 panic。  
- **把同一個 `conn` 給很多 goroutine 無鎖同寫**：會交錯；寫要串行（F6 同樣問題）。

## 延伸閱讀

- <https://pkg.go.dev/net>  
