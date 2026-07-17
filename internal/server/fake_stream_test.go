package server

import (
	"yolo-go-inference/internal/stream"
)


type fakeStream struct {
	started bool
	stopped bool
}


func (f *fakeStream) StartStream(
	name string,
	cfg stream.StreamConfig,
) error {

	f.started=true

	return nil
}



func (f *fakeStream) StopStream(
	name string,
){

	f.stopped=true

}



func (f *fakeStream) StopAll(){

}
