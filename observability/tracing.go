package observability

import (
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func ConfigureTracerProvider() *trace.TracerProvider {
	jaegerEndpoint := os.Getenv("JAEGER_ENDPOINT")
	if jaegerEndpoint == "" {
		jaegerEndpoint = fmt.Sprintf("%s/jaeger-api/api/traces", os.Getenv("GATEWAY_ADDR"))
	}

	exp, err := jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(jaegerEndpoint),
		),
	)
	if err != nil {
		panic(err)
	}

	tp := trace.NewTracerProvider(
		// WARNING: `tracesdk.WithSyncer` should be not used in production.
		// For production, you should use `tracesdk.WithBatcher`.
		trace.WithSyncer(exp),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("tickets"),
		)),
	)

	otel.SetTracerProvider(tp)

	// Don't forget this line! Omitting it will cause the trace to not be propagated via messages.
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp
}
