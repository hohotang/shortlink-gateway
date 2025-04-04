package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/hohotang/shortlink-gateway/internal/config"
	"github.com/hohotang/shortlink-gateway/internal/engine"
	"github.com/hohotang/shortlink-gateway/internal/handler"
	"github.com/hohotang/shortlink-gateway/internal/middleware"
	"github.com/hohotang/shortlink-gateway/internal/service"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router     *gin.Engine
	config     *config.Config
	httpServer *http.Server
	urlService service.URLService
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
		grpcClient, err := service.NewURLGrpcClient(cfg.GrpcServerAddr, cfg)
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
		router:     engine,
		config:     cfg,
		urlService: urlService,
	}
}

func (s *Server) Run() error {
	addr := fmt.Sprintf(":%d", s.config.Port)

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	// First shutdown the HTTP server
	var err error
	if s.httpServer != nil {
		err = s.httpServer.Shutdown(ctx)
	}

	// Then close the URL service
	if s.urlService != nil {
		closeErr := s.urlService.Close()
		if err == nil {
			err = closeErr
		}
	}

	return err
}
