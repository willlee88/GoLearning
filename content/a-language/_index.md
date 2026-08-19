---
lessonId: "A0"
title: "先看這卷要學什麼（語言地基）"
description: "A 卷在幹嘛：從安裝、package、零值，一路到 slice／map／struct，幫你把 Go 的資料與型別想清楚。"
volume: "a"
order: 0
level: "l1"
status: "ready"
path_required: true
tags: ["language"]
example: ""
prev: "P0.8"
next: "A1"
---

## 這章你會搞懂什麼

A 卷是整條路的**語言地基**。目標不是把 API 背起來，而是學完後你能用自己的話回答：

1. **資料放哪、怎麼共享**——值拷貝、指標、slice 底層陣列差在哪  
2. **誰能看見誰**——package 的大寫匯出（export）規則  
3. **編譯器幫你守什麼契約**——型別、方法、後面會碰到的介面  

簡單講：**A1–A14** 是主路徑，必走；**A15–A18** 是加廣／深潛（嵌入、any、unsafe、逃逸），有空再啃也沒關係。

## Python 對照

| 你熟悉的 | 在 A 卷會變成 |
|----------|----------------|
| 動態型別 + `None` | 靜態型別 + 零值（zero value）／`nil` |
| `list` / `dict` | slice / map（語意更貼「底層怎麼存」） |
| class 方法 | struct + method（值或指標 receiver） |
| 模組與 `__all__` 約定 | 名稱大寫＝匯出，編譯器直接擋 |

每一章都會再細講一次對照，這裡先建立「整卷地圖」。

## 怎麼寫

這卷沒有單一程式要寫。建議你照這個順序走，每章至少跑一次範例：

| 順序 | 章 | 你大概在練什麼 |
|------|----|----------------|
| 1 | A1–A3 | 安裝、module、package、零值與常數 |
| 2 | A4–A7 | 控制流、函式、defer、指標 |
| 3 | A8 | slice／map（本卷檢查點，很重要） |
| 4 | A9–A11 | 字串、struct、method |
| 5 | A12–A14 | interface、型別斷言、泛型（克制用） |

做完 A 卷，下一步是 **B（錯誤怎麼當一等公民）→ C（併發）**。

## 細節

為什麼要先把「資料與型別」講透？因為後面房間、連線、廣播全是：

- 一堆 struct（玩家、房間、指令）
- 一堆 map／slice（查人、排隊、快照）
- 偶爾指標與介面（共享狀態、可替換依賴）

你要是這裡含糊，C 卷以後的 race、鎖、channel 會更痛苦——不是語法難，是**搞不清楚誰擁有哪份資料**。

## 遊戲 Server 會用在哪

A 卷結束後，你應該能乾淨地表達這些「細胞」：

- 玩家實體（Player）
- 房間裡用 ID 查人的 map
- 廣播或存檔用的 slice buffer／快照

這些就是後面 Session、Room、tick 的原料。Arena Mini（`demo/arena-mini`）裡你也會一直看到它們。

## 請丟掉的舊習慣

1. 「先把語法掃過、細節以後再說」——A8 的 slice 共享、A12 的 nil 介面，拖到後面會連環炸。  
2. 把 Go 當成「比較快的 Python」硬套 class／例外心智。  
3. 只看文章不跑 `go run`／`go test`——手感要靠編譯器罵你建立。

## 動手練習

### 必做

1. 順讀 A1→A8，每章至少跑一個 `go run` 或 `go test`。  
2. 完成 `examples/a08-player-registry`（A8 檢查點）。  

### 選做

1. 把 P0 的 `p0-config-stats` 改寫成多 package（呼應 A2）。  

## 常見坑

- 跳章只看「有趣的」介面／泛型，結果 slice／指標基礎不穩。  
- 範例路徑搞錯目錄：一定要在有 `go.mod` 的資料夾跑指令。  
- 以為「看懂文字＝會了」——碰到編譯錯誤才是真正開始學。
