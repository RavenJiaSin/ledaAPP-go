package pipeline

import (
	"image"

	"yolo-go-inference/internal/postprocess"
	"yolo-go-inference/internal/preprocess"
	"yolo-go-inference/pkg/types"
)

type Runner interface {
	Run(input []float32, shape []int64) ([]float32, error)
}

type Pipeline struct {
	Session       Runner
	InputSize     int
	NumPreds int
	NumClasses    int
	OutputLayout postprocess.YOLOv8Layout
	ConfThreshold float32
	IouThreshold  float32
}

// 建立 pipeline
func NewPipeline(
	session Runner,
	inputSize int,
	numPreds int,
	numClasses int,
	outputLayout postprocess.YOLOv8Layout,
	confThreshold float32,
	iouThreshold float32,
) *Pipeline {
	return &Pipeline{
		Session:       session,
		InputSize:     inputSize,
		NumPreds:      numPreds,
		NumClasses:    numClasses,
		OutputLayout:  outputLayout,
		ConfThreshold: confThreshold,
		IouThreshold:  iouThreshold,
	}
}

func clampDetections(dets []types.Detection, width, height int) []types.Detection {
	for i := range dets {
		dets[i].X1 = clampFloat32(dets[i].X1, 0, float32(width))
		dets[i].Y1 = clampFloat32(dets[i].Y1, 0, float32(height))
		dets[i].X2 = clampFloat32(dets[i].X2, 0, float32(width))
		dets[i].Y2 = clampFloat32(dets[i].Y2, 0, float32(height))
	}
	return dets
}

func clampFloat32(v, minValue, maxValue float32) float32 {
	if v < minValue {
		return minValue
	}
	if v > maxValue {
		return maxValue
	}
	return v
}

// 單張影像推論
func (p *Pipeline) Infer(img image.Image) (types.InferenceResult, error) {
	pre := preprocess.Letterbox(img, p.InputSize)

	output, err := p.Session.Run(pre.Tensor, pre.Shape)
	if err != nil {
		return types.InferenceResult{}, err
	}

	dets := postprocess.DecodeYOLOv8(
		output,
		p.NumPreds,
		p.NumClasses,
		p.OutputLayout,
		p.ConfThreshold,
		pre.Scale,
		pre.PadX,
		pre.PadY,
	)

	dets = postprocess.NMS(dets, p.IouThreshold)
	dets = clampDetections(dets, img.Bounds().Dx(), img.Bounds().Dy())

	return types.InferenceResult{
		Detections: dets,
		Count:      len(dets),
	}, nil
}
