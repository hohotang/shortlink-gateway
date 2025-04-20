package otel

import (
	"context"
	"time"

	"github.com/hohotang/shortlink-gateway/internal/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/zap"
)

// Telemetry contains all OpenTelemetry components
type Telemetry struct {
	TracerProvider *trace.TracerProvider
	MeterProvider  *sdkmetric.MeterProvider
	Metrics        *Metrics
	Logger         *zap.Logger
}

// Metrics contains all metric instruments
type Metrics struct {
	RequestCounter  metric.Int64Counter
	RequestDuration metric.Float64Histogram
}

// New creates a new Telemetry instance with all components initialized
func New(ctx context.Context, cfg *config.Config, logger *zap.Logger) (*Telemetry, error) {
	// Create a shared resource for both trace and metrics
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(cfg.ServiceName),
	)

	// Initialize tracing
	tp, err := initTracing(ctx, cfg, res)
	if err != nil {
		return nil, err
	}

	// Initialize metrics
	mp, err := initMetrics(res)
	if err != nil {
		return nil, err
	}

	// Initialize metric instruments
	metrics, err := initMetricInstruments(mp)
	if err != nil {
		return nil, err
	}

	// Set up propagation for cross-service context transfer
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return &Telemetry{
		TracerProvider: tp,
		MeterProvider:  mp,
		Metrics:        metrics,
		Logger:         logger,
	}, nil
}

func initTracing(ctx context.Context, cfg *config.Config, res *resource.Resource) (*trace.TracerProvider, error) {
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(cfg.TracesEndpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	return tp, nil
}

func initMetrics(res *resource.Resource) (*sdkmetric.MeterProvider, error) {
	exporter, err := prometheus.New(prometheus.WithRegisterer(nil))
	if err != nil {
		return nil, err
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exporter),
		sdkmetric.WithResource(res),
	)

	otel.SetMeterProvider(mp)
	return mp, nil
}

func initMetricInstruments(mp *sdkmetric.MeterProvider) (*Metrics, error) {
	meter := mp.Meter("http-server")

	requestCounter, err := meter.Int64Counter(
		"http_server_requests_total",
		metric.WithDescription("Total number of HTTP requests"),
	)
	if err != nil {
		return nil, err
	}

	requestDuration, err := meter.Float64Histogram(
		"http_server_request_duration_seconds",
		metric.WithDescription("Duration of HTTP requests in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	return &Metrics{
		RequestCounter:  requestCounter,
		RequestDuration: requestDuration,
	}, nil
}

// Shutdown flushes and shuts down the trace and metric providers
func (t *Telemetry) Shutdown(ctx context.Context) error {
	var err error

	// Shutdown trace provider
	if t.TracerProvider != nil {
		tpCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if shutdownErr := t.TracerProvider.Shutdown(tpCtx); shutdownErr != nil {
			t.Logger.Error("Failed to shut down trace provider",
				zap.Error(shutdownErr))
			err = shutdownErr
		}
	}

	// Shutdown metric provider
	if t.MeterProvider != nil {
		mpCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if shutdownErr := t.MeterProvider.Shutdown(mpCtx); shutdownErr != nil && err == nil {
			t.Logger.Error("Failed to shut down metric provider",
				zap.Error(shutdownErr))
			err = shutdownErr
		}
	}

	return err
}
