// internal\postprocess\nms.go
package postprocess

import (
	"sort"

	"yolo-go-inference/pkg/types"
)

// 計算 IoU (Intersection over Union)
func IoU(a, b types.Box) float32 {
	// intersection
	x1 := max(a.X1, b.X1)
	y1 := max(a.Y1, b.Y1)
	x2 := min(a.X2, b.X2)
	y2 := min(a.Y2, b.Y2)

	interW := x2 - x1
	interH := y2 - y1

	if interW <= 0 || interH <= 0 {
		return 0.0
	}

	intersection := interW * interH

	// union
	areaA := (a.X2 - a.X1) * (a.Y2 - a.Y1)
	areaB := (b.X2 - b.X1) * (b.Y2 - b.Y1)

	union := areaA + areaB - intersection

	if union <= 0 {
		return 0.0
	}

	return intersection / union
}

// 將 Detection 轉成 Box（方便計算）
func toBox(d types.Detection) types.Box {
	return types.Box{
		X1: d.X1,
		Y1: d.Y1,
		X2: d.X2,
		Y2: d.Y2,
	}
}

// Non-Maximum Suppression (class-wise)
// input:
//   detections: decode 後的所有框
//   iouThreshold: IoU 超過此值會被抑制
//
// output:
//   經過 NMS 的 detections
func NMS(detections []types.Detection, iouThreshold float32) []types.Detection {

	// 1. 依 class 分組
	classMap := make(map[int][]types.Detection)

	for _, det := range detections {
		classMap[det.ClassID] = append(classMap[det.ClassID], det)
	}

	var finalDetections []types.Detection

	// 2. 對每個 class 做 NMS
	for _, dets := range classMap {

		// 2.1 依 confidence 排序 (descending)
		sort.Slice(dets, func(i, j int) bool {
			return dets[i].Confidence > dets[j].Confidence
		})

		var selected []types.Detection

		// 2.2 greedy NMS
		for len(dets) > 0 {
			// 取最高分
			current := dets[0]
			selected = append(selected, current)

			var remaining []types.Detection

			for i := 1; i < len(dets); i++ {
				iou := IoU(toBox(current), toBox(dets[i]))

				if iou < iouThreshold {
					remaining = append(remaining, dets[i])
				}
				// else: suppress
			}

			dets = remaining
		}

		finalDetections = append(finalDetections, selected...)
	}

	return finalDetections
}

// ----------------- helper -----------------

func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}
