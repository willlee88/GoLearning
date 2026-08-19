# GoLearning

給 **學過 Python** 的開發者：廣覆蓋、可深讀的 Go 學習網站，並以 **遊戲 Server** 為應用主軸。

| 項目 | 說明 |
|------|------|
| 路徑 | `F:\GoLearning` |
| 規劃書 | [`docs/規劃書.md`](docs/規劃書.md) |
| 內容策略 | 廣覆蓋 · L1/L2/L3 分層 · Python 對照 · 遊戲情境 |

## 怎麼看課（**不用 Docker**）

### 最省事：直接開網頁（GitHub Pages）

**https://willlee88.github.io/GoLearning/**

改完課文並跑完下方「作者端」發佈後，重新整理即可，不用裝東西、不用 `start.bat`。

### 本機離線：`git pull` + `看課.bat`

倉庫裡也有建好的靜態站：`site/`。

```powershell
git clone https://github.com/willlee88/GoLearning.git
cd GoLearning
.\看課.bat
```

瀏覽器開 **http://127.0.0.1:4321**。之後：`git pull` → 重新整理（或再跑一次 `看課.bat`）。

> 靜態版只有**課程閱讀**。線上跑 Go／Arena 仍要 Docker（可選）。  
> 詳見 [`docs/deploy-static.md`](docs/deploy-static.md)。

### 作者端：改完課文要更新 HTML

```powershell
# （建議）同步本機 site/，給 git pull / 離線看課用
powershell -ExecutionPolicy Bypass -File .\scripts\build-static.ps1
git add content site
git commit -m "Update lessons"
git push
```

**push 到 `main` 後，GitHub Actions 會自動建置並更新**  
https://willlee88.github.io/GoLearning/  
（約 1～2 分鐘；可在 repo 的 Actions 分頁看進度）

> 備用手動發佈（通常不用）：`scripts\deploy-pages.ps1`
## 快速開始（本機開發）

### 前置需求

- **Node.js** ≥ 22.12（建置／預覽網站）
- **Go** ≥ 1.22（可選：跑 examples／demo）

### 啟動學習網站

```powershell
cd F:\GoLearning\web
npm install
npm run dev
```

瀏覽器開終端機顯示的網址（通常 `http://localhost:4321`）。

### Docker 全套（可選：Playground + Arena）

```powershell
cd F:\GoLearning
docker compose up --build -d
```

| 服務 | 網址 |
|------|------|
| 學習站 + 線上跑 Go | http://localhost:8088 |
| Arena Mini | http://localhost:8080 |

詳見 [`docs/deploy-docker.md`](docs/deploy-docker.md)。

## 倉庫結構

```text
GoLearning/
├── 看課.bat              # 閱讀端：開 site/ 本機伺服器
├── site/                 # 已建置的靜態 HTML（進 git，pull 就能看）
├── docs/                 # 規劃書、ADR、部署說明
├── content/              # 課程 Markdown（真相來源）
├── web/                  # Astro 原始碼（改版面／建置用）
├── examples/             # 與章節對應的可執行範例
├── samples/              # 廣覆蓋精簡樣本
├── demo/arena-mini/      # Capstone 遊戲 Server
├── scripts/              # build-static 等
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
