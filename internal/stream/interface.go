package stream

import (
	
)


// StreamController defines stream operations
type StreamController interface {

	StartStream(
		cameraName string,
		cfg StreamConfig,
	) error


	StopStream(
		cameraName string,
	)


	StopAll()
}
