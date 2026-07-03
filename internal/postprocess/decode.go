package postprocess

import (
	"yolo-go-inference/pkg/types"
)

type YOLOv8Layout int

const (
	YOLOv8LayoutChannelsFirst YOLOv8Layout = iota // [1, 84, 8400]
	YOLOv8LayoutPredsFirst                        // [1, 8400, 84]
)

// Decode YOLOv8 ONNX output
func DecodeYOLOv8(
	output []float32,
	numPreds int,
	numClasses int,
	layout YOLOv8Layout,
	confThreshold float32,

	scale float64,
	padX float64,
	padY float64,
) []types.Detection {
	var detections []types.Detection

	valuesPerPred := 4 + numClasses
	expectedSize := valuesPerPred * numPreds
	if len(output) < expectedSize {
		return detections
	}

	valueAt := func(channel, pred int) float32 {
		switch layout {
		case YOLOv8LayoutChannelsFirst:
			return output[channel*numPreds+pred]
		case YOLOv8LayoutPredsFirst:
			return output[pred*(4+numClasses)+channel]
		default:
			return 0
		}
	}

	for i := 0; i < numPreds; i++ {
		cx := valueAt(0, i)
		cy := valueAt(1, i)
		w := valueAt(2, i)
		h := valueAt(3, i)

		maxProb := float32(0)
		classID := -1

		for c := 0; c < numClasses; c++ {
			score := valueAt(4+c, i)
			if score > maxProb {
				maxProb = score
				classID = c
			}
		}

		if maxProb < confThreshold {
			continue
		}

		x1 := cx - w/2
		y1 := cy - h/2
		x2 := cx + w/2
		y2 := cy + h/2

		x1 = float32((float64(x1) - padX) / scale)
		y1 = float32((float64(y1) - padY) / scale)
		x2 = float32((float64(x2) - padX) / scale)
		y2 = float32((float64(y2) - padY) / scale)

		detections = append(detections, types.Detection{
			X1:         x1,
			Y1:         y1,
			X2:         x2,
			Y2:         y2,
			Confidence: maxProb,
			ClassID:    classID,
		})
	}

	return detections
}
