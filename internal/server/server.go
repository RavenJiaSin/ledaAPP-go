package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"yolo-go-inference/internal/stream"
)


type Server struct {

	Router *gin.Engine

	Runtime *stream.Runtime
}


// 建立 Server

func New(
	runtime *stream.Runtime,
) *Server {
	s := &Server{
		Runtime: runtime,
	}
	s.Router = gin.Default()
	s.registerRoutes()

	return s
}




func (s *Server) registerRoutes() {


	// -----------------------------
	// health
	// -----------------------------

	s.Router.GET(
		"/health",
		s.HealthHandler,
	)



	// -----------------------------
	// camera
	// -----------------------------

	s.Router.POST(
		"/api/camera/open",
		s.OpenCamera,
	)


	s.Router.POST(
		"/api/camera/close",
		s.CloseCamera,
	)


	s.Router.GET(
		"/api/camera/check",
		s.CheckCamera,
	)


	s.Router.GET(
		"/api/camera/frame",
		s.CameraFrame,
	)


	s.Router.GET(
		"/api/camera/live",
		s.CameraLive,
	)



	// -----------------------------
	// stream inference
	// -----------------------------


	s.Router.POST(
		"/api/start_live_infer_od",
		s.StartLiveInferOD,
	)



	s.Router.POST(
		"/api/stop_live_infer_od",
		s.StopLiveInferOD,
	)



	s.Router.GET(
		"/api/infer_od/live",
		s.InferODLive,
	)



	s.Router.GET(
		"/api/infer_od/live_result",
		s.InferODLiveResult,
	)

}




func (s *Server) HealthHandler(
	c *gin.Context,
){

	c.JSON(
		http.StatusOK,
		gin.H{
			"status":"ok",
		},
	)

}
