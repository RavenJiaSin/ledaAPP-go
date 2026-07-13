package stream

import (
	"testing"
	"time"

	"gocv.io/x/gocv"

	"yolo-go-inference/internal/camera"
	"yolo-go-inference/internal/pipeline"
	"yolo-go-inference/internal/postprocess"
	"yolo-go-inference/internal/store"
)

// fake ONNX runtime
type integrationFakeRunner struct {
	output []float32
	err    error
}

func (f *integrationFakeRunner) Run(
	input []float32,
	shape []int64,
) ([]float32, error) {
	return f.output, f.err
}

func TestFullStreamingPipeline(t *testing.T) {

	cameraMgr := camera.NewTestCamera()

	pipelineStore := store.NewPipelineStore()

	streamMgr := NewStreamManager(
		cameraMgr,
		pipelineStore,
	)

	defer streamMgr.StopAll()
	defer cameraMgr.CloseAll()

	numPreds := 10
	numClasses := 2

	output := make(
		[]float32,
		(4+numClasses)*numPreds,
	)

	output[0*numPreds] = 5
	output[1*numPreds] = 5
	output[2*numPreds] = 2
	output[3*numPreds] = 2

	output[4*numPreds] = 0.95

	runner := &integrationFakeRunner{
		output: output,
	}

	pipe := pipeline.NewPipeline(
		runner,
		10,
		numPreds,
		numClasses,
		postprocess.YOLOv8LayoutChannelsFirst,
		0.5,
		0.45,
	)

	pipelineStore.RegisterPipeline(
		"test-model",
		pipe,
	)

	pipelineStore.BindCamera(
		"cam-test",
		"test-model",
	)

	mat := gocv.NewMatWithSize(
		10,
		10,
		gocv.MatTypeCV8UC3,
	)

	defer mat.Close()

	cameraMgr.InjectFrame(
		"cam-test",
		mat,
	)

	err := streamMgr.StartStream(
		"cam-test",
		StreamConfig{
			Interval: 10 * time.Millisecond,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	timeout := time.After(
		2 * time.Second,
	)

	for {

		select {

		case <-timeout:

			t.Fatal(
				"timeout",
			)

		default:

			result, ok :=
				pipelineStore.GetCameraResult(
					"cam-test",
				)

			if ok {

				if result.Result.Count != 1 {
					t.Fatalf(
						"expect 1 got %d",
						result.Result.Count,
					)
				}

				return
			}

			time.Sleep(
				10 * time.Millisecond,
			)
		}
	}
}
