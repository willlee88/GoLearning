# 無 Docker 部署：靜態 HTML 版

適用：另一台電腦**不能用 Docker**、也不想依賴 GitHub Pages。  
只要**教學閱讀網站**（不含線上跑 Go、不含 Arena 遊戲）。

## 推薦：打成 ZIP 拷過去

在有 Node 的電腦（你的開發機）：

```powershell
cd F:\GoLearning
powershell -ExecutionPolicy Bypass -File .\scripts\build-static.ps1
```

會產生：

| 路徑 | 說明 |
|------|------|
| `release/GoLearning-static/` | 可直接拷貝的資料夾 |
| `release/GoLearning-static.zip` | 同上，壓縮包 |

### 在另一台電腦

1. 解壓 `GoLearning-static.zip`  
2. 雙擊 **`start.bat`**  
3. 瀏覽器開 **http://127.0.0.1:4321**

`start.bat` 會用 Python 或 Node 開一個本機靜態伺服器（**不要**直接雙擊 `index.html`，路徑會壞）。

若那台完全沒有 Python / Node：把整個資料夾放到 NAS、IIS、或其他靜態網站空間即可。

## 可選：GitHub Pages

若帳號支援 Pages，也可自動發佈（需 Actions 的 `workflow` 權限）。  
網址形態：`https://<user>.github.io/GoLearning/`  
見 repo 內 `.github/workflows/pages.yml`（若尚未推上，可手動新增）。

## 和 Docker 差在哪

| | 靜態 HTML ZIP | Docker |
|--|---------------|--------|
| 讀課程 | ✅ | ✅ |
| 線上跑 Go | ❌ | ✅ |
| Arena | ❌ | ✅ |
| 對方需求 | 瀏覽器 +（建議）Python/Node 開本地伺服器 | Docker |
