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

- 完成 YOLOv8 ONNX 推論核心（CGO + ONNX Runtime）
- 建立完整 inference pipeline（preprocess → inference → postprocess）
- 實作 HTTP API `/infer` 支援單張 JPEG 推論
- 完成 NMS（Non-Maximum Suppression）與 detection decode
- 建立 camera manager（frame capture loop）
- 建立 store 架構（pipeline registry + result cache）
- 實作 polygon / class rule filtering logic（尚未接入 pipeline）
- 系統目前處於：
  > 單張圖片推論可用，但尚未進入 streaming pipeline 架構

---

## 八、系統總結

目前系統已完成：

> 「YOLOv8 單張圖片推論引擎」

下一階段目標為：

> 「從 request-based inference → 轉換為 streaming inference system」
