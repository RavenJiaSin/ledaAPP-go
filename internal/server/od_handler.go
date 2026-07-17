package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"yolo-go-inference/pkg/types"
)



// GET /api/infer_od/live?name=cam0
//
// OD inference result stream
func (s *Server) InferODLive(
	c *gin.Context,
) {


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


			frame,ok :=
				s.Runtime.Camera.LastFrame(
					cameraName,
				)


			if !ok {
				continue
			}



			data,err :=
				encodeJPEG(frame)



			frame.Close()



			if err != nil {
				continue
			}



			if _,err :=
				fmt.Fprintf(
					c.Writer,
					"--frame\r\nContent-Type: image/jpeg\r\nContent-Length: %d\r\n\r\n",
					len(data),
				);
				err != nil {

				return
			}



			if _,err :=
				c.Writer.Write(data);
				err != nil {

				return
			}



			if _,err :=
				fmt.Fprint(
					c.Writer,
					"\r\n",
				);
				err != nil {

				return
			}



			c.Writer.Flush()

		}

	}

}





// GET /api/infer_od/live_result?name=cam0
func (s *Server) InferODLiveResult(
	c *gin.Context,
) {


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


	result,ok :=
		s.Runtime.Store.GetCameraResult(
			cameraName,
		)



	if !ok {


		c.JSON(
			http.StatusOK,
			APIResponse{

				Code:0,

				Result:LiveInferResultResponse{

					LastResult:
						[]types.Detection{},

				},

			},
		)


		return

	}




	c.JSON(
		http.StatusOK,
		APIResponse{

			Code:0,


			Result:LiveInferResultResponse{

				LastResult:
					result.Result.Detections,


				Timestamp:
					result.Timestamp,

			},

		},
	)

}
