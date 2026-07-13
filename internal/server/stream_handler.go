package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"yolo-go-inference/internal/stream"
	"yolo-go-inference/pkg/types"
)


// -----------------------------
// Stream Handler
// -----------------------------

type APIResponse struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg,omitempty"`
	Result interface{} `json:"result,omitempty"`
}

type LiveInferResponse struct {
	ResultImage string `json:"result_image"`
}

type LiveInferOdResultResponse struct {
	LastResult interface{} `json:"last_result"`
	Timestamp  interface{} `json:"timestamp"`
}

// StartLiveInferOD starts YOLO streaming inference
// POST /api/start_live_infer_od
func (s *Server) StartLiveInferOD(c *gin.Context) {

	var req struct {
		CameraName string `json:"cam_name"`
		Interval   int    `json:"interval"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code: 400,
			Msg:  err.Error(),
		})
		return
	}


	interval := time.Duration(req.Interval) * time.Millisecond

	if interval <= 0 {
		interval = 33 * time.Millisecond
	}


	err := s.Runtime.Stream.StartStream(
		req.CameraName,
		stream.StreamConfig{
			Interval: interval,

			// 後續由 request 加入
			// Polygons:
			// ClassRules:
		},
	)


	if err != nil {
		c.JSON(
			http.StatusOK,
			APIResponse{
				Code: 500,
				Msg: err.Error(),
			},
		)
		return
	}


	c.JSON(
		http.StatusOK,
		APIResponse{
			Code: 0,
			Result: LiveInferResponse{
				ResultImage:
					fmt.Sprintf(
						"/api/infer_od/live?name=%s",
						req.CameraName,
					),
			},
		},
	)
}



// StopLiveInferOD stops inference
// POST /api/stop_live_infer_od
func (s *Server) StopLiveInferOD(c *gin.Context) {


	var req struct {
		CameraName string `json:"cam_name"`
	}


	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest,
			APIResponse{
				Code:400,
				Msg:err.Error(),
			},
		)
		return
	}


	s.Runtime.Stream.StopStream(
		req.CameraName,
	)


	c.JSON(
		http.StatusOK,
		APIResponse{},
	)
}



// InferODLive returns MJPEG stream
// GET /api/infer_od/live?name=cam0
func (s *Server) InferODLive(c *gin.Context) {

	cameraName := c.Query("name")

	if cameraName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing camera name"})
		return
	}

	c.Header("Content-Type", "multipart/x-mixed-replace; boundary=frame")

	ticker := time.NewTicker(40 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {

		case <-c.Request.Context().Done():
			return

		case <-ticker.C:

			frame, ok := s.Runtime.Camera.LastFrame(cameraName)
			if !ok {
				continue
			}

			data, err := encodeJPEG(frame)
			frame.Close()

			if err != nil {
				continue
			}

			if _, err := fmt.Fprintf(
				c.Writer,
				"--frame\r\nContent-Type: image/jpeg\r\nContent-Length: %d\r\n\r\n",
				len(data),
			); err != nil {
				return
			}

			if _, err := c.Writer.Write(data); err != nil {
				return
			}

			if _, err := fmt.Fprint(c.Writer, "\r\n"); err != nil {
				return
			}

			c.Writer.Flush()
		}
	}
}



// InferODLiveResult returns latest detection result
// GET /api/infer_od/live_result?name=cam0
func (s *Server) InferODLiveResult(c *gin.Context) {


	cameraName :=
		c.Query("name")


	result, ok :=
		s.Runtime.Store.GetCameraResult(cameraName)


	if !ok {

		c.JSON(
			http.StatusOK,
			APIResponse{
				Code:0,
				Result:LiveInferOdResultResponse{
					LastResult:[]types.Detection{},
				},
			},
		)

		return
	}



	c.JSON(
		http.StatusOK,
		APIResponse{
			Code:0,
			Result:LiveInferOdResultResponse{
				LastResult:
					result.Result.Detections,

				Timestamp:
					result.Timestamp,
			},
		},
	)
}
