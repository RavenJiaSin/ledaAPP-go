package postprocess

import "testing"

func TestDecodeYOLOv8CHWLayout(t *testing.T) {
	numPreds := 2
	numClasses := 3

	output := make([]float32, (4+numClasses)*numPreds)

	// pred 0: cx=50, cy=60, w=20, h=10
	output[0*numPreds+0] = 50
	output[1*numPreds+0] = 60
	output[2*numPreds+0] = 20
	output[3*numPreds+0] = 10

	// pred 0 class scores: class 1 wins
	output[4*numPreds+0] = 0.1
	output[5*numPreds+0] = 0.9
	output[6*numPreds+0] = 0.2

	// pred 1 below threshold
	output[0*numPreds+1] = 10
	output[1*numPreds+1] = 10
	output[2*numPreds+1] = 5
	output[3*numPreds+1] = 5
	output[4*numPreds+1] = 0.1
	output[5*numPreds+1] = 0.2
	output[6*numPreds+1] = 0.3

	dets := DecodeYOLOv8(
		output,
		numPreds,
		numClasses,
		YOLOv8LayoutChannelsFirst,
		0.5,
		1.0,
		0,
		0,
	)

	if len(dets) != 1 {
		t.Fatalf("expected 1 detection, got %d", len(dets))
	}

	got := dets[0]

	if got.ClassID != 1 {
		t.Fatalf("expected class 1, got %d", got.ClassID)
	}

	if got.Confidence != 0.9 {
		t.Fatalf("expected confidence 0.9, got %f", got.Confidence)
	}

	if got.X1 != 40 || got.Y1 != 55 || got.X2 != 60 || got.Y2 != 65 {
		t.Fatalf("unexpected box: %+v", got)
	}
}

func TestDecodeYOLOv8PredsFirstLayout(t *testing.T) {
	numPreds := 2
	numClasses := 3
	valuesPerPred := 4 + numClasses

	output := make([]float32, valuesPerPred*numPreds)

	// pred 0: cx=50, cy=60, w=20, h=10
	base := 0 * valuesPerPred
	output[base+0] = 50
	output[base+1] = 60
	output[base+2] = 20
	output[base+3] = 10

	// pred 0 class scores: class 1 wins
	output[base+4] = 0.1
	output[base+5] = 0.9
	output[base+6] = 0.2

	// pred 1 below threshold
	base = 1 * valuesPerPred
	output[base+0] = 10
	output[base+1] = 10
	output[base+2] = 5
	output[base+3] = 5
	output[base+4] = 0.1
	output[base+5] = 0.2
	output[base+6] = 0.3

	dets := DecodeYOLOv8(
		output,
		numPreds,
		numClasses,
		YOLOv8LayoutPredsFirst,
		0.5,
		1.0,
		0,
		0,
	)

	if len(dets) != 1 {
		t.Fatalf("expected 1 detection, got %d", len(dets))
	}

	got := dets[0]

	if got.ClassID != 1 {
		t.Fatalf("expected class 1, got %d", got.ClassID)
	}

	if got.Confidence != 0.9 {
		t.Fatalf("expected confidence 0.9, got %f", got.Confidence)
	}

	if got.X1 != 40 || got.Y1 != 55 || got.X2 != 60 || got.Y2 != 65 {
		t.Fatalf("unexpected box: %+v", got)
	}
}
