package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"yolo-go-inference/internal/stream"
	"yolo-go-inference/pkg/types"
	"yolo-go-inference/internal/store"
	"yolo-go-inference/internal/pipeline"
)



// --------------------------------------------------
// Fake Stream Controller
// --------------------------------------------------

type integrationFakeStream struct {

	store interface {
		SetCameraResult(
			string,
			types.InferenceResult,
		)
	}

	running map[string]bool
}



func (f *integrationFakeStream) StartStream(
	name string,
	cfg stream.StreamConfig,
) error {


	f.running[name] = true


	go func(){

		time.Sleep(
			50*time.Millisecond,
		)


		f.store.SetCameraResult(
			name,
			types.InferenceResult{

				Detections:
					[]types.Detection{

						{
							ClassID:0,

							ClassName:"person",

							Confidence:0.92,

							X1:100,

							Y1:100,

							X2:300,

							Y2:400,
						},

					},

			},
		)


	}()


	return nil
}



func (f *integrationFakeStream) StopStream(
	name string,
){

	delete(
		f.running,
		name,
	)

}



func (f *integrationFakeStream) StopAll(){

}



// --------------------------------------------------
// Full inference flow
//
// API
//  |
// Stream
//  |
// Store
//  |
// Query Result
// --------------------------------------------------

func TestServerInferenceFlow(
	t *testing.T,
){

	gin.SetMode(
		gin.TestMode,
	)



	// -------------------------
	// create server
	// -------------------------

	srv :=
		setupTestServer()

	srv.Runtime.RegisterPipeline(
		"yolov8",
		&pipeline.Pipeline{},
	)

	// -------------------------
	// replace stream manager
	// -------------------------

	fake :=
		&integrationFakeStream{

			store:
				srv.Runtime.Store,

			running:
				make(map[string]bool),
		}



	srv.Runtime.Stream =
		fake



	// -------------------------
	// routes
	// -------------------------

	srv.Router.POST(
		"/api/inference/start",
		srv.StartInference,
	)


	srv.Router.GET(
		"/api/infer_od/live_result",
		srv.InferODLiveResult,
	)



	// -------------------------
	// 1. start inference
	// -------------------------

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



	var resp APIResponse

	json.Unmarshal(
		rec.Body.Bytes(),
		&resp,
	)


	if resp.Code != 0 {
		t.Fatalf(
			"start inference failed: %s",
			resp.Msg,
		)
	}



	// -------------------------
	// 2. wait inference result
	// -------------------------

	waitInferenceResult(
		t,
		srv.Runtime.Store,
		"cam0",
	)



	// -------------------------
	// 3. query result API
	// -------------------------

	req =
		httptest.NewRequest(
			http.MethodGet,
			"/api/infer_od/live_result?name=cam0",
			nil,
		)


	rec =
		httptest.NewRecorder()



	srv.Router.ServeHTTP(
		rec,
		req,
	)



	if rec.Code != http.StatusOK {

		t.Fatalf(
			"unexpected status: %d",
			rec.Code,
		)

	}



	var response APIResponse


	err :=
		json.Unmarshal(
			rec.Body.Bytes(),
			&response,
		)


	if err != nil {

		t.Fatalf(
			"invalid json: %v",
			err,
		)

	}



	result,ok :=
		response.Result.(map[string]interface{})


	if !ok {

		t.Fatalf(
			"invalid result format",
		)

	}



	last,ok :=
		result["last_result"].([]interface{})


	if !ok {

		t.Fatalf(
			"missing last_result",
		)

	}



	if len(last)==0 {

		t.Fatal(
			"empty detection result",
		)

	}

}



// --------------------------------------------------
// query non-existing camera
// --------------------------------------------------

func TestInferODLiveResultNotFound(
	t *testing.T,
){

	srv :=
		setupTestServer()



	srv.Router.GET(
		"/api/infer_od/live_result",
		srv.InferODLiveResult,
	)



	req :=
		httptest.NewRequest(
			http.MethodGet,
			"/api/infer_od/live_result?name=unknown",
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
			"expected 200 got %d",
			rec.Code,
		)

	}

}



// --------------------------------------------------
// helper
// --------------------------------------------------

func waitInferenceResult(
	t *testing.T,
	store *store.PipelineStore,
	camera string,
) {

	timeout := time.After(
		2 * time.Second,
	)


	for {

		select {

		case <-timeout:

			t.Fatal(
				"timeout waiting inference result",
			)


		default:

			_, ok :=
				store.GetCameraResult(
					camera,
				)


			if ok {
				return
			}


			time.Sleep(
				20 * time.Millisecond,
			)
		}
	}
}
