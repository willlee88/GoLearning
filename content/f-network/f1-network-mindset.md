---
lessonId: "F1"
title: "網路心智模型"
description: "延遲、吞吐、丟包與遊戲手感；可靠傳輸不等於即時。"
volume: "f"
order: 1
level: "l2"
status: "ready"
path_required: true
tags: ["network", "game-server"]
example: ""
prev: "F0"
next: "F2"
---

## 本章你會建立的心智模型

網路不是「函式呼叫比較慢」，而是：**延遲、抖動、吞吐、丟包、重傳、中間盒**。TCP 保證有序可靠，但可能阻塞在隊頭（HOL）；UDP 快但不保證到達——本站實作主軸仍是 **TCP/WebSocket + 權威狀態**，UDP 只在概念對照。

## Python 對照

兩邊物理相同。差異在你如何用 API 表達非阻塞、超時與取消（Go：`context`、deadline）。

## L1 能用

記住三個數字直覺：

| 情境 | 量級（示意） |
|------|----------------|
| 同機 loopback | &lt; 1ms |
| 同城 RTT | 數 ms～數十 ms |
| 跨洲 RTT | 100ms+ |

遊戲 tick 20–60Hz 時，網路 RTT 會直接吃掉「手感預算」。

## L2 機制

- **延遲 vs 頻寬**：小狀態同步常是延遲／訊息率問題，不是把檔案傳完。  
- **廣播放大**：N 人互相全量同步 ≈ O(N²) 訊息風險。  
- **超時必設**：沒有 deadline 的 `Read` 是生產事故。  
- **背壓**：客戶端慢會反壓 Server 緩衝——要有界佇列。

## L3 深潛（可選）

- Nagle、delayed ACK 與小包遊戲。  
- 行動網路的 NAT 與連線保活。

## 請丟掉的 Python 習慣

1. 預設阻塞讀、不設 timeout。  
2. 認為「TCP 不會丟」所以不必做應用層重連與冪等。  

## 遊戲 Server 連結

權威 Server 選擇常因「可預期狀態」而非「最速 UDP」。先把 TCP/WS 做對，再談預測與補償（H 卷）。

## 練習

### 必做

1. 用文字估計：30Hz、50 人同房、每 tick 200B 全量廣播，粗算每秒從 Server 出去的流量。  
2. 列出你會為連線設定的三種 timeout（握手／讀／閒置）。  

## 延伸閱讀

- RFC 793（TCP，可略讀概念）  
