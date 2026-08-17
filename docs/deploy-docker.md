# Docker 部署指南

把 GoLearning **學習站** 與 **Arena Mini** 打包成映像，在任何已安裝 Docker 的電腦上運行。

## 你會得到什麼

| 服務 | 預設網址 | 說明 |
|------|----------|------|
| **`web`（主入口）** | **http://localhost:8088** | **完整教學站**（P0～K、路徑、搜尋、卷冊） |
| `arena` | http://localhost:8080 | **只有** Capstone 小遊戲 Arena Mini（容易誤以為是全部） |

> 若你只看到「Ready / 碰撞 / 分數」，代表開到了 **8080 遊戲**；請改開 **8088 課程站**。

對方電腦**不必**安裝 Node.js 或 Go。

## 前置

- [Docker Engine](https://docs.docker.com/engine/install/) 或 **Docker Desktop（Windows 需先啟動）**  
- 能存取專案原始碼（`git clone` 或複製整個 `GoLearning` 資料夾）  
- 確認 `docker version` 的 Server 區塊有輸出（僅 Client 表示引擎未開）

## 一鍵啟動（建議）

在專案根目錄（有 `docker-compose.yml` 的地方）：

```bash
docker compose up --build -d
```

查看狀態：

```bash
docker compose ps
docker compose logs -f
```

停止：

```bash
docker compose down
```

### 自訂埠

建立 `.env`（可選）：

```env
WEB_PORT=80
ARENA_PORT=8080
```

再 `docker compose up --build -d`。

## 只建映像、不啟動

```bash
docker compose build
docker images | findstr golearning
```

## 給其他電腦：三種方式

### 1. 整包原始碼（最簡單）

1. 複製或 `git clone` 整個 repo  
2. 在那台機器：`docker compose up --build -d`  
3. 防火牆放行 `WEB_PORT` / `ARENA_PORT`  
4. 瀏覽器開 `http://<主機IP>:8088` 與 `http://<主機IP>:8080`

### 2. 匯出映像檔（離線 / 內網）

在建置機器：

```bash
docker compose build
docker save golearning-web:latest golearning-arena:latest -o golearning-images.tar
```

在目標機器：

```bash
docker load -i golearning-images.tar
```

再放一份 `docker-compose.yml`（可改成 `image:` only、拿掉 `build:`），然後：

```bash
docker compose up -d
```

### 3. 推到映像倉庫（團隊）

```bash
docker tag golearning-web:latest  <registry>/golearning-web:latest
docker tag golearning-arena:latest <registry>/golearning-arena:latest
docker push <registry>/golearning-web:latest
docker push <registry>/golearning-arena:latest
```

目標機器改 compose 的 `image:` 後 `docker compose pull && docker compose up -d`。

## 區網存取注意

1. **綁定**：compose 已映射主機埠；確認監聽 `0.0.0.0`（Docker 預設會）。  
2. **防火牆**：Windows 需允許對應埠的入站連線。  
3. **WebSocket**：Arena 使用 `ws://主機:8080/ws`；若前面加反向代理，需開啟 Upgrade。  
4. **學習站與 Arena 不同埠**：正常；遊戲頁在 8080，課程在 8088。

## 健康檢查

```bash
curl http://127.0.0.1:8088/
curl http://127.0.0.1:8080/healthz
curl http://127.0.0.1:8080/metrics
```

## 開發用 Compose

本機改原始碼、熱重載：

```bash
docker compose -f docker-compose.dev.yml up
# Arena：
docker compose -f docker-compose.dev.yml --profile demo up
```

## 疑難排解

| 現象 | 處理 |
|------|------|
| web build 失敗找不到 content | 確認 build context 是 repo 根目錄；`.dockerignore` 未排除 `content/**/*.md` |
| arena 無法啟動 | `docker compose logs arena`；確認 `go.mod` 與原始碼完整 |
| 別台打不開 | 用主機 IP 而非 localhost；查防火牆與 Docker Desktop「允許區網」 |
| 架構不同（ARM Mac） | 目前 arena Dockerfile 預設 `GOARCH=amd64`；Apple Silicon 可改 Dockerfile 或用 `buildx` 多架構 |

### 多架構建置（可選）

```bash
docker buildx build --platform linux/amd64,linux/arm64 -f demo/arena-mini/Dockerfile -t golearning-arena:latest --push .
```

（需先設定 buildx；本機僅 amd64 時可維持預設。）

## 安全提醒

- 預設**無 TLS、無帳號**，適合內網學習。  
- 暴露公網前請加反向代理（Caddy/nginx）與 HTTPS，並限制管理埠。  
