package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"encoding/json"

	"github.com/gin-gonic/gin"

	"yolo-go-inference/internal/stream"
	"yolo-go-inference/internal/store"
	"yolo-go-inference/pkg/types"
)


// -----------------------------
// Fake StreamMgr
// -----------------------------


type fakeStreamMgr struct {
	store *store.PipelineStore
	running map[string]bool
}


func (f *fakeStreamMgr) StartStream(
	name string,
	cfg stream.StreamConfig,
) error {

	go func() {

		time.Sleep(
			50*time.Millisecond,
		)

		f.running[name] = true

		f.store.SetCameraResult(
			name,
			types.InferenceResult{
				Detections: []types.Detection{
					{
						ClassName:"person",
						Confidence:0.8,
					},
				},
			},
		)

	}()

	return nil
}



func (f *fakeStreamMgr) StopStream(name string){
	delete(f.running, name)
}


func (f *fakeStreamMgr) StopAll(){
}



// -----------------------------
// Test4
// Camera
//  |
// Stream
//  |
// Store
//  |
// API
// -----------------------------

func TestFullFlow(t *testing.T) {

	gin.SetMode(gin.TestMode)


	// -----------------
	// store
	// -----------------

	pipelineStore := store.NewPipelineStore()



	// -----------------
	// fake stream
	// -----------------

	streamMgr := &fakeStreamMgr{
		store: pipelineStore,
		running: make(map[string]bool),
	}



	// -----------------
	// server
	// -----------------

	runtime := stream.NewRuntime()

	runtime.Store = pipelineStore
	runtime.Stream = streamMgr


	srv := &Server{
		Runtime: runtime,
	}


	srv.Router = gin.New()


	srv.Router.POST(
		"/api/start_live_infer_od",
		srv.StartLiveInferOD,
	)


	srv.Router.GET(
		"/api/infer_od/live_result",
		srv.InferODLiveResult,
	)



	// -----------------
	// 1. start stream
	// -----------------

	body := []byte(`
	{
		"cam_name":"cam0",
		"interval":33
	}
	`)


	req := httptest.NewRequest(
		http.MethodPost,
		"/api/start_live_infer_od",
		bytes.NewBuffer(body),
	)


	req.Header.Set(
		"Content-Type",
		"application/json",
	)


	rec := httptest.NewRecorder()


	srv.Router.ServeHTTP(
		rec,
		req,
	)



	if rec.Code != http.StatusOK {

		t.Fatalf(
			"start stream failed: code=%d body=%s",
			rec.Code,
			rec.Body.String(),
		)
	}



	// -----------------
	// 2. wait inference
	// -----------------

	waitResult(
		t,
		pipelineStore,
		"cam0",
	)



	// -----------------
	// 3. query API
	// -----------------

	req = httptest.NewRequest(
		http.MethodGet,
		"/api/infer_od/live_result?name=cam0",
		nil,
	)


	rec = httptest.NewRecorder()


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

	var resp APIResponse

	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatalf("unexpected result format")
	}

	rawLast, ok := result["last_result"]
	if !ok {
		t.Fatalf("missing last_result field")
	}

	last, ok := rawLast.([]interface{})
	if !ok {
		t.Fatalf("last_result is not array")
	}

	if len(last) == 0 {
		t.Fatal("no detections")
	}

	// -----------------
	// 4. query non-exist camera
	// -----------------

	req = httptest.NewRequest(
		http.MethodGet,
		"/api/infer_od/live_result?name=not_exist",
		nil,
	)

	rec = httptest.NewRecorder()

	srv.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}

	// optional: 驗證 response 是空結果
	var resp2 APIResponse

	if err := json.Unmarshal(rec.Body.Bytes(), &resp2); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

}


// -----------------------------
// wait store result
// -----------------------------

func waitResult(
	t *testing.T,
	store *store.PipelineStore,
	camera string,
){

	timeout :=
		time.After(
			2*time.Second,
		)


	for {

		select {


		case <-timeout:

			t.Fatal(
				"timeout waiting inference result",
			)


		default:

			_, ok :=
				store.GetCameraResult(camera)


			if ok {
				return
			}


			time.Sleep(
				20*time.Millisecond,
			)

		}
	}

}



// -----------------------------
// contains
// -----------------------------

func contains(
	s string,
	sub string,
) bool {

	for i:=0; i+len(sub)<=len(s); i++ {

		if s[i:i+len(sub)] == sub {
			return true
		}

	}

	return false
}
