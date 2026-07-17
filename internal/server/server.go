package server

import (
	"github.com/gin-gonic/gin"
	"yolo-go-inference/internal/stream"
)


type Server struct {

	Router *gin.Engine

	Runtime *stream.Runtime

}


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
