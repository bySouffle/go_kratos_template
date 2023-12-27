package otlptrace

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"net/url"

	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
	"os"
	"time"
)

const (
	OTLPGRPC = "otlp.grpc"
	OTLPHTTP = "otlp.http"
)

type Config struct {
	ServiceName       string
	HostName          string
	CollectorEndpoint string
	Exporter          string
	TraceFilePath     string
	GrpcToken         string
}

// 设置应用资源
func (c *Config) newResource(ctx context.Context) *resource.Resource {
	// hostname默认值为本机主机名
	if len(c.HostName) == 0 {
		c.HostName, _ = os.Hostname()
	}

	r, err := resource.New(
		ctx,
		resource.WithFromEnv(),
		resource.WithProcess(), // runtime信息 process.runtime.name: go/gc, process.runtime.version: go1.20.1s
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(c.ServiceName),
			semconv.HostNameKey.String(c.HostName),
		),
	)

	if err != nil {
		panic("[otlp] Failed to create the GRPC OpenTelemetry resource")
	}
	return r
}

func (c *Config) newGrpcExporterAndSpanProcessor(ctx context.Context) (*otlptrace.Exporter, sdktrace.SpanProcessor) {
	headers := map[string]string{"Authentication": c.GrpcToken}

	traceExporter, err := otlptrace.New(
		ctx,
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(c.CollectorEndpoint),
			otlptracegrpc.WithHeaders(headers),
			otlptracegrpc.WithDialOption(grpc.WithBlock()),
			otlptracegrpc.WithCompressor(gzip.Name)),
	)

	if err != nil {
		panic("[otlp] Failed to create the GRPC OpenTelemetry trace exporter")
	}

	batchSpanProcessor := sdktrace.NewBatchSpanProcessor(traceExporter)

	return traceExporter, batchSpanProcessor
}

func (c *Config) newHTTPExporterAndSpanProcessor(ctx context.Context) (*otlptrace.Exporter, sdktrace.SpanProcessor) {
	u, err := url.Parse(c.CollectorEndpoint)
	if err != nil {
		panic("[otlp] CollectorEndpoint error")
	}
	traceExporter, err := otlptrace.New(ctx, otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(u.Host),
		otlptracehttp.WithURLPath(u.Path),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithCompression(1)))

	if err != nil {
		panic("[otlp] Failed to create the HTTP OpenTelemetry trace exporter")
	}

	batchSpanProcessor := sdktrace.NewBatchSpanProcessor(traceExporter)

	return traceExporter, batchSpanProcessor
}

// OpenTelemetry初始化方法
func (c *Config) InitOpenTelemetry() (*sdktrace.TracerProvider, func()) {
	ctx := context.Background()

	var traceExporter *otlptrace.Exporter
	var batchSpanProcessor sdktrace.SpanProcessor

	if c.Exporter == OTLPGRPC {
		traceExporter, batchSpanProcessor = c.newGrpcExporterAndSpanProcessor(ctx)
	} else if c.Exporter == OTLPHTTP {
		traceExporter, batchSpanProcessor = c.newHTTPExporterAndSpanProcessor(ctx)
	} else {
		panic("[otlp] Exporter type error")
	}

	otelResource := c.newResource(ctx)

	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(otelResource),
		sdktrace.WithSpanProcessor(batchSpanProcessor))

	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return traceProvider, func() {
		cxt, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		if err := traceExporter.Shutdown(cxt); err != nil {
			otel.Handle(err)
		}
	}
}
