package stream

import (
	"context"
	"errors"
	"sync"
	"time"

	"gocv.io/x/gocv"
	
	"yolo-go-inference/internal/logic"
	"yolo-go-inference/internal/pipeline"
	"yolo-go-inference/internal/store"
	"yolo-go-inference/pkg/types"
	"yolo-go-inference/internal/camera"
)

// StreamConfig defines rules and parameters for a camera stream
type StreamConfig struct {
	Polygons   []types.Polygon   // Alert regions for polygon filtering
	ClassRules []types.ClassRule // Size/class constraints for filtering
	Interval   time.Duration     // Frame sampling rate (e.g., 33ms for 30 FPS)
}

// StreamManager coordinates real-time inference pipelines for multiple cameras
type StreamManager struct {
	mu sync.RWMutex
	cameraMgr camera.CameraController
	store *store.PipelineStore
	cancels map[string]context.CancelFunc
	configs map[string]StreamConfig
}

// NewStreamManager creates a new StreamManager instance
func NewStreamManager(
	cameraMgr camera.CameraController,
	store *store.PipelineStore,
) *StreamManager {

	return &StreamManager{

		cameraMgr: cameraMgr,

		store: store,

		cancels: make(
			map[string]context.CancelFunc,
		),

		configs: make(
			map[string]StreamConfig,
		),
	}
}

// StartStream begins background dispatching and inference worker for a camera
func (sm *StreamManager) StartStream(cameraName string, config StreamConfig) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, running := sm.cancels[cameraName]; running {
		return errors.New("stream already running for camera: " + cameraName)
	}

	modelName, bound := sm.store.GetModelByCamera(cameraName)
	if !bound {
		return errors.New("no model bound to camera: " + cameraName)
	}

	pipe, exists := sm.store.GetPipeline(modelName)
	if !exists {
		return errors.New("pipeline not found for model: " + modelName)
	}

	ctx, cancel := context.WithCancel(context.Background())
	sm.cancels[cameraName] = cancel
	sm.configs[cameraName] = config

	// Size of 1 allows exactly 1 frame buffering with immediate dropping on backpressure
	frameCh := make(chan gocv.Mat, 1)

	// Start Goroutines for frame dispatching and pipeline inference
	go sm.dispatchLoop(ctx, cameraName, config.Interval, frameCh)
	go sm.workerLoop(ctx, cameraName, pipe, frameCh, config)

	return nil
}

// StopStream stops background processing for a camera
func (sm *StreamManager) StopStream(cameraName string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if cancel, running := sm.cancels[cameraName]; running {
		cancel()
		delete(sm.cancels, cameraName)
		delete(sm.configs, cameraName)
	}
}

// StopAll stops all active streaming pipelines
func (sm *StreamManager) StopAll() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	for name, cancel := range sm.cancels {
		cancel()
		delete(sm.cancels, name)
		delete(sm.configs, name)
	}
}

// IsRunning checks if a camera stream is active
func (sm *StreamManager) IsRunning(cameraName string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	_, running := sm.cancels[cameraName]
	return running
}

// dispatcherLoop handles sampling frames at the configured FPS/interval
func (sm *StreamManager) dispatchLoop(ctx context.Context, cameraName string, interval time.Duration, frameCh chan<- gocv.Mat) {
	if interval <= 0 {
		interval = 33 * time.Millisecond // Defaults to ~30 FPS
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	defer close(frameCh)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Fetch a cloned copy of the latest frame from the CameraManager
			mat, ok := sm.cameraMgr.LastFrame(cameraName)
			if !ok {
				continue
			}
			if mat.Empty() {
				mat.Close()
				continue
			}

			// Non-blocking channel push for backpressure/frame dropping
			select {
			case frameCh <- mat:
				// Worker accepted the frame; it now owns the mat and will call Close()
			default:
				// Worker is busy; immediately drop the frame and release CGO memory
				mat.Close()
			}
		}
	}
}

// workerLoop executes preprocessing, ONNX inference, postprocessing, filtering, and caching
func (sm *StreamManager) workerLoop(ctx context.Context, cameraName string, pipe *pipeline.Pipeline, frameCh <-chan gocv.Mat, config StreamConfig) {
	for {
		select {
		case <-ctx.Done():
			// Safely drain and clean up any remaining mats in the channel
			for mat := range frameCh {
				mat.Close()
			}
			return
		case mat, ok := <-frameCh:
			if !ok {
				return
			}

			// Convert GoCV Mat to native image.Image (CPU bound)
			img, err := mat.ToImage()
			mat.Close() // Safe release of the C++ resource immediately after conversion
			if err != nil {
				continue
			}

			// Run YOLO inference pipeline
			result, err := pipe.Infer(img)
			if err != nil {
				continue
			}

			// Apply polygon and class rules filter if specified
			if len(config.Polygons) > 0 || len(config.ClassRules) > 0 {
				filteredDets := logic.FilterDetections(result.Detections, config.Polygons, config.ClassRules)
				result.Detections = filteredDets
				result.Count = len(filteredDets)
			}

			// Store latest result under the camera cache
			sm.store.SetCameraResult(cameraName, result)
		}
	}
}
