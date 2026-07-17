package camera

import (
	"errors"
	"sync"

	"gocv.io/x/gocv"
)


type CameraManager struct {

	mu sync.RWMutex


	// cameraName -> VideoCapture
	cameras map[string]*gocv.VideoCapture


	// cameraName -> last frame
	lastFrames map[string]gocv.Mat


	// cameraName -> running
	running map[string]bool


	// cameraName -> capture worker waitgroup
	workers map[string]*sync.WaitGroup
}



func NewCameraManager() *CameraManager {

	return &CameraManager{

		cameras:
			make(map[string]*gocv.VideoCapture),


		lastFrames:
			make(map[string]gocv.Mat),


		running:
			make(map[string]bool),


		workers:
			make(map[string]*sync.WaitGroup),
	}
}




func (m *CameraManager) Open(
	cameraName string,
	source interface{},
) error {


	m.mu.Lock()
	defer m.mu.Unlock()



	if _, exists :=
		m.cameras[cameraName];
		exists {

		return errors.New(
			"camera already opened",
		)

	}



	var (
		cap *gocv.VideoCapture
		err error
	)



	switch v := source.(type) {

	case int:

		cap, err =
			gocv.OpenVideoCapture(v)


	case string:

		if v == "0" {

			cap, err =
				gocv.OpenVideoCapture(0)

		} else {

			cap, err =
				gocv.OpenVideoCapture(v)

		}


	default:

		return errors.New(
			"invalid camera source type",
		)
	}



	if err != nil {
		return err
	}



	if !cap.IsOpened() {

		return errors.New(
			"failed to open camera",
		)
	}



	m.cameras[cameraName] =
		cap


	m.lastFrames[cameraName] =
		gocv.NewMat()


	m.running[cameraName] =
		true



	wg := &sync.WaitGroup{}

	wg.Add(1)

	m.workers[cameraName] =
		wg



	go m.captureLoop(
		cameraName,
		cap,
		wg,
	)



	return nil
}





func (m *CameraManager) captureLoop(
	name string,
	cap *gocv.VideoCapture,
	wg *sync.WaitGroup,
){

	defer wg.Done()



	for {


		m.mu.RLock()

		running :=
			m.running[name]

		m.mu.RUnlock()



		if !running {
			return
		}



		mat :=
			gocv.NewMat()



		ok :=
			cap.Read(&mat)



		if !ok {

			mat.Close()

			return
		}



		if mat.Empty(){

			mat.Close()

			continue
		}



		clone :=
			mat.Clone()


		mat.Close()



		m.mu.Lock()



		if !m.running[name] {

			clone.Close()

			m.mu.Unlock()

			return
		}



		if old,ok :=
			m.lastFrames[name];
			ok {

			old.Close()

		}



		m.lastFrames[name] =
			clone



		m.mu.Unlock()

	}

}







func (m *CameraManager) LastFrame(
	cameraName string,
) (gocv.Mat,bool){


	m.mu.RLock()

	defer m.mu.RUnlock()



	frame,ok :=
		m.lastFrames[cameraName]


	if !ok || frame.Empty(){

		return gocv.Mat{},false
	}



	return frame.Clone(),true

}








func (m *CameraManager) Close(
	cameraName string,
){



	m.mu.Lock()



	delete(
		m.running,
		cameraName,
	)



	cap,ok :=
		m.cameras[cameraName]


	wg :=
		m.workers[cameraName]



	delete(
		m.cameras,
		cameraName,
	)


	delete(
		m.workers,
		cameraName,
	)



	m.mu.Unlock()



	if ok {

		// 立即釋放 device
		cap.Close()

	}



	if wg != nil {

		wg.Wait()

	}



	m.mu.Lock()

	defer m.mu.Unlock()



	if frame,ok :=
		m.lastFrames[cameraName];
		ok {


		frame.Close()


		delete(
			m.lastFrames,
			cameraName,
		)

	}

}







func (m *CameraManager) CloseAll(){


	m.mu.Lock()



	caps :=
		make([]*gocv.VideoCapture,0)


	for name,cap :=
		range m.cameras {


		caps =
			append(
				caps,
				cap,
			)


		delete(
			m.running,
			name,
		)


		delete(
			m.cameras,
			name,
		)

	}



	workers :=
		make([]*sync.WaitGroup,0)


	for _,wg :=
		range m.workers {

		workers =
			append(
				workers,
				wg,
			)

	}


	m.workers =
		make(map[string]*sync.WaitGroup)



	frames :=
		make([]gocv.Mat,0)



	for name,frame :=
		range m.lastFrames {


		frames =
			append(
				frames,
				frame,
			)


		delete(
			m.lastFrames,
			name,
		)

	}



	m.mu.Unlock()




	for _,cap :=
		range caps {

		cap.Close()

	}



	for _,wg :=
		range workers {

		wg.Wait()

	}



	for _,frame :=
		range frames {

		frame.Close()

	}

}







func (m *CameraManager) InjectFrame(
	cameraName string,
	mat gocv.Mat,
){

	m.mu.Lock()

	defer m.mu.Unlock()



	clone :=
		mat.Clone()



	if old,ok :=
		m.lastFrames[cameraName];
		ok {

		old.Close()

	}



	m.lastFrames[cameraName] =
		clone
}
