---
lessonId: "F3"
title: "封包邊界與長度前綴"
description: "TCP 是流：粘包/拆包與 length-prefix 幀。"
volume: "f"
order: 3
level: "l2"
status: "ready"
path_required: true
tags: ["protocol", "tcp"]
example: "examples/f03-framed"
prev: "F2"
next: "F4"
---

## 本章你會建立的心智模型

TCP 不保留「訊息邊界」。一次 `Read` 可能半包、多包混雜。應用層要做 **framing**：常見是 **4-byte 大端長度 + payload**。這是自訂遊戲協定的地基。

## Python 對照

| Python | Go |
|--------|-----|
| 自己拼 buffer | `encoding/binary` + `io.ReadFull` |
| `struct.pack('!I', n)` | `binary.BigEndian.PutUint32` |

## L1 能用

```go
// 寫
buf := make([]byte, 4+len(payload))
binary.BigEndian.PutUint32(buf[:4], uint32(len(payload)))
copy(buf[4:], payload)
_, err = conn.Write(buf)

// 讀
var hdr [4]byte
if _, err = io.ReadFull(r, hdr[:]); err != nil { ... }
n := binary.BigEndian.Uint32(hdr[:])
payload := make([]byte, n)
_, err = io.ReadFull(r, payload)
```

範例：`examples/f03-framed`。

## L2 機制

- 必用 `io.ReadFull` 湊滿 n bytes。  
- **限制最大幀** 防 OOM（例如 1MB）。  
- JSON 可做 payload；二進位見 G 卷。  
- WebSocket 已有幀，但仍要在文字/二進位訊息上定義應用信封。

## 請丟掉的 Python 習慣

1. 假設 `recv` 一次就是一條完整訊息。  
2. 用換行當協定卻不處理部分行（可用，但是另一種 framing）。  

## 遊戲 Server 連結

`Command` / `StateSnapshot` 都先成 payload，再套長度前綴或 WS message。

## 練習

### 必做

1. 跑 `examples/f03-framed` 測試。  
2. 拒絕 `n > maxFrame` 並回錯誤。  

## 延伸閱讀

- `io.ReadFull` 文件  
