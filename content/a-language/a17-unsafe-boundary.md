---
lessonId: "A17"
title: "unsafe 邊界（概覽）"
description: "知道它存在、知道平時不該用。"
volume: "a"
order: 17
level: "l3"
status: "ready"
path_required: false
tags: ["unsafe"]
example: ""
prev: "A16"
next: "A18"
---

## 本章你會建立的心智模型

`unsafe` 能打破型別與記憶體安全假設，用於極致優化或與 C 互通。教學要求是：**能認出、能拒絕在業務裡亂用不必要的 unsafe**。遊戲 Server 熱路徑優化前先 profile。

## Python 對照

類似 ctypes／直接玩記憶體——多數 web 後端不需要。

## L1 能用

記住：

- 標準庫與編譯器保證在 `unsafe` 外較強。  
- 一用 `unsafe`，可攜性、版本升級、race 推理都變難。  
- Code review 紅旗。  

## L2 機制

- `Pointer`、對齊、生命週期（勿保留已失效指標）。  
- `string`/`[]byte` 零拷貝轉換是經典（也危險）例子——優先可讀正確版。  

## 請丟掉的 Python 習慣

1. 「快一點就好」無測量就上黑魔法。  

## 遊戲 Server 連結

狀態同步優化順序：減少廣播 → delta → 二進位 →（最後才）unsafe。

## 練習

### 必做

1. 搜尋你依賴裡誰用了 unsafe（可選）。  
2. 寫下三條「禁止在 Room 規則使用 unsafe」的理由。  

## 延伸閱讀

- `unsafe` package doc  
