package server

import (
	"github.com/hohotang/shortlink-gateway/internal/handler"
	"github.com/hohotang/shortlink-gateway/internal/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

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
	// Metrics endpoint should be registered first and WITHOUT any middleware
	// that might interfere with Prometheus scraping
	r.engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API routes with middleware
	api := r.engine.Group("/")
	api.Use(r.middleware.Otel(), r.middleware.LoggingMiddleware(), r.middleware.MetricsMiddleware())
	{
		api.POST("v1/shorten", r.shortlinkHandler.Shorten)
		api.GET("v1/expand/:shortID", r.shortlinkHandler.Expand)
	}

	// Swagger documentation route
	r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
