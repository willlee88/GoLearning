---
lessonId: "A14"
title: "泛型（克制使用）"
description: "型別參數、約束，以及何時不該上泛型。"
volume: "a"
order: 14
level: "l2"
status: "ready"
path_required: true
tags: ["generics"]
example: "examples/a14-generics"
prev: "A13"
next: "A15"
---

## 本章你會建立的心智模型

Go 1.18+ 有泛型（type parameters）。它適合**資料結構與演算法**的重複消除，不適合把所有 API 都變成 `func F[T any]`。深刻理解包含：**約束（constraints）**、可讀性代價、以及與 interface 的分工。

## Python 對照

| Python | Go |
|--------|-----|
| `TypeVar` / 結構型 Protocol | `func F[T constraints.Ordered](…)` |
| 執行期泛型感覺 | 編譯期單態化／字典派發（實作細節） |
| 到處 `Any` | 到處 `any` 約束 ≈ 沒約束 |

## L1 能用

```go
func Keys[K comparable, V any](m map[K]V) []K {
	out := make([]K, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
```

```go
type Stack[T any] struct{ items []T }

func (s *Stack[T]) Push(v T) { s.items = append(s.items, v) }
```

範例：`examples/a14-generics`。

## L2 機制

- `comparable` 用於 map key、相等比較。  
- 自訂約束：`interface { ~int | ~int64; String() string }`。  
- 方法不能再宣告新的型別參數（在 type 上宣告）。  
- 介面仍是「行為抽象」主力；泛型是「型別參數化」。

## L3 深潛（可選）

- 何時編譯器字典 vs 單態。  
- 與 `any` + 型別斷言的取捨。

## 請丟掉的 Python 習慣

1. 為了「優雅」把業務流程寫成天書泛型。  
2. 用泛型逃避把 domain 型別建模清楚。  

## 遊戲 Server 連結

適合：通用 `Registry[ID, Entity]`、有界佇列 `Queue[Command]`。  
不適合：把 `Room` 規則本身做成難讀泛型階層。

## 練習

### 必做

1. 跑 `examples/a14-generics`。  
2. 寫 `func Filter[T any](s []T, keep func(T) bool) []T`。  

### 選做

1. 用泛型做一個小的 `Set[T comparable]`。  

## 常見坑與如何看見

- 錯誤訊息又臭又長——先縮小約束再除錯。  
- 過度抽象導致新人無法改房間邏輯。  

## 延伸閱讀

- <https://go.dev/doc/tutorial/generics>  

## A 卷檢查點（M2）

到 A14 你應能：

1. 用 package 切分可測程式  
2. 正確使用 slice/map/struct/method/interface  
3. 解釋 nil interface 與 slice 共享  
4. 在需要時克制使用泛型  

下一卷 **B** 深入 error 與 API 慣例。  
