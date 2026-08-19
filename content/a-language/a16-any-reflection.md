---
lessonId: "A16"
title: "any 與反射：邊界才碰，別當日常語言"
description: "any 等於放棄靜態資訊；reflect 能動但慢又脆。業務規則優先用具體型別與小介面。"
volume: "a"
order: 16
level: "l3"
status: "ready"
path_required: false
tags: ["any", "reflect"]
example: ""
prev: "A15"
next: "A17"
---

## 這章你會搞懂什麼

`any`（以前寫成 `interface{}`）的意思很白話：**我放棄編譯期對這個值的靜態資訊**。  
框架、解未知 JSON、除錯列印會碰到它。

`reflect`（反射）能在執行期偷看／改動型別與值，但通常**慢、難讀、易錯**。業務規則請優先用具體型別與小介面；這章教你認邊界，不是叫你愛上它。

## Python 對照

| Python | Go |
|--------|-----|
| 到處標 `Any` | 到處用 `any` 味道一樣糟 |
| `getattr`／高度動態 | `reflect`（往往更痛） |

Python 習慣動態屬性沒關係；Go 選靜態的理由之一，就是想把這類爆炸留在少數邊界。

## 怎麼寫

```go
func printAny(v any) {
	fmt.Printf("%T %#v\n", v, v)
}
```

處理未知 JSON 時，常見是 `map[string]any`，或先留 `json.RawMessage`，讀到 `type` 再解具體 struct（見 G 卷信封）。

需要拿回具體型別時，用 A13 的型別斷言／type switch，不要一開始就反射。

## 細節

### 介面值還是那兩欄

動態類型 + 動態值。`any` 只是「方法集為空」的介面，所以什麼都能裝，也什麼都不保證。

### 反射常見藏身處

- `encoding/json`、類似編碼套件  
- 少數 DI／ORM／產生碼工具  

你自己寫房間規則、匹配、戰鬥结算時，**少碰 reflect** 是美德。

### 為什麼熱路徑討厭它

反射多半走執行期查找，除錯堆疊也比較不友善。先把結構用型別表達清楚，通常又快又穩。

### 進階可先略過

- `reflect.Value` 的可設性、只能動匯出欄位等規則。  
- 真要動態，優先考慮產生碼，而不是手寫反射魔術。

## 遊戲 Server 會用在哪

訊息分派：`type` + 具體 payload struct，在規則層算邏輯。  
不要整包 `map[string]any` 一路算到血量——欄位拼錯只能執行期才知道，超適合半夜爆炸。

## 請丟掉的舊習慣

1. 用 `any` 逃避建模。  
2. 運行期「有這個 key／屬性就算」當正式 API 契約。  
3. 看到 json 就反射到天亮。

## 動手練習

### 必做

1. 把某個 `any` 參數 API 改成泛型或小介面。  
2. 用自己的話說明：為什麼 Room 狀態不該用 `map[string]any` 當主力。  

### 選做

1. 讀一段 `encoding/json` 文件，列出它幫你「藏起反射」的好處。  

## 常見坑

- **`map[string]any` 傳很深**：型別斷言地獄。  
- **反射改未匯出欄位**：失敗或只能在同 package 玩，別當跨層魔法。  
- **測試通過、線上 JSON 多一個欄位就歪**：契約要用 struct＋測試釘住。
