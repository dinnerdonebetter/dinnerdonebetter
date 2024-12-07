package otelgrpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dinnerdonebetter/backend/internal/observability/metrics"

	"github.com/hashicorp/go-multierror"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

type Config struct {
	BaseName           string        `json:"baseName"                  toml:"base_name"`
	CollectorEndpoint  string        `json:"metricsCollectorEndpoint"  toml:"metrics_collector_endpoint"`
	CollectionInterval time.Duration `json:"metricsCollectionInterval" toml:"metrics_collection_interval"`
}

func ProvideMetricsProvider(ctx context.Context, cfg *Config) (metrics.Provider, error) {
	if cfg == nil {
		return nil, errors.New("nil config")
	}

	// Set up propagator.
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Set up meter provider.
	options := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(cfg.CollectorEndpoint),
		otlpmetricgrpc.WithInsecure(),
	}

	exporter, err := otlpmetricgrpc.New(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create metric provider: %w", err)
	}

	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(sdkmetric.NewPeriodicReader(
		exporter,
		sdkmetric.WithInterval(cfg.CollectionInterval),
	)))
	// TODO: handle shutdown func somehow

	// Set the global meter provider
	otel.SetMeterProvider(provider)

	i := &providerImpl{
		mp:                provider.Meter(cfg.BaseName),
		shutdownFunctions: []func(context.Context) error{provider.Shutdown},
	}

	return i, nil
}

var _ metrics.Provider = (*providerImpl)(nil)

type providerImpl struct {
	mp                metric.Meter
	shutdownFunctions []func(context.Context) error
}

func (m *providerImpl) Shutdown(ctx context.Context) error {
	multierr := &multierror.Error{}

	for _, fn := range m.shutdownFunctions {
		if err := fn(ctx); err != nil {
			multierr = multierror.Append(multierr, err)
		}
	}

	return multierr.ErrorOrNil()
}

func (m *providerImpl) NewFloat64Counter(name string, options ...metric.Float64CounterOption) (metric.Float64Counter, error) {
	return m.mp.Float64Counter(name, options...)
}

func (m *providerImpl) NewFloat64Gauge(name string, options ...metric.Float64GaugeOption) (metric.Float64Gauge, error) {
	return m.mp.Float64Gauge(name, options...)
}

func (m *providerImpl) NewFloat64UpDownCounter(name string, options ...metric.Float64UpDownCounterOption) (metric.Float64UpDownCounter, error) {
	return m.mp.Float64UpDownCounter(name, options...)
}

func (m *providerImpl) NewFloat64Histogram(name string, options ...metric.Float64HistogramOption) (metric.Float64Histogram, error) {
	return m.mp.Float64Histogram(name, options...)
}

func (m *providerImpl) NewInt64Counter(name string, options ...metric.Int64CounterOption) (metric.Int64Counter, error) {
	return m.mp.Int64Counter(name, options...)
}

func (m *providerImpl) NewInt64Gauge(name string, options ...metric.Int64GaugeOption) (metric.Int64Gauge, error) {
	return m.mp.Int64Gauge(name, options...)
}

func (m *providerImpl) NewInt64UpDownCounter(name string, options ...metric.Int64UpDownCounterOption) (metric.Int64UpDownCounter, error) {
	return m.mp.Int64UpDownCounter(name, options...)
}

func (m *providerImpl) NewInt64Histogram(name string, options ...metric.Int64HistogramOption) (metric.Int64Histogram, error) {
	return m.mp.Int64Histogram(name, options...)
}
