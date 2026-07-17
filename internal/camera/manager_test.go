package camera
import (
	"testing"
	
)
func TestCameraCloseLifecycle(t *testing.T){

	cam :=
		NewCameraManager()


	err :=
		cam.Open(
			"cam0",
			0,
		)


	if err != nil {
		t.Skip("camera unavailable")
	}



	cam.Close(
		"cam0",
	)



	cam.mu.RLock()

	defer cam.mu.RUnlock()



	if _,ok :=
		cam.cameras["cam0"];
		ok {

		t.Fatal(
			"camera still exists",
		)

	}

}
