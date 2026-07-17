package camera

import (
	"sync"

	"gocv.io/x/gocv"
)


type FakeCamera struct {

	mu sync.RWMutex

	frames map[string]gocv.Mat
}



func NewFakeCamera() *FakeCamera {

	return &FakeCamera{
		frames: make(map[string]gocv.Mat),
	}
}



func (f *FakeCamera) Open(
	name string,
	source interface{},
) error {


	mat := gocv.NewMatWithSize(
		640,
		640,
		gocv.MatTypeCV8UC3,
	)


	f.mu.Lock()
	defer f.mu.Unlock()


	f.frames[name] = mat


	return nil
}



func (f *FakeCamera) Close(
	name string,
){

	f.mu.Lock()
	defer f.mu.Unlock()


	if frame, ok := f.frames[name]; ok {

		frame.Close()

		delete(
			f.frames,
			name,
		)
	}

}



func (f *FakeCamera) CloseAll(){

	f.mu.Lock()
	defer f.mu.Unlock()


	for name, frame := range f.frames {

		frame.Close()

		delete(
			f.frames,
			name,
		)
	}

}



func (f *FakeCamera) LastFrame(
	name string,
)(gocv.Mat,bool){


	f.mu.RLock()
	defer f.mu.RUnlock()


	frame,ok :=
		f.frames[name]


	if !ok || frame.Empty(){
		return gocv.Mat{},false
	}


	return frame.Clone(),true
}
