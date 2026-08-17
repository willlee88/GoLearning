---
lessonId: "G3"
title: "版本化與相容策略"
description: "v 欄位、容錯解碼、棄用路徑。"
volume: "g"
order: 3
level: "l2"
status: "ready"
path_required: true
tags: ["protocol"]
example: "examples/g01-envelope"
prev: "G2"
next: "H0"
---

## 本章你會建立的心智模型

客戶端不會同時升級。協定要：

1. **顯式版本** `v`  
2. **向後相容**：舊客戶端可玩或可被優雅拒絕  
3. **棄用窗口**：雙支援 → 告警 → 移除  

## Python 對照

API versioning（`/v1`）同一思想；即時協定更常在訊息內帶 `v`。

## L1 能用

```go
if env.V != 1 {
	return fmt.Errorf("unsupported version %d", env.V)
}
```

或：v1/v2 不同 handler。

## L2 機制

- 只加可選欄位通常安全；改語意要 bump v。  
- 伺服器可同時接受 v1+v2，轉換到內部模型。  
- 文件化「最低客戶端版本」。  

## 請丟掉的 Python 習慣

1. 默默改 JSON 形狀上線。  
2. 用「還沒人用」當相容策略。  

## 遊戲 Server 連結

App Store 審核延遲 → 強制雙版本視窗更長。

## 練習

### 必做

1. 為 envelope 拒絕 `v!=1` 並回 `{type:"error",...}`。  
2. 寫一頁 CHANGELOG：假設要加 `payload.skin`。  

## 延伸閱讀

- 語意化版本（概念）  

## M3 檢查點

你應能：

1. 解釋 TCP 流與 framing  
2. 寫 HTTP health + 簡單 API  
3. 跑 WS 多房間廣播  
4. 用 JSON 信封區分命令／事件  

下一階段 **H 卷**：Room tick 與權威狀態（Arena Mini 對戰化）。  
