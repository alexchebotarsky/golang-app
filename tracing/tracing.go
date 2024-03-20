package tracing

import (
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// If empty, default name will be used
const tracerName = ""

func Init(serviceName, environment string) error {
	// TODO: replace this exporter with actual exporter
	exporter, err := newFileExporter("traces.txt")
	if err != nil {
		return fmt.Errorf("error creating new file exporter: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.DeploymentEnvironment(environment),
			semconv.K8SDeploymentName(serviceName),
			semconv.K8SNamespaceName(environment),
		)),
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
