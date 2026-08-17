---
lessonId: "J5"
title: "安全基線"
description: "輸入驗證、速率限制、秘密與信任邊界。"
volume: "j"
order: 5
level: "l2"
status: "ready"
path_required: true
tags: ["security"]
example: ""
prev: "J4"
next: "J6"
---

## 本章你會建立的心智模型

客戶端**不可信**。基線：

1. 驗證所有命令（範圍、phase、身份）  
2. 限制訊息大小與速率  
3. 認證 token（生產）  
4. TLS 終端  
5. 不在日誌打秘密  

權威 Server 是反作弊第一層，不是全部。

## Python 對照

Pydantic 驗證 + rate limit middleware 同層。

## L1 能用

- 最大幀 / 最大 JSON  
- 每連線每秒 N 個 input  
- 暱稱長度與字元集  

## L2 機制

- 房號枚舉攻擊 → 隨機 id / 權限。  
- Origin 檢查（瀏覽器 WS）。  
- 管理 API 與玩家口分離。  

## 請丟掉的 Python 習慣

1. 前端隱藏的「管理員按鈕」當安全。  
2. 預設 debug 開在公網。  

## 遊戲 Server 連結

Arena Mini 已拒絕非法 input；速率限制可作練習。

## 練習

### 必做

1. 提出 3 個你會加的 rate limit。  
2. 說明為何 `pos` 直傳座標危險（呼應 H4）。  

## 延伸閱讀

- OWASP API Security（概念）  
