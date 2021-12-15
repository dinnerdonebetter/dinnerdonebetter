package cloudwatch

import (
	"context"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/unit"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"google.golang.org/grpc"

	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/metrics"
)

const (
	// defaultNamespace is the default namespace under which we register metrics.
	defaultNamespace                 = "prixfixe_server"
	instrumentationVersion           = "1.0.0"
	minimumMetricsCollectionInterval = time.Second
	minimumRuntimeCollectionInterval = time.Second
)

type (
	// Config contains settings related to Prometheus.
	Config struct {
		_ struct{}

		CollectorEndpoint                string        `json:"collector_endpoint,omitempty" mapstructure:"collector_endpoint" toml:"collector_endpoint,omitempty"`
		MetricsCollectionInterval        time.Duration `json:"runtimeMetricsCollectionInterval,omitempty" mapstructure:"runtime_metrics_collection_interval" toml:"runtime_metrics_collection_interval,omitempty"`
		RuntimeMetricsCollectionInterval time.Duration `json:"metricsCollectionInterval,omitempty" mapstructure:"metrics_collection_interval" toml:"metrics_collection_interval,omitempty"`
	}
)

// ValidateWithContext validates the config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, cfg,
		validation.Field(&cfg.CollectorEndpoint, validation.Required),
		validation.Field(&cfg.MetricsCollectionInterval, validation.Min(minimumRuntimeCollectionInterval)),
		validation.Field(&cfg.RuntimeMetricsCollectionInterval, validation.Min(minimumRuntimeCollectionInterval)),
	)
}

func initiateExporter(ctx context.Context, cfg *Config) (metric.MeterProvider, error) {
	metricExporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(cfg.CollectorEndpoint),
		otlpmetricgrpc.WithDialOption(grpc.WithBlock()),
	)
	if err != nil {
		return nil, err
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("prixfixe_api"),
	)

	cont := controller.New(
		processor.NewFactory(
			simple.NewWithHistogramDistribution(),
			metricExporter,
		),
		controller.WithExporter(metricExporter),
		controller.WithCollectPeriod(cfg.MetricsCollectionInterval),
		controller.WithResource(res),
	)

	global.SetMeterProvider(cont)

	if metricStartErr := cont.Start(ctx); metricStartErr != nil {
		return nil, metricStartErr
	}

	return cont, nil
}

// ProvideInstrumentationHandler provides an instrumentation handler.
func (cfg *Config) ProvideInstrumentationHandler(ctx context.Context, logger logging.Logger) (metric.MeterProvider, error) {
	logger.Debug("setting metrics provider")

	if err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(cfg.RuntimeMetricsCollectionInterval)); err != nil {
		return nil, fmt.Errorf("failed to start runtime metrics collection: %w", err)
	}

	return initiateExporter(ctx, cfg)
}

// ProvideUnitCounterProvider provides an instrumentation handler.
func (cfg *Config) ProvideUnitCounterProvider(ctx context.Context, logger logging.Logger) (metrics.UnitCounterProvider, error) {
	logger.Debug("setting up meter")

	meterProvider, err := initiateExporter(ctx, cfg)
	if err != nil {
		return nil, err
	}

	mustMeter := metric.Must(meterProvider.Meter(defaultNamespace, metric.WithInstrumentationVersion(instrumentationVersion)))
	logger.Debug("meter initialized successfully")

	uc := func(name, description string) metrics.UnitCounter {
		l := logger.WithValue("name", name)

		counter := mustMeter.NewInt64Counter(
			fmt.Sprintf("%s_count", name),
			metric.WithUnit(unit.Dimensionless),
			metric.WithDescription(description),
		)

		l.Debug("returning wrapped unit counter")

		return metrics.NewUnitCounter(counter)
	}

	return uc, nil
}
