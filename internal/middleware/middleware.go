package middleware

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/hohotang/shortlink-gateway/internal/config"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type contextKey struct{}

var loggerKey = contextKey{}

// responseBodyWriter is a struct used to capture response content
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write overrides the original Write method to simultaneously write to the response and capture content
func (r *responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// WriteString overrides the original WriteString method
func (r *responseBodyWriter) WriteString(s string) (int, error) {
	r.body.WriteString(s)
	return r.ResponseWriter.WriteString(s)
}

type middleware struct {
	config *config.Config
	logger *zap.Logger
}

type Middleware interface {
	Otel() gin.HandlerFunc
	LoggingMiddleware() gin.HandlerFunc
}

func NewMiddleware(
	config *config.Config,
	logger *zap.Logger,
) Middleware {
	return &middleware{
		config: config,
		logger: logger,
	}
}

func (m *middleware) Otel() gin.HandlerFunc {
	return otelgin.Middleware(m.config.ServiceName)
}

func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func GetLogger(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(loggerKey).(*zap.Logger)
	if !ok {
		return zap.L() // fallback to global logger
	}
	return logger
}

func (m *middleware) LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Capture request body
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// Reset the request body so subsequent handlers can read it normally
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Create response capture writer
		responseBody := &bytes.Buffer{}
		responseWriter := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           responseBody,
		}
		c.Writer = responseWriter

		// Before handler execution
		path := c.Request.URL.Path
		method := c.Request.Method
		traceID := ""
		spanID := ""

		// Extract trace info from context
		if span := trace.SpanFromContext(c.Request.Context()); span != nil {
			sc := span.SpanContext()
			if sc.HasTraceID() {
				traceID = sc.TraceID().String()
			}
			if sc.HasSpanID() {
				spanID = sc.SpanID().String()
			}
			span.SetAttributes(
				attribute.String("http.request.body", string(requestBody)),
			)
		}

		ctx := WithLogger(c.Request.Context(), m.logger)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		// After handler execution
		latency := time.Since(start)
		status := c.Writer.Status()

		if span := trace.SpanFromContext(c.Request.Context()); span != nil {
			span.SetAttributes(
				attribute.String("http.response.body", responseBody.String()),
			)
		}

		// Log complete request and response information
		m.logger.Info("HTTP request",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.Duration("latency_ms", latency),
			zap.String("trace_id", traceID),
			zap.String("span_id", spanID),
			zap.String("request_headers", headersToString(c.Request.Header)),
			zap.String("request_body", string(requestBody)),
			zap.String("response_headers", headersToString(c.Writer.Header())),
			zap.String("response_body", responseBody.String()),
		)
	}
}

// headersToString converts HTTP headers to string
func headersToString(headers http.Header) string {
	buf := &bytes.Buffer{}
	for name, values := range headers {
		for _, value := range values {
			buf.WriteString(name)
			buf.WriteString(": ")
			buf.WriteString(value)
			buf.WriteString("\n")
		}
	}
	return buf.String()
}
