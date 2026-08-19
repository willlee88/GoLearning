---
lessonId: "E5"
title: "encoding/json 實務"
description: "Marshal/Unmarshal、tag、RawMessage、常見陷阱。"
volume: "e"
order: 5
level: "l2"
status: "ready"
path_required: false
tags: ["json", "stdlib"]
example: "examples/e05-json"
prev: "E4"
next: "E6"
---

## 這章你會搞懂什麼

遊戲前後端、HTTP API、WebSocket 訊息，十有八九會碰到 **JSON**（一種文字格式的資料交換，JavaScript Object Notation）。

Go 標準庫是 `encoding/json`。你會搞懂：

- 結構體怎麼跟 JSON 欄位對上（**struct tag**）  
- 只有**匯出欄位**（大寫開頭）才會進出 JSON  
- `Marshal`／`Unmarshal` 基本用法  
- 進階一點：`json.RawMessage` 先留著 payload，之後再二次解碼（G1 信封會用）

## 先跟 Python 對一下

| Python | Go | 注意 |
|--------|-----|------|
| `json.loads` / `dumps` | `json.Unmarshal` / `Marshal` | 名字對調很常背反，看一眼就好 |
| `dict` 到處傳 | 優先用 struct | 有 schema、編譯器幫你抓錯字欄位 |
| pydantic／自行驗證 | struct + 手動檢查 | Go 不幫你驗證「HP 不能是負的」，要自己寫 |
| `orjson` 等加速庫 | 先標準庫，熱路徑再評估 | 先測再換 |

## 怎麼寫（能跑的最小例子）

```go
type Player struct {
	Name   string `json:"name"`
	HP     int    `json:"hp"`
	secret string // 小寫：不會出現在 JSON
}

p := Player{Name: "Ada", HP: 10, secret: "x"}
b, err := json.Marshal(p)
// b == []byte(`{"name":"Ada","hp":10}`)

var q Player
err = json.Unmarshal([]byte(`{"name":"Lin","hp":7,"extra":true}`), &q)
// q.Name=="Lin"；多餘的 extra 預設會被忽略
```

範例：`examples/e05-json`。

### 常用 tag

| tag | 意思 |
|-----|------|
| `` `json:"name"` `` | JSON 鍵叫 `name` |
| `` `json:"hp,omitempty"` `` | 若是**零值**就省略這個欄位 |
| `` `json:"-"` `` | 永遠忽略（不輸出、不解入） |

### 先看 type 再解 payload：`RawMessage`

```go
type Envelope struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"` // 先當原始 JSON 位元組留著
}
```

解完信封後，依 `Type` 再 `json.Unmarshal(env.Payload, &具體結構)`。這是協定層超常見手法（見 G1）。

## 為什麼這樣設計／底層在幹嘛

1. **匯出規則與反射**  
   `encoding/json` 用反射看結構體。未匯出欄位外面的套件（包含 `json`）摸不到，所以自然不會序列化——這不是 bug，是可見性規則。

2. **未知欄位預設忽略**  
   對「向前相容」友善：伺服器加欄位，舊客戶端不解也不炸。若你要嚴格：用 `Decoder` 並 `DisallowUnknownFields()`。

3. **數字進 `any`／`map[string]any` 會變 `float64`**  
   JSON 數字沒有 int／float 之分。解到空介面時要小心，別直接當 `int` 型斷言。

4. **成本**  
   每次 Marshal／Unmarshal 都有解析與配置成本。大廳聊天沒差；若每秒成千上萬小包，要 **benchmark**（D 卷）再決定要不要換二進位協定。

5. **錯誤別吞**  
   範例裡為了短會寫 `_ =`；正式碼請檢查 `err`，否則壞訊息會變成「零值結構體」，超難查。

## 遊戲 Server 會用在哪

- WebSocket 文字訊息：`{"type":"move","payload":{...}}`  
- HTTP `/rooms` 回傳房間列表  
- 狀態快照 `type=state`  
- 設定檔若用 JSON，同一套 API

## 請丟掉的舊習慣

1. **整條管線只用 `dict`／`map[string]any`**——型別與驗證會散落各處。  
2. **信任客戶端 JSON**：還是要檢查範圍（例如 dx 只能 -1..1）。  
3. **熱迴圈反覆 Marshal 超大快照卻從不測量**——先量，再優化。

## 動手練習

### 必做

1. 跑 `examples/e05-json`，看 `secret` 有沒有進 JSON。  
2. 自己加一個未匯出欄位，證明它不會出現。  

### 選做

1. 幫 `Player` 加 `omitempty` 的 `Title string`，空字串時輸出什麼。  
2. 用 `json.Decoder` + `DisallowUnknownFields` 拒絕多餘鍵。  

## 常見坑

- **Unmarshal 忘了傳指標**（要 `&q`）：會錯。  
- **把時間當字串格式搞錯**：`time.Time` 預設 RFC3339；前後端要講好。  
- **`omitempty` 對 `false`／`0`／`""` 都會省略**——布林「假」可能傳不出去，必要時用指標或自訂型別。  
- **stream 用 `Marshal` 拼字串**：HTTP 可直接 `json.NewEncoder(w).Encode(v)`。

## 延伸閱讀

- <https://pkg.go.dev/encoding/json>  
