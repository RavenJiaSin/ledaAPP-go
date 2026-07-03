package postprocess

import (
	"testing"

	"yolo-go-inference/pkg/types"
)

func TestNMSRemovesOverlappingBoxesSameClass(t *testing.T) {
	dets := []types.Detection{
		{X1: 0, Y1: 0, X2: 100, Y2: 100, Confidence: 0.9, ClassID: 0},
		{X1: 10, Y1: 10, X2: 110, Y2: 110, Confidence: 0.8, ClassID: 0},
		{X1: 200, Y1: 200, X2: 300, Y2: 300, Confidence: 0.7, ClassID: 0},
	}

	got := NMS(dets, 0.45)

	if len(got) != 2 {
		t.Fatalf("expected 2 detections, got %d", len(got))
	}

	if got[0].Confidence != 0.9 {
		t.Fatalf("expected highest confidence detection first, got %f", got[0].Confidence)
	}
}
