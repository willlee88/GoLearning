---
lessonId: "J2"
title: "Metrics：現在系統怎樣，一眼看數字"
description: "用原子計數做連線數、房間數、訊息速率；先提供 JSON /metrics，之後再換成 Prometheus。"
volume: "j"
order: 2
level: "l2"
status: "ready"
path_required: true
tags: ["metrics"]
example: "examples/j02-metrics"
prev: "J1"
next: "J3"
---

## 這章你會搞懂什麼

日誌告訴你「發生了什麼故事」。  
**Metrics（指標）** 告訴你「現在整體怎樣」——連線多少、房間多少、訊息多快、錯誤多兇、服務活了多久。

遊戲 Server 很需要這種「儀表板直覺」：

- 剛開活動，連線是爬升還是卡住？  
- 廣播有沒有把 `messages_out` 打爆？  
- 關房後 `rooms`／`connections` 有沒有掉下來（還是洩漏了）？

教學路線先做：**原子計數（atomic）+ JSON 的 `GET /metrics`**。  
概念對了，之後换成 Prometheus 文字格式只是輸出層的事。

讀完你要能：

1. 解釋 counter（計數器）跟 gauge（可上可下的量）差在哪  
2. 為連線／訊息加計數，並從 `/metrics` 讀回來  
3. 知道「每個玩家一個標籤」為什麼是災難

## Python 對照

| Python | Go |
|--------|-----|
| `prometheus_client` Counter／Gauge | 自己用 `atomic` + JSON，或 Prometheus client 庫 |
| Flask／FastAPI 掛 `/metrics` | `net/http` 掛 `GET /metrics` |
| 先亂加一堆 label | 一樣危險：標籤基數爆炸會打爆時序庫 |

口訣一樣：**先有正確、低基數的指標，再談漂亮儀表板。**

## 怎麼寫

最小形狀：

```go
type Metrics struct {
	Connections atomic.Int64
	Rooms       atomic.Int64
	MessagesIn  atomic.Int64
}

func (m *Metrics) snapshot() map[string]int64 {
	return map[string]int64{
		"connections": m.Connections.Load(),
		"rooms":       m.Rooms.Load(),
		"messages_in": m.MessagesIn.Load(),
	}
}
```

HTTP 端點回 JSON：

```go
mux.HandleFunc("GET /metrics", func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(m.snapshot())
})
```

在連線建立／斷開、進房／關房、收到訊息的路徑上 `Add(1)` 或 `Add(-1)`。

範例：`examples/j02-metrics`（`:8099`，另有 `POST /hit` 用來增加 `messages_in`）。

Arena Mini 的實作在 `demo/arena-mini/server/internal/metrics`，還多了 `messages_out`、`inputs_applied`、`errors`、`uptime_sec` 等欄位。瀏覽器打開：

`http://localhost:8080/metrics`

## 為什麼這樣設計／底層在幹嘛

### 為什麼不用 Mutex 包一個普通 int？

可以，但 metrics 更新非常頻繁（每則訊息、每次廣播）。  
`sync/atomic` 的 `Int64` 專門做這種「單調加減與讀取」，開銷小、也不容易寫錯成「讀改寫 race」。  
記住：atomic 適合**獨立計數**；若要「先讀一堆欄位做成一致快照再算複雜邏輯」，那是另一層設計。教學上 snapshot 時分別 `Load()` 就夠用。

### Counter vs Gauge

| 類型 | 意思 | 例子 |
|------|------|------|
| Counter | 大致只增（或重啟歸零後再增） | `messages_in`、`errors` |
| Gauge | 可上可下的當前值 | `connections`、`rooms` |

搞混的後果：你把連線數做成「只加不加」，關線後數字永遠往上——儀表板會騙你。

### 為什麼教學先 JSON，不先 Prometheus？

因為你要先建立：**量什麼、在哪加、怎麼驗證數字會動**。  
Prometheus 格式、histogram bucket、抓取間隔，都是下一步。概念錯了，換格式也救不了。

### 標籤（label）基數爆炸

進階監控常幫指標加標籤，例如 `room_id`、`weapon_type`。  
若標籤值幾乎無限（每個玩家 id、每條聊天內容），時序庫會產生海量時間序列 → 記憶體與查詢一起垮。

新手規則：

- 進程級總數：連線、房間、訊息、錯誤 → OK  
- 少量固定枚舉：phase、API 路徑群組 → 小心但通常 OK  
- 玩家 id、房間 id 當標籤 → 教學與多數遊戲後端 **不要**

要查「某一房」，用日誌的 `room` 欄位；metrics 看整機健康。

### USE／RED 這類口訣（有感覺即可）

- **USE**：Utilization、Saturation、Errors（資源視角）  
- **RED**：Rate、Errors、Duration（請求視角）  

遊戲裡常看：連線數、訊息 rate、錯誤數、tick／廣播耗時。延遲分位（histogram）進階再補。

### 進階可先略過

- Prometheus histogram／summary  
- Exemplar、exemplar 對 trace  
- 多進程 metrics 彙總與 relabel

## 遊戲 Server 會用在哪

Arena Mini 典型觀察：

1. 開兩個瀏覽器進同房 → `connections`、`rooms` 變化  
2. 開始移動 → `messages_in`／`inputs_applied`／`messages_out` 上升  
3. 關掉分頁 → `connections` 應降下來  

壓測時（J4／J8）你更該盯著這些數字，而不是只看終端有沒有噴錯。

廣播放大時（J6），`messages_out` 往往比 `messages_in` 成長得兇——這正是 metrics 要幫你抓到的現象。

## 請丟掉的舊習慣

1. 只靠日誌當監控（又慢又貴又難畫曲線）。  
2. 每個玩家／每個房間一個 metric 標籤。  
3. 加了計數卻從不在關線／關房路徑減少 gauge。  
4. 指標名稱今天英文明天中文，沒有穩定契約。

## 動手練習

### 必做

1. 跑 `examples/j02-metrics`，對 `GET http://127.0.0.1:8099/metrics`，再 `POST /hit` 幾次，看 `messages_in` 是否增加。  
2. 啟動 Arena Mini，打開 `/metrics`，開房前後對照 `connections`／`rooms`。  
3. 想一下：若要加 `messages_out`，應在「每次成功送出／廣播」的哪一點 `Add`？（Arena Mini 已有此欄，可對照原始碼。）

### 選做

1. 在自己的小伺服器加 `uptime_sec`（用啟動時間算）。  
2. 畫一張紙：列出你認為遊戲 Server 最少要有的 5 個指標，並標註哪個是 counter、哪個是 gauge。

## 常見坑

- **在熱路徑配字串做複雜格式化才更新 metrics**：指標更新要便宜；重活放到 snapshot／對外輸出。  
- **gauge 只加不減**：連線洩漏和計數 bug 看起來一樣，先確認斷線路徑有沒有 `Add(-1)`。  
- **把 `/metrics` 當公開嬉戲頁又暴露內部細節**：教學本機沒關係；生產常要網路政策或獨立 listen。  
- **沒有基線就優化**：先記「20 連線時 messages_out 大概多少」，再談改碼（接 J4）。

## 延伸閱讀

- USE／RED 方法論（搜尋關鍵字即可，先抓概念）  
- 範例：`examples/j02-metrics`  
- Arena Mini：`demo/arena-mini/server/internal/metrics`
