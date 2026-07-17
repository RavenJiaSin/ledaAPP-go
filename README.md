# YOLOv8 ONNX 推論服務（Go）

## 一、專案簡介

本專案是一個基於 Go + ONNX Runtime 的影像辨識推論系統，使用 YOLOv8 模型進行物件偵測，並透過 HTTP API 提供單張圖片推論服務。

系統目前已完成「單張圖片推論完整流程」，並具備可擴展為即時串流（RTSP / Camera pipeline）的架構雛形。

---

## 二、系統架構

整體推論流程如下：
```
Client (HTTP Upload Image)
↓
HTTP Server (server)
↓
Image Decode (JPEG)
↓
Inference Pipeline
├─ Preprocess (Letterbox)
├─ ONNX Runtime Inference
├─ Postprocess (Decode YOLOv8)
├─ NMS (去重)
└─ Clamp (邊界修正)
↓
Detection Result
↓
JSON Response
```


---

## 三、核心模組說明

### 1. inference（推論核心）
- ONNX Runtime CGO 封裝
- 負責模型載入與推論執行
- 提供 `Run(input, shape)` API

---

### 2. pipeline（推論流程控制）
- 管理整個 AI 推論流程
- 包含：
  - 前處理（Letterbox）
  - 模型推論（ONNX）
  - 後處理（Decode + NMS）
- 輸出標準化 Detection Result

---

### 3. preprocess
- YOLOv8 input 格式處理
- resize + padding（letterbox）
- 輸出 tensor 與 scale / padding 資訊

---

### 4. postprocess
- YOLOv8 output decode
- bounding box 還原
- NMS（Non-Maximum Suppression）

---

### 5. logic（規則過濾）
- IoU 計算
- Polygon 區域判斷
- Class size rule filtering
- （目前未整合進 pipeline）

---

### 6. camera（影像來源）
- 使用 GoCV 進行影像擷取
- 支援 webcam / RTSP
- 持續更新最新 frame
- （目前未接入推論流程）

---

### 7. store（狀態管理）
- Pipeline registry（模型管理）
- Camera ↔ Model 綁定
- 推論結果 cache
- （目前未在 runtime 使用）

---

### 8. server（HTTP API）
- 提供 `/infer` 推論 API
- 目前僅支援 JPEG 單張圖片
- 同步推論模式（request → inference → response）

---

## 四、資料流（目前版本）

### 單張圖片推論流程
```
HTTP Request
↓
JPEG Decode
↓
Pipeline.Infer()
↓
Preprocess (Letterbox)
↓
ONNX Runtime (YOLOv8n.onnx)
↓
Postprocess (Decode + NMS)
↓
Detection Result
↓
JSON Response
```


---

## 五、目前完成狀態

### ✔ 已完成（可運行功能）

- YOLOv8 ONNX 模型推論（CPU）
- 完整 inference pipeline（pre + inference + post）
- HTTP `/infer` API
- Detection JSON 輸出
- NMS + Bounding Box 修正
- CGO ONNX Runtime binding

---

### ⚠ 部分完成（未接入 runtime）

- Camera Manager（影像擷取已完成，但未接 pipeline）
- Filter logic（Polygon / Class rule 已完成但未使用）
- Store（Pipeline registry 已完成但未串接）

---

### ❌ 未完成（第一階段目標）

- Camera → Pipeline streaming 串流
- inference worker pool（併發推論）
- store-based result serving（非同步 API）
- queue / buffer architecture
- RTSP / MJPEG streaming output
- pipeline 動態控制 API

---

## 六、系統特性

### 已具備能力

- 單張圖片物件偵測
- YOLOv8 推論
- 模組化 pipeline 架構
- 可擴展 ONNX 模型替換

---

### 未來可擴展方向

- 即時影像串流推論（RTSP / webcam）
- 多 camera pipeline
- GPU 加速（CUDA / TensorRT）
- 事件觸發式分析（alert system）
- WebSocket 即時回傳結果

---

## 七、開發日誌（Development Log）

### 2026-07-03

#### 完成項目

- 完成 YOLOv8 ONNX 推論核心（CGO + ONNX Runtime）
- 建立完整 inference pipeline：
  - preprocess
  - inference
  - postprocess

- 實作 HTTP API `/infer` 支援單張 JPEG 推論
- 完成 YOLOv8 detection decode 與 NMS（Non-Maximum Suppression）
- 建立 camera manager（frame capture loop）
- 建立 store 架構：
  - pipeline registry
  - inference result cache

- 實作 polygon / class rule filtering logic（尚未接入 pipeline）

#### 系統狀態

目前：

> 單張圖片推論可用，但尚未進入 streaming pipeline 架構


---

### 2026-07-13

#### 完成項目

- 新增 `internal/app` 層：
  - Runtime 初始化
  - Model 載入
  - Camera 初始化
  - Pipeline 註冊

- 修改 config 架構：
  - 支援 camera 設定
  - 支援 model binding
  - 支援 stream interval

- 修正 ONNX Runtime 啟動問題：
  - 啟用 onnxruntime build tag
  - 更新 ONNX Runtime 1.26 DLL
  - 排除舊版本 DLL 衝突

- 完成 streaming inference 流程：
```
Camera
↓
Stream Worker
↓
YOLO Pipeline
↓
ONNX Runtime
↓
Store
↓
API Response
```

- 驗證：
  - `/health`
  - `/api/camera/check`
  - `/api/infer_od/live_result`

#### 測試結果

確認：

- Camera frame capture 正常
- YOLOv8 streaming inference 正常
- Result API 可取得 detection

#### 尚未完成

- Camera 與 inference lifecycle 尚未完全分離
- API handler 尚未模組化
- 尚未完成 integration test


#### 下一步

1. 分離 camera open 與 inference start
2. 重構 server handler 架構
3. 完成 camera / inference API integration test


---

### 2026-07-16

#### 完成項目

- 移除 App 啟動時自動開啟 camera / stream
- Camera lifecycle 改由 API 控制：
```
camera/open
camera/close
camera/check
camera/frame
camera/live
```
- 完成 Server handler 模組化：
```
internal/server

server.go
routes.go
health_handler.go
camera_handler.go
inference_handler.go
od_handler.go
```
- 分離 Camera 與 Inference lifecycle：
```
Camera

↓

Inference Stream

↓

Result Store

↓

API
```

- 完成 OD API：
  - inference start / stop
  - live result query
  - inference stream

- 補充 Class Name Mapping：
  - ClassID → ClassName

- 新增 server integration test：
  - inference start
  - result storage
  - API response

- 通過 race test：
```
go test -race ./...
```

#### 測試結果

確認流程：
```
Camera API
↓
Stream Manager
↓
YOLO Pipeline
↓
PipelineStore
↓
Inference API
```

可正常運作。

#### 尚未完成

- inference API 仍以 OD 為主要設計
- Pipeline 架構尚未抽象化
- 尚未支援 CLS / SEG / Pose 等任務


#### 下一步

1. 泛化 inference task 架構，支援 OD / CLS / SEG
2. 建立 task / pipeline registry
3. 抽象化 inference result structure
4. 完成 Docker 化部署流程

### 2026-07-17

#### 完成項目

- 完成 Runtime Pipeline Binding 架構：
  - Runtime 負責管理 Pipeline 註冊與 Camera 綁定
  - PipelineStore 增加：
    - modelName → Pipeline registry
    - cameraName → modelName binding
    - cameraName → latest inference result cache

架構：
```
Runtime

↓

PipelineStore

├── Pipeline Registry
│ model → pipeline
│
├── Camera Binding
│ camera → model
│
└── Result Cache
camera → inference result
```

- 完成 Stream Manager 與 Pipeline Binding 整合：
```
Camera

↓

StreamManager

↓

Get camera binding

↓

Get Pipeline

↓

YOLO Inference

↓

Store Result

↓

API Query
```

- 修正 Camera Lifecycle 資源釋放問題：
  - 增加 capture worker 管理
  - 使用 per-camera WaitGroup 等待 capture goroutine 結束
  - 修正 VideoCapture 關閉時 device 未釋放問題
  - 改善 gocv Mat 資源管理

Camera lifecycle：
```
camera/open

↓

VideoCapture open

↓

captureLoop start

↓

camera/close

↓

stop captureLoop

↓

release VideoCapture

↓

release Mat
```


- 修正 Camera API 測試流程：
  - 支援 video file 作為 camera source
  - 使用 test.mp4 模擬 camera input
  - 驗證完整 inference flow

測試流程：
```
test.mp4

↓

CameraManager

↓

StreamManager

↓

YOLOv8 Pipeline

↓

Result Store

↓

Inference API
```


- 新增 API integration test：

測試內容：

- camera/open
- camera/check
- camera/frame
- inference/start
- inference result generation
- inference/live_result
- inference/stop
- camera/close

- 完成實際 Build 測試：

#### 測試結果

確認完整 Runtime Flow：
```
Camera API

↓

CameraManager

↓

StreamManager

↓

Pipeline Binding

↓

YOLOv8 Pipeline

↓

PipelineStore

↓

Inference API
```

可正常運作。


實際 API 測試：
```
camera/open
✓

inference/start
✓

infer_od/live_result
✓

inference/stop
✓

camera/close
✓
```


確認：

- Camera 與 Inference lifecycle 可獨立控制
- Video source 可取代實體 camera 進行測試
- Pipeline binding 可正常將 camera 導向指定模型
- Result cache 可提供即時 inference 結果


#### 尚未完成

- Pipeline interface 尚未抽象化
  - 目前 StreamManager 仍直接依賴：
    ```
    *pipeline.Pipeline
    ```
  - 尚無法直接替換不同 inference task


- Inference Result 結構仍偏向 Object Detection：
  - Detection
  - Bounding Box
  - Class ID


尚未支援：
```
Classification

Segmentation

Pose Estimation

Tracking
```
- Live Streaming API 尚未完整驗證：
  - camera/live 僅完成 route
  - 尚未確認長時間串流穩定性
  - 尚未整合前端或 WebSocket/MJPEG output


- Deployment 尚未開始：
  - Docker image
  - runtime config
  - model volume management

#### 下一步

1. Pipeline Interface 化
   - 抽象 inference pipeline
   - 支援 OD / CLS / SEG 等不同任務


2. 建立 Pipeline Registry
   - 統一管理不同模型 pipeline


3. 泛化 Inference Result
   - 支援不同任務輸出格式


4. 完成 Pipeline 抽象後，再進行：
   - camera/live 串流輸出
   - 前端即時畫面整合
