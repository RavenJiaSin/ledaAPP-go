package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"yolo-go-inference/internal/stream"
)


type StartInferenceRequest struct {

	CameraName string `json:"cam_name"`

	ModelName string `json:"model"`

	Interval int `json:"interval"`

}



func (s *Server) StartInference(
	c *gin.Context,
){


	var req StartInferenceRequest



	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(
			http.StatusBadRequest,
			APIResponse{
				Code:400,
				Msg:err.Error(),
			},
		)

		return
	}



	if req.CameraName == "" {

		c.JSON(
			http.StatusBadRequest,
			APIResponse{
				Code:400,
				Msg:"missing camera name",
			},
		)

		return
	}



	if req.ModelName == "" {

		req.ModelName = "yolov8"

	}



	// bind camera -> model

	err := s.Runtime.BindCamera(
		req.CameraName,
		req.ModelName,
	)


	if err != nil {

		c.JSON(
			http.StatusOK,
			APIResponse{
				Code:500,
				Msg:err.Error(),
			},
		)

		return

	}




	interval :=
		time.Duration(req.Interval) *
		time.Millisecond



	if interval <= 0 {

		interval =
			33*time.Millisecond

	}




	err =
		s.Runtime.Stream.StartStream(
			req.CameraName,
			stream.StreamConfig{
				Interval:interval,
			},
		)



	if err != nil {

		c.JSON(
			http.StatusOK,
			APIResponse{
				Code:500,
				Msg:err.Error(),
			},
		)

		return

	}



	c.JSON(
		http.StatusOK,
		APIResponse{
			Code:0,
			Result:gin.H{
				"camera":req.CameraName,
				"model":req.ModelName,
			},
		},
	)

}

func (s *Server) StopInference(
	c *gin.Context,
){

	var req struct {

		CameraName string `json:"cam_name"`

	}



	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(
			http.StatusBadRequest,
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
		APIResponse{
			Code:0,
		},
	)

}
