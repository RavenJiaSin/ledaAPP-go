package preprocess

import (
	"image"
	"image/color"
	"math"
)

// 結果包含 tensor + resize資訊（後處理要用）
type PreprocessResult struct {
	Tensor []float32
	Shape  []int64 // [1,3,H,W]

	Scale float64
	PadX  float64
	PadY  float64

	InputW int
	InputH int
}

// 主入口
func Letterbox(img image.Image, targetSize int) PreprocessResult {
	origW := img.Bounds().Dx()
	origH := img.Bounds().Dy()

	// scale ratio
	r := math.Min(float64(targetSize)/float64(origW), float64(targetSize)/float64(origH))

	newW := int(math.Round(float64(origW) * r))
	newH := int(math.Round(float64(origH) * r))

	// resize
	resized := Resize(img, newW, newH)

	// padding
	padW := targetSize - newW
	padH := targetSize - newH

	padLeft := padW / 2
	padTop := padH / 2

	canvas := image.NewRGBA(image.Rect(0, 0, targetSize, targetSize))

	// 填充灰色（YOLO 預設 114）
	fillColor := color.RGBA{114, 114, 114, 255}
	for y := 0; y < targetSize; y++ {
		for x := 0; x < targetSize; x++ {
			canvas.Set(x, y, fillColor)
		}
	}

	// 貼圖
	for y := 0; y < newH; y++ {
		for x := 0; x < newW; x++ {
			canvas.Set(x+padLeft, y+padTop, resized.At(x, y))
		}
	}

	// 轉 tensor
	tensor := ImageToTensor(canvas)

	return PreprocessResult{
		Tensor: tensor,
		Shape:  []int64{1, 3, int64(targetSize), int64(targetSize)},

		Scale: r,
		PadX:  float64(padLeft),
		PadY:  float64(padTop),

		InputW: targetSize,
		InputH: targetSize,
	}
}

//
// 工具：Resize（最簡單版本，之後可優化成 bilinear）
//
func Resize(img image.Image, newW, newH int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))

	srcW := img.Bounds().Dx()
	srcH := img.Bounds().Dy()

	for y := 0; y < newH; y++ {
		for x := 0; x < newW; x++ {

			srcX := int(float64(x) * float64(srcW) / float64(newW))
			srcY := int(float64(y) * float64(srcH) / float64(newH))

			dst.Set(x, y, img.At(srcX, srcY))
		}
	}
	return dst
}

//
// 核心：Image → Tensor
// HWC → CHW + normalize + RGB
//
func ImageToTensor(img image.Image) []float32 {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	tensor := make([]float32, 3*w*h)

	// channel offset
	chSize := w * h

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {

			r, g, b, _ := img.At(x, y).RGBA()

			// RGBA 回傳 0~65535，要轉 0~255
			rf := float32(r>>8) / 255.0
			gf := float32(g>>8) / 255.0
			bf := float32(b>>8) / 255.0

			idx := y*w + x

			// CHW（注意順序：RGB）
			tensor[0*chSize+idx] = rf
			tensor[1*chSize+idx] = gf
			tensor[2*chSize+idx] = bf
		}
	}

	return tensor
}
