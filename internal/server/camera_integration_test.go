package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"yolo-go-inference/internal/camera"
	"yolo-go-inference/internal/stream"
)


func TestCameraAPIIntegration(t *testing.T) {


	gin.SetMode(gin.TestMode)


	runtime := stream.NewRuntime()

	defer runtime.Stop()


	runtime.Camera =
		camera.NewFakeCamera()



	srv := New(runtime)


	router := srv.Router



	// -----------------
	// 1. open camera
	// -----------------

	body := []byte(
		`{
			"name":"cam0",
			"uri":"0"
		}`,
	)


	req := httptest.NewRequest(
		http.MethodPost,
		"/api/camera/open",
		bytes.NewBuffer(body),
	)


	req.Header.Set(
		"Content-Type",
		"application/json",
	)


	rec := httptest.NewRecorder()


	router.ServeHTTP(
		rec,
		req,
	)



	if rec.Code != http.StatusOK {

		t.Fatalf(
			"open camera failed: %d %s",
			rec.Code,
			rec.Body.String(),
		)
	}



	// -----------------
	// 2. check camera
	// -----------------

	req = httptest.NewRequest(
		http.MethodGet,
		"/api/camera/check?name=cam0",
		nil,
	)


	rec = httptest.NewRecorder()


	router.ServeHTTP(
		rec,
		req,
	)



	if rec.Code != http.StatusOK {

		t.Fatalf(
			"check failed: %d",
			rec.Code,
		)
	}



	if rec.Body.String() != `{"alive":true}` {

		t.Fatalf(
			"unexpected check result: %s",
			rec.Body.String(),
		)
	}



	// -----------------
	// 3. get frame
	// -----------------

	req = httptest.NewRequest(
		http.MethodGet,
		"/api/camera/frame?name=cam0",
		nil,
	)


	rec = httptest.NewRecorder()


	router.ServeHTTP(
		rec,
		req,
	)



	if rec.Code != http.StatusOK {

		t.Fatalf(
			"frame failed: %d %s",
			rec.Code,
			rec.Body.String(),
		)
	}



	contentType :=
		rec.Header().Get(
			"Content-Type",
		)


	if contentType != "image/jpeg" {

		t.Fatalf(
			"unexpected content type: %s",
			contentType,
		)
	}



	if rec.Body.Len() == 0 {

		t.Fatal(
			"empty jpeg",
		)
	}

}
