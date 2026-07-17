package stream

import (
	"context"
	"sync"
	"fmt"

	"yolo-go-inference/internal/camera"
	"yolo-go-inference/internal/pipeline"
	"yolo-go-inference/internal/store"
	
)


type Runtime struct {

	Camera camera.CameraController

	Store *store.PipelineStore

	Stream StreamController



	ctx context.Context

	cancel context.CancelFunc

	once sync.Once
}



// NewRuntime 建立 Runtime
func NewRuntime() *Runtime {

	cameraMgr :=
		camera.NewCameraManager()


	pipelineStore :=
		store.NewPipelineStore()


	ctx, cancel :=
		context.WithCancel(
			context.Background(),
		)


	streamMgr :=
		NewStreamManager(
			cameraMgr,
			pipelineStore,
		)


	return &Runtime{

		Camera: cameraMgr,

		Store: pipelineStore,

		Stream: streamMgr,

		ctx: ctx,

		cancel: cancel,
	}
}



// Start 啟動 runtime
func (r *Runtime) Start() error {

	return nil
}



// Stop 關閉 runtime
func (r *Runtime) Stop() {

	r.once.Do(func(){

		r.Stream.StopAll()

		r.Camera.CloseAll()

		r.cancel()

	})
}



// Context
func (r *Runtime) Context() context.Context {

	return r.ctx
}



// -----------------------------
// Pipeline management
// -----------------------------


// RegisterPipeline
// 註冊模型 pipeline

func (r *Runtime) RegisterPipeline(
	name string,
	p *pipeline.Pipeline,
) {
	r.Store.RegisterPipeline(name, p)
}


func (r *Runtime) BindCamera(
	cameraName string,
	modelName string,
) error {

	_, ok := r.Store.GetPipeline(modelName)
	if !ok {
		return fmt.Errorf(
			"pipeline not found for model: %s",
			modelName,
		)
	}

	// ✅ 寫入 store（唯一來源）
	r.Store.BindCamera(cameraName, modelName)

	return nil
}
