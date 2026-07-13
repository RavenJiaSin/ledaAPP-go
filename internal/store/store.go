package store

import (
	"sync"
	"time"

	"yolo-go-inference/internal/pipeline"
	"yolo-go-inference/pkg/types"
)

// -----------------------------
// Detection Result (NEW)
// -----------------------------

type DetectionResult struct {
	CameraName string
	Timestamp  time.Time
	Result     types.InferenceResult
}

// -----------------------------
// Pipeline Registry
// -----------------------------

type PipelineStore struct {
	mu sync.RWMutex

	// modelName → pipeline
	pipelines map[string]*pipeline.Pipeline

	// cameraName → modelName
	cameraBind map[string]string

	// cameraName → last result  ❗改這裡（原本是 modelName）
	lastResult map[string]DetectionResult
}

// 建立 store
func NewPipelineStore() *PipelineStore {
	return &PipelineStore{
		pipelines:  make(map[string]*pipeline.Pipeline),
		cameraBind: make(map[string]string),
		lastResult: make(map[string]DetectionResult),
	}
}

//
// -----------------------------
// Pipeline 管理
// -----------------------------

// 註冊 pipeline
func (s *PipelineStore) RegisterPipeline(name string, p *pipeline.Pipeline) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.pipelines[name] = p
}

// 取得 pipeline
func (s *PipelineStore) GetPipeline(name string) (*pipeline.Pipeline, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p, ok := s.pipelines[name]
	return p, ok
}

//
// -----------------------------
// camera → model binding
// -----------------------------

func (s *PipelineStore) BindCamera(cameraName, modelName string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cameraBind[cameraName] = modelName
}

func (s *PipelineStore) GetModelByCamera(cameraName string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	m, ok := s.cameraBind[cameraName]
	return m, ok
}

//
// -----------------------------
// inference result cache (UPDATED)
// -----------------------------

// ❗ streaming worker 用
func (s *PipelineStore) SetCameraResult(cameraName string, result types.InferenceResult) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.lastResult[cameraName] = DetectionResult{
		CameraName: cameraName,
		Timestamp:  time.Now(),
		Result:     result,
	}
}

// ❗ API 用
func (s *PipelineStore) GetCameraResult(cameraName string) (DetectionResult, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	r, ok := s.lastResult[cameraName]
	return r, ok
}

//
// -----------------------------
// helper
// -----------------------------

func (s *PipelineStore) HasPipeline(name string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.pipelines[name]
	return ok
}
