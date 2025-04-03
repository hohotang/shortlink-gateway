package main

import (
	"context"
	"log"

	"github.com/hohotang/shortlink-gateway/internal/config"
	"github.com/hohotang/shortlink-gateway/internal/logger"
	"github.com/hohotang/shortlink-gateway/internal/otel"
	"github.com/hohotang/shortlink-gateway/internal/server"

	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	logger.Init("api-gateway", cfg.Env)
	defer logger.Sync()

	otel.Init(cfg)
	defer otel.Shutdown(context.Background())

	logger.L().Info("🚀 Starting API Gateway...",
		zap.Int("port", cfg.Port),
		zap.String("env", cfg.Env),
	)

	srv := server.New(cfg)
	if err := srv.Run(); err != nil {
		log.Fatalf("❌ Server exited with error: %v", err)
	}
}
