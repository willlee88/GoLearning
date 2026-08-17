# ADR 0001 — M0 技術選型

- 狀態：Accepted
- 日期：2026-08-13

## 決策

- 網站：Astro 7 + Tailwind 4
- 課程內容：repo 根目錄 `content/**/*.md`（gray-matter frontmatter）
- 第一個範例：`examples/p0-config-stats`
- Demo：`demo/arena-mini`（標準庫 HTTP + `golang.org/x/net/websocket`）
- 內容策略：廣覆蓋、Python 背景、L1/L2/L3（見規劃書 v0.2）

## 後果

- 未安裝 Go 時仍可瀏覽網站與閱讀 Markdown。
- WebSocket 選用 x/net 以減少框架噪音；日後可換實作但不改協定 JSON 形狀。
