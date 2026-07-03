package server

import (
	"encoding/json"
	"image"
	"image/jpeg"
	"net/http"

	"yolo-go-inference/pkg/types"
)

type Inferencer interface {
	Infer(img image.Image) (types.InferenceResult, error)
}

type Server struct {
	Inferencer Inferencer
}

func New(inferencer Inferencer) *Server {
	return &Server{
		Inferencer: inferencer,
	}
}

func (s *Server) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *Server) InferHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "failed to read image", http.StatusBadRequest)
		return
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		http.Error(w, "invalid image format (jpeg only)", http.StatusBadRequest)
		return
	}

	result, err := s.Inferencer.Infer(img)
	if err != nil {
		http.Error(w, "inference failed", http.StatusInternalServerError)
		return
	}

	resp := toResponse(result)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type DetectionResp struct {
	X1         float32 `json:"x1"`
	Y1         float32 `json:"y1"`
	X2         float32 `json:"x2"`
	Y2         float32 `json:"y2"`
	Confidence float32 `json:"confidence"`
	ClassID    int     `json:"class_id"`
	ClassName  string  `json:"class_name"`
}

func toResponse(result types.InferenceResult) []DetectionResp {
	out := make([]DetectionResp, 0, len(result.Detections))

	for _, d := range result.Detections {
		out = append(out, DetectionResp{
			X1:         d.X1,
			Y1:         d.Y1,
			X2:         d.X2,
			Y2:         d.Y2,
			Confidence: d.Confidence,
			ClassID:    d.ClassID,
			ClassName:  d.ClassName,
		})
	}

	return out
}
