---
lessonId: "A15"
title: "嵌入與小型組合"
description: "struct 嵌入、方法提升，以及何時不要嵌。"
volume: "a"
order: 15
level: "l2"
status: "ready"
path_required: false
tags: ["struct", "composition"]
example: "examples/a15-embedding"
prev: "A14"
next: "A16"
---

## 本章你會建立的心智模型

嵌入（embedding）讓外層型別**提升**內層的方法與欄位，是 Go 的組合工具，**不是繼承**。適合「has-a 且想少寫轉發」；不適合深層 is-a 樹。

## Python 對照

| Python | Go |
|--------|-----|
| 多重繼承 / mixin | 嵌入多個 struct（謹慎） |
| 繼承覆寫 | 外層定義同名方法遮罩 |

## L1 能用

```go
type Logger struct{}
func (Logger) Log(msg string) { fmt.Println(msg) }

type Server struct {
	Logger // 匿名嵌入
	Addr string
}

s := Server{Addr: ":8080"}
s.Log("up") // 提升
```

範例：`examples/a15-embedding`。

## L2 機制

- 提升的是方法集；呼叫時 receiver 是內層。  
- JSON 嵌入會影響輸出形狀（欄位提升）。  
- 衝突時需顯式 `s.Logger.Log`。  
- 嵌入 `sync.Mutex` 常見，但**不要匯出**含鎖的 struct 副本。  

## L3 深潛

- 介面滿足：外層是否因嵌入而實作介面。  

## 請丟掉的 Python 習慣

1. 為了「以後擴充」嵌三層 Base。  
2. 把嵌入當子型別多型（不能把 `Server` 當 `Logger` 傳入）。  

## 遊戲 Server 連結

`type roomRuntime struct { game *game.Room; /* conn map */ }` 是持有，不是嵌入；規則物件保持獨立更清晰。

## 練習

### 必做

1. 跑 `examples/a15-embedding`。  
2. 外層覆寫 `Log` 觀察行為。  

## 延伸閱讀

- Effective Go: Embedding  
