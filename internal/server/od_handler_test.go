package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"yolo-go-inference/internal/pipeline"

)



func TestInferODLiveResultEmpty(
	t *testing.T,
){

	srv :=
		setupTestServer()

	srv.Runtime.RegisterPipeline(
		"yolov8",
		&pipeline.Pipeline{}, // 或 fake pipeline
	)

	srv.Router.GET(
		"/api/infer_od/live_result",
		srv.InferODLiveResult,
	)

	req :=
		httptest.NewRequest(
			http.MethodGet,
			"/api/infer_od/live_result?name=cam0",
			nil,
		)


	rec :=
		httptest.NewRecorder()



	srv.Router.ServeHTTP(
		rec,
		req,
	)



	if rec.Code != http.StatusOK {

		t.Fatalf(
			"expected 200",
		)

	}

}
