package server

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"yolo-go-inference/pkg/types"
)

type fakeInferencer struct {
	result types.InferenceResult
	err    error
}

func (f fakeInferencer) Infer(img image.Image) (types.InferenceResult, error) {
	return f.result, f.err
}

func TestHealthHandler(t *testing.T) {
	srv := New(fakeInferencer{})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	srv.HealthHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if rec.Body.String() != "OK" {
		t.Fatalf("expected body OK, got %q", rec.Body.String())
	}
}

func TestInferHandler(t *testing.T) {
	srv := New(fakeInferencer{
		result: types.InferenceResult{
			Detections: []types.Detection{
				{
					X1:         10,
					Y1:         20,
					X2:         30,
					Y2:         40,
					Confidence: 0.9,
					ClassID:    1,
					ClassName:  "person",
				},
			},
			Count: 1,
		},
	})

	body, contentType := makeJPEGMultipartBody(t)

	req := httptest.NewRequest(http.MethodPost, "/infer", body)
	req.Header.Set("Content-Type", contentType)

	rec := httptest.NewRecorder()

	srv.InferHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", rec.Code, rec.Body.String())
	}

	expected := `[{"x1":10,"y1":20,"x2":30,"y2":40,"confidence":0.9,"class_id":1,"class_name":"person"}]` + "\n"
	if rec.Body.String() != expected {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestInferHandlerRejectsNonPost(t *testing.T) {
	srv := New(fakeInferencer{})

	req := httptest.NewRequest(http.MethodGet, "/infer", nil)
	rec := httptest.NewRecorder()

	srv.InferHandler(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rec.Code)
	}
}

func TestInferHandlerRejectsMissingImage(t *testing.T) {
	srv := New(fakeInferencer{})

	req := httptest.NewRequest(http.MethodPost, "/infer", nil)
	rec := httptest.NewRecorder()

	srv.InferHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestInferHandlerRejectsInvalidJPEG(t *testing.T) {
	srv := New(fakeInferencer{})

	body, contentType := makeTextMultipartBody(t)

	req := httptest.NewRequest(http.MethodPost, "/infer", body)
	req.Header.Set("Content-Type", contentType)

	rec := httptest.NewRecorder()

	srv.InferHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestInferHandlerReturnsInternalServerError(t *testing.T) {
	srv := New(fakeInferencer{
		err: errors.New("inference failed"),
	})

	body, contentType := makeJPEGMultipartBody(t)

	req := httptest.NewRequest(http.MethodPost, "/infer", body)
	req.Header.Set("Content-Type", contentType)

	rec := httptest.NewRecorder()

	srv.InferHandler(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}
}

func makeJPEGMultipartBody(t *testing.T) (*bytes.Buffer, string) {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("image", "test.jpg")
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}

	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	fillImage(img, color.RGBA{255, 255, 255, 255})

	if err := jpeg.Encode(part, img, nil); err != nil {
		t.Fatalf("failed to encode jpeg: %v", err)
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close multipart writer: %v", err)
	}

	return body, writer.FormDataContentType()
}

func makeTextMultipartBody(t *testing.T) (*bytes.Buffer, string) {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("image", "test.txt")
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}

	if _, err := part.Write([]byte("not a jpeg")); err != nil {
		t.Fatalf("failed to write invalid image: %v", err)
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close multipart writer: %v", err)
	}

	return body, writer.FormDataContentType()
}

func fillImage(img *image.RGBA, c color.RGBA) {
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			img.Set(x, y, c)
		}
	}
}
