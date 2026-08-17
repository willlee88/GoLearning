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

## 本章你會建立的心智模型

`encoding/json` 是 Go 最常用的序列化。struct tag 控制欄位名；只有**匯出欄位**會被處理。熱路徑要注意分配；協定層常用 `RawMessage` 二次解碼（見 G1）。

## Python 對照

| Python | Go |
|--------|-----|
| `json.loads/dumps` | `Unmarshal/Marshal` |
| pydantic | struct + 手動驗證 |

## L1 能用

```go
type P struct {
	Name string `json:"name"`
	HP   int    `json:"hp"`
}
b, _ := json.Marshal(P{Name: "Ada", HP: 10})
var p P
_ = json.Unmarshal(b, &p)
```

範例：`examples/e05-json`。

## L2 機制

- `omitempty` 省略零值。  
- `json:"-"` 忽略。  
- 未知欄位預設忽略。  
- `Decoder.DisallowUnknownFields` 可嚴格。  
- 數字進 `any` 會變 `float64`。  

## 請丟掉的 Python 習慣

1. 依賴 dict 到處傳、無 schema。  
2. 熱迴圈反覆 Marshal 大物件不測量。  

## 遊戲 Server 連結

`state` 快照、信封 payload 都靠 json。

## 練習

### 必做

1. 跑 `examples/e05-json`。  
2. 加一個未匯出欄位證明不會出現在 JSON。  

## 延伸閱讀

- `encoding/json`  
