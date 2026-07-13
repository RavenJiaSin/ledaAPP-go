package stream

import (
	"testing"
	"time"

	"gocv.io/x/gocv"
	"yolo-go-inference/internal/store"
	"yolo-go-inference/internal/pipeline"
	"yolo-go-inference/internal/postprocess"
	"yolo-go-inference/internal/camera"

)

// fakeRunner is a mock implementation of pipeline.Runner
type fakeRunner struct {
	output []float32
	err    error
}

func (f *fakeRunner) Run(input []float32, shape []int64) ([]float32, error) {
	return f.output, f.err
}

func TestStreamManager_LifecycleAndFiltering(t *testing.T) {

	cameraMgr := camera.NewTestCamera()

	pipelineStore := store.NewPipelineStore()


	sm := NewStreamManager(
		cameraMgr,
		pipelineStore,
	)

	defer sm.StopAll()



	// ============================
	// fake YOLO pipeline
	// ============================

	numPreds := 10
	numClasses := 2

	output := make(
		[]float32,
		(4+numClasses)*numPreds,
	)


	// bbox
	output[0*numPreds] = 5
	output[1*numPreds] = 5
	output[2*numPreds] = 2
	output[3*numPreds] = 2


	// class confidence
	output[4*numPreds] = 0.95



	runner := &fakeRunner{
		output: output,
	}


	p := pipeline.NewPipeline(
		runner,
		10,
		numPreds,
		numClasses,
		postprocess.YOLOv8LayoutChannelsFirst,
		0.5,
		0.45,
	)



	// ============================
	// register pipeline
	// ============================

	pipelineStore.RegisterPipeline(
		"yolov8n",
		p,
	)


	pipelineStore.BindCamera(
		"cam_test",
		"yolov8n",
	)



	// ============================
	// inject fake frame
	// ============================

	mat := gocv.NewMatWithSize(
		10,
		10,
		gocv.MatTypeCV8UC3,
	)

	cameraMgr.InjectFrame(
		"cam_test",
		mat,
	)

	mat.Close()



	// ============================
	// start stream
	// ============================

	err := sm.StartStream(
		"cam_test",
		StreamConfig{
			Interval:
				10 * time.Millisecond,
		},
	)


	if err != nil {
		t.Fatal(err)
	}



	// ============================
	// wait inference
	// ============================

	time.Sleep(
		200*time.Millisecond,
	)



	result, ok :=
		pipelineStore.GetCameraResult(
			"cam_test",
		)


	if !ok {
		t.Fatal(
			"no inference result",
		)
	}


	if result.Result.Count != 1 {

		t.Fatalf(
			"expect detection got %d",
			result.Result.Count,
		)
	}



	sm.StopStream(
		"cam_test",
	)
}
