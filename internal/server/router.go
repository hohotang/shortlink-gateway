package server

import (
	"shortlink-gateway/internal/handler"

	"github.com/gin-gonic/gin"
)

func registerRoutes(r *gin.Engine) {
	r.POST("/shorten", handler.Shorten)
	r.GET("/:shortID", handler.Expand)
}
