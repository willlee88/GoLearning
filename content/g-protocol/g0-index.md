---
lessonId: "G0"
title: "G 卷導讀 · 序列化與協定"
description: "訊息形狀、版本化、命令與事件分離。"
volume: "g"
order: 0
level: "l1"
status: "ready"
path_required: true
tags: ["protocol"]
example: ""
prev: "F8"
next: "G1"
---

## 這章你會搞懂什麼

F 卷把**水管**接好了（TCP／HTTP／WebSocket）。但水管裡流的是什麼？

**協定（protocol）**就是雙方的約定：訊息長什麼樣、誰可以說什麼、版本升級怎麼辦。協定寫得好，你才方便除錯、擴欄位、防客戶端作弊式亂送。

G 卷聚焦四件事（本批先做前三）：

1. JSON **信封**長什麼樣、成本與好處  
2. **命令**（client→server）vs **事件／狀態**（server→client）  
3. **版本欄位**與相容策略  
4. （概念）以後若要上 Protobuf 等二進位，心智怎麼接

## 先跟 Python 對一下

| 你在 Python 常做的 | 在 Go／本卷 |
|--------------------|-------------|
| pydantic 模型當 API schema | struct + `encoding/json` + 手動驗證 |
| REST `/v1/...` 版本化 | 即時訊息常在 JSON 內帶 `"v": 1` |
| 客戶端算完結果 POST 存檔 | 遊戲要**權威伺服器**：客戶端只送意圖 |

## 章節地圖

| 章 | 主題 |
|----|------|
| G1 | JSON 信封實務（`v`／`type`／`payload`） |
| G2 | 命令與事件分離 |
| G3 | 版本化與相容；並作為進 H 卷前檢查點 |

## 遊戲 Server 會用在哪

```text
Client  --command-->  Server(validate/apply)  --event/state-->  Clients
```

Arena Mini 的 `move`／`chat`／`state` 都該能對上這個圖。協定穩了，H 卷的 Room／tick 才不會跟「字串長什麼樣」纏成一團。

## 動手練習

1. 打開 F8／Arena 的一則實際訊息，標出哪些是「意圖」、哪些是「伺服器認定的結果」。  
2. 順讀 G1→G3，跑 `examples/g01-envelope`。  

## 常見坑

- **每種訊息頂層形狀都不同**：沒有 `type`，分派與日誌都會痛。  
- **先追求二進位極限優化**：教學與早期產品用 JSON 通常更賺除錯時間。  

## 延伸閱讀

- 下一章 **G1**。  
