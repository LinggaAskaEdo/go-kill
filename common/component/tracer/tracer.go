package tracer

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
)

type Config struct {
	ServiceName    string        `yaml:"service_name"`
	ServiceVersion string        `yaml:"service_version"`
	Endpoint       string        `yaml:"endpoint"`
	Insecure       bool          `yaml:"insecure"`
	Timeout        time.Duration `yaml:"timeout"`
}

type Tracer interface {
	Stop(ctx context.Context) error
}

type tracerImpl struct {
	log      zerolog.Logger
	provider *sdktrace.TracerProvider
}

func Init(log zerolog.Logger, cfg Config) Tracer {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
		),
	)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to create resource for tracer")
		return nil
	}

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(cfg.Endpoint),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithTimeout(cfg.Timeout),
	)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to create OTLP exporter for tracer")
		return nil
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter,
			sdktrace.WithBatchTimeout(5*time.Second),
			sdktrace.WithMaxExportBatchSize(512),
		),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(tp)

	log.Print("Tracer initialized successfully")

	return &tracerImpl{
		log:      log,
		provider: tp,
	}
}

func (t *tracerImpl) Stop(ctx context.Context) error {
	if t.provider == nil {
		t.log.Print("Tracer provider is nil, nothing to shut down...")
		return nil
	}

	t.log.Print("Shutting down tracer...")
	if err := t.provider.Shutdown(ctx); err != nil {
		t.log.Printf("Error shutting down tracer: %v", err)
		return err
	}

	t.log.Print("Tracer shutdown complete...")

	return nil
}
