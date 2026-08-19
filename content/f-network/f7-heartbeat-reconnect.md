---
lessonId: "F7"
title: "心跳、重連與會話"
description: "連線≠玩家；閒置檢測與重連綁定概念。"
volume: "f"
order: 7
level: "l2"
status: "ready"
path_required: true
tags: ["session", "websocket"]
example: "examples/f08-ws-room"
prev: "F6"
next: "F8"
---

## 這章你會搞懂什麼

手機進電梯、筆電蓋上、咖啡店 Wi‑Fi 抽風——**連線會斷**。若你的設計是「socket 斷＝玩家毀滅」，體驗會很脆。

請在腦子裡拆成兩層：

1. **連線（Connection）**：當下這條 TCP／WebSocket  
2. **會話（Session）**：玩家是誰、能不能在短時間內恢復  

再加上：

- **心跳（heartbeat／ping）**：定期探活，找出「看起來還連著但其實假死」的連線  
- **重連**：用 session id／token 綁回同一玩家，而不是變成另一個匿名路人

## 先跟 Python 對一下

概念完全一樣。Python asyncio 常用 cancel scope／背景 task 做 ping；Go 用 `time.Ticker` + `context` 取消。重連與冪等是遊戲／即時系統共通難題，不是某一語言獨有。

## 怎麼寫（能跑的最小概念）

```go
// 概念：ping loop；與讀 loop 共享寫鎖
ticker := time.NewTicker(15 * time.Second)
defer ticker.Stop()
for {
	select {
	case <-ticker.C:
		_ = send(conn, pingMsg) // 例如 {"type":"ping"}
	case <-ctx.Done():
		return
	}
}
```

讀側：

- 每次收到訊息就刷新「最後活動時間」  
- 或設 `SetReadDeadline`：多久沒讀到東西就當斷線清理  

應用層 `ping`／`pong` 的好處是：**你在 log 與協定裡都看得到**，比只靠 TCP keepalive 好除錯。

## 為什麼這樣設計／底層在幹嘛

1. **TCP keepalive 不夠靠譜當唯一手段**  
   中間的 NAT／負載平衡器常有**閒置逾時**（例如雲廠商 LB 几十秒到幾分鐘）。應用層心跳能主動撐住或更快發現死連線。

2. **連線 ≠ 玩家**  
   同一玩家可能：斷線 → 3 秒內重連 → 應回到同一座位。實作上 server 保留 session 一段 **TTL**（存活時間），過期才清狀態。

3. **重送與冪等**  
   重連後客戶端可能重送「開火」「購買」。伺服器要有序號或去重，否則一槍變兩槍（G／H 卷會再碰到）。

4. **局內能不能恢復看遊戲**  
   聊天室：容易。競技對戰：可能要暫停、觀戰位、或乾脆重開。先**設計表**再寫碼。

5. **安全**  
   session token 要難以猜；不要只拿「暱稱」當身份。HTTPS／WSS 傳 token。

## 遊戲 Server 會用在哪

Arena Mini 可以先預留：

- 訊息型別 `ping`／`pong`  
- session id 欄位  
- 短暫停牌、重連回房（進階分支）  

教學 M3 階段：先把概念與訊息型別想清楚，不必一次做滿。

## 請丟掉的舊習慣

1. **斷線立刻清光玩家，不留任何恢復窗口。**  
2. **完全依賴 TCP keepalive，不做應用心跳。**  
3. **重連只看暱稱字串**——撞名與假冒。  

## 動手練習

### 必做

1. 在 `f08-ws-room` 或 arena-mini 加 `type=ping`／`pong`，兩邊印 log。  
2. 設計一張表：哪些狀態可重連恢復、哪些必須重開局（紙上即可）。  

### 選做

1. 讀 deadline：90 秒沒任何訊息就踢；有 ping 就延期。  

## 常見坑

- **心跳太勤**：浪費電與流量；太疏：假死發現太慢。15–45 秒是常見討論區間，依 LB 閒置時間調。  
- **只 ping 不踢**：server 還掛著一堆死 session。  
- **重連開新房間座位**：變成「分身」。  
- **忽略雲廠商文件的 idle timeout**：本機測永遠好好的。

## 延伸閱讀

- 你使用的雲／反向代理文件中，搜尋 WebSocket idle／timeout  
- 接續 **F8** lab 把連線與房間清收取齊  
