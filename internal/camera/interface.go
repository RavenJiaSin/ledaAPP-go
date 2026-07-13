package camera

import (
	"gocv.io/x/gocv"
)


// CameraController defines camera operations
type CameraController interface {

	Open(
		cameraName string,
		source interface{},
	) error


	Close(
		cameraName string,
	)


	CloseAll()


	LastFrame(
		cameraName string,
	) (gocv.Mat, bool)
}


