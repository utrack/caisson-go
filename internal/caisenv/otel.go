package caisenv

import (
	"context"
	"log"
	"log/slog"

	"github.com/utrack/caisson-go/errors"
	"github.com/utrack/caisson-go/pkg/plconfig"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func initTracer() func(context.Context) error {

	cfg := plconfig.Get()

	var exporter sdktrace.SpanExporter
	var err error

	if !cfg.Otel.Enable {
		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
		slog.Warn("trace telemetry is disabled, falling back to stdout", "service_name", cfg.ServiceName)
	} else {

		var secureOption otlptracegrpc.Option

		if !cfg.Otel.CollectorInsecure {
			secureOption = otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
		} else {
			secureOption = otlptracegrpc.WithInsecure()
		}

		exporter, err = otlptrace.New(
			context.Background(),
			otlptracegrpc.NewClient(
				secureOption,
				otlptracegrpc.WithEndpoint(cfg.Otel.CollectorEndpoint),
				otlptracegrpc.WithDialOption(grpc.WithBlock()),
			),
		)
	}
	if err != nil {
		panic(errors.Wrap(err, "failed to create trace exporter"))
	}

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", cfg.ServiceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		log.Fatalf("Could not set resources: %v", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return exporter.Shutdown
}

func initMetrics() func(context.Context) error {

	cfg := plconfig.Get()

	var exporter metric.Exporter
	var err error
	if !cfg.Otel.Enable {
		exporter, err = stdoutmetric.New(stdoutmetric.WithPrettyPrint())
		slog.Warn("metric telemetry is disabled, falling back to stdout", "service_name", cfg.ServiceName)
	} else {

		var secureOption otlpmetricgrpc.Option

		if !cfg.Otel.CollectorInsecure {
			secureOption = otlpmetricgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
		} else {
			secureOption = otlpmetricgrpc.WithInsecure()
		}

		exporter, err = otlpmetricgrpc.New(
			context.Background(),
			otlpmetricgrpc.WithEndpoint(cfg.Otel.CollectorEndpoint),
			secureOption,
		)
	}

	if err != nil {
		panic(errors.Wrap(err, "failed to create metric exporter"))
	}

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", cfg.ServiceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		panic(errors.Wrap(err, "failed to create OTLP metric resource"))
	}

	otel.SetMeterProvider(
		metric.NewMeterProvider(
			metric.WithResource(resources),
		),
	)

	return exporter.Shutdown
}
