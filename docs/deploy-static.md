# 無 Docker 部署（純教學站）

另一台電腦**不能跑 Docker**時，只要部署**靜態教學網站**即可（閱讀課程）。  
Playground 線上跑 Go、Arena 遊戲需要後端／Docker，靜態站不會包含它們。

## 方式一：GitHub Pages（推薦）

推送到 `main` 後，GitHub Actions 會自動建置並發佈。

### 第一次請開啟 Pages

1. 打開 https://github.com/willlee88/GoLearning/settings/pages  
2. **Source** 選 **GitHub Actions**  
3. 等 Actions 跑完（或手動跑 workflow「Deploy GitHub Pages」）

### 網址

**https://willlee88.github.io/GoLearning/**

之後每次 `git push` 到 `main` 就會更新網站。對方電腦只需瀏覽器。

## 方式二：本機建置後拷貝 `dist`

在有 Node.js 的電腦：

```powershell
cd F:\GoLearning\web
npm install
npm run build
```

把 `web/dist/` 整包複製到另一台，用任一靜態伺服器開啟，例如：

```powershell
# 若那台有 Node
npx --yes serve dist -p 4321
```

或用 IIS / nginx / 檔案分享軟體提供 `dist` 目錄。  
不要直接用 `file://` 開（路徑與搜尋索引容易壞）。

> 注意：拷貝本機 `npm run build`（未設 `GITHUB_PAGES`）時，站點根路徑是 `/`。  
> GitHub Pages 建置會設 `GITHUB_PAGES=true`，根路徑是 `/GoLearning/`。

## 和 Docker 的對照

| | GitHub Pages / 靜態 dist | Docker |
|--|--------------------------|--------|
| 讀課程 | ✅ | ✅ |
| 線上跑 Go | ❌ | ✅ Playground |
| Arena 遊戲 | ❌ | ✅ |
| 對方要裝什麼 | 瀏覽器 | Docker |
