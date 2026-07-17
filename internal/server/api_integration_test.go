package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"yolo-go-inference/internal/pipeline"
)


// --------------------------------------------------
// Fake Pipeline
// --------------------------------------------------

type fakeRunner struct{}


func (f *fakeRunner) Run(
	input []float32,
	shape []int64,
) ([]float32,error){

	return make([]float32,0),nil
}


// --------------------------------------------------
// helper
// --------------------------------------------------

func doRequest(
	router *gin.Engine,
	method string,
	path string,
	body []byte,
) *httptest.ResponseRecorder {


	req :=
		httptest.NewRequest(
			method,
			path,
			bytes.NewBuffer(body),
		)


	req.Header.Set(
		"Content-Type",
		"application/json",
	)


	rec :=
		httptest.NewRecorder()


	router.ServeHTTP(
		rec,
		req,
	)


	return rec
}



// --------------------------------------------------
// Full lifecycle test
// --------------------------------------------------

func TestCameraInferenceLifecycle(
	t *testing.T,
){


	gin.SetMode(
		gin.TestMode,
	)



	// -------------------------
	// check video
	// -------------------------

	if _,err :=
		os.Stat("../../test.mp4");
		err != nil {

		t.Skip(
			"test.mp4 not found",
		)
	}



	srv :=
		setupTestServer()



	// -------------------------
	// register fake pipeline
	// -------------------------

	fakePipeline :=
		&pipeline.Pipeline{}


	srv.Runtime.RegisterPipeline(
		"yolov8",
		fakePipeline,
	)



	// -------------------------
	// routes
	// -------------------------

	srv.Router.POST(
		"/api/camera/open",
		srv.OpenCamera,
	)


	srv.Router.POST(
		"/api/camera/close",
		srv.CloseCamera,
	)


	srv.Router.GET(
		"/api/camera/check",
		srv.CheckCamera,
	)


	srv.Router.GET(
		"/api/camera/frame",
		srv.CameraFrame,
	)


	srv.Router.POST(
		"/api/inference/start",
		srv.StartInference,
	)


	srv.Router.POST(
		"/api/inference/stop",
		srv.StopInference,
	)


	srv.Router.GET(
		"/api/infer_od/live_result",
		srv.InferODLiveResult,
	)



	// -------------------------
	// 1. open camera
	// -------------------------

	openBody :=
		[]byte(`
		{
			"cam_name":"cam0",
			"uri":"D:\\NCU\\intern\\workspace\\LedaAPP_go\\ledaAPP-go\\test.mp4"
		}
		`)


	rec :=
		doRequest(
			srv.Router,
			http.MethodPost,
			"/api/camera/open",
			openBody,
		)



	if rec.Code != http.StatusOK {

		t.Fatalf(
			"open camera failed: %s",
			rec.Body.String(),
		)

	}




	// -------------------------
	// 2. check camera
	// -------------------------

	rec =
		doRequest(
			srv.Router,
			http.MethodGet,
			"/api/camera/check?name=cam0",
			nil,
		)



	if rec.Code != http.StatusOK {

		t.Fatalf(
			"camera check failed",
		)

	}





	// -------------------------
	// 3. start inference
	// -------------------------

	startBody :=
		[]byte(`
		{
			"cam_name":"cam0",
			"model":"yolov8",
			"interval":33
		}
		`)



	rec =
		doRequest(
			srv.Router,
			http.MethodPost,
			"/api/inference/start",
			startBody,
		)



	if rec.Code != http.StatusOK {

		t.Fatalf(
			"start inference failed: %s",
			rec.Body.String(),
		)

	}



	// 等待 stream
	time.Sleep(
		2*time.Second,
	)



	// -------------------------
	// 4. query result
	// -------------------------

	rec =
		doRequest(
			srv.Router,
			http.MethodGet,
			"/api/infer_od/live_result?name=cam0",
			nil,
		)



	if rec.Code != http.StatusOK {

		t.Fatalf(
			"live result failed",
		)

	}



	var response APIResponse


	err :=
		json.Unmarshal(
			rec.Body.Bytes(),
			&response,
		)


	if err != nil {

		t.Fatal(err)

	}



	t.Log(
		"result:",
		response,
	)




	// -------------------------
	// 5. stop inference
	// -------------------------

	rec =
		doRequest(
			srv.Router,
			http.MethodPost,
			"/api/inference/stop",
			[]byte(`
			{
				"cam_name":"cam0"
			}
			`),
		)



	if rec.Code != http.StatusOK {

		t.Fatal(
			"stop inference failed",
		)

	}



	// -------------------------
	// 6. close camera
	// -------------------------

	rec =
		doRequest(
			srv.Router,
			http.MethodPost,
			"/api/camera/close",
			[]byte(`
			{
				"cam_name":"cam0"
			}
			`),
		)



	if rec.Code != http.StatusOK {

		t.Fatalf(
			"close camera failed: %s",
			rec.Body.String(),
		)

	}


}
