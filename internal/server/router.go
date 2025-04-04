package server

import (
	"github.com/hohotang/shortlink-gateway/internal/handler"
	"github.com/hohotang/shortlink-gateway/internal/middleware"

	"github.com/gin-gonic/gin"
	_ "github.com/hohotang/shortlink-gateway/docs" // swagger generated docs
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	// API routes
	root := r.engine.Group("/")
	root.Use(r.middleware.Otel(), r.middleware.LoggingMiddleware())
	{
		root.POST("/shorten", r.shortlinkHandler.Shorten)
		root.GET("/:shortID", r.shortlinkHandler.Expand)
	}

	// Swagger documentation route
	r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
