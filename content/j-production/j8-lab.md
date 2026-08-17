---
lessonId: "J8"
title: "Lab：關得乾淨 + 壓一下"
description: "M5 檢查點：metrics、shutdown、load-client。"
volume: "j"
order: 8
level: "l2"
status: "ready"
path_required: true
tags: ["lab", "production"]
example: "examples/j04-load-client"
prev: "J7"
next: "K0"
---

## 本章你會建立的心智模型

把 I/J 落到操作：

1. 啟動 Arena Mini M5  
2. 看 `/healthz` `/metrics`  
3. 跑 load-client  
4. Ctrl+C 觀察優雅關閉日誌  
5. （選）pprof 自行加掛  

## L1 能用

```powershell
cd F:\GoLearning\demo\arena-mini\server
go run ./cmd/server

# 另一終端
cd F:\GoLearning\examples\j04-load-client
go run . -addr http://127.0.0.1:8080 -n 30 -seconds 15
```

## L2 檢查清單

- [ ] SIGINT 後進程退出碼 0 或明確  
- [ ] metrics 反映連線變化  
- [ ] 關閉時 tick 停止  
- [ ] 無界 goroutine 洩漏（反覆開關壓測觀察）  

## 練習

### 必做

1. 完成上述流程並截一段 metrics JSON。  
2. 寫下你觀察到的第一個瓶頸猜測。  

### 選做

1. 為 input 加每秒速率限制。  
2. Ended 時丟結算到 worker（I5）。  

## 延伸閱讀

- `demo/arena-mini/README.md`  

## 往 K 卷

Capstone 複習與擴充任務清單見 K0。  
