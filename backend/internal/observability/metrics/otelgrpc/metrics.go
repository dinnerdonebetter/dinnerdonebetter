package otelgrpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/sdk/metric/exemplar"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/metrics"

	"github.com/hashicorp/go-multierror"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

const (
	servicePrefix = "dinner-done-better-api."
)

var (
	ErrNilConfig = errors.New("nil config")
)

type Config struct {
	BaseName           string        `json:"baseName"                 toml:"base_name"`
	CollectorEndpoint  string        `json:"metricsCollectorEndpoint" toml:"metrics_collector_endpoint"`
	CollectionInterval time.Duration `json:"collectionInterval"       toml:"collection_interval"`
	Insecure           bool          `json:"insecure"                 toml:"insecure"`
	CollectionTimeout  time.Duration `json:"collectionTimeout"        toml:"collection_timeout"`
}

func setupMetricsProvider(ctx context.Context, logger logging.Logger, cfg *Config) (metric.MeterProvider, func(context.Context) error, error) {
	if cfg == nil {
		return nil, nil, ErrNilConfig
	}

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithOSType(),
		resource.WithAttributes(
			attribute.KeyValue{
				Key:   "service.name",
				Value: attribute.StringValue("demo-server"),
			},
		),
	)
	if err != nil {
		return nil, nil, err
	}

	options := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(cfg.CollectorEndpoint),
	}

	if cfg.Insecure {
		options = append(options, otlpmetricgrpc.WithInsecure())
	}

	metricExp, err := otlpmetricgrpc.New(
		ctx,
		options...,
	)
	if err != nil {
		return nil, nil, err
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithExemplarFilter(exemplar.AlwaysOnFilter),
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(
				metricExp,
				sdkmetric.WithInterval(cfg.CollectionInterval),
			),
		),
	)
	otel.SetMeterProvider(meterProvider)

	return meterProvider, meterProvider.Shutdown, nil
}

func ProvideMetricsProvider(ctx context.Context, logger logging.Logger, cfg *Config) (metrics.Provider, error) {
	if cfg == nil {
		return nil, ErrNilConfig
	}

	logger.WithValue("interval", cfg.CollectionInterval.String()).Info("setting up period metric reader")

	meterProvider, shutdown, err := setupMetricsProvider(ctx, logger, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create metric provider: %w", err)
	}

	// Set the global meter provider
	otel.SetMeterProvider(meterProvider)

	i := &providerImpl{
		mp: meterProvider.Meter(cfg.BaseName),
		shutdownFunctions: []func(context.Context) error{
			shutdown,
		},
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
