package server

import (
	"yolo-go-inference/internal/store"
	"yolo-go-inference/internal/stream"

	"github.com/gin-gonic/gin"
)


func setupTestServer() *Server {

	gin.SetMode(gin.TestMode)


	runtime := stream.NewRuntime()


	runtime.Store =
		store.NewPipelineStore()


	srv := &Server{
		Runtime: runtime,
	}


	srv.Router = gin.New()


	return srv
}
