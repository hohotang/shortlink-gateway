package server

import (
	"fmt"
	"log"

	"github.com/hohotang/shortlink-gateway/internal/config"
	"github.com/hohotang/shortlink-gateway/internal/engine"
	"github.com/hohotang/shortlink-gateway/internal/handler"
	"github.com/hohotang/shortlink-gateway/internal/middleware"
	"github.com/hohotang/shortlink-gateway/internal/service"

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
	var urlService service.URLService

	// Choose between mock or real gRPC client based on configuration
	if cfg.UseGrpc {
		grpcClient, err := service.NewURLGrpcClient(cfg.GrpcServerAddr)
		if err != nil {
			log.Printf("Failed to create gRPC client: %v, falling back to mock", err)
			urlService = service.NewURLService()
		} else {
			urlService = service.NewURLServiceWithClient(grpcClient)
		}
	} else {
		urlService = service.NewURLService() // Use default mock implementation
	}

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
