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

	loggerInstance := logger.L()

	// Initialize OpenTelemetry with timeout
	initCtx, initCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer initCancel()

	telemetry, err := otel.New(initCtx, cfg, loggerInstance)
	if err != nil {
		loggerInstance.Fatal("Failed to initialize telemetry", zap.Error(err))
	}
	defer func() {
		if err := telemetry.Shutdown(context.Background()); err != nil {
			loggerInstance.Error("Failed to shut down telemetry", zap.Error(err))
		}
	}()

	loggerInstance.Info("🚀 Starting API Gateway...",
		zap.Int("port", cfg.Port),
		zap.String("env", cfg.Env),
	)

	srv := server.New(cfg, loggerInstance, telemetry)

	// Create a channel to listen for OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Run server in a goroutine
	go func() {
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			loggerInstance.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	loggerInstance.Info("Server started and listening", zap.Int("port", cfg.Port))

	// Block until we receive a signal
	sig := <-sigChan
	loggerInstance.Info("Received shutdown signal", zap.String("signal", sig.String()))

	// Create a context with timeout for shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(shutdownCtx); err != nil {
		loggerInstance.Error("Server forced to shutdown", zap.Error(err))
		os.Exit(1)
	}

	loggerInstance.Info("Server gracefully stopped")
}
