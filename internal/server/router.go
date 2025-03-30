package server

import (
	"shortlink-gateway/internal/handler"
	"shortlink-gateway/internal/middleware"

	"github.com/gin-gonic/gin"
)

func registerRoutes(r *gin.Engine, mw middleware.Middleware) {
	root := r.Group("/")
	root.Use(mw.Otel(), mw.LoggingMiddleware())
	{
		root.POST("/shorten", handler.Shorten)
		root.GET("/:shortID", handler.Expand)
	}
}
