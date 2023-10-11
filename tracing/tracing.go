package tracing

import (
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// If empty, default name will be used
const tracerName = ""

func Init(serviceName string) error {
	exporter, err := newFileExporter("traces.txt")
	if err != nil {
		return fmt.Errorf("error creating new file exporter: %v", err)
	}

	tp, err := newTracerProvider(serviceName, exporter)
	if err != nil {
		return fmt.Errorf("error creating new tracer provider: %v", err)
	}

	otel.SetTracerProvider(tp)

	return nil
}

func newTracerProvider(name string, exporter tracesdk.SpanExporter) (*tracesdk.TracerProvider, error) {
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(name),
		)),
	)

	return tp, nil
}

func newFileExporter(filePath string) (tracesdk.SpanExporter, error) {
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
