package preprocess

import (
	"image"
	"image/color"
	"testing"
)

func TestLetterboxWideImage(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 320, 160))
	fillImage(img, color.RGBA{255, 0, 0, 255})

	got := Letterbox(img, 640)

	if len(got.Tensor) != 3*640*640 {
		t.Fatalf("expected tensor length %d, got %d", 3*640*640, len(got.Tensor))
	}

	if len(got.Shape) != 4 {
		t.Fatalf("expected shape rank 4, got %d", len(got.Shape))
	}

	expectedShape := []int64{1, 3, 640, 640}
	for i := range expectedShape {
		if got.Shape[i] != expectedShape[i] {
			t.Fatalf("expected shape %v, got %v", expectedShape, got.Shape)
		}
	}

	if got.Scale != 2.0 {
		t.Fatalf("expected scale 2.0, got %f", got.Scale)
	}

	if got.PadX != 0 {
		t.Fatalf("expected pad x 0, got %f", got.PadX)
	}

	if got.PadY != 160 {
		t.Fatalf("expected pad y 160, got %f", got.PadY)
	}
}

func TestLetterboxTallImage(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 160, 320))
	fillImage(img, color.RGBA{0, 255, 0, 255})

	got := Letterbox(img, 640)

	if got.Scale != 2.0 {
		t.Fatalf("expected scale 2.0, got %f", got.Scale)
	}

	if got.PadX != 160 {
		t.Fatalf("expected pad x 160, got %f", got.PadX)
	}

	if got.PadY != 0 {
		t.Fatalf("expected pad y 0, got %f", got.PadY)
	}
}

func TestImageToTensorUsesCHWRGBNormalized(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 2, 1))
	img.Set(0, 0, color.RGBA{255, 0, 0, 255})
	img.Set(1, 0, color.RGBA{0, 128, 255, 255})

	got := ImageToTensor(img)

	if len(got) != 6 {
		t.Fatalf("expected tensor length 6, got %d", len(got))
	}

	// R channel
	if got[0] != 1.0 || got[1] != 0.0 {
		t.Fatalf("unexpected R channel: %v", got[0:2])
	}

	// G channel
	if got[2] != 0.0 || got[3] < 0.50 || got[3] > 0.51 {
		t.Fatalf("unexpected G channel: %v", got[2:4])
	}

	// B channel
	if got[4] != 0.0 || got[5] != 1.0 {
		t.Fatalf("unexpected B channel: %v", got[4:6])
	}
}

func fillImage(img *image.RGBA, c color.RGBA) {
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			img.Set(x, y, c)
		}
	}
}
