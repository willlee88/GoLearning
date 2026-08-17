---
lessonId: "I4"
title: "快取一致性直覺"
description: "旁路快取、TTL、與遊戲可接受的過期。"
volume: "i"
order: 4
level: "l2"
status: "ready"
path_required: true
tags: ["cache"]
example: ""
prev: "I3"
next: "I5"
---

## 本章你會建立的心智模型

快取會**錯**。問題是「錯多久、錯多大」。

| 模式 | 說明 |
|------|------|
| Cache-aside | 讀 miss 載入；寫時刪/更快取 |
| TTL | 簡單，允許短暫舊資料 |
| 寫穿 | 寫同時更快取，仍有競態 |

遊戲設定表可長 TTL；錢包餘額要強一致或版本號。

## Python 對照

Django cache / functools.lru_cache 同權衡。

## L1 能用

```text
讀: cache.get → miss → db → cache.set(ttl)
寫: db.write → cache.delete(key)
```

## L2 機制

- 驚群：miss 時單飛（singleflight）。  
- 版本號／CAS 更新。  
- 不要快取「正在進行的 Room 權威狀態」當跨服真相（除非精心設計）。  

## 請丟掉的 Python 習慣

1. 無限 TTL 快取使用者資產。  
2. 更新 DB 卻忘了失效快取。  

## 遊戲 Server 連結

角色外觀可快取；對局血量以 Room 為準。

## 練習

### 必做

1. 選兩個欄位：一個可 TTL 60s，一個必須強一致，說明原因。  

## 延伸閱讀

- singleflight 套件概念  
