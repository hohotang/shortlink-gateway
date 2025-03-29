package server

import (
	"fmt"

	"shortlink-gateway/internal/config"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	config *config.Config
}

func New(cfg *config.Config) *Server {
	r := gin.New()

	// 中介層：自訂 log middleware + recovery
	r.Use(otelgin.Middleware("api-gateway")) // 這會自動產生 root span
	r.Use(LoggingMiddleware())
	r.Use(gin.Recovery())

	registerRoutes(r)

	return &Server{
		router: r,
		config: cfg,
	}
}

func (s *Server) Run() error {
	addr := fmt.Sprintf(":%d", s.config.Port)
	return s.router.Run(addr)
}
