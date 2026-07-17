package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)



// GET /health
func (s *Server) HealthHandler(
	c *gin.Context,
) {

	c.JSON(
		http.StatusOK,
		gin.H{
			"status":"ok",
		},
	)

}
