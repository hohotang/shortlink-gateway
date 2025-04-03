package server

import (
	"shortlink-gateway/internal/handler"
	"shortlink-gateway/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Router struct {
	engine           *gin.Engine
	middleware       middleware.Middleware
	shortlinkHandler *handler.ShortlinkHandler
}

func NewRouter(engine *gin.Engine, mw middleware.Middleware, shortlinkHandler *handler.ShortlinkHandler) *Router {
	return &Router{
		engine:           engine,
		middleware:       mw,
		shortlinkHandler: shortlinkHandler,
	}
}

// InitRoute registers all the routes
func (r *Router) InitRoute() {
	root := r.engine.Group("/")
	root.Use(r.middleware.Otel(), r.middleware.LoggingMiddleware())
	{
		root.POST("/shorten", r.shortlinkHandler.Shorten)
		root.GET("/:shortID", r.shortlinkHandler.Expand)
	}
}
