---
lessonId: "A14"
title: "泛型：好用，但請克制"
description: "型別參數適合容器與演算法；別把房間規則寫成天書。到這裡 A 卷主路徑告一段落。"
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

## 這章你會搞懂什麼

Go 1.18+ 有**泛型（generics）**：函式或型別可以帶型別參數。它很適合消除「資料結構／演算法」的重複，例如 stack、filter、map 的 keys。

但它不適合把所有 API 都變成 `func F[T any](...)`。深刻理解包含：**約束（constraints）**、可讀性代價、以及跟 interface 怎麼分工。

## Python 對照

| Python | Go |
|--------|-----|
| `TypeVar`／結構型 Protocol | `func F[T constraints.Ordered](...)` |
| 執行期比較像「感覺有泛型」 | 編譯期檢查；實作細節另說 |
| 到處 `Any` | 到處用 `any` 當約束 ≈ 幾乎沒約束 |

## 怎麼寫

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

func (s *Stack[T]) Push(v T) {
	s.items = append(s.items, v)
}
```

範例：`examples/a14-generics`。

## 細節

### 常見約束

- `any`：什麼型別都行（也代表你幾乎沒限制）  
- `comparable`：可比較，才能當 map key、才能用 `==`  
- 自訂約束：用介面寫「一組型別」或「需要哪些方法」

### 跟介面怎麼選

- **介面**：抽象「行為」（能 Notify、能 Read）  
- **泛型**：參數化「資料長什麼型別」（Stack of T）

很多時候普通介面就夠了；別為了顯得現代硬上泛型。

### 方法不能再自己宣告一組新的型別參數

型別參數宣告在型別／函式上。這點跟某些語言不一樣，寫到編譯器報錯時先回來查這條。

### 進階可先略過

- 編譯器何時單態化、何時用字典派發。  
- 跟 `any` + 型別斷言的效能／可讀性取捨。

## 遊戲 Server 會用在哪

適合：

- 通用 `Registry[ID, Entity]`  
- 有界佇列 `Queue[Command]`  
- 小工具函式：`Filter`、`Keys`

不適合：

- 把 `Room` 規則本體做成難讀的泛型階層  
- 為了「一個函式打天下」讓錯誤訊息長到沒人敢改

## 請丟掉的舊習慣

1. 為了優雅把業務流程寫成天書泛型。  
2. 用泛型逃避把 domain 型別建模清楚。  
3. 以為「有泛型了就不需要介面」。

## 動手練習

### 必做

1. 跑 `examples/a14-generics`。  
2. 寫 `func Filter[T any](s []T, keep func(T) bool) []T`。  

### 選做

1. 用泛型做一個小的 `Set[T comparable]`。  

## 常見坑

- **錯誤訊息又臭又長**：先把約束縮小、範例變單純，再除錯。  
- **過度抽象**：新人改不了進房邏輯，抽象就失敗了。  
- **該用介面卻用泛型**（或反过来）：先問「我在抽象行為還是資料形狀」。

## A 卷檢查點（主路徑到這裡）

走到 A14，你應該能：

1. 用 package 切分可測程式  
2. 正確使用 slice／map／struct／method／interface  
3. 解釋 nil 介面陷阱與 slice 底層共享  
4. 需要時**克制**使用泛型  

下一卷 **B** 會專講 error 與 API 慣例。A15–A18 是加廣／深潛，時間夠再看。
