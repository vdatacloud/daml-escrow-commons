// Package telemetry sets up OpenTelemetry tracing and metrics for a
// platform service and traces its HTTP traffic in and out (daml-escrow
// PLAN.md Phase 96): OTLP over HTTP to a collector, W3C trace-context
// propagation so one trace follows a request across services, and spans
// whose attributes keep to the logging rules -- the route pattern and a
// redacted path, never the query string, a body or a header.
package telemetry

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Config is a service's telemetry configuration; the caller loads it.
type Config struct {
	// Enabled exports traces and metrics; off, the global providers stay
	// no-ops (trace context still propagates).
	Enabled bool
	// Endpoint is the collector's OTLP/HTTP address: "host:port", or a URL
	// ("https://collector:4318"). "" leaves it to the OTEL_EXPORTER_OTLP_*
	// environment variables, else localhost:4318.
	Endpoint string
	// Insecure sends without TLS (a collector on the same host or in the
	// same cluster network only).
	Insecure bool
	// SampleRatio is the share of new traces kept, in (0, 1]; 0 or less
	// keeps all. A trace started upstream keeps its parent's decision.
	// Metrics are never sampled.
	SampleRatio float64
	// MetricInterval is how often metrics are exported; 0 is 15s.
	MetricInterval time.Duration
	Service        string
	Environment    string
	Version        string
}

// Propagator is the W3C trace-context and baggage propagator Setup installs.
func Propagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
}

// Setup installs the propagator and, when cfg.Enabled, the global tracer
// and meter providers exporting to the collector. The returned shutdown
// flushes and stops them; call it before exit.
func Setup(ctx context.Context, cfg Config) (shutdown func(context.Context) error, err error) {
	otel.SetTextMapPropagator(Propagator())
	if !cfg.Enabled {
		return func(context.Context) error { return nil }, nil
	}
	// Schemaless: merging a versioned schema onto resource.Default fails
	// whenever the two semconv versions differ (escrow-api's own setup did,
	// silently, for months -- PLAN.md Phase 96).
	res, err := resource.Merge(resource.Default(), resource.NewSchemaless(
		attribute.String("service.name", cfg.Service),
		attribute.String("service.version", cfg.Version),
		attribute.String("deployment.environment", cfg.Environment),
	))
	if err != nil {
		return nil, fmt.Errorf("telemetry: resource: %w", err)
	}

	traceExporter, err := otlptracehttp.New(ctx, traceOptions(cfg)...)
	if err != nil {
		return nil, fmt.Errorf("telemetry: trace exporter: %w", err)
	}
	ratio := cfg.SampleRatio
	if ratio <= 0 || ratio > 1 {
		ratio = 1
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(ratio))),
	)

	metricExporter, err := otlpmetrichttp.New(ctx, metricOptions(cfg)...)
	if err != nil {
		_ = tp.Shutdown(ctx)
		return nil, fmt.Errorf("telemetry: metric exporter: %w", err)
	}
	interval := cfg.MetricInterval
	if interval <= 0 {
		interval = 15 * time.Second
	}
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter, sdkmetric.WithInterval(interval))),
		sdkmetric.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetMeterProvider(mp)
	return func(ctx context.Context) error {
		return errors.Join(tp.Shutdown(ctx), mp.Shutdown(ctx))
	}, nil
}

func isURL(endpoint string) bool {
	return strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://")
}

func traceOptions(cfg Config) []otlptracehttp.Option {
	var o []otlptracehttp.Option
	switch {
	case isURL(cfg.Endpoint):
		o = append(o, otlptracehttp.WithEndpointURL(strings.TrimSuffix(cfg.Endpoint, "/")+"/v1/traces"))
	case cfg.Endpoint != "":
		o = append(o, otlptracehttp.WithEndpoint(cfg.Endpoint))
	}
	if cfg.Insecure {
		o = append(o, otlptracehttp.WithInsecure())
	}
	return o
}

func metricOptions(cfg Config) []otlpmetrichttp.Option {
	var o []otlpmetrichttp.Option
	switch {
	case isURL(cfg.Endpoint):
		o = append(o, otlpmetrichttp.WithEndpointURL(strings.TrimSuffix(cfg.Endpoint, "/")+"/v1/metrics"))
	case cfg.Endpoint != "":
		o = append(o, otlpmetrichttp.WithEndpoint(cfg.Endpoint))
	}
	if cfg.Insecure {
		o = append(o, otlpmetrichttp.WithInsecure())
	}
	return o
}
