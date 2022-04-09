package providers

import (
	"go-web-api/internal/globals"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func CreateJaegerProvider(url string) (*trace.TracerProvider, error) {

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))

	if err != nil {
		return nil, err
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(globals.TracerAppName),
		attribute.String("environment", "DEV"),
	)
	tp := trace.NewTracerProvider(trace.WithBatcher(exp), trace.WithResource(res))

	return tp, nil
}
