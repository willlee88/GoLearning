---
lessonId: "I2"
title: "database/sql：連線池、context，還有 Repo 這層皮"
description: "標準庫資料庫介面怎麼用；為什麼不要每個請求 Open 一次；Room 熱路徑為什麼不該直接握 DB。"
volume: "i"
order: 2
level: "l2"
status: "ready"
path_required: true
tags: ["sql"]
example: "examples/i02-sql-memory"
prev: "I1"
next: "I3"
---

## 這章你會搞懂什麼

Go 標準庫的 **`database/sql`** 是一層「跟具體資料庫無關」的介面：你透過驅動（driver）連 Postgres、MySQL、SQLite……但業務碼多半只看到 `*sql.DB`、`QueryContext`、`ExecContext`。

重點先建立正確形狀：

- `*sql.DB` 代表的是 **連線池**（pool），不是「一條連線用完就丟」  
- 查詢要帶 **`context`**（上下文）：逾時、取消才能從上層一路傳下來  
- 遊戲碼裡常再包一層 **Repository（倉庫）介面**，讓 Room／服務依賴「能存玩家」，而不是依賴「正在用哪家 SQL」

本站範例 `examples/i02-sql-memory` 用**記憶體 map** 實作同一個 `PlayerRepo` 介面——這樣你不用先裝 Postgres 也能練「邊界長怎樣」。以後換真 SQL，是換實作，不是改遍 Room。

## 先跟 Python 對一下

| Python | Go |
|--------|-----|
| SQLAlchemy `Session` | `*sql.DB` + 你自己寫的 SQL／輕量包 |
| `with engine.connect() as conn:` | `QueryContext`／`ExecContext` + `defer rows.Close()` |
| `session.query(User).get(id)` | `QueryRowContext` + `Scan`，或封在 Repo 裡 |
| 依賴注入一個 `Session` factory | 注入 `PlayerRepo` 介面（測試塞 memory） |

Go 生態也有 ORM／程式碼產生器，但教材先把標準庫與介面邊界講清楚：你會比較知道 ORM「幫你藏了什麼」。

## 怎麼寫

### 1. 先定倉庫介面（port）

```go
type PlayerRepo interface {
	Upsert(ctx context.Context, p Player) error
	Get(ctx context.Context, id string) (Player, error)
}
```

誰要存檔就依賴這個介面。Room **不需要**知道後面是 map 還是 Postgres。

### 2. 記憶體實作（本站範例）

`examples/i02-sql-memory` 大致是：

- `MemoryRepo` 內含 `sync.RWMutex` + `map[string]Player`  
- `Upsert`／`Get` 先看 `ctx.Err()`，取消了就別做  
- `Get` 找不到回 `ErrNotFound`（哨兵錯誤，方便 `errors.Is`）

你在範例目錄跑測試即可驗證行為；之後若加 `ListByScore`，一樣先改介面再改實作。

### 3. 真的接上 `database/sql` 時長怎樣（概念）

```go
db, err := sql.Open("pgx", dsn) // 驅動名依你選的套件
if err != nil {
	return err
}
db.SetMaxOpenConns(20)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(time.Hour)

ctx, cancel := context.WithTimeout(parent, 3*time.Second)
defer cancel()

_, err = db.ExecContext(ctx,
	`INSERT INTO players (id, name, score) VALUES ($1, $2, $3)
	 ON CONFLICT (id) DO UPDATE SET name = $2, score = $3`,
	p.ID, p.Name, p.Score,
)
```

佔位符寫法依驅動而異（Postgres 常用 `$1`，有的用 `?`）。**不要**用字串拼接把玩家輸入黏進 SQL。

查多列時：

```go
rows, err := db.QueryContext(ctx, `SELECT id, name, score FROM players WHERE score >= $1`, min)
if err != nil {
	return err
}
defer rows.Close()

for rows.Next() {
	var p Player
	if err := rows.Scan(&p.ID, &p.Name, &p.Score); err != nil {
		return err
	}
	// ...
}
return rows.Err() // 別忘了迴圈後的錯誤
```

## 為什麼這樣設計

### 為什麼 `*sql.DB` 要長期共用？

`sql.Open` 之後，DB 物件會幫你管一池連線。每個 HTTP 請求或每個 tick 都 `Open`／`Close` 一次，等於不斷建連、拆連，又慢又容易把資料庫端的連線數打滿。

實務：**行程啟動時 Open 一次**（或經由組態建構），注入到 Repo／服務，關服時再 `Close`。

### 為什麼要 `XxxContext`？

遊戲與 API 都有逾時：客戶端斷了、閘道取消了、結算不能卡 30 秒。  
`context` 取消時，驅動有機會中止等待中的查詢，連線回池，上層才能快速失敗而不是空等。

I1 說的「冷路徑」仍然要有截止時間——慢查詢拖住 worker，佇列一樣會堆起來。

### 為什麼教學先用 memory Repo？

- 你能先把「介面、錯誤、concurrency（若多 goroutine 碰同一 store）」練完  
- CI／新手環境不必先裝 DB  
- 真 SQL 的掃描、null、migrate 可以下一層再加，不阻塞心智模型  

等你要接真庫時，再為同一個介面寫 `PostgresRepo` 即可。

### 進階可先略過

- 資料庫 schema 變更用 migrate 工具，跟應用版本一起發版。  
- `database/sql` 的 `NullString` 等型別處理 SQL `NULL`。  
- 事務：`db.BeginTx(ctx, opts)`，買道具這類「扣款＋入袋」要在同一交易。

## 遊戲 Server 會用在哪

| 資料 | 建議 |
|------|------|
| 帳號、庫存、賽季結算 | Repo → SQL（冷路徑／worker） |
| 進行中 Room 狀態 | 記憶體；**不要**讓 Room 在 tick 持有並呼叫 `*sql.DB` |
| 測試 | 塞 `MemoryRepo`，測規則不測 SQL 方言 |

HTTP 後台（商城、營運工具）可以直接在 handler 調 Repo；對局迴圈則維持 I1 的事件邊界。

## 請丟掉的舊習慣

1. **字串拼接 SQL**（SQL 注入）——參數佔位符是底線。  
2. **忽略 `rows.Close()`／`rows.Err()`**——洩漏連線或吞掉迭代錯誤。  
3. **每個請求 `sql.Open`**——池的意義被你廢掉。  
4. **把 ORM session 綁進 Room 模擬物件**——熱路徑一卡，整房一起卡。

## 動手練習

### 必做

1. 跑通 `examples/i02-sql-memory`（看 `go.mod` 所在目錄，執行套件測試）。  
2. 為 repo 加 `ListByScore(ctx, min int) ([]Player, error)`（或同等語意），補測試。  

### 選做

1. 同一介面再寫一個「故意超時」的假實作：`ctx` 取消時 `Upsert` 立刻回錯，確認呼叫端有處理。  

## 常見坑

- **`Scan` 欄位順序／型別對不上**：執行期才爆；測試要覆蓋。  
- **查完不 `Close` rows**：高併發下池被借光，看起來像「DB 掛了」。  
- **把連線池參數設無限大**：資料庫先被連線數打掛，應用還以為自己很猛。  
- **在 tick 裡同步 `QueryContext`「先做著」**：I1 白念了。

## 延伸閱讀

- <https://go.dev/doc/database/>  

## 下一章

SQL 解決「持久且可查」；接下來 Redis 解決「跨進程、要快、可帶 TTL」——但用途地圖要比「全部塞進去」精細得多。
