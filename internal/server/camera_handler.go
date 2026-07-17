package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gocv.io/x/gocv"
)


// POST /api/camera/open
func (s *Server) OpenCamera(c *gin.Context) {

	var req struct {

		Name string `json:"name"`

		URI string `json:"uri"`

	}


	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}



	err := s.Runtime.Camera.Open(
		req.Name,
		req.URI,
	)


	if err != nil {

		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}



	c.JSON(
		http.StatusOK,
		gin.H{
			"status":"opened",
		},
	)

}



// POST /api/camera/close
func (s *Server) CloseCamera(c *gin.Context) {


	var req struct {

		Name string `json:"name"`

	}


	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error":err.Error(),
			},
		)

		return
	}



	s.Runtime.Camera.Close(
		req.Name,
	)



	c.JSON(
		http.StatusOK,
		gin.H{
			"status":"closed",
		},
	)

}



// GET /api/camera/check?name=cam0
func (s *Server) CheckCamera(c *gin.Context) {


	name := c.Query("name")


	if name == "" {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error":"missing camera name",
			},
		)

		return
	}



	_, ok :=
		s.Runtime.Camera.LastFrame(name)



	c.JSON(
		http.StatusOK,
		gin.H{
			"alive":ok,
		},
	)

}



// GET /api/camera/frame?name=cam0
func (s *Server) CameraFrame(c *gin.Context) {


	cameraName :=
		c.Query("name")


	if cameraName == "" {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error":"missing camera name",
			},
		)

		return
	}



	frame, ok :=
		s.Runtime.Camera.LastFrame(
			cameraName,
		)


	if !ok {

		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error":"frame unavailable",
			},
		)

		return
	}



	defer frame.Close()



	data, err :=
		encodeJPEG(frame)



	if err != nil {

		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error":err.Error(),
			},
		)

		return
	}



	c.Data(
		http.StatusOK,
		"image/jpeg",
		data,
	)

}



// GET /api/camera/live?name=cam0
func (s *Server) CameraLive(c *gin.Context) {


	cameraName :=
		c.Query("name")


	if cameraName == "" {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error":"missing camera name",
			},
		)

		return
	}



	c.Header(
		"Content-Type",
		"multipart/x-mixed-replace; boundary=frame",
	)


	c.Header(
		"Cache-Control",
		"no-cache",
	)



	ticker :=
		time.NewTicker(
			40*time.Millisecond,
		)


	defer ticker.Stop()



	for {

		select {


		case <-c.Request.Context().Done():

			return



		case <-ticker.C:


			frame, ok :=
				s.Runtime.Camera.LastFrame(
					cameraName,
				)


			if !ok {
				continue
			}



			data, err :=
				encodeJPEG(frame)


			frame.Close()



			if err != nil {
				continue
			}



			_, err =
				fmt.Fprintf(
					c.Writer,
					"--frame\r\nContent-Type: image/jpeg\r\nContent-Length: %d\r\n\r\n",
					len(data),
				)


			if err != nil {
				return
			}



			_, err =
				c.Writer.Write(data)


			if err != nil {
				return
			}



			fmt.Fprint(
				c.Writer,
				"\r\n",
			)



			c.Writer.Flush()

		}

	}

}



// encode OpenCV Mat -> JPEG bytes
func encodeJPEG(
	mat gocv.Mat,
) ([]byte,error) {


	jpegMat,err :=
		gocv.IMEncode(
			".jpg",
			mat,
		)


	if err != nil {

		return nil,err

	}



	defer jpegMat.Close()



	// copy C memory -> Go memory

	data :=
		jpegMat.GetBytes()



	out :=
		make([]byte,len(data))


	copy(
		out,
		data,
	)



	return out,nil

}
