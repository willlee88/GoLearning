# Arena Mini

GoLearning Capstone — 權威即時對戰教學 demo。

## 玩法

1. 兩人同 `room`、不同 `name`
2. 雙方 **Ready** → `playing`（20Hz tick）
3. WASD / 方向鍵移動（`input`）
4. **碰撞得分**：距離 &lt; hit radius 時，較積極移動的一方 +1  
5. 先到 **3 分** 獲勝 → `ended`（約 3 秒）→ 自動回 **lobby** 可再 Ready  
6. 或 **60 秒**超時，分數高者勝

座標只由 Server 計算（權威狀態）。

## 為什麼有這個遊戲？

見課程 **P0.0**：它是 Go 主路徑的應用錨點（長連線、共享狀態、tick、權威、觀測），不是美術課。

## 架構

```text
cmd/server           HTTP + 優雅關閉
internal/hub         WS、廣播、tick loop
internal/game        純規則（可 go test）
internal/metrics     /metrics 計數
```

## 執行

```powershell
cd F:\GoLearning\demo\arena-mini\server
go mod tidy
go test ./...
go run ./cmd/server
```

- 遊戲：<http://localhost:8080/>
- 指標：<http://localhost:8080/metrics>

### 壓測

```powershell
cd F:\GoLearning\examples\j04-load-client
go run . -addr http://127.0.0.1:8080 -n 30 -seconds 15 -room load
```

## 協定 v1

| type | 說明 |
|------|------|
| `ready` | 大廳準備 |
| `input` | `"dx,dy"` ∈ -1..1 |
| `state` | 權威快照（含 score、winner、phase） |
| `ping`/`pong` | 心跳 |
| `chat` / `system` / `error` | 其他 |

## state 摘要

```json
{
  "phase": "playing",
  "tick": 42,
  "score_to_win": 3,
  "winner": "",
  "players": [{"name":"a","x":10,"y":20,"ready":true,"score":1}]
}
```
