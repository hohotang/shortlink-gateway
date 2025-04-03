package server

import (
	"fmt"

	"shortlink-gateway/internal/config"
	"shortlink-gateway/internal/engine"
	"shortlink-gateway/internal/handler"
	"shortlink-gateway/internal/middleware"
	"shortlink-gateway/internal/service"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	config *config.Config
}

func New(cfg *config.Config) *Server {
	// Create engine
	engine := engine.NewEngine(cfg)

	// Create middleware
	mw := middleware.NewMiddleware(cfg)

	// Create services
	urlService := service.NewURLService()

	// Create handlers
	shortlinkHandler := handler.NewShortlinkHandler(urlService)

	// Create and initialize router
	router := NewRouter(engine, mw, shortlinkHandler)
	router.InitRoute()

	return &Server{
		router: engine,
		config: cfg,
	}
}

func (s *Server) Run() error {
	addr := fmt.Sprintf(":%d", s.config.Port)
	return s.router.Run(addr)
}
