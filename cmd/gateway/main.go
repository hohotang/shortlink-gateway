package main

// @title           Shortlink Gateway API
// @version         1.0
// @description     A URL shortening service API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.hohotang.com/support
// @contact.email  support@hohotang.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hohotang/shortlink-gateway/internal/config"
	"github.com/hohotang/shortlink-gateway/internal/logger"
	"github.com/hohotang/shortlink-gateway/internal/otel"
	"github.com/hohotang/shortlink-gateway/internal/server"

	_ "github.com/hohotang/shortlink-gateway/docs" // Import swagger docs
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	logger.Init("api-gateway", cfg.Env)
	defer logger.Sync()

	otel.Init(cfg)
	defer otel.Shutdown(context.Background())
	logger := logger.L()

	logger.Info("🚀 Starting API Gateway...",
		zap.Int("port", cfg.Port),
		zap.String("env", cfg.Env),
	)

	srv := server.New(cfg, logger)

	// Create a channel to listen for OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Run server in a goroutine
	go func() {
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	logger.Info("Server started and listening", zap.Int("port", cfg.Port))

	// Block until we receive a signal
	sig := <-sigChan
	logger.Info("Received shutdown signal", zap.String("signal", sig.String()))

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
		os.Exit(1)
	}

	logger.Info("Server gracefully stopped")
}
