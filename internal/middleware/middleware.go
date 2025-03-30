package middleware

import (
	"shortlink-gateway/internal/config"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type middleware struct {
	config *config.Config
}

type Middleware interface {
	Otel() gin.HandlerFunc
	LoggingMiddleware() gin.HandlerFunc
}

func NewMiddleware(
	config *config.Config,
) Middleware {
	return &middleware{
		config: config,
	}
}

func (m *middleware) Otel() gin.HandlerFunc {
	return otelgin.Middleware(m.config.ServiceName)
}

func (m *middleware) LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 執行 handler 前
		path := c.Request.URL.Path
		method := c.Request.Method
		traceID := ""
		spanID := ""

		// 從 context 中取出 trace info
		if span := trace.SpanFromContext(c.Request.Context()); span != nil {
			sc := span.SpanContext()
			if sc.HasTraceID() {
				traceID = sc.TraceID().String()
			}
			if sc.HasSpanID() {
				spanID = sc.SpanID().String()
			}
		}

		c.Next()

		// 執行 handler 後
		latency := time.Since(start)
		status := c.Writer.Status()

		zap.L().Info("HTTP request",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.Duration("latency_ms", latency),
			zap.String("trace_id", traceID),
			zap.String("span_id", spanID),
		)
	}
}
