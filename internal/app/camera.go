package app

import (
	"time"

	"yolo-go-inference/internal/stream"
)


func initCameras(
	runtime *stream.Runtime,
	cfg interface{},
) error {


	err:=runtime.Camera.Open(
		"cam0",
		0,
	)

	if err!=nil{
		return err
	}


	runtime.BindCamera(
		"cam0",
		"yolov8",
	)


	return runtime.Stream.StartStream(
		"cam0",
		stream.StreamConfig{
			Interval:
				33*time.Millisecond,
		},
	)
}
