package instruments

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/skit-ai/vcore/env"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	OtelEnable        = env.Bool("OTEL_ENABLE", false)
	serviceName       = env.String("OTEL_SERVICE_NAME", "")
	collectorEndpoint = env.String("OTEL_COLLECTOR_ENDPOINT", "")
	useTls            = env.Bool("OTEL_USE_TLS", false)
)

// Initializes an OTLP exporter, and configures the corresponding trace and
// metric providers.
func InitProvider() (func(context.Context) error, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to create resource: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	opts := []grpc.DialOption{
		grpc.WithBlock(),
	}
	if useTls {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12})))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	conn, err := grpc.DialContext(ctx, collectorEndpoint, opts...)
	if err != nil {
		return nil, fmt.Errorf("Failed to create gRPC connection to collector: %w", err)
	}

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("Failed to create trace exporter: %w", err)
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Shutdown will flush any remaining spans and shut down the exporter.
	return tracerProvider.Shutdown, nil
}

// Extracts traceID from a Context
func ExtractTraceID(ctx context.Context) trace.TraceID {
	return trace.SpanFromContext(ctx).SpanContext().TraceID()
}

// Set a custom trace span name.
func SpanNameFormatter(_ string, r *http.Request) string {
	return fmt.Sprintf("%s %s %s", r.Method, r.Host, r.URL.Path)
}
