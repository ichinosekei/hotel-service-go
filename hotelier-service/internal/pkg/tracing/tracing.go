package tracing

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func StartTracing(serviceName string) trace.Tracer {
	tracer := otel.Tracer(serviceName)
	return tracer
}
