package server

import (
	"fmt"

	"shortlink-gateway/internal/config"
	"shortlink-gateway/internal/engine"
	"shortlink-gateway/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	config *config.Config
}

func New(cfg *config.Config) *Server {

	// 中介層：自訂 log middleware + recovery
	// r.Use(otelgin.Middleware(cfg.ServiceName)) // 這會自動產生 root span
	// r.Use(LoggingMiddleware())
	// r.Use(gin.Recovery())

	mw := middleware.NewMiddleware(cfg)
	// engine
	engine := engine.NewEngine(cfg)

	registerRoutes(engine, mw)

	return &Server{
		router: engine,
		config: cfg,
	}
}

func (s *Server) Run() error {
	addr := fmt.Sprintf(":%d", s.config.Port)
	return s.router.Run(addr)
}
