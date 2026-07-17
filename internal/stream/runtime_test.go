package stream

import (
	"testing"
	"time"

	"yolo-go-inference/internal/pipeline"
)


func TestRuntimeLifecycle(t *testing.T) {

	runtime := NewRuntime()


	if runtime.Camera == nil {
		t.Fatal("camera manager is nil")
	}


	if runtime.Store == nil {
		t.Fatal("store is nil")
	}


	if runtime.Stream == nil {
		t.Fatal("stream manager is nil")
	}



	err := runtime.Start()

	if err != nil {
		t.Fatalf(
			"runtime start failed: %v",
			err,
		)
	}



	ctx := runtime.Context()


	select {

	case <-ctx.Done():

		t.Fatal(
			"context canceled before stop",
		)


	default:

	}



	runtime.Stop()



	select {

	case <-ctx.Done():

	case <-time.After(
		time.Second,
	):

		t.Fatal(
			"context not cancelled",
		)
	}

}

func TestRuntimePipelineBinding(t *testing.T){

	r := NewRuntime()

	defer r.Stop()


	r.RegisterPipeline(
		"yolov8",
		&pipeline.Pipeline{},
	)


	err := r.BindCamera(
		"cam0",
		"yolov8",
	)


	if err != nil {
		t.Fatalf(
			"bind camera failed: %v",
			err,
		)
	}


	model, ok :=
		r.Store.GetModelByCamera(
			"cam0",
		)


	if !ok {
		t.Fatal(
			"camera binding missing",
		)
	}


	if model != "yolov8" {

		t.Fatalf(
			"unexpected model %s",
			model,
		)
	}

}
