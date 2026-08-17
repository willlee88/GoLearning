---
lessonId: "A16"
title: "any 與反射邊界"
description: "空介面的代價；reflect 何時才碰。"
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

## 本章你會建立的心智模型

`any`（`interface{}`）是「我放棄靜態資訊」。框架、JSON 解未知結構、列印除錯會碰到它。`reflect` 能在執行期看型別，但**慢、難讀、易錯**——業務規則優先用具體型別與小介面。

## Python 對照

| Python | Go |
|--------|-----|
| 到處 `Any` | 到處 `any` 同樣味道 |
| `getattr` / 動態 | `reflect`（更痛） |

## L1 能用

```go
func printAny(v any) {
	fmt.Printf("%T %#v\n", v, v)
}
```

JSON：`map[string]any` 或 `json.RawMessage` 先拆 type 再解（見 G1）。

## L2 機制

- 介面值 = 動態類型 + 動態值。  
- 型別斷言取回具體型別。  
- 反射常見於：encoding、DI 容器、ORM——你寫業務時少碰。  

## L3 深潛

- `reflect.Value` 可設性、可導出欄位規則。  
- 效能：熱路徑避免 reflect。  

## 請丟掉的 Python 習慣

1. 用 `any` 逃避建模。  
2. 運行期「有這個屬性就算」當 API 契約。  

## 遊戲 Server 連結

訊息分派：`type` + 具體 payload struct，而不是整包 `map[string]any` 算邏輯。

## 練習

### 必做

1. 把一個 `any` 參數 API 改成泛型或介面。  
2. 說明為何 Room 狀態不用 `map[string]any`。  

## 延伸閱讀

- Laws of Reflection（Go Blog）  
