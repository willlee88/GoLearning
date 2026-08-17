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

## 本章你會建立的心智模型

教學與初期產品用 JSON 很合適：人眼可讀、瀏覽器友善。統一信封：

```json
{ "v": 1, "type": "chat", "payload": { "text": "hi" } }
```

`type` 分派；`v` 管相容；`payload` 可再解到具體 struct。

## Python 對照

| Python | Go |
|--------|-----|
| `pydantic` 模型 | struct + `encoding/json` |
| `model_dump_json` | `json.Marshal` |

## L1 能用

```go
type Envelope struct {
	V       int             `json:"v"`
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
```

`json.RawMessage` 讓你先讀 type 再解 payload。

範例：`examples/g01-envelope`。

## L2 機制

- GC 與 CPU：高頻小包 JSON 有成本——先測再換 Protobuf。  
- 未知 type：忽略或回 `error` 事件，勿 panic。  
- 欄位 `omitempty` 與零值語意要文件化。  

## 請丟掉的 Python 習慣

1. 每種訊息完全不同頂層形狀、無 type。  
2. 信任客戶端任意 JSON 不驗證。  

## 遊戲 Server 連結

`move` / `ready` / `chat` 進；`state` / `system` 出。

## 練習

### 必做

1. 跑 `examples/g01-envelope`。  
2. 為 `move` payload 定義 struct 並驗證範圍。  

## 延伸閱讀

- `encoding/json`  
