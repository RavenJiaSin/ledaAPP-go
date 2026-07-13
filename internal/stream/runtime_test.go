package stream

import (
	"testing"
	"time"
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


	if r.Store == nil {
		t.Fatal("store missing")
	}


	// 目前只測 API 存在
	// 真正 pipeline 測試在 integration test

	r.BindCamera(
		"cam0",
		"yolov8",
	)


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
