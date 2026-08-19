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

## 這章你會搞懂什麼

TCP 給你的是**一條位元組河流**，不幫你記住「這是第幾則訊息」。

所以會發生：

- **粘包**：一次 `Read` 拿到兩則訊息黏在一起  
- **拆包／半包**：一則訊息分兩次才到齊  

應用層必須做 **framing（組幀）**：定義訊息邊界。遊戲自訂協定最常見的一種是 **長度前綴（length-prefix）**：

```text
[4 byte 大端長度 N][N byte payload]
```

## 先跟 Python 對一下

| Python | Go | 注意 |
|--------|-----|------|
| `struct.pack('!I', n)` | `binary.BigEndian.PutUint32` | `!`／BigEndian＝網路位元組序 |
| 自己拼 `bytearray` | `[]byte` + `io.ReadFull` | 讀滿 N byte 是關鍵 |
| `asyncio` Stream 讀精確長度 | 同上 | 觀念相同 |

人話：

- **大端（Big-Endian）**：先送高位 byte，網路協定慣用。  
- **`io.ReadFull`**：硬是讀到夠 N 個 byte，或回錯誤——不要只呼叫一次 `Read`。

## 怎麼寫（能跑的最小例子）

```go
// 寫一幀
buf := make([]byte, 4+len(payload))
binary.BigEndian.PutUint32(buf[:4], uint32(len(payload)))
copy(buf[4:], payload)
_, err = conn.Write(buf)

// 讀一幀
var hdr [4]byte
if _, err = io.ReadFull(r, hdr[:]); err != nil { /* … */ }
n := binary.BigEndian.Uint32(hdr[:])
payload := make([]byte, n)
_, err = io.ReadFull(r, payload)
```

完整一點、含最大長度保護：見 `examples/f03-framed`（`WriteFrame`／`ReadFrame`，`MaxFrame = 1MiB`）。

```bash
cd examples/f03-framed
go test .
```

## 為什麼這樣設計／底層在幹嘛

1. **為什麼不能「Read 一次當一則」**  
   OS 與 TCP 以 segment／緩衝交資料，跟你的業務訊息無關。對端寫兩次、你可能一次讀完；對端寫一次，你也可能分兩次讀。

2. **為什麼要 `ReadFull`**  
   長度頭剛好 4 byte，也必須凑滿；payload 的 N byte 同理。半包時 `ReadFull` 會繼續讀到滿或失敗。

3. **一定要限制最大幀**  
   惡意或 bug 送來「長度 = 2GB」，你 `make([]byte, n)` 會直接炸記憶體。範例用 `ErrFrameTooLarge` 拒絕。

4. **其他 framing**  
   - 換行分隔（文字協定）  
   - 固定長度  
   - WebSocket：底層已有訊息幀，但**應用層**仍要定義 JSON 信封（F6／G1）

5. **payload 裡放什麼**  
   可以是 JSON、Protobuf、自訂二進位。組幀負責「切開」，不負責語意。

## 遊戲 Server 會用在哪

- 桌面／手機 TCP 客戶端：`Command`、`StateSnapshot` 先成 payload，再套長度前綴。  
- WebSocket 教學 demo：一則 WS message ≈ 一幀，但仍要 `type`／`payload` 信封。  
- 任何「黏在一起的 byte 流」上的多工訊息。

## 請丟掉的舊習慣

1. **假設 `recv`／`Read` 一次就是完整訊息。**  
2. **用換行當協定卻不處理「半行」**——用 `bufio.Scanner` 可以，但要設 `MaxScanTokenSize`，並理解它是另一種幀。  
3. **相信長度欄位不驗證上限。**

## 動手練習

### 必做

1. 跑 `examples/f03-framed` 的測試。  
2. 確認當 `n > maxFrame` 時會回錯誤（讀文件或自己寫測試呼叫）。  

### 選做

1. 在記憶體 `bytes.Buffer` 連續寫兩幀，再連續讀兩幀，證明邊界正確。  

## 常見坑

- **長度欄位用有號整數亂轉**：用 `uint32`，並檢查與 `int` 轉換溢位。  
- **寫了 header 忘了寫完 payload**：對端會卡在 `ReadFull`。  
- **混用 LittleEndian**：前後端位元組序不一致會解出天書長度。  
- **把 framing 跟加密／壓縮順序搞亂**：通常先有明文幀格式，再決定外層包什麼。

## 延伸閱讀

- <https://pkg.go.dev/io#ReadFull>  
- <https://pkg.go.dev/encoding/binary>  
