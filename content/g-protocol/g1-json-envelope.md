---
lessonId: "G1"
title: "JSON 信封實務"
description: "type + payload + v；成本與除錯優勢。"
volume: "g"
order: 1
level: "l2"
status: "ready"
path_required: true
tags: ["json", "protocol"]
example: "examples/g01-envelope"
prev: "G0"
next: "G2"
---

## 這章你會搞懂什麼

教學與早期產品用 **JSON** 當訊息格式很合適：人眼看得懂、瀏覽器友善、抓包好除錯。

建議大家統一用一種**信封（envelope）**外包裝，例如：

```json
{ "v": 1, "type": "chat", "payload": { "text": "hi" } }
```

- **`type`**：這則訊息是什麼，用來分派  
- **`v`**：協定版本，用來相容（G3）  
- **`payload`**：真正內容；可以先當原始 JSON 留著，再二次解碼  

## 先跟 Python 對一下

| Python | Go | 注意 |
|--------|-----|------|
| pydantic 模型 | struct + tag | |
| `model_dump_json()` | `json.Marshal` | |
| 先讀 `type` 再 parse union | `json.RawMessage` | Go 很常見、很好用 |
| 任意 `dict` | 少用 `map[string]any` 當主路徑 | 邊界才用 |

## 怎麼寫（能跑的最小例子）

```go
type Envelope struct {
	V       int             `json:"v"`
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
```

`json.RawMessage` 本質是 `[]byte`，Unmarshal 信封時**不會**急著把 payload 解成具體型別。之後：

```go
switch env.Type {
case "move":
	var m MovePayload
	if err := json.Unmarshal(env.Payload, &m); err != nil { /* … */ }
	// 驗證範圍再套用
default:
	// 未知 type：回 error 事件或忽略，別 panic
}
```

範例：`examples/g01-envelope`（含 `Parse`、`ParseMove`、範圍檢查）。

```bash
cd examples/g01-envelope
go test .
```

## 為什麼這樣設計／底層在幹嘛

1. **為什麼要信封，而不是每種訊息完全不同頂層**  
   路由、log、度量、版本檢查可以共用一層。沒有 `type` 的世界裡，你會寫一堆「猜 JSON 形狀」的爛碼。

2. **為什麼 `RawMessage`**  
   一次解到 `map[string]any` 再轉，又慢又容易型別錯（數字變 `float64`）。先讀 `type`，再解到正確 struct，乾淨。

3. **成本要心裡有數**  
   JSON 有 CPU 與 GC 成本。大廳、聊天、每秒幾十則完全沒差；若變成「每秒成千上萬小包」，再 **benchmark** 決定要不要 Protobuf。先別過早優化。

4. **未知 type／壞 payload**  
   客戶端會亂送。伺服器應回 `{type:"error", ...}` 或切斷，**不要 panic**。

5. **`omitempty` 與零值**  
   文件化「0／false／空字串」省略後的語意，避免前端以為欄位沒送。

## 遊戲 Server 會用在哪

進來（命令）：`move`／`ready`／`chat`…  
出去（事件／狀態）：`state`／`system`／`error`…

F8 的房間廣播，本質就是在送這些信封。

## 請丟掉的舊習慣

1. **每種訊息頂層長得完全不一樣、沒有 type。**  
2. **信任客戶端任意 JSON，不驗證。**  
3. **用字串比對整包 JSON 當協定測試**——應解成 struct 測欄位。

## 動手練習

### 必做

1. 跑 `examples/g01-envelope`。  
2. 為 `move` payload 定義 struct，並驗證 dx／dy 範圍（範例已有可讀）。  

### 選做

1. 加 `chat` payload：`text` 長度上限 200，超過回錯。  

## 常見坑

- **payload 写成字串再包一層**：變成 `"payload": "{\"dx\":1}"`——通常不是你要的。  
- **Marshal 失敗還照樣送**：檢查 err。  
- **大小寫／欄位名跟前端不一致**：靠 struct tag 對齊，並寫一則契約測試。  

## 延伸閱讀

- <https://pkg.go.dev/encoding/json>  
- 複習 **E5**  
