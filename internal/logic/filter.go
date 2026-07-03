package logic

import (
	"math"

	"yolo-go-inference/pkg/types"
)

// ============================
// 工具：計算 IoU
// ============================
func IoU(a, b types.Box) float32 {
	x1 := float32(math.Max(float64(a.X1), float64(b.X1)))
	y1 := float32(math.Max(float64(a.Y1), float64(b.Y1)))
	x2 := float32(math.Min(float64(a.X2), float64(b.X2)))
	y2 := float32(math.Min(float64(a.Y2), float64(b.Y2)))

	w := x2 - x1
	h := y2 - y1

	if w <= 0 || h <= 0 {
		return 0
	}

	inter := w * h

	areaA := (a.X2 - a.X1) * (a.Y2 - a.Y1)
	areaB := (b.X2 - b.X1) * (b.Y2 - b.Y1)

	union := areaA + areaB - inter
	if union <= 0 {
		return 0
	}

	return inter / union
}

// ============================
// 工具：點是否在多邊形內
// Ray Casting Algorithm
// ============================
func pointInPolygon(pt types.Point, poly types.Polygon) bool {
	inside := false
	n := len(poly.Points)

	for i, j := 0, n-1; i < n; j, i = i, i+1 {
		xi, yi := poly.Points[i].X, poly.Points[i].Y
		xj, yj := poly.Points[j].X, poly.Points[j].Y

		intersect := ((yi > pt.Y) != (yj > pt.Y)) &&
			(pt.X < (xj-xi)*(pt.Y-yi)/(yj-yi+1e-6)+xi)

		if intersect {
			inside = !inside
		}
	}

	return inside
}

// ============================
// 檢查 bbox 是否與 polygon 相交
// ============================
func boxIntersectsPolygon(box types.Box, poly types.Polygon) bool {
	// 1. 如果任一角在 polygon 內
	corners := []types.Point{
		{X: box.X1, Y: box.Y1},
		{X: box.X2, Y: box.Y1},
		{X: box.X2, Y: box.Y2},
		{X: box.X1, Y: box.Y2},
	}

	for _, c := range corners {
		if pointInPolygon(c, poly) {
			return true
		}
	}

	// 2. polygon 任一點在 box 內
	for _, p := range poly.Points {
		if p.X >= box.X1 && p.X <= box.X2 &&
			p.Y >= box.Y1 && p.Y <= box.Y2 {
			return true
		}
	}

	return false
}

// ============================
// Class rule filter
// ============================
func matchClassRule(det types.Detection, rules []types.ClassRule) bool {
	for _, r := range rules {
		if det.ClassName != r.ClassName {
			continue
		}

		w := det.X2 - det.X1
		h := det.Y2 - det.Y1

		if w < r.MinWidth || w > r.MaxWidth {
			return false
		}
		if h < r.MinHeight || h > r.MaxHeight {
			return false
		}
		return true
	}
	return true
}

// ============================
// 主 Filter Pipeline
// ============================
func FilterDetections(
	dets []types.Detection,
	polygons []types.Polygon,
	rules []types.ClassRule,
) []types.Detection {

	if len(dets) == 0 {
		return nil
	}

	filtered := make([]types.Detection, 0, len(dets))

	for _, det := range dets {

		// 1. class size rule
		if len(rules) > 0 && !matchClassRule(det, rules) {
			continue
		}

		box := types.Box{
			X1: det.X1,
			Y1: det.Y1,
			X2: det.X2,
			Y2: det.Y2,
		}

		// 2. polygon filtering
		if len(polygons) > 0 {
			ok := false
			for _, poly := range polygons {
				if boxIntersectsPolygon(box, poly) {
					ok = true
					break
				}
			}
			if !ok {
				continue
			}
		}

		filtered = append(filtered, det)
	}

	return filtered
}
