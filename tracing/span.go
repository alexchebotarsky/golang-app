package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type Span struct {
	span trace.Span
}

func (s *Span) End() {
	s.span.End()
}

func (s *Span) SetName(name string) {
	s.span.SetName(name)
}

func (s *Span) SetTag(key, value string) {
	s.span.SetAttributes(attribute.String(key, value))
}

func (s *Span) RecordError(err error) {
	s.span.RecordError(err)
	s.span.SetStatus(codes.Error, fmt.Sprintf("%v", err))
}

func StartSpan(ctx context.Context, name string) (context.Context, Span) {
	var s Span

	ctx, s.span = otel.Tracer(tracerName).Start(ctx, name)

	return ctx, s
}

func SpanFromContext(ctx context.Context) Span {
	var s Span

	s.span = trace.SpanFromContext(ctx)

	return s
}

func NewCarrier(ctx context.Context) propagation.MapCarrier {
	carrier := propagation.MapCarrier{}

	otel.GetTextMapPropagator().Inject(ctx, &carrier)

	return carrier
}

func StartSpanFromCarrier(ctx context.Context, carrier propagation.MapCarrier, name string) (context.Context, Span) {
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	return StartSpan(ctx, name)
}
