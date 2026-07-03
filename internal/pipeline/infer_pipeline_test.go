package pipeline

import (
	"image"
	"image/color"
	"testing"
	"yolo-go-inference/internal/postprocess"
)

type fakeRunner struct {
	output []float32
	err    error

	gotInputLen int
	gotShape    []int64
}

func (f *fakeRunner) Run(input []float32, shape []int64) ([]float32, error) {
	f.gotInputLen = len(input)
	f.gotShape = append([]int64(nil), shape...)
	return f.output, f.err
}

func TestPipelineInfer(t *testing.T) {
	numPreds := 8400
	numClasses := 3

	output := make([]float32, (4+numClasses)*numPreds)

	// Detection at model input space:
	// cx=100, cy=120, w=40, h=20, class 2 score=0.95
	output[0*numPreds+0] = 100
	output[1*numPreds+0] = 120
	output[2*numPreds+0] = 40
	output[3*numPreds+0] = 20
	output[(4+2)*numPreds+0] = 0.95

	runner := &fakeRunner{output: output}

	p := NewPipeline(
		runner,
		640,
		numPreds,
		numClasses,
		postprocess.YOLOv8LayoutChannelsFirst,
		0.5,
		0.45,
	)

	img := image.NewRGBA(image.Rect(0, 0, 640, 640))
	fillImage(img, color.RGBA{255, 255, 255, 255})

	got, err := p.Infer(img)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if runner.gotInputLen != 3*640*640 {
		t.Fatalf("expected input len %d, got %d", 3*640*640, runner.gotInputLen)
	}

	expectedShape := []int64{1, 3, 640, 640}
	for i := range expectedShape {
		if runner.gotShape[i] != expectedShape[i] {
			t.Fatalf("expected shape %v, got %v", expectedShape, runner.gotShape)
		}
	}

	if got.Count != 1 {
		t.Fatalf("expected 1 detection, got %d", got.Count)
	}

	det := got.Detections[0]

	if det.ClassID != 2 {
		t.Fatalf("expected class 2, got %d", det.ClassID)
	}

	if det.Confidence != 0.95 {
		t.Fatalf("expected confidence 0.95, got %f", det.Confidence)
	}

	if det.X1 != 80 || det.Y1 != 110 || det.X2 != 120 || det.Y2 != 130 {
		t.Fatalf("unexpected box: %+v", det)
	}
}

func fillImage(img *image.RGBA, c color.RGBA) {
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			img.Set(x, y, c)
		}
	}
}

func TestPipelineInferClampsBoxesToImageBounds(t *testing.T) {
	numPreds := 8400
	numClasses := 3

	output := make([]float32, (4+numClasses)*numPreds)

	// cx=320, cy=320, w=800, h=900 => box extends outside 640x640
	output[0*numPreds+0] = 320
	output[1*numPreds+0] = 320
	output[2*numPreds+0] = 800
	output[3*numPreds+0] = 900
	output[(4+1)*numPreds+0] = 0.9

	runner := &fakeRunner{output: output}

	p := NewPipeline(
		runner,
		640,
		numPreds,
		numClasses,
		postprocess.YOLOv8LayoutChannelsFirst,
		0.5,
		0.45,
	)

	img := image.NewRGBA(image.Rect(0, 0, 640, 640))

	got, err := p.Infer(img)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got.Count != 1 {
		t.Fatalf("expected 1 detection, got %d", got.Count)
	}

	det := got.Detections[0]
	if det.X1 != 0 || det.Y1 != 0 || det.X2 != 640 || det.Y2 != 640 {
		t.Fatalf("expected clamped box to image bounds, got %+v", det)
	}
}
