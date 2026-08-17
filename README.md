# GoLearning

給 **學過 Python** 的開發者：廣覆蓋、可深讀的 Go 學習網站，並以 **遊戲 Server** 為應用主軸。

| 項目 | 說明 |
|------|------|
| 路徑 | `F:\GoLearning` |
| 規劃書 | [`docs/規劃書.md`](docs/規劃書.md) |
| 內容策略 | 廣覆蓋 · L1/L2/L3 分層 · Python 對照 · 遊戲情境 |

## 快速開始

### 前置需求

- **Node.js** ≥ 22.12（已用於學習網站）
- **Go** ≥ 1.22（執行 `examples/`、`demo/`；尚未安裝可先只看網站）

安裝 Go（Windows）：<https://go.dev/dl/>

### 啟動學習網站

```powershell
cd F:\GoLearning\web
npm install
npm run dev
```

瀏覽器開啟終端機顯示的本機網址（通常是 `http://localhost:4321`）。

### 執行範例（需已安裝 Go）

```powershell
cd F:\GoLearning\examples\p0-config-stats
go test ./...
go run .
```

### Arena Mini Demo（需 Go）

```powershell
cd F:\GoLearning\demo\arena-mini\server
go run .
# 另一個終端
# 開啟 demo\arena-mini\web\index.html 或之後由 compose 提供
```

### Docker 打包部署（給本機或其他電腦）

```powershell
cd F:\GoLearning
docker compose up --build -d
```

| 服務 | 網址 |
|------|------|
| 學習站 | http://localhost:8088 |
| Arena Mini | http://localhost:8080 |

- 對方電腦只需 **Docker**，不必裝 Node/Go。  
- 完整說明：[`docs/deploy-docker.md`](docs/deploy-docker.md)  
- 開發熱重載：`docker compose -f docker-compose.dev.yml up`

## 倉庫結構

```text
GoLearning/
├── docs/                 # 規劃書、ADR
├── content/              # 課程 Markdown（真相來源）
├── web/                  # Astro 學習網站
├── examples/             # 與章節對應的可執行範例
├── samples/              # 廣覆蓋精簡樣本
├── demo/arena-mini/      # Capstone 遊戲 Server
├── scripts/              # 內容檢查等
├── docker-compose.yml
└── Makefile
```

## 學習路徑（主路徑）

1. **P0** Python → Go 心智遷移  
2. **A–E** 語言、錯誤、併發、工程化、標準庫  
3. **F–G** 網路與協定  
4. **H–J** 遊戲 Server、資料、生產化  
5. **K** Capstone：Arena Mini  

詳見網站「學習路徑」或 [`content/paths/main.yaml`](content/paths/main.yaml)。

## 章節寫作約定

見 [`content/_templates/lesson.md`](content/_templates/lesson.md)。

每章 frontmatter 必填：`id`、`title`、`volume`、`level`、`status`。

## 目前進度

### M0
- [x] monorepo 骨架、Astro 站、P0 全卷、p0-config-stats、Arena Mini 骨架  

### M1
- [x] A0–A8、高亮渲染、站內搜尋  

### M2
- [x] A9–A14、B0–B4、C0–C6 + race lab  

### M3
- [x] F0–F8、G0–G3、網路範例、Arena Mini 信封  

### M4
- [x] H0–H10、Arena Mini 權威 tick  

### M5
- [x] I/J/K + Arena Mini 生產化基線  

### 廣覆蓋 / 玩法
- [x] P0.0、A–K 主內容、Arena 碰撞得分  

### Docker 打包
- [x] `web/Dockerfile`（Astro → nginx）  
- [x] `demo/arena-mini/Dockerfile`（靜態 binary）  
- [x] `docker compose up --build` 生產向  
- [x] `docs/deploy-docker.md`  
- [ ] 推到公開 registry…（可選）  

## 授權

暫為本機學習專案；開源授權待定。
