# 無 Docker：用 git 看靜態 HTML

適用：另一台電腦**不想裝 Docker**，但可以用 git。  
只要**教學閱讀網站**（不含線上跑 Go、不含 Arena 遊戲）。

## 推薦：clone / pull + `看課.bat`

倉庫已提交建置好的 HTML（`site/`），閱讀端流程：

```powershell
git clone https://github.com/willlee88/GoLearning.git
cd GoLearning
.\看課.bat
```

瀏覽器開 **http://127.0.0.1:4321**。

之後只要：

```powershell
git pull
```

再開一次 `看課.bat`（或伺服器還開著就直接重新整理頁面）。  
**不必**再去 Releases 下載 ZIP。

`看課.bat` 會用 Python 或 Node 開本機靜態伺服器（**不要**直接雙擊 `index.html`，路徑會壞）。

### 需求

| 項目 | 說明 |
|------|------|
| git | clone / pull |
| Python 3 或 Node.js | 只為了開本機小伺服器 |
| 都沒有？ | 把 `site/` 丟到 NAS / IIS / 任何靜態網站空間 |

## 作者端：改課文後怎麼更新 HTML

在有 Node 的開發機：

```powershell
cd F:\GoLearning
powershell -ExecutionPolicy Bypass -File .\scripts\build-static.ps1
git add content/ site/
git commit -m "Update lessons + static site"
git push
```

腳本會寫入：

| 路徑 | 進 git？ | 說明 |
|------|----------|------|
| `site/` | ✅ | 給閱讀端 pull 用 |
| `release/GoLearning-static.zip` | ❌（gitignore） | 可選，給不走 git 的人 |

**記得 `content/` 與 `site/` 一起 commit**，否則別人 pull 到的 HTML 還是舊的。

## 可選：仍用 ZIP 拷貝

不方便裝 git 時：

1. 開發機跑 `scripts/build-static.ps1`
2. 把 `release/GoLearning-static.zip` 拷過去
3. 解壓 → `start.bat` → http://127.0.0.1:4321

也可直接下載 GitHub Release 的 zip（較舊時請以 `site/` + git 為準）。

## 推薦線上：GitHub Pages（連 start.bat 都不用）

網址：**https://willlee88.github.io/GoLearning/**

**push 到 `main` 就會自動更新**（workflow：`.github/workflows/static.yml`）。

流程：checkout → `web/` 用 `GITHUB_PAGES=true` 建置 Astro → 上傳 `web/dist` → Deploy Pages。

作者端只要：

```powershell
# 可選但建議：同步本機 site/
powershell -ExecutionPolicy Bypass -File .\scripts\build-static.ps1
git add content site
git commit -m "Update lessons"
git push
```

等 Actions 變綠燈（約 1～2 分鐘）後重新整理網站即可。

備用（通常不用）：`scripts/deploy-pages.ps1` 推到 `gh-pages` 分支。  
正確的 workflow 範本也備份在 `docs/github-pages-workflow.example.yml`。
## 和 Docker 差在哪

| | `site/` + git | Docker |
|--|---------------|--------|
| 讀課程 | ✅ | ✅ |
| 更新方式 | `git pull` | `git pull` + compose 重建（或 volume） |
| 線上跑 Go | ❌ | ✅ |
| Arena | ❌ | ✅ |
| 對方需求 | 瀏覽器 + Python/Node（開伺服器） | Docker |
