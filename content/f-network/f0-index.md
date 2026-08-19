---
lessonId: "F0"
title: "F 卷導讀 · 網路"
description: "從 TCP/HTTP 到 WebSocket：遊戲 Server 的連線層。"
volume: "f"
order: 0
level: "l1"
status: "ready"
path_required: true
tags: ["network"]
example: ""
prev: "C6"
next: "F1"
---

## 這章你會搞懂什麼

前面的併發（C 卷）讓你會開很多 goroutine；E 卷讓你會用時間、I/O、JSON。F 卷要把這些接到**真實網路**：有延遲、會斷線、一次 `Read` 不一定是一則完整訊息。

遊戲後端幾乎永遠是兩種面：

- **長連線**：玩家持續連著，雙向推訊息（TCP 自訂協定或 WebSocket）  
- **控制面 HTTP**：健康檢查、登入發 token、房間列表、metrics

本站 Arena Mini 的主軸是：**HTTP + WebSocket**。

## 先跟 Python 對一下

| Python | Go |
|--------|-----|
| `socket`／`asyncio.open_connection` | `net` |
| Flask／FastAPI | 先學標準庫 `net/http`，再考慮框架 |
| `websockets`／Starlette WS | `golang.org/x/net/websocket` 或 gorilla／nhooyr 等 |
| 自己拼緩衝、擔心粘包 | 一樣要做——見 F3 framing |

物理網路兩邊相同；差在 API 怎麼表達**非阻塞、超時、取消**（Go：deadline、`context`、每連線 goroutine）。

## 建議路徑

| 章 | 主題（人話） |
|----|----------------|
| F1 | 網路心智：延遲、吞吐、丟包、手感預算 |
| F2 | TCP server／client：`Listen`、`Accept` |
| F3 | 封包邊界：長度前綴，別再假設一次 Read 一則訊息 |
| F4 | `net/http`：Handler、JSON API |
| F5 | Middleware：logging／auth／recover 怎麼疊 |
| F6 | WebSocket：HTTP 升級成長連線 |
| F7 | 心跳、重連、會話（連線 ≠ 玩家） |
| F8 | Lab：WS 房間廣播（F 卷檢查點） |

做完 F8 → **G 卷（協定）**，再進遊戲房邏輯 H 卷。也請跑通 `examples/f08-ws-room` 與 `demo/arena-mini`。

## 遊戲 Server 會用在哪

```text
瀏覽器 ──WS──► Gateway/Server ──► Room（狀態、tick）
         └──HTTP──► healthz / rooms / 靜態頁
```

本卷先把「水管」接穩；訊息長什麼樣交給 G 卷，房間怎麼跑交給 H 卷。

## 動手練習

1. 掃過上表，標出你聽過／沒聽過的詞（粘包、心跳、Upgrade…）。  
2. 確認本機可編譯 Go 後，預覽一下 `examples/f02-tcp-echo` 目錄結構。  

## 常見坑

- **一上來就找「遊戲網路框架」**：先把 TCP／HTTP／WS 標準路徑走一次，除錯能力差很多。  
- **只測本機 loopback**：真機延遲與 NAT 會打臉，F1／F7 先建立預期。  

## 延伸閱讀

- 本卷順讀 **F1** 開始。  
