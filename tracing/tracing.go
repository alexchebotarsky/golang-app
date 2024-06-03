package tracing

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func Init(ctx context.Context, serviceName, environment string, sampleRate float64) error {
	exporter, err := newFileExporter("traces.json")
	if err != nil {
		return fmt.Errorf("error creating new file exporter: %v", err)
	}

	res, err := resource.New(ctx,
		resource.WithTelemetrySDK(),
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.DeploymentEnvironment(environment),
		),
	)
	if err != nil {
		return fmt.Errorf("error creating new resource: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(sampleRate)),
	)

	otel.SetTracerProvider(tp)

	return nil
}

func newFileExporter(filePath string) (sdktrace.SpanExporter, error) {
	traceFile, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("error creating trace file: %v", err)
	}

	exporter, err := stdouttrace.New(stdouttrace.WithWriter(traceFile), stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, fmt.Errorf("error creating stdout trace exporter: %v", err)
	}

	return exporter, nil
}
