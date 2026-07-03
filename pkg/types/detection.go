package types

type Detection struct {
	X1 float32
	Y1 float32
	X2 float32
	Y2 float32

	Confidence float32
	ClassID    int
	ClassName  string
}

// 用於 NMS / 計算 IoU
type Box struct {
	X1 float32
	Y1 float32
	X2 float32
	Y2 float32
}

// cls_w_h 規則
type ClassRule struct {
	ClassName string

	MinWidth  float32
	MaxWidth  float32
	MinHeight float32
	MaxHeight float32
}

// Polygon alert area
type Point struct {
	X float32
	Y float32
}

type Polygon struct {
	Points []Point
}

// 推論結果（對齊你 Python API）
type InferenceResult struct {
	Detections []Detection
	Count      int
}
