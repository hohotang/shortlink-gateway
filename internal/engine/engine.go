package engine

import (
	"github.com/hohotang/shortlink-gateway/internal/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func NewEngine(cfg *config.Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	// 設定 CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "DELETE", "PUT", "PATCH"}
	corsConfig.AllowHeaders = []string{
		"RecaptchaToken",
		"AccessToken",
		"Authorization",
		"Content-Type",
		"Upgrade",
		"Origin",
		"Connection",
		"Accept-Encoding",
		"Accept-Language",
		"Host",
		"Access-Control-Request-Method",
		"Access-Control-Request-Headers",
	}

	// 初始化 Gin 引擎
	server := gin.New()
	server.Use(cors.New(corsConfig))
	server.Use(gzip.Gzip(gzip.DefaultCompression))

	return server
}
