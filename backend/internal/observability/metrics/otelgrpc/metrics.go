package otelgrpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/metrics"

	"github.com/hashicorp/go-multierror"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

const (
	servicePrefix = "dinner-done-better-api."
)

type Config struct {
	BaseName           string        `json:"baseName"                 toml:"base_name"`
	CollectorEndpoint  string        `json:"metricsCollectorEndpoint" toml:"metrics_collector_endpoint"`
	CollectionInterval time.Duration `json:"collectionInterval"       toml:"collection_interval"`
	Insecure           bool          `json:"insecure"                 toml:"insecure"`
	CollectionTimeout  time.Duration `json:"collectionTimeout"        toml:"collection_timeout"`
}

func ProvideMetricsProvider(ctx context.Context, logger logging.Logger, cfg *Config) (metrics.Provider, error) {
	if cfg == nil {
		return nil, errors.New("nil config")
	}

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			attribute.KeyValue{
				Key:   "service_name",
				Value: attribute.StringValue(servicePrefix + cfg.BaseName),
			},
		),
	)

	// Set up meter provider.
	options := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(cfg.CollectorEndpoint),
	}

	if cfg.Insecure {
		options = append(options, otlpmetricgrpc.WithInsecure())
	}

	exporter, err := otlpmetricgrpc.New(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create metric provider: %w", err)
	}

	logger.WithValue("interval", cfg.CollectionInterval.String()).Info("setting up period metric reader")

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(
				exporter,
				sdkmetric.WithInterval(cfg.CollectionInterval),
				sdkmetric.WithTimeout(time.Second),
			),
		),
		sdkmetric.WithResource(res),
	)

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
	return m.mp.Float64Counter(fmt.Sprintf("%s.%s", servicePrefix, name), options...)
}

func (m *providerImpl) NewFloat64Gauge(name string, options ...metric.Float64GaugeOption) (metric.Float64Gauge, error) {
	return m.mp.Float64Gauge(fmt.Sprintf("%s.%s", servicePrefix, name), options...)
}

func (m *providerImpl) NewFloat64UpDownCounter(name string, options ...metric.Float64UpDownCounterOption) (metric.Float64UpDownCounter, error) {
	return m.mp.Float64UpDownCounter(fmt.Sprintf("%s.%s", servicePrefix, name), options...)
}

func (m *providerImpl) NewFloat64Histogram(name string, options ...metric.Float64HistogramOption) (metric.Float64Histogram, error) {
	return m.mp.Float64Histogram(fmt.Sprintf("%s.%s", servicePrefix, name), options...)
}

func (m *providerImpl) NewInt64Counter(name string, options ...metric.Int64CounterOption) (metric.Int64Counter, error) {
	return m.mp.Int64Counter(fmt.Sprintf("%s.%s", servicePrefix, name), options...)
}

func (m *providerImpl) NewInt64Gauge(name string, options ...metric.Int64GaugeOption) (metric.Int64Gauge, error) {
	return m.mp.Int64Gauge(fmt.Sprintf("%s.%s", servicePrefix, name), options...)
}

func (m *providerImpl) NewInt64UpDownCounter(name string, options ...metric.Int64UpDownCounterOption) (metric.Int64UpDownCounter, error) {
	return m.mp.Int64UpDownCounter(fmt.Sprintf("%s.%s", servicePrefix, name), options...)
}

func (m *providerImpl) NewInt64Histogram(name string, options ...metric.Int64HistogramOption) (metric.Int64Histogram, error) {
	return m.mp.Int64Histogram(fmt.Sprintf("%s.%s", servicePrefix, name), options...)
}
