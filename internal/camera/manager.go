package camera

import (
	"errors"
	"sync"

	"gocv.io/x/gocv"
)

// -----------------------------
// Camera Manager
// -----------------------------



type CameraManager struct {
	mu sync.RWMutex

	// cameraName → VideoCapture
	cameras map[string]*gocv.VideoCapture

	// cameraName → last frame
	lastFrames map[string]gocv.Mat

	// running state
	running map[string]bool
}

// 建立 manager
func NewCameraManager() *CameraManager {
	return &CameraManager{
		cameras:    make(map[string]*gocv.VideoCapture),
		lastFrames: make(map[string]gocv.Mat),
		running:    make(map[string]bool),
	}
}

//
// -----------------------------
// Open camera
// -----------------------------
// source: 0 / RTSP / video path
//

func (m *CameraManager) Open(cameraName string, source interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.cameras[cameraName]; exists {
		return errors.New("camera already opened")
	}

	var cap *gocv.VideoCapture
	var err error

	switch v := source.(type) {

	case int:
		cap, err = gocv.OpenVideoCapture(v)

	case string:
		cap, err = gocv.OpenVideoCapture(v)

	default:
		return errors.New("invalid camera source type")
	}

	if err != nil {
		return err
	}

	if !cap.IsOpened() {
		return errors.New("failed to open camera")
	}

	m.cameras[cameraName] = cap
	m.lastFrames[cameraName] = gocv.NewMat()
	m.running[cameraName] = true

	// start goroutine capture loop
	go m.captureLoop(cameraName, cap)

	return nil
}

//
// -----------------------------
// capture loop
// -----------------------------
// continuous frame update
//

func (m *CameraManager) captureLoop(name string, cap *gocv.VideoCapture) {
	for {

		m.mu.RLock()
		if !m.running[name] {
			m.mu.RUnlock()
			return
		}
		m.mu.RUnlock()

		mat := gocv.NewMat()
		if ok := cap.Read(&mat); !ok || mat.Empty() {
			mat.Close()
			continue
		}

		// ✅ clone 一份給 store（避免 race）
		clone := mat.Clone()
		mat.Close()

		m.mu.Lock()

		if old, ok := m.lastFrames[name]; ok {
			old.Close()
		}

		m.lastFrames[name] = clone

		m.mu.Unlock()
	}
}

//
// -----------------------------
// Get last frame
// -----------------------------
// non-blocking
//

func (m *CameraManager) LastFrame(cameraName string) (gocv.Mat, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	frame, ok := m.lastFrames[cameraName]
	if !ok || frame.Empty() {
		return gocv.Mat{}, false
	}

	return frame.Clone(), true
}

//
// -----------------------------
// Stop camera
// -----------------------------

func (m *CameraManager) Close(cameraName string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if cap, ok := m.cameras[cameraName]; ok {
		cap.Close()
		delete(m.cameras, cameraName)
	}

	if frame, ok := m.lastFrames[cameraName]; ok {
		frame.Close()
		delete(m.lastFrames, cameraName)
	}

	delete(m.running, cameraName)
}

//
// -----------------------------
// Stop all
// -----------------------------

func (m *CameraManager) CloseAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for name, cap := range m.cameras {
		cap.Close()
		delete(m.cameras, name)
	}

	for name, frame := range m.lastFrames {
		frame.Close()
		delete(m.lastFrames, name)
	}

	m.running = make(map[string]bool)
}

// InjectFrame injects a gocv.Mat into the camera manager for testing purposes.
func (m *CameraManager) InjectFrame(cameraName string, mat gocv.Mat) {
	m.mu.Lock()
	defer m.mu.Unlock()

	clone := mat.Clone()

	if old, ok := m.lastFrames[cameraName]; ok {
		old.Close()
	}

	m.lastFrames[cameraName] = clone
}
