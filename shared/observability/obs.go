package observability

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

type Telemetry struct {
	cfg           *Config
	loggerHandler slog.Handler
	tracer        trace.Tracer
	shutdown      func(context.Context) error
}

type Config struct {
	Environment    string
	ServiceName    string
	ServiceVersion string
	Enabled        bool
	OtlpEndpoint   string
	SampleRatio    float64
}

func New(ctx context.Context, cfg Config) (*Telemetry, error) {

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{},
	))
	host, _ := os.Hostname()
	res, _ := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.ServiceVersion),
			semconv.DeploymentEnvironmentKey.String(cfg.Environment),
			semconv.ServiceInstanceID(host),
		),
	)

	// Setup tracing
	endpoint := cfg.OtlpEndpoint
	if endpoint == "" {
		endpoint = "http://otel-collector:4317"
	}
	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithRetry(otlptracegrpc.RetryConfig{
			Enabled:        true,
			MaxInterval:    2 * time.Second,
			MaxElapsedTime: 10 * time.Second,
		}),
	)
	if err != nil {
		return nil, err
	}
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(cfg.SampleRatio)),
	)
	otel.SetTracerProvider(traceProvider)
	tracer := traceProvider.Tracer(cfg.ServiceName)

	// Setup logging
	logExporter, _ := otlploggrpc.New(
		ctx,
		otlploggrpc.WithEndpoint("0.0.0.0:4317"),
	)
	logProvider := log.NewLoggerProvider(
		log.WithResource(res),
		log.WithProcessor(log.NewBatchProcessor(logExporter)),
	)
	loggerHandler := otelslog.NewHandler(
		cfg.ServiceName,
		otelslog.WithLoggerProvider(logProvider),
		otelslog.WithSchemaURL(semconv.SchemaURL),
	)

	global.SetLoggerProvider(logProvider)

	return &Telemetry{
		loggerHandler: loggerHandler,
		tracer:        tracer,
		shutdown: func(ctx context.Context) error { //TODO: replace with error group
			_ = logProvider.Shutdown(ctx)
			return traceProvider.Shutdown(ctx)
		},
	}, nil
}

func (t *Telemetry) LoggerHandler() slog.Handler {
	return t.loggerHandler
}

func (t *Telemetry) Tracer() trace.Tracer {
	return t.tracer
}

func (t *Telemetry) Shutdown(ctx context.Context) error {
	return t.shutdown(ctx)
}

func WrapHTTP(h any, operation string) any {
	if handler, ok := h.(interface {
		ServeHTTP(http.ResponseWriter, *http.Request)
	}); ok {
		return otelhttp.NewHandler(handler, operation)
	}
	return h
}
