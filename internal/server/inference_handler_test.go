package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"yolo-go-inference/internal/pipeline"
)


func TestStartInference(t *testing.T){

	gin.SetMode(gin.TestMode)


	stream :=
		&fakeStream{}


	srv :=
		setupTestServer()

	srv.Runtime.RegisterPipeline(
		"yolov8",
		&pipeline.Pipeline{}, // 或 fake pipeline
	)


	srv.Runtime.Stream =
		stream




	srv.Router.POST(
		"/api/inference/start",
		srv.StartInference,
	)



	body :=
		[]byte(`
		{
			"cam_name":"cam0",
			"interval":33
		}
		`)


	req :=
		httptest.NewRequest(
			http.MethodPost,
			"/api/inference/start",
			bytes.NewBuffer(body),
		)


	req.Header.Set(
		"Content-Type",
		"application/json",
	)


	rec :=
		httptest.NewRecorder()



	srv.Router.ServeHTTP(
		rec,
		req,
	)



	if rec.Code != http.StatusOK {

		t.Fatalf(
			"unexpected status %d body=%s",
			rec.Code,
			rec.Body.String(),
		)

	}


	if !stream.started {

		t.Fatal(
			"stream not started",
		)
	}

}
